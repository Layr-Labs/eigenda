package benchmark

import (
	"encoding/binary"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Layr-Labs/eigenda/litt/util"
)

// CohortFileExtension is the file extension used for cohort files.
const CohortFileExtension = ".cohort"

// CohortSwapFileExtension is the file extension used for cohort swap files. Used to atomically update cohort files.
const CohortSwapFileExtension = ".cohort.swap"

// A Cohort is a grouping of key-value pairs used for benchmarking.
//
// Key-value pairs each have unique indices, and knowing the index of a key-value pair allows the data to be
// regenerated deterministically. All key-value pairs in a cohort have sequential indices.
//
// The benchmarking engine records on disk several properties for each cohort:
// - the index numbers that belong to the cohort
// - whether or not all key-value pairs in the cohort have been written to the database
// - the timestamp when the first write occurred
//
// The goal of this record keeping is to allow the benchmarking engine to keep track of values that are eligible to
// be read. If it wants to perform X writes per second, it needs to have a way to get an eligible list of keys.
// Ideally it would just track all the keys in the database, but this is cost prohibitive. Tracking at the cohort
// granularity allows key tracking to be done performance.
//
// The benchmark engine will not attempt to read any value from a cohort unless it knows for certain that all values
// in the cohort have been written to the database. It will also not attempt to read from a cohort if the timestamp
// of the first value written is too close to its expected expiration time.
type Cohort struct {
	// The directory where the cohort file is stored.
	parentDirectory string

	// The unique ID of this cohort.
	cohortIndex uint64

	// The index (inclusive) of the first key-value pair in the cohort.
	lowIndex uint64

	// The index (exclusive) of the last key-value pair in the cohort.
	highIndex uint64

	// True iff all key-value pairs in the cohort have been written to the database.
	allValuesWritten bool

	// A timestamp that is guaranteed to come before the first value in the cohort is written to the database.
	firstValueTimestamp time.Time
}

// NewCohort creates a new cohort with the given index range.
func NewCohort(
	parentDirectory string,
	cohortIndex uint64,
	lowIndex uint64,
	highIndex uint64) (*Cohort, error) {

	cohort := &Cohort{
		parentDirectory:     parentDirectory,
		cohortIndex:         cohortIndex,
		lowIndex:            lowIndex,
		highIndex:           highIndex,
		allValuesWritten:    false,
		firstValueTimestamp: time.Time{},
	}

	err := cohort.Write()
	if err != nil {
		return nil, fmt.Errorf("failed to write cohort file: %w", err)
	}

	return cohort, nil
}

func LoadCohort(
	parentDirectory string,
	cohortIndex uint64) (*Cohort, error) {

	cohort := &Cohort{
		parentDirectory: parentDirectory,
		cohortIndex:     cohortIndex,
	}

	filePath := cohort.Path(false)
	exists, err := util.Exists(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to check if cohort file exists: %w", err)
	}

	if !exists {
		return nil, fmt.Errorf("cohort file does not exist: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open cohort file: %w", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cohort file: %w", err)
	}

	err = cohort.deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize cohort file: %w", err)
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close cohort file: %w", err)
	}

	return cohort, nil
}

// MarkComplete marks that all key-value pairs in the cohort have been written to the database. Once done,
// all key-value pairs in the cohort become safe to read, so long as the cohort has not yet expired.
func (c *Cohort) MarkComplete() error {
	c.allValuesWritten = true
	err := c.Write()
	if err != nil {
		return fmt.Errorf("failed to mark cohort complete: %w", err)
	}
	return nil
}

// Path returns the file path of the cohort file.
func (c *Cohort) Path(swap bool) string {

	var extension string
	if swap {
		extension = CohortSwapFileExtension
	} else {
		extension = CohortFileExtension
	}

	return path.Join(c.parentDirectory, fmt.Sprintf("%d%s", c.cohortIndex, extension))
}

func (c *Cohort) Write() error {
	swapPath := c.Path(true)
	targetPath := c.Path(false)

	swapFile, err := os.Create(swapPath)
	if err != nil {
		return fmt.Errorf("failed to create swap file: %w", err)
	}

	bytes := c.serialize()
	_, err = swapFile.Write(bytes)
	if err != nil {
		return fmt.Errorf("failed to write to swap file: %w", err)
	}

	err = swapFile.Sync()
	if err != nil {
		return fmt.Errorf("failed to sync swap file: %w", err)
	}

	err = swapFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close swap file: %w", err)
	}

	err = os.Rename(swapPath, targetPath)
	if err != nil {
		return fmt.Errorf("failed to rename swap file: %w", err)
	}

	return nil
}

// serialize serializes the cohort to a byte array.
func (c *Cohort) serialize() []byte {
	// Data size:
	//  - cohortIndex (8 bytes)
	//  - lowIndex (8 bytes)
	//  - highIndex (8 bytes)
	//  - firstValueTimestamp (8 bytes)
	//  - allValuesWritten (1 byte)
	// Total: 33 bytes

	data := make([]byte, 33)
	binary.BigEndian.PutUint64(data[0:8], c.cohortIndex)
	binary.BigEndian.PutUint64(data[8:16], c.lowIndex)
	binary.BigEndian.PutUint64(data[16:24], c.highIndex)
	binary.BigEndian.PutUint64(data[24:32], uint64(c.firstValueTimestamp.Unix()))
	if c.allValuesWritten {
		data[32] = 1
	} else {
		data[32] = 0
	}

	return data
}

func (c *Cohort) deserialize(data []byte) error {
	if len(data) != 33 {
		return fmt.Errorf("invalid data length: %d", len(data))
	}

	cohortIndex := binary.BigEndian.Uint64(data[0:8])
	if cohortIndex != c.cohortIndex {
		return fmt.Errorf("cohort index mismatch: %d != %d", cohortIndex, c.cohortIndex)
	}

	c.lowIndex = binary.BigEndian.Uint64(data[8:16])
	c.highIndex = binary.BigEndian.Uint64(data[16:24])
	if c.lowIndex >= c.highIndex {
		return fmt.Errorf("invalid index range: %d >= %d", c.lowIndex, c.highIndex)
	}

	c.firstValueTimestamp = time.Unix(int64(binary.BigEndian.Uint64(data[24:32])), 0)
	c.allValuesWritten = data[32] == 1

	return nil
}
