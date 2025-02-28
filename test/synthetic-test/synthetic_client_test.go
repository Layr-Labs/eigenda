//go:build k8s_func_test
// +build k8s_func_test

package k8s_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/google/uuid"
	"github.com/shurcooL/graphql"

	"github.com/Layr-Labs/eigenda/api/clients"
	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	retriever_rpc "github.com/Layr-Labs/eigenda/api/grpc/retriever"
	common "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	rollupbindings "github.com/Layr-Labs/eigenda/contracts/bindings/MockRollup"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type ClientType string

const (
	Disperser ClientType = "Disperser"
	Retriever ClientType = "Retriever"
)

type GrpcClient struct {
	Hostname string
	GrpcPort string
	Timeout  time.Duration
	Client   interface{} // This can be a specific client type (disperser_rpc.DisperserClient, retriever_rpc.RetrieverClient, etc.)
}

type RetrieverClientConfig struct {
	Bls_Operator_State_Retriever     string
	EigenDA_ServiceManager_Retriever string
	ChurnerGraphUrl                  string
	RetrieverSrsOrder                string
	RetrieverG1Path                  string
	RetrieverG2Path                  string
	RetrieverCachePath               string
	RetrieverG2PointPowerOf2Path     string
}

type TestClients struct {
	Clients map[ClientType]*GrpcClient
}

// TestSuite Struct.
type SyntheticTestSuite struct {
	Clients             *TestClients
	EthClient           common.EthClient
	MockRollUp          *rollupbindings.ContractMockRollup
	Logger              *log.Logger
	RetrievalClient     clients.RetrievalClient
	TestRunID           string
	BatcherPullInterval string
}

var (
	testSuite                  *SyntheticTestSuite
	isRetrieverClientEnabled   bool = false
	validateOnchainTransaction bool = false
	retrievalClient            clients.RetrievalClient
	logger                     logging.Logger
)

func setUpClients(pk string, rpcUrl string, mockRollUpContractAddress string, retrieverClientConfig RetrieverClientConfig, batcherPullInterval string) *SyntheticTestSuite {
	testRunID := uuid.New().String()
	logger := log.New(os.Stdout, "EigenDA SyntheticClient:"+testRunID+" ", log.Ldate|log.Ltime)

	// Initialize your TestClients and other suite fields here
	clients, err := newTestClients(map[ClientType]*GrpcClient{
		Disperser: {
			Hostname: "disperser.disperser.svc.cluster.local",
			GrpcPort: "32001",
			Timeout:  10 * time.Second,
		},
		Retriever: {
			Hostname: "retriever.retriever.svc.cluster.local",
			GrpcPort: "retriever-port",
			Timeout:  10 * time.Second,
		},
	})
	if err != nil {
		logger.Printf("Error initializing clients: %v", err)
	}

	loggerConfig := common.DefaultLoggerConfig()
	ethLogger, err := common.NewLogger(loggerConfig)
	if err != nil {
		logger.Printf("Error: %v", err)
		return nil
	}

	pk = strings.TrimPrefix(pk, "0X")
	pk = strings.TrimPrefix(pk, "0x")

	ethConfig := geth.EthClientConfig{
		RPCURLs:          []string{rpcUrl},
		PrivateKeyString: pk,
		NumConfirmations: 0,
	}
	ethClient, err := geth.NewClient(ethConfig, gcommon.Address{}, 0, ethLogger)
	if err != nil {
		log.Printf("Error: failed to create eth client: %v", err)
	}

	var mockRollup *rollupbindings.ContractMockRollup
	if validateOnchainTransaction {
		log.Printf("Create instance of MockRollUp with Contract Address %s", mockRollUpContractAddress)
		mockRollup, err = rollupbindings.NewContractMockRollup(gcommon.HexToAddress(mockRollUpContractAddress), ethClient)
		if err != nil {
			logger.Printf("Error: %v", err)
			return nil
		}
	}

	if isRetrieverClientEnabled {
		err = setupRetrievalClient(ethClient, &retrieverClientConfig, ethLogger)

		if err != nil {
			logger.Printf("Error: %v", err)
			return nil
		}
	}

	// Assign client connections to pointers in TestClients struct
	return &SyntheticTestSuite{
		Clients:             clients,
		EthClient:           ethClient,
		MockRollUp:          mockRollup,
		RetrievalClient:     retrievalClient,
		Logger:              logger,
		TestRunID:           testRunID,
		BatcherPullInterval: batcherPullInterval,
	}
}

func TestMain(m *testing.M) {

	// These are set in deployment/cronjob yaml
	privateKey := os.Getenv("ETHCLIENT_PRIVATE_KEY")
	rpcUrl := os.Getenv("ETHCLIENT_RPC_URL")
	mockRollUpContractAddress := os.Getenv("MOCKROLLUP_CONTRACT_ADDRESS")
	isRetrieverClientEnabled = os.Getenv("RETRIEVER_CLIENT_ENABLE") == strings.ToLower("true")
	validateOnchainTransaction = os.Getenv("VALIDATE_ONCHAIN_TRANSACTION") == strings.ToLower("true")
	blsOperatorStateRetriever := os.Getenv("BLS_OPERATOR_STATE_RETRIEVER")
	eigenDAServiceManagerRetreiever := os.Getenv("EIGENDA_SERVICE_MANAGER_RETRIEVER")
	churnerGraphUrl := os.Getenv("CHURNER_GRAPH_URL")
	retrieverSrsOrder := os.Getenv("RETRIEVER_SRS_ORDER")
	retrieverG1Path := os.Getenv("RETRIEVER_G1_PATH")
	retrieverG2Path := os.Getenv("RETRIEVER_G2_PATH")
	retrieverG2PoinPowerOf2Path := os.Getenv("RETRIEVER_G2_POINT_POWER_OF_2_PATH")
	retrieverCachePath := os.Getenv("RETRIEVER_CACHE_PATH")
	batcherPullInterval := os.Getenv("BATCHER_PULL_INTERVAL")

	// Retriever Config
	retrieverClientConfig := &RetrieverClientConfig{
		Bls_Operator_State_Retriever:     blsOperatorStateRetriever,
		EigenDA_ServiceManager_Retriever: eigenDAServiceManagerRetreiever,
		ChurnerGraphUrl:                  churnerGraphUrl,
		RetrieverSrsOrder:                retrieverSrsOrder,
		RetrieverG1Path:                  retrieverG1Path,
		RetrieverG2Path:                  retrieverG2Path,
		RetrieverG2PointPowerOf2Path:     retrieverG2PoinPowerOf2Path,
		RetrieverCachePath:               retrieverCachePath,
	}

	// Initialize Clients
	testSuite = setUpClients(privateKey, rpcUrl, mockRollUpContractAddress, *retrieverClientConfig, batcherPullInterval)
	logger := testSuite.Logger

	// Check if testSuite is nil
	if testSuite == nil {
		logger.Println("Error: TestSuite initialization failed.")
		os.Exit(1) // Exit the program with a non-zero status code indicating failure
	}
	logger.Println("RPC_URL for Chain...", rpcUrl)
	logger.Println("Mock RollUp Contract Address...", mockRollUpContractAddress)
	logger.Println("Retriever Client Enabled...", isRetrieverClientEnabled)

	logger.Println("Running Test Client...")
	// Run the tests and get the exit code
	exitCode := m.Run()
	logger.Printf("Exiting Test Client Run with Code:%d", exitCode)
	// Exit with the test result code
	os.Exit(exitCode)
}

// SetUp RetrievalClient to retriever blob from Operator Node
func setupRetrievalClient(ethClient common.EthClient, retrievalClientConfig *RetrieverClientConfig, logger logging.Logger) error {
	// https://github.com/Layr-Labs/eigenda/blob/b8c151436ecefc8046e4aefcdcfee67abf9e8faa/inabox/tests/integration_suite_test.go#L124
	tx, err := eth.NewTransactor(logger, ethClient, retrievalClientConfig.Bls_Operator_State_Retriever, retrievalClientConfig.EigenDA_ServiceManager_Retriever)
	if err != nil {
		return err
	}

	cs := eth.NewChainState(tx, ethClient)
	querier := graphql.NewClient(retrievalClientConfig.ChurnerGraphUrl, nil)
	indexedChainStateClient := thegraph.NewIndexedChainState(cs, querier, logger)
	agn := &core.StdAssignmentCoordinator{}

	// TODO: What should be the value here?
	nodeClient := clients.NewNodeClient(20 * time.Second)
	srsOrder, err := strconv.Atoi(retrievalClientConfig.RetrieverSrsOrder)
	if err != nil {
		return err
	}

	v, err := verifier.NewVerifier(&kzg.KzgConfig{
		G1Path:          retrievalClientConfig.RetrieverG1Path,
		G2Path:          retrievalClientConfig.RetrieverG2Path,
		G2PowerOf2Path:  retrievalClientConfig.RetrieverG2PointPowerOf2Path,
		CacheDir:        retrievalClientConfig.RetrieverCachePath,
		NumWorker:       1,
		SRSOrder:        uint64(srsOrder),
		SRSNumberToLoad: uint64(srsOrder),
		Verbose:         true,
		PreloadEncoder:  false,
		LoadG2Points:    true,
	}, false)
	if err != nil {
		return err
	}

	retrievalClient, err = clients.NewRetrievalClient(logger, indexedChainStateClient, agn, nodeClient, v, 10)
	if err != nil {
		return err
	}

	return indexedChainStateClient.Start(context.Background())
}

// TODO: This file contains some code that can be refactored and shared across some other tests ex:Integration Test.
// This should be done at a later time.
// https://github.com/Layr-Labs/eigenda-internal/issues/221

// Before this test is called by Github Action. expectation is the new updated components are already deployed on K8s cluster
// Disperse Blob over GRPC
// Get Blob Status
// Retrieve Blob using Retriever Client
// Compare Dispersed and Retrieved Blobs
func TestDisperseBlobEndToEnd(t *testing.T) {
	logger := testSuite.Logger

	// Define Different DataSizes
	dataSize := []int{100000, 200000, 1000, 80, 30000}
	// Disperse Blob with different DataSizes
	rand.Seed(time.Now().UnixNano())
	data := make([]byte, dataSize[rand.Intn(len(dataSize))])
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)

	disperseBlobStartTime := time.Now()
	ctx := context.Background()

	logger.Printf("Start to Disperse New Blob")
	// Disperse Blob with QuorumID 0, ConfirmationThreshold 50, AdversaryThreshold 25
	// TODO: Set random values for QuorumID, ConfirmationThreshold, AdversaryThreshold
	key, err := disperseBlob(data, 0, 50, 25)
	logger.Printf("Blob Key After Dispersing %s", hex.EncodeToString(key))
	assert.Nil(t, err)
	assert.NotNil(t, key)
	if key == nil {
		logger.Printf("Blob Key Dispersing Error %s", err.Error())
		t.Fail()
		return
	}

	disperseBlobStopTime := time.Since(disperseBlobStartTime)
	// For now log....later we can define a baseline value for this
	logger.Printf("Time to Disperse Blob %s", disperseBlobStopTime.String())

	// Set Confirmation Deadline For Confirmation of Dispersed Blob
	// Update this to a minute over Batcher_Pull_Interval
	confirmationDeadline, err := time.ParseDuration(testSuite.BatcherPullInterval)

	// Start the loop with a timeout mechanism
	confirmationTicker := time.NewTicker(5 * time.Second)
	defer confirmationTicker.Stop()

	// Blob Confirm Timer
	disperseBlobConfirmStartTime := time.Now()

loop:
	for {
		select {
		case <-ctx.Done():
			logger.Printf("timed out")
		case <-confirmationTicker.C:
			bst, blobReply, err := getBlobStatus(key)
			assert.Nil(t, err)
			assert.NotEqual(t, disperser_rpc.BlobStatus_UNKNOWN, bst)
			// Retrieve Blob if Confirmed
			if bst == disperser_rpc.BlobStatus_CONFIRMED {
				disperseBlobConfirmTime := time.Since(disperseBlobConfirmStartTime)
				logger.Printf("Time to Confirm Dispersed Blob %s", disperseBlobConfirmTime.String())

				// Retrieve Blob from Disperser Client
				logger.Printf("Start to Retrieve Dispersed Blob Using Disperser Endpoint")
				disperserClientBlobReply, err := disperserClientBlobRetrieve(blobReply)
				assert.Nil(t, err)
				assert.NotNil(t, disperserClientBlobReply)

				// Verify DisperserClientData Matches input data
				logger.Printf("Verify Retrieve Dispersed Blob Using Disperser Endpoint Against Dispersersed Blob")
				assert.Equal(t, data, disperserClientBlobReply.Data)

				// Retrieve Blob from Retriever Client
				// Retrieval Client Iterates Over Operators to get the specific Blob
				if isRetrieverClientEnabled {
					logger.Printf("Retrieve Blob using Retrieval Client for %v", blobReply)
					retrieverClientCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
					defer cancel()
					logger.Printf("RetrievalClient:GetBatchHeaderHash() %v", blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash())
					logger.Printf("RetrievalClient:GetBlobIndex() %v", blobReply.GetInfo().GetBlobVerificationProof().GetBlobIndex())
					logger.Printf("RetrievalClient:GetReferenceBlockNumber() %v", blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber())
					logger.Printf("RetrievalClient:GetBatchRoot() %v", blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot())

					retrieved, err := retrievalClient.RetrieveBlob(retrieverClientCtx,
						[32]byte(blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
						blobReply.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
						uint(blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
						[32]byte(blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
						0,
					)
					if err != nil {
						logger.Printf("Error Retrieving Blob %v", err)
					}

					if err == nil {
						logger.Printf("Validate BlobReply %v is equal to inputData %v", bytes.TrimRight(retrieved, "\x00"), data)
						assert.Equal(t, data, bytes.TrimRight(retrieved, "\x00"))

						logger.Printf("Blob Retrieved with Retrieval Client for %v", blobReply)
					}
				}

				if validateOnchainTransaction {
					logger.Printf("Validating OnChain Transaction for Blob with header %v", blobReply.Info.BlobHeader)

					// Verify Blob OnChain
					blobHeader := blobHeaderFromProto(blobReply.Info.BlobHeader)
					logger.Printf("BlobHeader %v", blobHeader)
					verificationProof := blobVerificationProofFromProto(blobReply.GetInfo().GetBlobVerificationProof())
					logger.Printf("VerificationProof %v", verificationProof)

					// Get MockRollUp And EthClient
					mockRollup := testSuite.MockRollUp
					ethClient := testSuite.EthClient

					ethClientCtx, cancel := context.WithTimeout(context.Background(), time.Second*60)
					defer cancel()

					opts, err := ethClient.GetNoSendTransactOpts()
					assert.Nil(t, err)
					tx, err := mockRollup.PostCommitment(opts, blobHeader, verificationProof)
					assert.Nil(t, err)
					assert.NotNil(t, tx)
					logger.Printf("PostCommitment Tx %v", tx)
					_, err = ethClient.EstimateGasPriceAndLimitAndSendTx(ethClientCtx, tx, "PostCommitment", nil)
					assert.Nil(t, err)
				}

				break loop
			}

			// Check if the confirmation process has exceeded the maximum duration
			if time.Now().After(time.Now().Add(confirmationDeadline)) {
				logger.Println("Dispersing Blob Confirmation is taking longer than the specified timeout")
				logger.Println("Failing the test")
				t.Fail()
				return
			}
		}
	}
}

func disperseBlob(data []byte, quorumID uint32, quorumThreshold uint32, adversaryThreshold uint32) ([]byte, error) {
	disperserTestClient := testSuite.Clients.Clients[Disperser]
	logger := testSuite.Logger
	ctxTimeout, cancel := context.WithTimeout(context.Background(), disperserTestClient.Timeout)
	defer cancel()

	request := &disperser_rpc.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: nil,
	}

	disperserClient := disperserTestClient.Client.(disperser_rpc.DisperserClient)

	reply, err := disperserClient.DisperseBlob(ctxTimeout, request)
	if err != nil {
		return nil, err
	}

	logger.Printf("Disperse Blob Reply %v", reply)

	return reply.RequestId, nil
}

func getBlobStatus(key []byte) (disperser_rpc.BlobStatus, *disperser_rpc.BlobStatusReply, error) {
	// Get DisperserTestClient
	disperserTestClient := testSuite.Clients.Clients[Disperser]
	logger := testSuite.Logger

	ctxTimeout, cancel := context.WithTimeout(context.Background(), disperserTestClient.Timeout)
	defer cancel()

	request := &disperser_rpc.BlobStatusRequest{
		RequestId: key,
	}

	disperserClient := disperserTestClient.Client.(disperser_rpc.DisperserClient)

	reply, err := disperserClient.GetBlobStatus(ctxTimeout, request)
	if err != nil {
		return disperser_rpc.BlobStatus_UNKNOWN, nil, err
	}

	logger.Printf("Get Blob Status %v", reply.GetStatus())
	return reply.GetStatus(), reply, nil
}

func disperserClientBlobRetrieve(blobStatusReply *disperser_rpc.BlobStatusReply) (*disperser_rpc.RetrieveBlobReply, error) {
	// Get Disperser TestClient
	disperserTestClient := testSuite.Clients.Clients[Disperser]
	ctxTimeout, cancel := context.WithTimeout(context.Background(), disperserTestClient.Timeout)
	defer cancel()

	batchHeaderHash := blobStatusReply.Info.BlobVerificationProof.BatchMetadata.BatchHeaderHash
	blobIndex := blobStatusReply.Info.BlobVerificationProof.BlobIndex

	request := &disperser_rpc.RetrieveBlobRequest{
		BatchHeaderHash: batchHeaderHash,
		BlobIndex:       blobIndex,
	}

	disperserClient := disperserTestClient.Client.(disperser_rpc.DisperserClient)
	blobReply, err := disperserClient.RetrieveBlob(ctxTimeout, request)

	if err != nil {
		return nil, err
	}

	return blobReply, err
}

func retrieverClientBlobRetrieve(blobStatusReply *disperser_rpc.BlobStatusReply, quorumID uint32) (*retriever_rpc.BlobReply, error) {
	logger := testSuite.Logger
	// Get Retriever TestClient
	retrieverTestClient := testSuite.Clients.Clients[Retriever]
	ctxTimeout, cancel := context.WithTimeout(context.Background(), retrieverTestClient.Timeout)
	defer cancel()

	batchHeaderHash := blobStatusReply.Info.BlobVerificationProof.BatchMetadata.BatchHeaderHash
	blobIndex := blobStatusReply.Info.BlobVerificationProof.BlobIndex
	referenceBlockNumber := blobStatusReply.Info.BlobVerificationProof.BatchMetadata.BatchHeader.ReferenceBlockNumber

	request := &retriever_rpc.BlobRequest{
		BatchHeaderHash:      batchHeaderHash,
		BlobIndex:            blobIndex,
		ReferenceBlockNumber: referenceBlockNumber,
		QuorumId:             quorumID,
	}

	retrieverClient := retrieverTestClient.Client.(retriever_rpc.RetrieverClient)
	blobReply, err := retrieverClient.RetrieveBlob(ctxTimeout, request)

	if err != nil {
		return nil, err
	}
	logger.Printf("Retrieve Blob %v", blobReply)
	return blobReply, err

}

func blobHeaderFromProto(blobHeader *disperser_rpc.BlobHeader) rollupbindings.IEigenDAServiceManagerBlobHeader {
	commitX := new(fp.Element).SetBytes(blobHeader.GetCommitment().GetX())
	commitY := new(fp.Element).SetBytes(blobHeader.GetCommitment().GetY())
	commitment := &encoding.G1Commitment{
		X: *commitX,
		Y: *commitY,
	}

	quorums := make([]rollupbindings.IEigenDAServiceManagerQuorumBlobParam, len(blobHeader.GetBlobQuorumParams()))
	for i, quorum := range blobHeader.GetBlobQuorumParams() {
		quorums[i] = rollupbindings.IEigenDAServiceManagerQuorumBlobParam{
			QuorumNumber:                    uint8(quorum.GetQuorumNumber()),
			AdversaryThresholdPercentage:    uint8(quorum.GetAdversaryThresholdPercentage()),
			ConfirmationThresholdPercentage: uint8(quorum.GetConfirmationThresholdPercentage()),
			ChunkLength:                     quorum.GetChunkLength(),
		}
	}

	return rollupbindings.IEigenDAServiceManagerBlobHeader{
		Commitment: rollupbindings.BN254G1Point{
			X: commitment.X.BigInt(new(big.Int)),
			Y: commitment.Y.BigInt(new(big.Int)),
		},
		DataLength:       blobHeader.GetDataLength(),
		QuorumBlobParams: quorums,
	}
}

func blobVerificationProofFromProto(verificationProof *disperser_rpc.BlobVerificationProof) rollupbindings.EigenDARollupUtilsBlobVerificationProof {
	logger := testSuite.Logger
	batchMetadataProto := verificationProof.GetBatchMetadata()
	batchHeaderProto := verificationProof.GetBatchMetadata().GetBatchHeader()
	var batchRoot [32]byte
	copy(batchRoot[:], batchHeaderProto.GetBatchRoot())
	batchHeader := rollupbindings.IEigenDAServiceManagerBatchHeader{
		BlobHeadersRoot:       batchRoot,
		QuorumNumbers:         batchHeaderProto.GetQuorumNumbers(),
		SignedStakeForQuorums: batchHeaderProto.GetQuorumSignedPercentages(),
		ReferenceBlockNumber:  batchHeaderProto.GetReferenceBlockNumber(),
	}
	var sig [32]byte
	copy(sig[:], batchMetadataProto.GetSignatoryRecordHash())
	logger.Printf("VerificationProof:SignatoryRecordHash: %v\n", sig)
	logger.Printf("VerificationProof:ConfirmationBlockNumber: %v\n", batchMetadataProto.GetConfirmationBlockNumber())
	batchMetadata := rollupbindings.IEigenDAServiceManagerBatchMetadata{
		BatchHeader:             batchHeader,
		SignatoryRecordHash:     sig,
		ConfirmationBlockNumber: batchMetadataProto.GetConfirmationBlockNumber(),
	}

	logger.Printf("VerificationProof:BatchId: %v\n", verificationProof.GetBatchId())
	logger.Printf("VerificationProof:BlobIndex: %v\n", uint8(verificationProof.GetBlobIndex()))
	logger.Printf("VerificationProof:BatchMetadata: %v\n", batchMetadata)
	logger.Printf("VerificationProof:InclusionProof: %v\n", verificationProof.GetInclusionProof())
	logger.Printf("VerificationProof:QuorumThresholdIndexes: %v\n", verificationProof.GetQuorumIndexes())

	return rollupbindings.EigenDARollupUtilsBlobVerificationProof{
		BatchId:        verificationProof.GetBatchId(),
		BlobIndex:      verificationProof.GetBlobIndex(),
		BatchMetadata:  batchMetadata,
		InclusionProof: verificationProof.GetInclusionProof(),
		QuorumIndices:  verificationProof.GetQuorumIndexes(),
	}
}

func createClient(conn *grpc.ClientConn, clientType ClientType) interface{} {
	switch clientType {
	case Disperser:
		return disperser_rpc.NewDisperserClient(conn)
	case Retriever:
		return retriever_rpc.NewRetrieverClient(conn)
	default:
		return nil
	}
}

func newTestClients(clients map[ClientType]*GrpcClient) (*TestClients, error) {
	grpcClients := make(map[ClientType]*GrpcClient)

	for clientType, grpcClient := range clients {
		addr := fmt.Sprintf("%s:%s", grpcClient.Hostname, grpcClient.GrpcPort)
		conn, err := grpc.Dial(addr, grpc.WithInsecure()) // Using insecure connection for simplicity, consider using secure connection in production
		if err != nil {
			return nil, fmt.Errorf("failed to dial %s server: %v", clientType, err)
		}

		grpcClient.Client = createClient(conn, clientType) // Create specific client based on client type
		grpcClients[clientType] = grpcClient
	}

	return &TestClients{
		Clients: clients,
	}, nil
}
