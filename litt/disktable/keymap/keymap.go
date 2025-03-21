package keymap

import (
	"github.com/Layr-Labs/eigenda/litt/types"
)

// KeymapDirectoryName is the name of the directory where the keymap stores its files.
const KeymapDirectoryName = "keymap"

// KeymapDataDirectoryName is the name of the directory where the keymap implementation stores its data files.
const KeymapDataDirectoryName = "data"

// Keymap maintains a mapping between keys and addresses.
type Keymap interface {
	// Put adds keys to the keymap as a batch.
	Put(pairs []*types.KAPair) error

	// Get returns the address for a key. Returns true if the key exists, and false otherwise (i.e. does not
	// return an error if the key does not exist).
	Get(key []byte) (types.Address, bool, error)

	// Delete removes keys from the keymap.
	Delete(keys []*types.KAPair) error

	// Stop stops the keymap.
	Stop() error

	// Destroy stops the keymap and permanently deletes all data.
	Destroy() error
}
