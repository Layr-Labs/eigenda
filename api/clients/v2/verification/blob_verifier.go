package verification

import (
	"errors"
	"fmt"

	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

// VerifyBlobAgainstCert verifies the blob against a kzg commitment
//
// If verification succeeds, the method returns nil. Otherwise, it returns an error.
//
// TODO: in the future, this will be optimized to use fiat shamir transformation for verification, rather than
//
//	regenerating the commitment: https://github.com/Layr-Labs/eigenda/issues/1037
func VerifyBlobAgainstCert(
	blobKey *core.BlobKey,
	blobBytes []byte,
	kzgCommitment *encoding.G1Commitment,
	g1Srs []bn254.G1Affine) error {

	valid, err := GenerateAndCompareBlobCommitment(g1Srs, blobBytes, kzgCommitment)
	if err != nil {
		return fmt.Errorf("generate and compare commitment for blob %v: %w", blobKey.Hex(), err)
	}

	if !valid {
		return fmt.Errorf("commitment for blob %v is invalid", blobKey.Hex())
	}

	return nil
}

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
