package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"os"
	"path"
	"path/filepath"
	"sync/atomic"

	"github.com/Layr-Labs/eigenda/litt/disktable"
	"github.com/Layr-Labs/eigenda/litt/disktable/keymap"
	"github.com/Layr-Labs/eigenda/litt/disktable/segment"
	"github.com/Layr-Labs/eigenda/litt/littbuilder"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/urfave/cli/v2"
)

// rebaseCommand is the command to rebase a LittDB database.
func rebaseCommand(ctx *cli.Context) error {
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

	deep := !ctx.Bool("shallow")
	preserveOriginal := ctx.Bool("preserve")
	verbose := !ctx.Bool("quiet")

	// TODO figure out how rebasing a snapshot directory works

	return rebase(sources, destinations, deep, preserveOriginal, true, verbose)
}

// Files to manage during a rebase:
// - litt.lock: delete if discovered in directory that is going to be deleted
// - table/keymap: copy/move if source goes away
// - table/segments: copy/move if source goes away
// - table/segments/{metadata/keys/values}: copy/move if source goes away
// - table/table.metadata: move/copy if source goes away
// - table/snapshot: if it exists and the source goes away, ensure that there are equivalent snapshots in the destination

// rebase moves LittDB database files from one location to another (locally). This function is idempotent. If it
// crashes part of the way through, just run it again and it will continue where it left off.
func rebase(
	sources []string,
	destinations []string,
	deep bool,
	preserveOriginal bool,
	fsync bool,
	verbose bool,
) error {

	sourceSet := make(map[string]struct{})
	destinationSet := make(map[string]struct{})

	for _, src := range sources {
		exists, err := util.Exists(src)
		if err != nil {
			return fmt.Errorf("error checking if source path %s exists: %w", src, err)
		}
		// Ignore non-existent source paths. They could have been deleted by a prior run of this command.
		if exists {
			sourceSet[src] = struct{}{}
		}
	}

	for _, dest := range destinations {
		destinationSet[dest] = struct{}{}

		exists, err := util.Exists(dest)
		if err != nil {
			return fmt.Errorf("error checking if destination path %s exists: %w", dest, err)
		}
		if !exists {
			err = os.MkdirAll(dest, 0755)
			if err != nil {
				return fmt.Errorf("error creating destination path %s: %w", dest, err)
			}
		}

		// Acquire locks on all destination directories.
		lockPath := path.Join(dest, littbuilder.LockfileName)
		lock, err := util.NewFileLock(lockPath, fsync)
		if err != nil {
			return fmt.Errorf("failed to acquire lock on %s: %v", dest, err)
		}
		defer func() {
			_ = lock.Release()
		}()
	}

	directoriesGoingAway := make([]string, 0, len(sourceSet))
	for source := range sourceSet {
		// If the source directory is not in the destination set, it is going away.
		if _, ok := destinationSet[source]; !ok {
			directoriesGoingAway = append(directoriesGoingAway, source)
		}
	}

	var totalSegmentFileCount int64
	var segmentFileCount atomic.Int64
	if verbose {
		var err error
		totalSegmentFileCount, err = countSegmentFiles(directoriesGoingAway)
		if err != nil {
			return fmt.Errorf("failed to count segment files in sources %v: %w", sources, err)
		}
	}

	// For each directory that is going away, transfer its data to the new destination.
	for _, source := range directoriesGoingAway {
		err := transferDataInDirectory(
			source,
			destinations,
			deep,
			preserveOriginal,
			fsync,
			verbose,
			totalSegmentFileCount,
			&segmentFileCount)
		if err != nil {
			return fmt.Errorf("error transferring data from %s to %v: %w",
				source, destinations, err)
		}
	}

	if verbose {
		fmt.Println()
	}

	return nil
}

// Get a count of the segment files in the source directories.
func countSegmentFiles(sources []string) (int64, error) {
	count := int64(0)

	for _, source := range sources {
		exists, err := util.Exists(source)
		if err != nil {
			return 0, fmt.Errorf("failed to check if source directory %s exists: %w", source, err)
		}
		if !exists {
			continue
		}

		// Walk the file tree to find all files ending with .metadata, .keys, or .values.
		err = filepath.WalkDir(source, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error walking directory %s: %w", path, err)
			}

			if d.IsDir() {
				// Skip directories
				return nil
			}

			// Ignore "table.metadata" files, as they are not segment files.
			if d.Name() == disktable.TableMetadataFileName {
				return nil
			}

			// Check if the file is a segment file.
			extension := filepath.Ext(path)
			if extension == segment.MetadataFileExtension ||
				extension == segment.KeyFileExtension ||
				extension == segment.ValuesFileExtension {
				count++
			}

			return nil
		})

		if err != nil {
			return 0, fmt.Errorf("error counting segment files in source directories: %w", err)
		}
	}

	return count, nil
}

// transfers all data in a directory to the specified destinations.
func transferDataInDirectory(
	source string,
	destinations []string,
	deep bool,
	preserveOriginal bool,
	fsync bool,
	verbose bool,
	totalSegmentFileCount int64,
	segmentFileCount *atomic.Int64,
) error {
	exists, err := util.Exists(source)
	if err != nil {
		return fmt.Errorf("failed to check if source directory %s exists: %w", source, err)
	}
	if !exists {
		return nil
	}

	// Acquire a lock on the source directory.
	lockPath := path.Join(source, littbuilder.LockfileName)
	lock, err := util.NewFileLock(lockPath, fsync)
	if err != nil {
		return fmt.Errorf("failed to acquire lock on %s: %w", source, err)
	}
	defer func() {
		// double release is a no-op
		_ = lock.Release()
	}()

	// Transfer each table stored in this directory.
	children, err := os.ReadDir(source)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", source, err)
	}
	for _, child := range children {
		if !child.IsDir() {
			continue
		}

		err = transferDataInTable(
			source,
			child.Name(),
			destinations,
			deep,
			preserveOriginal,
			fsync,
			verbose,
			totalSegmentFileCount,
			segmentFileCount)
		if err != nil {
			return fmt.Errorf("error transferring data in table %s: %w", child.Name(), err)
		}
	}

	// Release the lock so we can delete the directory.
	err = lock.Release()
	if err != nil {
		return fmt.Errorf("failed to release lock on source directory %s: %w", source, err)
	}

	if !preserveOriginal {
		// Delete the directory.
		err = os.Remove(source)
		if err != nil {
			return fmt.Errorf("failed to remove source directory %s: %w", source, err)
		}
	}

	return nil
}

func transferDataInTable(
	source string,
	tableName string,
	destinations []string,
	deep bool,
	preserveOriginal bool,
	fsync bool,
	verbose bool,
	totalSegmentFileCount int64,
	segmentFileCount *atomic.Int64,
) error {

	err := createDestinationTableDirectories(destinations, tableName)
	if err != nil {
		return fmt.Errorf("failed to create destination table directories for table %s: %w", tableName, err)
	}

	err = transferKeymap(source, tableName, destinations, deep, preserveOriginal, fsync, verbose)
	if err != nil {
		return fmt.Errorf("failed to transfer keymap for table %s: %w", tableName, err)
	}

	err = transferTableMetadata(source, tableName, destinations, deep, preserveOriginal, fsync, verbose)
	if err != nil {
		return fmt.Errorf("failed to transfer table metadata for table %s: %w", tableName, err)
	}

	err = transferSegmentData(
		source,
		tableName,
		destinations,
		deep,
		preserveOriginal,
		fsync,
		verbose,
		totalSegmentFileCount,
		segmentFileCount)
	if err != nil {
		return fmt.Errorf("failed to transfer segment data for table %s: %w", tableName, err)
	}

	if !preserveOriginal {
		err = deleteSnapshotDirectory(source, tableName, true)
		if err != nil {
			return fmt.Errorf("failed to delete snapshot directory for table %s: %w", tableName, err)
		}

		err = deleteBoundaryFiles(source, tableName, verbose)
		if err != nil {
			return fmt.Errorf("failed to delete boundary files for table %s: %w", tableName, err)
		}

		// Once all data in a table is transferred, delete the table directory.
		sourceTableDir := filepath.Join(source, tableName)
		err = os.Remove(sourceTableDir)
		if err != nil {
			return fmt.Errorf("failed to remove table directory %s: %w", sourceTableDir, err)
		}
	}

	return nil
}

// deleteBoundaryFiles deletes the boundary files for a table. Only will be present if the source
// directory contains symlink snapshots.
func deleteBoundaryFiles(source string, tableName string, verbose bool) error {
	lowerBoundPath := path.Join(source, tableName, disktable.LowerBoundFileName)
	exists, err := util.Exists(lowerBoundPath)
	if err != nil {
		return fmt.Errorf("failed to check if lower bound file %s exists: %w", lowerBoundPath, err)
	}
	if exists {
		if verbose {
			fmt.Printf("Deleting lower bound file: %s\n", lowerBoundPath)
		}
		err = os.Remove(lowerBoundPath)
		if err != nil {
			return fmt.Errorf("failed to remove lower bound file %s: %w", lowerBoundPath, err)
		}
	}

	upperBoundPath := path.Join(source, tableName, disktable.UpperBoundFileName)
	exists, err = util.Exists(upperBoundPath)
	if err != nil {
		return fmt.Errorf("failed to check if upper bound file %s exists: %w", upperBoundPath, err)
	}
	if exists {
		if verbose {
			fmt.Printf("Deleting upper bound file: %s\n", upperBoundPath)
		}
		err = os.Remove(upperBoundPath)
		if err != nil {
			return fmt.Errorf("failed to remove upper bound file %s: %w", upperBoundPath, err)
		}
	}

	return nil
}

// delete the old snapshot directory for a table. This will be reconstructed the next time the DB is loaded.
func deleteSnapshotDirectory(source string, tableName string, verbose bool) error {
	snapshotDir := filepath.Join(source, tableName, segment.HardLinkDirectory)

	exists, err := util.Exists(snapshotDir)
	if err != nil {
		return fmt.Errorf("failed to check if snapshot directory %s exists: %w", snapshotDir, err)
	}
	if !exists {
		return nil
	}

	if verbose {
		fmt.Printf("Deleting snapshot directory: %s\n", snapshotDir)
	}

	err = os.RemoveAll(snapshotDir)
	if err != nil {
		return fmt.Errorf("failed to remove snapshot directory %s: %w", snapshotDir, err)
	}

	return nil
}

// In the destination directories, create directories for the tables (if they don't exist).
func createDestinationTableDirectories(destinations []string, tableName string) error {
	for _, destination := range destinations {
		destinationTableDir := filepath.Join(destination, tableName)
		exists, err := util.Exists(destinationTableDir)
		if err != nil {
			return fmt.Errorf("failed to check if destination table directory %s exists: %w",
				destinationTableDir, err)
		}
		if !exists {
			err = os.MkdirAll(destinationTableDir, 0755)
			if err != nil {
				return fmt.Errorf("failed to create destination table directory %s: %w",
					destinationTableDir, err)
			}
		}
	}

	return nil
}

// Transfer the keymap (if it is present in the source directory).
func transferKeymap(
	source string,
	tableName string,
	destinations []string,
	deep bool,
	preserveOriginal bool,
	fsync bool,
	verbose bool,
) error {

	sourceKeymapPath := filepath.Join(source, tableName, keymap.KeymapDirectoryName)
	exists, err := util.Exists(sourceKeymapPath)
	if err != nil {
		return fmt.Errorf("failed to check if keymap directory %s exists: %w", sourceKeymapPath, err)
	}
	if !exists {
		return nil
	}

	destination, err := determineDestination(sourceKeymapPath, destinations)
	if err != nil {
		return fmt.Errorf("failed to determine destination for keymap %s: %w", sourceKeymapPath, err)
	}

	destinationKeymapPath := filepath.Join(destination, tableName, keymap.KeymapDirectoryName)

	if verbose {
		text := fmt.Sprintf("Transferring table '%s' keymap", tableName)
		writer := bufio.NewWriter(os.Stdout)
		_, _ = fmt.Fprintf(writer, "\r%-100s", text)
		_ = writer.Flush()
	}

	err = util.RecursiveMove(sourceKeymapPath, destinationKeymapPath, deep, preserveOriginal, fsync)
	if err != nil {
		return fmt.Errorf("failed to copy keymap from %s to %s: %w",
			sourceKeymapPath, destinationKeymapPath, err)
	}

	return nil
}

// transfers data in the segments/ directory
func transferSegmentData(
	source string,
	tableName string,
	destinations []string,
	deep bool,
	preserveOriginal bool,
	fsync bool,
	verbose bool,
	totalSegmentFileCount int64,
	segmentFileCount *atomic.Int64,
) error {

	sourceTableDir := filepath.Join(source, tableName)

	sourceSegmentDir := filepath.Join(sourceTableDir, segment.SegmentDirectory)
	exists, err := util.Exists(sourceSegmentDir)
	if err != nil {
		return fmt.Errorf("failed to check if segment directory %s exists: %w", sourceSegmentDir, err)
	}
	if !exists {
		return nil
	}

	segmentFiles, err := os.ReadDir(sourceSegmentDir)
	if err != nil {
		return fmt.Errorf("failed to read segment directory %s: %w", sourceSegmentDir, err)
	}

	for _, segmentFile := range segmentFiles {
		segmentFilePath := filepath.Join(sourceSegmentDir, segmentFile.Name())
		err = transferSegmentFile(
			segmentFile.Name(),
			segmentFilePath,
			tableName,
			destinations,
			deep,
			preserveOriginal,
			fsync,
			verbose,
			totalSegmentFileCount,
			segmentFileCount)
		if err != nil {
			return fmt.Errorf("failed to transfer segment file %s for table %s: %w",
				segmentFilePath, tableName, err)
		}
	}

	if !preserveOriginal {
		// Now that we've copied the segment files, we can delete the original directory.
		err = os.Remove(sourceSegmentDir)
		if err != nil {
			return fmt.Errorf("failed to remove segment directory %s: %w", sourceSegmentDir, err)
		}
	}

	return nil
}

// Transfer a single segment file (i.e. *.metadata, *.keys, *.values).
func transferSegmentFile(
	segmentName string,
	segmentFilePath string,
	tableName string,
	destinations []string,
	deep bool,
	preserveOriginal bool,
	fsync bool,
	verbose bool,
	totalSegmentFileCount int64,
	segmentFileCount *atomic.Int64,
) error {

	destination, err := determineDestination(segmentFilePath, destinations)
	if err != nil {
		return fmt.Errorf("failed to determine destination for segment file %s: %w", segmentFilePath, err)
	}

	destinationSegmentPath := filepath.Join(destination, tableName, segment.SegmentDirectory, segmentName)

	if verbose {
		count := segmentFileCount.Add(1)
		text := fmt.Sprintf("Transferring Segment File %d/%d: %s from table '%s'",
			count, totalSegmentFileCount, filepath.Base(segmentFilePath), tableName)
		writer := bufio.NewWriter(os.Stdout)
		_, _ = fmt.Fprintf(writer, "\r%-100s", text)
		_ = writer.Flush()
	}

	err = util.RecursiveMove(segmentFilePath, destinationSegmentPath, deep, preserveOriginal, fsync)
	if err != nil {
		return fmt.Errorf("failed to copy segment file from %s to %s: %w",
			segmentFilePath, destinationSegmentPath, err)
	}

	return nil
}

// transfers the table metadata file, if it is present.
func transferTableMetadata(
	source string,
	tableName string,
	destinations []string,
	deep bool,
	preserveOriginal bool,
	fsync bool,
	verbose bool,
) error {

	sourceTableDir := filepath.Join(source, tableName)

	sourceMetadataPath := filepath.Join(sourceTableDir, disktable.TableMetadataFileName)
	exists, err := util.Exists(sourceMetadataPath)
	if err != nil {
		return fmt.Errorf("failed to check if table metadata file %s exists: %w", sourceMetadataPath, err)
	}

	if !exists {
		return nil
	}

	destination, err := determineDestination(sourceTableDir, destinations)
	if err != nil {
		return fmt.Errorf("failed to determine destination for table metadata %s: %w", sourceMetadataPath, err)
	}

	destinationMetadataPath := filepath.Join(destination, tableName, disktable.TableMetadataFileName)

	if verbose {
		text := fmt.Sprintf("Transferring table '%s' metadata", tableName)
		writer := bufio.NewWriter(os.Stdout)
		_, _ = fmt.Fprintf(writer, "\r%-100s", text)
		_ = writer.Flush()
	}

	err = util.RecursiveMove(sourceMetadataPath, destinationMetadataPath, deep, preserveOriginal, fsync)
	if err != nil {
		return fmt.Errorf("failed to copy table metadata from %s to %s: %w",
			sourceMetadataPath, destinationMetadataPath, err)
	}

	return nil
}

// Determines the location where a file should be transferred given a list of options.
// This function is deterministic. This is important! If a rebase is interrupted, the
// second attempt should always transfer the file to the same location as the first attempt.
func determineDestination(source string, destinations []string) (string, error) {
	hasher := fnv.New64a()
	_, err := hasher.Write([]byte(source))
	if err != nil {
		return "", fmt.Errorf("failed to hash source path %s: %w", source, err)
	}

	return destinations[hasher.Sum64()%uint64(len(destinations))], nil
}

// TODO make sure rebase works if the source contains symlinks
