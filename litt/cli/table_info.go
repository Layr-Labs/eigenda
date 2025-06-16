package main

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/litt"
	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

// TableInfo contains high level information about a table in LittDB.
type TableInfo struct {
	// The number of key-value pairs in the table.
	KeyCount uint64
	// The size of the table in bytes.
	Size uint64
	// If true, the table at the specified path is a snapshot of another table.
	IsSnapshot bool
	// The time when the oldest segment was sealed.
	OldestSegmentSealTime time.Time
	// The time when the newest segment was sealed. This is used to determine the age of the oldest segment.
	NewestSegmentSealTime time.Time
	// The index of the oldest segment in the table.
	LowestSegmentIndex uint32
	// The index of the newest segment in the table.
	HighestSegmentIndex uint32
	// The type of the keymap used by the table. If "", then this table doesn't have a keymap (i.e. it will rebuild
	// a keymap the next time it is loaded).
	KeymapType string
}

// TODO ensure boundary files are respected

// tableInfoCommand is the CLI command handler for the "table-info" command.
func tableInfoCommand(ctx *cli.Context) error {
	if ctx.NArg() != 1 {
		return fmt.Errorf(
			"table-info command requires exactly at least one argument: <table-name>")
	}

	tableName := ctx.Args().Get(0)

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

	info, err := tableInfo(tableName, sources, true)
	if err != nil {
		return fmt.Errorf("failed to get table info for table %s at paths %v: %v", tableName, sources, err)
	}

	oldestSegmentAge := uint64(time.Since(info.OldestSegmentSealTime).Nanoseconds())
	newestSegmentAge := uint64(time.Since(info.NewestSegmentSealTime).Nanoseconds())
	segmentSpan := oldestSegmentAge - newestSegmentAge

	// Print table information in a human-readable format
	fmt.Printf("Table:                       %s\n", tableName)
	fmt.Printf("Key count:                   %s\n", util.CommaOMatic(info.KeyCount))
	fmt.Printf("Size:                        %s\n", util.PrettyPrintBytes(info.Size))
	fmt.Printf("Is snapshot:                 %t\n", info.IsSnapshot)
	fmt.Printf("Oldest segment age:          %s\n", util.PrettyPrintTime(oldestSegmentAge))
	fmt.Printf("Oldest segment seal time:    %s\n", info.OldestSegmentSealTime.Format(time.RFC3339))
	fmt.Printf("Newest segment age:          %s\n", util.PrettyPrintTime(newestSegmentAge))
	fmt.Printf("Newest segment seal time:    %s\n", info.NewestSegmentSealTime.Format(time.RFC3339))
	fmt.Printf("Segment span:                %s\n", util.PrettyPrintTime(segmentSpan))
	fmt.Printf("Lowest segment index:        %d\n", info.LowestSegmentIndex)
	fmt.Printf("Highest segment index:       %d\n", info.HighestSegmentIndex)
	fmt.Printf("Key map type:                %s\n", info.KeymapType)

	return nil
}

// tableInfo retrieves information about a table at the specified path.
func tableInfo(tableName string, paths []string, fsync bool) (*TableInfo, error) {
	if !litt.IsTableNameValid(tableName) {
		return nil, fmt.Errorf("table name '%s' is invalid, "+
			"must be at least one character long and contain only letters, numbers, underscores, and dashes",
			tableName)
	}

	// Forbid touching tables in active use.
	for _, rootPath := range paths {
		lockPath := path.Join(rootPath, util.LockfileName)
		lock, err := util.NewFileLock(lockPath, fsync)
		if err != nil {
			return nil, fmt.Errorf("failed to acquire lock on %s: %v", rootPath, err)
		}
		defer func() {
			_ = lock.Release()
		}()
	}

	segmentPaths, err := segment.BuildSegmentPaths(paths, "", tableName)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to build segment paths for table %s at paths %v: %v", tableName, paths, err)
	}

	for _, segmentPath := range segmentPaths {
		exists, err := util.Exists(segmentPath.SegmentDirectory())
		if err != nil {
			return nil, fmt.Errorf("failed to check if segment directory %s exists: %v",
				segmentPath.SegmentDirectory(), err)
		}
		if !exists {
			return nil, fmt.Errorf("segment directory %s does not exist", segmentPath.SegmentDirectory())
		}
	}

	logger, err := common.NewLogger(common.DefaultTextLoggerConfig())
	require.NoError(nil, err)
	errorMonitor := util.NewErrorMonitor(context.Background(), logger, nil)

	lowestSegmentIndex, highestSegmentIndex, segments, err := segment.GatherSegmentFiles(
		logger,
		errorMonitor,
		segmentPaths,
		time.Now(),
		false,
		true)

	if err != nil {
		return nil, fmt.Errorf("failed to gather segment files for table %s at paths %v: %v",
			tableName, paths, err)
	}
	if ok, err := errorMonitor.IsOk(); !ok {
		// This should be impossible since we aren't doing anything on background threads that report to the
		// error monitor, but it doesn't hurt to check.
		return nil, fmt.Errorf("error monitor reports errors: %v", err)
	}

	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments found for table %s at paths %v", tableName, paths)
	}

	keyCount := uint64(0)
	size := uint64(0)
	for _, seg := range segments {
		keyCount += uint64(seg.KeyCount())
		size += seg.Size()
	}

	_, _, keymapTypeFile, err := littbuilder.FindKeymapLocation(paths, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to find keymap location for table %s at paths %v: %v",
			tableName, paths, err)
	}

	keymapType := "none (will be rebuilt on next LittDB startup)"
	if keymapTypeFile != nil {
		keymapType = (string)(keymapTypeFile.Type())
	}

	isSnapshot, err := segments[lowestSegmentIndex].IsSnapshot()
	if err != nil {
		return nil, fmt.Errorf("failed to check if segment %d is a snapshot: %v", lowestSegmentIndex, err)
	}

	return &TableInfo{
		KeyCount:              keyCount,
		Size:                  size,
		IsSnapshot:            isSnapshot,
		OldestSegmentSealTime: segments[lowestSegmentIndex].GetSealTime(),
		NewestSegmentSealTime: segments[highestSegmentIndex].GetSealTime(),
		LowestSegmentIndex:    lowestSegmentIndex,
		HighestSegmentIndex:   highestSegmentIndex,
		KeymapType:            keymapType,
	}, nil
}
