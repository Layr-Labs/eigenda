package keymap

import (
	"fmt"
	"os"

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
	exists := false
	_, err := os.Stat(keymapPath)
	if err == nil {
		exists = true
	} else if os.IsNotExist(err) {
		err = os.MkdirAll(keymapPath, 0755)
		if err != nil {
			return nil, false, fmt.Errorf("error creating keymap directory: %w", err)
		}
	} else {
		return nil, false, fmt.Errorf("error checking for keymap directory: %w", err)
	}

	keymap, err := NewLevelDBKeymap(logger, keymapPath, doubleWriteProtection, true)
	if err != nil {
		return nil, false, fmt.Errorf("error creating LevelDBKeymap: %w", err)
	}

	requiresReload := !exists
	return keymap, requiresReload, nil
}

func (b *LevelDBKeymapBuilder) DeleteFiles(logger logging.Logger, keymapPath string) error {
	_, err := os.Stat(keymapPath)
	if err == nil {
		logger.Infof("deleting keymap directory: %s", keymapPath)
		err = os.RemoveAll(keymapPath)
		if err != nil {
			return fmt.Errorf("error deleting keymap directory: %w", err)
		}
	}

	return nil
}
