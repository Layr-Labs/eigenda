package keymap

import "github.com/Layr-Labs/eigensdk-go/logging"

// KeymapBuilder is an interface for building and managing keymaps.
type KeymapBuilder interface {

	// Type returns the type of the keymap builder.
	Type() KeymapType

	// Build creates a new keymap. Keymap implementations that do not
	// store files on disk can instantiate an empty keymap and return it.
	//
	// If the returned boolean is true, then the keymap requires a reload of key/address pairs. For in-memory
	// implementations, this will always return true. For disk-based implementations, this will return true if the
	// Keymap's files are present on disk, and false otherwise.
	Build(logger logging.Logger, keymapPath string, doubleWriteProtection bool) (Keymap, bool, error)

	// DeleteFiles deletes all files associated with the keymap that are located in any of the provided path.
	DeleteFiles(logger logging.Logger, keymapPath string) error
}
