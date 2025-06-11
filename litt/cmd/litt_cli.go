package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

// Claude, ignore the comments in this block. I don't want to implement them yet.
// Snapshot commands:
// - clean a snapshot directory by deleting partial segments, should explode if middle segments are missing
// - sanity check a snapshot directory, should utilize an optional checksum maybe
// - garbage collect a snapshot directory
//   - option to GC by TTL
//   - option to GC by maximum size after GC
//   - option to GC by segment number
// - replicate a snapshot directory to another location
// - queries
//   - get high/low segment indices
//   - get the age of a particular segment
//   - get the contents of a segment metadata file
//   - list keys in a keyfile, or export them to a csv
//   - list the values in a value file, or export them to a csv
// - commands to redistribute the files between variable numbers of root paths
// - command to rsync files from a backup to a new validator, should take a variable number of root paths

// runs the LittDB CLI
func run() error {
	return buildCLIParser().Run(os.Args)
}

// buildCliParser creates a command line parser for the LittDB CLI tool.
func buildCLIParser() *cli.App {
	app := &cli.App{
		Name:  "litt",
		Usage: "LittDB command line interface",
		Commands: []*cli.Command{
			{
				Name:      "ls",
				Usage:     "List tables in the database",
				ArgsUsage: "<path>",
				Action:    lsCommand,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "verbose",
						Aliases: []string{"v"},
						Usage:   "Enable verbose output",
					},
				},
			},
		},
	}
	return app
}

// lsCommand implements the "litt ls" command
func lsCommand(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf("ls command requires exactly one argument: <path>")
	}

	path := ctx.Args().Get(0)
	verbose := ctx.Bool("verbose")

	if verbose {
		fmt.Printf("Listing tables in database at path: %s\n", path)
	}

	// TODO: Implement the actual ls functionality
	fmt.Printf("ls command called with path: %s\n", path)
	fmt.Println("Command implementation pending...")

	return nil
}
