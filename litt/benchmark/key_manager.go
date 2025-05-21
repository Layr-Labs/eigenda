package benchmark

// KeyManager is responsible for tracking key-value pairs that have been written to the database.
//
// Public methods on this object are NOT thread safe, and should not be called from multiple threads concurrently.
type KeyManager struct {
	// The directory where cohort files are stored.
	cohortDirectory string

	// A map from cohort index to information about the cohort.
	cohorts map[uint64]*Cohort

	// The index of the oldest cohort being tracked.
	lowestCohortIndex uint64

	// The index of the newest cohort being tracked.
	highestCohortIndex uint64
}

// NewKeyManager creates a new key manager.
func NewKeyManager() (*KeyManager, error) {
	// TODO
	return nil, nil
}

// GetWritableIndex returns a key-value index that is safe to write to (i.e. it has never been written before).
// Each key-value pair written by this benchmark has a unique index, and knowing the index permits the key and
// value to be deterministically generated.
func (m *KeyManager) GetWritableIndex() uint64 {
	// TODO
	return 0
}

// MarkHighestIndexWritten marks the given index as having been written. It is assumed that writes happen in index,
// meaning that calling MarkHighestIndexWritten(X) implies that index X-1 has also been written.
func (m *KeyManager) MarkHighestIndexWritten(index uint64) error {
	// TODO
	return nil
}

// GetReadableIndex returns a random key-value index that is safe to read from (i.e. it has been written before).
func (m *KeyManager) GetReadableIndex(maxIndex uint64) uint64 {
	// TODO
	return 0
}
