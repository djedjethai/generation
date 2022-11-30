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
	setSrv           setter.Setter
	delSrv           deleter.Deleter
	fileLoggerActive bool
	dbLoggerActive   bool
	postgresConfig   config.PostgresDBParams
	encryptK         string
}

type LoggerFacade struct {
	fileLogger   TransactionLogger
	dbLogger     TransactionLogger
	isFileRecord bool
	isDBRecord   bool
}

func NewLoggerFacade(setSrv setter.Setter, delSrv deleter.Deleter, fileLoggerActive bool, dbLoggerActive bool, postgresConfig config.PostgresDBParams, encryptK string) (*LoggerFacade, error) {

	fileLogger, dbLogger, err := NewTransactionLoggerFactory(setSrv, delSrv, fileLoggerActive, dbLoggerActive, postgresConfig, encryptK).Start()

	return &LoggerFacade{
		fileLogger:   fileLogger,
		dbLogger:     dbLogger,
		isFileRecord: fileLoggerActive,
		isDBRecord:   dbLoggerActive,
	}, err
}

func (lf *LoggerFacade) WriteSet(key, value string) {
	if lf.isFileRecord {
		lf.fileLogger.WriteSet(key, value)
	}
	if lf.isDBRecord {
		lf.dbLogger.WriteSet(key, value)
	}
}

func (lf *LoggerFacade) WriteDelete(key string) {
	if lf.isFileRecord {
		lf.fileLogger.WriteDelete(key)
	}
	if lf.isDBRecord {
		lf.dbLogger.WriteDelete(key)
	}
}

func (lf *LoggerFacade) CloseFileLogger() func() {
	if lf.isFileRecord && lf.isDBRecord {
		return func() {
			defer lf.fileLogger.CloseFileLogger()
			defer lf.dbLogger.CloseFileLogger()
		}
	}

	if !lf.isFileRecord && lf.isDBRecord {
		return func() {
			defer lf.dbLogger.CloseFileLogger()
		}
	}

	if lf.isFileRecord && !lf.isDBRecord {
		return func() {
			defer lf.fileLogger.CloseFileLogger()
		}
	}
	return nil
}

func NewTransactionLoggerFactory(setSrv setter.Setter, delSrv deleter.Deleter, fileLoggerActive, dbLoggerActive bool, postgresConfig config.PostgresDBParams, encryptK string) *TransactionLoggerFactory {
	return &TransactionLoggerFactory{
		setSrv:           setSrv,
		delSrv:           delSrv,
		fileLoggerActive: fileLoggerActive,
		dbLoggerActive:   dbLoggerActive,
		postgresConfig:   postgresConfig,
		encryptK:         encryptK,
	}
}

func (tlf *TransactionLoggerFactory) Start() (TransactionLogger, TransactionLogger, error) {
	var err error
	var fileLogger TransactionLogger
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

	if tlf.fileLoggerActive {
		fileLogger, err = NewFileTransactionLogger("transaction.log", tlf.encryptK)
		if err != nil {
			log.Println("Err when initialize fileTransactionLogger", err)
		}

		err = tlf.runner(fileLogger)
		if err != nil {
			log.Println("Err when run fileTransactionLogger", err)
		}
	}

	return fileLogger, dbLogger, err
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
