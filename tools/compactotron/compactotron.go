package main

import (
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/docker/go-units"
	"github.com/syndtr/goleveldb/leveldb"
)

// The maximum size of a batch to write to LevelDB.
const maxBatchSize = 100 * units.MiB

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: compactotron <source_path> <destination_path>")
		os.Exit(1)
	}

	sourcePath := os.Args[1]
	destinationPath := os.Args[2]

	err := CompactLevelDB(sourcePath, destinationPath)
	if err != nil {
		fmt.Printf("Error compacting LevelDB: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Compaction completed successfully.")
}

// Compacts LevelDB database at the given source path and writes the compacted data to the destination path.
func CompactLevelDB(source string, destination string) error {
	var err error

	source, err = util.SanitizePath(source)
	if err != nil {
		return fmt.Errorf("failed to sanitize source path: %w", err)
	}

	destination, err = util.SanitizePath(destination)
	if err != nil {
		return fmt.Errorf("failed to sanitize destination path: %w", err)
	}

	if source == destination {
		return fmt.Errorf("source and destination paths are both the same: %s", source)
	}

	err = util.ErrIfNotExists(source)
	if err != nil {
		return fmt.Errorf("source path does not exist: %w", err)
	}

	err = util.ErrIfExists(destination)
	if err != nil {
		return fmt.Errorf("destination path already exists: %w", err)
	}

	err = os.MkdirAll(destination, 0755)
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	sourceDB, err := leveldb.OpenFile(source, nil)
	if err != nil {
		return fmt.Errorf("failed to open source LevelDB: %w", err)
	}
	defer func() {
		_ = sourceDB.Close()
	}()

	destinationDB, err := leveldb.OpenFile(destination, nil)
	if err != nil {
		return fmt.Errorf("failed to open destination LevelDB: %w", err)
	}
	defer func() {
		_ = destinationDB.Close()
	}()

	iterator := sourceDB.NewIterator(nil, nil)
	defer iterator.Release()

	batch := new(leveldb.Batch)
	batchSize := 0
	totalSize := 0

	for iterator.Next() {
		key := iterator.Key()
		value := iterator.Value()
		batchSize += len(key) + len(value)

		batch.Put(key, value)

		if batchSize >= maxBatchSize {
			err = destinationDB.Write(batch, nil)
			if err != nil {
				return fmt.Errorf("failed to write batch to destination LevelDB: %w", err)
			}

			totalSize += batchSize
			fmt.Printf("%s copied so far\n", common.PrettyPrintBytes(uint64(totalSize)))

			batch = new(leveldb.Batch)
			batchSize = 0
		}
	}

	if batchSize > 0 {
		err = destinationDB.Write(batch, nil)
		if err != nil {
			return fmt.Errorf("failed to write final batch to destination LevelDB: %w", err)
		}

		totalSize += batchSize
		fmt.Printf("%s copied in total\n", common.PrettyPrintBytes(uint64(totalSize)))
	}

	return nil
}
