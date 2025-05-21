package benchmark

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/Layr-Labs/eigenda/litt/util"
)

// CohortFileExtension is the file extension used for cohort files.
const CohortFileExtension = ".cohort"

// CohortSwapFileExtension is the file extension used for cohort swap files. Used to atomically update cohort files.
const CohortSwapFileExtension = ".cohort.swap"

/* The lifecycle of a cohort:

    +-----+     +-----------+     +----------+
    | new | --> | exhausted | --> | complete |
    +-----+     +-----------+     +----------+
       |              |
       v              |
    +------------+    |
    | incomplete | <--|
    +------------+

- new: the cohort was just created and is currently being used to supply keys for writing.
- exhausted: all keys in the cohort have taken to be written, but the DB may not have ingested them all yet.
- complete: all keys in the cohort have been written to the DB and are safe to read.
- incomplete: before becoming complete, the benchmark was restarted. It will never be thread safe to read or write
              any keys in this cohort.
*/

// A Cohort is a grouping of key-value pairs used for benchmarking.
//
// Key-value pairs each have unique indices, and knowing the index of a key-value pair allows the data to be
// regenerated deterministically. All key-value pairs in a cohort have sequential indices.
type Cohort struct {
	// The directory where the cohort file is stored.
	parentDirectory string

	// The unique ID of this cohort.
	cohortIndex uint64

	// The index of the first key-value pair in the cohort.
	lowKeyIndex uint64

	// The index of the last key-value pair in the cohort.
	highKeyIndex uint64

	// The next available index to be written. Only relevant for a new cohort that is currently being written to
	// the DB. This value is undefined for cohorts that have been completely written or loaded from disk. This value
	// is NOT serialized to disk.
	nextKeyIndex uint64

	// True iff all key-value pairs in the cohort have been written to the database.
	allValuesWritten bool

	// A timestamp that is guaranteed to come before the first value in the cohort is written to the database.
	firstValueTimestamp time.Time

	// True iff the cohort has been loaded from disk. This value is NOT serialized to disk.
	loadedFromDisk bool
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
		lowKeyIndex:         lowIndex,
		highKeyIndex:        highIndex,
		nextKeyIndex:        lowIndex,
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
		loadedFromDisk:  true,
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

// CohortIndex returns the index of the cohort.
func (c *Cohort) CohortIndex() uint64 {
	return c.cohortIndex
}

// LowKeyIndex returns the index of the first key in the cohort.
func (c *Cohort) LowKeyIndex() uint64 {
	return c.lowKeyIndex
}

// HighKeyIndex returns the index of the last key in the cohort.
func (c *Cohort) HighKeyIndex() uint64 {
	return c.highKeyIndex
}

// FirstValueTimestamp returns the timestamp of the first value in the cohort.
func (c *Cohort) FirstValueTimestamp() time.Time {
	return c.firstValueTimestamp
}

// IsComplete returns true if all key-value pairs in the cohort have been written to the database. Only complete
// cohorts are safe to read from.
func (c *Cohort) IsComplete() bool {
	return c.allValuesWritten
}

// IsExhausted returns true if the cohort has been exhausted, i.e. it has produced all keys for writing that it is
// capable of producing. Once exhausted, a cohort should be marked as completed once all key-value pairs have been
// written to the database, thus making all keys in the cohort safe to read.
func (c *Cohort) IsExhausted() bool {
	return c.nextKeyIndex > c.highKeyIndex
}

// GetKeyIndexForWriting gets the next key to be written to the database.
func (c *Cohort) GetKeyIndexForWriting() (uint64, error) {
	if c.loadedFromDisk {
		return 0, fmt.Errorf("cannot allocate key for writing: cohort has been loaded from disk")
	}
	if c.allValuesWritten {
		return 0, fmt.Errorf("cannot allocate key for writing: cohort is already complete")
	}
	if c.nextKeyIndex > c.highKeyIndex {
		return 0, fmt.Errorf("cannot allocate key for writing: cohort is exhausted")
	}

	key := c.nextKeyIndex
	c.nextKeyIndex++

	return key, nil
}

// GetKeyIndexForReading gets a random key from the cohort that is safe to read. This function should only be called
// after the cohort has been marked as complete.
func (c *Cohort) GetKeyIndexForReading(rand *rand.Rand) (uint64, error) {
	if !c.allValuesWritten {
		return 0, fmt.Errorf("cannot allocate key for reading: cohort is not complete")
	}

	choice := (rand.Uint64() % (c.highKeyIndex - c.lowKeyIndex + 1)) + c.lowKeyIndex

	// sanity check
	if choice < c.lowKeyIndex || choice > c.highKeyIndex {
		return 0, fmt.Errorf("invalid choice: %d not in range [%d, %d]", choice, c.lowKeyIndex, c.highKeyIndex)
	}

	return choice, nil
}

// MarkComplete marks that all key-value pairs in the cohort have been written to the database. Once done,
// all key-value pairs in the cohort become safe to read, so long as the cohort has not yet expired. A cohort
// is said to have expired when it is possible that at least one key in the cohort may be deleted from the DB
// due to the TTL.
func (c *Cohort) MarkComplete() error {
	if c.allValuesWritten == true {
		return fmt.Errorf("cannot mark cohort complete: cohort is already complete")
	}
	if c.loadedFromDisk {
		return fmt.Errorf("cannot mark cohort complete: cohort has been loaded from disk")
	}
	if c.nextKeyIndex <= c.highKeyIndex {
		return fmt.Errorf("cannot mark cohort complete: cohort is not exhausted")
	}

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
	//  - lowKeyIndex (8 bytes)
	//  - highKeyIndex (8 bytes)
	//  - firstValueTimestamp (8 bytes)
	//  - allValuesWritten (1 byte)
	// Total: 33 bytes

	data := make([]byte, 33)
	binary.BigEndian.PutUint64(data[0:8], c.cohortIndex)
	binary.BigEndian.PutUint64(data[8:16], c.lowKeyIndex)
	binary.BigEndian.PutUint64(data[16:24], c.highKeyIndex)
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

	c.lowKeyIndex = binary.BigEndian.Uint64(data[8:16])
	c.highKeyIndex = binary.BigEndian.Uint64(data[16:24])
	if c.lowKeyIndex >= c.highKeyIndex {
		return fmt.Errorf("invalid index range: %d >= %d", c.lowKeyIndex, c.highKeyIndex)
	}

	c.firstValueTimestamp = time.Unix(int64(binary.BigEndian.Uint64(data[24:32])), 0)
	c.allValuesWritten = data[32] == 1

	return nil
}
