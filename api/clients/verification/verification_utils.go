package verification

import (
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/encoding"
	"slices"

	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"

	disperserpb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

// VerifyBlobStatusReply verifies all blob data, as required to trust that the blob is valid
func VerifyBlobStatusReply(
	kzgVerifier *verifier.Verifier,
	reply *disperserpb.BlobStatusReply,
// these bytes must represent the blob in coefficient form (i.e. IFFTed)
	blobBytes []byte) error {

	blobHeader, err := core.BlobHeaderFromProtobuf(reply.BlobVerificationInfo.BlobCertificate.BlobHeader)
	if err != nil {
		return fmt.Errorf("blob header from protobuf: %w", err)
	}
	blobKey, err := blobHeader.BlobKey()
	if err != nil {
		return fmt.Errorf("compute blob key: %w", err)
	}

	err = VerifyMerkleProof(
		&blobKey,
		reply.BlobVerificationInfo.InclusionProof,
		reply.SignedBatch.Header.BatchRoot,
		reply.BlobVerificationInfo.BlobIndex)
	if err != nil {
		return fmt.Errorf("verify merkle proof: %w", err)
	}

	commitments, err := encoding.BlobCommitmentsFromProtobuf(reply.BlobVerificationInfo.BlobCertificate.BlobHeader.Commitment)
	if err != nil {
		return fmt.Errorf("commitments from protobuf: %w", err)
	}

	err = VerifyKzgCommitment(kzgVerifier, commitments.Commitment, blobBytes)
	if err != nil {
		return fmt.Errorf("verify commitment: %w", err)
	}

	if uint(len(blobBytes)) != commitments.Length {
		return fmt.Errorf("actual blob length (%d) doesn't match claimed length in commitment (%d)", len(blobBytes), commitments.Length)
	}

	err = kzgVerifier.VerifyBlobLength(*commitments)
	if err != nil {
		return fmt.Errorf("verify blob length: %w", err)
	}

	return nil
}

// VerifyKzgCommitment asserts that the claimed commitment from the certificate matches the commitment computed
// from the blob.
//
// TODO: Optimize implementation by opening a point on the commitment instead (don’t compute the g2 point. ask Bowen)
func VerifyKzgCommitment(
	kzgVerifier *verifier.Verifier,
	claimedCommitment *encoding.G1Commitment,
	blobBytes []byte) error {

	computedCommitment, err := GenerateBlobCommitment(kzgVerifier, blobBytes)
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

// GenerateBlobCommitment computes kzg-bn254 commitment of blob data using SRS
//
// The blob data input to this method should be in coefficient form (IFFTed)
func GenerateBlobCommitment(
	kzgVerifier *verifier.Verifier,
	blob []byte) (*bn254.G1Affine, error) {

	inputFr, err := rs.ToFrArray(blob)
	if err != nil {
		return nil, fmt.Errorf("cannot convert bytes to field elements, %w", err)
	}

	if len(kzgVerifier.Srs.G1) < len(inputFr) {
		return nil, fmt.Errorf(
			"cannot verify commitment because the number of stored srs in the memory is insufficient, have %v need %v",
			len(kzgVerifier.Srs.G1),
			len(inputFr))
	}

	config := ecc.MultiExpConfig{}
	var commitment bn254.G1Affine
	_, err = commitment.MultiExp(kzgVerifier.Srs.G1[:len(inputFr)], inputFr, config)
	if err != nil {
		return nil, err
	}

	return &commitment, nil
}

// TODO: I don't think that any of the checks done in `verifySecurityParams` from the proxy code are necessary in v2. confirm this.

// VerifyMerkleProof verifies the blob batch inclusion proof against the claimed batch root
func VerifyMerkleProof(
// the key of the blob, which functions as the leaf in the batch merkle tree
	blobKey *core.BlobKey,
// the inclusion proof, which contains the sibling hashes necessary to compute the root hash starting at the leaf
	inclusionProof []byte,
// the claimed merkle root hash, which must be verified
	claimedRoot []byte,
// the index of the blob in the batch. this informs whether the leaf being verified is a left or right child
	blobIndex uint32) error {

	generatedRoot, err := ProcessInclusionProof(inclusionProof, blobKey, uint64(blobIndex))
	if err != nil {
		return err
	}

	equal := slices.Equal(claimedRoot, generatedRoot.Bytes())
	if !equal {
		return fmt.Errorf("root hash mismatch, expected: %x, got: %x", claimedRoot, generatedRoot)
	}

	return nil
}

// ProcessInclusionProof computes the root hash, using the leaf and relevant inclusion proof
//
// This logic is implemented here, rather than using merkletree.VerifyProofUsing, in order to exactly mirror
// https://github.com/Layr-Labs/eigenlayer-contracts/blob/dev/src/contracts/libraries/Merkle.sol#L49-L76, which is the
// verification method that executes on chain. It is important that the offchain verification result exactly matches
// the onchain result
func ProcessInclusionProof(proof []byte, blobKey *core.BlobKey, index uint64) (common.Hash, error) {
	if len(proof)%32 != 0 {
		return common.Hash{}, errors.New("proof length should be a multiple of 32 bytes or 256 bits")
	}

	// computedHash starts out equal to the hash of the leaf (the blob key)
	var computedHash common.Hash
	copy(computedHash[:], blobKey[:])

	// we then work our way up the merkle tree, to compute the root hash
	for i := 0; i < len(proof); i += 32 {
		var proofElement common.Hash
		copy(proofElement[:], proof[i:i+32])

		var combined []byte
		if index%2 == 0 { // right
			combined = append(computedHash.Bytes(), proofElement.Bytes()...)
		} else { // left
			combined = append(proofElement.Bytes(), computedHash.Bytes()...)
		}

		computedHash = crypto.Keccak256Hash(combined)
		index /= 2
	}

	return computedHash, nil
}
