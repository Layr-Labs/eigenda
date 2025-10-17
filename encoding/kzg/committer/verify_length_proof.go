package committer

import (
	"errors"
	"fmt"
	"math/bits"

	"github.com/Layr-Labs/eigenda/common/math"
	eigenbn254 "github.com/Layr-Labs/eigenda/crypto/ecc/bn254"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/resources/srs"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
)

// VerifyLengthProof by itself is not sufficient to verify the length of a blob commitment!
// It must be used in conjunction with VerifyCommitEquivalenceBatch to ensure that the
// blob commitment on G1 and blob commitment on G2 (LengthCommitment) are equivalent.
func VerifyLengthProof(commitments encoding.BlobCommitments) error {
	return verifyLengthProof(
		(*bn254.G2Affine)(commitments.LengthCommitment),
		(*bn254.G2Affine)(commitments.LengthProof),
		uint64(commitments.Length),
	)
}

// verifyLengthProof verifies the length proof (low degree proof).
// See https://layr-labs.github.io/eigenda/protocol/architecture/encoding.html#validation-via-kzg
//
// This function verifies a low degree proof against a poly commitment.
// We wish to show x^shift poly = shiftedPoly, with shift = 2^28 - blob_length.
// We verify this by checking the pairing equation:
// e( s^shift G1, p(s)G2 ) = e( G1, p(s^shift)G2 )
// Note that we also need to verify that the blob_commitment and length_commitment are equivalent,
// by verifying the other pairing equation: e(blob_commitment,G2) = e(length_commitment,C2)
// This is done in [VerifyCommitEquivalenceBatch].
// TODO(samlaf): can we move combine the 2 pairings into a single function?
func verifyLengthProof(
	lengthCommit *bn254.G2Affine, lengthProof *bn254.G2Affine, commitmentLength uint64,
) error {
	// This also prevents commitmentLength=0.
	if !math.IsPowerOfTwo(commitmentLength) {
		return fmt.Errorf("commitment length %d is not a power of 2", commitmentLength)
	}
	// Because commitmentLength is power of 2, we know its represented as 100..0 in binary,
	// so counting the number of trailing zeros gives us log2(commitmentLength).
	// We need commitmentLengthLog <= 27 because we have hardcoded SRS points only for that range.
	commitmentLengthLog := bits.TrailingZeros64(commitmentLength)
	if commitmentLengthLog > 27 {
		return fmt.Errorf("commitment length %d is > max possible 2^28", commitmentLength)
	}
	// g1Challenge = [tau^(2^28 - commitmentLength)]_1
	// G1ReversePowerOf2SRS contains the 28 hardcoded points that we need.
	g1Challenge := srs.G1ReversePowerOf2SRS[commitmentLengthLog]

	err := eigenbn254.PairingsVerify(&g1Challenge, lengthCommit, &kzg.GenG1, lengthProof)
	if err != nil {
		return fmt.Errorf("verify pairing: %w", err)
	}
	return nil
}

type CommitmentPair struct {
	Commitment       bn254.G1Affine
	LengthCommitment bn254.G2Affine
}

// VerifyCommitEquivalenceBatch is conceptually part of VerifyLengthProof.
// It's currently a separate function for historical reasons, from the times when we were batching.
// Now that we no longer are batching, we could verify a single commitmentEquivalence at a time,
// and do so as part of VerifyLengthProof.
// TODO(samlaf): refactor into a single VerifyLengthProof function.
func VerifyCommitEquivalenceBatch(commitments []encoding.BlobCommitments) error {
	commitmentsPair := make([]CommitmentPair, len(commitments))

	for i, c := range commitments {
		commitmentsPair[i] = CommitmentPair{
			Commitment:       (bn254.G1Affine)(*c.Commitment),
			LengthCommitment: (bn254.G2Affine)(*c.LengthCommitment),
		}
	}
	return batchVerifyCommitEquivalence(commitmentsPair)
}

func batchVerifyCommitEquivalence(commitmentsPair []CommitmentPair) error {

	g1commits := make([]bn254.G1Affine, len(commitmentsPair))
	g2commits := make([]bn254.G2Affine, len(commitmentsPair))
	for i := 0; i < len(commitmentsPair); i++ {
		g1commits[i] = commitmentsPair[i].Commitment
		g2commits[i] = commitmentsPair[i].LengthCommitment
	}

	randomsFr, err := eigenbn254.RandomFrs(len(g1commits))
	if err != nil {
		return fmt.Errorf("create randomness vector: %w", err)
	}

	var lhsG1 bn254.G1Affine
	_, err = lhsG1.MultiExp(g1commits, randomsFr, ecc.MultiExpConfig{})
	if err != nil {
		return fmt.Errorf("compute lhsG1: %w", err)
	}

	lhsG2 := &kzg.GenG2

	var rhsG2 bn254.G2Affine
	_, err = rhsG2.MultiExp(g2commits, randomsFr, ecc.MultiExpConfig{})
	if err != nil {
		return fmt.Errorf("compute rhsG2: %w", err)
	}
	rhsG1 := &kzg.GenG1

	err = eigenbn254.PairingsVerify(&lhsG1, lhsG2, rhsG1, &rhsG2)
	if err == nil {
		return nil
	} else {
		return errors.New("incorrect universal batch verification")
	}
}
