package core_test

import (
	"encoding/json"
	"math/big"
	"testing"

	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/encoding"
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
	encodedBlobHeader      = "0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000000064000000000000000000000000000000000000000000000000000000000000000a"
	blobHeaderHash         = "0xd14b018fcb05ce94b21782c5d3a9c469cb8fcf66926139fee11ceaf0ab7d7c11"
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

	onchainBatchHeader := binding.BatchHeader{
		BlobHeadersRoot:       batchRoot,
		QuorumNumbers:         []byte{0},
		SignedStakeForQuorums: []byte{100},
		ReferenceBlockNumber:  1,
	}
	hash, err = core.HashBatchHeader(onchainBatchHeader)
	assert.NoError(t, err)
	assert.Equal(t, hexutil.Encode(hash[:]), batchHeaderHash)
}

func TestBlobHeaderEncoding(t *testing.T) {

	var commitX, commitY fp.Element
	commitX = *commitX.SetBigInt(big.NewInt(1))
	commitY = *commitY.SetBigInt(big.NewInt(2))

	commitment := &encoding.G1Commitment{
		X: commitX,
		Y: commitY,
	}

	var lengthXA0, lengthXA1, lengthYA0, lengthYA1 fp.Element
	_, err := lengthXA0.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	assert.NoError(t, err)
	_, err = lengthXA1.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	assert.NoError(t, err)
	_, err = lengthYA0.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	assert.NoError(t, err)
	_, err = lengthYA1.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531")
	assert.NoError(t, err)

	var lengthProof, lengthCommitment bn254.G2Affine
	lengthProof.X.A0 = lengthXA0
	lengthProof.X.A1 = lengthXA1
	lengthProof.Y.A0 = lengthYA0
	lengthProof.Y.A1 = lengthYA1

	lengthCommitment = lengthProof

	blobHeader := &core.BlobHeader{
		BlobCommitments: encoding.BlobCommitments{
			Commitment:       commitment,
			LengthCommitment: (*encoding.G2Commitment)(&lengthCommitment),
			LengthProof:      (*encoding.G2Commitment)(&lengthProof),
			Length:           10,
		},
		QuorumInfos: []*core.BlobQuorumInfo{
			{
				SecurityParam: core.SecurityParam{
					QuorumID:              1,
					AdversaryThreshold:    80,
					ConfirmationThreshold: 100,
				},
				ChunkLength: 10,
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

	commitment := &encoding.G1Commitment{
		X: commitX,
		Y: commitY,
	}

	marshalled, err := json.Marshal(commitment)
	assert.NoError(t, err)

	recovered := new(encoding.G1Commitment)
	err = json.Unmarshal(marshalled, recovered)
	assert.NoError(t, err)
	assert.Equal(t, recovered, commitment)
}

func TestQuorumParamsHash(t *testing.T) {
	blobHeader := &core.BlobHeader{
		QuorumInfos: []*core.BlobQuorumInfo{
			{
				SecurityParam: core.SecurityParam{
					QuorumID:              0,
					AdversaryThreshold:    80,
					ConfirmationThreshold: 100,
				},
				ChunkLength: 10,
			},
		},
	}
	hash, err := blobHeader.GetQuorumBlobParamsHash()
	assert.NoError(t, err)
	expected := "89b336cf7ea7dcd13e275b541843175165a1f7dd94ddfa82282be3d7ab402ba2"
	assert.Equal(t, common.Bytes2Hex(hash[:]), expected)
}

func TestHashPubKeyG1(t *testing.T) {
	x, ok := new(big.Int).SetString("166951537990155304646296676950704619272379920143528795571830693741626950865", 10)
	assert.True(t, ok)
	y, ok := new(big.Int).SetString("1787567470127357668828096785064424339221076501074969235378695359686742067296", 10)
	assert.True(t, ok)
	pk := &core.G1Point{
		G1Affine: &bn254.G1Affine{
			X: *new(fp.Element).SetBigInt(x),
			Y: *new(fp.Element).SetBigInt(y),
		},
	}
	hash := eth.HashPubKeyG1(pk)
	assert.Equal(t, common.Bytes2Hex(hash[:]), "426d1a0363fbdcd0c8d33b643252164057193ca022958fa0da99d9e70c980dd7")
}

func TestParseOperatorSocket(t *testing.T) {
	operatorSocket := "localhost:1234;5678"
	host, dispersalPort, retrievalPort, err := core.ParseOperatorSocket(operatorSocket)
	assert.NoError(t, err)
	assert.Equal(t, "localhost", host)
	assert.Equal(t, "1234", dispersalPort)
	assert.Equal(t, "5678", retrievalPort)

	_, _, _, err = core.ParseOperatorSocket("localhost:12345678")
	assert.NotNil(t, err)
	assert.Equal(t, "invalid socket address format, missing retrieval port: localhost:12345678", err.Error())

	_, _, _, err = core.ParseOperatorSocket("localhost1234;5678")
	assert.NotNil(t, err)
	assert.Equal(t, "invalid socket address format: localhost1234;5678", err.Error())
}

func TestSignatureBytes(t *testing.T) {
	sig := &core.Signature{
		G1Point: core.NewG1Point(big.NewInt(1), big.NewInt(2)),
	}
	bytes := sig.Bytes()
	recovered := new(bn254.G1Affine)
	_, err := recovered.SetBytes(bytes[:])
	assert.NoError(t, err)
	assert.Equal(t, recovered, sig.G1Point.G1Affine)
}
