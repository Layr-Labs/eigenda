package integration_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"time"

	clientsmock "github.com/Layr-Labs/eigenda/api/clients/mock"
	nodepb "github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common"
	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	commonmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/common/pubip"
	"github.com/Layr-Labs/eigenda/common/version"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	dispatcher "github.com/Layr-Labs/eigenda/disperser/batcher/grpc"
	batchermock "github.com/Layr-Labs/eigenda/disperser/batcher/mock"
	"github.com/Layr-Labs/eigenda/disperser/common/inmem"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/codec"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/node"
	nodegrpc "github.com/Layr-Labs/eigenda/node/grpc"
	"github.com/Layr-Labs/eigenda/retriever"
	retrievermock "github.com/Layr-Labs/eigenda/retriever/mock"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/Layr-Labs/eigensdk-go/metrics"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/docker/go-units"
	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gammazero/workerpool"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/peer"
)

var (
	logger = test.GetLogger()
	p      *prover.Prover
	v      *verifier.Verifier
	asn    core.AssignmentCoordinator

	gettysburgAddressBytes  = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
	serviceManagerAddress   = gethcommon.HexToAddress("0x0000000000000000000000000000000000000000")
	handleBatchLivenessChan = make(chan time.Time, 1)

	localstackContainer *testbed.LocalStackContainer
	clientConfig        commonaws.ClientConfig
	deployLocalStack    bool
	localStackPort      = "4565"
)

const (
	// Operator configuration
	numOperators = 10

	// Service ports
	disperserGrpcPort = 4000
	encoderPort       = "3100"

	// Quorum parameters
	q0AdversaryThreshold = uint8(80)
	q0QuorumThreshold    = uint8(100)

	// Blob configuration
	testMaxBlobSize = 2 * 1024 * 1024
)

type TestOperator struct {
	Node     *node.Node
	ServerV1 *nodegrpc.Server
	ServerV2 *nodegrpc.ServerV2
}

type TestRetriever struct {
	Server *retriever.Server
}

func TestMain(m *testing.M) {
	// Setup test components
	p, v = mustMakeTestComponents()
	asn = &core.StdAssignmentCoordinator{}

	// Run tests
	code := m.Run()

	// Cleanup global resources
	if deployLocalStack && localstackContainer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := localstackContainer.Terminate(ctx); err != nil {
			logger.Error("Failed to terminate localstack container", "error", err)
		}
	}

	os.Exit(code)
}

func TestDispersalAndRetrieval(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	baseCtx := t.Context()

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 3000,
		},
	}
	ctx := peer.NewContext(baseCtx, p)

	cst, err := coremock.MakeChainDataMock(map[uint8]int{
		0: numOperators,
		1: numOperators,
		2: numOperators,
	})
	require.NoError(t, err)

	cst.On("GetCurrentBlockNumber").Return(uint(10), nil)

	store := inmem.NewBlobStore()
	dis := mustMakeDisperser(t, cst, store)
	go func() {
		_ = dis.encoderServer.Start()
	}()
	t.Cleanup(func() {
		dis.encoderServer.Close()
	})

	ops := mustMakeOperators(t, cst)

	gethClient, _ := mustMakeRetriever()

	for _, op := range ops {
		idStr := hexutil.Encode(op.Node.Config.ID[:])
		fmt.Println("Operator: ", idStr)

		fmt.Println("Starting server")
		err = nodegrpc.RunServers(op.ServerV1, op.ServerV2, op.Node.Config, logger)
		require.NoError(t, err)
	}

	blob := makeTestBlob()
	requestedAt := uint64(time.Now().UnixNano())
	metadataKey, err := store.StoreBlob(ctx, &blob, requestedAt)
	require.NoError(t, err)
	out := make(chan batcher.EncodingResultOrStatus)
	err = dis.batcher.EncodingStreamer.RequestEncoding(baseCtx, out)
	require.NoError(t, err)
	err = dis.batcher.EncodingStreamer.ProcessEncodedBlobs(baseCtx, <-out)
	require.NoError(t, err)
	dis.batcher.EncodingStreamer.Pool.StopWait()

	txn := types.NewTransaction(0, gethcommon.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	dis.transactor.On("BuildConfirmBatchTxn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(txn, nil)
	dis.txnManager.On("ProcessTransaction").Return(nil)

	err = dis.batcher.HandleSingleBatch(ctx)
	require.NoError(t, err)
	require.Greater(t, len(dis.txnManager.Requests), 0)
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
	err = dis.batcher.ProcessConfirmedBatch(ctx, &batcher.ReceiptOrErr{
		Receipt:  receipt,
		Err:      nil,
		Metadata: dis.txnManager.Requests[len(dis.txnManager.Requests)-1].Metadata,
	})
	require.NoError(t, err)

	// Check that the blob was processed
	metadata, err := store.GetBlobMetadata(ctx, metadataKey)
	require.NoError(t, err)
	require.Equal(t, metadataKey, metadata.GetBlobKey())
	require.Equal(t, disperser.Confirmed, metadata.BlobStatus)

	isConfirmed, err := metadata.IsConfirmed()
	require.NoError(t, err)
	require.True(t, isConfirmed)
	batchHeaderHash := metadata.ConfirmationInfo.BatchHeaderHash[:]
	txHash := gethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000")
	topics := [][]gethcommon.Hash{
		{common.BatchConfirmedEventSigHash},
		{gethcommon.BytesToHash(batchHeaderHash)},
	}
	calldata, err := hex.DecodeString("7794965a000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000560000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000016400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000001c01b4136a161225e9cebe4e2c561148043b2fde423fc5b64e01d897d0fb7970a142d5474fb609bda1b747bdb5c47375d5819000e3c5cbc75baf55b19849410a2610de9c40eb95b49aca940e0bec6ae8b2868855a6324d04d864cbfa61128cf06a51c069e5a0c490c5a359086b0a3660c2ea2e4fb50722bec1ef593c5245413e4cd0a3c7e490348fb279ccb58f91a3bd494511c2ab0321e3922a0cd26012ef3133c043acb758e735db805d360196f3fc89a6395a4b174c19b981afb7f64c2b1193e0000000000000000000000000000000000000000000000000000000000000220000000000000000000000000000000000000000000000000000000000000026000000000000000000000000000000000000000000000000000000000000002a0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001170c867415fef7db6d88e37598228f43de085616a25939dacbb6b5900f680c7f1d582c9ea38023afb08f368ea93692d17946619d9cf5f3c4d7b3c0cff1a92dff0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000")
	require.NoError(t, err)
	r, ok := new(big.Int).SetString("8ad2b300a012fb0e90dceb8b66fa564717a2d218ca0fd25f11a1875e0153d1d8", 16)
	require.True(t, ok)
	s, ok := new(big.Int).SetString("1accb1e1c69fa07bd4237d92143275960b24eec780862a673d54ffaaa5e77f9b", 16)
	require.True(t, ok)
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
	require.NoError(t, err)

	blobLength := encoding.GetBlobLength(uint32(len(blob.Data)))
	chunkLength, err := asn.CalculateChunkLength(operatorState, uint(blobLength), 0, blob.RequestHeader.SecurityParams[0])
	require.NoError(t, err)

	blobQuorumInfo := &core.BlobQuorumInfo{
		SecurityParam: core.SecurityParam{
			QuorumID:              0,
			AdversaryThreshold:    q0AdversaryThreshold,
			ConfirmationThreshold: q0QuorumThreshold,
		},
		ChunkLength: chunkLength,
	}

	assignments, info, err := asn.GetAssignments(operatorState, uint(blobLength), blobQuorumInfo)
	require.NoError(t, err)

	var indices []encoding.ChunkNumber
	var chunks []*encoding.Frame
	var blobHeader *core.BlobHeader
	for _, op := range ops {

		fmt.Println("Processing operator: ", hexutil.Encode(op.Node.Config.ID[:]))

		// check that blob headers can be retrieved from operators
		headerReply, err := op.ServerV1.GetBlobHeader(ctx, &nodepb.GetBlobHeaderRequest{
			BatchHeaderHash: batchHeaderHash,
			BlobIndex:       metadata.ConfirmationInfo.BlobIndex,
			QuorumId:        uint32(0),
		})
		require.NoError(t, err)
		actualCommitment := &encoding.G1Commitment{
			X: *new(fp.Element).SetBytes(headerReply.GetBlobHeader().GetCommitment().GetX()),
			Y: *new(fp.Element).SetBytes(headerReply.GetBlobHeader().GetCommitment().GetY()),
		}
		var actualLengthCommitment, actualLengthProof encoding.G2Commitment
		actualLengthCommitment.X.A0.SetBytes(headerReply.GetBlobHeader().GetLengthCommitment().GetXA0())
		actualLengthCommitment.X.A1.SetBytes(headerReply.GetBlobHeader().GetLengthCommitment().GetXA1())
		actualLengthCommitment.Y.A0.SetBytes(headerReply.GetBlobHeader().GetLengthCommitment().GetYA0())
		actualLengthCommitment.Y.A1.SetBytes(headerReply.GetBlobHeader().GetLengthCommitment().GetYA1())
		actualLengthProof.X.A0.SetBytes(headerReply.GetBlobHeader().GetLengthProof().GetXA0())
		actualLengthProof.X.A1.SetBytes(headerReply.GetBlobHeader().GetLengthProof().GetXA1())
		actualLengthProof.Y.A0.SetBytes(headerReply.GetBlobHeader().GetLengthProof().GetYA0())
		actualLengthProof.Y.A1.SetBytes(headerReply.GetBlobHeader().GetLengthProof().GetYA1())

		require.Equal(t, metadata.ConfirmationInfo.BlobCommitment.Commitment, actualCommitment)
		require.Equal(t, metadata.ConfirmationInfo.BlobCommitment.LengthCommitment, &actualLengthCommitment)
		require.Equal(t, metadata.ConfirmationInfo.BlobCommitment.LengthProof, &actualLengthProof)
		require.Equal(t, uint32(metadata.ConfirmationInfo.BlobCommitment.Length), headerReply.GetBlobHeader().GetLength())
		require.Len(t, headerReply.GetBlobHeader().GetQuorumHeaders(), 1)
		require.Equal(t, uint32(0), headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetQuorumId())
		require.Equal(t, uint32(q0QuorumThreshold),
			headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetConfirmationThreshold())
		require.Equal(t, uint32(q0AdversaryThreshold),
			headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetAdversaryThreshold())
		require.Greater(t, headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetChunkLength(), uint32(0))

		if blobHeader == nil {
			blobHeader, err = core.BlobHeaderFromProtobuf(headerReply.GetBlobHeader())
			require.NoError(t, err)
		}

		// check that chunks can be retrieved from operators
		chunksReply, err := op.ServerV1.RetrieveChunks(ctx, &nodepb.RetrieveChunksRequest{
			BatchHeaderHash: batchHeaderHash,
			BlobIndex:       metadata.ConfirmationInfo.BlobIndex,
			QuorumId:        uint32(0),
		})

		require.NoError(t, err)
		assignment, ok := assignments[op.Node.Config.ID]
		require.True(t, ok)
		for _, data := range chunksReply.GetChunks() {
			chunk, err := new(encoding.Frame).DeserializeGob(data)
			require.NoError(t, err)
			chunks = append(chunks, chunk)
		}
		require.Len(t, chunksReply.GetChunks(), int(assignments[op.Node.Config.ID].NumChunks))
		indices = append(indices, assignment.GetIndices()...)
	}

	encodingParams := encoding.ParamsFromMins(uint64(chunkLength), info.TotalChunks)
	require.NoError(t, err)
	recovered, err := v.Decode(chunks, indices, encodingParams, uint64(blobHeader.Length)*encoding.BYTES_PER_SYMBOL)
	require.NoError(t, err)

	restored := codec.RemoveEmptyByteFromPaddedBytes(recovered)

	restored = bytes.TrimRight(restored, "\x00")
	require.Equal(t, gettysburgAddressBytes, restored[:len(gettysburgAddressBytes)])
}

func mustMakeTestComponents() (*prover.Prover, *verifier.Verifier) {
	config := &kzg.KzgConfig{
		G1Path:          "../../resources/srs/g1.point",
		G2Path:          "../../resources/srs/g2.point",
		CacheDir:        "../../resources/srs/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}

	p, err := prover.NewProver(config, nil)
	if err != nil {
		panic(fmt.Errorf("failed to create prover: %w", err))
	}
	v, err := verifier.NewVerifier(config, nil)
	if err != nil {
		panic(fmt.Errorf("failed to create verifier: %w", err))
	}

	return p, v
}

func makeTestBlob() core.Blob {
	return core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: []*core.SecurityParam{
				{
					QuorumID:              0,
					AdversaryThreshold:    q0AdversaryThreshold,
					ConfirmationThreshold: q0QuorumThreshold,
				},
			},
		},
		Data: codec.ConvertByPaddingEmptyByte(gettysburgAddressBytes),
	}
}

type TestDisperser struct {
	batcher       *batcher.Batcher
	server        *apiserver.DispersalServer
	encoderServer *encoder.EncoderServer
	transactor    *coremock.MockWriter
	txnManager    *batchermock.MockTxnManager
}

func mustMakeDisperser(t *testing.T, cst core.IndexedChainState, store disperser.BlobStore) *TestDisperser {
	ctx := t.Context()

	dispatcherConfig := &dispatcher.Config{
		Timeout: time.Second,
	}
	batcherMetrics := batcher.NewMetrics("9100", logger)
	dispatcher := dispatcher.NewDispatcher(dispatcherConfig, logger, batcherMetrics.DispatcherMetrics)

	transactor := &coremock.MockWriter{}
	transactor.On("OperatorIDToAddress").Return(gethcommon.Address{}, nil)
	agg, err := core.NewStdSignatureAggregator(logger, transactor)
	require.NoError(t, err)

	batcherConfig := batcher.Config{
		PullInterval:             5 * time.Second,
		NumConnections:           1,
		EncoderSocket:            fmt.Sprintf("localhost:%s", encoderPort),
		EncodingRequestQueueSize: 100,
		SRSOrder:                 3000,
	}
	timeoutConfig := batcher.TimeoutConfig{
		EncodingTimeout:         10 * time.Second,
		AttestationTimeout:      10 * time.Second,
		BatchAttestationTimeout: 12 * time.Second,
		ChainReadTimeout:        10 * time.Second,
		ChainWriteTimeout:       10 * time.Second,
		TxnBroadcastTimeout:     10 * time.Second,
	}

	p0, _ := mustMakeTestComponents()

	metrics := encoder.NewMetrics(prometheus.NewRegistry(), "9000", logger)
	grpcEncoder := encoder.NewEncoderServer(encoder.ServerConfig{
		GrpcPort:              encoderPort,
		MaxConcurrentRequests: 16,
		RequestPoolSize:       32,
	}, logger, p0, metrics, grpcprom.NewServerMetrics())

	encoderClient, err := encoder.NewEncoderClient(batcherConfig.EncoderSocket, 10*time.Second)
	require.NoError(t, err, "failed to create encoder client")

	finalizer := batchermock.NewFinalizer()
	disperserMetrics := disperser.NewMetrics(prometheus.NewRegistry(), "9100", logger)
	txnManager := batchermock.NewTxnManager()

	batcher, err := batcher.NewBatcher(batcherConfig, timeoutConfig, store, dispatcher, cst, asn, encoderClient, agg, &commonmock.MockEthClient{}, finalizer, transactor, txnManager, logger, batcherMetrics, handleBatchLivenessChan)
	require.NoError(t, err, "failed to create batcher")

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
	tx := &coremock.MockWriter{}
	tx.On("GetCurrentBlockNumber").Return(uint64(100), nil)
	tx.On("GetQuorumCount").Return(1, nil)

	// this is disperser client's private key used in tests
	privateKey, err := crypto.HexToECDSA("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded") // Remove "0x" prefix
	require.NoError(t, err, "failed to create ECDSA private key from hex string")

	publicKey := crypto.PubkeyToAddress(privateKey.PublicKey)
	mockState := &coremock.MockOnchainPaymentState{}
	reservationLimit := uint64(1024)
	paymentLimit := big.NewInt(512)
	mockState.On("GetReservedPaymentByAccount", mock.Anything, mock.MatchedBy(func(account gethcommon.Address) bool {
		return account == publicKey
	})).Return(&core.ReservedPayment{SymbolsPerSecond: reservationLimit, StartTimestamp: 0, EndTimestamp: math.MaxUint32, QuorumSplits: []byte{50, 50}, QuorumNumbers: []uint8{0, 1}}, nil)
	mockState.On("GetReservedPaymentByAccount", mock.Anything, mock.Anything).Return(&core.ReservedPayment{}, errors.New("reservation not found"))

	mockState.On("GetOnDemandPaymentByAccount", mock.Anything, mock.MatchedBy(func(account gethcommon.Address) bool {
		return account == publicKey
	})).Return(&core.OnDemandPayment{CumulativePayment: paymentLimit}, nil)
	mockState.On("GetOnDemandPaymentByAccount", mock.Anything, mock.Anything).Return(&core.OnDemandPayment{}, errors.New("payment not found"))
	mockState.On("GetOnDemandQuorumNumbers", mock.Anything).Return([]uint8{0, 1}, nil)
	mockState.On("GetGlobalSymbolsPerSecond", mock.Anything).Return(uint64(1024), nil)
	mockState.On("GetPricePerSymbol", mock.Anything).Return(uint32(1), nil)
	mockState.On("GetMinNumSymbols", mock.Anything).Return(uint32(128), nil)
	mockState.On("GetReservationWindow", mock.Anything).Return(uint32(60), nil)
	mockState.On("RefreshOnchainPaymentState", mock.Anything).Return(nil).Maybe()

	deployLocalStack = (os.Getenv("DEPLOY_LOCALSTACK") != "false")
	if !deployLocalStack {
		localStackPort = os.Getenv("LOCALSTACK_PORT")
	}

	if deployLocalStack {
		var err error
		localstackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       localStackPort,
			Logger:         logger,
		})
		require.NoError(t, err, "failed to start localstack container")
	}

	clientConfig = commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}

	table_names := []string{"reservations_integration", "ondemand_integration", "global_integration"}

	err = meterer.CreateReservationTable(clientConfig, table_names[0])
	require.NoError(t, err, "failed to create reservation table")

	err = meterer.CreateOnDemandTable(clientConfig, table_names[1])
	require.NoError(t, err, "failed to create ondemand table")

	err = meterer.CreateGlobalReservationTable(clientConfig, table_names[2])
	require.NoError(t, err, "failed to create global reservation table")

	meteringStore, err := meterer.NewDynamoDBMeteringStore(
		clientConfig,
		table_names[0],
		table_names[1],
		table_names[2],
		logger,
	)
	require.NoError(t, err, "failed to create metering store")

	mockState.On("RefreshOnchainPaymentState", mock.Anything).Return(nil).Maybe()
	err = mockState.RefreshOnchainPaymentState(ctx)
	require.NoError(t, err, "failed to refresh on-chain payment state")

	mt := meterer.NewMeterer(meterer.Config{}, mockState, meteringStore, logger)
	server := apiserver.NewDispersalServer(serverConfig, store, tx, logger, disperserMetrics, grpcprom.NewServerMetrics(), mt, ratelimiter, rateConfig, testMaxBlobSize)

	return &TestDisperser{
		batcher:       batcher,
		server:        server,
		encoderServer: grpcEncoder,
		transactor:    transactor,
		txnManager:    txnManager,
	}
}

func mustMakeOperators(t *testing.T, cst *coremock.ChainDataMock) map[core.OperatorID]TestOperator {
	ctx := t.Context()
	bn := uint(0)
	state := cst.GetTotalOperatorState(ctx, bn)

	ops := make(map[core.OperatorID]TestOperator, len(state.IndexedOperators))

	setRegisteredQuorums := true
	for id, op := range state.PrivateOperators {

		idStr := hexutil.Encode(id[:])
		fmt.Println("Operator: ", idStr)

		dbPath := fmt.Sprintf("testdata/%v/db", idStr)
		logPath := fmt.Sprintf("testdata/%v/log", idStr)

		err := os.RemoveAll(dbPath)
		require.NoError(t, err, "failed to remove db path")

		err = os.RemoveAll(logPath)
		require.NoError(t, err, "failed to remove log path")

		err = os.MkdirAll(dbPath, os.ModePerm)
		require.NoError(t, err, "failed to create db path")

		err = os.MkdirAll(logPath, os.ModePerm)
		require.NoError(t, err, "failed to create log path")

		registeredQuorums := []core.QuorumID{}
		if setRegisteredQuorums {
			registeredQuorums = []core.QuorumID{0}
			setRegisteredQuorums = false
		}

		config := &node.Config{
			Hostname:                            op.Host,
			DispersalPort:                       op.DispersalPort,
			RetrievalPort:                       op.RetrievalPort,
			InternalRetrievalPort:               op.RetrievalPort,
			InternalDispersalPort:               op.DispersalPort,
			V2DispersalPort:                     op.V2DispersalPort,
			V2RetrievalPort:                     op.V2RetrievalPort,
			EnableMetrics:                       false,
			Timeout:                             10,
			ExpirationPollIntervalSec:           10,
			DbPath:                              dbPath,
			LogPath:                             logPath,
			ID:                                  id,
			QuorumIDList:                        registeredQuorums,
			DispersalAuthenticationKeyCacheSize: 1024,
			DisableDispersalAuthentication:      false,
			RelayMaxMessageSize:                 units.GiB,
			EnableV1:                            true,
			EnableV2:                            false,
		}

		// creating a new instance of encoder instead of sharing enc because enc is not thread safe
		_, v0 := mustMakeTestComponents()
		val := core.NewShardValidator(v0, asn, cst, id)

		tx := &coremock.MockWriter{}
		tx.On("RegisterBLSPublicKey").Return(nil)
		tx.On("RegisterOperator").Return(nil)
		tx.On("GetRegisteredQuorumIdsForOperator").Return(registeredQuorums, nil)
		tx.On("UpdateOperatorSocket").Return(nil)
		tx.On("GetBlockStaleMeasure").Return(nil)
		tx.On("GetStoreDurationBlocks").Return(nil)
		tx.On("OperatorIDToAddress").Return(gethcommon.Address{1}, nil)
		socket := core.MakeOperatorSocket(config.Hostname, config.DispersalPort, config.RetrievalPort, config.V2DispersalPort, config.V2RetrievalPort)
		tx.On("GetOperatorSocket", mock.Anything, mock.Anything).Return(socket.String(), nil)

		noopMetrics := metrics.NewNoopMetrics()
		reg := prometheus.NewRegistry()
		metrics := node.NewMetrics(noopMetrics, reg, logger, ":9090", config.ID, -1, tx, cst)
		store, err := node.NewLevelDBStore(
			config.DbPath+"/chunk",
			logger,
			metrics,
			1e9,
			true,
			false,
			1e9)
		require.NoError(t, err, "failed to create leveldb store")

		mockOperatorSocketsFilterer := &coremock.MockOperatorSocketsFilterer{}

		mockSocketChan := make(chan string)
		mockOperatorSocketsFilterer.On("WatchOperatorSocketUpdate").Return(mockSocketChan, nil)

		pubIPProvider := pubip.CustomProvider(
			pubip.RequestDoerFunc(func(req *http.Request) (*http.Response, error) {
				w := httptest.NewRecorder()
				_, _ = w.WriteString("8.8.8.8")
				return w.Result(), nil
			}), "custom", "")

		node := &node.Node{
			Config:                  config,
			Logger:                  logger,
			KeyPair:                 op.KeyPair,
			BLSSigner:               op.Signer,
			Metrics:                 metrics,
			Store:                   store,
			ChainState:              cst,
			Validator:               val,
			Transactor:              tx,
			PubIPProvider:           pubIPProvider,
			OperatorSocketsFilterer: mockOperatorSocketsFilterer,
			ValidationPool:          workerpool.New(1),
		}

		rateLimiter := &commonmock.NoopRatelimiter{}

		// TODO(cody-littley): Once we switch this test to use the v2 disperser, we will need to properly set up
		//  the disperser's public/private keys for signing StoreChunks() requests
		disperserAddress := gethcommon.Address{}
		reader := &coremock.MockWriter{}
		reader.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

		// Create listeners with OS-allocated ports for testing
		v1DispersalListener, err := net.Listen("tcp", "0.0.0.0:0")
		require.NoError(t, err)
		v1RetrievalListener, err := net.Listen("tcp", "0.0.0.0:0")
		require.NoError(t, err)
		v2DispersalListener, err := net.Listen("tcp", "0.0.0.0:0")
		require.NoError(t, err)
		v2RetrievalListener, err := net.Listen("tcp", "0.0.0.0:0")
		require.NoError(t, err)

		serverV1 := nodegrpc.NewServer(
			config,
			node,
			logger,
			rateLimiter,
			version.DefaultVersion(),
			v1DispersalListener,
			v1RetrievalListener,
		)
		serverV2, err := nodegrpc.NewServerV2(
			ctx,
			config,
			node,
			logger,
			rateLimiter,
			prometheus.NewRegistry(),
			reader,
			version.DefaultVersion(),
			v2DispersalListener,
			v2RetrievalListener)
		require.NoError(t, err)

		ops[id] = TestOperator{
			Node:     node,
			ServerV1: serverV1,
			ServerV2: serverV2,
		}
	}

	return ops
}

func mustMakeRetriever() (*commonmock.MockEthClient, *TestRetriever) {
	config := &retriever.Config{
		Timeout: 5 * time.Second,
	}
	gethClient := &commonmock.MockEthClient{}
	retrievalClient := &clientsmock.MockRetrievalClient{}
	chainClient := retrievermock.NewMockChainClient()
	server := retriever.NewServer(config, logger, retrievalClient, chainClient)

	return gethClient, &TestRetriever{
		Server: server,
	}
}
