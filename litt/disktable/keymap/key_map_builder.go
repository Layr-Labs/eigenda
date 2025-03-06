package keymap

import "github.com/Layr-Labs/eigensdk-go/logging"

// KeyMapBuilder is an interface for building and managing key maps.
type KeyMapBuilder interface {

	// Type returns the type of the key map builder.
	Type() KeyMapType

	// Build creates a new key map based on the provided paths. KeyMap implementations that do not
	// store files on disk can instantiate an empty key map and return it.
	//
	// If the returned boolean is true, then the keymap requires a reload of key/address pairs. For in-memory
	// implementations, this will always return true. For disk-based implementations, this will return true if the
	// KeyMap's files are present on disk, and false otherwise.
	Build(logger logging.Logger, paths []string) (KeyMap, bool, error)

	// DeleteFiles deletes all files associated with the key map that are located in any of the provided paths.
	// This may be called even if there is no key map in the provided paths (this method should be a no-op in
	// that case).
	DeleteFiles(logger logging.Logger, paths []string) error
}
