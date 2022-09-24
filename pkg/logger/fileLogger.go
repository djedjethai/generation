package logger

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/djedjethai/generation0/pkg/internal"
)

type FileTransactionLogger struct {
	events       chan<- Event // Write-only channel for sending events
	errors       <-chan error // Read-only channel for receiving errors
	lastSequence uint64       // The last used event sequence number
	file         *os.File     // The location of the transaction log
	encryptK     string
	active       bool
}

func NewFileTransactionLogger(filename, encryptK string, active bool) (TransactionLogger, error) {
	if active {
		file, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			return nil, fmt.Errorf("cannot open transaction log file: %w", err)
		}
		return &FileTransactionLogger{
			file:     file,
			encryptK: encryptK,
			active:   active,
		}, nil
	} else {
		return &FileTransactionLogger{
			active: active,
		}, nil
	}
}

func (l *FileTransactionLogger) WritePut(key, value string) {
	if l.active {
		key, err := internal.Encrypt(key, l.encryptK)
		value, err = internal.Encrypt(value, l.encryptK)
		if err != nil {
			log.Println("Error encrypting the key, value wheb WritePut")
		}
		l.events <- Event{EventType: EventPut, Key: key, Value: value}
	}
}

func (l *FileTransactionLogger) WriteDelete(key string) {
	if l.active {
		key, err := internal.Encrypt(key, l.encryptK)
		value, err := internal.Encrypt("", l.encryptK)
		if err != nil {
			log.Println("Error encrypting the key when WriteDelete")
		}
		l.events <- Event{EventType: EventDelete, Key: key, Value: value}
	}
}

func (l *FileTransactionLogger) CloseFileLogger() {
	if l.active {
		if err := l.file.Close(); err != nil {
			log.Println("error closing the fileLogger")
		}
	}
}

func (l *FileTransactionLogger) Run() {
	events := make(chan Event, 16) // Make an events channel
	l.events = events

	errors := make(chan error, 1) // Make an errors channel
	l.errors = errors

	go func() {

		for e := range events { // Retrieve the next Event

			l.lastSequence++ // Increment sequence number

			_, err := fmt.Fprintf( // Write the event to the log
				l.file,
				"%d\t%d\t%s\t%s\n",
				l.lastSequence, e.EventType, e.Key, e.Value)

			if err != nil {
				errors <- err
				return
			}
		}

	}()
}

func (l *FileTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *FileTransactionLogger) ReadEvents() (<-chan Event, <-chan error) {
	scanner := bufio.NewScanner(l.file)
	outEvent := make(chan Event)
	outError := make(chan error, 1)

	go func() {
		var e Event

		defer close(outEvent)
		defer close(outError)

		for scanner.Scan() {
			line := scanner.Text()

			if _, err := fmt.Sscanf(line, "%d\t%d\t%s\t%s",
				&e.Sequence, &e.EventType, &e.Key, &e.Value); err != nil {
				outError <- fmt.Errorf("input parse error: %w", err)
				return
			}

			if l.lastSequence >= e.Sequence {
				outError <- fmt.Errorf("transaction numbers out of sequence")
				return
			}

			l.lastSequence = e.Sequence

			var err error
			e.Key, err = internal.Decrypt(e.Key, l.encryptK)
			e.Value, err = internal.Decrypt(e.Value, l.encryptK)
			if err != nil {
				outError <- fmt.Errorf("err when decrypting the key, value")
			}

			outEvent <- e
		}

		if err := scanner.Err(); err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}
