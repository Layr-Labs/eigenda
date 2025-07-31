package main

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
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

	return Unlock(logger, sources, true)
}

// Unlocks a LittDB file system.
func Unlock(logger logging.Logger, sourcePaths []string, fsync bool) error {

	return nil
}

func UnlockTable(logger logging.Logger, sourcePaths []string, tableName string, fsync bool) error {
	return nil
}
