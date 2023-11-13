package core_test

import (
	"encoding/json"
	"math/big"
	"testing"

	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/core"
	kzgbn254 "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

const (
	encodedBatchHeader     = "0x31000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001"
	reducedBatchHeaderHash = "0x891d0936da4627f445ef193aad63afb173409af9e775e292e4e35aff790a45e2"
	batchHeaderHash        = "0xa48219ff51a67bf779c6f7858e3bf9760ef10a766e5dc5d461318c8e9d5607b6"
	encodedBlobHeader      = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000005000000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000014"
	blobHeaderHash         = "0x48b3e6540820e2f2c185764e22d438d5ff03551a867299b26cccf57fa2c3f237"
)

func TestBatchHeaderEncoding(t *testing.T) {
	batchRoot := [32]byte{}
	copy(batchRoot[:], []byte("1"))
	batchHeader := &core.BatchHeader{
		ReferenceBlockNumber: 1,
		BatchRoot:            batchRoot,
	}

	data, err := batchHeader.Encode()
	assert.NoError(t, err)
	assert.Equal(t, hexutil.Encode(data), encodedBatchHeader)

	hash, err := batchHeader.GetBatchHeaderHash()
	assert.NoError(t, err)
	assert.Equal(t, hexutil.Encode(hash[:]), reducedBatchHeaderHash)

	onchainBatchHeader := binding.IEigenDAServiceManagerBatchHeader{
		BlobHeadersRoot:            batchRoot,
		QuorumNumbers:              []byte{0},
		QuorumThresholdPercentages: []byte{100},
		ReferenceBlockNumber:       1,
	}
	hash, err = core.HashBatchHeader(onchainBatchHeader)
	assert.NoError(t, err)
	assert.Equal(t, hexutil.Encode(hash[:]), batchHeaderHash)
}

func TestBlobHeaderEncoding(t *testing.T) {

	var commitX, commitY, lengthX, lengthY fp.Element
	commitX = *commitX.SetBigInt(big.NewInt(1))
	commitY = *commitY.SetBigInt(big.NewInt(2))
	lengthX = *lengthX.SetBigInt(big.NewInt(1))
	lengthY = *lengthY.SetBigInt(big.NewInt(2))

	commitment := &kzgbn254.G1Point{
		X: commitX,
		Y: commitY,
	}
	lengthProof := &kzgbn254.G1Point{
		X: lengthX,
		Y: lengthY,
	}
	blobHeader := &core.BlobHeader{
		BlobCommitments: core.BlobCommitments{
			Commitment: &core.Commitment{
				commitment,
			},
			LengthProof: &core.Commitment{
				lengthProof,
			},
			Length: 10,
		},
		QuorumInfos: []*core.BlobQuorumInfo{
			{
				SecurityParam: core.SecurityParam{
					QuorumID:           1,
					AdversaryThreshold: 80,
					QuorumThreshold:    100,
				},
				QuantizationFactor: 20,
			},
		},
	}
	data, err := blobHeader.Encode()
	assert.NoError(t, err)
	assert.Equal(t, encodedBlobHeader, hexutil.Encode(data))

	h, err := blobHeader.GetBlobHeaderHash()
	assert.NoError(t, err)
	assert.Equal(t, blobHeaderHash, hexutil.Encode(h[:]))
}

func TestSignatoryRecord(t *testing.T) {

	var X1, Y1, X2, Y2 fp.Element
	X1 = *X1.SetBigInt(big.NewInt(1))
	Y1 = *Y1.SetBigInt(big.NewInt(2))
	X2 = *X2.SetBigInt(big.NewInt(3))
	Y2 = *Y2.SetBigInt(big.NewInt(4))

	key1 := &core.G1Point{
		G1Affine: &bn254.G1Affine{
			X: X1,
			Y: Y1,
		},
	}
	key2 := &core.G1Point{
		G1Affine: &bn254.G1Affine{
			X: X2,
			Y: Y2,
		},
	}

	operatorID1 := key1.GetOperatorID()
	operatorID2 := key2.GetOperatorID()
	assert.Equal(t, common.Bytes2Hex(operatorID1[:]), "e90b7bceb6e7df5418fb78d8ee546e97c83a08bbccc01a0644d599ccd2a7c2e0")
	assert.Equal(t, common.Bytes2Hex(operatorID2[:]), "2e174c10e159ea99b867ce3205125c24a42d128804e4070ed6fcc8cc98166aa0")
	hash := core.ComputeSignatoryRecordHash(123, []*core.G1Point{
		key1, key2,
	})

	expected := "f60f497b0f816a24c750d818c538f7eb2131a6c3bf487053042914021a671023"
	assert.Equal(t, common.Bytes2Hex(hash[:]), expected)
}

func TestCommitmentMarshaling(t *testing.T) {

	var commitX, commitY fp.Element
	commitX = *commitX.SetBigInt(big.NewInt(1))
	commitY = *commitY.SetBigInt(big.NewInt(2))

	commitment := &core.Commitment{
		G1Point: &kzgbn254.G1Point{
			X: commitX,
			Y: commitY,
		},
	}

	marshalled, err := json.Marshal(commitment)
	assert.NoError(t, err)

	recovered := new(core.Commitment)
	err = json.Unmarshal(marshalled, recovered)
	assert.NoError(t, err)
	assert.Equal(t, recovered, commitment)
}

func TestQuorumParamsHash(t *testing.T) {
	blobHeader := &core.BlobHeader{
		QuorumInfos: []*core.BlobQuorumInfo{
			{
				SecurityParam: core.SecurityParam{
					QuorumID:           0,
					AdversaryThreshold: 80,
					QuorumThreshold:    100,
				},
				QuantizationFactor: 1,
				EncodedBlobLength:  48,
			},
		},
	}
	hash, err := blobHeader.GetQuorumBlobParamsHash()
	assert.NoError(t, err)
	expected := "90a8cc415c00b8bc3dcc3b21f240277e93ef712327e0001094b045ec60dff65c"
	assert.Equal(t, common.Bytes2Hex(hash[:]), expected)
}
