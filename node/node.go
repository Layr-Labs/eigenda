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
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/Layr-Labs/eigenda/common/pubip"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wealdtech/go-merkletree/v2"
	"github.com/wealdtech/go-merkletree/v2/keccak256"
	"google.golang.org/protobuf/proto"

	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/indexer"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/metrics"
	rpccalls "github.com/Layr-Labs/eigensdk-go/metrics/collectors/rpc_calls"
	"github.com/Layr-Labs/eigensdk-go/nodeapi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gammazero/workerpool"
)

const (
	// The percentage of time in garbage collection in a GC cycle.
	gcPercentageTime = 0.1
)

var (
	// eigenDAUIMap is a mapping for ChainID to the EigenDA UI url.
	eigenDAUIMap = map[string]string{
		"17000": "https://holesky.eigenlayer.xyz/avs/eigenda",
		"1":     "https://app.eigenlayer.xyz/avs/eigenda",
	}
)

type Node struct {
	Config                  *Config
	Logger                  logging.Logger
	KeyPair                 *core.KeyPair
	Metrics                 *Metrics
	NodeApi                 *nodeapi.NodeApi
	Store                   *Store
	ChainState              core.ChainState
	Validator               core.ShardValidator
	Transactor              core.Writer
	PubIPProvider           pubip.Provider
	OperatorSocketsFilterer indexer.OperatorSocketsFilterer
	ChainID                 *big.Int

	mu            sync.Mutex
	CurrentSocket string
}

// NewNode creates a new Node with the provided config.
func NewNode(reg *prometheus.Registry, config *Config, pubIPProvider pubip.Provider, logger logging.Logger) (*Node, error) {
	// Setup metrics
	// sdkClients, err := buildSdkClients(config, logger)
	// if err != nil {
	// 	return nil, err
	// }

	eigenMetrics := metrics.NewEigenMetrics(AppName, ":"+config.MetricsPort, reg, logger.With("component", "EigenMetrics"))
	rpcCallsCollector := rpccalls.NewCollector(AppName, reg)

	// Generate BLS keys
	keyPair, err := core.MakeKeyPairFromString(config.PrivateBls)
	if err != nil {
		return nil, err
	}

	config.ID = keyPair.GetPubKeyG1().GetOperatorID()

	// Make sure config folder exists.
	err = os.MkdirAll(config.DbPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("could not create db directory at %s: %w", config.DbPath, err)
	}

	client, err := geth.NewInstrumentedEthClient(config.EthClientConfig, rpcCallsCollector, logger)
	if err != nil {
		return nil, fmt.Errorf("cannot create chain.Client: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chainID: %w", err)
	}

	// Create Transactor
	tx, err := eth.NewWriter(logger, client, config.BLSOperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return nil, err
	}

	// Create ChainState Client
	cst := eth.NewChainState(tx, client)

	// Setup Node Api
	nodeApi := nodeapi.NewNodeApi(AppName, SemVer, ":"+config.NodeApiPort, logger.With("component", "NodeApi"))

	metrics := NewMetrics(eigenMetrics, reg, logger, ":"+config.MetricsPort, config.ID, config.OnchainMetricsInterval, tx, cst)

	// Make validator
	v, err := verifier.NewVerifier(&config.EncoderConfig, false)
	if err != nil {
		return nil, err
	}
	asgn := &core.StdAssignmentCoordinator{}
	validator := core.NewShardValidator(v, asgn, cst, config.ID)

	// Resolve the BLOCK_STALE_MEASURE and STORE_DURATION_BLOCKS.
	var blockStaleMeasure, storeDurationBlocks uint32
	if config.EnableTestMode && config.OverrideBlockStaleMeasure > 0 {
		blockStaleMeasure = uint32(config.OverrideBlockStaleMeasure)
	} else {
		staleMeasure, err := tx.GetBlockStaleMeasure(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get BLOCK_STALE_MEASURE: %w", err)
		}
		blockStaleMeasure = staleMeasure
	}
	if config.EnableTestMode && config.OverrideStoreDurationBlocks > 0 {
		storeDurationBlocks = uint32(config.OverrideStoreDurationBlocks)
	} else {
		storeDuration, err := tx.GetStoreDurationBlocks(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to get STORE_DURATION_BLOCKS: %w", err)
		}
		storeDurationBlocks = storeDuration
	}
	// Create new store
	store, err := NewLevelDBStore(config.DbPath+"/chunk", logger, metrics, blockStaleMeasure, storeDurationBlocks)
	if err != nil {
		return nil, fmt.Errorf("failed to create new store: %w", err)
	}

	eigenDAServiceManagerAddr := gethcommon.HexToAddress(config.EigenDAServiceManagerAddr)
	socketsFilterer, err := indexer.NewOperatorSocketsFilterer(eigenDAServiceManagerAddr, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create new operator sockets filterer: %w", err)
	}
	nodeLogger := logger.With("component", "Node")
	nodeLogger.Info("Creating node", "chainID", chainID.String(), "operatorID", config.ID.Hex(),
		"dispersalPort", config.DispersalPort, "retrievalPort", config.RetrievalPort, "churnerUrl", config.ChurnerUrl,
		"quorumIDs", fmt.Sprint(config.QuorumIDList), "registerNodeAtStart", config.RegisterNodeAtStart, "pubIPCheckInterval", config.PubIPCheckInterval,
		"eigenDAServiceManagerAddr", config.EigenDAServiceManagerAddr, "blockStaleMeasure", blockStaleMeasure, "storeDurationBlocks", storeDurationBlocks, "enableGnarkBundleEncoding", config.EnableGnarkBundleEncoding)

	return &Node{
		Config:                  config,
		Logger:                  nodeLogger,
		KeyPair:                 keyPair,
		Metrics:                 metrics,
		NodeApi:                 nodeApi,
		Store:                   store,
		ChainState:              cst,
		Transactor:              tx,
		Validator:               validator,
		PubIPProvider:           pubIPProvider,
		OperatorSocketsFilterer: socketsFilterer,
		ChainID:                 chainID,
	}, nil
}

// Starts the Node. If the node is not registered, register it on chain, otherwise just
// update its socket on chain.
func (n *Node) Start(ctx context.Context) error {
	if n.Config.EnableMetrics {
		n.Metrics.Start()
		n.Logger.Info("Enabled metrics", "socket", n.Metrics.socketAddr)
	}
	if n.Config.EnableNodeApi {
		n.NodeApi.Start()
		n.Logger.Info("Enabled node api", "port", n.Config.NodeApiPort)
	}

	go n.expireLoop()
	go n.checkNodeReachability()

	// Build the socket based on the hostname/IP provided in the CLI
	socket := string(core.MakeOperatorSocket(n.Config.Hostname, n.Config.DispersalPort, n.Config.RetrievalPort))
	var operator *Operator
	if n.Config.RegisterNodeAtStart {
		n.Logger.Info("Registering node on chain with the following parameters:", "operatorId",
			n.Config.ID.Hex(), "hostname", n.Config.Hostname, "dispersalPort", n.Config.DispersalPort,
			"retrievalPort", n.Config.RetrievalPort, "churnerUrl", n.Config.ChurnerUrl, "quorumIds", fmt.Sprint(n.Config.QuorumIDList))
		socket := string(core.MakeOperatorSocket(n.Config.Hostname, n.Config.DispersalPort, n.Config.RetrievalPort))
		privateKey, err := crypto.HexToECDSA(n.Config.EthClientConfig.PrivateKeyString)
		if err != nil {
			return fmt.Errorf("NewClient: cannot parse private key: %w", err)
		}
		operator = &Operator{
			Address:             crypto.PubkeyToAddress(privateKey.PublicKey).Hex(),
			Socket:              socket,
			Timeout:             10 * time.Second,
			PrivKey:             privateKey,
			KeyPair:             n.KeyPair,
			OperatorId:          n.Config.ID,
			QuorumIDs:           n.Config.QuorumIDList,
			RegisterNodeAtStart: n.Config.RegisterNodeAtStart,
		}
		churnerClient := NewChurnerClient(n.Config.ChurnerUrl, n.Config.UseSecureGrpc, n.Config.Timeout, n.Logger)
		err = RegisterOperator(ctx, operator, n.Transactor, churnerClient, n.Logger)
		if err != nil {
			return fmt.Errorf("failed to register the operator: %w", err)
		}
	} else {
		eigenDAUrl, ok := eigenDAUIMap[n.ChainID.String()]
		if ok {
			n.Logger.Infof("The node has successfully started. Note: if it's not opted in on %s, then please follow the EigenDA operator guide section in docs.eigenlayer.xyz to register", eigenDAUrl)
		} else {
			n.Logger.Infof("The node has started but the network with chainID %s is not supported yet", n.ChainID.String())
		}
	}

	if operator != nil && operator.Address != "" {
		operatorID, err := n.Transactor.OperatorAddressToID(ctx, gethcommon.HexToAddress(operator.Address))
		if err != nil {
			return fmt.Errorf("failed to get operator ID: %w", err)
		}
		if operatorID != operator.OperatorId {
			return fmt.Errorf("operator ID mismatch: expected %s, got %s", operator.OperatorId.Hex(), operatorID.Hex())
		}
	}

	n.CurrentSocket = socket
	// Start the Node IP updater only if the PUBLIC_IP_PROVIDER is greater than 0.
	if n.Config.PubIPCheckInterval > 0 {
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
		numBatchesDeleted, numMappingsDeleted, numBlobsDeleted, err := n.Store.DeleteExpiredEntries(time.Now().Unix(), timeLimitSec)
		n.Logger.Info("Complete an expiration cycle to remove expired batches", "num expired batches found and removed", numBatchesDeleted, "num expired mappings found and removed", numMappingsDeleted, "num expired blobs found and removed", numBlobsDeleted)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				n.Logger.Error("Expiration cycle exited with ContextDeadlineExceed, meaning more expired batches need to be removed, which will continue in next cycle", "time limit (sec)", timeLimitSec)
			} else {
				n.Logger.Error("Expiration cycle encountered error when removing expired batches, which will be retried in next cycle", "err", err)
			}
		}
	}
}

// ProcessBatch validates the batch is correct, stores data into the node's Store, and then returns a signature for the entire batch.
//
// The batch will be itemized into batch header, header and chunks of each blob in the batch. These items will
// be stored atomically to the database.
//
// Notes:
//   - If the batch is stored already, it's no-op to store it more than once
//   - If the batch is stored, but the processing fails after that, these data items will not be rollback
//   - These data items will be garbage collected eventually when they become stale.
func (n *Node) ProcessBatch(ctx context.Context, header *core.BatchHeader, blobs []*core.BlobMessage, rawBlobs []*node.Blob) (*core.Signature, error) {
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
	log.Debug("Start processing a batch", "batchHeaderHash", batchHeaderHashHex, "batchSize (in bytes)", batchSize, "num of blobs", len(blobs), "referenceBlockNumber", header.ReferenceBlockNumber)

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
				storeChan <- storeResult{err: fmt.Errorf("failed to store batch: %w", err), keys: nil, latency: 0}
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
			log.Debug("Batch validation failed, rolling back the key/value entries stored in database", "number of entires", len(*result.keys), "batchHeaderHash", batchHeaderHashHex)
			if deleteKeysErr := n.Store.DeleteKeys(ctx, result.keys); deleteKeysErr != nil {
				log.Error("Failed to delete the invalid batch that should be rolled back", "batchHeaderHash", batchHeaderHashHex, "err", deleteKeysErr)
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
		n.Logger.Debug("Store batch succeeded", "batchHeaderHash", batchHeaderHashHex, "duration:", result.latency)
	} else {
		n.Logger.Warn("Store batch skipped because the batch already exists in the store", "batchHeaderHash", batchHeaderHashHex)
	}

	// Sign batch header hash if all validation checks pass and data items are written to database.
	stageTimer = time.Now()
	sig := n.KeyPair.SignMessage(batchHeaderHash)
	n.Metrics.RecordStoreChunksStage("signed", batchSize, time.Since(stageTimer))
	log.Debug("Sign batch succeeded", "pubkey", hexutil.Encode(n.KeyPair.GetPubKeyG2().Serialize()), "duration", time.Since(stageTimer))

	log.Debug("Exiting process batch", "duration", time.Since(start))
	return sig, nil
}

// ProcessBlobs validates the blobs are correct, stores data into the node's Store, and then returns a signature for each blob.
// This method is similar to ProcessBatch method except that it doesn't require a batch.
//
// Notes:
//   - If the blob is stored already, it's no-op to store it more than once
//   - If the blob is stored, but the processing fails after that, these data items will not be rollback
//   - These data items will be garbage collected eventually when they become stale.
func (n *Node) ProcessBlobs(ctx context.Context, blobs []*core.BlobMessage, rawBlobs []*node.Blob) ([]*core.Signature, error) {
	start := time.Now()
	log := n.Logger

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

	log.Debug("Start processing blobs", "batchSizeBytes", batchSize, "numBlobs", len(blobs))

	// Store the blobs
	// Run this in a goroutine so we can parallelize the blob storing and blob verifaction work.
	// This should be able to improve latency without needing more CPUs, because blob
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
		keys, err := n.Store.StoreBlobs(ctx, blobs, rawBlobs)
		if err != nil {
			// If batch already exists, we don't store it again, but we should not
			// error out in such case.
			if errors.Is(err, ErrBatchAlreadyExist) {
				storeChan <- storeResult{err: nil, keys: nil, latency: 0}
			} else {
				storeChan <- storeResult{err: fmt.Errorf("failed to store batch: %w", err), keys: nil, latency: 0}
			}
			return
		}
		storeChan <- storeResult{err: nil, keys: keys, latency: time.Since(start)}
	}(n)

	// Validate batch
	// Assumes that all blobs have been validated to have the same reference block number.
	stageTimer := time.Now()
	referenceBlockNumber := uint(0)
	for _, blob := range rawBlobs {
		blobRefBlock := blob.GetHeader().GetReferenceBlockNumber()
		if referenceBlockNumber == 0 && blobRefBlock != 0 {
			referenceBlockNumber = uint(blobRefBlock)
			break
		}
	}
	if referenceBlockNumber == 0 {
		return nil, errors.New("reference block number is not set")
	}

	err := n.ValidateBlobs(ctx, blobs, referenceBlockNumber)
	if err != nil {
		// If we have already stored the batch into database, but it's not valid, we
		// revert all the keys for that batch.
		result := <-storeChan
		if result.keys != nil {
			log.Debug("Batch validation failed, rolling back the key/value entries stored in database", "number of entires", len(*result.keys), "referenceBlockNumber", referenceBlockNumber)
			if deleteKeysErr := n.Store.DeleteKeys(ctx, result.keys); deleteKeysErr != nil {
				log.Error("Failed to delete the invalid batch that should be rolled back", "err", deleteKeysErr)
			}
		}
		return nil, err
	}
	n.Metrics.RecordStoreChunksStage("validated", batchSize, time.Since(stageTimer))
	log.Debug("Validate blobs took", "duration:", time.Since(stageTimer))

	// Before we sign the blobs, we should first complete the batch storing successfully.
	result := <-storeChan
	if result.err != nil {
		return nil, fmt.Errorf("failed to store blobs: %w", result.err)
	}
	if result.keys != nil {
		n.Metrics.RecordStoreChunksStage("stored", batchSize, result.latency)
		n.Logger.Debug("StoreBlobs succeeded", "duration:", result.latency)
	} else {
		n.Logger.Warn("StoreBlobs skipped because the batch already exists in the store")
	}

	// Sign all blobs if all validation checks pass and data items are written to database.
	stageTimer = time.Now()
	signatures, err := n.SignBlobs(blobs, referenceBlockNumber)
	if err != nil {
		return nil, fmt.Errorf("failed to sign blobs: %w", err)
	}
	n.Metrics.RecordStoreChunksStage("signed", batchSize, time.Since(stageTimer))
	log.Debug("SignBlobs succeeded", "pubkey", hexutil.Encode(n.KeyPair.GetPubKeyG2().Serialize()), "duration", time.Since(stageTimer))

	log.Debug("Exiting ProcessBlobs", "duration", time.Since(start))
	return signatures, nil
}

func (n *Node) ValidateBatch(ctx context.Context, header *core.BatchHeader, blobs []*core.BlobMessage) error {
	start := time.Now()
	operatorState, err := n.ChainState.GetOperatorStateByOperator(ctx, header.ReferenceBlockNumber, n.Config.ID)
	if err != nil {
		return err
	}
	getStateDuration := time.Since(start)

	pool := workerpool.New(n.Config.NumBatchValidators)
	err = n.Validator.ValidateBatch(header, blobs, operatorState, pool)
	if err != nil {
		h, hashErr := operatorState.Hash()
		if hashErr != nil {
			n.Logger.Error("failed to get operator state hash", "err", hashErr)
		}

		hStr := make([]string, 0, len(h))
		for q, hash := range h {
			hStr = append(hStr, fmt.Sprintf("%d: %x", q, hash))
		}
		return fmt.Errorf("failed to validate batch with operator state %x: %w", strings.Join(hStr, ","), err)
	}
	n.Logger.Debug("ValidateBatch completed", "get operator state duration", getStateDuration, "total duration", time.Since(start))
	return nil
}

// ValidateBlobs validates the blob commitments are correct
func (n *Node) ValidateBlobs(ctx context.Context, blobs []*core.BlobMessage, referenceBlockNumber uint) error {
	start := time.Now()
	operatorState, err := n.ChainState.GetOperatorStateByOperator(ctx, referenceBlockNumber, n.Config.ID)
	if err != nil {
		return err
	}
	getStateDuration := time.Since(start)

	pool := workerpool.New(n.Config.NumBatchValidators)
	err = n.Validator.ValidateBlobs(blobs, operatorState, pool)
	if err != nil {
		h, hashErr := operatorState.Hash()
		if hashErr != nil {
			n.Logger.Error("failed to get operator state hash", "err", hashErr)
		}

		hStr := make([]string, 0, len(h))
		for q, hash := range h {
			hStr = append(hStr, fmt.Sprintf("%d: %x", q, hash))
		}
		return fmt.Errorf("failed to validate batch with operator state %x: %w", strings.Join(hStr, ","), err)
	}
	n.Logger.Debug("ValidateBlob completed", "get operator state duration", getStateDuration, "total duration", time.Since(start))
	return nil
}

func (n *Node) SignBlobs(blobs []*core.BlobMessage, referenceBlockNumber uint) ([]*core.Signature, error) {
	start := time.Now()
	signatures := make([]*core.Signature, len(blobs))
	for i, blob := range blobs {
		if blob == nil || blob.BlobHeader == nil {
			signatures[i] = nil
			continue
		}
		batchHeader := &core.BatchHeader{
			ReferenceBlockNumber: referenceBlockNumber,
			BatchRoot:            [32]byte{},
		}
		_, err := batchHeader.SetBatchRoot([]*core.BlobHeader{blob.BlobHeader})
		if err != nil {
			return nil, fmt.Errorf("failed to set batch root: %w", err)
		}
		batchHeaderHash, err := batchHeader.GetBatchHeaderHash()
		if err != nil {
			return nil, fmt.Errorf("failed to get batch header hash: %w", err)
		}
		sig := n.KeyPair.SignMessage(batchHeaderHash)
		signatures[i] = sig
	}

	n.Logger.Debug("SignBlobs completed", "duration", time.Since(start))
	return signatures, nil
}

// ValidateBlobHeadersRoot validates the blob headers root hash
// by comparing it with the merkle tree root hash of the blob headers.
// It also checks if all blob headers have the same reference block number
func (n *Node) ValidateBatchContents(ctx context.Context, blobHeaderHashes [][32]byte, batchHeader *core.BatchHeader) error {
	leafs := make([][]byte, 0)
	for _, blobHeaderHash := range blobHeaderHashes {
		blobHeaderBytes, err := n.Store.GetBlobHeaderByHeaderHash(ctx, blobHeaderHash)
		if err != nil {
			return fmt.Errorf("failed to get blob header by hash: %w", err)
		}
		if blobHeaderBytes == nil {
			return fmt.Errorf("blob header not found for hash %x", blobHeaderHash)
		}

		var protoBlobHeader node.BlobHeader
		err = proto.Unmarshal(blobHeaderBytes, &protoBlobHeader)
		if err != nil {
			return fmt.Errorf("failed to unmarshal blob header: %w", err)
		}
		if uint32(batchHeader.ReferenceBlockNumber) != protoBlobHeader.GetReferenceBlockNumber() {
			return errors.New("blob headers have different reference block numbers")
		}

		blobHeader, err := GetBlobHeaderFromProto(&protoBlobHeader)
		if err != nil {
			return fmt.Errorf("failed to get blob header from proto: %w", err)
		}

		blobHeaderHash, err := blobHeader.GetBlobHeaderHash()
		if err != nil {
			return fmt.Errorf("failed to get blob header hash: %w", err)
		}
		leafs = append(leafs, blobHeaderHash[:])
	}

	if len(leafs) == 0 {
		return errors.New("no blob headers found")
	}

	tree, err := merkletree.NewTree(merkletree.WithData(leafs), merkletree.WithHashType(keccak256.New()))
	if err != nil {
		return fmt.Errorf("failed to create merkle tree: %w", err)
	}

	if !reflect.DeepEqual(tree.Root(), batchHeader.BatchRoot[:]) {
		return errors.New("invalid batch header")
	}

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
	n.Logger.Info("Start checkRegisteredNodeIpOnChain goroutine in background to subscribe the operator socket change events onchain")

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
				n.Logger.Info("Detected socket registered onchain which is different than the socket kept at the DA Node", "socket kept at DA Node", n.CurrentSocket, "socket registered onchain", socket, "the action taken", "update the socket kept at DA Node")
				n.CurrentSocket = socket
			}
			n.mu.Unlock()
		}
	}
}

func (n *Node) checkCurrentNodeIp(ctx context.Context) {
	n.Logger.Info("Start checkCurrentNodeIp goroutine in background to detect the current public IP of the operator node")

	t := time.NewTimer(n.Config.PubIPCheckInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			newSocketAddr, err := SocketAddress(ctx, n.PubIPProvider, n.Config.DispersalPort, n.Config.RetrievalPort)
			if err != nil {
				n.Logger.Error("failed to get socket address", "err", err)
				continue
			}
			n.updateSocketAddress(ctx, newSocketAddr)
		}
	}
}

// OperatorReachabilityResponse is the response object for the reachability check
type OperatorReachabilityResponse struct {
	OperatorID      string `json:"operator_id"`
	DispersalSocket string `json:"dispersal_socket"`
	RetrievalSocket string `json:"retrieval_socket"`
	DispersalOnline bool   `json:"dispersal_online"`
	RetrievalOnline bool   `json:"retrieval_online"`
}

func (n *Node) checkNodeReachability() {
	if n.Config.ReachabilityPollIntervalSec == 0 {
		n.Logger.Warn("Node reachability checks disabled!!! ReachabilityPollIntervalSec set to 0")
		return
	}

	if n.Config.DataApiUrl == "" {
		n.Logger.Error("Unable to perform reachability check - NODE_DATAAPI_URL is not defined in .env")
		return
	}

	checkURL, err := GetReachabilityURL(n.Config.DataApiUrl, n.Config.ID.Hex())
	if err != nil {
		n.Logger.Error("Failed to get reachability check URL", err)
		return
	}

	n.Logger.Info("Start nodeReachabilityCheck goroutine in background to check the reachability of the operator node")
	ticker := time.NewTicker(time.Duration(n.Config.ReachabilityPollIntervalSec) * time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		n.Logger.Debug("Calling reachability check", "url", checkURL)

		resp, err := http.Get(checkURL)
		if err != nil {
			n.Logger.Error("Reachability check request failed", err)
			continue
		} else if resp.StatusCode == 404 {
			body, _ := io.ReadAll(resp.Body)
			if string(body) == "404 page not found" {
				n.Logger.Error("Invalid reachability check url", "checkUrl", checkURL)
			} else {
				n.Logger.Warn("Reachability check operator id not found", "status", resp.StatusCode, "operator_id", n.Config.ID.Hex())
			}
			continue
		} else if resp.StatusCode != 200 {
			n.Logger.Error("Reachability check request failed", "status", resp.StatusCode)
			continue
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			n.Logger.Error("Failed to read reachability check response", err)
			continue
		}

		var responseObject OperatorReachabilityResponse
		err = json.Unmarshal(data, &responseObject)
		if err != nil {
			n.Logger.Error("Reachability check failed to unmarshal json response", err)
			continue
		}

		if responseObject.DispersalOnline {
			n.Logger.Info("Reachability check - dispersal socket is ONLINE", "socket", responseObject.DispersalSocket)
			n.Metrics.ReachabilityGauge.WithLabelValues("dispersal").Set(1.0)
		} else {
			n.Logger.Error("Reachability check - dispersal socket is UNREACHABLE", "socket", responseObject.DispersalSocket)
			n.Metrics.ReachabilityGauge.WithLabelValues("dispersal").Set(0.0)
		}
		if responseObject.RetrievalOnline {
			n.Logger.Info("Reachability check - retrieval socket is ONLINE", "socket", responseObject.RetrievalSocket)
			n.Metrics.ReachabilityGauge.WithLabelValues("retrieval").Set(1.0)
		} else {
			n.Logger.Error("Reachability check - retrieval socket is UNREACHABLE", "socket", responseObject.RetrievalSocket)
			n.Metrics.ReachabilityGauge.WithLabelValues("retrieval").Set(0.0)
		}
	}
}

func GetReachabilityURL(dataApiUrl, operatorID string) (string, error) {
	checkURLString, err := url.JoinPath(dataApiUrl, "/api/v1/operators-info/port-check")
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
