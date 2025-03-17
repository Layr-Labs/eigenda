package keymap

import (
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ Keymap = &memKeymap{}

// An in-memory keymap implementation. When a table using a memKeymap is restarted, it loads all keys from
// the segment files.
//
// - potentially high memory usage for large keymaps
// - potentially slow startup time for large keymaps
// - very fast after startup
type memKeymap struct {
	logger logging.Logger
	data   map[string]types.Address
	lock   sync.RWMutex
}

// NewMemKeymap creates a new in-memory keymap.
func NewMemKeymap(logger logging.Logger) Keymap {

	return &memKeymap{
		logger: logger,
		data:   make(map[string]types.Address),
	}
}

func (m *memKeymap) Put(pairs []*types.KAPair) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, pair := range pairs {

		// TODO make this check optional!!
		// TODO: also add a similar but optional check to the LevelDBKeymap
		_, ok := m.data[string(pair.Key)]
		if ok {
			return fmt.Errorf("key %s already exists", pair.Key)
		}

		m.data[string(pair.Key)] = pair.Address
	}
	return nil
}

func (m *memKeymap) Get(key []byte) (types.Address, bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	address, ok := m.data[string(key)]
	return address, ok, nil
}

func (m *memKeymap) Delete(keys []*types.KAPair) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, key := range keys {
		delete(m.data, string(key.Key))
	}

	return nil
}

func (m *memKeymap) Stop() error {
	// nothing to do here
	return nil
}

func (m *memKeymap) Destroy() error {
	// nothing to do here
	return nil
}
