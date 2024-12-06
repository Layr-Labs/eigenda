package clients_test

import (
	"bytes"
	"context"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients"
	clientsmock "github.com/Layr-Labs/eigenda/api/clients/mock"
	"github.com/Layr-Labs/eigenda/core"
	coreindexer "github.com/Layr-Labs/eigenda/core/indexer"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	indexermock "github.com/Layr-Labs/eigenda/indexer/mock"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wealdtech/go-merkletree/v2"
	"github.com/wealdtech/go-merkletree/v2/keccak256"
)

const numOperators = 10

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

var (
	indexedChainState core.IndexedChainState
	chainState        core.ChainState
	indexer           *indexermock.MockIndexer
	operatorState     *core.OperatorState
	nodeClient        *clientsmock.MockNodeClient
	coordinator       *core.StdAssignmentCoordinator
	retrievalClient   clients.RetrievalClient
	blobHeader        *core.BlobHeader
	encodedBlob       core.EncodedBlob = core.EncodedBlob{
		BlobHeader:               nil,
		EncodedBundlesByOperator: make(map[core.OperatorID]core.EncodedBundles),
	}
	batchHeaderHash        [32]byte
	batchRoot              [32]byte
	gettysburgAddressBytes = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
)

func setup(t *testing.T) {

	var err error
	chainState, err = coremock.MakeChainDataMock(map[uint8]int{
		0: numOperators,
		1: numOperators,
		2: numOperators,
	})
	if err != nil {
		t.Fatalf("failed to create new mocked chain data: %s", err)
	}

	indexedChainState, err = coremock.MakeChainDataMock(map[uint8]int{
		0: numOperators,
		1: numOperators,
		2: numOperators,
	})
	if err != nil {
		t.Fatalf("failed to create new mocked indexed chain data: %s", err)
	}

	nodeClient = clientsmock.NewNodeClient()
	coordinator = &core.StdAssignmentCoordinator{}
	p, v, err := makeTestComponents()
	if err != nil {
		t.Fatal(err)
	}
	logger := logging.NewNoopLogger()
	indexer = &indexermock.MockIndexer{}
	indexer.On("Index").Return(nil).Once()

	ics, err := coreindexer.NewIndexedChainState(chainState, indexer)
	if err != nil {
		panic("failed to create a new indexed chain state")
	}

	retrievalClient, err = clients.NewRetrievalClient(logger, ics, coordinator, nodeClient, v, 2)
	if err != nil {
		panic("failed to create a new retrieval client")
	}
	err = indexer.Index(context.Background())
	if err != nil {
		panic("failed to start indexing")
	}

	var (
		quorumID           core.QuorumID = 0
		adversaryThreshold uint8         = 80
		quorumThreshold    uint8         = 90
	)
	securityParams := []*core.SecurityParam{
		{
			QuorumID:              quorumID,
			ConfirmationThreshold: quorumThreshold,
			AdversaryThreshold:    adversaryThreshold,
		},
	}
	blob := core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: securityParams,
		},
		Data: codec.ConvertByPaddingEmptyByte(gettysburgAddressBytes),
	}
	operatorState, err = indexedChainState.GetOperatorState(context.Background(), (0), []core.QuorumID{quorumID})
	if err != nil {
		t.Fatalf("failed to get operator state: %s", err)
	}

	blobSize := uint(len(blob.Data))
	blobLength := encoding.GetBlobLength(uint(blobSize))

	chunkLength, err := coordinator.CalculateChunkLength(operatorState, blobLength, 0, securityParams[0])
	if err != nil {
		t.Fatal(err)
	}

	quorumHeader := &core.BlobQuorumInfo{
		SecurityParam: core.SecurityParam{
			QuorumID:              quorumID,
			AdversaryThreshold:    adversaryThreshold,
			ConfirmationThreshold: quorumThreshold,
		},
		ChunkLength: chunkLength,
	}

	assignments, info, err := coordinator.GetAssignments(operatorState, blobLength, quorumHeader)
	if err != nil {
		t.Fatal(err)
	}

	params := encoding.ParamsFromMins(chunkLength, info.TotalChunks)

	commitments, chunks, err := p.EncodeAndProve(blob.Data, params)
	if err != nil {
		t.Fatal(err)
	}

	blobHeader = &core.BlobHeader{
		BlobCommitments: encoding.BlobCommitments{
			Commitment:       commitments.Commitment,
			LengthCommitment: commitments.LengthCommitment,
			LengthProof:      commitments.LengthProof,
			Length:           commitments.Length,
		},
		QuorumInfos: []*core.BlobQuorumInfo{quorumHeader},
	}

	blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
	if err != nil {
		t.Fatal(err)
	}

	tree, err := merkletree.NewTree(merkletree.WithData([][]byte{blobHeaderHash[:]}), merkletree.WithHashType(keccak256.New()))
	if err != nil {
		t.Fatal(err)
	}

	copy(batchRoot[:], tree.Root())
	batchHeaderHash, err = core.BatchHeader{
		BatchRoot:            batchRoot,
		ReferenceBlockNumber: 0,
	}.GetBatchHeaderHash()
	if err != nil {
		t.Fatal(err)
	}

	for id, assignment := range assignments {
		bundles := make(map[core.QuorumID]core.Bundle, len(blobHeader.QuorumInfos))
		bundles[quorumID] = chunks[assignment.StartIndex : assignment.StartIndex+assignment.NumChunks]
		encodedBlob.BlobHeader = blobHeader
		eb, err := core.Bundles(bundles).ToEncodedBundles()
		if err != nil {
			t.Fatal(err)
		}
		encodedBlob.EncodedBundlesByOperator[id] = eb
	}

}

func mustMakeOpertatorPubKeysPair(t *testing.T) *coreindexer.OperatorPubKeys {
	operators := make(map[core.OperatorID]coreindexer.OperatorPubKeysPair, len(operatorState.Operators))
	for operatorId := range operatorState.Operators[0] {
		keyPair, err := core.GenRandomBlsKeys()
		if err != nil {
			t.Fatalf("Generating random BLS keys Error: %s", err.Error())
		}
		operators[operatorId] = coreindexer.OperatorPubKeysPair{
			PubKeyG1: keyPair.PubKey.G1Affine,
			PubKeyG2: keyPair.GetPubKeyG2().G2Affine,
		}
	}
	keyPair, err := core.GenRandomBlsKeys()
	if err != nil {
		t.Fatalf("Generating random BLS keys Error: %s", err.Error())
	}
	return &coreindexer.OperatorPubKeys{
		Operators: operators,
		QuorumTotals: map[core.QuorumID]*bn254.G1Affine{
			0: keyPair.PubKey.G1Affine,
		},
	}
}

func musMakeOperatorSocket(t *testing.T) coreindexer.OperatorSockets {
	operatorSocket := make(coreindexer.OperatorSockets, len(operatorState.Operators))
	for operatorId := range operatorState.Operators[0] {
		operatorSocket[operatorId] = "test"
	}
	return operatorSocket
}

func TestInvalidBlobHeader(t *testing.T) {

	setup(t)

	// TODO: add the blob proof to the response
	nodeClient.On("GetBlobHeader", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(blobHeader, [][]byte{{1}}, uint64(0), nil).Times(numOperators)
	nodeClient.
		On("GetChunks", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(encodedBlob)

	operatorPubKeys := mustMakeOpertatorPubKeysPair(t)
	operatorSocket := musMakeOperatorSocket(t)

	indexer.On("GetObject", mock.Anything, 0).Return(operatorPubKeys, nil).Once()
	indexer.On("GetObject", mock.Anything, 1).Return(operatorSocket, nil).Once()

	_, err := retrievalClient.RetrieveBlob(context.Background(), batchHeaderHash, 0, 0, batchRoot, 0)
	assert.ErrorContains(t, err, "failed to get blob header from all operators")

}

func TestValidBlobHeader(t *testing.T) {

	setup(t)

	// TODO: add the blob proof to the response
	nodeClient.On("GetBlobHeader", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(blobHeader, [][]byte{}, uint64(0), nil).Once()
	nodeClient.
		On("GetChunks", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(encodedBlob)

	operatorPubKeys := mustMakeOpertatorPubKeysPair(t)
	operatorSocket := musMakeOperatorSocket(t)

	indexer.On("GetObject", mock.Anything, 0).Return(operatorPubKeys, nil).Once()
	indexer.On("GetObject", mock.Anything, 1).Return(operatorSocket, nil).Once()

	data, err := retrievalClient.RetrieveBlob(context.Background(), batchHeaderHash, 0, 0, batchRoot, 0)
	assert.NoError(t, err)

	restored := codec.RemoveEmptyByteFromPaddedBytes(data)
	assert.Len(t, restored, 1488) // 48*31
	restored = bytes.TrimRight(restored, "\x00")
	assert.Equal(t, gettysburgAddressBytes, restored[:len(gettysburgAddressBytes)])

}
