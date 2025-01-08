package verification

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"

	"github.com/Layr-Labs/eigenda/encoding/rs"
)

// GenerateBlobCommitment computes a kzg-bn254 commitment of blob data using SRS
func GenerateBlobCommitment(
	g1Srs []bn254.G1Affine,
	blobBytes []byte) (*encoding.G1Commitment, error) {

	inputFr, err := rs.ToFrArray(blobBytes)
	if err != nil {
		return nil, fmt.Errorf("convert bytes to field elements, %w", err)
	}

	if len(g1Srs) < len(inputFr) {
		return nil, fmt.Errorf(
			"insufficient SRS in memory: have %v, need %v",
			len(g1Srs),
			len(inputFr))
	}

	var commitment bn254.G1Affine
	_, err = commitment.MultiExp(g1Srs[:len(inputFr)], inputFr, ecc.MultiExpConfig{})
	if err != nil {
		return nil, fmt.Errorf("MultiExp: %w", err)
	}

	return &encoding.G1Commitment{X: commitment.X, Y: commitment.Y}, nil
}

// GenerateAndCompareBlobCommitment generates the kzg-bn254 commitment of the blob, and compares it with a claimed
// commitment. An error is returned if there is a problem generating the commitment, or if the comparison fails.
func GenerateAndCompareBlobCommitment(
	g1Srs []bn254.G1Affine,
	blobBytes []byte,
	claimedCommitment *encoding.G1Commitment) error {

	computedCommitment, err := GenerateBlobCommitment(g1Srs, blobBytes)
	if err != nil {
		return fmt.Errorf("compute commitment: %w", err)
	}

	if claimedCommitment.X.Equal(&computedCommitment.X) &&
		claimedCommitment.Y.Equal(&computedCommitment.Y) {
		return nil
	}

	return fmt.Errorf(
		"commitment field elements do not match. computed commitment: (x: %x, y: %x), claimed commitment (x: %x, y: %x)",
		computedCommitment.X, computedCommitment.Y, claimedCommitment.X, claimedCommitment.Y)
}
