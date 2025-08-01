package node

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	"github.com/Layr-Labs/eigenda/common/pprof"
	"github.com/Layr-Labs/eigenda/common/pubip"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/indexer"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/metrics"
	"github.com/Layr-Labs/eigensdk-go/nodeapi"
	blssigner "github.com/Layr-Labs/eigensdk-go/signer/bls"

	"github.com/gammazero/workerpool"
)

const (
	// The percentage of time in garbage collection in a GC cycle.
	gcPercentageTime = 0.1

	v1CheckPath = "api/v1/operators-info/port-check"
	v2CheckPath = "api/v2/operators/liveness"
)

var (
	// eigenDAUIMap is a mapping for ChainID to the EigenDA UI url.
	eigenDAUIMap = map[string]string{
		"1":     "https://app.eigenlayer.xyz/avs/0x870679e138bcdf293b7ff14dd44b70fc97e12fc0",
		"17000": "https://holesky.eigenlayer.xyz/avs/0xd4a7e1bd8015057293f0d0a557088c286942e84b/operator-set/4294967295",
	}
)

type Node struct {
	Config                  *Config
	Logger                  logging.Logger
	KeyPair                 *core.KeyPair
	Metrics                 *Metrics
	NodeApi                 *nodeapi.NodeApi
	Store                   *Store
	BlacklistStore          BlacklistStore
	ValidatorStore          ValidatorStore
	ChainState              core.ChainState
	Validator               core.ShardValidator
	ValidatorV2             corev2.ShardValidator
	Transactor              core.Writer
	PubIPProvider           pubip.Provider
	OperatorSocketsFilterer indexer.OperatorSocketsFilterer
	ChainID                 *big.Int
	// a worker pool used to download chunk data from the relays
	DownloadPool *workerpool.WorkerPool
	// a worker pool used to validate batches
	ValidationPool *workerpool.WorkerPool

	BLSSigner blssigner.Signer

	RelayClient atomic.Value

	mu            sync.Mutex
	CurrentSocket string

	// BlobVersionParams is a map of blob version parameters loaded from the chain.
	// It is used to determine blob parameters based on the version number.
	BlobVersionParams atomic.Pointer[corev2.BlobVersionParameterMap]

	// TODO: utilize meterer onchain state later to check quorum ID and minimum payments
	// QuorumCount is the number of quorums in the network.
	QuorumCount atomic.Uint32
}

// NewNode creates a new Node with the provided config.
// TODO: better context management, don't just use context.Background() everywhere in here.
func NewNode(
	reg *prometheus.Registry,
	config *Config,
	pubIPProvider pubip.Provider,
	client *geth.InstrumentedEthClient,
	logger logging.Logger,
) (*Node, error) {
	nodeLogger := logger.With("component", "Node")

	socketAddr := fmt.Sprintf(":%d", config.MetricsPort)
	eigenMetrics := metrics.NewEigenMetrics(AppName, socketAddr, reg, logger.With("component", "EigenMetrics"))

	// Make sure config folder exists.
	err := os.MkdirAll(config.DbPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not create DB directory at %s: %w", config.DbPath, err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chainID: %w", err)
	}

	// Create Transactor
	tx, err := eth.NewWriter(logger, client, config.EigenDADirectory, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create writer: %w", err)
	}

	// Create ChainState Client
	cst := eth.NewChainState(tx, client)

	blsSigner, err := blssigner.NewSigner(config.BlsSignerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create BLS signer: %w", err)
	}
	operatorID, err := blsSigner.GetOperatorId()
	if err != nil {
		return nil, fmt.Errorf("failed to get operator ID: %w", err)
	}
	config.ID, err = core.OperatorIDFromHex(operatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to convert operator ID: %w", err)
	}

	// Setup Node Api
	nodeApi := nodeapi.NewNodeApi(AppName, SemVer, ":"+config.NodeApiPort, logger.With("component", "NodeApi"))

	metrics := NewMetrics(eigenMetrics, reg, logger, socketAddr, config.ID, config.OnchainMetricsInterval, tx, cst)

	// Make validator
	config.EncoderConfig.LoadG2Points = false
	v, err := verifier.NewVerifier(&config.EncoderConfig, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create verifier: %w", err)
	}
	asgn := &core.StdAssignmentCoordinator{}
	validator := core.NewShardValidator(v, asgn, cst, config.ID)
	validatorV2 := corev2.NewShardValidator(v, config.ID, logger)

	// Resolve the BLOCK_STALE_MEASURE and STORE_DURATION_BLOCKS.
	var blockStaleMeasure, storeDurationBlocks uint32
	if config.EnableTestMode && config.OverrideBlockStaleMeasure > 0 {
		blockStaleMeasure = uint32(config.OverrideBlockStaleMeasure)
		logger.Info("Test Mode Override!", "blockStaleMeasure", blockStaleMeasure)
	} else {
		staleMeasure, err := tx.GetBlockStaleMeasure(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get BLOCK_STALE_MEASURE: %w", err)
		}
		blockStaleMeasure = staleMeasure
	}
	if config.EnableTestMode && config.OverrideStoreDurationBlocks > 0 {
		storeDurationBlocks = uint32(config.OverrideStoreDurationBlocks)
		logger.Info("Test Mode Override!", "storeDurationBlocks", storeDurationBlocks)
	} else {
		storeDuration, err := tx.GetStoreDurationBlocks(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get STORE_DURATION_BLOCKS: %w", err)
		}
		storeDurationBlocks = storeDuration
	}
	// Create new chunk store
	store, err := NewLevelDBStore(
		config.DbPath+"/chunk",
		logger,
		metrics,
		blockStaleMeasure,
		config.LevelDBDisableSeeksCompactionV1,
		config.LevelDBSyncWritesV1,
		storeDurationBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to create new store: %w", err)
	}

	// If EigenDADirectory is provided, use it to get service manager addresses
	// Otherwise, use the provided address (legacy support; will be removed as a breaking change)
	eigenDAServiceManagerAddr := gethcommon.HexToAddress(config.EigenDAServiceManagerAddr)
	if config.EigenDADirectory != "" && gethcommon.IsHexAddress(config.EigenDADirectory) {
		addressReader, err := eth.NewEigenDADirectoryReader(config.EigenDADirectory, client)
		if err != nil {
			return nil, fmt.Errorf("failed to create address directory reader: %w", err)
		}
		eigenDAServiceManagerAddr, err = addressReader.GetServiceManagerAddress(&bind.CallOpts{Context: context.Background()})
		if err != nil {
			return nil, fmt.Errorf("failed to get service manager address from EigenDADirectory: %w", err)
		}
		if config.EigenDAServiceManagerAddr != "" && eigenDAServiceManagerAddr.String() != config.EigenDAServiceManagerAddr {
			return nil, fmt.Errorf("EigenDAServiceManagerAddr passed in as config (%v) does not match the one retrieved from EigenDADirectory (%v)",
				config.EigenDADirectory, eigenDAServiceManagerAddr.Hex())
		}
	} else {
		logger.Warn("EigenDADirectory is not set or is not a valid address, using provided EigenDAServiceManagerAddr. "+
			"This is deprecated and will be removed in a future release. Please switch to using EigenDADirectory.",
			"EigenDAServiceManagerAddr", eigenDAServiceManagerAddr.Hex(), "EigenDADirectory", config.EigenDADirectory)
	}
	socketsFilterer, err := indexer.NewOperatorSocketsFilterer(eigenDAServiceManagerAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create new operator sockets filterer: %w", err)
	}

	nodeLogger.Info("Creating node",
		"chainID", chainID.String(),
		"operatorID", config.ID.Hex(),
		"dispersalPort", config.DispersalPort,
		"internalDispersalPort", config.InternalDispersalPort,
		"v2DispersalPort", config.V2DispersalPort,
		"internalV2DispersalPort", config.InternalV2DispersalPort,
		"retrievalPort", config.RetrievalPort,
		"internalRetrievalPort", config.InternalRetrievalPort,
		"v2RetrievalPort", config.V2RetrievalPort,
		"internalV2RetrievalPort", config.InternalV2RetrievalPort,
		"churnerUrl", config.ChurnerUrl,
		"quorumIDs", fmt.Sprint(config.QuorumIDList), //nolint:staticcheck // QF1010
		"registerNodeAtStart", config.RegisterNodeAtStart,
		"pubIPCheckInterval", config.PubIPCheckInterval,
		"eigenDAServiceManagerAddr", eigenDAServiceManagerAddr.Hex(),
		"blockStaleMeasure", blockStaleMeasure,
		"storeDurationBlocks", storeDurationBlocks,
		"enableGnarkBundleEncoding", config.EnableGnarkBundleEncoding)

	downloadPoolSize := config.DownloadPoolSize
	if downloadPoolSize < 1 {
		downloadPoolSize = 1
	}
	downloadPool := workerpool.New(downloadPoolSize)

	validationPoolSize := config.NumBatchValidators
	if validationPoolSize < 1 {
		validationPoolSize = 1
	}
	validationPool := workerpool.New(validationPoolSize)

	n := &Node{
		Config:                  config,
		Logger:                  nodeLogger,
		Metrics:                 metrics,
		NodeApi:                 nodeApi,
		Store:                   store,
		BlacklistStore:          nil,
		ChainState:              cst,
		Transactor:              tx,
		Validator:               validator,
		ValidatorV2:             validatorV2,
		PubIPProvider:           pubIPProvider,
		OperatorSocketsFilterer: socketsFilterer,
		ChainID:                 chainID,
		BLSSigner:               blsSigner,
		DownloadPool:            downloadPool,
		ValidationPool:          validationPool,
	}

	if !config.EnableV2 {
		return n, nil
	}

	var blobVersionParams *corev2.BlobVersionParameterMap
	if config.EnableV2 {
		ctx := context.Background()
		// 12s per block
		ttl := time.Duration(blockStaleMeasure+storeDurationBlocks) * 12 * time.Second
		n.ValidatorStore, err = NewValidatorStore(logger, config, time.Now, ttl, reg)
		if err != nil {
			return nil, fmt.Errorf("failed to create new store v2: %w", err)
		}

		blobParams, err := tx.GetAllVersionedBlobParams(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get versioned blob parameters: %w", err)
		}
		blobVersionParams = corev2.NewBlobVersionParameterMap(blobParams)

		relayClientConfig := &relay.RelayClientConfig{
			UseSecureGrpcFlag:  config.RelayUseSecureGrpc,
			OperatorID:         &config.ID,
			MessageSigner:      n.SignMessage,
			MaxGRPCMessageSize: n.Config.RelayMaxMessageSize,
		}

		relayUrlProvider, err := relay.NewRelayUrlProvider(client, tx.GetRelayRegistryAddress())
		if err != nil {
			return nil, fmt.Errorf("create relay url provider: %w", err)
		}

		relayClient, err := relay.NewRelayClient(relayClientConfig, logger, relayUrlProvider)
		if err != nil {
			return nil, fmt.Errorf("failed to create new relay client: %w", err)
		}

		n.RelayClient.Store(relayClient)

		blockNumber, err := tx.GetCurrentBlockNumber(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get block number: %w", err)
		}
		quorumCount, err := tx.GetQuorumCount(ctx, blockNumber)
		if err != nil {
			return nil, fmt.Errorf("failed to get quorum count: %w", err)
		}
		n.QuorumCount.Store(uint32(quorumCount))
	}

	n.BlobVersionParams.Store(blobVersionParams)
	return n, nil
}

// Start starts the Node. If the node is not registered, register it on chain, otherwise just
// update its socket on chain.
func (n *Node) Start(ctx context.Context) error {
	pprofProfiler := pprof.NewPprofProfiler(n.Config.PprofHttpPort, n.Logger)
	if n.Config.EnablePprof {
		go pprofProfiler.Start()
		n.Logger.Info("Enabled pprof for Node", "port", n.Config.PprofHttpPort)
	}
	if n.Config.EnableMetrics {
		n.Metrics.Start()
		n.Logger.Info("Enabled metrics", "socket", n.Metrics.socketAddr)
	}
	if n.Config.EnableNodeApi {
		n.NodeApi.Start()
		n.Logger.Info("Enabled node api", "port", n.Config.NodeApiPort)
	}

	if n.Config.EnableV1 {
		go n.expireLoop()
		go n.checkNodeReachability(v1CheckPath)
	}

	if n.Config.EnableV2 {
		go func() {
			_ = n.RefreshOnchainState(ctx)
		}()
		go n.checkNodeReachability(v2CheckPath)
	}

	// Build the socket based on the hostname/IP provided in the CLI
	socket := string(core.MakeOperatorSocket(
		n.Config.Hostname,
		n.Config.DispersalPort,
		n.Config.RetrievalPort,
		n.Config.V2DispersalPort,
		n.Config.V2RetrievalPort))
	var operator *Operator
	if n.Config.RegisterNodeAtStart {
		n.Logger.Info("Registering node on chain with the following parameters:",
			"operatorId", n.Config.ID.Hex(),
			"hostname", n.Config.Hostname,
			"dispersalPort", n.Config.DispersalPort,
			"v2DispersalPort", n.Config.V2DispersalPort,
			"retrievalPort", n.Config.RetrievalPort,
			"v2RetrievalPort", n.Config.V2RetrievalPort,
			"churnerUrl", n.Config.ChurnerUrl,
			"quorumIds", fmt.Sprintf("%v", n.Config.QuorumIDList))
		privateKey, err := crypto.HexToECDSA(n.Config.EthClientConfig.PrivateKeyString)
		if err != nil {
			return fmt.Errorf("NewClient: cannot parse private key: %w", err)
		}
		operator = &Operator{
			Address:             crypto.PubkeyToAddress(privateKey.PublicKey).Hex(),
			Socket:              socket,
			Timeout:             10 * time.Second,
			PrivKey:             privateKey,
			Signer:              n.BLSSigner,
			OperatorId:          n.Config.ID,
			QuorumIDs:           n.Config.QuorumIDList,
			RegisterNodeAtStart: n.Config.RegisterNodeAtStart,
		}
		churnerClient := NewChurnerClient(n.Config.ChurnerUrl, n.Config.ChurnerUseSecureGrpc, n.Config.Timeout, n.Logger)
		err = RegisterOperator(ctx, operator, n.Transactor, churnerClient, n.Logger)
		if err != nil {
			return fmt.Errorf("failed to register the operator: %w", err)
		}
	} else {
		registeredSocket, err := n.Transactor.GetOperatorSocket(ctx, n.Config.ID)
		// Error out if registration on-chain is a requirement
		if err != nil {
			n.Logger.Warnf("failed to get operator socket: %w", err)
		}
		if registeredSocket != socket {
			n.Logger.Warnf("registered socket %s does not match expected socket %s", registeredSocket, socket)
		}

		eigenDAUrl, ok := eigenDAUIMap[n.ChainID.String()]
		if ok {
			n.Logger.Infof("The node has successfully started. Note: if it's not opted in on %s, "+
				"then please follow the EigenDA operator guide section in https://docs.eigencloud.xyz/products/eigenda/operator-guides/run-a-node/registration to register", eigenDAUrl)
		} else {
			n.Logger.Infof("The node has started but the network with chainID %s is not supported yet",
				n.ChainID.String())
		}
	}

	if operator != nil && operator.Address != "" {
		operatorID, err := n.Transactor.OperatorAddressToID(ctx, gethcommon.HexToAddress(operator.Address))
		if err != nil {
			return fmt.Errorf("failed to get operator ID: %w", err)
		}
		if operatorID != operator.OperatorId {
			return fmt.Errorf("operator ID mismatch: expected %s, got %s",
				operator.OperatorId.Hex(), operatorID.Hex())
		}
	}

	n.CurrentSocket = socket
	// Start the Node IP updater only if the PUBLIC_IP_PROVIDER is greater than 0.
	if n.Config.PubIPCheckInterval > 0 && n.Config.EnableV1 && n.Config.EnableV2 {
		go n.checkRegisteredNodeIpOnChain(ctx)
		go n.checkCurrentNodeIp(ctx)
	}

	return nil
}

// The expireLoop is a loop that is run once per configured second(s) while the node
// is running. It scans for expired batches and removes them from the local database.
func (n *Node) expireLoop() {
	n.Logger.Info("Start expireLoop goroutine in background to periodically remove expired batches on the node")
	ticker := time.NewTicker(time.Duration(n.Config.ExpirationPollIntervalSec) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		// We cap the time the deletion function can run, to make sure there is no overlapping
		// between loops and the garbage collection doesn't take too much resource.
		// The heuristic is to cap the GC time to a percentage of the poll interval, but at
		// least have 1 second.
		timeLimitSec := uint64(math.Max(float64(n.Config.ExpirationPollIntervalSec)*gcPercentageTime, 1.0))
		numBatchesDeleted, numMappingsDeleted, numBlobsDeleted, err := n.Store.DeleteExpiredEntries(
			time.Now().Unix(), timeLimitSec)
		n.Logger.Info("Complete an expiration cycle to remove expired batches",
			"num expired batches found and removed", numBatchesDeleted,
			"num expired mappings found and removed", numMappingsDeleted,
			"num expired blobs found and removed", numBlobsDeleted)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				n.Logger.Error("Expiration cycle exited with ContextDeadlineExceed, meaning more expired "+
					"batches need to be removed, which will continue in next cycle", "time limit (sec)",
					timeLimitSec)
			} else {
				n.Logger.Error("Expiration cycle encountered error when removing expired batches, "+
					"which will be retried in next cycle", "err", err)
			}
		}
	}
}

// RefreshOnchainState refreshes the onchain state of the node.
// It fetches the latest blob parameters from the chain and updates the BlobVersionParams.
// It runs periodically based on the OnchainStateRefreshInterval.
// WARNING: this method is not thread-safe and should not be called concurrently.
func (n *Node) RefreshOnchainState(ctx context.Context) error {
	if !n.Config.EnableV2 || n.Config.OnchainStateRefreshInterval <= 0 {
		return nil
	}
	ticker := time.NewTicker(n.Config.OnchainStateRefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			n.Logger.Info("Refreshing onchain state")
			existingBlobParams := n.BlobVersionParams.Load()
			blobParams, err := n.Transactor.GetAllVersionedBlobParams(ctx)
			if err == nil {
				if existingBlobParams == nil || !existingBlobParams.Equal(blobParams) {
					n.BlobVersionParams.Store(corev2.NewBlobVersionParameterMap(blobParams))
				}
			} else {
				n.Logger.Error("error fetching blob params", "err", err)
			}
			blockNumber, err := n.Transactor.GetCurrentBlockNumber(ctx)
			if err == nil {
				quorumCount, err := n.Transactor.GetQuorumCount(ctx, blockNumber)
				if err == nil {
					n.QuorumCount.Store(uint32(quorumCount))
				} else {
					n.Logger.Error("error fetching quorum count", "err", err)
				}
			} else {
				n.Logger.Error("error fetching block number", "err", err)
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// ProcessBatch validates the batch is correct, stores data into the node's Store, and then returns a
// signature for the entire batch.
//
// The batch will be itemized into batch header, header and chunks of each blob in the batch. These items will
// be stored atomically to the database.
//
// Notes:
//   - If the batch is stored already, it's no-op to store it more than once
//   - If the batch is stored, but the processing fails after that, these data items will not be rollback
//   - These data items will be garbage collected eventually when they become stale.
func (n *Node) ProcessBatch(
	ctx context.Context,
	header *core.BatchHeader,
	blobs []*core.BlobMessage,
	rawBlobs []*node.Blob,
) (*core.Signature, error) {

	start := time.Now()
	log := n.Logger

	batchHeaderHash, err := header.GetBatchHeaderHash()
	if err != nil {
		return nil, err
	}

	if len(blobs) == 0 {
		return nil, errors.New("number of blobs must be greater than zero")
	}

	if len(blobs) != len(rawBlobs) {
		return nil, errors.New("number of parsed blobs must be the same as number of blobs from protobuf request")
	}

	// Measure num batches received and its size in bytes
	batchSize := uint64(0)
	for _, blob := range blobs {
		for quorumID, bundle := range blob.Bundles {
			n.Metrics.AcceptBlobs(quorumID, bundle.Size())
		}
		batchSize += blob.Bundles.Size()
	}
	n.Metrics.AcceptBatches("received", batchSize)

	batchHeaderHashHex := hex.EncodeToString(batchHeaderHash[:])
	log.Debug("Start processing a batch",
		"batchHeaderHash", batchHeaderHashHex,
		"batchSize (in bytes)", batchSize,
		"num of blobs", len(blobs),
		"referenceBlockNumber", header.ReferenceBlockNumber)

	// Store the batch.
	// Run this in a goroutine so we can parallelize the batch storing and batch
	// verifaction work.
	// This should be able to improve latency without needing more CPUs, because batch
	// storing is an IO operation.
	type storeResult struct {
		// Whether StoreBatch failed.
		err error

		// The keys that are stored to database for a single batch.
		// Defined only if the batch not already exists and gets stored to database successfully.
		keys *[][]byte

		// Latency to store the batch.
		// Defined only if the batch not already exists and gets stored to database successfully.
		latency time.Duration
	}
	storeChan := make(chan storeResult)
	go func(n *Node) {
		start := time.Now()
		keys, err := n.Store.StoreBatch(ctx, header, blobs, rawBlobs)
		if err != nil {
			// If batch already exists, we don't store it again, but we should not
			// error out in such case.
			if errors.Is(err, ErrBatchAlreadyExist) {
				storeChan <- storeResult{err: nil, keys: nil, latency: 0}
			} else {
				storeChan <- storeResult{
					err:     fmt.Errorf("failed to store batch: %w", err),
					keys:    nil,
					latency: 0,
				}
			}
			return
		}
		storeChan <- storeResult{err: nil, keys: keys, latency: time.Since(start)}
	}(n)

	// Validate batch.
	stageTimer := time.Now()
	err = n.ValidateBatch(ctx, header, blobs)
	if err != nil {
		// If we have already stored the batch into database, but it's not valid, we
		// revert all the keys for that batch.
		result := <-storeChan
		if result.keys != nil {
			log.Debug("Batch validation failed, rolling back the key/value entries stored in database",
				"number of entries", len(*result.keys),
				"batchHeaderHash", batchHeaderHashHex)
			if deleteKeysErr := n.Store.DeleteKeys(ctx, result.keys); deleteKeysErr != nil {
				log.Error("Failed to delete the invalid batch that should be rolled back",
					"batchHeaderHash", batchHeaderHashHex,
					"err", deleteKeysErr)
			}
		}
		return nil, err
	}
	n.Metrics.RecordStoreChunksStage("validated", batchSize, time.Since(stageTimer))
	log.Debug("Validate batch took", "duration:", time.Since(stageTimer))

	// Before we sign the batch, we should first complete the batch storing successfully.
	result := <-storeChan
	if result.err != nil {
		log.Error("Store batch failed", "batchHeaderHash", batchHeaderHashHex, "err", result.err)
		return nil, err
	}
	if result.keys != nil {
		n.Metrics.RecordStoreChunksStage("stored", batchSize, result.latency)
		n.Logger.Debug("Store batch succeeded",
			"batchHeaderHash", batchHeaderHashHex,
			"duration:", result.latency)
	} else {
		n.Logger.Warn("Store batch skipped because the batch already exists in the store",
			"batchHeaderHash", batchHeaderHashHex)
	}

	// Sign batch header hash if all validation checks pass and data items are written to database.
	stageTimer = time.Now()
	signature, err := n.SignMessage(ctx, batchHeaderHash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign batch: %w", err)
	}

	n.Metrics.RecordStoreChunksStage("signed", batchSize, time.Since(stageTimer))
	log.Debug("Sign batch succeeded",
		"pubkey", n.BLSSigner.GetPublicKeyG1(),
		"duration", time.Since(stageTimer))

	log.Debug("Exiting process batch", "duration", time.Since(start))
	return signature, nil
}

func (n *Node) SignMessage(ctx context.Context, data [32]byte) (*core.Signature, error) {
	signature, err := n.BLSSigner.Sign(ctx, data[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}
	sig := new(core.Signature)
	g, err := sig.Deserialize(signature)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize signature: %w", err)
	}
	return &core.Signature{
		G1Point: g,
	}, nil
}

func (n *Node) ValidateBatch(ctx context.Context, header *core.BatchHeader, blobs []*core.BlobMessage) error {
	start := time.Now()
	operatorState, err := n.ChainState.GetOperatorStateByOperator(ctx, header.ReferenceBlockNumber, n.Config.ID)
	if err != nil {
		return err
	}
	getStateDuration := time.Since(start)

	err = n.Validator.ValidateBatch(header, blobs, operatorState, n.ValidationPool)
	if err != nil {
		h, hashErr := operatorState.Hash()
		if hashErr != nil {
			n.Logger.Error("failed to get operator state hash", "err", hashErr)
		}

		hStr := make([]string, 0, len(h))
		for q, hash := range h {
			hStr = append(hStr, fmt.Sprintf("%d: %x", q, hash))
		}
		return fmt.Errorf("failed to validate batch with operator state %x: %w",
			strings.Join(hStr, ","), err)
	}
	n.Logger.Debug("ValidateBatch completed",
		"get operator state duration", getStateDuration,
		"total duration", time.Since(start))
	return nil
}

func (n *Node) updateSocketAddress(ctx context.Context, newSocketAddr string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if newSocketAddr == n.CurrentSocket {
		return
	}

	if err := n.Transactor.UpdateOperatorSocket(ctx, newSocketAddr); err != nil {
		n.Logger.Error("failed to update operator's socket", err)
		return
	}

	n.Logger.Info("Socket update", "old socket", n.CurrentSocket, "new socket", newSocketAddr)
	n.Metrics.RecordSocketAddressChange()
	n.CurrentSocket = newSocketAddr
}

func (n *Node) checkRegisteredNodeIpOnChain(ctx context.Context) {
	n.Logger.Info("Start checkRegisteredNodeIpOnChain goroutine in background to subscribe the " +
		"operator socket change events onchain")

	socketChan, err := n.OperatorSocketsFilterer.WatchOperatorSocketUpdate(ctx, n.Config.ID)
	if err != nil {
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case socket := <-socketChan:
			n.mu.Lock()
			if socket != n.CurrentSocket {
				n.Logger.Info(
					"Detected socket registered onchain which is different than the socket kept at the DA Node",
					"socket kept at DA Node", n.CurrentSocket,
					"socket registered onchain", socket,
					"the action taken", "update the socket kept at DA Node")
				n.CurrentSocket = socket
			}
			n.mu.Unlock()
		}
	}
}

func (n *Node) checkCurrentNodeIp(ctx context.Context) {
	n.Logger.Info(
		"Start checkCurrentNodeIp goroutine in background to detect the current public IP of the operator node")

	t := time.NewTimer(n.Config.PubIPCheckInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			newSocketAddr, err := SocketAddress(
				ctx,
				n.PubIPProvider,
				n.Config.DispersalPort,
				n.Config.RetrievalPort,
				n.Config.V2DispersalPort,
				n.Config.V2RetrievalPort)
			if err != nil {
				n.Logger.Error("failed to get socket address", "err", err)
				continue
			}
			n.updateSocketAddress(ctx, newSocketAddr)
		}
	}
}

// OperatorReachabilityResponse is the response object for the reachability check
// For v1 endpoints
type OperatorReachabilityResponse struct {
	OperatorID      string `json:"operator_id"`
	DispersalSocket string `json:"dispersal_socket"`
	RetrievalSocket string `json:"retrieval_socket"`
	DispersalOnline bool   `json:"dispersal_online"`
	RetrievalOnline bool   `json:"retrieval_online"`
	DispersalStatus string `json:"dispersal_status"`
	RetrievalStatus string `json:"retrieval_status"`
}

// OperatorV2ReachabilityResponse is the response object for the v2 reachability check
type OperatorV2ReachabilityResponse struct {
	Operators []OperatorReachabilityResponse `json:"operators"`
}

func (n *Node) checkNodeReachability(checkPath string) {
	if n.Config.ReachabilityPollIntervalSec == 0 {
		n.Logger.Warn("Node reachability checks disabled!")
		return
	}

	if n.Config.DataApiUrl == "" {
		n.Logger.Error("Unable to perform reachability check - NODE_DATAAPI_URL is not defined in .env")
		return
	}

	version := "v1"
	if strings.Contains(checkPath, "v2") {
		version = "v2"
	}

	checkURL, err := GetReachabilityURL(n.Config.DataApiUrl, checkPath, n.Config.ID.Hex())
	if err != nil {
		n.Logger.Error("Failed to get reachability check URL", err)
		return
	}

	n.Logger.Info(
		"Start nodeReachabilityCheck goroutine in background to check the reachability of the operator node")
	ticker := time.NewTicker(time.Duration(n.Config.ReachabilityPollIntervalSec) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		n.Logger.Debug(fmt.Sprintf("Calling %s reachability check", version), "url", checkURL)

		resp, err := http.Get(checkURL)
		if err != nil {
			n.Logger.Error(fmt.Sprintf("Reachability check %s - request failed", version), err)
			continue
		} else if resp.StatusCode == 404 {
			body, _ := io.ReadAll(resp.Body)
			if string(body) == "404 page not found" {
				n.Logger.Error("Invalid reachability check url", "checkUrl", checkURL)
			} else {
				n.Logger.Warn("Reachability check operator id not found",
					"status", resp.StatusCode,
					"operator_id", n.Config.ID.Hex())
			}
			continue
		} else if resp.StatusCode != 200 {
			n.Logger.Error(fmt.Sprintf("Reachability check %s - request failed", version),
				"status", resp.StatusCode)
			continue
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			n.Logger.Error(fmt.Sprintf("Failed to read %s reachability check response", version), err)
			continue
		}

		if version == "v1" {
			var responseObject OperatorReachabilityResponse
			err = json.Unmarshal(data, &responseObject)
			if err != nil {
				n.Logger.Error("Reachability check failed to unmarshal json response", err)
				continue
			}

			n.processReachabilityResponse(version, responseObject)
		} else {
			var v2ResponseObject OperatorV2ReachabilityResponse
			err = json.Unmarshal(data, &v2ResponseObject)
			if err != nil {
				n.Logger.Error("Reachability check v2 failed to unmarshal json response", err)
				continue
			}

			if len(v2ResponseObject.Operators) > 0 {
				// Process the first operator from the array
				n.processReachabilityResponse(version, v2ResponseObject.Operators[0])
			} else {
				n.Logger.Error("Reachability check v2 returned empty operators array")
			}
		}
	}
}

// processReachabilityResponse handles the response for a single operator
func (n *Node) processReachabilityResponse(version string, responseObject OperatorReachabilityResponse) {
	if responseObject.DispersalOnline {
		n.Logger.Info(fmt.Sprintf("Reachability check %s - dispersal socket ONLINE", version),
			"status", responseObject.DispersalStatus,
			"socket", responseObject.DispersalSocket)
		n.Metrics.ReachabilityGauge.WithLabelValues(fmt.Sprintf("dispersal-%s", version)).Set(1.0)
	} else {
		n.Logger.Error(fmt.Sprintf("Reachability check %s - dispersal socket UNREACHABLE", version),
			"status", responseObject.DispersalStatus,
			"socket", responseObject.DispersalSocket)
		n.Metrics.ReachabilityGauge.WithLabelValues(fmt.Sprintf("dispersal-%s", version)).Set(0.0)
	}
	if responseObject.RetrievalOnline {
		n.Logger.Info(fmt.Sprintf("Reachability check %s - retrieval socket ONLINE", version),
			"status", responseObject.RetrievalStatus,
			"socket", responseObject.RetrievalSocket)
		n.Metrics.ReachabilityGauge.WithLabelValues(fmt.Sprintf("retrieval-%s", version)).Set(1.0)
	} else {
		n.Logger.Error(fmt.Sprintf("Reachability check %s - retrieval socket UNREACHABLE", version),
			"status", responseObject.RetrievalStatus,
			"socket", responseObject.RetrievalSocket)
		n.Metrics.ReachabilityGauge.WithLabelValues(fmt.Sprintf("retrieval-%s", version)).Set(0.0)
	}
}

func GetReachabilityURL(dataApiUrl, path, operatorID string) (string, error) {
	checkURLString, err := url.JoinPath(dataApiUrl, path)
	if err != nil {
		return "", err
	}
	checkURL, err := url.Parse(checkURLString)
	if err != nil {
		return "", err
	}

	q := checkURL.Query()
	q.Set("operator_id", operatorID)
	checkURL.RawQuery = q.Encode()

	return checkURL.String(), nil
}
