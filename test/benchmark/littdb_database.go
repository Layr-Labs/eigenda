package benchmark

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
)

var _ Database = (*LittDBDatabase)(nil)

// LittDBDatabase is a thin wrapper around LittDB that enables it to be used for benchmarking.
type LittDBDatabase struct {
	db    litt.DB
	table litt.Table
}

func NewLittDBDatabase(config *litt.Config) (*LittDBDatabase, error) {
	db, err := littbuilder.NewDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to build DB: %w", err)
	}

	table, err := db.GetTable("benchmark")
	if err != nil {
		return nil, fmt.Errorf("failed to get table: %w", err)
	}

	return &LittDBDatabase{
		db:    db,
		table: table,
	}, nil
}

func (l *LittDBDatabase) Write(key []byte, value []byte) error {
	return l.table.Put(key, value)
}

func (l *LittDBDatabase) Flush() error {
	return l.table.Flush()
}

func (l *LittDBDatabase) Read(key []byte) (value []byte, exists bool, err error) {
	return l.table.Get(key)
}

func (l *LittDBDatabase) Close() error {
	return l.db.Close()
}
