package retriever_test

import (
	"context"
	"log"
	"runtime"
	"testing"

	clientsmock "github.com/Layr-Labs/eigenda/api/clients/mock"
	pb "github.com/Layr-Labs/eigenda/api/grpc/retriever"
	binding "github.com/Layr-Labs/eigenda/contracts/bindings/EigenDAServiceManager"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/retriever"
	"github.com/Layr-Labs/eigenda/retriever/mock"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/assert"
)

const numOperators = 10

var (
	indexedChainState      core.IndexedChainState
	retrievalClient        *clientsmock.MockRetrievalClient
	chainClient            *mock.MockChainClient
	batchHeaderHash        [32]byte
	batchRoot              [32]byte
	gettysburgAddressBytes = codec.ConvertByPaddingEmptyByte([]byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth."))
)

func makeTestComponents() (encoding.Prover, encoding.Verifier, error) {
	config := &kzg.KzgConfig{
		G1Path:          "../inabox/resources/kzg/g1.point",
		G2Path:          "../inabox/resources/kzg/g2.point",
		CacheDir:        "../inabox/resources/kzg/SRSTables",
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
	if err != nil {
		log.Fatalf("failed to create new mocked chain data: %s", err)
	}

	_, _, err = makeTestComponents()
	if err != nil {
		log.Fatal(err)
	}

	retrievalClient = &clientsmock.MockRetrievalClient{}
	chainClient = mock.NewMockChainClient()
	return retriever.NewServer(config, logger, retrievalClient, indexedChainState, chainClient)
}

func TestRetrieveBlob(t *testing.T) {
	server := newTestServer(t)
	chainClient.On("FetchBatchHeader").Return(&binding.BatchHeader{
		BlobHeadersRoot:       batchRoot,
		QuorumNumbers:         []byte{0},
		SignedStakeForQuorums: []byte{90},
		ReferenceBlockNumber:  0,
	}, nil)

	retrievalClient.On("RetrieveBlob").Return(gettysburgAddressBytes, nil)

	retrievalReply, err := server.RetrieveBlob(context.Background(), &pb.BlobRequest{
		BatchHeaderHash:      batchHeaderHash[:],
		BlobIndex:            0,
		ReferenceBlockNumber: 0,
		QuorumId:             0,
	})
	assert.NoError(t, err)
	assert.Equal(t, gettysburgAddressBytes, retrievalReply.Data)
}
