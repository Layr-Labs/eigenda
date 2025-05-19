package benchmark

import (
	"encoding/binary"
	"fmt"
	"path"
	"time"
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

	return cohort, nil
}

// Seal marks that all key-value pairs in the cohort have been written to the database.
func (cohort *Cohort) Seal() error {
	return nil
}

// Path returns the file path of the cohort file.
func (cohort *Cohort) Path() string {
	return path.Join(cohort.parentDirectory, fmt.Sprintf("%d%s", cohort.cohortIndex, CohortFileExtension))
}

// serialize serializes the cohort to a byte array.
func (cohort *Cohort) serialize() []byte {
	// Data size:
	//  - cohortIndex (8 bytes)
	//  - lowIndex (8 bytes)
	//  - highIndex (8 bytes)
	//  - firstValueTimestamp (8 bytes)
	//  - allValuesWritten (1 byte)
	// Total: 33 bytes

	data := make([]byte, 33)
	binary.BigEndian.PutUint64(data[0:8], cohort.cohortIndex)
	binary.BigEndian.PutUint64(data[8:16], cohort.lowIndex)
	binary.BigEndian.PutUint64(data[16:24], cohort.highIndex)
	binary.BigEndian.PutUint64(data[24:32], uint64(cohort.firstValueTimestamp.Unix()))
	if cohort.allValuesWritten {
		data[32] = 1
	} else {
		data[32] = 0
	}

	return data
}
