package keymap

import (
	"github.com/Layr-Labs/eigenda/litt/types"
)

// KeymapDirectoryName is the name of the directory where the keymap stores its files.
const KeymapDirectoryName = "keymap"

// KeymapDataDirectoryName is the name of the directory where the keymap implementation stores its data files.
// This directory will be created inside the keymap directory.
const KeymapDataDirectoryName = "data"

// Keymap maintains a mapping between keys and addresses. Implementations of this interface are goroutine safe.
type Keymap interface {
	// Put adds keys to the keymap as a batch.
	//
	// It is not thread safe to modify the contents of any slices passed to this function after the call.
	// This includes the byte slices containing the keys.
	Put(pairs []*types.KAPair) error

	// Get returns the address for a key. Returns true if the key exists, and false otherwise (i.e. does not
	// return an error if the key does not exist).
	//
	// It is not thread safe to modify key byte slice after it is passed to this method.
	Get(key []byte) (types.Address, bool, error)

	// Delete removes keys from the keymap. Deleting non-existent keys is a no-op.
	//
	// It is not thread safe to modify the contents of any slices passed to this function after the call.
	// This includes the byte slices containing the keys.
	Delete(keys []*types.KAPair) error

	// Stop stops the keymap.
	Stop() error

	// Destroy stops the keymap and permanently deletes all data.
	Destroy() error
}
