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

	require.Equal(t, cohort.cohortIndex, cohortIndex)
	require.Equal(t, cohort.lowIndex, lowIndex)
	require.Equal(t, cohort.highIndex, highIndex)
	require.Equal(t, cohort.allValuesWritten, false)

	// Check if the cohort file exists
	filePath := cohort.Path(false)
	exists, err := util.Exists(filePath)
	require.NoError(t, err)
	require.True(t, exists)

	// Initialize a copy cohort from the file
	loadedCohort, err := LoadCohort(testDirectory, cohortIndex)
	require.NoError(t, err)
	require.Equal(t, loadedCohort.cohortIndex, cohortIndex)
	require.Equal(t, loadedCohort.lowIndex, lowIndex)
	require.Equal(t, loadedCohort.highIndex, highIndex)
	require.Equal(t, loadedCohort.allValuesWritten, false)

	// Mark the cohort as written
	err = loadedCohort.MarkComplete()
	require.NoError(t, err)
	require.True(t, loadedCohort.allValuesWritten)

	// Load the cohort again.
	loadedCohort, err = LoadCohort(testDirectory, cohortIndex)
	require.NoError(t, err)
	require.Equal(t, loadedCohort.cohortIndex, cohortIndex)
	require.Equal(t, loadedCohort.lowIndex, lowIndex)
	require.Equal(t, loadedCohort.highIndex, highIndex)
	require.Equal(t, loadedCohort.allValuesWritten, true)
}
