package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/urfave/cli/v2"
)

// TODO unit test

// Returns a list of LittDB tables at the specified LittDB path. Tables are alphabetically sorted by their names.
// Returns an error if the path does not exist or if no tables are found.
func ls(path string) ([]string, error) {

	// LittDB has one directory under the root directory per table, with the name
	// of the table being the name of the directory.
	possibleTables, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %v", path, err)
	}

	// Each table directory will contain a "segments" directory. Infer that any directory containing this directory
	// is a table. If we are looking at a real LittDB instance, there shouldn't be any other directories, but
	// there is no need to enforce that here.
	tables := make([]string, 0, len(possibleTables))
	for _, entry := range possibleTables {
		segmentPath := filepath.Join(path, entry.Name(), segment.SegmentDirectory)
		exists, err := util.Exists(segmentPath)
		if err != nil {
			return nil, fmt.Errorf("failed to check if segment path %s exists: %v", segmentPath, err)
		}
		if exists && entry.IsDir() {
			tables = append(tables, entry.Name())
		}
	}

	// Alphabetically sort the tables.
	sort.Strings(tables)

	return tables, nil
}

func lsCommand(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf("ls command requires exactly one argument: <path>")
	}
	path := ctx.Args().Get(0)
	path, err := util.SanitizePath(path)
	if err != nil {
		return fmt.Errorf("failed to sanitize path %s: %v", path, err)
	}

	tables, err := ls(path)
	if err != nil {
		return fmt.Errorf("failed to list tables in path %s: %v", path, err)
	}

	for _, table := range tables {
		fmt.Println(table)
	}

	return nil
}
