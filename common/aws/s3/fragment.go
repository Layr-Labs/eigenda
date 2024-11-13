package s3

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// getFragmentCount returns the number of fragments that a file of the given size will be broken into.
func getFragmentCount(fileSize int, fragmentSize int) int {
	if fileSize < fragmentSize {
		return 1
	} else if fileSize%fragmentSize == 0 {
		return fileSize / fragmentSize
	} else {
		return fileSize/fragmentSize + 1
	}
}

// getFragmentKey returns the key for the fragment at the given index.
//
// Fragment keys take the form of "body-index[f]". The index is the index of the fragment. The character "f" is
// appended to the key of the last fragment in the series.
//
// Example: fileKey="abc123", fragmentCount=3
// The keys will be "abc123-0", "abc123-1", "abc123-2f"
func getFragmentKey(fileKey string, fragmentCount int, index int) (string, error) {

	postfix := ""
	if fragmentCount-1 == index {
		postfix = "f"
	}

	if index >= fragmentCount {
		return "", fmt.Errorf("index %d is too high for fragment count %d", index, fragmentCount)
	}

	return fmt.Sprintf("%s-%d%s", fileKey, index, postfix), nil
}

// Fragment is a subset of a file.
type Fragment struct {
	FragmentKey string
	Data        []byte
	Index       int
}

// BreakIntoFragments breaks a file into fragments of the given size.
func BreakIntoFragments(fileKey string, data []byte, fragmentSize int) ([]*Fragment, error) {
	fragmentCount := getFragmentCount(len(data), fragmentSize)
	fragments := make([]*Fragment, fragmentCount)
	for i := 0; i < fragmentCount; i++ {
		start := i * fragmentSize
		end := start + fragmentSize
		if end > len(data) {
			end = len(data)
		}

		fragmentKey, err := getFragmentKey(fileKey, fragmentCount, i)
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
func GetFragmentKeys(fileKey string, fragmentCount int) ([]string, error) {
	keys := make([]string, fragmentCount)
	for i := 0; i < fragmentCount; i++ {
		fragmentKey, err := getFragmentKey(fileKey, fragmentCount, i)
		if err != nil {
			return nil, err
		}
		keys[i] = fragmentKey
	}
	return keys, nil
}

// AllFragmentsExist returns true if all keys are fragment keys.
// It checks if all fragment keys starting with 0th fragment and ending with nth fragment (marked by "f" postfix) exist.
// Warning: this function sorts the input slice in place and therefore mutates the ordering of the input slice.
func SortAndCheckAllFragmentsExist(keys []string) bool {
	sort.Strings(keys)
	for i, key := range keys {
		if !strings.HasSuffix(key, "-"+strconv.Itoa(i)) {
			if strings.HasSuffix(key, "-"+strconv.Itoa(i)+"f") {
				return i == len(keys)-1
			}
			return false
		}
	}

	return false
}

// recombineFragments recombines fragments into a single file.
// Returns an error if any fragments are missing.
func recombineFragments(fragments []*Fragment) ([]byte, error) {

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
