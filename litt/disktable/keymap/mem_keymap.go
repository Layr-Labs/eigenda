package keymap

import (
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ Keymap = &memKeymap{}

// An in-memory keymap implementation. When a table using a memKeymap is restarted, it loads all keys from
// the segment files.  Methods on this struct are goroutine safe.
//
// - potentially high memory usage for large keymaps
// - potentially slow startup time for large keymaps
// - very fast after startup
type memKeymap struct {
	logger logging.Logger
	data   map[string]types.Address
	// if true, then return an error if an update would overwrite an existing key
	doubleWriteProtection bool
	lock                  sync.RWMutex
}

// NewMemKeymap creates a new in-memory keymap.
func NewMemKeymap(logger logging.Logger,
	keymapPath string,
	doubleWriteProtection bool) (kmap Keymap, requiresReload bool, err error) {

	return &memKeymap{
		logger:                logger,
		data:                  make(map[string]types.Address),
		doubleWriteProtection: doubleWriteProtection,
	}, true, nil
}

func (m *memKeymap) Put(pairs []*types.KAPair) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, pair := range pairs {
		stringKey := util.UnsafeBytesToString(pair.Key)

		if m.doubleWriteProtection {
			_, ok := m.data[stringKey]
			if ok {
				return fmt.Errorf("key %s already exists", pair.Key)
			}
		}

		m.data[stringKey] = pair.Address
	}
	return nil
}

func (m *memKeymap) Get(key []byte) (types.Address, bool, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	address, ok := m.data[util.UnsafeBytesToString(key)]
	return address, ok, nil
}

func (m *memKeymap) Delete(keys []*types.KAPair) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, key := range keys {
		delete(m.data, util.UnsafeBytesToString(key.Key))
	}

	return nil
}

func (m *memKeymap) Stop() error {
	// nothing to do here
	return nil
}

func (m *memKeymap) Destroy() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.data = nil
	return nil
}
