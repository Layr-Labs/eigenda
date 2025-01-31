package verification

import (
	"errors"
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"
)

// CheckBlobLength accepts bytes representing a blob, and a claimed length in symbols. Note that claimed length is an
// upper bound, not a precise length. Two length checks are performed:
//
// 1. Blob doesn't have length 0
// 2. Blob length is <= the claimed blob length. Claimed blob length is from the BlobCommitment
func CheckBlobLength(blobBytes []byte, claimedBlobLength uint) error {
	if len(blobBytes) == 0 {
		return errors.New("blob has length 0")
	}

	if uint(len(blobBytes)) > claimedBlobLength*encoding.BYTES_PER_SYMBOL {
		return fmt.Errorf(
			"length (%d bytes) is greater than claimed blob length (%d bytes)",
			len(blobBytes),
			claimedBlobLength*encoding.BYTES_PER_SYMBOL)
	}

	return nil
}
