package segment

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

// verifyFilePermissions checks if a file has read/write permissions and is a regular file (if it exists),
// returning an error if it does not if the file permissions or file type is not as expected.
// Also returns a boolean indicating if the file exists and its size (to save on additional os.Stat calls).
//
// A file is considered to have the correct permissions/type if:
// - it exists and is a standard file with read+write permissions
// - if it does not exist but its parent directory has read+write permissions.
//
// The arguments for the function are the result of os.Stat(path). There is no need to do error checking on the
// result of os.Stat in the calling context (this method does it for you).
func verifyFilePermissions(path string) (exists bool, size int64, err error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// The file does not exist. Check the parent.
			parentPath := filepath.Dir(path)
			parentInfo, err := os.Stat(parentPath)
			if err != nil {
				if os.IsNotExist(err) {
					return false, -1, fmt.Errorf("parent directory %s does not exist", parentPath)
				}
				return false, -1, fmt.Errorf("failed to stat parent directory %s: %w", parentPath, err)
			}

			if !parentInfo.IsDir() {
				return false, -1, fmt.Errorf("parent directory %s is not a directory", parentPath)
			}

			if parentInfo.Mode()&0700 != 0700 {
				return false, -1, fmt.Errorf("parent directory %s has insufficent permissions", parentPath)
			}
		}

		return false, 0, nil
	}

	// File exists. Check if it is a regular file and that it is readable+writeable.
	if info.IsDir() {
		return false, -1, fmt.Errorf("file %s is a directory", path)
	}
	if info.Mode()&0600 != 0600 {
		return false, -1, fmt.Errorf("file %s has insufficent permissions", path)
	}

	return true, info.Size(), nil
}

// perm64 computes A permutation (invertible function) on 64 bits.
// The constants were found by automated search, to
// optimize avalanche. Avalanche means that for a
// random number x, flipping bit i of x has about a
// 50 percent chance of flipping bit j of perm64(x).
// For each possible pair (i,j), this function achieves
// a probability between 49.8 and 50.2 percent.
//
// Warning: this is not a cryptographic hash function. This hash function may be suitable for hash tables, but not for
// cryptographic purposes. It is trivially easy to reverse this function.
//
// Algorithm borrowed from https://github.com/hiero-ledger/hiero-consensus-node/blob/main/platform-sdk/swirlds-common/src/main/java/com/swirlds/common/utility/NonCryptographicHashing.java
// (original implementation is under Apache 2.0 license, algorithm designed by Leemon Baird)
func perm64(x uint64) uint64 {
	// This is necessary so that 0 does not hash to 0.
	// As a side effect this constant will hash to 0.
	x ^= 0x5e8a016a5eb99c18

	x += x << 30
	x ^= x >> 27
	x += x << 16
	x ^= x >> 20
	x += x << 5
	x ^= x >> 18
	x += x << 10
	x ^= x >> 24
	x += x << 30
	return x
}

// perm64Bytes hashes a byte slice using perm64.
func perm64Bytes(b []byte) uint64 {
	x := uint64(0)

	for i := 0; i < len(b); i += 8 {
		var next uint64
		if i+8 > len(b) {
			// grab the next 8 bytes
			next = binary.BigEndian.Uint64(b[i:])
		} else {
			// insufficient bytes, pad with zeros
			nextBytes := make([]byte, 8)
			copy(nextBytes, b[i:])
			next = binary.BigEndian.Uint64(nextBytes)
		}
		x ^= perm64(next)
	}

	return x
}

// hashKey hashes a key using perm64 and a salt.
func hashKey(key []byte, salt uint32) uint32 {
	return uint32(perm64Bytes(key) ^ uint64(salt))
}
