package benchmark

import (
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/litt/util"
	"github.com/stretchr/testify/require"
)

func TestCohortSerialization(t *testing.T) {
	rand := random.NewTestRandom()
	testDirectory := t.TempDir()

	cohortIndex := rand.Uint64()
	lowIndex := rand.Uint64Range(1, 1000)
	highIndex := rand.Uint64Range(1000, 2000)
	cohort, err := NewCohort(
		testDirectory,
		cohortIndex,
		lowIndex,
		highIndex)
	require.NoError(t, err)

	require.Equal(t, cohort.CohortIndex(), cohortIndex)
	require.Equal(t, cohort.LowKeyIndex(), lowIndex)
	require.Equal(t, cohort.HighKeyIndex(), highIndex)
	require.Equal(t, cohort.IsComplete(), false)

	// Check if the cohort file exists
	filePath := cohort.Path(false)
	exists, err := util.Exists(filePath)
	require.NoError(t, err)
	require.True(t, exists)

	// Initialize a copy cohort from the file
	loadedCohort, err := LoadCohort(testDirectory, cohortIndex)
	require.NoError(t, err)
	require.Equal(t, loadedCohort.CohortIndex(), cohortIndex)
	require.Equal(t, loadedCohort.LowKeyIndex(), lowIndex)
	require.Equal(t, loadedCohort.HighKeyIndex(), highIndex)
	require.Equal(t, loadedCohort.IsComplete(), false)

	// Mark the cohort as written
	loadedCohort.allValuesWritten = true
	require.NoError(t, err)
	require.True(t, loadedCohort.IsComplete())
	err = loadedCohort.Write()
	require.NoError(t, err)

	// Load the cohort again.
	loadedCohort, err = LoadCohort(testDirectory, cohortIndex)
	require.NoError(t, err)
	require.Equal(t, loadedCohort.CohortIndex(), cohortIndex)
	require.Equal(t, loadedCohort.LowKeyIndex(), lowIndex)
	require.Equal(t, loadedCohort.HighKeyIndex(), highIndex)
	require.Equal(t, loadedCohort.IsComplete(), true)
}

func TestStandardCohortLifecycle(t *testing.T) {
	rand := random.NewTestRandom()
	testDirectory := t.TempDir()

	cohortIndex := rand.Uint64()
	lowIndex := rand.Uint64Range(1, 1000)
	highIndex := rand.Uint64Range(1000, 2000)
	cohort, err := NewCohort(
		testDirectory,
		cohortIndex,
		lowIndex,
		highIndex)
	require.NoError(t, err)

	require.Equal(t, cohort.CohortIndex(), cohortIndex)
	require.Equal(t, cohort.LowKeyIndex(), lowIndex)
	require.Equal(t, cohort.HighKeyIndex(), highIndex)
	require.Equal(t, cohort.IsComplete(), false)

	// Extract all keys from the cohort.
	for i := lowIndex; i <= highIndex; i++ {
		key, err := cohort.GetKeyIndexForWriting()
		require.NoError(t, err)
		require.Equal(t, i, key)

		shouldBeExhausted := i == highIndex
		require.Equal(t, shouldBeExhausted, cohort.IsExhausted())

		if i < highIndex {
			// Attempting to mark as complete now should fail.
			err = cohort.MarkComplete()
			require.Error(t, err)
		}
		require.Equal(t, false, cohort.IsComplete())

		// Attempting to get a key for reading should fail.
		_, err = cohort.GetKeyIndexForReading(rand.Rand)
		require.Error(t, err)
	}

	// Attempting to allocate another key for writing should fail.
	_, err = cohort.GetKeyIndexForWriting()
	require.Error(t, err)

	// We can now mark the cohort as complete.
	err = cohort.MarkComplete()
	require.NoError(t, err)
	require.Equal(t, true, cohort.IsComplete())

	// We can now get keys for reading.
	for i := 0; i < 100; i++ {
		key, err := cohort.GetKeyIndexForReading(rand.Rand)
		require.NoError(t, err)
		require.GreaterOrEqual(t, key, lowIndex)
		require.LessOrEqual(t, key, highIndex)
	}

	// Marking complete again should fail.
	err = cohort.MarkComplete()
	require.Error(t, err)
}

func TestIncompleteCohortAllKeysExtractedLifecycle(t *testing.T) {
	rand := random.NewTestRandom()
	testDirectory := t.TempDir()

	cohortIndex := rand.Uint64()
	lowIndex := rand.Uint64Range(1, 1000)
	highIndex := rand.Uint64Range(1000, 2000)
	cohort, err := NewCohort(
		testDirectory,
		cohortIndex,
		lowIndex,
		highIndex)
	require.NoError(t, err)

	require.Equal(t, cohort.CohortIndex(), cohortIndex)
	require.Equal(t, cohort.LowKeyIndex(), lowIndex)
	require.Equal(t, cohort.HighKeyIndex(), highIndex)
	require.Equal(t, cohort.IsComplete(), false)

	// Extract all keys from the cohort.
	for i := lowIndex; i <= highIndex; i++ {
		key, err := cohort.GetKeyIndexForWriting()
		require.NoError(t, err)
		require.Equal(t, i, key)

		shouldBeExhausted := i == highIndex
		require.Equal(t, shouldBeExhausted, cohort.IsExhausted())

		if i < highIndex {
			// Attempting to mark as complete now should fail.
			err = cohort.MarkComplete()
			require.Error(t, err)
		}
		require.Equal(t, false, cohort.IsComplete())

		// Attempting to get a key for reading should fail.
		_, err = cohort.GetKeyIndexForReading(rand.Rand)
		require.Error(t, err)
	}

	// Simulate a benchmark restart by reloading the cohort from disk.
	loadedCohort, err := LoadCohort(testDirectory, cohortIndex)
	require.NoError(t, err)

	require.Equal(t, loadedCohort.CohortIndex(), cohortIndex)
	require.False(t, loadedCohort.IsComplete())

	// Attempting to allocate another key for writing should fail.
	_, err = loadedCohort.GetKeyIndexForWriting()
	require.Error(t, err)

	// Attempting to get a key for reading should fail.
	_, err = loadedCohort.GetKeyIndexForReading(rand.Rand)
	require.Error(t, err)

	// We shouldn't be able to mark the cohort as complete.
	err = loadedCohort.MarkComplete()
	require.Error(t, err)
}

func TestIncompleteCohortSomeKeysExtractedLifecycle(t *testing.T) {
	rand := random.NewTestRandom()
	testDirectory := t.TempDir()

	cohortIndex := rand.Uint64()
	lowIndex := rand.Uint64Range(1, 1000)
	highIndex := rand.Uint64Range(1000, 2000)
	cohort, err := NewCohort(
		testDirectory,
		cohortIndex,
		lowIndex,
		highIndex)
	require.NoError(t, err)

	require.Equal(t, cohort.CohortIndex(), cohortIndex)
	require.Equal(t, cohort.LowKeyIndex(), lowIndex)
	require.Equal(t, cohort.HighKeyIndex(), highIndex)
	require.Equal(t, cohort.IsComplete(), false)

	// Extract all keys from the cohort.
	for i := lowIndex; i <= (lowIndex+highIndex)/2; i++ {
		key, err := cohort.GetKeyIndexForWriting()
		require.NoError(t, err)
		require.Equal(t, i, key)

		shouldBeExhausted := i == highIndex
		require.Equal(t, shouldBeExhausted, cohort.IsExhausted())
		
		// Attempting to mark as complete now should fail.
		err = cohort.MarkComplete()
		require.Error(t, err)
		require.Equal(t, false, cohort.IsComplete())

		// Attempting to get a key for reading should fail.
		_, err = cohort.GetKeyIndexForReading(rand.Rand)
		require.Error(t, err)
	}

	// Simulate a benchmark restart by reloading the cohort from disk.
	loadedCohort, err := LoadCohort(testDirectory, cohortIndex)
	require.NoError(t, err)

	require.Equal(t, loadedCohort.CohortIndex(), cohortIndex)
	require.False(t, loadedCohort.IsComplete())

	// Attempting to allocate another key for writing should fail.
	_, err = loadedCohort.GetKeyIndexForWriting()
	require.Error(t, err)

	// Attempting to get a key for reading should fail.
	_, err = loadedCohort.GetKeyIndexForReading(rand.Rand)
	require.Error(t, err)

	// We shouldn't be able to mark the cohort as complete.
	err = loadedCohort.MarkComplete()
	require.Error(t, err)
}
