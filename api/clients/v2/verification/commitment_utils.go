package verification

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// GenerateBlobCommitment computes a kzg-bn254 commitment from field element coefficients using SRS
func GenerateBlobCommitment(g1Srs []bn254.G1Affine, coefficients []fr.Element) (*encoding.G1Commitment, error) {

	if len(g1Srs) < len(coefficients) {
		return nil, fmt.Errorf(
			"insufficient SRS in memory: have %v, need %v",
			len(g1Srs),
			len(coefficients))
	}

	var commitment bn254.G1Affine
	_, err := commitment.MultiExp(g1Srs[:len(coefficients)], coefficients, ecc.MultiExpConfig{})
	if err != nil {
		return nil, fmt.Errorf("MultiExp: %w", err)
	}

	return &encoding.G1Commitment{X: commitment.X, Y: commitment.Y}, nil
}

// GenerateAndCompareBlobCommitment generates the kzg-bn254 commitment of the blob, and compares it with a claimed
// commitment. An error is returned if there is a problem generating the commitment. True is returned if the commitment
// is successfully generated, and is equal to the claimed commitment, otherwise false.
func GenerateAndCompareBlobCommitment(
	g1Srs []bn254.G1Affine,
	blob *coretypes.Blob,
	claimedCommitment *encoding.G1Commitment,
) (bool, error) {

	computedCommitment, err := GenerateBlobCommitment(g1Srs, blob.GetCoefficients())
	if err != nil {
		return false, fmt.Errorf("compute commitment: %w", err)
	}

	if claimedCommitment.X.Equal(&computedCommitment.X) &&
		claimedCommitment.Y.Equal(&computedCommitment.Y) {
		return true, nil
	}

	return false, nil
}
