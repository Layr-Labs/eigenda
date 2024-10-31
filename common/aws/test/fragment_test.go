package test

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"strings"
	"testing"
)

func TestGetFragmentCount(t *testing.T) {
	tu.InitializeRandom()

	// Test a file smaller than a fragment
	fileSize := rand.Intn(100) + 100
	fragmentSize := fileSize * 2
	fragmentCount := s3.GetFragmentCount(fileSize, fragmentSize)
	assert.Equal(t, 1, fragmentCount)

	// Test a file that can fit in a single fragment
	fileSize = rand.Intn(100) + 100
	fragmentSize = fileSize
	fragmentCount = s3.GetFragmentCount(fileSize, fragmentSize)
	assert.Equal(t, 1, fragmentCount)

	// Test a file that is one byte larger than a fragment
	fileSize = rand.Intn(100) + 100
	fragmentSize = fileSize - 1
	fragmentCount = s3.GetFragmentCount(fileSize, fragmentSize)
	assert.Equal(t, 2, fragmentCount)

	// Test a file that is one less than a multiple of the fragment size
	fragmentSize = rand.Intn(100) + 100
	expectedFragmentCount := rand.Intn(10) + 1
	fileSize = fragmentSize*expectedFragmentCount - 1
	fragmentCount = s3.GetFragmentCount(fileSize, fragmentSize)
	assert.Equal(t, expectedFragmentCount, fragmentCount)

	// Test a file that is a multiple of the fragment size
	fragmentSize = rand.Intn(100) + 100
	expectedFragmentCount = rand.Intn(10) + 1
	fileSize = fragmentSize * expectedFragmentCount
	fragmentCount = s3.GetFragmentCount(fileSize, fragmentSize)
	assert.Equal(t, expectedFragmentCount, fragmentCount)

	// Test a file that is one more than a multiple of the fragment size
	fragmentSize = rand.Intn(100) + 100
	expectedFragmentCount = rand.Intn(10) + 2
	fileSize = fragmentSize*(expectedFragmentCount-1) + 1
	fragmentCount = s3.GetFragmentCount(fileSize, fragmentSize)
	assert.Equal(t, expectedFragmentCount, fragmentCount)
}

// Fragment keys take the form of "prefix/body-index[f]". Verify the prefix part of the key.
func TestPrefix(t *testing.T) {
	tu.InitializeRandom()

	keyLength := rand.Intn(10) + 10
	key := tu.RandomString(keyLength)

	for i := 0; i < keyLength*2; i++ {
		fragmentCount := rand.Intn(10) + 10
		fragmentIndex := rand.Intn(fragmentCount)
		fragmentKey, err := s3.GetFragmentKey(key, i, fragmentCount, fragmentIndex)
		assert.NoError(t, err)

		parts := strings.Split(fragmentKey, "/")
		assert.Equal(t, 2, len(parts))
		prefix := parts[0]

		if i >= keyLength {
			assert.Equal(t, key, prefix)
		} else {
			assert.Equal(t, key[:i], prefix)
		}
	}
}

// Fragment keys take the form of "prefix/body-index[f]". Verify the body part of the key.
func TestKeyBody(t *testing.T) {
	tu.InitializeRandom()

	for i := 0; i < 10; i++ {
		keyLength := rand.Intn(10) + 10
		key := tu.RandomString(keyLength)
		fragmentCount := rand.Intn(10) + 10
		fragmentIndex := rand.Intn(fragmentCount)
		fragmentKey, err := s3.GetFragmentKey(key, rand.Intn(10), fragmentCount, fragmentIndex)
		assert.NoError(t, err)

		parts := strings.Split(fragmentKey, "/")
		assert.Equal(t, 2, len(parts))
		parts = strings.Split(parts[1], "-")
		assert.Equal(t, 2, len(parts))
		body := parts[0]

		assert.Equal(t, key, body)
	}
}

// Fragment keys take the form of "prefix/body-index[f]". Verify the index part of the key.
func TestKeyIndex(t *testing.T) {
	tu.InitializeRandom()

	for i := 0; i < 10; i++ {
		fragmentCount := rand.Intn(10) + 10
		index := rand.Intn(fragmentCount)
		fragmentKey, err := s3.GetFragmentKey(tu.RandomString(10), rand.Intn(10), fragmentCount, index)
		assert.NoError(t, err)

		parts := strings.Split(fragmentKey, "/")
		assert.Equal(t, 2, len(parts))
		parts = strings.Split(parts[1], "-")
		assert.Equal(t, 2, len(parts))
		indexStr := parts[1]
		assert.True(t, strings.HasPrefix(indexStr, fmt.Sprintf("%d", index)))
	}
}

// Fragment keys take the form of "prefix/body-index[f]".
// Verify the postfix part of the key, which should be "f" for the last fragment.
func TestKeyPostfix(t *testing.T) {
	tu.InitializeRandom()

	segmentCount := rand.Intn(10) + 10

	for i := 0; i < segmentCount; i++ {
		fragmentKey, err := s3.GetFragmentKey(tu.RandomString(10), rand.Intn(10), segmentCount, i)
		assert.NoError(t, err)

		if i == segmentCount-1 {
			assert.True(t, strings.HasSuffix(fragmentKey, "f"))
		} else {
			assert.False(t, strings.HasSuffix(fragmentKey, "f"))
		}
	}
}

func TestExampleInGodoc(t *testing.T) {
	fileKey := "abc123"
	prefixLength := 2
	fragmentCount := 3
	fragmentKeys, err := s3.GetFragmentKeys(fileKey, prefixLength, fragmentCount)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(fragmentKeys))
	assert.Equal(t, "ab/abc123-0", fragmentKeys[0])
	assert.Equal(t, "ab/abc123-1", fragmentKeys[1])
	assert.Equal(t, "ab/abc123-2f", fragmentKeys[2])
}

func TestGetFragmentKeys(t *testing.T) {
	tu.InitializeRandom()

	fileKey := tu.RandomString(10)
	prefixLength := rand.Intn(3) + 1
	fragmentCount := rand.Intn(10) + 10

	fragmentKeys, err := s3.GetFragmentKeys(fileKey, prefixLength, fragmentCount)
	assert.NoError(t, err)
	assert.Equal(t, fragmentCount, len(fragmentKeys))

	for i := 0; i < fragmentCount; i++ {
		expectedKey, err := s3.GetFragmentKey(fileKey, prefixLength, fragmentCount, i)
		assert.NoError(t, err)
		assert.Equal(t, expectedKey, fragmentKeys[i])

		parts := strings.Split(fragmentKeys[i], "/")
		assert.Equal(t, 2, len(parts))
		parsedPrefix := parts[0]
		assert.Equal(t, fileKey[:prefixLength], parsedPrefix)
		parts = strings.Split(parts[1], "-")
		assert.Equal(t, 2, len(parts))
		parsedKey := parts[0]
		assert.Equal(t, fileKey, parsedKey)
		index := parts[1]

		if i == fragmentCount-1 {
			assert.Equal(t, fmt.Sprintf("%d", i)+"f", index)
		} else {
			assert.Equal(t, fmt.Sprintf("%d", i), index)
		}
	}
}

func TestGetFragments(t *testing.T) {
	tu.InitializeRandom()

	fileKey := tu.RandomString(10)
	data := tu.RandomBytes(1000)
	prefixLength := rand.Intn(3) + 1
	fragmentSize := rand.Intn(100) + 100

	fragments, err := s3.BreakIntoFragments(fileKey, data, prefixLength, fragmentSize)
	assert.NoError(t, err)
	assert.Equal(t, s3.GetFragmentCount(len(data), fragmentSize), len(fragments))

	totalSize := 0

	for i, fragment := range fragments {
		fragmentKey, err := s3.GetFragmentKey(fileKey, prefixLength, len(fragments), i)
		assert.NoError(t, err)
		assert.Equal(t, fragmentKey, fragment.FragmentKey)

		start := i * fragmentSize
		end := start + fragmentSize
		if end > len(data) {
			end = len(data)
		}
		assert.Equal(t, data[start:end], fragment.Data)
		assert.Equal(t, i, fragment.Index)
		totalSize += len(fragment.Data)
	}

	assert.Equal(t, len(data), totalSize)
}

func TestGetFragmentsSmallFile(t *testing.T) {
	tu.InitializeRandom()

	fileKey := tu.RandomString(10)
	data := tu.RandomBytes(10)
	prefixLength := rand.Intn(3) + 1
	fragmentSize := rand.Intn(100) + 100

	fragments, err := s3.BreakIntoFragments(fileKey, data, prefixLength, fragmentSize)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(fragments))

	fragmentKey, err := s3.GetFragmentKey(fileKey, prefixLength, 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, fragmentKey, fragments[0].FragmentKey)
	assert.Equal(t, data, fragments[0].Data)
	assert.Equal(t, 0, fragments[0].Index)
}

func TestGetFragmentsExactlyOnePerfectlySizedFile(t *testing.T) {
	tu.InitializeRandom()

	fileKey := tu.RandomString(10)
	fragmentSize := rand.Intn(100) + 100
	data := tu.RandomBytes(fragmentSize)
	prefixLength := rand.Intn(3) + 1

	fragments, err := s3.BreakIntoFragments(fileKey, data, prefixLength, fragmentSize)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(fragments))

	fragmentKey, err := s3.GetFragmentKey(fileKey, prefixLength, 1, 0)
	assert.NoError(t, err)
	assert.Equal(t, fragmentKey, fragments[0].FragmentKey)
	assert.Equal(t, data, fragments[0].Data)
	assert.Equal(t, 0, fragments[0].Index)
}

func TestRecombineFragments(t *testing.T) {
	tu.InitializeRandom()

	fileKey := tu.RandomString(10)
	data := tu.RandomBytes(1000)
	prefixLength := rand.Intn(3) + 1
	fragmentSize := rand.Intn(100) + 100

	fragments, err := s3.BreakIntoFragments(fileKey, data, prefixLength, fragmentSize)
	assert.NoError(t, err)
	recombinedData, err := s3.RecombineFragments(fragments)
	assert.NoError(t, err)
	assert.Equal(t, data, recombinedData)

	// Shuffle the fragments
	for i := range fragments {
		j := rand.Intn(i + 1)
		fragments[i], fragments[j] = fragments[j], fragments[i]
	}

	recombinedData, err = s3.RecombineFragments(fragments)
	assert.NoError(t, err)
	assert.Equal(t, data, recombinedData)
}

func TestRecombineFragmentsSmallFile(t *testing.T) {
	tu.InitializeRandom()

	fileKey := tu.RandomString(10)
	data := tu.RandomBytes(10)
	prefixLength := rand.Intn(3) + 1
	fragmentSize := rand.Intn(100) + 100

	fragments, err := s3.BreakIntoFragments(fileKey, data, prefixLength, fragmentSize)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(fragments))
	recombinedData, err := s3.RecombineFragments(fragments)
	assert.NoError(t, err)
	assert.Equal(t, data, recombinedData)
}

func TestMissingFragment(t *testing.T) {
	tu.InitializeRandom()

	fileKey := tu.RandomString(10)
	data := tu.RandomBytes(1000)
	prefixLength := rand.Intn(3) + 1
	fragmentSize := rand.Intn(100) + 100

	fragments, err := s3.BreakIntoFragments(fileKey, data, prefixLength, fragmentSize)
	assert.NoError(t, err)

	fragmentIndexToSkip := rand.Intn(len(fragments))
	fragments = append(fragments[:fragmentIndexToSkip], fragments[fragmentIndexToSkip+1:]...)

	_, err = s3.RecombineFragments(fragments[:len(fragments)-1])
	assert.Error(t, err)
}

func TestMissingFinalFragment(t *testing.T) {
	tu.InitializeRandom()

	fileKey := tu.RandomString(10)
	data := tu.RandomBytes(1000)
	prefixLength := rand.Intn(3) + 1
	fragmentSize := rand.Intn(100) + 100

	fragments, err := s3.BreakIntoFragments(fileKey, data, prefixLength, fragmentSize)
	assert.NoError(t, err)
	fragments = fragments[:len(fragments)-1]

	_, err = s3.RecombineFragments(fragments)
	assert.Error(t, err)
}
