package apiserver_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
	"google.golang.org/grpc/peer"
)

var (
	queue           disperser.BlobStore
	dispersalServer *apiserver.DispersalServer
	transactor      *mock.MockTransactor

	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	UUID               = uuid.New()
	metadataTableName  = fmt.Sprintf("test-BlobMetadata-%v", UUID)
	bucketTableName    = fmt.Sprintf("test-BucketStore-%v", UUID)

	deployLocalStack bool
	localStackPort   = "4568"
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestDisperseBlob(t *testing.T) {
	data := make([]byte, 3*1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	status, _, key := disperseBlob(t, dispersalServer, data)
	assert.Equal(t, status, pb.BlobStatus_PROCESSING)
	assert.NotNil(t, key)
}

func TestDisperseBlobWithRequiredQuorums(t *testing.T) {
	data := make([]byte, 1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	quorumParams := []*core.SecurityParam{
		{QuorumID: 0, AdversaryThreshold: 50, ConfirmationThreshold: 100},
		{QuorumID: 1, AdversaryThreshold: 50, ConfirmationThreshold: 100},
	}
	transactor.On("GetQuorumSecurityParams", tmock.Anything).Return(quorumParams, nil)
	transactor.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{0, 1}, nil)

	reply, err := dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:          data,
		QuorumNumbers: []uint32{1},
	})
	assert.NoError(t, err)

	assert.Equal(t, reply.GetResult(), pb.BlobStatus_PROCESSING)

	requestID := reply.GetRequestId()
	assert.NotNil(t, requestID)

}

func TestDisperseBlobWithInvalidQuorum(t *testing.T) {
	data := make([]byte, 1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	quorumParams := []*core.SecurityParam{
		{QuorumID: 0, AdversaryThreshold: 50, ConfirmationThreshold: 100},
		{QuorumID: 1, AdversaryThreshold: 50, ConfirmationThreshold: 100},
	}
	transactor.On("GetQuorumSecurityParams", tmock.Anything).Return(quorumParams, nil)
	transactor.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{}, nil)

	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:          data,
		QuorumNumbers: []uint32{2},
	})
	assert.ErrorContains(t, err, "invalid request: the quorum_id must be in range [0, 1], but found 2")

	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:          data,
		QuorumNumbers: []uint32{0, 0},
	})
	assert.ErrorContains(t, err, "invalid request: security_params must not contain duplicate quorum_id")
}

func TestGetBlobStatus(t *testing.T) {
	data := make([]byte, 1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	status, blobSize, requestID := disperseBlob(t, dispersalServer, data)
	assert.Equal(t, status, pb.BlobStatus_PROCESSING)
	assert.NotNil(t, requestID)

	reply, err := dispersalServer.GetBlobStatus(context.Background(), &pb.BlobStatusRequest{
		RequestId: requestID,
	})
	assert.NoError(t, err)
	assert.Equal(t, reply.GetStatus(), pb.BlobStatus_PROCESSING)

	// simulate blob confirmation
	securityParams := []*core.SecurityParam{
		{
			QuorumID:              0,
			AdversaryThreshold:    80,
			ConfirmationThreshold: 100,
		},
		{
			QuorumID:              1,
			AdversaryThreshold:    80,
			ConfirmationThreshold: 100,
		},
	}
	confirmedMetadata := simulateBlobConfirmation(t, requestID, blobSize, securityParams, 0)

	reply, err = dispersalServer.GetBlobStatus(context.Background(), &pb.BlobStatusRequest{
		RequestId: requestID,
	})
	assert.NoError(t, err)
	assert.Equal(t, reply.GetStatus(), pb.BlobStatus_CONFIRMED)
	actualCommitX := reply.GetInfo().GetBlobHeader().GetCommitment().X
	actualCommitY := reply.GetInfo().GetBlobHeader().GetCommitment().Y
	assert.Equal(t, actualCommitX, confirmedMetadata.ConfirmationInfo.BlobCommitment.Commitment.X.Marshal())
	assert.Equal(t, actualCommitY, confirmedMetadata.ConfirmationInfo.BlobCommitment.Commitment.Y.Marshal())
	assert.Equal(t, reply.GetInfo().GetBlobHeader().GetDataLength(), uint32(confirmedMetadata.ConfirmationInfo.BlobCommitment.Length))

	actualBlobQuorumParams := make([]*pb.BlobQuorumParam, len(securityParams))
	quorumNumbers := make([]byte, len(securityParams))
	quorumPercentSigned := make([]byte, len(securityParams))
	quorumIndexes := make([]byte, len(securityParams))
	for i, sp := range securityParams {
		actualBlobQuorumParams[i] = &pb.BlobQuorumParam{
			QuorumNumber:                    uint32(sp.QuorumID),
			AdversaryThresholdPercentage:    uint32(sp.AdversaryThreshold),
			ConfirmationThresholdPercentage: uint32(sp.ConfirmationThreshold),
			ChunkLength:                     10,
		}
		quorumNumbers[i] = sp.QuorumID
		quorumPercentSigned[i] = confirmedMetadata.ConfirmationInfo.QuorumResults[sp.QuorumID].PercentSigned
		quorumIndexes[i] = byte(i)
	}
	assert.Equal(t, reply.GetInfo().GetBlobHeader().GetBlobQuorumParams(), actualBlobQuorumParams)

	assert.Equal(t, reply.GetInfo().GetBlobVerificationProof().GetBatchId(), confirmedMetadata.ConfirmationInfo.BatchID)
	assert.Equal(t, reply.GetInfo().GetBlobVerificationProof().GetBlobIndex(), confirmedMetadata.ConfirmationInfo.BlobIndex)
	assert.Equal(t, reply.GetInfo().GetBlobVerificationProof().GetBatchMetadata(), &pb.BatchMetadata{
		BatchHeader: &pb.BatchHeader{
			BatchRoot:               confirmedMetadata.ConfirmationInfo.BatchRoot,
			QuorumNumbers:           quorumNumbers,
			QuorumSignedPercentages: quorumPercentSigned,
			ReferenceBlockNumber:    confirmedMetadata.ConfirmationInfo.ReferenceBlockNumber,
		},
		SignatoryRecordHash:     confirmedMetadata.ConfirmationInfo.SignatoryRecordHash[:],
		Fee:                     confirmedMetadata.ConfirmationInfo.Fee,
		ConfirmationBlockNumber: confirmedMetadata.ConfirmationInfo.ConfirmationBlockNumber,
		BatchHeaderHash:         confirmedMetadata.ConfirmationInfo.BatchHeaderHash[:],
	})
	assert.Equal(t, reply.GetInfo().GetBlobVerificationProof().GetInclusionProof(), confirmedMetadata.ConfirmationInfo.BlobInclusionProof)
	assert.Equal(t, reply.GetInfo().GetBlobVerificationProof().GetQuorumIndexes(), quorumIndexes)
}

func TestRetrieveBlob(t *testing.T) {
	// Create random data
	data := make([]byte, 1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	// Disperse the random data
	status, blobSize, requestID := disperseBlob(t, dispersalServer, data)
	assert.Equal(t, status, pb.BlobStatus_PROCESSING)
	assert.NotNil(t, requestID)

	reply, err := dispersalServer.GetBlobStatus(context.Background(), &pb.BlobStatusRequest{
		RequestId: requestID,
	})
	assert.NoError(t, err)
	assert.Equal(t, reply.GetStatus(), pb.BlobStatus_PROCESSING)

	// Simulate blob confirmation so that we can retrieve the blob
	securityParams := []*core.SecurityParam{
		{
			QuorumID:              0,
			AdversaryThreshold:    80,
			ConfirmationThreshold: 100,
		},
		{
			QuorumID:              1,
			AdversaryThreshold:    80,
			ConfirmationThreshold: 100,
		},
	}
	_ = simulateBlobConfirmation(t, requestID, blobSize, securityParams, 1)

	reply, err = dispersalServer.GetBlobStatus(context.Background(), &pb.BlobStatusRequest{
		RequestId: requestID,
	})
	assert.NoError(t, err)
	assert.Equal(t, reply.GetStatus(), pb.BlobStatus_CONFIRMED)

	// Retrieve the blob and compare it with the original data
	retrieveData, err := retrieveBlob(t, dispersalServer, 1)
	assert.NoError(t, err)

	assert.Equal(t, data, retrieveData)
}

func TestRetrieveBlobFailsWhenBlobNotConfirmed(t *testing.T) {
	// Create random data
	data := make([]byte, 1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	// Disperse the random data
	status, _, requestID := disperseBlob(t, dispersalServer, data)
	assert.Equal(t, status, pb.BlobStatus_PROCESSING)
	assert.NotNil(t, requestID)

	reply, err := dispersalServer.GetBlobStatus(context.Background(), &pb.BlobStatusRequest{
		RequestId: requestID,
	})
	assert.NoError(t, err)
	assert.Equal(t, reply.GetStatus(), pb.BlobStatus_PROCESSING)

	// Try to retrieve the blob before it is confirmed
	_, err = retrieveBlob(t, dispersalServer, 2)
	assert.Error(t, err)
}

func TestDisperseBlobWithExceedSizeLimit(t *testing.T) {
	data := make([]byte, 2*1024*1024+10)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	quorumParams := []*core.SecurityParam{
		{QuorumID: 0, AdversaryThreshold: 80, ConfirmationThreshold: 100},
		{QuorumID: 1, AdversaryThreshold: 80, ConfirmationThreshold: 100},
	}
	transactor.On("GetQuorumSecurityParams", tmock.Anything).Return(quorumParams, nil)
	transactor.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{}, nil)

	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:          data,
		QuorumNumbers: []uint32{0, 1},
	})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "blob size cannot exceed 2 MiB")
}

func setup(m *testing.M) {

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

	err := deploy.DeployResources(dockertestPool, localStackPort, metadataTableName, bucketTableName)
	if err != nil {
		teardown()
		panic("failed to deploy AWS resources")
	}

	dispersalServer = newTestServer(m)
}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func newTestServer(m *testing.M) *apiserver.DispersalServer {
	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	if err != nil {
		panic("failed to create a new logger")
	}

	bucketName := "test-eigenda-blobstore"
	awsConfig := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}
	s3Client, err := s3.NewClient(context.Background(), awsConfig, logger)
	if err != nil {
		panic("failed to create s3 client")
	}
	dynamoClient, err := dynamodb.NewClient(awsConfig, logger)
	if err != nil {
		panic("failed to create dynamoDB client")
	}
	blobMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, metadataTableName, time.Hour)

	globalParams := common.GlobalRateParams{
		CountFailed: false,
		BucketSizes: []time.Duration{3 * time.Second},
		Multipliers: []float32{1},
	}
	bucketStore, err := store.NewLocalParamStore[common.RateBucketParams](1000)
	if err != nil {
		panic("failed to create bucket store")
	}
	ratelimiter := ratelimit.NewRateLimiter(globalParams, bucketStore, logger)

	rateConfig := apiserver.RateConfig{
		QuorumRateInfos: map[core.QuorumID]apiserver.QuorumRateInfo{
			0: {
				PerUserUnauthThroughput: 20 * 1024,
				TotalUnauthThroughput:   1048576,
				PerUserUnauthBlobRate:   3 * 1e6,
				TotalUnauthBlobRate:     100 * 1e6,
			},
			1: {
				PerUserUnauthThroughput: 20 * 1024,
				TotalUnauthThroughput:   1048576,
				PerUserUnauthBlobRate:   3 * 1e6,
				TotalUnauthBlobRate:     100 * 1e6,
			},
		},
		ClientIPHeader: "",
		Allowlist: apiserver.Allowlist{
			"1.2.3.4": map[uint8]apiserver.PerUserRateInfo{
				0: {
					Throughput: 100 * 1024,
					BlobRate:   5 * 1e6,
				},
				1: {
					Throughput: 1024 * 1024,
					BlobRate:   5 * 1e6,
				},
			},
			"0x1aa8226f6d354380dDE75eE6B634875c4203e522": map[uint8]apiserver.PerUserRateInfo{
				0: {
					Throughput: 100 * 1024,
					BlobRate:   5 * 1e6,
				},
				1: {
					Throughput: 1024 * 1024,
					BlobRate:   5 * 1e6,
				},
			},
		},
	}

	queue = blobstore.NewSharedStorage(bucketName, s3Client, blobMetadataStore, logger)
	transactor = &mock.MockTransactor{}
	transactor.On("GetCurrentBlockNumber").Return(uint32(100), nil)
	transactor.On("GetQuorumCount").Return(uint8(2), nil)

	return apiserver.NewDispersalServer(disperser.ServerConfig{
		GrpcPort: "51001",
	}, queue, transactor, logger, disperser.NewMetrics("9001", logger), ratelimiter, rateConfig)
}

func disperseBlob(t *testing.T, server *apiserver.DispersalServer, data []byte) (pb.BlobStatus, uint, []byte) {
	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	quorumParams := []*core.SecurityParam{
		{QuorumID: 0, AdversaryThreshold: 80, ConfirmationThreshold: 100},
		{QuorumID: 1, AdversaryThreshold: 80, ConfirmationThreshold: 100},
	}
	transactor.On("GetQuorumSecurityParams", tmock.Anything).Return(quorumParams, nil)
	transactor.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{}, nil)

	reply, err := server.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:          data,
		QuorumNumbers: []uint32{0, 1},
	})
	assert.NoError(t, err)
	return reply.GetResult(), uint(len(data)), reply.GetRequestId()
}

func retrieveBlob(t *testing.T, server *apiserver.DispersalServer, blobIndex uint32) ([]byte, error) {
	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	reply, err := server.RetrieveBlob(ctx, &pb.RetrieveBlobRequest{
		BatchHeaderHash: []byte{1, 2, 3},
		BlobIndex:       blobIndex,
	})
	if err != nil {
		return nil, err
	}

	return reply.GetData(), nil
}

func simulateBlobConfirmation(t *testing.T, requestID []byte, blobSize uint, securityParams []*core.SecurityParam, blobIndex uint32) *disperser.BlobMetadata {
	ctx := context.Background()

	metadataKey, err := disperser.ParseBlobKey(string(requestID))
	assert.NoError(t, err)

	// simulate processing
	err = queue.MarkBlobProcessing(ctx, metadataKey)
	assert.NoError(t, err)

	// simulate blob confirmation
	batchHeaderHash := [32]byte{1, 2, 3}
	requestedAt := uint64(time.Now().Nanosecond())
	var commitX, commitY fp.Element
	_, err = commitX.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	assert.NoError(t, err)

	_, err = commitY.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	assert.NoError(t, err)

	commitment := &encoding.G1Commitment{
		X: commitX,
		Y: commitY,
	}
	dataLength := 32
	batchID := uint32(99)
	batchRoot := []byte("hello")
	referenceBlockNumber := uint32(132)
	confirmationBlockNumber := uint32(150)
	sigRecordHash := [32]byte{0}
	fee := []byte{0}
	inclusionProof := []byte{1, 2, 3, 4, 5}
	quorumResults := make(map[core.QuorumID]*core.QuorumResult, len(securityParams))
	quorumInfos := make([]*core.BlobQuorumInfo, len(securityParams))
	for i, sp := range securityParams {
		quorumResults[sp.QuorumID] = &core.QuorumResult{
			QuorumID:      sp.QuorumID,
			PercentSigned: 100,
		}
		quorumInfos[i] = &core.BlobQuorumInfo{
			SecurityParam: *sp,
			ChunkLength:   10,
		}
	}

	confirmationInfo := &disperser.ConfirmationInfo{
		BatchHeaderHash:      batchHeaderHash,
		BlobIndex:            blobIndex,
		SignatoryRecordHash:  sigRecordHash,
		ReferenceBlockNumber: referenceBlockNumber,
		BatchRoot:            batchRoot,
		BlobInclusionProof:   inclusionProof,
		BlobCommitment: &encoding.BlobCommitments{
			Commitment: commitment,
			Length:     uint(dataLength),
		},
		BatchID:                 batchID,
		ConfirmationTxnHash:     gethcommon.HexToHash("0x123"),
		ConfirmationBlockNumber: confirmationBlockNumber,
		Fee:                     fee,
		QuorumResults:           quorumResults,
		BlobQuorumInfos:         quorumInfos,
	}
	metadata := &disperser.BlobMetadata{
		BlobHash:     metadataKey.BlobHash,
		MetadataHash: metadataKey.MetadataHash,
		BlobStatus:   disperser.Processing,
		Expiry:       0,
		NumRetries:   0,
		RequestMetadata: &disperser.RequestMetadata{
			BlobRequestHeader: core.BlobRequestHeader{
				SecurityParams: securityParams,
			},
			RequestedAt: requestedAt,
			BlobSize:    blobSize,
		},
	}
	updated, err := queue.MarkBlobConfirmed(ctx, metadata, confirmationInfo)
	assert.NoError(t, err)

	return updated
}
