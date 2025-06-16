package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/urfave/cli/v2"
)

func lsCommand(ctx *cli.Context) error {
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

	tables, err := lsPaths(sources, true)
	if err != nil {
		return fmt.Errorf("failed to list tables in paths %v: %v", sources, err)
	}

	for _, table := range tables {
		fmt.Println(table)
	}

	return nil
}

// Similar to ls, but searches for tables in multiple paths.
func lsPaths(rootPaths []string, fsync bool) ([]string, error) {
	tableSet := make(map[string]struct{})

	for _, rootPath := range rootPaths {
		tables, err := ls(rootPath, fsync)
		if err != nil {
			return nil, fmt.Errorf("error finding tables: %v", err)
		}
		for _, table := range tables {
			tableSet[table] = struct{}{}
		}
	}

	tableNames := make([]string, 0, len(tableSet))
	for tableName := range tableSet {
		tableNames = append(tableNames, tableName)
	}

	sort.Strings(tableNames)

	return tableNames, nil
}

// Returns a list of LittDB tables at the specified LittDB path. Tables are alphabetically sorted by their names.
// Returns an error if the path does not exist or if no tables are found.
func ls(rootPath string, fsync bool) ([]string, error) {

	// Forbid touching tables in active use.
	lockPath := path.Join(rootPath, util.LockfileName)
	lock, err := util.NewFileLock(lockPath, fsync)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock on %s: %v", rootPath, err)
	}
	defer func() {
		_ = lock.Release()
	}()

	// LittDB has one directory under the root directory per table, with the name
	// of the table being the name of the directory.
	possibleTables, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %v", rootPath, err)
	}

	// Each table directory will contain a "segments" directory. Infer that any directory containing this directory
	// is a table. If we are looking at a real LittDB instance, there shouldn't be any other directories, but
	// there is no need to enforce that here.
	tables := make([]string, 0, len(possibleTables))
	for _, entry := range possibleTables {
		if !entry.IsDir() {
			continue
		}

		segmentPath := filepath.Join(rootPath, entry.Name(), segment.SegmentDirectory)
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
