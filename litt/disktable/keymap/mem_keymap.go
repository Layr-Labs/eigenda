package keymap

import (
	"sync"

	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ KeyMap = &memKeyMap{}

// An in-memory keymap implementation. When a table using a memKeyMap is restarted, it loads all keys from
// the segment files.
//
// - potentially high memory usage for large keymaps
// - potentially slow startup time for large keymaps
// - very fast after startup
type memKeyMap struct {
	logger logging.Logger
	data   map[string]types.Address
	lock   sync.RWMutex
}

// NewMemKeyMap creates a new in-memory keymap.
func NewMemKeyMap(logger logging.Logger) KeyMap {
	return &memKeyMap{
		logger: logger,
		data:   make(map[string]types.Address),
	}
}

func (m *memKeyMap) Put(pairs []*types.KAPair) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, pair := range pairs {
		m.data[string(pair.Key)] = pair.Address
	}
	return nil
}

func (m *memKeyMap) Get(key []byte) (types.Address, bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	address, ok := m.data[string(key)]
	return address, ok, nil
}

func (m *memKeyMap) Delete(keys []*types.KAPair) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, key := range keys {
		delete(m.data, string(key.Key))
	}
	return nil
}

func (m *memKeyMap) Stop() error {
	// nothing to do here
	return nil
}

func (m *memKeyMap) Destroy() error {
	// nothing to do here
	return nil
}
