package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/litt/disktable"
	"github.com/Layr-Labs/eigenda/litt/disktable/keymap"
	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/urfave/cli/v2"
)

// pruneCommand can be used to remove data from a LittDB instance/snapshot.
func pruneCommand(ctx *cli.Context) error {
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

	tables := ctx.StringSlice("table")

	maxAgeSeconds := ctx.Uint64("max-age")

	return prune(sources, tables, maxAgeSeconds, true)
}

// prune deletes data from a littDB database/snapshot.
func prune(sources []string, allowedTables []string, maxAgeSeconds uint64, fsync bool) error {
	allowedTablesSet := make(map[string]struct{})
	for _, table := range allowedTables {
		allowedTablesSet[table] = struct{}{}
	}

	// Determine which tables to prune.
	var tables []string
	foundTables, err := lsPaths(sources, fsync)
	if err != nil {
		return fmt.Errorf("failed to list tables in paths %v: %v", sources, err)
	}
	if len(allowedTables) == 0 {
		tables = foundTables
	} else {
		for _, table := range foundTables {
			if _, ok := allowedTablesSet[table]; ok {
				tables = append(tables, table)
			}
		}
	}

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	if err != nil {
		return fmt.Errorf("failed to create logger: %v", err)
	}

	// Prune each table.
	for _, table := range tables {

		fmt.Printf("Pruning table %s\n", table) // TODO

		bytesDeleted, err := pruneTable(logger, sources, table, maxAgeSeconds, fsync)
		if err != nil {
			return fmt.Errorf("failed to prune table %s in paths %v: %v", table, sources, err)
		}

		fmt.Printf("Deleted %s from table '%s'.", util.PrettyPrintBytes(bytesDeleted), table)
	}

	return nil
}

// pruneTable performs offline garbage collection on a LittDB database/snapshot.
func pruneTable(
	logger logging.Logger,
	sources []string,
	tableName string,
	maxAgeSeconds uint64,
	fsync bool) (uint64, error) {

	// TODO grab lock files!

	errorMonitor := util.NewErrorMonitor(context.Background(), logger, nil)

	segmentPaths, err := segment.BuildSegmentPaths(sources, "", tableName)
	if err != nil {
		return 0, fmt.Errorf("failed to build segment paths for table %s at paths %v: %v",
			tableName, sources, err)
	}

	lowestSegmentIndex, highestSegmentIndex, segments, err := segment.GatherSegmentFiles(
		logger,
		errorMonitor,
		segmentPaths,
		time.Now(),
		false,
		fsync)
	if err != nil {
		return 0, fmt.Errorf("failed to gather segment files for table %s at paths %v: %v",
			tableName, sources, err)
	}

	if len(segments) == 0 {
		return 0, fmt.Errorf("no segments found for table %s at paths %v", tableName, sources)
	}

	// Determine if we are working on the snapshot directory (i.e. the directory with symlinks to the segments).
	isSnapshot, err := segments[lowestSegmentIndex].IsSnapshot()
	if err != nil {
		return 0, fmt.Errorf("failed to check if segment %d is a snapshot: %v", lowestSegmentIndex, err)
	}

	if isSnapshot {
		// If we are dealing with a snapshot, respect the snapshot boundary file.

		if len(sources) > 1 {
			return 0, fmt.Errorf("this is a symlinked snapshot directory, " +
				"snapshot directory cannot be spread across multiple sources.")
		}

		boundaryFile, err := disktable.LoadBoundaryFile(false, path.Join(sources[0], tableName))
		if err != nil {
			return 0, fmt.Errorf("failed to load boundary file for table %s at path %s: %v",
				tableName, sources[0], err)
		}

		if boundaryFile.IsDefined() {
			highestSegmentIndex = boundaryFile.BoundaryIndex()
		}
	}

	// Delete old segments.
	bytesDeleted := uint64(0)
	deletedSegments := make([]*segment.Segment, 0)
	for segmentIndex := lowestSegmentIndex; segmentIndex <= highestSegmentIndex; segmentIndex++ {
		seg := segments[segmentIndex]
		segmentAge := time.Since(seg.GetSealTime())

		if segmentAge < time.Duration(maxAgeSeconds)*time.Second {
			// We've pruned all segments that we can.
			break
		}

		deletedSegments = append(deletedSegments, seg)
		bytesDeleted += seg.Size()
		seg.Release()
	}

	// Wait for deletion to complete.
	for _, seg := range deletedSegments {
		err = seg.BlockUntilFullyDeleted()
		if err != nil {
			return 0, fmt.Errorf("failed to block until segment %d is fully deleted: %v",
				seg.SegmentIndex(), err)
		}
	}

	if ok, err := errorMonitor.IsOk(); !ok {
		return 0, fmt.Errorf("error monitor reports errors: %v", err)
	}

	if !isSnapshot {
		// If we are doing GC on a table that isn't a snapshot, then we need to delete the snapshots/keymap
		// for the table. The DB will automatically rebuild the snapshots directory & keymap on the next startup.

		for _, source := range sources {
			snapshotsPath := path.Join(source, tableName, segment.HardLinkDirectory)
			exists, err := util.Exists(snapshotsPath)
			if err != nil {
				return 0, fmt.Errorf("failed to check if snapshots path %s exists: %v", snapshotsPath, err)
			}
			if exists {
				err = os.RemoveAll(snapshotsPath)
				if err != nil {
					return 0, fmt.Errorf("failed to remove snapshots path %s: %v", snapshotsPath, err)
				}
			}

			keymapPath := path.Join(source, tableName, keymap.KeymapDirectoryName)
			exists, err = util.Exists(keymapPath)
			if err != nil {
				return 0, fmt.Errorf("failed to check if keymap path %s exists: %v", keymapPath, err)
			}
			if exists {
				err = os.RemoveAll(keymapPath)
				if err != nil {
					return 0, fmt.Errorf("failed to remove keymap path %s: %v", keymapPath, err)
				}
			}
		}
	}

	return bytesDeleted, nil
}
