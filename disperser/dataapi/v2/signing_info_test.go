package v2

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/stretchr/testify/require"
)

func TestComputeBlobQuorumSigningStats(t *testing.T) {
	ethOnlyPubkey := core.NewG1Point(big.NewInt(1), big.NewInt(2))
	dualQuorumPubkey := core.NewG1Point(big.NewInt(3), big.NewInt(4))
	ethOnlyID := ethOnlyPubkey.GetOperatorID().Hex()
	dualQuorumID := dualQuorumPubkey.GetOperatorID().Hex()

	operatorQuorumIntervals := dataapi.OperatorQuorumIntervals{
		ethOnlyID: {
			0: {{StartBlock: 1, EndBlock: 10}},
		},
		dualQuorumID: {
			0: {{StartBlock: 1, EndBlock: 10}},
			1: {{StartBlock: 1, EndBlock: 10}},
		},
	}

	profiles := make(map[string]*batchQuorumProfile)
	attestations := []*corev2.Attestation{
		// A phantom ETH quorum in the attestation must not make an ETH-only
		// operator responsible for an EIGEN-only blob.
		newSigningInfoTestBatch(t, 1, 1, [][]core.QuorumID{{1}}, []core.QuorumID{0, 1},
			[]*core.G1Point{ethOnlyPubkey}, profiles),
		// An ETH-only batch is part of the ETH-only operator's denominator.
		newSigningInfoTestBatch(t, 2, 2, [][]core.QuorumID{{0}}, []core.QuorumID{0},
			[]*core.G1Point{ethOnlyPubkey}, profiles),
		// The ETH-only operator cannot validate every blob in this batch.
		newSigningInfoTestBatch(t, 3, 3, [][]core.QuorumID{{0}, {1}}, []core.QuorumID{0, 1},
			[]*core.G1Point{ethOnlyPubkey}, profiles),
		// One multi-quorum blob is valid for an ETH-only operator. Quorum 0 is
		// intentionally absent from the attestation to verify that absence is
		// not treated as a missed signature.
		newSigningInfoTestBatch(t, 4, 4, [][]core.QuorumID{{0, 1}}, []core.QuorumID{1},
			nil, profiles),
		// Empty attestations have no reliable non-signer set and are excluded.
		{
			BatchHeader: &corev2.BatchHeader{BatchRoot: [32]byte{5}, ReferenceBlockNumber: 5},
			AttestedAt:  5,
		},
	}

	numFailed, numResponsible, totalBatches, err := computeBlobQuorumSigningStats(
		attestations,
		profiles,
		operatorQuorumIntervals,
	)
	require.NoError(t, err)

	require.Equal(t, map[uint8]int{0: 3, 1: 3}, totalBatches)
	require.Equal(t, map[uint8]int{0: 2}, numResponsible[ethOnlyID])
	require.Equal(t, map[uint8]int{0: 1}, numFailed[ethOnlyID])
	require.Equal(t, map[uint8]int{0: 3, 1: 3}, numResponsible[dualQuorumID])
	require.NotContains(t, numFailed, dualQuorumID)
}

func TestNewBatchQuorumProfileFromBlobQuorumsValidation(t *testing.T) {
	_, err := newBatchQuorumProfileFromBlobQuorums(nil)
	require.ErrorContains(t, err, "at least one blob quorum set")

	_, err = newBatchQuorumProfileFromBlobQuorums([][]core.QuorumID{{}})
	require.ErrorContains(t, err, "has no quorums")
}

func TestDeduplicateAttestationsKeepsLatestUpdate(t *testing.T) {
	header := &corev2.BatchHeader{BatchRoot: [32]byte{1}, ReferenceBlockNumber: 1}
	attestations := []*corev2.Attestation{
		{BatchHeader: header, AttestedAt: 1},
		{BatchHeader: header, AttestedAt: 3, QuorumNumbers: []core.QuorumID{0}},
		{BatchHeader: header, AttestedAt: 2, QuorumNumbers: []core.QuorumID{0, 1}},
	}

	deduplicated, err := deduplicateAttestations(attestations)
	require.NoError(t, err)
	require.Len(t, deduplicated, 1)
	require.Same(t, attestations[1], deduplicated[0])
}

func newSigningInfoTestBatch(
	t *testing.T,
	batchRoot byte,
	attestedAt uint64,
	blobQuorums [][]core.QuorumID,
	attestationQuorums []core.QuorumID,
	nonSigners []*core.G1Point,
	profiles map[string]*batchQuorumProfile,
) *corev2.Attestation {
	t.Helper()

	header := &corev2.BatchHeader{
		BatchRoot:            [32]byte{batchRoot},
		ReferenceBlockNumber: attestedAt,
	}
	certificates := make([]*corev2.BlobCertificate, len(blobQuorums))
	for i, quorums := range blobQuorums {
		certificates[i] = &corev2.BlobCertificate{
			BlobHeader: &corev2.BlobHeader{QuorumNumbers: quorums},
		}
	}
	profile, err := newBatchQuorumProfile(&corev2.Batch{
		BatchHeader:      header,
		BlobCertificates: certificates,
	})
	require.NoError(t, err)
	batchHeaderHash, err := header.Hash()
	require.NoError(t, err)
	profiles[hex.EncodeToString(batchHeaderHash[:])] = profile

	return &corev2.Attestation{
		BatchHeader:      header,
		AttestedAt:       attestedAt,
		NonSignerPubKeys: nonSigners,
		QuorumNumbers:    attestationQuorums,
	}
}
