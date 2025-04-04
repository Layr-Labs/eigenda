package keymap

import (
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ KeymapBuilder = &LevelDBKeymapBuilder{}

// UnsafeLevelDBKeymapBuilder is a KeymapBuilder that builds LevelDBKeymap instances. It runs with sync writes disabled.
// This is much faster than the default LevelDBKeymapBuilder, but it is not safe to use in production. This is only
// intended for use in tests.
type UnsafeLevelDBKeymapBuilder struct {
}

// NewUnsafeLevelDBKeymapBuilder creates a new UnsafeLevelDBKeymapBuilder.
func NewUnsafeLevelDBKeymapBuilder() *LevelDBKeymapBuilder {
	return &LevelDBKeymapBuilder{}
}

func (b *UnsafeLevelDBKeymapBuilder) Type() KeymapType {
	return UnsafeLevelDBKeymapType
}

func (b *UnsafeLevelDBKeymapBuilder) Build(
	logger logging.Logger,
	keymapPath string,
	doubleWriteProtection bool) (Keymap, bool, error) {

	// check to see if the keymap directory exists in one of the provided paths
	exists, err := util.Exists(keymapPath)
	if err != nil {
		return nil, false, fmt.Errorf("error checking for keymap directory: %w", err)
	}

	if !exists {
		err = os.MkdirAll(keymapPath, 0755)
		if err != nil {
			return nil, false, fmt.Errorf("error creating keymap directory: %w", err)
		}
	}

	keymap, err := NewUnsafeLevelDBKeymap(logger, keymapPath, doubleWriteProtection)
	if err != nil {
		return nil, false, fmt.Errorf("error creating LevelDBKeymap: %w", err)
	}

	requiresReload := !exists
	return keymap, requiresReload, nil
}
