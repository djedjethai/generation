package logger

import (
	"fmt"
	"github.com/djedjethai/generation/internal/config"
	"golang.org/x/net/context"
	"log"
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
	services       *config.Services
	dbLoggerActive bool
	postgresConfig config.PostgresDBParams
}

type LoggerFacade struct {
	// fileLogger   TransactionLogger
	dbLogger   TransactionLogger
	isDBRecord bool
}

func NewLoggerFacade(srv *config.Services, dbLoggerActive bool, postgresConfig config.PostgresDBParams) (*LoggerFacade, error) {

	dbLogger, err := NewTransactionLoggerFactory(srv, dbLoggerActive, postgresConfig).Start()

	fmt.Println("in logger NewLoggerFacade")

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

func NewTransactionLoggerFactory(srv *config.Services, dbLoggerActive bool, postgresConfig config.PostgresDBParams) *TransactionLoggerFactory {
	return &TransactionLoggerFactory{
		services:       srv,
		dbLoggerActive: dbLoggerActive,
		postgresConfig: postgresConfig,
	}
}

func (tlf *TransactionLoggerFactory) Start() (TransactionLogger, error) {
	var err error
	var dbLogger TransactionLogger

	fmt.Println("see postgres config: ", tlf.postgresConfig)

	if tlf.dbLoggerActive {
		dbLogger, err = NewPostgresTransactionLogger(tlf.postgresConfig)
		if err != nil {
			log.Println("Err when initialize PostgresTransactionLogger", err)
		}

		err = tlf.runner(dbLogger)
		if err != nil {
			log.Println("Err when run PostgresTransactionLogger", err)
		}
		fmt.Println("tame mmememe la puuute okkk")
	}

	return dbLogger, err
}

func (tlf *TransactionLoggerFactory) runner(logger TransactionLogger) error {
	// TODO here add ctx
	ctx := context.Background()

	fmt.Println("tame mmememe")
	var err error
	events, errors := logger.ReadEvents()
	e, ok := Event{}, true

	fmt.Println("tame mmememe2")
	for ok {
		select {
		case err, ok = <-errors:

			fmt.Println("tame mmememe6", err)
		case e, ok = <-events:
			switch e.EventType {
			case EventDelete:
				fmt.Println("tame mmememe4")
				err = tlf.services.Deleter.Delete(ctx, e.Key)
				fmt.Println("tame mmememe44")
			case EventPut:
				fmt.Println("tame mmememe5")
				err = tlf.services.Setter.Set(ctx, e.Key, []byte(e.Value))
				fmt.Println("tame mmememe55")
			}

		}
	}

	fmt.Println("tame mmememe3")

	logger.Run()

	return err
}
