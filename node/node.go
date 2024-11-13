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
	"time"

	"github.com/Layr-Labs/eigenda/common/pubip"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"

	"github.com/prometheus/client_golang/prometheus"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/grpc/node"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/indexer"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/Layr-Labs/eigensdk-go/metrics"
	rpccalls "github.com/Layr-Labs/eigensdk-go/metrics/collectors/rpc_calls"
	"github.com/Layr-Labs/eigensdk-go/nodeapi"
	"github.com/gammazero/workerpool"

	blssignerV1 "github.com/Layr-Labs/cerberus-api/pkg/api/v1"
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
	StoreV2                 StoreV2
	ChainState              core.ChainState
	Validator               core.ShardValidator
	ValidatorV2             corev2.ShardValidator
	Transactor              core.Writer
	PubIPProvider           pubip.Provider
	OperatorSocketsFilterer indexer.OperatorSocketsFilterer
	ChainID                 *big.Int
	BLSSigner               blssignerV1.SignerClient

	RelayClient clients.RelayClient

	mu            sync.Mutex
	CurrentSocket string
}

// NewNode creates a new Node with the provided config.
func NewNode(
	reg *prometheus.Registry,
	config *Config,
	pubIPProvider pubip.Provider,
	logger logging.Logger,
) (*Node, error) {
	// Setup metrics
	// sdkClients, err := buildSdkClients(config, logger)
	// if err != nil {
	// 	return nil, err
	// }

	nodeLogger := logger.With("component", "Node")

	eigenMetrics := metrics.NewEigenMetrics(AppName, ":"+config.MetricsPort, reg, logger.With("component", "EigenMetrics"))
	rpcCallsCollector := rpccalls.NewCollector(AppName, reg)

	// Make sure config folder exists.
	err := os.MkdirAll(config.DbPath, os.ModePerm)
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

	var keyPair *core.KeyPair
	var blsClient blssignerV1.SignerClient
	if config.PrivateBls != "" {
		nodeLogger.Info("using local keystore private key for BLS signing")
		// Generate BLS keys
		keyPair, err = core.MakeKeyPairFromString(config.PrivateBls)
		if err != nil {
			return nil, err
		}

		config.ID = keyPair.GetPubKeyG1().GetOperatorID()
	} else {
		pkBytes, err := hex.DecodeString(config.BLSPublicKeyHex)
		if err != nil {
			return nil, fmt.Errorf("failed to decode BLS public key: %w", err)
		}
		pubkey := new(core.G1Point)
		publicKey, err := pubkey.Deserialize(pkBytes)
		if err != nil {
			return nil, err
		}

		config.ID = publicKey.GetOperatorID()

		nodeLogger.Info("creating signer client", "url", config.BLSRemoteSignerUrl)
		creds := insecure.NewCredentials()
		if config.BLSSignerTLSCertFilePath != "" {
			creds, err = credentials.NewClientTLSFromFile(config.BLSSignerTLSCertFilePath, "")
			if err != nil {
				return nil, err
			}
		}
		conn, err := grpc.NewClient(
			config.BLSRemoteSignerUrl, grpc.WithTransportCredentials(creds),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create new BLS remote signer client: %w", err)
		}
		blsClient = blssignerV1.NewSignerClient(conn)
	}

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
	validatorV2 := corev2.NewShardValidator(v, config.ID)

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

	nodeLogger.Info("Creating node", "chainID", chainID.String(), "operatorID", config.ID.Hex(),
		"dispersalPort", config.DispersalPort, "retrievalPort", config.RetrievalPort, "churnerUrl", config.ChurnerUrl,
		"quorumIDs", fmt.Sprint(config.QuorumIDList), "registerNodeAtStart", config.RegisterNodeAtStart, "pubIPCheckInterval", config.PubIPCheckInterval,
		"eigenDAServiceManagerAddr", config.EigenDAServiceManagerAddr, "blockStaleMeasure", blockStaleMeasure, "storeDurationBlocks", storeDurationBlocks, "enableGnarkBundleEncoding", config.EnableGnarkBundleEncoding)

	var relayClient clients.RelayClient
	// Create a new relay client with relay addresses onchain
	return &Node{
		Config:                  config,
		Logger:                  nodeLogger,
		KeyPair:                 keyPair,
		Metrics:                 metrics,
		NodeApi:                 nodeApi,
		Store:                   store,
		StoreV2:                 nil,
		ChainState:              cst,
		Transactor:              tx,
		Validator:               validator,
		ValidatorV2:             validatorV2,
		PubIPProvider:           pubIPProvider,
		OperatorSocketsFilterer: socketsFilterer,
		ChainID:                 chainID,
		RelayClient:             relayClient,
		BLSSigner:               blsClient,
	}, nil
}

// Start starts the Node. If the node is not registered, register it on chain, otherwise just
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
	signature, err := n.SignMessage(ctx, batchHeaderHash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign batch: %w", err)
	}

	n.Metrics.RecordStoreChunksStage("signed", batchSize, time.Since(stageTimer))
	log.Debug("Sign batch succeeded", "pubkey", n.Config.BLSPublicKeyHex, "duration", time.Since(stageTimer))

	log.Debug("Exiting process batch", "duration", time.Since(start))
	return signature, nil
}

func (n *Node) SignMessage(ctx context.Context, data [32]byte) (*core.Signature, error) {
	if n.Config.BLSRemoteSignerEnabled {
		sigResp, err := n.BLSSigner.SignGeneric(
			ctx,
			&blssignerV1.SignGenericRequest{
				PublicKey: n.Config.BLSPublicKeyHex,
				Password:  n.Config.BLSKeyPassword,
				Data:      data[:],
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to sign data: %w", err)
		}
		sig := new(core.Signature)
		g, err := sig.Deserialize(sigResp.Signature)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize signature: %w", err)
		}
		return &core.Signature{
			G1Point: g,
		}, nil
	}
	return n.KeyPair.SignMessage(data), nil
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
