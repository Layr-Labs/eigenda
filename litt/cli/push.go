package main

import (
	"fmt"
	"strings"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/urfave/cli/v2"
)

func pushCommand(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return fmt.Errorf("not enough arguments provided, must provide USER@HOST")
	}

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	if err != nil {
		return fmt.Errorf("failed to create logger: %v", err)
	}

	sources := ctx.StringSlice("src")
	if len(sources) == 0 {
		return fmt.Errorf("no sources provided")
	}
	for i, src := range sources {
		var err error
		sources[i], err = util.SanitizePath(src)
		if err != nil {
			return fmt.Errorf("Invalid source path: %s", src)
		}
	}

	destinations := ctx.StringSlice("dest")
	if len(destinations) == 0 {
		return fmt.Errorf("no destinations provided")
	}
	for i, dest := range destinations {
		var err error
		destinations[i], err = util.SanitizePath(dest)
		if err != nil {
			return fmt.Errorf("Invalid source path: %s", dest)
		}
	}

	userHost := ctx.Args().First()
	parts := strings.Split(userHost, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid USER@HOST format: %s", userHost)
	}
	user := parts[0]
	host := parts[1]

	port := ctx.Uint64("port")

	keyPath := ctx.String("key")
	keyPath, err = util.SanitizePath(keyPath)
	if err != nil {
		return fmt.Errorf("Invalid key path: %s", keyPath)
	}

	deleteAfterTransfer := !ctx.Bool("no-gc")

	verbose := !ctx.Bool("quiet")

	return Push(logger, sources, destinations, user, host, port, keyPath, deleteAfterTransfer, verbose)
}

// Push uses rsync to transfer LittDB data to a remote location(s)
func Push(
	logger logging.Logger,
	sources []string,
	destinations []string,
	user string,
	host string,
	port uint64,
	keyPath string,
	deleteAfterTransfer bool,
	verbose bool) error {

	// TODO lock source dirs

	return nil

}
