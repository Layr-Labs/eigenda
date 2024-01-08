package integration_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/pubip"

	clientsmock "github.com/Layr-Labs/eigenda/clients/mock"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	dispatcher "github.com/Layr-Labs/eigenda/disperser/batcher/grpc"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/retriever"
	retrievermock "github.com/Layr-Labs/eigenda/retriever/mock"
	"github.com/Layr-Labs/eigensdk-go/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/peer"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/logging"
	commonmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/encoding"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	batchermock "github.com/Layr-Labs/eigenda/disperser/batcher/mock"
	"github.com/Layr-Labs/eigenda/disperser/common/inmem"
	dispersermock "github.com/Layr-Labs/eigenda/disperser/mock"
	"github.com/Layr-Labs/eigenda/node"
	nodegrpc "github.com/Layr-Labs/eigenda/node/grpc"
	"github.com/Layr-Labs/eigenda/pkg/encoding/kzgEncoder"
	"github.com/Layr-Labs/eigenda/pkg/kzg/bn254"

	nodepb "github.com/Layr-Labs/eigenda/api/grpc/node"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	enc core.Encoder
	asn core.AssignmentCoordinator

	gettysburgAddressBytes = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
	serviceManagerAddress  = gethcommon.HexToAddress("0x0000000000000000000000000000000000000000")
)

const (
	numOperators         = 10
	disperserGrpcPort    = 4000
	encoderPort          = "3100"
	q0AdversaryThreshold = uint8(80)
	q0QuorumThreshold    = uint8(100)
)

func init() {
	enc = mustMakeTestEncoder()
	asn = &core.StdAssignmentCoordinator{}
}

// makeTestEncoder makes an encoder currently using the only supported backend.
func mustMakeTestEncoder() core.Encoder {
	config := kzgEncoder.KzgConfig{
		G1Path:    "../inabox/resources/kzg/g1.point",
		G2Path:    "../inabox/resources/kzg/g2.point",
		CacheDir:  "../inabox/resources/kzg/SRSTables",
		SRSOrder:  3000,
		NumWorker: uint64(runtime.GOMAXPROCS(0)),
	}

	encoder, err := encoding.NewEncoder(
		encoding.EncoderConfig{KzgConfig: config},
	)
	if err != nil {
		log.Fatalln("failed to initialize new encoder")
	}

	return encoder
}

func mustMakeTestBlob() core.Blob {
	blob := core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: []*core.SecurityParam{
				{
					QuorumID:           0,
					AdversaryThreshold: q0AdversaryThreshold,
					QuorumThreshold:    q0QuorumThreshold,
				},
			},
		},
		Data: gettysburgAddressBytes,
	}
	return blob
}

type TestDisperser struct {
	Batcher       *batcher.Batcher
	Server        *apiserver.DispersalServer
	EncoderServer *encoder.Server
}

func mustMakeDisperser(t *testing.T, cst core.IndexedChainState, store disperser.BlobStore, logger common.Logger) TestDisperser {
	dispatcherConfig := &dispatcher.Config{
		Timeout: time.Second,
	}
	dispatcher := dispatcher.NewDispatcher(dispatcherConfig, logger)

	transactor := &coremock.MockTransactor{}
	transactor.On("OperatorIDToAddress").Return(gethcommon.Address{}, nil)
	agg, err := core.NewStdSignatureAggregator(logger, transactor)
	assert.NoError(t, err)

	confirmer := dispersermock.NewBatchConfirmer()
	// should be encoding 3 and 0
	logData, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatal(err)
	}

	receipt := &types.Receipt{
		Logs: []*types.Log{
			{
				Topics: []gethcommon.Hash{common.BatchConfirmedEventSigHash, gethcommon.HexToHash("1234")},
				Data:   logData,
			},
		},
		BlockNumber: big.NewInt(123),
	}
	confirmer.On("ConfirmBatch", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(receipt, nil)

	batcherConfig := batcher.Config{
		PullInterval:             5 * time.Second,
		NumConnections:           1,
		EncoderSocket:            fmt.Sprintf("localhost:%s", encoderPort),
		EncodingRequestQueueSize: 100,
		SRSOrder:                 3000,
	}
	timeoutConfig := batcher.TimeoutConfig{
		EncodingTimeout:    10 * time.Second,
		AttestationTimeout: 10 * time.Second,
		ChainReadTimeout:   10 * time.Second,
		ChainWriteTimeout:  10 * time.Second,
	}

	enc0 := mustMakeTestEncoder()
	metrics := encoder.NewMetrics("9000", logger)
	grpcEncoder := encoder.NewServer(encoder.ServerConfig{
		GrpcPort:              encoderPort,
		MaxConcurrentRequests: 16,
		RequestPoolSize:       32,
	}, logger, enc0, metrics)

	encoderClient, err := encoder.NewEncoderClient(batcherConfig.EncoderSocket, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	finalizer := batchermock.NewFinalizer()

	disperserMetrics := disperser.NewMetrics("9100", logger)
	batcherMetrics := batcher.NewMetrics("9100", logger)

	batcher, err := batcher.NewBatcher(batcherConfig, timeoutConfig, store, dispatcher, confirmer, cst, asn, encoderClient, agg, &commonmock.MockEthClient{}, finalizer, logger, batcherMetrics)
	if err != nil {
		t.Fatal(err)
	}

	ratelimiter := &commonmock.NoopRatelimiter{}
	rateConfig := apiserver.RateConfig{
		QuorumRateInfos: map[core.QuorumID]apiserver.QuorumRateInfo{
			0: {
				PerUserUnauthThroughput: 0,
				TotalUnauthThroughput:   0,
			},
		},
	}

	serverConfig := disperser.ServerConfig{
		GrpcPort: fmt.Sprint(disperserGrpcPort),
	}
	tx := &coremock.MockTransactor{}
	tx.On("GetCurrentBlockNumber").Return(uint64(100), nil)
	tx.On("GetQuorumCount").Return(1, nil)
	server := apiserver.NewDispersalServer(serverConfig, store, tx, logger, disperserMetrics, ratelimiter, rateConfig)

	return TestDisperser{
		Batcher:       batcher,
		Server:        server,
		EncoderServer: grpcEncoder,
	}
}

type TestOperator struct {
	Node   *node.Node
	Server *nodegrpc.Server
}

func mustMakeOperators(t *testing.T, cst *coremock.ChainDataMock, logger common.Logger) map[core.OperatorID]TestOperator {
	bn := uint(0)
	state := cst.GetTotalOperatorState(context.Background(), bn)

	ops := make(map[core.OperatorID]TestOperator, len(state.IndexedOperators))

	setRegisteredQuorums := true
	for id, op := range state.PrivateOperators {

		idStr := hexutil.Encode(id[:])
		fmt.Println("Operator: ", idStr)

		dbPath := fmt.Sprintf("testdata/%v/db", idStr)
		logPath := fmt.Sprintf("testdata/%v/log", idStr)

		err := os.RemoveAll(dbPath)
		if err != nil {
			t.Fatal(err)
		}
		err = os.RemoveAll(logPath)
		if err != nil {
			t.Fatal(err)
		}
		err = os.MkdirAll(dbPath, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}
		err = os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			t.Fatal(err)
		}

		registeredQuorums := []core.QuorumID{}
		if setRegisteredQuorums {
			registeredQuorums = []core.QuorumID{0}
			setRegisteredQuorums = false
		}

		config := &node.Config{
			Hostname:                  op.Host,
			DispersalPort:             op.DispersalPort,
			RetrievalPort:             op.RetrievalPort,
			InternalRetrievalPort:     op.RetrievalPort,
			InternalDispersalPort:     op.DispersalPort,
			EnableMetrics:             false,
			Timeout:                   10,
			ExpirationPollIntervalSec: 10,
			DbPath:                    dbPath,
			LogPath:                   logPath,
			PrivateBls:                string(op.KeyPair.GetPubKeyG1().Serialize()),
			ID:                        id,
			QuorumIDList:              registeredQuorums,
		}

		// creating a new instance of encoder instead of sharing enc because enc is not thread safe
		encoder := mustMakeTestEncoder()
		val := core.NewChunkValidator(encoder, asn, cst, id)

		noopMetrics := metrics.NewNoopMetrics()
		reg := prometheus.NewRegistry()
		metrics := node.NewMetrics(noopMetrics, reg, logger, ":9090")
		store, err := node.NewLevelDBStore(config.DbPath+"/chunk", logger, metrics, 1e9, 1e9)
		if err != nil {
			t.Fatal(err)
		}

		tx := &coremock.MockTransactor{}
		tx.On("RegisterBLSPublicKey").Return(nil)
		tx.On("RegisterOperator").Return(nil)
		tx.On("GetRegisteredQuorumIdsForOperator").Return(registeredQuorums, nil)
		tx.On("UpdateOperatorSocket").Return(nil)
		tx.On("GetBlockStaleMeasure").Return(nil)
		tx.On("GetStoreDurationBlocks").Return(nil)

		mockOperatorSocketsFilterer := &coremock.MockOperatorSocketsFilterer{}

		mockSocketChan := make(chan string)
		mockOperatorSocketsFilterer.On("WatchOperatorSocketUpdate").Return(mockSocketChan, nil)

		pubIPProvider := &pubip.SimpleProvider{
			RequestDoer: pubip.RequestDoerFunc(func(req *http.Request) (*http.Response, error) {
				w := httptest.NewRecorder()
				_, _ = w.WriteString("8.8.8.8")
				return w.Result(), nil
			}),
			Name: "",
			URL:  "",
		}

		n := &node.Node{
			Config:                  config,
			Logger:                  logger,
			KeyPair:                 op.KeyPair,
			Metrics:                 metrics,
			Store:                   store,
			ChainState:              cst,
			Validator:               val,
			Transactor:              tx,
			PubIPProvider:           pubIPProvider,
			OperatorSocketsFilterer: mockOperatorSocketsFilterer,
		}

		ratelimiter := &commonmock.NoopRatelimiter{}

		s := nodegrpc.NewServer(config, n, logger, ratelimiter)

		ops[id] = TestOperator{
			Node:   n,
			Server: s,
		}
	}

	return ops
}

type TestRetriever struct {
	Server *retriever.Server
}

func mustMakeRetriever(cst core.IndexedChainState, logger common.Logger) (*commonmock.MockEthClient, TestRetriever) {
	config := &retriever.Config{
		Timeout: 5 * time.Second,
	}
	gethClient := &commonmock.MockEthClient{}
	retrievalClient := &clientsmock.MockRetrievalClient{}
	chainClient := retrievermock.NewMockChainClient()
	server := retriever.NewServer(config, logger, retrievalClient, enc, cst, chainClient)

	return gethClient, TestRetriever{
		Server: server,
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestDispersalAndRetrieval(t *testing.T) {

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 3000,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	cst, err := coremock.MakeChainDataMock(numOperators)
	assert.NoError(t, err)

	cst.On("GetCurrentBlockNumber").Return(uint(10), nil)

	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	assert.NoError(t, err)
	store := inmem.NewBlobStore()
	dis := mustMakeDisperser(t, cst, store, logger)
	go func() {
		_ = dis.EncoderServer.Start()
	}()
	t.Cleanup(func() {
		dis.EncoderServer.Close()
	})
	ops := mustMakeOperators(t, cst, logger)
	gethClient, _ := mustMakeRetriever(cst, logger)

	for _, op := range ops {
		idStr := hexutil.Encode(op.Node.Config.ID[:])
		fmt.Println("Operator: ", idStr)

		fmt.Println("Starting node")
		err = op.Node.Start(ctx)
		assert.NoError(t, err)

		fmt.Println("Starting server")
		go op.Server.Start()
	}

	blob := mustMakeTestBlob()
	requestedAt := uint64(time.Now().UnixNano())
	metadataKey, err := store.StoreBlob(ctx, &blob, requestedAt)
	assert.NoError(t, err)
	out := make(chan batcher.EncodingResultOrStatus)
	err = dis.Batcher.EncodingStreamer.RequestEncoding(context.Background(), out)
	assert.NoError(t, err)
	err = dis.Batcher.EncodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.NoError(t, err)
	dis.Batcher.EncodingStreamer.Pool.StopWait()

	err = dis.Batcher.HandleSingleBatch(ctx)
	assert.NoError(t, err)

	// Check that the blob was processed
	metadata, err := store.GetBlobMetadata(ctx, metadataKey)
	assert.NoError(t, err)
	assert.Equal(t, metadataKey, metadata.GetBlobKey())
	assert.Equal(t, disperser.Confirmed, metadata.BlobStatus)

	isConfirmed, err := metadata.IsConfirmed()
	assert.NoError(t, err)
	assert.True(t, isConfirmed)
	batchHeaderHash := metadata.ConfirmationInfo.BatchHeaderHash[:]
	txHash := gethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	topics := [][]gethcommon.Hash{
		{common.BatchConfirmedEventSigHash},
		{gethcommon.BytesToHash(batchHeaderHash)},
	}
	calldata, err := hex.DecodeString("7794965a000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000560000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000001c01b4136a161225e9cebe4e2c561148043b2fde423fc5b64e01d897d0fb7970a142d5474fb609bda1b747bdb5c47375d5819000e3c5cbc75baf55b19849410a2610de9c40eb95b49aca940e0bec6ae8b2868855a6324d04d864cbfa61128cf06a51c069e5a0c490c5a359086b0a3660c2ea2e4fb50722bec1ef593c5245413e4cd0a3c7e490348fb279ccb58f91a3bd494511c2ab0321e3922a0cd26012ef3133c043acb758e735db805d360196f3fc89a6395a4b174c19b981afb7f64c2b1193e0000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000026000000000000000000000000000000000000000000000000000000000000002a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001170c867415fef7db6d88e37598228f43de085616a25939dacbb6b5900f680c7f1d582c9ea38023afb08f368ea93692d17946619d9cf5f3c4d7b3c0cff1a92dff0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000")
	assert.NoError(t, err)
	r, ok := new(big.Int).SetString("8ad2b300a012fb0e90dceb8b66fa564717a2d218ca0fd25f11a1875e0153d1d8", 16)
	assert.True(t, ok)
	s, ok := new(big.Int).SetString("1accb1e1c69fa07bd4237d92143275960b24eec780862a673d54ffaaa5e77f9b", 16)
	assert.True(t, ok)
	gethClient.On("TransactionByHash", txHash).Return(
		types.NewTx(&types.DynamicFeeTx{
			ChainID:    big.NewInt(1),
			Nonce:      1,
			GasTipCap:  big.NewInt(1_000_000),
			GasFeeCap:  big.NewInt(1_000_000),
			Gas:        298617,
			To:         &serviceManagerAddress,
			Value:      big.NewInt(0),
			Data:       calldata,
			AccessList: types.AccessList{},
			V:          big.NewInt(0x1),
			R:          r,
			S:          s,
		}), false, nil)
	gethClient.On("FilterLogs", ethereum.FilterQuery{
		Addresses: []gethcommon.Address{serviceManagerAddress},
		Topics:    topics,
	}).Return([]types.Log{
		{
			Address: serviceManagerAddress,
			Topics: []gethcommon.Hash{
				topics[0][0], topics[1][0],
			},
			Data:        []byte{},
			BlockHash:   gethcommon.HexToHash("0x0"),
			BlockNumber: 123,
			TxHash:      txHash,
			TxIndex:     0,
			Index:       0,
		},
	}, nil)

	operatorState, err := cst.GetOperatorState(ctx, 0, []core.QuorumID{0})
	assert.NoError(t, err)

	blobLength := core.GetBlobLength(uint(len(blob.Data)))
	chunkLength, err := asn.CalculateChunkLength(operatorState, blobLength, 0, blob.RequestHeader.SecurityParams[0])
	assert.NoError(t, err)

	blobQuorumInfo := &core.BlobQuorumInfo{
		SecurityParam: core.SecurityParam{
			QuorumID:           0,
			AdversaryThreshold: q0AdversaryThreshold,
			QuorumThreshold:    q0QuorumThreshold,
		},
		ChunkLength: chunkLength,
	}

	assignments, info, err := asn.GetAssignments(operatorState, blobLength, blobQuorumInfo)
	assert.NoError(t, err)

	var indices []core.ChunkNumber
	var chunks []*core.Chunk
	var blobHeader *core.BlobHeader
	for _, op := range ops {

		// check that blob headers can be retrieved from operators
		headerReply, err := op.Server.GetBlobHeader(ctx, &nodepb.GetBlobHeaderRequest{
			BatchHeaderHash: batchHeaderHash,
			BlobIndex:       metadata.ConfirmationInfo.BlobIndex,
			QuorumId:        uint32(0),
		})
		assert.NoError(t, err)
		actualCommitment, err := new(core.Commitment).Deserialize(headerReply.GetBlobHeader().GetCommitment())
		assert.NoError(t, err)
		actualLengthProof, err := new(core.Commitment).Deserialize(headerReply.GetBlobHeader().GetLengthProof())
		assert.NoError(t, err)
		assert.Equal(t, metadata.ConfirmationInfo.BlobCommitment.Commitment, actualCommitment)
		assert.Equal(t, metadata.ConfirmationInfo.BlobCommitment.LengthProof, actualLengthProof)
		assert.Equal(t, uint32(metadata.ConfirmationInfo.BlobCommitment.Length), headerReply.GetBlobHeader().GetLength())
		assert.Len(t, headerReply.GetBlobHeader().GetQuorumHeaders(), 1)
		assert.Equal(t, uint32(0), headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetQuorumId())
		assert.Equal(t, uint32(q0QuorumThreshold), headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetQuorumThreshold())
		assert.Equal(t, uint32(q0AdversaryThreshold), headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetAdversaryThreshold())
		assert.Greater(t, headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetChunkLength(), uint32(0))

		if blobHeader == nil {
			blobHeader, err = nodegrpc.GetBlobHeaderFromProto(headerReply.GetBlobHeader())
			assert.NoError(t, err)
		}

		// check that chunks can be retrieved from operators
		chunksReply, err := op.Server.RetrieveChunks(ctx, &nodepb.RetrieveChunksRequest{
			BatchHeaderHash: batchHeaderHash,
			BlobIndex:       metadata.ConfirmationInfo.BlobIndex,
			QuorumId:        uint32(0),
		})

		assert.NoError(t, err)
		assignment, ok := assignments[op.Node.Config.ID]
		assert.True(t, ok)
		for _, data := range chunksReply.GetChunks() {
			chunk, err := new(core.Chunk).Deserialize(data)
			assert.NoError(t, err)
			chunks = append(chunks, chunk)
		}
		assert.Len(t, chunksReply.GetChunks(), int(assignments[op.Node.Config.ID].NumChunks))
		indices = append(indices, assignment.GetIndices()...)
	}

	encodingParams, err := core.GetEncodingParams(chunkLength, info.TotalChunks)
	assert.NoError(t, err)
	recovered, err := enc.Decode(chunks, indices, encodingParams, uint64(blobHeader.Length)*bn254.BYTES_PER_COEFFICIENT)
	assert.NoError(t, err)
	recovered = bytes.TrimRight(recovered, "\x00")
	assert.Equal(t, gettysburgAddressBytes, recovered)
}
