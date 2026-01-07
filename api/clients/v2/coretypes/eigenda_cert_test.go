package coretypes_test

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/v2/coretypes"
	contractEigenDACertVerifierV2 "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifierV2"
	certTypesBinding "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDACertTypeBindings"
	coreV2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/require"
)

// TestEigenDACertV3_RLPEncodeDecode tests that V3 certificates can be RLP encoded and decoded successfully
func TestEigenDACertV3_RLPEncodeDecode(t *testing.T) {
	// Create a sample V3 certificate
	cert := createSampleEigenDACertV3()

	// Serialize using RLP
	encoded, err := cert.Serialize(coretypes.CertSerializationRLP)
	require.NoError(t, err)
	require.NotEmpty(t, encoded)

	// Deserialize using RLP
	decoded, err := coretypes.DeserializeEigenDACertV3(encoded, coretypes.CertSerializationRLP)
	require.NoError(t, err)
	require.NotNil(t, decoded)

	// Verify the decoded certificate matches the original
	assertCertV3Equal(t, cert, decoded)
}

// TestEigenDACertV3_ABIEncodeDecode tests that V3 certificates can be ABI encoded and decoded successfully
func TestEigenDACertV3_ABIEncodeDecode(t *testing.T) {
	// Create a sample V3 certificate
	cert := createSampleEigenDACertV3()

	// Serialize using ABI
	encoded, err := cert.Serialize(coretypes.CertSerializationABI)
	require.NoError(t, err)
	require.NotEmpty(t, encoded)

	// Deserialize using ABI
	decoded, err := coretypes.DeserializeEigenDACertV3(encoded, coretypes.CertSerializationABI)
	require.NoError(t, err)
	require.NotNil(t, decoded)

	// Verify the decoded certificate matches the original
	assertCertV3Equal(t, cert, decoded)
}

// TestDeserializeEigenDACert tests the generic deserialization function
func TestDeserializeEigenDACert(t *testing.T) {
	tests := []struct {
		name        string
		version     coretypes.CertificateVersion
		createCert  func() coretypes.EigenDACert
		serialType  coretypes.CertSerializationType
		shouldError bool
	}{
		{
			name:        "V3 RLP",
			version:     coretypes.VersionThreeCert,
			createCert:  func() coretypes.EigenDACert { return createSampleEigenDACertV3() },
			serialType:  coretypes.CertSerializationRLP,
			shouldError: false,
		},
		{
			name:        "V3 ABI",
			version:     coretypes.VersionThreeCert,
			createCert:  func() coretypes.EigenDACert { return createSampleEigenDACertV3() },
			serialType:  coretypes.CertSerializationABI,
			shouldError: false,
		},
		{
			name:        "V2 ABI",
			version:     coretypes.VersionTwoCert,
			createCert:  func() coretypes.EigenDACert { return createSampleEigenDACertV2() },
			shouldError: false,
		},
		{
			name:        "Unsupported version",
			version:     0xFF,
			createCert:  func() coretypes.EigenDACert { return createSampleEigenDACertV3() },
			serialType:  coretypes.CertSerializationRLP,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldError {
				// Test unsupported version
				_, err := coretypes.DeserializeEigenDACert([]byte{}, tt.version, tt.serialType)
				require.Error(t, err)
				require.Contains(t, err.Error(), "unsupported certificate version")
				return
			}

			cert := tt.createCert()
			encoded, err := cert.Serialize(tt.serialType)
			require.NoError(t, err)

			decoded, err := coretypes.DeserializeEigenDACert(encoded, tt.version, tt.serialType)
			require.NoError(t, err)
			require.NotNil(t, decoded)
		})
	}
}

// Helper functions to create sample certificates for testing
func createSampleEigenDACertV2() *coretypes.EigenDACertV2 {
	return &coretypes.EigenDACertV2{
		BlobInclusionInfo: contractEigenDACertVerifierV2.EigenDATypesV2BlobInclusionInfo{
			BlobCertificate: contractEigenDACertVerifierV2.EigenDATypesV2BlobCertificate{
				BlobHeader: contractEigenDACertVerifierV2.EigenDATypesV2BlobHeaderV2{
					Version:       1,
					QuorumNumbers: []byte{0, 1},
					Commitment: contractEigenDACertVerifierV2.EigenDATypesV2BlobCommitment{
						Commitment: contractEigenDACertVerifierV2.BN254G1Point{
							X: big.NewInt(12345),
							Y: big.NewInt(67890),
						},
						LengthCommitment: contractEigenDACertVerifierV2.BN254G2Point{
							X: [2]*big.Int{big.NewInt(111), big.NewInt(222)},
							Y: [2]*big.Int{big.NewInt(333), big.NewInt(444)},
						},
						LengthProof: contractEigenDACertVerifierV2.BN254G2Point{
							X: [2]*big.Int{big.NewInt(555), big.NewInt(666)},
							Y: [2]*big.Int{big.NewInt(777), big.NewInt(888)},
						},
						Length: 1024,
					},
					PaymentHeaderHash: [32]byte{1, 2, 3},
				},
				Signature: []byte{10, 20, 30},
				RelayKeys: []coreV2.RelayKey{1, 2, 3},
			},
			BlobIndex:      5,
			InclusionProof: []byte{40, 50, 60},
		},
		BatchHeader: contractEigenDACertVerifierV2.EigenDATypesV2BatchHeaderV2{
			BatchRoot:            [32]byte{4, 5, 6},
			ReferenceBlockNumber: 12345,
		},
		NonSignerStakesAndSignature: contractEigenDACertVerifierV2.EigenDATypesV1NonSignerStakesAndSignature{
			NonSignerQuorumBitmapIndices: []uint32{0, 1},
			NonSignerPubkeys: []contractEigenDACertVerifierV2.BN254G1Point{
				{X: big.NewInt(100), Y: big.NewInt(200)},
			},
			QuorumApks: []contractEigenDACertVerifierV2.BN254G1Point{
				{X: big.NewInt(300), Y: big.NewInt(400)},
			},
			ApkG2: contractEigenDACertVerifierV2.BN254G2Point{
				X: [2]*big.Int{big.NewInt(500), big.NewInt(600)},
				Y: [2]*big.Int{big.NewInt(700), big.NewInt(800)},
			},
			Sigma: contractEigenDACertVerifierV2.BN254G1Point{
				X: big.NewInt(900),
				Y: big.NewInt(1000),
			},
			QuorumApkIndices:      []uint32{0},
			TotalStakeIndices:     []uint32{0},
			NonSignerStakeIndices: [][]uint32{{0}},
		},
		SignedQuorumNumbers: []byte{0, 1},
	}
}

func createSampleEigenDACertV3() *coretypes.EigenDACertV3 {
	return &coretypes.EigenDACertV3{
		BlobInclusionInfo: certTypesBinding.EigenDATypesV2BlobInclusionInfo{
			BlobCertificate: certTypesBinding.EigenDATypesV2BlobCertificate{
				BlobHeader: certTypesBinding.EigenDATypesV2BlobHeaderV2{
					Version:       1,
					QuorumNumbers: []byte{0, 1},
					Commitment: certTypesBinding.EigenDATypesV2BlobCommitment{
						Commitment: certTypesBinding.BN254G1Point{
							X: big.NewInt(12345),
							Y: big.NewInt(67890),
						},
						LengthCommitment: certTypesBinding.BN254G2Point{
							X: [2]*big.Int{big.NewInt(111), big.NewInt(222)},
							Y: [2]*big.Int{big.NewInt(333), big.NewInt(444)},
						},
						LengthProof: certTypesBinding.BN254G2Point{
							X: [2]*big.Int{big.NewInt(555), big.NewInt(666)},
							Y: [2]*big.Int{big.NewInt(777), big.NewInt(888)},
						},
						Length: 1024,
					},
					PaymentHeaderHash: [32]byte{1, 2, 3},
				},
				Signature: []byte{10, 20, 30},
				RelayKeys: []coreV2.RelayKey{1, 2, 3},
			},
			BlobIndex:      5,
			InclusionProof: []byte{40, 50, 60},
		},
		BatchHeader: certTypesBinding.EigenDATypesV2BatchHeaderV2{
			BatchRoot:            [32]byte{4, 5, 6},
			ReferenceBlockNumber: 12345,
		},
		NonSignerStakesAndSignature: certTypesBinding.EigenDATypesV1NonSignerStakesAndSignature{
			NonSignerQuorumBitmapIndices: []uint32{0, 1},
			NonSignerPubkeys: []certTypesBinding.BN254G1Point{
				{X: big.NewInt(100), Y: big.NewInt(200)},
			},
			QuorumApks: []certTypesBinding.BN254G1Point{
				{X: big.NewInt(300), Y: big.NewInt(400)},
			},
			ApkG2: certTypesBinding.BN254G2Point{
				X: [2]*big.Int{big.NewInt(500), big.NewInt(600)},
				Y: [2]*big.Int{big.NewInt(700), big.NewInt(800)},
			},
			Sigma: certTypesBinding.BN254G1Point{
				X: big.NewInt(900),
				Y: big.NewInt(1000),
			},
			QuorumApkIndices:      []uint32{0},
			TotalStakeIndices:     []uint32{0},
			NonSignerStakeIndices: [][]uint32{{0}},
		},
		SignedQuorumNumbers: []byte{0, 1},
	}
}

func assertCertV3Equal(t *testing.T, expected, actual *coretypes.EigenDACertV3) {
	require.True(t, reflect.DeepEqual(expected, actual))
}
