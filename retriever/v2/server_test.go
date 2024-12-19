package v2_test

import (
	"context"
	"math/big"
	"runtime"
	"testing"

	clientsmock "github.com/Layr-Labs/eigenda/api/clients/v2/mock"
	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"
	commonpbv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	pb "github.com/Layr-Labs/eigenda/api/grpc/retriever/v2"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/retriever/mock"
	retriever "github.com/Layr-Labs/eigenda/retriever/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/stretchr/testify/require"
)

const numOperators = 10

var (
	indexedChainState      core.IndexedChainState
	retrievalClient        *clientsmock.MockRetrievalClient
	chainClient            *mock.MockChainClient
	gettysburgAddressBytes = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
)

func makeTestComponents() (encoding.Prover, encoding.Verifier, error) {
	config := &kzg.KzgConfig{
		G1Path:          "../../inabox/resources/kzg/g1.point",
		G2Path:          "../../inabox/resources/kzg/g2.point",
		CacheDir:        "../../inabox/resources/kzg/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}

	p, err := prover.NewProver(config, nil)
	if err != nil {
		return nil, nil, err
	}

	v, err := verifier.NewVerifier(config, nil)
	if err != nil {
		return nil, nil, err
	}

	return p, v, nil
}

func newTestServer(t *testing.T) *retriever.Server {
	var err error
	config := &retriever.Config{}

	logger := logging.NewNoopLogger()

	indexedChainState, err = coremock.MakeChainDataMock(map[uint8]int{
		0: numOperators,
		1: numOperators,
		2: numOperators,
	})
	require.NoError(t, err)

	_, _, err = makeTestComponents()
	require.NoError(t, err)

	retrievalClient = &clientsmock.MockRetrievalClient{}
	chainClient = mock.NewMockChainClient()
	return retriever.NewServer(config, logger, retrievalClient, indexedChainState)
}

func TestRetrieveBlob(t *testing.T) {
	server := newTestServer(t)
	data := codec.ConvertByPaddingEmptyByte(gettysburgAddressBytes)
	retrievalClient.On("GetBlob").Return(data, nil)

	var X1, Y1 fp.Element
	X1 = *X1.SetBigInt(big.NewInt(1))
	Y1 = *Y1.SetBigInt(big.NewInt(2))

	var lengthXA0, lengthXA1, lengthYA0, lengthYA1 fp.Element
	_, err := lengthXA0.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	require.NoError(t, err)
	_, err = lengthXA1.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	require.NoError(t, err)
	_, err = lengthYA0.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	require.NoError(t, err)
	_, err = lengthYA1.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531")
	require.NoError(t, err)

	var lengthProof, lengthCommitment bn254.G2Affine
	lengthProof.X.A0 = lengthXA0
	lengthProof.X.A1 = lengthXA1
	lengthProof.Y.A0 = lengthYA0
	lengthProof.Y.A1 = lengthYA1

	lengthCommitment = lengthProof

	mockCommitment := encoding.BlobCommitments{
		Commitment: &encoding.G1Commitment{
			X: X1,
			Y: Y1,
		},
		LengthCommitment: (*encoding.G2Commitment)(&lengthCommitment),
		LengthProof:      (*encoding.G2Commitment)(&lengthProof),
		Length:           16,
	}
	c, err := mockCommitment.ToProtobuf()
	require.NoError(t, err)
	retrievalReply, err := server.RetrieveBlob(context.Background(), &pb.BlobRequest{
		BlobHeader: &commonpbv2.BlobHeader{
			Version:       0,
			QuorumNumbers: []uint32{0},
			Commitment:    c,
			PaymentHeader: &commonpb.PaymentHeader{
				AccountId: "account_id",
			},
		},
		ReferenceBlockNumber: 100,
		QuorumId:             0,
	})
	require.NoError(t, err)
	require.Equal(t, gettysburgAddressBytes, retrievalReply.Data)
}
