package main

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/litt/disktable"
	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/urfave/cli/v2"
)

func pushCommand(ctx *cli.Context) error {
	if ctx.NArg() < 1 {
		return fmt.Errorf("not enough arguments provided, must provide USER@HOST")
	}

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
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

	return Push(logger, sources, destinations, user, host, port, keyPath, deleteAfterTransfer, true, verbose)
}

// Push uses rsync to transfer LittDB data to the remote location(s)
func Push(
	logger logging.Logger,
	sources []string,
	destinations []string,
	user string,
	host string,
	port uint64,
	keyPath string,
	deleteAfterTransfer bool,
	fsync bool,
	verbose bool) error {

	if len(sources) == 0 {
		return fmt.Errorf("no source paths provided")
	}
	if len(destinations) == 0 {
		return fmt.Errorf("no destination paths provided")
	}

	// Lock source files. It would be nice to also lock the remote directories, but that's tricky given that
	// we are interacting with the remote machine via SSH and rsync.
	releaseSourceLocks, err := util.LockDirectories(logger, sources, util.LockfileName, fsync)
	if err != nil {
		return fmt.Errorf("failed to lock source directories: %v", err)
	}
	defer releaseSourceLocks()

	// Create an SSH session to the remote host.
	connection, err := util.NewSSHSession(logger, user, host, port, keyPath, verbose)
	if err != nil {
		return fmt.Errorf("failed to create SSH session to %s@%s port %d: %v", user, host, port, err)
	}

	// Figure out where data currently exists at the destination(s). We don't want this operation to cause a file
	// to exist in multiple places.
	// TODO make sure this handles when there are multiple tables.
	existingFilesMap, err := mapExistingFiles(logger, destinations, connection)
	if err != nil {
		return fmt.Errorf("failed to map existing files at destinations: %v", err)
	}

	tables, err := lsPaths(logger, sources, false, fsync)
	if err != nil {
		return fmt.Errorf("failed to list tables in source paths %v: %v", sources, err)
	}

	for _, tableName := range tables {
		err = pushTable(
			logger,
			tableName,
			sources,
			destinations,
			connection,
			existingFilesMap,
			deleteAfterTransfer,
			fsync,
			verbose,
		)

		if err != nil {
			return fmt.Errorf("failed to push table %s: %v", tableName, err)
		}
	}

	return nil
}

// Figure out which files are already present at the destination(s). Although these files may be partial, we always
// want to preserve any pre-existing arrangements of files at the destination(s).
//
// The returned map is a map from file name (e.g. 1234.metadata) to the destination path (e.g. /path/to/remote/dir).
func mapExistingFiles(
	logger logging.Logger,
	destinations []string,
	connection *util.SSHSession) (map[string]string, error) {

	existingFiles := make(map[string]string)

	extensions := []string{segment.MetadataFileExtension, segment.KeyFileExtension, segment.ValuesFileExtension}

	for _, dest := range destinations {
		filePaths, err := connection.FindFiles(dest, extensions)
		if err != nil {
			return nil, fmt.Errorf("failed to list files in destination %s: %v", dest, err)
		}

		for _, filePath := range filePaths {
			// Extract the file name from the path.
			fileName := path.Base(filePath)
			if _, exists := existingFiles[fileName]; !exists {
				existingFiles[fileName] = dest
			} else {
				logger.Warnf("File %s already exists in destination %s, skipping", fileName, dest)
			}
		}
	}

	return existingFiles, nil
}

// Push the data in a single table to the remote location(s).
func pushTable(
	logger logging.Logger,
	tableName string,
	sources []string,
	destinations []string,
	connection *util.SSHSession,
	existingFilesMap map[string]string,
	deleteAfterTransfer bool,
	fsync bool,
	verbose bool) error {

	segmentPaths, err := segment.BuildSegmentPaths(sources, "", tableName)
	if err != nil {
		return fmt.Errorf("failed to build segment paths for table %s at paths %v: %v", tableName, sources, err)
	}

	errorMonitor := util.NewErrorMonitor(context.Background(), logger, nil)

	// Gather segment files to send.
	lowestSegmentIndex, highestSegmentIndex, segments, err := segment.GatherSegmentFiles(
		logger,
		errorMonitor,
		segmentPaths,
		time.Now(),
		false,
		fsync)
	if err != nil {
		return fmt.Errorf("failed to gather segment files for table %s at paths %v: %v",
			tableName, sources, err)
	}

	if len(segments) == 0 {
		logger.Infof("No segments found for table %s", tableName)
		return nil
	}

	// Special handling if we are transferring data from a snapshot.
	isSnapshot, err := segments[lowestSegmentIndex].IsSnapshot()
	if err != nil {
		return fmt.Errorf("failed to check if segment %d is a snapshot: %v", lowestSegmentIndex, err)
	}
	if isSnapshot {
		if len(sources) > 1 {
			return fmt.Errorf("table %s is a snapshot, but source directories found: %v", tableName, sources)
		}

		boundaryFile, err := disktable.LoadBoundaryFile(false, path.Join(sources[0], tableName))
		if err != nil {
			return fmt.Errorf("failed to load boundary file for table %s at path %s: %v",
				tableName, sources[0], err)
		}

		if boundaryFile.IsDefined() {
			highestSegmentIndex = boundaryFile.BoundaryIndex()
		}
	}

	// Ensure the remote segment directories exists.
	for _, dest := range destinations {
		segmentDir := path.Join(dest, tableName, segment.SegmentDirectory)
		err = connection.Mkdirs(segmentDir)
		if err != nil {
			return fmt.Errorf("failed to create segment directory %s at destination %s: %v",
				segmentDir, dest, err)
		}
	}

	// Transfer the files.
	for i := lowestSegmentIndex; i <= highestSegmentIndex; i++ {
		seg := segments[i]
		filesToTransfer := seg.GetFilePaths()

		for _, filePath := range filesToTransfer {
			fileName := path.Base(filePath)

			destination := ""
			if existingDest, exists := existingFilesMap[fileName]; exists {
				destination = existingDest
			} else {
				destination, err = determineDestination(fileName, destinations)
				if err != nil {
					return fmt.Errorf("failed to determine destination for file %s: %v", fileName, err)
				}
			}

			targetLocation := path.Join(destination, tableName, segment.SegmentDirectory, fileName)

			err = connection.Rsync(filePath, targetLocation)
			if err != nil {
				return fmt.Errorf("failed to rsync file %s to %s: %v", filePath, targetLocation, err)
			}
		}
	}

	// Now that we have transferred the files, we can delete them if requested.
	if deleteAfterTransfer {
		for _, seg := range segments {
			seg.Release()
		}
		for _, seg := range segments {
			err = seg.BlockUntilFullyDeleted()
			if err != nil {
				return fmt.Errorf("failed to delete segment %d for table %s: %v",
					seg.SegmentIndex(), tableName, err)
			}
		}

		if isSnapshot {
			// If we are dealing with a snapshot, update the lower bound file.
			boundaryFile, err := disktable.LoadBoundaryFile(true, path.Join(destinations[0], tableName))
			if err != nil {
				return fmt.Errorf("failed to load boundary file for table %s at path %s: %v",
					tableName, destinations[0], err)
			}

			err = boundaryFile.Update(highestSegmentIndex)
			if err != nil {
				return fmt.Errorf("failed to update boundary file for table %s at path %s: %v",
					tableName, destinations[0], err)
			}
		}
	}

	return nil
}
