package verify

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	grpccommon "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/encoding/v1/kzg"
	kzgverifier "github.com/Layr-Labs/eigenda/encoding/v1/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/v1/rs"
	"github.com/stretchr/testify/require"
)

func TestCommitmentVerification(t *testing.T) {
	t.Parallel()

	var data = []byte("inter-subjective and not objective!")

	x, err := hex.DecodeString("1021d699eac68ce312196d480266e8b82fd5fe5c4311e53313837b64db6df178")
	require.NoError(t, err)

	y, err := hex.DecodeString("02efa5a7813233ae13f32bae9b8f48252fa45c1b06a5d70bed471a9bea8d98ae")
	require.NoError(t, err)

	c := &grpccommon.G1Commitment{
		X: x,
		Y: y,
	}

	kzgConfig := kzg.KzgConfig{
		G1Path:          "../../../../resources/g1.point",
		G2Path:          "../../../../resources/g2.point",
		G2TrailingPath:  "../../../../resources/g2.trailing.point",
		CacheDir:        "../../../../resources/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    false,
	}

	kzgVerifier, err := kzgverifier.NewVerifier(&kzgConfig, nil)
	require.NoError(t, err)

	cfg := &Config{
		VerifyCerts: false,
	}

	v, err := NewVerifier(cfg, kzgVerifier, nil)
	require.NoError(t, err)

	// Happy path verification
	codec := codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec())
	blob, err := codec.EncodeBlob(data)
	require.NoError(t, err)
	err = v.VerifyCommitment(c, blob)
	require.NoError(t, err)

	// failure with wrong data
	fakeData, err := codec.EncodeBlob([]byte("I am an imposter!!"))
	require.NoError(t, err)
	err = v.VerifyCommitment(c, fakeData)
	require.Error(t, err)
}

func TestCommitmentWithTooLargeBlob(t *testing.T) {

	var dataRand [2000 * 32]byte
	_, err := rand.Read(dataRand[:])
	require.NoError(t, err)
	data := dataRand[:]

	kzgConfig := kzg.KzgConfig{
		G1Path:          "../../../../resources/g1.point",
		G2Path:          "../../../../resources/g2.point",
		G2TrailingPath:  "../../../../resources/g2.trailing.point",
		CacheDir:        "../../../../resources/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    false,
	}

	kzgVerifier, err := kzgverifier.NewVerifier(&kzgConfig, nil)
	require.NoError(t, err)

	cfg := &Config{
		VerifyCerts: false,
	}

	v, err := NewVerifier(cfg, kzgVerifier, nil)
	require.NoError(t, err)

	// Some wrong commitment just to pass in function
	x, err := hex.DecodeString("1021d699eac68ce312196d480266e8b82fd5fe5c4311e53313837b64db6df178")
	require.NoError(t, err)

	y, err := hex.DecodeString("02efa5a7813233ae13f32bae9b8f48252fa45c1b06a5d70bed471a9bea8d98ae")
	require.NoError(t, err)

	c := &grpccommon.G1Commitment{
		X: x,
		Y: y,
	}

	// Happy path verification
	codec := codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec())
	blob, err := codec.EncodeBlob(data)
	require.NoError(t, err)

	inputFr, err := rs.ToFrArray(blob)
	require.NoError(t, err)

	err = v.VerifyCommitment(c, blob)
	msg := fmt.Sprintf(
		"cannot verify commitment because the number of stored srs in the memory is insufficient, have %v need %v",
		kzgConfig.SRSNumberToLoad,
		len(inputFr),
	)
	require.EqualError(t, err, msg)

}

func TestVerifySecurityParams(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		blobHeader  BlobHeader
		batchHeader *disperser.BatchHeader
		setupCV     func() *CertVerifier
		holesky     bool
		expectError bool
		errorMsg    string
	}{
		{
			name: "blob has more quorum parameters than available quorums",
			blobHeader: BlobHeader{
				QuorumBlobParams: []QuorumBlobParam{
					{QuorumNumber: 0, AdversaryThresholdPercentage: 33, ConfirmationThresholdPercentage: 50},
					{QuorumNumber: 1, AdversaryThresholdPercentage: 33, ConfirmationThresholdPercentage: 50},
					{QuorumNumber: 2, AdversaryThresholdPercentage: 33, ConfirmationThresholdPercentage: 50},
				},
			},
			batchHeader: &disperser.BatchHeader{
				QuorumNumbers: []byte{0, 1}, // Only 2 quorums available
			},
			setupCV:     func() *CertVerifier { return nil },
			expectError: true,
			errorMsg:    "blob has more quorum parameters than available quorums: got 3 quorum params, available quorums: 2",
		},
		{
			name: "equal number of quorum parameters and available quorums - valid",
			blobHeader: BlobHeader{
				QuorumBlobParams: []QuorumBlobParam{
					{QuorumNumber: 0, AdversaryThresholdPercentage: 33, ConfirmationThresholdPercentage: 50},
					{QuorumNumber: 1, AdversaryThresholdPercentage: 33, ConfirmationThresholdPercentage: 50},
				},
			},
			batchHeader: &disperser.BatchHeader{
				QuorumNumbers:           []byte{0, 1},
				QuorumSignedPercentages: []byte{60, 60}, // Above confirmation threshold
				ReferenceBlockNumber:    3000000,
			},
			setupCV: func() *CertVerifier {
				return &CertVerifier{
					quorumAdversaryThresholds: map[uint8]uint8{0: 33, 1: 33},
					quorumsRequired:           []uint8{0, 1},
				}
			},
			expectError: false,
		},
		{
			name: "fewer quorum parameters than available quorums - valid when required quorums are met",
			blobHeader: BlobHeader{
				QuorumBlobParams: []QuorumBlobParam{
					{QuorumNumber: 0, AdversaryThresholdPercentage: 33, ConfirmationThresholdPercentage: 50},
				},
			},
			batchHeader: &disperser.BatchHeader{
				QuorumNumbers:           []byte{0, 1, 2},
				QuorumSignedPercentages: []byte{60, 60, 60},
				ReferenceBlockNumber:    3000000,
			},
			setupCV: func() *CertVerifier {
				return &CertVerifier{
					quorumAdversaryThresholds: map[uint8]uint8{0: 33},
					quorumsRequired:           []uint8{0},
				}
			},
			expectError: false,
		},
		{
			name: "quorum number mismatch",
			blobHeader: BlobHeader{
				QuorumBlobParams: []QuorumBlobParam{
					{QuorumNumber: 1, AdversaryThresholdPercentage: 33, ConfirmationThresholdPercentage: 50},
				},
			},
			batchHeader: &disperser.BatchHeader{
				QuorumNumbers:           []byte{0}, // Mismatch: expects 0, got 1
				QuorumSignedPercentages: []byte{60},
				ReferenceBlockNumber:    3000000,
			},
			setupCV: func() *CertVerifier {
				return &CertVerifier{
					quorumAdversaryThresholds: map[uint8]uint8{0: 33},
					quorumsRequired:           []uint8{0},
				}
			},
			expectError: true,
			errorMsg:    "quorum number mismatch, expected: 0, got: 1",
		},
		{
			name: "adversary threshold exceeds confirmation threshold",
			blobHeader: BlobHeader{
				QuorumBlobParams: []QuorumBlobParam{
					{QuorumNumber: 0, AdversaryThresholdPercentage: 60, ConfirmationThresholdPercentage: 50},
				},
			},
			batchHeader: &disperser.BatchHeader{
				QuorumNumbers:           []byte{0},
				QuorumSignedPercentages: []byte{70},
				ReferenceBlockNumber:    3000000,
			},
			setupCV: func() *CertVerifier {
				return &CertVerifier{
					quorumAdversaryThresholds: map[uint8]uint8{0: 33},
					quorumsRequired:           []uint8{0},
				}
			},
			expectError: true,
			errorMsg:    "adversary threshold percentage must be greater than or equal to confirmation threshold percentage",
		},
		{
			name: "adversary threshold below quorum adversary threshold",
			blobHeader: BlobHeader{
				QuorumBlobParams: []QuorumBlobParam{
					{QuorumNumber: 0, AdversaryThresholdPercentage: 25, ConfirmationThresholdPercentage: 50},
				},
			},
			batchHeader: &disperser.BatchHeader{
				QuorumNumbers:           []byte{0},
				QuorumSignedPercentages: []byte{60},
				ReferenceBlockNumber:    3000000,
			},
			setupCV: func() *CertVerifier {
				return &CertVerifier{
					quorumAdversaryThresholds: map[uint8]uint8{0: 33},
					quorumsRequired:           []uint8{0},
				}
			},
			expectError: true,
			errorMsg:    "adversary threshold percentage must be >= quorum adversary threshold percentage",
		},
		{
			name: "signed stake below confirmation threshold",
			blobHeader: BlobHeader{
				QuorumBlobParams: []QuorumBlobParam{
					{QuorumNumber: 0, AdversaryThresholdPercentage: 33, ConfirmationThresholdPercentage: 50},
				},
			},
			batchHeader: &disperser.BatchHeader{
				QuorumNumbers:           []byte{0},
				QuorumSignedPercentages: []byte{40}, // Below confirmation threshold of 50
				ReferenceBlockNumber:    3000000,
			},
			setupCV: func() *CertVerifier {
				return &CertVerifier{
					quorumAdversaryThresholds: map[uint8]uint8{0: 33},
					quorumsRequired:           []uint8{0},
				}
			},
			expectError: true,
			errorMsg:    "signed stake for quorum must be >= to confirmation threshold percentage",
		},
		{
			name: "required quorum not present in confirmed quorums",
			blobHeader: BlobHeader{
				QuorumBlobParams: []QuorumBlobParam{
					{QuorumNumber: 0, AdversaryThresholdPercentage: 33, ConfirmationThresholdPercentage: 50},
				},
			},
			batchHeader: &disperser.BatchHeader{
				QuorumNumbers:           []byte{0, 1},
				QuorumSignedPercentages: []byte{60, 60},
				ReferenceBlockNumber:    3000000,
			},
			setupCV: func() *CertVerifier {
				return &CertVerifier{
					quorumAdversaryThresholds: map[uint8]uint8{0: 33, 1: 33},
					quorumsRequired:           []uint8{0, 1}, // Requires both 0 and 1
				}
			},
			expectError: true,
			errorMsg:    "quorum 1 is required but not present in confirmed quorums",
		},
		{
			name: "holesky special case - only quorum 0 required for specific block range",
			blobHeader: BlobHeader{
				QuorumBlobParams: []QuorumBlobParam{
					{QuorumNumber: 0, AdversaryThresholdPercentage: 33, ConfirmationThresholdPercentage: 50},
				},
			},
			batchHeader: &disperser.BatchHeader{
				QuorumNumbers:           []byte{0},
				QuorumSignedPercentages: []byte{60},
				ReferenceBlockNumber:    2955000, // In the special block range [2950000, 2960000)
			},
			setupCV: func() *CertVerifier {
				return &CertVerifier{
					quorumAdversaryThresholds: map[uint8]uint8{0: 33, 1: 33},
					quorumsRequired:           []uint8{0, 1}, // Would normally require both
				}
			},
			holesky:     true,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Create a minimal verifier with just the CertVerifier set
			v := &Verifier{
				cv: tc.setupCV(),
			}

			err := v.verifySecurityParams(tc.blobHeader, tc.batchHeader)

			if tc.expectError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
