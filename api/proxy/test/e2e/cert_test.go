package e2e

import (
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
	certbindings "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
)

func TestCert(t *testing.T) {
	certs.NewVersionedCert([]byte{0}, certs.V2VersionByte)

	cert := coretypes.EigenDACertV3{
		BatchHeader: certbindings.EigenDATypesV2BatchHeaderV2{
			BatchRoot:            [32]byte{}, // 32 zero bytes
			ReferenceBlockNumber: 0,
		},
		BlobInclusionInfo: certbindings.EigenDATypesV2BlobInclusionInfo{
			BlobCertificate: certbindings.EigenDATypesV2BlobCertificate{
				BlobHeader: certbindings.EigenDATypesV2BlobHeaderV2{
					Version:       0,
					QuorumNumbers: []byte{}, // Empty slice
					Commitment: certbindings.EigenDATypesV2BlobCommitment{
						Commitment: certbindings.BN254G1Point{
							X: big.NewInt(0),
							Y: big.NewInt(0),
						},
						LengthCommitment: certbindings.BN254G2Point{
							X: [2]*big.Int{big.NewInt(0), big.NewInt(0)},
							Y: [2]*big.Int{big.NewInt(0), big.NewInt(0)},
						},
						LengthProof: certbindings.BN254G2Point{
							X: [2]*big.Int{big.NewInt(0), big.NewInt(0)},
							Y: [2]*big.Int{big.NewInt(0), big.NewInt(0)},
						},
						Length: 0,
					},
					PaymentHeaderHash: [32]byte{}, // 32 zero bytes
				},
				Signature: []byte{},   // Empty slice
				RelayKeys: []uint32{}, // Empty slice
			},
			BlobIndex:      0,
			InclusionProof: []byte{}, // Empty slice
		},
		NonSignerStakesAndSignature: certbindings.EigenDATypesV1NonSignerStakesAndSignature{
			NonSignerQuorumBitmapIndices: []uint32{},                    // Empty slice
			NonSignerPubkeys:             []certbindings.BN254G1Point{}, // Empty slice
			QuorumApks:                   []certbindings.BN254G1Point{}, // Empty slice
			ApkG2: certbindings.BN254G2Point{
				X: [2]*big.Int{big.NewInt(0), big.NewInt(0)},
				Y: [2]*big.Int{big.NewInt(0), big.NewInt(0)},
			},
			Sigma: certbindings.BN254G1Point{
				X: big.NewInt(0),
				Y: big.NewInt(0),
			},
			QuorumApkIndices:      []uint32{},   // Empty slice
			TotalStakeIndices:     []uint32{},   // Empty slice
			NonSignerStakeIndices: [][]uint32{}, // Empty slice of slices
		},
		SignedQuorumNumbers: []byte{}, // Empty slice
	}
	_ = cert
}
