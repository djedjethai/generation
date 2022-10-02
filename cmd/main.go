package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/djedjethai/generation0/pkg/config"
	"github.com/djedjethai/generation0/pkg/deleter"
	"github.com/djedjethai/generation0/pkg/getter"
	"github.com/djedjethai/generation0/pkg/handlers/rest"
	lgr "github.com/djedjethai/generation0/pkg/logger"
	"github.com/djedjethai/generation0/pkg/setter"
	storage "github.com/djedjethai/generation0/pkg/storage"
)

// var store = make(map[string]*storage.Node)
var encryptK = "PX9PHFrdn79ljrjLDZHlV1t+BdxHRFf5"
var postgresConfig = config.PostgresDBParams{
	Host:     "localhost",
	DbName:   "transactions",
	User:     "postgres",
	Password: "password",
}

var port = ":8080"

// default value
var fileLoggerActive = false
var dbLoggerActive = false

func main() {

	// storage(infra layer)
	// the first arg is the number of shard, the second the number of item/shard
	shardedMap := storage.NewShardedMap(3, 2)

	setSrv := setter.NewSetter(shardedMap)
	getSrv := getter.NewGetter(shardedMap)
	delSrv := deleter.NewDeleter(shardedMap)

	// in case the srv crash, when start back it will read the logger and recover its state
	logger, err := initializeTransactionLog(setSrv, delSrv, fileLoggerActive)
	if err != nil {
		log.Panic("FileLogger initialization failed")
	}
	defer logger.CloseFileLogger()

	dbLogger, err := initializeTransactionLogDB(setSrv, delSrv, dbLoggerActive)
	if err != nil {
		log.Panic("dbLogger initialization failed")
	}
	defer logger.CloseFileLogger()

	// handler(application layer)
	router := rest.Handler(setSrv, getSrv, delSrv, logger, dbLogger)

	fmt.Printf("***** Service listening on port %s *****", port)
	log.Fatal(http.ListenAndServe(port, router))
}

func initializeTransactionLogDB(setSrv setter.Setter, delSrv deleter.Deleter, active bool) (lgr.TransactionLogger, error) {
	var err error

	dbLogger, err := lgr.NewPostgresTransactionLogger(postgresConfig, active)
	if err != nil {
		return nil, fmt.Errorf("failed to create db event logger: %w", err)
	}

	if active {
		events, errors := dbLogger.ReadEvents()
		e, ok := lgr.Event{}, true

		for ok && err == nil {
			select {
			case err, ok = <-errors:
			case e, ok = <-events:
				switch e.EventType {
				case lgr.EventDelete:
					err = delSrv.Delete(e.Key)
				case lgr.EventPut:
					err = setSrv.Set(e.Key, []byte(e.Value))
				}

			}
		}

		dbLogger.Run()
	}

	return dbLogger, err

}

func initializeTransactionLog(setSrv setter.Setter, delSrv deleter.Deleter, active bool) (lgr.TransactionLogger, error) {
	var err error

	fileLogger, err := lgr.NewFileTransactionLogger("transaction.log", encryptK, active)
	if err != nil {
		return nil, fmt.Errorf("failed to create event logger: %w", err)
	}

	if active {
		events, errors := fileLogger.ReadEvents()
		e, ok := lgr.Event{}, true

		for ok && err == nil {
			select {
			case err, ok = <-errors: // retrieve any error
			case e, ok = <-events:
				switch e.EventType {
				case lgr.EventDelete:
					err = delSrv.Delete(e.Key)
				case lgr.EventPut:
					err = setSrv.Set(e.Key, []byte(e.Value))
				}
			}
		}

		fileLogger.Run()
	}

	return fileLogger, err
}
