package keymap

import (
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/types"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ KeyMap = &memKeymap{}

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

// NewMemKeyMap creates a new in-memory keymap.
func NewMemKeyMap(logger logging.Logger) KeyMap {
	return &memKeymap{
		logger: logger,
		data:   make(map[string]types.Address),
	}
}

func (m *memKeymap) Put(pairs []*types.KAPair) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, pair := range pairs {
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

func (m *memKeymap) LoadFromSegments(
	segments map[uint32]*segment.Segment,
	lowestSegmentIndex uint32,
	highestSegmentIndex uint32) error {

	// It's possible that some of the data written near the end of the previous session was corrupted.
	// Read data from the end until the first valid key/value pair is found.
	isValid := false

	for segmentIndex := highestSegmentIndex; segmentIndex >= lowestSegmentIndex && segmentIndex+1 != 0; segmentIndex-- {
		if !segments[segmentIndex].IsSealed() {
			// ignore unsealed segment, this will have been created in the current session and will not
			// yet contain any data.
			continue
		}

		keys, err := segments[segmentIndex].GetKeys()
		if err != nil {
			return fmt.Errorf("failed to get keys from segment: %v", err)
		}
		for keyIndex := len(keys) - 1; keyIndex >= 0; keyIndex-- {
			key := keys[keyIndex]

			if !isValid {
				_, err = segments[segmentIndex].Read(key.Address)
				if err == nil {
					// we found a valid key/value pair. All subsequent keys are valid.
					isValid = true
				} else {
					// This is not cause for alarm (probably).
					// This can happen when the database is not cleanly shut down,
					// and just means that some data near the end was not fully committed.
					m.logger.Infof("truncated value for key %s with address %s", key.Key, key.Address)
				}
			}

			if isValid {
				m.data[string(key.Key)] = key.Address
			}
		}
	}

	return nil
}
