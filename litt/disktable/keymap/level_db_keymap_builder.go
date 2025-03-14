package keymap

import (
	"fmt"
	"os"
	"path"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// directoryName is the name of the directory where the LevelDBKeyMap stores its files.
const directoryName = "ldb-keymap"

var _ KeyMapBuilder = &LevelDBKeyMapBuilder{}

// LevelDBKeyMapBuilder is a KeyMapBuilder that builds LevelDBKeyMap instances.
type LevelDBKeyMapBuilder struct {
}

// NewLevelDBKeyMapBuilder creates a new LevelDBKeyMapBuilder.
func NewLevelDBKeyMapBuilder() *LevelDBKeyMapBuilder {
	return &LevelDBKeyMapBuilder{}
}

func (b *LevelDBKeyMapBuilder) Type() KeyMapType {
	return LevelDBKeyMapType
}

func (b *LevelDBKeyMapBuilder) Build(logger logging.Logger, paths []string) (KeyMap, bool, error) {
	if len(paths) == 0 {
		return nil, false, fmt.Errorf("no paths provided")
	}

	// check to see if the keymap directory exists in one of the provided paths
	exists := false
	targetPath := ""
	for _, potentialRoot := range paths {
		potentialPath := path.Join(potentialRoot, directoryName)
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
		targetPath = path.Join(paths[0], directoryName)
	}

	keymap, err := NewLevelDBKeyMap(logger, targetPath)
	if err != nil {
		return nil, false, fmt.Errorf("error creating LevelDBKeyMap: %w", err)
	}

	requiresReload := !exists
	return keymap, requiresReload, nil
}

func (b *LevelDBKeyMapBuilder) DeleteFiles(logger logging.Logger, paths []string) error {
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
