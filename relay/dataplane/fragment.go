package dataplane

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// GetFragmentCount returns the number of fragments that a file of the given size will be broken into.
func GetFragmentCount(fileSize int, fragmentSize int) int {
	if fileSize < fragmentSize {
		return 1
	} else if fileSize%fragmentSize == 0 {
		return fileSize / fragmentSize
	} else {
		return fileSize/fragmentSize + 1
	}
}

// GetFragmentKey returns the key for the fragment at the given index.
//
// Fragment keys take the form of "prefix/body-index[f]". The prefix is the first prefixLength characters
// of the file key. The body is the file key. The index is the index of the fragment. The character "f" is appended
// to the key of the last fragment in the series.
//
// Example: fileKey="abc123", prefixLength=2, fragmentCount=3
// The keys will be "ab/abc123-0", "ab/abc123-1", "ab/abc123-2f"
func GetFragmentKey(fileKey string, prefixLength int, fragmentCount int, index int) (string, error) {
	var prefix string
	if prefixLength > len(fileKey) {
		prefix = fileKey
	} else {
		prefix = fileKey[:prefixLength]
	}

	postfix := ""
	if fragmentCount-1 == index {
		postfix = "f"
	}

	if index >= fragmentCount {
		return "", errors.New(fmt.Sprintf("index %d is too high for fragment count %d", index, fragmentCount))
	}

	return fmt.Sprintf("%s/%s-%d%s", prefix, fileKey, index, postfix), nil
}

// Fragment is a subset of a file.
type Fragment struct {
	FragmentKey string
	Data        []byte
	Index       int
}

// BreakIntoFragments breaks a file into fragments of the given size.
func BreakIntoFragments(fileKey string, data []byte, prefixLength int, fragmentSize int) ([]*Fragment, error) {
	fragmentCount := GetFragmentCount(len(data), fragmentSize)
	fragments := make([]*Fragment, fragmentCount)
	for i := 0; i < fragmentCount; i++ {
		start := i * fragmentSize
		end := start + fragmentSize
		if end > len(data) {
			end = len(data)
		}

		fragmentKey, err := GetFragmentKey(fileKey, prefixLength, fragmentCount, i)
		if err != nil {
			return nil, err
		}
		fragments[i] = &Fragment{
			FragmentKey: fragmentKey,
			Data:        data[start:end],
			Index:       i,
		}
	}
	return fragments, nil
}

// GetFragmentKeys returns the keys for all fragments of a file.
func GetFragmentKeys(fileKey string, prefixLength int, fragmentCount int) ([]string, error) {
	keys := make([]string, fragmentCount)
	for i := 0; i < fragmentCount; i++ {
		fragmentKey, err := GetFragmentKey(fileKey, prefixLength, fragmentCount, i)
		if err != nil {
			return nil, err
		}
		keys[i] = fragmentKey
	}
	return keys, nil
}

// RecombineFragments recombines fragments into a single file.
// Returns an error if any fragments are missing.
func RecombineFragments(fragments []*Fragment) ([]byte, error) {

	if len(fragments) == 0 {
		return nil, fmt.Errorf("no fragments")
	}

	// Sort the fragments by index
	sort.Slice(fragments, func(i, j int) bool {
		return fragments[i].Index < fragments[j].Index
	})

	// Make sure there aren't any gaps in the fragment indices
	dataSize := 0
	for i, fragment := range fragments {
		if fragment.Index != i {
			return nil, fmt.Errorf("missing fragment with index %d", i)
		}
		dataSize += len(fragment.Data)
	}

	// Make sure we have the last fragment
	if !strings.HasSuffix(fragments[len(fragments)-1].FragmentKey, "f") {
		return nil, fmt.Errorf("missing final fragment")
	}

	fragmentSize := len(fragments[0].Data)

	// Concatenate the data
	result := make([]byte, dataSize)
	for _, fragment := range fragments {
		copy(result[fragment.Index*fragmentSize:], fragment.Data)
	}

	return result, nil
}
