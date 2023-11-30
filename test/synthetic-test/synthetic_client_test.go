//go:build k8s_func_test
// +build k8s_func_test

package k8s_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shurcooL/graphql"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	retriever_rpc "github.com/Layr-Labs/eigenda/api/grpc/retriever"
	"github.com/Layr-Labs/eigenda/clients"
	common "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	rollupbindings "github.com/Layr-Labs/eigenda/contracts/bindings/MockRollup"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/encoding"
	"github.com/Layr-Labs/eigenda/core/eth"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	encoder_rpc "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/pkg/encoding/kzgEncoder"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type ClientType string

const (
	Disperser ClientType = "Disperser"
	Retriever ClientType = "Retriever"
	Encoder   ClientType = "Encoder"
)

type GrpcClient struct {
	Hostname string
	GrpcPort string
	Timeout  time.Duration
	Client   interface{} // This can be a specific client type (disperser_rpc.DisperserClient, retriever_rpc.RetrieverClient, etc.)
}

type RetrieverClientConfig struct {
	ChurnerGraphUrl    string
	RetrieverSrsOrder  string
	RetrieverG1Path    string
	RetrieverG2Path    string
	RetrieverCachePath string
}

type TestClients struct {
	Clients map[ClientType]*GrpcClient
}

// TestSuite Struct.
type SyntheticTestSuite struct {
	Clients         *TestClients
	EthClient       common.EthClient
	MockRollUp      *rollupbindings.ContractMockRollup
	RetrievalClient *retrievalClient
	Logger          *log.Logger
	TestRunID       string
}

var (
	testSuite                  *SyntheticTestSuite
	isRetrieverClientDeployed  bool = false
	validateOnchainTransaction bool = false
)

func setUpClients(pk string, rpcUrl string, mockRollUpContractAddress string, retrieverClientConfig RetrieverClientConfig) *SyntheticTestSuite {
	testRunID := uuid.New().String()
	logger := log.New(os.Stdout, "EigenDA SyntheticClient:"+testRunID+" ", log.Ldate|log.Ltime)

	// Initialize your TestClients and other suite fields here
	clients, err := newTestClients(map[ClientType]*GrpcClient{
		Disperser: {
			Hostname: "disperser-goerli.eigenda-testnet.eigenops.xyz",
			GrpcPort: "443",
			Timeout:  10 * time.Second,
		},
		Retriever: {
			Hostname: "retriever.retriever.svc.cluster.local",
			GrpcPort: "retriever-port",
			Timeout:  10 * time.Second,
		},
		Encoder: {
			Hostname: "encoder.encoder.svc.cluster.local",
			GrpcPort: "34000",
			Timeout:  10 * time.Second,
		},
	})
	if err != nil {
		logger.Printf("Error initializing clients: %v", err)
	}

	ethLogger, err := logging.GetLogger(logging.DefaultCLIConfig())
	if err != nil {
		// Handle the error
		logger.Printf("Error: %v", err)
		return nil
	}
	pk = strings.TrimPrefix(pk, "0X")
	pk = strings.TrimPrefix(pk, "0x")
	ethClient, err := geth.NewClient(geth.EthClientConfig{
		RPCURL:           rpcUrl,
		PrivateKeyString: pk,
	}, ethLogger)

	if err != nil {
		logger.Printf("Error: %v", err)
	}

	mockRollup, err := rollupbindings.NewContractMockRollup(gcommon.HexToAddress(mockRollUpContractAddress), ethClient)

	retrievalClient, err := setupRetrievalClient(ethClient, &deploy.Config{

	// Assign client connections to pointers in TestClients struct
	return &SyntheticTestSuite{
		Clients:    clients,
		EthClient:  ethClient,
		MockRollUp: mockRollup,
		Logger:     logger,
	}
}

func TestMain(m *testing.M) {

	// These are set in deployment/cronjob yaml
	privateKey := os.Getenv("ETHCLIENT_PRIVATE_KEY")
	rpcUrl := os.Getenv("ETHCLIENT_RPC_URL")
	mockRollUpContractAddress := os.Getenv("MOCKROLLUP_CONTRACT_ADDRESS")
	isRetrieverClientDeployed = os.Getenv("RETRIEVER_CLIENT_DEPLOYED") == strings.ToLower("true")
	validateOnchainTransaction = os.Getenv("VALIDATE_ONCHAIN_TRANSACTION") == strings.ToLower("true")

	// Retriever Config
	retrieverClientConfig := &RetrieverClientConfig{
		ChurnerGraphUrl:    os.Getenv("RETRIEVER_CHURNER_GRAPH_URL"),
		RetrieverSrsOrder:  os.Getenv("RETRIEVER_SRS_ORDER"),
		RetrieverG1Path:    os.Getenv("RETRIEVER_G1_PATH"),
		RetrieverG2Path:    os.Getenv("RETRIEVER_G2_PATH"),
		RetrieverCachePath: os.Getenv("RETRIEVER_CACHE_PATH"),
	}

	// Initialize Clients
	testSuite = setUpClients(privateKey, rpcUrl, mockRollUpContractAddress, retrieverClientConfig)
	logger := testSuite.Logger

	// Check if testSuite is nil
	if testSuite == nil {
		logger.Println("Error: TestSuite initialization failed.")
		os.Exit(1) // Exit the program with a non-zero status code indicating failure
	}
	logger.Println("RPC_URL for Chain...", rpcUrl)
	logger.Println("Mock RollUp Contract Address...", mockRollUpContractAddress)
	logger.Println("Retriever Client Deployed...", isRetrieverClientDeployed)

	logger.Println("Running Test Client...")
	// Run the tests and get the exit code
	exitCode := m.Run()
	logger.Printf("Exiting Test Client Run with Code:%d", exitCode)
	// Exit with the test result code
	os.Exit(exitCode)
}

func setupRetrievalClient(ethClient common.EthClient, retrievalClientConfig *RetrieverClientConfig, logger *log.Logger) (retrievalClient, error) {
	// https://github.com/Layr-Labs/eigenda/blob/b8c151436ecefc8046e4aefcdcfee67abf9e8faa/inabox/tests/integration_suite_test.go#L124
	tx, err := eth.NewTransactor(logger, ethClient, testConfig.Retriever.RETRIEVER_BLS_OPERATOR_STATE_RETRIVER, testConfig.Retriever.RETRIEVER_EIGENDA_SERVICE_MANAGER)
	if err != nil {
		return err
	}

	cs := eth.NewChainState(tx, client)
	querier := graphql.NewClient(retrievalClientConfig.CHURNER_GRAPH_URL, nil)
	indexedChainStateClient := thegraph.NewIndexedChainState(cs, querier, logger)
	agn := &core.StdAssignmentCoordinator{}

	// TODO: What should be the value here?
	nodeClient := clients.NewNodeClient(20 * time.Second)
	srsOrder, err := strconv.Atoi(retrievalClientConfig.RETRIEVER_SRS_ORDER)
	if err != nil {
		return nil, err
	}
	encoder, err := encoding.NewEncoder(encoding.EncoderConfig{
		KzgConfig: kzgEncoder.KzgConfig{
			G1Path:         retrievalClientConfig.RETRIEVER_G1_PATH,
			G2Path:         retrievalClientConfig.RETRIEVER_G2_PATH,
			CacheDir:       retrievalClientConfig.RETRIEVER_CACHE_PATH,
			NumWorker:      1,
			SRSOrder:       uint64(srsOrder),
			Verbose:        true,
			PreloadEncoder: true,
		},
	})
	if err != nil {
		return nil, err
	}

	retrievalClient = clients.NewRetrievalClient(logger, indexedChainStateClient, agn, nodeClient, encoder, 10)
	return retrievalClient, nil
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
	disperseBlobStartTime := time.Now()
	ctx := context.Background()

	logger.Printf("Start to Disperse New Blob")
	// Disperse Blob with QuorumID 0, QuorumThreshold 50, AdversaryThreshold 25
	// TODO: Set random values for QuorumID, QuorumThreshold, AdversaryThreshold
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

	// Set Confirmation DeaLine For Confirmation of Dispersed Blob
	confirmationDeadline := time.Now().Add(240 * time.Second)

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

					tx, err := mockRollup.PostCommitment(ethClient.GetNoSendTransactOpts(), blobHeader, verificationProof)
					assert.Nil(t, err)
					assert.NotNil(t, tx)
					logger.Printf("PostCommitment Tx %v", tx)
					_, err = ethClient.EstimateGasPriceAndLimitAndSendTx(ethClientCtx, tx, "PostCommitment", nil)
					assert.Nil(t, err)
				}

				// Retrieve Blob from Retriever Client
				if isRetrieverClientEnabled {
					logger.Printf("Retry Blob using Retrieval Client %v", blobReply)
					retrieverClientCtx, cancel = context.WithTimeout(context.Background(), time.Second*5)
					defer cancel()
					logger.Printf("RetrievalClient:GetBatchHeaderHash()", blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash())
					logger.Printf("RetrievalClient:GetBlobIndex()", blobReply.GetInfo().GetBlobVerificationProof().GetBlobIndex())
					logger.Printf("RetrievalClient:GetReferenceBlockNumber()", blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber())
					logger.Printf("RetrievalClient:GetBatchRoot()", blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot())
					
					retrieved, err := retrievalClient.RetrieveBlob(retrieverClientCtx,
						[32]byte(blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeaderHash()),
						blobReply.GetInfo().GetBlobVerificationProof().GetBlobIndex(),
						uint(blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetReferenceBlockNumber()),
						[32]byte(blobReply.GetInfo().GetBlobVerificationProof().GetBatchMetadata().GetBatchHeader().GetBatchRoot()),
						0,
					)
					assert.Nil(t, err)
					assert.Equal(t, data, bytes.TrimRight(retrieved, "\x00"))
				
					logger.Printf("Retry Blob using Retrieval Client %v", blobReply)
				}

				break loop
			}

			// Check if the confirmation process has exceeded the maximum duration
			if time.Now().After(confirmationDeadline) {
				logger.Println("Dispersing Blob Confirmation is taking longer than the specified timeout of 4 minutes")
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
		Data: data,
		SecurityParams: []*disperser_rpc.SecurityParams{
			{
				QuorumId:           quorumID,
				QuorumThreshold:    quorumThreshold,
				AdversaryThreshold: adversaryThreshold,
			},
		},
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
	logger := testSuite.Logger
	commitmentBytes := blobHeader.GetCommitment()
	commitment, err := new(core.Commitment).Deserialize(commitmentBytes)
	if err != nil {
		logger.Printf("failed to deserialize commitment: %s", err)
		return rollupbindings.IEigenDAServiceManagerBlobHeader{}
	}
	type IEigenDAServiceManagerQuorumBlobParam struct {
		QuorumNumber                 uint8
		AdversaryThresholdPercentage uint8
		QuorumThresholdPercentage    uint8
		QuantizationParameter        uint8
	}
	quorums := make([]rollupbindings.IEigenDAServiceManagerQuorumBlobParam, len(blobHeader.GetBlobQuorumParams()))
	for i, quorum := range blobHeader.GetBlobQuorumParams() {
		quorums[i] = rollupbindings.IEigenDAServiceManagerQuorumBlobParam{
			QuorumNumber:                 uint8(quorum.GetQuorumNumber()),
			AdversaryThresholdPercentage: uint8(quorum.GetAdversaryThresholdPercentage()),
			QuorumThresholdPercentage:    uint8(quorum.GetQuorumThresholdPercentage()),
			QuantizationParameter:        uint8(quorum.GetQuantizationParam()),
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

func blobVerificationProofFromProto(verificationProof *disperser_rpc.BlobVerificationProof) rollupbindings.EigenDABlobUtilsBlobVerificationProof {
	logger := testSuite.Logger
	batchMetadataProto := verificationProof.GetBatchMetadata()
	batchHeaderProto := verificationProof.GetBatchMetadata().GetBatchHeader()
	var batchRoot [32]byte
	copy(batchRoot[:], batchHeaderProto.GetBatchRoot())
	batchHeader := rollupbindings.IEigenDAServiceManagerBatchHeader{
		BlobHeadersRoot:            batchRoot,
		QuorumNumbers:              batchHeaderProto.GetQuorumNumbers(),
		QuorumThresholdPercentages: batchHeaderProto.GetQuorumSignedPercentages(),
		ReferenceBlockNumber:       batchHeaderProto.GetReferenceBlockNumber(),
	}
	var sig [32]byte
	copy(sig[:], batchMetadataProto.GetSignatoryRecordHash())
	fee := new(big.Int).SetBytes(batchMetadataProto.GetFee())
	logger.Printf("VerificationProof:SignatoryRecordHash: %v\n", sig)
	logger.Printf("VerificationProof:ConfirmationBlockNumber: %v\n", batchMetadataProto.GetConfirmationBlockNumber())
	batchMetadata := rollupbindings.IEigenDAServiceManagerBatchMetadata{
		BatchHeader:             batchHeader,
		SignatoryRecordHash:     sig,
		Fee:                     fee,
		ConfirmationBlockNumber: batchMetadataProto.GetConfirmationBlockNumber(),
	}

	logger.Printf("VerificationProof:BatchId: %v\n", verificationProof.GetBatchId())
	logger.Printf("VerificationProof:BlobIndex: %v\n", uint8(verificationProof.GetBlobIndex()))
	logger.Printf("VerificationProof:BatchMetadata: %v\n", batchMetadata)
	logger.Printf("VerificationProof:InclusionProof: %v\n", verificationProof.GetInclusionProof())
	logger.Printf("VerificationProof:QuorumThresholdIndexes: %v\n", verificationProof.GetQuorumIndexes())

	return rollupbindings.EigenDABlobUtilsBlobVerificationProof{
		BatchId:                verificationProof.GetBatchId(),
		BlobIndex:              uint8(verificationProof.GetBlobIndex()),
		BatchMetadata:          batchMetadata,
		InclusionProof:         verificationProof.GetInclusionProof(),
		QuorumThresholdIndexes: verificationProof.GetQuorumIndexes(),
	}
}

func TestEncodeBlob(t *testing.T) {
	t.Skip("Skipping this test")

	var (
		gettysburgAddressBytes = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
	)
	encoderReply, encodingParams, err := encodeBlob(gettysburgAddressBytes)
	assert.NoError(t, err)
	assert.NotNil(t, encoderReply.Chunks)

	// Decode Server Data
	var chunksData []*core.Chunk

	for i := range encoderReply.Chunks {
		chunkSerialized, _ := new(core.Chunk).Deserialize(encoderReply.GetChunks()[i])
		// perform an operation
		chunksData = append(chunksData, chunkSerialized)
	}
	assert.NotNil(t, chunksData)

	// Indices obtained from Encoder_Test
	indices := []core.ChunkNumber{
		0, 1, 2, 3, 4, 5, 6, 7,
	}

	// Test Assumes below params set for Encoder
	kzgConfig := kzgEncoder.KzgConfig{
		G1Path:    "/data/kzg/g1.point",
		G2Path:    "/data/kzg/g2.point",
		CacheDir:  "/data/kzg/SRSTables",
		SRSOrder:  300000,
		NumWorker: uint64(runtime.GOMAXPROCS(0)),
	}

	encodingConfig := encoding.EncoderConfig{KzgConfig: kzgConfig}

	encoder, _ := encoding.NewEncoder(encodingConfig)

	maxInputSize := uint64(len(gettysburgAddressBytes)) + 10
	decoded, err := encoder.Decode(chunksData, indices, *encodingParams, maxInputSize)
	assert.Nil(t, err)
	assert.Equal(t, decoded, gettysburgAddressBytes)
}

func encodeBlob(data []byte) (*encoder_rpc.EncodeBlobReply, *core.EncodingParams, error) {
	logger := testSuite.Logger
	var quantizationFactor uint = 2
	var adversaryThreshold uint8 = 80
	var quorumThreshold uint8 = 90

	encoderTestClient := testSuite.Clients.Clients[Encoder]
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(10*float64(time.Second)))
	defer cancel()

	var quorumID core.QuorumID = 0

	securityParams := []*core.SecurityParam{
		{
			QuorumID:           quorumID,
			AdversaryThreshold: adversaryThreshold,
		},
	}

	testBlob := core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: securityParams,
		},
		Data: data,
	}
	// TODO: Refactor this code using indexed chain state by using retrieval client
	// Issue: https://github.com/Layr-Labs/eigenda-internal/issues/220
	indexedChainState, _ := coremock.NewChainDataMock(core.OperatorIndex(10))
	operatorState, err := indexedChainState.GetOperatorState(context.Background(), uint(0), []core.QuorumID{quorumID})
	if err != nil {
		logger.Printf("failed to get operator state: %s", err)
	}
	coordinator := &core.StdAssignmentCoordinator{}
	_, info, err := coordinator.GetAssignments(operatorState, quorumID, quantizationFactor)
	if err != nil {
		logger.Printf("failed to get assignments: %s", err)
	}

	testBlob = core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: securityParams,
		},
		Data: data,
	}

	blobSize := uint(len(testBlob.Data))

	blobLength := core.GetBlobLength(uint(blobSize))

	chunkLength, _ := coordinator.GetMinimumChunkLength(4, blobLength, 1, quorumThreshold, adversaryThreshold)

	testEncodingParams, _ := core.GetEncodingParams(chunkLength, info.TotalChunks)

	testEncodingParamsProto := &encoder_rpc.EncodingParams{
		ChunkLength: uint32(testEncodingParams.ChunkLength),
		NumChunks:   uint32(testEncodingParams.NumChunks),
	}

	encodeBlobRequestProto := &encoder_rpc.EncodeBlobRequest{
		Data:           []byte(testBlob.Data),
		EncodingParams: testEncodingParamsProto,
	}
	encoderClient := encoderTestClient.Client.(encoder_rpc.EncoderClient)

	reply, err := encoderClient.EncodeBlob(ctxTimeout, encodeBlobRequestProto)
	return reply, &testEncodingParams, err
}

func createClient(conn *grpc.ClientConn, clientType ClientType) interface{} {
	switch clientType {
	case Disperser:
		return disperser_rpc.NewDisperserClient(conn)
	case Retriever:
		return retriever_rpc.NewRetrieverClient(conn)
	case Encoder:
		return encoder_rpc.NewEncoderClient(conn)
	default:
		return nil
	}
}

func newTestClients(clients map[ClientType]*GrpcClient) (*TestClients, error) {
	grpcClients := make(map[ClientType]*GrpcClient)

	for clientType, grpcClient := range clients {
		addr := fmt.Sprintf("%s:%s", grpcClient.Hostname, grpcClient.GrpcPort)
		fmt.Println("Address: ", addr, " ClientType: ", clientType)

		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(nil)), grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(math.MaxInt32), // Increase frame size limit
		))

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
