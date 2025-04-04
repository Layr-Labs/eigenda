package keymap

import (
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

var _ KeymapBuilder = &LevelDBKeymapBuilder{}

// LevelDBKeymapBuilder is a KeymapBuilder that builds LevelDBKeymap instances.
type LevelDBKeymapBuilder struct {
}

// NewLevelDBKeymapBuilder creates a new LevelDBKeymapBuilder.
func NewLevelDBKeymapBuilder() *LevelDBKeymapBuilder {
	return &LevelDBKeymapBuilder{}
}

func (b *LevelDBKeymapBuilder) Type() KeymapType {
	return LevelDBKeymapType
}

func (b *LevelDBKeymapBuilder) Build(
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

	keymap, err := NewLevelDBKeymap(logger, keymapPath, doubleWriteProtection)
	if err != nil {
		return nil, false, fmt.Errorf("error creating LevelDBKeymap: %w", err)
	}

	requiresReload := !exists
	return keymap, requiresReload, nil
}
