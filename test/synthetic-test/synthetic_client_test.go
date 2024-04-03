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
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/google/uuid"
	"github.com/shurcooL/graphql"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	retriever_rpc "github.com/Layr-Labs/eigenda/api/grpc/retriever"
	"github.com/Layr-Labs/eigenda/clients"
	common "github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	rollupbindings "github.com/Layr-Labs/eigenda/contracts/bindings/MockRollup"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	coremock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	encoder_rpc "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder"
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
	Encoder   ClientType = "Encoder"
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
}

type TestClients struct {
	Clients map[ClientType]*GrpcClient
}

// TestSuite Struct.
type SyntheticTestSuite struct {
	Clients         *TestClients
	EthClient       common.EthClient
	MockRollUp      *rollupbindings.ContractMockRollup
	Logger          *log.Logger
	RetrievalClient clients.RetrievalClient
	TestRunID       string
}

var (
	testSuite                  *SyntheticTestSuite
	isRetrieverClientDeployed  bool = false
	validateOnchainTransaction bool = false
	retrievalClient            clients.RetrievalClient
	logger                     logging.Logger
)

func setUpClients(pk string, rpcUrl string, mockRollUpContractAddress string, retrieverClientConfig RetrieverClientConfig) *SyntheticTestSuite {
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
		Encoder: {
			Hostname: "encoder.encoder.svc.cluster.local",
			GrpcPort: "34000",
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
	log.Printf("Error: failed to create eth client: %v", err)
	if err != nil {
		log.Printf("Error: failed to create eth client: %v", err)
	}

	mockRollup, err := rollupbindings.NewContractMockRollup(gcommon.HexToAddress(mockRollUpContractAddress), ethClient)
	if err != nil {
		logger.Printf("Error: %v", err)
		return nil
	}

	err = setupRetrievalClient(ethClient, &retrieverClientConfig, ethLogger)

	if err != nil {
		logger.Printf("Error: %v", err)
		return nil
	}

	// Assign client connections to pointers in TestClients struct
	return &SyntheticTestSuite{
		Clients:         clients,
		EthClient:       ethClient,
		MockRollUp:      mockRollup,
		RetrievalClient: retrievalClient,
		Logger:          logger,
	}
}

func TestMain(m *testing.M) {

	// These are set in deployment/cronjob yaml
	privateKey := os.Getenv("ETHCLIENT_PRIVATE_KEY")
	rpcUrl := os.Getenv("ETHCLIENT_RPC_URL")
	mockRollUpContractAddress := os.Getenv("MOCKROLLUP_CONTRACT_ADDRESS")
	isRetrieverClientDeployed = os.Getenv("RETRIEVER_CLIENT_DEPLOYED") == strings.ToLower("true")
	validateOnchainTransaction = os.Getenv("VALIDATE_ONCHAIN_TRANSACTION") == strings.ToLower("true")
	blsOperatorStateRetriever := os.Getenv("BLS_OPERATOR_STATE_RETRIEVER")
	eigenDAServiceManagerRetreiever := os.Getenv("EIGENDA_SERVICE_MANAGER_RETRIEVER")
	churnerGraphUrl := os.Getenv("CHURNER_GRAPH_URL")
	retrieverSrsOrder := os.Getenv("RETRIEVER_SRS_ORDER")
	retrieverG1Path := os.Getenv("RETRIEVER_G1_PATH")
	retrieverG2Path := os.Getenv("RETRIEVER_G2_PATH")
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
		RetrieverCachePath:               retrieverCachePath,
	}

	// Initialize Clients
	testSuite = setUpClients(privateKey, rpcUrl, mockRollUpContractAddress, *retrieverClientConfig)
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
		G1Path:         retrievalClientConfig.RetrieverG1Path,
		G2Path:         retrievalClientConfig.RetrieverG2Path,
		CacheDir:       retrievalClientConfig.RetrieverCachePath,
		NumWorker:      1,
		SRSOrder:       uint64(srsOrder),
		Verbose:        true,
		PreloadEncoder: true,
	}, true)
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

	// Set Confirmation DeaLine For Confirmation of Dispersed Blob
	// Update this to a minute over Batcher_Pull_Interval
	confirmationDeadline := time.Now().Add(batcherPullInterval * time.Second)

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

					opts, err := ethClient.GetNoSendTransactOpts()
					assert.Nil(t, err)
					tx, err := mockRollup.PostCommitment(opts, blobHeader, verificationProof)
					assert.Nil(t, err)
					assert.NotNil(t, tx)
					logger.Printf("PostCommitment Tx %v", tx)
					_, err = ethClient.EstimateGasPriceAndLimitAndSendTx(ethClientCtx, tx, "PostCommitment", nil)
					assert.Nil(t, err)
				}

				// Retrieve Blob from Retriever Client
				// Retrieval Client Iterates Over Operators to get the specific Blob
				if isRetrieverClientDeployed {
					logger.Printf("Try Blob using Retrieval Client %v", blobReply)
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
					assert.Nil(t, err)
					logger.Printf("Validate BlobReply %v is equal to inputData %v", bytes.TrimRight(retrieved, "\x00"), data)
					assert.Equal(t, data, bytes.TrimRight(retrieved, "\x00"))

					logger.Printf("Blob using Retrieval Client %v", blobReply)
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
		BlobIndex:      uint8(verificationProof.GetBlobIndex()),
		BatchMetadata:  batchMetadata,
		InclusionProof: verificationProof.GetInclusionProof(),
		QuorumIndices:  verificationProof.GetQuorumIndexes(),
	}
}

func TestEncodeBlob(t *testing.T) {
	t.Skip("Skipping this test")

	var (
		gettysburgAddressBytes = codec.ConvertByPaddingEmptyByte([]byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth."))
	)
	encoderReply, encodingParams, err := encodeBlob(gettysburgAddressBytes)
	assert.NoError(t, err)
	assert.NotNil(t, encoderReply.Chunks)

	// Decode Server Data
	var chunksData []*encoding.Frame

	for i := range encoderReply.Chunks {
		chunkSerialized, _ := new(encoding.Frame).Deserialize(encoderReply.GetChunks()[i])
		// perform an operation
		chunksData = append(chunksData, chunkSerialized)
	}
	assert.NotNil(t, chunksData)

	// Indices obtained from Encoder_Test
	indices := []encoding.ChunkNumber{
		0, 1, 2, 3, 4, 5, 6, 7,
	}

	// Test Assumes below params set for Encoder
	kzgConfig := kzg.KzgConfig{
		G1Path:          "/data/kzg/g1.point",
		G2Path:          "/data/kzg/g2.point",
		CacheDir:        "/data/kzg/SRSTables",
		SRSOrder:        300000,
		SRSNumberToLoad: 300000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	v, _ := verifier.NewVerifier(&kzgConfig, false)

	maxInputSize := uint64(len(gettysburgAddressBytes)) + 10
	decoded, err := v.Decode(chunksData, indices, *encodingParams, maxInputSize)
	assert.Nil(t, err)
	assert.Equal(t, decoded, gettysburgAddressBytes)
}

func encodeBlob(data []byte) (*encoder_rpc.EncodeBlobReply, *encoding.EncodingParams, error) {
	logger := testSuite.Logger
	var adversaryThreshold uint8 = 80
	var quorumThreshold uint8 = 90

	encoderTestClient := testSuite.Clients.Clients[Encoder]
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(10*float64(time.Second)))
	defer cancel()

	var quorumID core.QuorumID = 0

	param := &core.SecurityParam{
		QuorumID:              quorumID,
		ConfirmationThreshold: quorumThreshold,
		AdversaryThreshold:    adversaryThreshold,
	}

	testBlob := core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: []*core.SecurityParam{param},
		},
		Data: data,
	}
	// TODO: Refactor this code using indexed chain state by using retrieval client
	// Issue: https://github.com/Layr-Labs/eigenda-internal/issues/220
	indexedChainState, _ := coremock.MakeChainDataMock(core.OperatorIndex(10))
	operatorState, err := indexedChainState.GetOperatorState(context.Background(), uint(0), []core.QuorumID{quorumID})
	if err != nil {
		logger.Printf("failed to get operator state: %s", err)
	}
	coordinator := &core.StdAssignmentCoordinator{}

	blobSize := uint(len(testBlob.Data))
	blobLength := encoding.GetBlobLength(uint(blobSize))

	chunkLength, err := coordinator.CalculateChunkLength(operatorState, blobLength, 0, param)
	if err != nil {
		logger.Printf("failed to calculate chunk length: %s", err)
	}

	quorumInfo := &core.BlobQuorumInfo{
		SecurityParam: *param,
		ChunkLength:   chunkLength,
	}

	_, info, err := coordinator.GetAssignments(operatorState, blobLength, quorumInfo)
	if err != nil {
		logger.Printf("failed to get assignments: %s", err)
	}
	testEncodingParams := encoding.ParamsFromMins(chunkLength, info.TotalChunks)

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
