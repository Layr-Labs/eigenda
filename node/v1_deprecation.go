package node

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Layr-Labs/eigensdk-go/logging"
)

// The subdirectory name where v1 chunk data is stored.
const V1ChunkSubdir = "chunk"

// Deletes the v1 data directory if it exists.
//
// Returns nil if an error occurs while deleting
func DeleteV1Data(logger logging.Logger, dbPath string) error {
	v1DataPath := filepath.Join(dbPath, V1ChunkSubdir)

	info, err := os.Stat(v1DataPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info("No v1 data found to delete", "path", v1DataPath)
			return nil
		}
		return fmt.Errorf("stat v1 data path %s: %w", v1DataPath, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("v1 data path %s exists but is not a directory", v1DataPath)
	}

	if err := os.RemoveAll(v1DataPath); err != nil {
		return fmt.Errorf("delete v1 data at %s: %w", v1DataPath, err)
	}

	logger.Info("Deleted v1 data", "path", v1DataPath)
	return nil
}
