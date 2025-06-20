package commitments

import (
	"fmt"

	"github.com/Layr-Labs/eigenda-proxy/common/types/certs"
)

type CommitmentMode string

const (
	OptimismKeccakCommitmentMode  CommitmentMode = "optimism_keccak256"
	OptimismGenericCommitmentMode CommitmentMode = "optimism_generic"
	StandardCommitmentMode        CommitmentMode = "standard"
)

// EncodeCommitment serializes the versionedCert prepends commitmentMode-related header bytes.
// The returned byte array is the final "commitment" which is returned to POST requests,
// and can be passed back to the same-mode GET routes to retrieve the original payload.
// The commitment is so called because it is typically sent as-is (or with an extra additional byte in the case of op)
// to the batcher inbox, as an "altda commitment".
// See https://specs.optimism.io/experimental/alt-da.html#input-commitment-submission
//
// See the Encode() function of each commitment type for more details on each encoding:
// standard mode: no extra prefixed bytes
// op keccak mode: 0x00 prefix byte
// op generic mode: 0x01 + 0x00 prefix bytes
func EncodeCommitment(
	versionedCert certs.VersionedCert,
	commitmentMode CommitmentMode,
) ([]byte, error) {
	switch commitmentMode {
	case OptimismKeccakCommitmentMode:
		return OPKeccak256Commitment(versionedCert.SerializedCert).Encode(), nil
	case OptimismGenericCommitmentMode:
		// Proxy returns an altDACommitment, which doesn't contain the first op version_byte
		// (from https://specs.optimism.io/experimental/alt-da.html#example-commitments)
		// This is because the version_byte is added by op-alt-da when calling TxData() right before submitting the tx:
		// https://github.com/Layr-Labs/optimism/blob/89ac40d0fddba2e06854b253b9f0266f36350af2/op-alt-da/commitment.go#L158-L160
		return NewOPEigenDAGenericCommitment(versionedCert).Encode(), nil
	case StandardCommitmentMode:
		return NewStandardCommitment(versionedCert).Encode(), nil
	}
	return nil, fmt.Errorf("unknown commitment mode")
}
