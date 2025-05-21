package benchmark

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/docker/go-units"
)

// WriteInfo contains information needed to perform a write operation.
type WriteInfo struct {
	// The index of the key to write.
	Index uint64
	// The key to write.
	Key []byte
	// The value to write.
	Value []byte
}

// ReadInfo contains information needed to perform a read operation.
type ReadInfo struct {
	// The key to read.
	Key []byte
	// The value we expect to read.
	Value []byte
}

// DataTracker is responsible for tracking key-value pairs that have been written to the database, and for generating
// new key-value pairs to be written.
type DataTracker struct {
	ctx    context.Context
	cancel context.CancelFunc

	// A source of randomness.
	rand *rand.Rand

	// The configuration for the benchmark.
	config *BenchmarkConfig

	// The directory where cohort files are stored.
	cohortDirectory string

	// A map from cohort index to information about the cohort.
	cohorts map[uint64]*Cohort

	// The cohort that is currently being used to generate keys for writing.
	activeCohort *Cohort

	// A set of cohorts that have been completely written to the database (i.e. cohorts that are safe to read).
	completeCohortSet map[uint64]struct{}

	// The index of the oldest cohort being tracked.
	lowestCohortIndex uint64

	// The index of the newest cohort being tracked.
	highestCohortIndex uint64

	// The index of the oldest cohort currently being written to. All older cohorts are either complete or abandoned.
	highestCohortBeingWrittenIndex uint64

	// A channel containing keys-value pairs that are ready to be written.
	writeInfoChan chan *WriteInfo

	// A channel containing keys that are ready to be read.
	readInfoChan chan *ReadInfo

	// A channel containing information about keys that have been written to the database.
	writtenKeyIndicesChan chan uint64

	// Responsible for producing "random" data for key-value pairs.
	generator *DataGenerator

	// The TTL minus a safety margin. Cohorts are considered to be expired if keys in them are older than this.
	safeTTL time.Duration

	// The size of the values in bytes for new cohorts.
	valueSize uint64
}

// NewDataTracker creates a new DataTracker instance, loading all relevant cohorts from disk.
func NewDataTracker(ctx context.Context, config *BenchmarkConfig) (*DataTracker, error) {
	cohortDirectory := path.Join(config.MetadataDirectory, "cohorts")

	lowestCohortIndex, highestCohortIndex, cohorts, err := gatherCohorts(cohortDirectory)
	if err != nil {
		return nil, fmt.Errorf("failed to gather cohorts: %w", err)
	}

	completeCohortSet := make(map[uint64]struct{})
	for i := lowestCohortIndex; i <= highestCohortIndex; i++ {
		if cohorts[i].IsComplete() {
			completeCohortSet[i] = struct{}{}
		}
	}

	valueSize := uint64(config.ValueSizeMB * float64(units.MiB))

	// Create an initial active cohort.
	var activeCohort *Cohort
	if len(cohorts) == 0 {
		// Starting fresh, create a new cohort starting from key index 0.
		activeCohort, err = NewCohort(
			cohortDirectory,
			0,
			0,
			config.CohortSize,
			valueSize)
		if err != nil {
			return nil, fmt.Errorf("failed to create genesis cohort: %w", err)
		}
	} else {
		activeCohort, err = cohorts[highestCohortIndex].NextCohort(config.CohortSize, valueSize)
		if err != nil {
			return nil, fmt.Errorf("failed to create next cohort: %w", err)
		}
	}
	highestCohortIndex = activeCohort.CohortIndex()
	cohorts[highestCohortIndex] = activeCohort

	writeInfoChan := make(chan *WriteInfo, config.WriteInfoChanelSize)
	readInfoChan := make(chan *ReadInfo, config.ReadInfoChanelSize)
	writtenKeyIndicesChan := make(chan uint64, 64)

	ttl := time.Duration(config.TTLHours * float64(time.Hour))
	safetyMargin := time.Duration(config.ReadSafetyMarginMinutes * float64(time.Minute))
	safeTTL := ttl - safetyMargin

	ctx, cancel := context.WithCancel(ctx)

	tracker := &DataTracker{
		ctx:                            ctx,
		cancel:                         cancel,
		rand:                           rand.New(rand.NewSource(time.Now().UnixNano())),
		config:                         config,
		cohortDirectory:                cohortDirectory,
		cohorts:                        cohorts,
		completeCohortSet:              completeCohortSet,
		writeInfoChan:                  writeInfoChan,
		readInfoChan:                   readInfoChan,
		writtenKeyIndicesChan:          writtenKeyIndicesChan,
		activeCohort:                   activeCohort,
		lowestCohortIndex:              lowestCohortIndex,
		highestCohortIndex:             highestCohortIndex,
		highestCohortBeingWrittenIndex: highestCohortIndex, // only the active cohort initially being written
		safeTTL:                        safeTTL,
		valueSize:                      valueSize,
	}

	go tracker.dataGenerator()

	return tracker, nil
}

func gatherCohorts(path string) (
	lowestCohortIndex uint64,
	highestCohortIndex uint64,
	cohorts map[uint64]*Cohort,
	err error) {

	cohorts = make(map[uint64]*Cohort)

	// walk over files in path
	// for each file, check if it is a cohort file
	// if it is, load the cohort and add it to the map
	// if it is not, ignore it
	files, err := os.ReadDir(path)
	if err != nil {
		return 0,
			0,
			nil,
			fmt.Errorf("failed to read directory: %w", err)
	}

	lowestCohortIndex = math.MaxUint64
	highestCohortIndex = 0

	for _, file := range files {
		fileName := file.Name()

		if strings.HasSuffix(fileName, CohortFileExtension) {
			cohort, err := LoadCohort(path)
			if err != nil {
				return 0,
					0,
					nil,
					fmt.Errorf("failed to load cohort: %w", err)
			}
			cohorts[cohort.CohortIndex()] = cohort

			if cohort.CohortIndex() < lowestCohortIndex {
				lowestCohortIndex = cohort.CohortIndex()
			}
			if cohort.cohortIndex > highestCohortIndex {
				highestCohortIndex = cohort.cohortIndex
			}
		} else if strings.HasSuffix(fileName, CohortSwapFileExtension) {
			// Delete any swap files discovered
			err = os.Remove(fileName)
			if err != nil {
				return 0,
					0,
					nil,
					fmt.Errorf("failed to delete swap file: %w", err)
			}
		}
	}

	if len(cohorts) == 0 {
		// Special case, no cohorts found.
		return 0, 0, cohorts, nil
	}

	return lowestCohortIndex, highestCohortIndex, cohorts, nil
}

// GetWriteInfo returns information required to perform a write operation. It returns the key index (which is needed to
// call MarkHighestIndexWritten()), the key, and the value. Data is generated on background goroutines in order to
// make this method very fast. Will not block as long as data can be generated in the background fast enough.
// May return nil if the context is cancelled.
func (t *DataTracker) GetWriteInfo() *WriteInfo {
	select {
	case info := <-t.writeInfoChan:
		return info
	case <-t.ctx.Done():
		return nil
	}
}

// MarkHighestIndexWritten marks the given index as having been written. It is assumed that writes happen in order,
// meaning that calling MarkHighestIndexWritten(X) implies that index X-1 has also been written (as long as X-1 is not
// an index that was allocated before the benchmark most recently restarted).
func (t *DataTracker) MarkHighestIndexWritten(index uint64) {
	select {
	case t.writtenKeyIndicesChan <- index:
		return
	case <-t.ctx.Done():
		return
	}
}

// GetReadInfo returns information required to perform a read operation. If this returns nil, then there are currently
// no keys that are safe to read.
func (t *DataTracker) GetReadInfo() *ReadInfo {
	select {
	case info := <-t.readInfoChan:
		return info
	case <-t.ctx.Done():
		return nil
	}
}

// Close stops the key manager's background tasks.
func (t *DataTracker) Close() {
	t.cancel()
}

// dataGenerator is responsible for generating data in the background.
func (t *DataTracker) dataGenerator() {
	ticker := time.NewTicker(time.Duration(t.config.CohortGCPeriodSeconds * float64(time.Second)))
	defer ticker.Stop()

	nextWriteInfo := t.generateNextWriteInfo()
	nextReadInfo := t.generateNextReadInfo()

	for {
		select {

		case <-t.ctx.Done():
			// abort when context is cancelled
			return

		case keyIndex := <-t.writtenKeyIndicesChan:
			// track keys that have been written so that we can read them in the future
			t.handleWrittenKey(keyIndex)

		case t.writeInfoChan <- nextWriteInfo:
			// prepare a value to be eventually written
			nextWriteInfo = t.generateNextWriteInfo()

		case t.readInfoChan <- nextReadInfo:
			// prepare a value to be eventually read
			nextReadInfo = t.generateNextReadInfo()

		case <-ticker.C:
			// perform garbage collection on cohorts
			t.DoCohortGC()
		}
	}
}

// handleWrittenKey handles a key that has been written to the database.
func (t *DataTracker) handleWrittenKey(keyIndex uint64) {
	// Iterate over cohorts currently being written. If any cohort has all keys less than or equal to
	// this keyIndex, then mark that cohort as complete.

	for i := t.highestCohortBeingWrittenIndex; i <= t.highestCohortIndex; i++ {
		cohort := t.cohorts[i]

		if cohort.HighKeyIndex() > keyIndex {
			// Once we find the first cohort without all keys written, we can stop checking.
			break
		}

		// All keys in this cohort have been written.
		err := cohort.MarkComplete()
		if err != nil {
			panic(fmt.Sprintf("failed to mark cohort complete: %v", err)) // TODO not clean
		}
		t.completeCohortSet[keyIndex] = struct{}{}
	}
}

// generateNextWriteInfo generates the next write info to be placed into the writeInfoChan.
func (t *DataTracker) generateNextWriteInfo() *WriteInfo {
	var err error

	if t.activeCohort.IsExhausted() {
		t.activeCohort, err = t.cohorts[t.highestCohortIndex].NextCohort(t.config.CohortSize, t.valueSize)
		if err != nil {
			panic(fmt.Sprintf("failed to generate next cohort for highest cohort: %v", err)) // TODO not clean
		}
		t.highestCohortIndex = t.activeCohort.CohortIndex()
	}

	keyIndex, err := t.activeCohort.GetKeyIndexForWriting()
	if err != nil {
		panic(fmt.Sprintf("failed to get key index for writing: %v", err)) // TODO not clean
	}

	return &WriteInfo{
		Index: keyIndex,
		Key:   t.generator.Key(keyIndex),
		Value: t.generator.Value(keyIndex, t.activeCohort.valueSize),
	}
}

// generateNextReadInfo generates the next read info to be placed into the readInfoChan.
func (t *DataTracker) generateNextReadInfo() *ReadInfo {
	var cohortIndexToRead uint64
	for cohortIndexToRead = range t.completeCohortSet {
		// map iteration is random in golang, so this will yield a random complete cohort.
		break
	}
	cohortToRead, ok := t.cohorts[cohortIndexToRead]
	if !ok {
		// This cohort has been removed from the set of complete cohorts.
		return nil
	}

	keyIndex, err := cohortToRead.GetKeyIndexForReading(t.rand)
	if err != nil {
		panic(fmt.Sprintf("failed to get key index for reading: %v", err)) // TODO not clean
	}

	return &ReadInfo{
		Key:   t.generator.Key(keyIndex),
		Value: t.generator.Value(keyIndex, cohortToRead.ValueSize()),
	}
}

// DoCohortGC performs garbage collection on the cohorts, removing cohorts with entries that are nearing expiration.
func (t *DataTracker) DoCohortGC() {
	now := time.Now()

	for i := t.lowestCohortIndex; i <= t.highestCohortIndex; i++ {
		cohort := t.cohorts[i]

		if cohort.IsExpired(now, t.safeTTL) {
			err := cohort.Delete()
			if err != nil {
				panic(fmt.Sprintf("failed to delete expired cohort: %v", err)) // TODO not clean
			}
			t.lowestCohortIndex++
		}
	}

	if len(t.cohorts) == 0 {
		// Edge case: we've been writing data slow enough that the active cohort has expired.
		// Create a new active cohort.
		activeCohort, err := t.activeCohort.NextCohort(t.config.CohortSize, t.valueSize)
		if err != nil {
			panic(fmt.Sprintf("failed to create new active cohort: %v", err)) // TODO not clean
		}

		t.activeCohort = activeCohort
		t.highestCohortIndex = activeCohort.CohortIndex()
		t.highestCohortBeingWrittenIndex = activeCohort.CohortIndex()
		t.cohorts[activeCohort.CohortIndex()] = activeCohort
	}
}
