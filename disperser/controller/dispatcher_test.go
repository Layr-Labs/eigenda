package controller_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	clientsmock "github.com/Layr-Labs/eigenda/api/clients/mock"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/encoding"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/gammazero/workerpool"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wealdtech/go-merkletree/v2"
	"github.com/wealdtech/go-merkletree/v2/keccak256"
)

var (
	opId0, _          = core.OperatorIDFromHex("e22dae12a0074f20b8fc96a0489376db34075e545ef60c4845d264a732568311")
	opId1, _          = core.OperatorIDFromHex("e23cae12a0074f20b8fc96a0489376db34075e545ef60c4845d264b732568312")
	mockChainState, _ = coremock.NewChainDataMock(map[uint8]map[core.OperatorID]int{
		0: {
			opId0: 1,
			opId1: 1,
		},
		1: {
			opId0: 1,
			opId1: 3,
		},
	})
	finalizationBlockDelay = uint64(10)
)

type dispatcherComponents struct {
	Dispatcher        *controller.Dispatcher
	BlobMetadataStore *blobstore.BlobMetadataStore
	Pool              common.WorkerPool
	ChainReader       *coremock.MockWriter
	ChainState        *coremock.ChainDataMock
	SigAggregator     *core.StdSignatureAggregator
	NodeClientManager *controller.MockClientManager
}

func TestDispatcherHandleBatch(t *testing.T) {
	components := newDispatcherComponents(t)
	objs := setupBlobCerts(t, components.BlobMetadataStore, 2)
	ctx := context.Background()

	// Get batch header hash to mock signatures
	merkleTree, err := corev2.BuildMerkleTree(objs.blobCerts)
	require.NoError(t, err)
	require.NotNil(t, merkleTree)
	require.NotNil(t, merkleTree.Root())
	batchHeader := &corev2.BatchHeader{
		ReferenceBlockNumber: blockNumber - finalizationBlockDelay,
	}
	copy(batchHeader.BatchRoot[:], merkleTree.Root())
	bhh, err := batchHeader.Hash()
	require.NoError(t, err)

	mockClient0 := clientsmock.NewNodeClientV2()
	sig0 := mockChainState.KeyPairs[opId0].SignMessage(bhh)
	mockClient0.On("StoreChunks", mock.Anything, mock.Anything).Return(sig0, nil)
	op0Port := mockChainState.GetTotalOperatorState(ctx, uint(blockNumber)).PrivateOperators[opId0].DispersalPort
	op1Port := mockChainState.GetTotalOperatorState(ctx, uint(blockNumber)).PrivateOperators[opId1].DispersalPort
	require.NotEqual(t, op0Port, op1Port)
	components.NodeClientManager.On("GetClient", mock.Anything, op0Port).Return(mockClient0, nil)
	mockClient1 := clientsmock.NewNodeClientV2()
	sig1 := mockChainState.KeyPairs[opId1].SignMessage(bhh)
	mockClient1.On("StoreChunks", mock.Anything, mock.Anything).Return(sig1, nil)
	components.NodeClientManager.On("GetClient", mock.Anything, op1Port).Return(mockClient1, nil)

	sigChan, batchData, err := components.Dispatcher.HandleBatch(ctx)
	require.NoError(t, err)
	err = components.Dispatcher.HandleSignatures(ctx, batchData, sigChan)
	require.NoError(t, err)

	// Test that the blob metadata status are updated
	bm0, err := components.BlobMetadataStore.GetBlobMetadata(ctx, objs.blobKeys[0])
	require.NoError(t, err)
	require.Equal(t, v2.Certified, bm0.BlobStatus)
	bm1, err := components.BlobMetadataStore.GetBlobMetadata(ctx, objs.blobKeys[1])
	require.NoError(t, err)
	require.Equal(t, v2.Certified, bm1.BlobStatus)

	// Get batch header
	vis, err := components.BlobMetadataStore.GetBlobVerificationInfos(ctx, objs.blobKeys[0])
	require.NoError(t, err)
	require.Len(t, vis, 1)
	bhh, err = vis[0].BatchHeader.Hash()
	require.NoError(t, err)

	// Test that attestation is written
	att, err := components.BlobMetadataStore.GetAttestation(ctx, bhh)
	require.NoError(t, err)
	require.NotNil(t, att)
	require.Equal(t, vis[0].BatchHeader, att.BatchHeader)
	require.Greater(t, att.AttestedAt, uint64(0))
	require.Len(t, att.NonSignerPubKeys, 0)
	require.NotNil(t, att.APKG2)
	require.Len(t, att.QuorumAPKs, 2)
	require.NotNil(t, att.Sigma)
	require.ElementsMatch(t, att.QuorumNumbers, []core.QuorumID{0, 1})
}

func TestDispatcherNewBatch(t *testing.T) {
	components := newDispatcherComponents(t)
	objs := setupBlobCerts(t, components.BlobMetadataStore, 2)
	require.Len(t, objs.blobHedaers, 2)
	require.Len(t, objs.blobKeys, 2)
	require.Len(t, objs.blobMetadatas, 2)
	require.Len(t, objs.blobCerts, 2)
	ctx := context.Background()

	batchData, err := components.Dispatcher.NewBatch(ctx, blockNumber)
	require.NoError(t, err)
	batch := batchData.Batch
	bhh, keys, state := batchData.BatchHeaderHash, batchData.BlobKeys, batchData.OperatorState
	require.NotNil(t, batch)
	require.NotNil(t, batch.BatchHeader)
	require.NotNil(t, bhh)
	require.NotNil(t, keys)
	require.NotNil(t, state)
	require.ElementsMatch(t, keys, objs.blobKeys)

	// Test that the batch header hash is correct
	hash, err := batch.BatchHeader.Hash()
	require.NoError(t, err)
	require.Equal(t, bhh, hash)

	// Test that the batch header is correct
	require.Equal(t, blockNumber, batch.BatchHeader.ReferenceBlockNumber)
	require.NotNil(t, batch.BatchHeader.BatchRoot)

	// Test that the batch header is written
	bh, err := components.BlobMetadataStore.GetBatchHeader(ctx, bhh)
	require.NoError(t, err)
	require.NotNil(t, bh)
	require.Equal(t, bh, batch.BatchHeader)

	// Test that blob verification infos are written
	vi0, err := components.BlobMetadataStore.GetBlobVerificationInfo(ctx, objs.blobKeys[0], bhh)
	require.NoError(t, err)
	require.NotNil(t, vi0)
	cert := batch.BlobCertificates[vi0.BlobIndex]
	require.Equal(t, objs.blobHedaers[0], cert.BlobHeader)
	require.Equal(t, objs.blobKeys[0], vi0.BlobKey)
	require.Equal(t, bh, vi0.BatchHeader)
	certHash, err := cert.Hash()
	require.NoError(t, err)
	proof, err := core.DeserializeMerkleProof(vi0.InclusionProof, uint64(vi0.BlobIndex))
	require.NoError(t, err)
	verified, err := merkletree.VerifyProofUsing(certHash[:], false, proof, [][]byte{vi0.BatchRoot[:]}, keccak256.New())
	require.NoError(t, err)
	require.True(t, verified)

	// Attempt to create a batch with the same blobs
	_, err = components.Dispatcher.NewBatch(ctx, blockNumber)
	require.ErrorContains(t, err, "no blobs to dispatch")
}

func TestDispatcherBuildMerkleTree(t *testing.T) {
	certs := []*corev2.BlobCertificate{
		{
			BlobHeader: &corev2.BlobHeader{
				BlobVersion:     0,
				QuorumNumbers:   []core.QuorumID{0},
				BlobCommitments: mockCommitment,
				PaymentMetadata: core.PaymentMetadata{
					AccountID:         "account 1",
					BinIndex:          0,
					CumulativePayment: big.NewInt(532),
				},
				Signature: []byte("signature"),
			},
			RelayKeys: []corev2.RelayKey{0},
		},
		{
			BlobHeader: &corev2.BlobHeader{
				BlobVersion:     0,
				QuorumNumbers:   []core.QuorumID{0, 1},
				BlobCommitments: mockCommitment,
				PaymentMetadata: core.PaymentMetadata{
					AccountID:         "account 2",
					BinIndex:          0,
					CumulativePayment: big.NewInt(532),
				},
				Signature: []byte("signature"),
			},
			RelayKeys: []corev2.RelayKey{0, 1, 2},
		},
	}
	merkleTree, err := corev2.BuildMerkleTree(certs)
	require.NoError(t, err)
	require.NotNil(t, merkleTree)
	require.NotNil(t, merkleTree.Root())

	proof, err := merkleTree.GenerateProofWithIndex(uint64(0), 0)
	require.NoError(t, err)
	require.NotNil(t, proof)
	hash, err := certs[0].Hash()
	require.NoError(t, err)
	verified, err := merkletree.VerifyProofUsing(hash[:], false, proof, [][]byte{merkleTree.Root()}, keccak256.New())
	require.NoError(t, err)
	require.True(t, verified)

	proof, err = merkleTree.GenerateProofWithIndex(uint64(1), 0)
	require.NoError(t, err)
	require.NotNil(t, proof)
	hash, err = certs[1].Hash()
	require.NoError(t, err)
	verified, err = merkletree.VerifyProofUsing(hash[:], false, proof, [][]byte{merkleTree.Root()}, keccak256.New())
	require.NoError(t, err)
	require.True(t, verified)
}

type testObjects struct {
	blobHedaers   []*corev2.BlobHeader
	blobKeys      []corev2.BlobKey
	blobMetadatas []*v2.BlobMetadata
	blobCerts     []*corev2.BlobCertificate
}

func setupBlobCerts(t *testing.T, blobMetadataStore *blobstore.BlobMetadataStore, numObjects int) *testObjects {
	ctx := context.Background()
	headers := make([]*corev2.BlobHeader, numObjects)
	keys := make([]corev2.BlobKey, numObjects)
	metadatas := make([]*v2.BlobMetadata, numObjects)
	certs := make([]*corev2.BlobCertificate, numObjects)
	for i := 0; i < numObjects; i++ {
		randomBytes := make([]byte, 16)
		_, err := rand.Read(randomBytes)
		require.NoError(t, err)
		randomBinIndex, err := rand.Int(rand.Reader, big.NewInt(1000))
		require.NoError(t, err)
		binIndex := uint32(randomBinIndex.Uint64())
		headers[i] = &corev2.BlobHeader{
			BlobVersion:     0,
			QuorumNumbers:   []core.QuorumID{0, 1},
			BlobCommitments: mockCommitment,
			PaymentMetadata: core.PaymentMetadata{
				AccountID:         hex.EncodeToString(randomBytes),
				BinIndex:          binIndex,
				CumulativePayment: big.NewInt(532),
			},
		}
		key, err := headers[i].BlobKey()
		require.NoError(t, err)
		keys[i] = key
		now := time.Now()
		metadatas[i] = &v2.BlobMetadata{
			BlobHeader: headers[i],
			BlobStatus: v2.Encoded,
			Expiry:     uint64(now.Add(time.Hour).Unix()),
			NumRetries: 0,
			UpdatedAt:  uint64(now.UnixNano()) - uint64(i),
		}
		err = blobMetadataStore.PutBlobMetadata(ctx, metadatas[i])
		require.NoError(t, err)

		certs[i] = &corev2.BlobCertificate{
			BlobHeader: headers[i],
			RelayKeys:  []corev2.RelayKey{0, 1, 2},
		}
		err = blobMetadataStore.PutBlobCertificate(ctx, certs[i], &encoding.FragmentInfo{})
		require.NoError(t, err)
	}

	return &testObjects{
		blobHedaers:   headers,
		blobKeys:      keys,
		blobMetadatas: metadatas,
		blobCerts:     certs,
	}
}

func newDispatcherComponents(t *testing.T) *dispatcherComponents {
	// logger := logging.NewNoopLogger()
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)
	pool := workerpool.New(5)

	chainReader := &coremock.MockWriter{}
	chainReader.On("OperatorIDToAddress").Return(gethcommon.Address{0}, nil)
	agg, err := core.NewStdSignatureAggregator(logger, chainReader)
	require.NoError(t, err)
	nodeClientManager := &controller.MockClientManager{}
	mockChainState.On("GetCurrentBlockNumber").Return(uint(blockNumber), nil)
	d, err := controller.NewDispatcher(controller.DispatcherConfig{
		PullInterval:           1 * time.Second,
		FinalizationBlockDelay: finalizationBlockDelay,
		NodeRequestTimeout:     1 * time.Second,
		NumRequestRetries:      3,
	}, blobMetadataStore, pool, mockChainState, agg, nodeClientManager, logger)
	require.NoError(t, err)
	return &dispatcherComponents{
		Dispatcher:        d,
		BlobMetadataStore: blobMetadataStore,
		Pool:              pool,
		ChainReader:       chainReader,
		ChainState:        mockChainState,
		SigAggregator:     agg,
		NodeClientManager: nodeClientManager,
	}
}
