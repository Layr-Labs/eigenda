package integration_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Layr-Labs/eigenda/common/pubip"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ory/dockertest/v3"

	clientsmock "github.com/Layr-Labs/eigenda/api/clients/mock"
	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/core/meterer"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	dispatcher "github.com/Layr-Labs/eigenda/disperser/batcher/grpc"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/retriever"
	retrievermock "github.com/Layr-Labs/eigenda/retriever/mock"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/metrics"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/peer"

	"github.com/Layr-Labs/eigenda/common"
	commonmock "github.com/Layr-Labs/eigenda/common/mock"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	batchermock "github.com/Layr-Labs/eigenda/disperser/batcher/mock"
	"github.com/Layr-Labs/eigenda/disperser/common/inmem"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/node"
	"github.com/Layr-Labs/eigenda/node/grpc"
	nodegrpc "github.com/Layr-Labs/eigenda/node/grpc"

	nodepb "github.com/Layr-Labs/eigenda/api/grpc/node"

	"github.com/ethereum/go-ethereum"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	p   encoding.Prover
	v   encoding.Verifier
	asn core.AssignmentCoordinator

	gettysburgAddressBytes  = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
	serviceManagerAddress   = gethcommon.HexToAddress("0x0000000000000000000000000000000000000000")
	handleBatchLivenessChan = make(chan time.Time, 1)

	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	clientConfig       commonaws.ClientConfig

	deployLocalStack bool
	localStackPort   = "4565"
)

const (
	numOperators         = 10
	disperserGrpcPort    = 4000
	encoderPort          = "3100"
	q0AdversaryThreshold = uint8(80)
	q0QuorumThreshold    = uint8(100)
	testMaxBlobSize      = 2 * 1024 * 1024
)

func init() {
	p, v = mustMakeTestComponents()
	asn = &core.StdAssignmentCoordinator{}
}

// makeTestEncoder makes an encoder currently using the only supported backend.
func mustMakeTestComponents() (encoding.Prover, encoding.Verifier) {
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
		log.Fatal(err)
	}

	v, err := verifier.NewVerifier(config, nil)
	if err != nil {
		log.Fatal(err)
	}

	return p, v
}

func mustMakeTestBlob() core.Blob {
	blob := core.Blob{
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
	return blob
}

type TestDisperser struct {
	batcher       *batcher.Batcher
	server        *apiserver.DispersalServer
	encoderServer *encoder.EncoderServer
	transactor    *coremock.MockWriter
	txnManager    *batchermock.MockTxnManager
}

func mustMakeDisperser(t *testing.T, cst core.IndexedChainState, store disperser.BlobStore, logger logging.Logger) TestDisperser {
	dispatcherConfig := &dispatcher.Config{
		Timeout: time.Second,
	}
	batcherMetrics := batcher.NewMetrics("9100", logger)
	dispatcher := dispatcher.NewDispatcher(dispatcherConfig, logger, batcherMetrics.DispatcherMetrics)

	transactor := &coremock.MockWriter{}
	transactor.On("OperatorIDToAddress").Return(gethcommon.Address{}, nil)
	agg, err := core.NewStdSignatureAggregator(logger, transactor)
	assert.NoError(t, err)

	batcherConfig := batcher.Config{
		PullInterval:             5 * time.Second,
		NumConnections:           1,
		EncoderSocket:            fmt.Sprintf("localhost:%s", encoderPort),
		EncodingRequestQueueSize: 100,
		SRSOrder:                 3000,
	}
	timeoutConfig := batcher.TimeoutConfig{
		EncodingTimeout:     10 * time.Second,
		AttestationTimeout:  10 * time.Second,
		ChainReadTimeout:    10 * time.Second,
		ChainWriteTimeout:   10 * time.Second,
		TxnBroadcastTimeout: 10 * time.Second,
	}

	p0, _ := mustMakeTestComponents()
	metrics := encoder.NewMetrics(prometheus.NewRegistry(), "9000", logger)
	grpcEncoder := encoder.NewEncoderServer(encoder.ServerConfig{
		GrpcPort:              encoderPort,
		MaxConcurrentRequests: 16,
		RequestPoolSize:       32,
	}, logger, p0, metrics, grpcprom.NewServerMetrics())

	encoderClient, err := encoder.NewEncoderClient(batcherConfig.EncoderSocket, 10*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	finalizer := batchermock.NewFinalizer()

	disperserMetrics := disperser.NewMetrics(prometheus.NewRegistry(), "9100", logger)
	txnManager := batchermock.NewTxnManager()

	batcher, err := batcher.NewBatcher(batcherConfig, timeoutConfig, store, dispatcher, cst, asn, encoderClient, agg, &commonmock.MockEthClient{}, finalizer, transactor, txnManager, logger, batcherMetrics, handleBatchLivenessChan)
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
	tx := &coremock.MockWriter{}
	tx.On("GetCurrentBlockNumber").Return(uint64(100), nil)
	tx.On("GetQuorumCount").Return(1, nil)

	// this is disperser client's private key used in tests
	privateKey, err := crypto.HexToECDSA("0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded") // Remove "0x" prefix
	if err != nil {
		panic("failed to convert hex to ECDSA")
	}
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

	deployLocalStack = !(os.Getenv("DEPLOY_LOCALSTACK") == "false")
	if !deployLocalStack {
		localStackPort = os.Getenv("LOCALSTACK_PORT")
	}

	if deployLocalStack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localStackPort)
		if err != nil {
			teardown()
			panic("failed to start localstack container")
		}
	}

	clientConfig = commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}

	table_names := []string{"reservations_integration", "ondemand_integration", "global_integration"}

	err = meterer.CreateReservationTable(clientConfig, table_names[0])
	if err != nil {
		teardown()
		panic("failed to create reservation table")
	}
	err = meterer.CreateOnDemandTable(clientConfig, table_names[1])
	if err != nil {
		teardown()
		panic("failed to create ondemand table")
	}
	err = meterer.CreateGlobalReservationTable(clientConfig, table_names[2])
	if err != nil {
		teardown()
		panic("failed to create global reservation table")
	}

	offchainStore, err := meterer.NewOffchainStore(
		clientConfig,
		table_names[0],
		table_names[1],
		table_names[2],
		logger,
	)
	if err != nil {
		panic("failed to create offchain store")
	}

	mockState.On("RefreshOnchainPaymentState", mock.Anything).Return(nil).Maybe()
	if err := mockState.RefreshOnchainPaymentState(context.Background()); err != nil {
		panic("failed to make initial query to the on-chain state")
	}

	mt := meterer.NewMeterer(meterer.Config{}, mockState, offchainStore, logger)
	server := apiserver.NewDispersalServer(serverConfig, store, tx, logger, disperserMetrics, grpcprom.NewServerMetrics(), mt, ratelimiter, rateConfig, testMaxBlobSize)

	return TestDisperser{
		batcher:       batcher,
		server:        server,
		encoderServer: grpcEncoder,
		transactor:    transactor,
		txnManager:    txnManager,
	}
}

type TestOperator struct {
	Node     *node.Node
	ServerV1 *nodegrpc.Server
	ServerV2 *nodegrpc.ServerV2
}

func mustMakeOperators(t *testing.T, cst *coremock.ChainDataMock, logger logging.Logger) map[core.OperatorID]TestOperator {
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
			Hostname:                            op.Host,
			DispersalPort:                       op.DispersalPort,
			RetrievalPort:                       op.RetrievalPort,
			InternalRetrievalPort:               op.RetrievalPort,
			InternalDispersalPort:               op.DispersalPort,
			EnableMetrics:                       false,
			Timeout:                             10,
			ExpirationPollIntervalSec:           10,
			DbPath:                              dbPath,
			LogPath:                             logPath,
			PrivateBls:                          string(op.KeyPair.GetPubKeyG1().Serialize()),
			ID:                                  id,
			QuorumIDList:                        registeredQuorums,
			DispersalAuthenticationKeyCacheSize: 1024,
			DisableDispersalAuthentication:      false,
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
		socket := core.MakeOperatorSocket(config.Hostname, config.DispersalPort, config.RetrievalPort)
		tx.On("GetOperatorSocket", mock.Anything, mock.Anything).Return(socket.String(), nil)

		noopMetrics := metrics.NewNoopMetrics()
		reg := prometheus.NewRegistry()
		metrics := node.NewMetrics(noopMetrics, reg, logger, ":9090", config.ID, -1, tx, cst)
		store, err := node.NewLevelDBStore(config.DbPath+"/chunk", logger, metrics, 1e9, 1e9)
		if err != nil {
			t.Fatal(err)
		}

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

		rateLimiter := &commonmock.NoopRatelimiter{}

		// TODO(cody-littley): Once we switch this test to use the v2 disperser, we will need to properly set up
		//  the disperser's public/private keys for signing StoreChunks() requests
		disperserAddress := gethcommon.Address{}
		reader := &coremock.MockWriter{}
		reader.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

		serverV1 := nodegrpc.NewServer(config, n, logger, rateLimiter)
		serverV2, err := nodegrpc.NewServerV2(
			context.Background(),
			config,
			n,
			logger,
			rateLimiter,
			prometheus.NewRegistry(),
			reader)
		require.NoError(t, err)

		ops[id] = TestOperator{
			Node:     n,
			ServerV1: serverV1,
			ServerV2: serverV2,
		}
	}

	return ops
}

type TestRetriever struct {
	Server *retriever.Server
}

func mustMakeRetriever(cst core.IndexedChainState, logger logging.Logger) (*commonmock.MockEthClient, TestRetriever) {
	config := &retriever.Config{
		Timeout: 5 * time.Second,
	}
	gethClient := &commonmock.MockEthClient{}
	retrievalClient := &clientsmock.MockRetrievalClient{}
	chainClient := retrievermock.NewMockChainClient()
	server := retriever.NewServer(config, logger, retrievalClient, cst, chainClient)

	return gethClient, TestRetriever{
		Server: server,
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestDispersalAndRetrieval(t *testing.T) {

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 3000,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	cst, err := coremock.MakeChainDataMock(map[uint8]int{
		0: numOperators,
		1: numOperators,
		2: numOperators,
	})
	assert.NoError(t, err)

	cst.On("GetCurrentBlockNumber").Return(uint(10), nil)

	logger := logging.NewNoopLogger()
	assert.NoError(t, err)
	store := inmem.NewBlobStore()
	dis := mustMakeDisperser(t, cst, store, logger)
	go func() {
		_ = dis.encoderServer.Start()
	}()
	t.Cleanup(func() {
		dis.encoderServer.Close()
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
		err = grpc.RunServers(op.ServerV1, op.ServerV2, op.Node.Config, logger)
		assert.NoError(t, err)
	}

	blob := mustMakeTestBlob()
	requestedAt := uint64(time.Now().UnixNano())
	metadataKey, err := store.StoreBlob(ctx, &blob, requestedAt)
	assert.NoError(t, err)
	out := make(chan batcher.EncodingResultOrStatus)
	err = dis.batcher.EncodingStreamer.RequestEncoding(context.Background(), out)
	assert.NoError(t, err)
	err = dis.batcher.EncodingStreamer.ProcessEncodedBlobs(context.Background(), <-out)
	assert.NoError(t, err)
	dis.batcher.EncodingStreamer.Pool.StopWait()

	txn := types.NewTransaction(0, gethcommon.Address{}, big.NewInt(0), 0, big.NewInt(0), nil)
	dis.transactor.On("BuildConfirmBatchTxn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(txn, nil)
	dis.txnManager.On("ProcessTransaction").Return(nil)

	err = dis.batcher.HandleSingleBatch(ctx)
	assert.NoError(t, err)
	assert.Greater(t, len(dis.txnManager.Requests), 0)
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

	blobLength := encoding.GetBlobLength(uint(len(blob.Data)))
	chunkLength, err := asn.CalculateChunkLength(operatorState, blobLength, 0, blob.RequestHeader.SecurityParams[0])
	assert.NoError(t, err)

	blobQuorumInfo := &core.BlobQuorumInfo{
		SecurityParam: core.SecurityParam{
			QuorumID:              0,
			AdversaryThreshold:    q0AdversaryThreshold,
			ConfirmationThreshold: q0QuorumThreshold,
		},
		ChunkLength: chunkLength,
	}

	assignments, info, err := asn.GetAssignments(operatorState, blobLength, blobQuorumInfo)
	assert.NoError(t, err)

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
		assert.NoError(t, err)
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

		assert.Equal(t, metadata.ConfirmationInfo.BlobCommitment.Commitment, actualCommitment)
		assert.Equal(t, metadata.ConfirmationInfo.BlobCommitment.LengthCommitment, &actualLengthCommitment)
		assert.Equal(t, metadata.ConfirmationInfo.BlobCommitment.LengthProof, &actualLengthProof)
		assert.Equal(t, uint32(metadata.ConfirmationInfo.BlobCommitment.Length), headerReply.GetBlobHeader().GetLength())
		assert.Len(t, headerReply.GetBlobHeader().GetQuorumHeaders(), 1)
		assert.Equal(t, uint32(0), headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetQuorumId())
		assert.Equal(t, uint32(q0QuorumThreshold), headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetConfirmationThreshold())
		assert.Equal(t, uint32(q0AdversaryThreshold), headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetAdversaryThreshold())
		assert.Greater(t, headerReply.GetBlobHeader().GetQuorumHeaders()[0].GetChunkLength(), uint32(0))

		if blobHeader == nil {
			blobHeader, err = core.BlobHeaderFromProtobuf(headerReply.GetBlobHeader())
			assert.NoError(t, err)
		}

		// check that chunks can be retrieved from operators
		chunksReply, err := op.ServerV1.RetrieveChunks(ctx, &nodepb.RetrieveChunksRequest{
			BatchHeaderHash: batchHeaderHash,
			BlobIndex:       metadata.ConfirmationInfo.BlobIndex,
			QuorumId:        uint32(0),
		})

		assert.NoError(t, err)
		assignment, ok := assignments[op.Node.Config.ID]
		assert.True(t, ok)
		for _, data := range chunksReply.GetChunks() {
			chunk, err := new(encoding.Frame).Deserialize(data)
			assert.NoError(t, err)
			chunks = append(chunks, chunk)
		}
		assert.Len(t, chunksReply.GetChunks(), int(assignments[op.Node.Config.ID].NumChunks))
		indices = append(indices, assignment.GetIndices()...)
	}

	encodingParams := encoding.ParamsFromMins(chunkLength, info.TotalChunks)
	assert.NoError(t, err)
	recovered, err := v.Decode(chunks, indices, encodingParams, uint64(blobHeader.Length)*encoding.BYTES_PER_SYMBOL)
	assert.NoError(t, err)

	restored := codec.RemoveEmptyByteFromPaddedBytes(recovered)

	restored = bytes.TrimRight(restored, "\x00")
	assert.Equal(t, gettysburgAddressBytes, restored[:len(gettysburgAddressBytes)])
}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}
