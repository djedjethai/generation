package logger

import (
	"database/sql"
	"fmt"
	cfg "github.com/djedjethai/generation0/pkg/config"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"time"
)

const (
	maxOpenDbConn = 25
	maxIdleDbConn = 25
	maxDbLifetime = 5 * time.Minute
)

type PostgresTransactionLogger struct {
	events chan<- Event // Write-only channel for sending events
	errors <-chan error // Read-only channel for receiving errors
	db     *sql.DB      // The database access interface
	active bool
}

func NewPostgresTransactionLogger(config cfg.PostgresDBParams, active bool) (TransactionLogger,
	error) {

	if active {
		dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s",
			config.Host, config.DbName, config.User, config.Password)

		db, err := sql.Open("pgx", dsn)
		if err != nil {
			panic(err)
		}
		db.SetMaxOpenConns(maxOpenDbConn)
		db.SetMaxIdleConns(maxIdleDbConn)
		db.SetConnMaxLifetime(maxDbLifetime)

		err = db.Ping()
		if err != nil {
			log.Println("Error! ", err)
		} else {
			log.Println("**** Pinged postgres successfuly ****")
		}

		logger := &PostgresTransactionLogger{
			db:     db,
			active: active,
		}
		exists, err := logger.verifyTableExists()
		if err != nil {
			return nil, fmt.Errorf("failed to verify table exists: %w", err)
		}
		if !exists {
			if err = logger.createTable(); err != nil {
				return nil, fmt.Errorf("failed to create table: %w", err)
			}
		}
		return logger, nil
	} else {

		logger := &PostgresTransactionLogger{active: active}
		return logger, nil
	}
}

func (l *PostgresTransactionLogger) CloseFileLogger() {
	if l.active {
		if err := l.db.Close(); err != nil {
			log.Println("error closing the fileLogger")
		}
	}
}

func (l *PostgresTransactionLogger) WritePut(key, value string) {
	if l.active {
		l.events <- Event{EventType: EventPut, Key: key, Value: value}
	}
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	if l.active {
		l.events <- Event{EventType: EventDelete, Key: key}
	}
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan Event, 16)
	l.events = events

	errors := make(chan error, 1)
	l.errors = errors

	go func() {
		query := `INSERT INTO transactions
			(event_type, key, value)
			VALUES($1, $2, $3)`

		for e := range events {

			_, err := l.db.Exec(query, e.EventType, e.Key, e.Value)
			if err != nil {
				errors <- err
			}
		}
	}()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		defer close(outEvent)
		defer close(outError)

		query := `SELECT sequence, event_type, key, value FROM transactions ORDER BY sequence`

		rows, err := l.db.Query(query)
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}
		defer rows.Close()

		e := Event{}

		for rows.Next() {
			err = rows.Scan(&e.Sequence, &e.EventType, &e.Key, &e.Value)
			if err != nil {
				outError <- fmt.Errorf("error reading row: %w", err)
			}

			outEvent <- e
		}
		if err = rows.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	const table = "transactions"

	var result string

	rows, err := l.db.Query(fmt.Sprintf("SELECT to_regclass('public.%s');", table))
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() && result != table {
		rows.Scan(&result)
	}

	return result == table, rows.Err()
}

func (l *PostgresTransactionLogger) createTable() error {
	var err error

	createQuery := `CREATE TABLE transactions (
		sequence      BIGSERIAL PRIMARY KEY,
		event_type    SMALLINT,
		key 		  TEXT,
		value         TEXT
	  );`

	_, err = l.db.Exec(createQuery)
	if err != nil {
		return err
	}

	return nil
}
