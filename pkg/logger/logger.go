package logger

import (
	"log"

	"github.com/djedjethai/generation/pkg/config"
	"github.com/djedjethai/generation/pkg/deleter"
	"github.com/djedjethai/generation/pkg/setter"
	"golang.org/x/net/context"
)

type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type TransactionLogger interface {
	CloseFileLogger()
	WriteDelete(key string)
	WriteSet(key, value string)
	Err() <-chan error
	Run()
	ReadEvents() (<-chan Event, <-chan error)
}

type Event struct {
	Sequence  uint64
	EventType EventType
	Key       string
	Value     string
}

type TransactionLoggerFactory struct {
	setSrv         setter.Setter
	delSrv         deleter.Deleter
	dbLoggerActive bool
	postgresConfig config.PostgresDBParams
}

type LoggerFacade struct {
	// fileLogger   TransactionLogger
	dbLogger   TransactionLogger
	isDBRecord bool
}

func NewLoggerFacade(setSrv setter.Setter, delSrv deleter.Deleter, dbLoggerActive bool, postgresConfig config.PostgresDBParams) (*LoggerFacade, error) {

	dbLogger, err := NewTransactionLoggerFactory(setSrv, delSrv, dbLoggerActive, postgresConfig).Start()

	return &LoggerFacade{
		dbLogger:   dbLogger,
		isDBRecord: dbLoggerActive,
	}, err
}

func (lf *LoggerFacade) WriteSet(key, value string) {
	if lf.isDBRecord {
		lf.dbLogger.WriteSet(key, value)
	}
}

func (lf *LoggerFacade) WriteDelete(key string) {
	if lf.isDBRecord {
		lf.dbLogger.WriteDelete(key)
	}
}

// func (lf *LoggerFacade) CloseFileLogger() func() {
// 	if lf.isFileRecord && lf.isDBRecord {
// 		return func() {
// 			defer lf.fileLogger.CloseFileLogger()
// 			defer lf.dbLogger.CloseFileLogger()
// 		}
// 	}
//
// 	if !lf.isFileRecord && lf.isDBRecord {
// 		return func() {
// 			defer lf.dbLogger.CloseFileLogger()
// 		}
// 	}
//
// 	if lf.isFileRecord && !lf.isDBRecord {
// 		return func() {
// 			defer lf.fileLogger.CloseFileLogger()
// 		}
// 	}
// 	return nil
// }

func NewTransactionLoggerFactory(setSrv setter.Setter, delSrv deleter.Deleter, dbLoggerActive bool, postgresConfig config.PostgresDBParams) *TransactionLoggerFactory {
	return &TransactionLoggerFactory{
		setSrv:         setSrv,
		delSrv:         delSrv,
		dbLoggerActive: dbLoggerActive,
		postgresConfig: postgresConfig,
	}
}

func (tlf *TransactionLoggerFactory) Start() (TransactionLogger, error) {
	var err error
	var dbLogger TransactionLogger

	if tlf.dbLoggerActive {
		dbLogger, err = NewPostgresTransactionLogger(tlf.postgresConfig)
		if err != nil {
			log.Println("Err when initialize PostgresTransactionLogger", err)
		}

		err = tlf.runner(dbLogger)
		if err != nil {
			log.Println("Err when run PostgresTransactionLogger", err)
		}
	}

	return dbLogger, err
}

func (tlf *TransactionLoggerFactory) runner(logger TransactionLogger) error {
	// TODO here add ctx
	ctx := context.Background()

	var err error
	events, errors := logger.ReadEvents()
	e, ok := Event{}, true

	for ok {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				err = tlf.delSrv.Delete(ctx, e.Key)
			case EventPut:
				err = tlf.setSrv.Set(ctx, e.Key, []byte(e.Value))
			}

		}
	}

	logger.Run()

	return err
}
