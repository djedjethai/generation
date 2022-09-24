package logger

type EventType byte

const (
	_                     = iota
	EventDelete EventType = iota
	EventPut
)

type TransactionLogger interface {
	CloseFileLogger()
	WriteDelete(key string)
	WritePut(key, value string)
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
