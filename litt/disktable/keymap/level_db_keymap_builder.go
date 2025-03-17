package keymap

import (
	"fmt"
	"os"
	"path"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// directoryName is the name of the directory where the LevelDBKeymap stores its files.
const directoryName = "ldb"

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

func (b *LevelDBKeymapBuilder) Build(logger logging.Logger, paths []string) (Keymap, bool, error) {
	if len(paths) == 0 {
		return nil, false, fmt.Errorf("no paths provided")
	}

	// check to see if the keymap directory exists in one of the provided paths
	exists := false
	targetPath := ""
	for _, potentialRoot := range paths {
		potentialPath := path.Join(potentialRoot, KeymapDirectoryName, directoryName)
		_, err := os.Stat(potentialPath)
		if err == nil {
			exists = true
			targetPath = potentialPath
			break
		} else if !os.IsNotExist(err) {
			return nil, false, fmt.Errorf("error checking for keymap directory: %w", err)
		}
	}

	if !exists {
		// if the keymap directory does not exist, create it in the first path.
		targetPath = path.Join(paths[0], KeymapDirectoryName, directoryName)
	}

	keymap, err := NewLevelDBKeymap(logger, targetPath)
	if err != nil {
		return nil, false, fmt.Errorf("error creating LevelDBKeymap: %w", err)
	}

	requiresReload := !exists
	return keymap, requiresReload, nil
}

func (b *LevelDBKeymapBuilder) DeleteFiles(logger logging.Logger, paths []string) error {
	for _, potentialRoot := range paths {
		potentialPath := path.Join(potentialRoot, directoryName)
		_, err := os.Stat(potentialPath)
		if err == nil {
			logger.Infof("deleting keymap directory: %s", potentialPath)
			err = os.RemoveAll(potentialPath)
			if err != nil {
				return fmt.Errorf("error deleting keymap directory: %w", err)
			}
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("error checking for keymap directory: %w", err)
		}
	}
	return nil
}
