package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/urfave/cli/v2"
)

// called by the CLI to unlock a LittDB file system.
func unlockCommand(ctx *cli.Context) error {
	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	sources := ctx.StringSlice(srcFlag.Name)

	if len(sources) == 0 {
		return fmt.Errorf("at least one source path is required")
	}

	return Unlock(logger, sources)
}

// Unlocks a LittDB file system.
//
// DANGER: calling this method opens the door for unsafe concurrent operations on LittDB files.
// With great power comes great responsibility.
func Unlock(logger logging.Logger, sourcePaths []string) error {
	for _, sourcePath := range sourcePaths {
		err := filepath.WalkDir(sourcePath, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			if strings.HasSuffix(path, util.LockfileName) {
				logger.Infof("Removing lock file %s", path)
				if removeErr := os.Remove(path); removeErr != nil {
					logger.Error("Failed to remove lock file", "path", path, "error", removeErr)
					return fmt.Errorf("failed to remove lock file %s: %w", path, removeErr)
				}
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to walk directory %s: %w", sourcePath, err)
		}
	}

	return nil
}
