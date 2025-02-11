package mem

import (
	"github.com/Layr-Labs/eigenda/litt"
	"sync"
	"time"
)

var _ litt.Table = &memTable{}

// memTable is a simple implementation of a Table that stores its data in memory.
type memTable struct {
	name string
	data sync.Map

	// TODO: queue for TTL
}

// NewMemTable creates a new in-memory table.
func newMemTable(name string) litt.Table {
	return &memTable{}
}

func (m *memTable) Name() string {
	//TODO implement me
	panic("implement me")
}

func (m *memTable) Put(key []byte, value []byte) error {
	//TODO implement me
	panic("implement me")
}

func (m *memTable) Get(key []byte) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (m *memTable) Flush() error {
	// This is a no-op for a memory table. Memory tables are ephemeral by nature.
	return nil
}

func (m *memTable) SetTTL(ttl time.Duration) {
	//TODO implement me
	panic("implement me")
}
