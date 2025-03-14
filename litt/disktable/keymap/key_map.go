package keymap

import (
	"github.com/Layr-Labs/eigenda/litt/types"
)

// KeyMap maintains a mapping between keys and addresses.
type KeyMap interface {
	// Put adds keys to the key map as a batch.
	Put(pairs []*types.KAPair) error

	// Get returns the address for a key. Returns true if the key exists, and false otherwise (i.e. does not
	// return an error if the key does not exist).
	Get(key []byte) (types.Address, bool, error)

	// Delete removes keys from the key map.
	Delete(keys []*types.KAPair) error

	// Stop stops the key map.
	Stop() error

	// Destroy stops the key map and permanently deletes all data.
	Destroy() error
}
