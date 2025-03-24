package verify

import (
	"math/big"
	"testing"

	eigenda_common "github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/api/grpc/disperser"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDACertVerifier"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

func TestHashBatchHashedMetadata(t *testing.T) {
	batchHeaderHash := crypto.Keccak256Hash([]byte("batchHeader"))
	sigRecordHash := crypto.Keccak256Hash([]byte("signatoryRecord"))

	// 1 - Test using uint32 MAX
	var blockNum uint32 = 4294967295

	expected := "0x687b60d8b30b6aaddf6413728fb66fb7a7554601c2cc8e17a37fa94ad0818500"
	actual, err := HashBatchHashedMetadata(batchHeaderHash, sigRecordHash, blockNum)
	require.NoError(t, err)

	require.Equal(t, expected, actual.String())

	// 2 - Test using uint32 value
	blockNum = 4294967294

	expected = "0x94d77be4d3d180d32d61ec8037e687b71e7996feded39b72a6dc3f9ff6406b30"
	actual, err = HashBatchHashedMetadata(batchHeaderHash, sigRecordHash, blockNum)
	require.NoError(t, err)

	require.Equal(t, expected, actual.String())

	// 3 - Testing using uint32 0 value
	blockNum = 0

	expected = "0x482dfb1545a792b6d118a045033143d0cc28b0e5a4b2e1924decf27e4fc8c250"
	actual, err = HashBatchHashedMetadata(batchHeaderHash, sigRecordHash, blockNum)
	require.NoError(t, err)

	require.Equal(t, expected, actual.String())
}

func TestHashBatchMetadata(t *testing.T) {
	testHash := crypto.Keccak256Hash([]byte("batchHeader"))

	header := &binding.BatchHeader{
		BlobHeadersRoot:       testHash,
		QuorumNumbers:         testHash.Bytes(),
		SignedStakeForQuorums: testHash.Bytes(),
		ReferenceBlockNumber:  1,
	}

	expected := "0x746f8a453586621d12e41d097eab089b1f25beca44c434281d68d4be0484b7e8"

	actual, err := HashBatchMetadata(header, testHash, 1)
	require.NoError(t, err)
	require.Equal(t, actual.String(), expected)

}

func TestHashBlobHeader(t *testing.T) {
	expected := "0xba4675a31c9bf6b2f7abfdcedd34b74645cb7332b35db39bff00ae8516a67393"

	// [[1,1],2,[[2,4,5,6]]]
	header := &disperser.BlobHeader{
		Commitment: &eigenda_common.G1Commitment{
			X: big.NewInt(1).Bytes(),
			Y: big.NewInt(1).Bytes(),
		},
		DataLength: 2,
		BlobQuorumParams: []*disperser.BlobQuorumParam{
			{
				QuorumNumber:                    2,
				AdversaryThresholdPercentage:    4,
				ConfirmationThresholdPercentage: 5,
				ChunkLength:                     6,
			},
			{
				QuorumNumber:                    2,
				AdversaryThresholdPercentage:    4,
				ConfirmationThresholdPercentage: 5,
				ChunkLength:                     6,
			},
		},
	}

	cert := &Certificate{
		BlobHeader: header,
	}

	actual, err := HashBlobHeader(cert.ReadBlobHeader())

	require.NoError(t, err)
	require.Equal(t, expected, actual.String())
}

func TestHashEncodeBlobHeader(t *testing.T) {
	expected := "0xf15f43fa44bae9b74cd2f88f8f838e09ff7ab5d50f2170f07b98479eb7da98ba"

	// [[1,1],2,[[2,4,5,6]]]
	header := &disperser.BlobHeader{
		Commitment: &eigenda_common.G1Commitment{
			X: big.NewInt(1).Bytes(),
			Y: big.NewInt(1).Bytes(),
		},
		DataLength: 2,
		BlobQuorumParams: []*disperser.BlobQuorumParam{
			{
				QuorumNumber:                    2,
				AdversaryThresholdPercentage:    4,
				ConfirmationThresholdPercentage: 5,
				ChunkLength:                     6,
			},
			{
				QuorumNumber:                    2,
				AdversaryThresholdPercentage:    4,
				ConfirmationThresholdPercentage: 5,
				ChunkLength:                     6,
			},
		},
	}

	cert := &Certificate{
		BlobHeader: header,
	}

	actual, err := HashEncodeBlobHeader(cert.ReadBlobHeader())

	require.NoError(t, err)
	require.Equal(t, expected, actual.String())
}
