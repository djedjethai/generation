package models

type KeysValues struct {
	Key   string
	Value string
}

type Record struct {
	Value  []byte
	Offset uint64
	Term   uint64
	Type   uint32
}
