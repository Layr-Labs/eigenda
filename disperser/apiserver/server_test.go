package apiserver_test

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/core/mock"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/core"
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

	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	UUID               = uuid.New()
	metadataTableName  = fmt.Sprintf("test-BlobMetadata-%v", UUID)
	bucketTableName    = fmt.Sprintf("test-BucketStore-%v", UUID)

	deployLocalStack bool
	localStackPort   = "4568"
	allowlistFile    *os.File
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestDisperseBlob(t *testing.T) {
	data := make([]byte, 3*1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)

	status, _, key := disperseBlob(t, dispersalServer, data)
	assert.Equal(t, status, pb.BlobStatus_PROCESSING)
	assert.NotNil(t, key)
}

func TestDisperseBlobAuth(t *testing.T) {

	data1KiB := make([]byte, 1024)
	_, err := rand.Read(data1KiB)
	assert.NoError(t, err)

	data1KiB = codec.ConvertByPaddingEmptyByte(data1KiB)

	// Use an unauthenticated signer
	privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdeb"
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

	errorChan := make(chan error, 10)

	// Should fail with account throughput limit because unauth throughput limit is 20 KiB/s for quorum 0
	simulateClient(t, signer, "0.0.0.0", data1KiB, []uint32{0}, 0, errorChan, false)

	err = <-errorChan
	assert.NoError(t, err)

}

func TestDisperseBlobAuthTimeout(t *testing.T) {

	data1KiB := make([]byte, 1024)
	_, err := rand.Read(data1KiB)
	assert.NoError(t, err)

	data1KiB = codec.ConvertByPaddingEmptyByte(data1KiB)

	// Use an unauthenticated signer
	privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdeb"
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

	errorChan := make(chan error, 10)

	simulateClient(t, signer, "0.0.0.0", data1KiB, []uint32{0}, 2*time.Second, errorChan, false)

	err = <-errorChan
	assert.ErrorContains(t, err, "context deadline exceeded")

	errorChan = make(chan error, 10)
	simulateClient(t, signer, "0.0.0.0", data1KiB, []uint32{0}, 0, errorChan, false)

	err = <-errorChan
	assert.NoError(t, err)

}

func TestDisperseBlobWithRequiredQuorums(t *testing.T) {

	transactor := &mock.MockTransactor{}
	transactor.On("GetCurrentBlockNumber").Return(uint32(100), nil)
	transactor.On("GetQuorumCount").Return(uint8(2), nil)
	quorumParams := []core.SecurityParam{
		{QuorumID: 0, AdversaryThreshold: 80, ConfirmationThreshold: 100},
		{QuorumID: 1, AdversaryThreshold: 80, ConfirmationThreshold: 100},
	}
	transactor.On("GetQuorumSecurityParams", tmock.Anything).Return(quorumParams, nil)

	dispersalServer := newTestServer(transactor)

	data := make([]byte, 1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	transactor.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{0, 1}, nil).Twice()

	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: []uint32{1},
	})
	assert.Error(t, err)

	reply, err := dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: []uint32{},
	})
	assert.NoError(t, err)
	assert.Equal(t, reply.GetResult(), pb.BlobStatus_PROCESSING)
	assert.NotNil(t, reply.GetRequestId())

	transactor.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{0}, nil).Twice()
	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: []uint32{0},
	})
	assert.Error(t, err)

	reply, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: []uint32{1},
	})
	assert.NoError(t, err)
	assert.Equal(t, pb.BlobStatus_PROCESSING, reply.GetResult())
	assert.NotNil(t, reply.GetRequestId())

	transactor.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{}, nil).Once()
	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: []uint32{},
	})
	assert.Error(t, err)
}

func TestDisperseBlobWithInvalidQuorum(t *testing.T) {
	data := make([]byte, 1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: []uint32{2},
	})
	assert.ErrorContains(t, err, "custom_quorum_numbers must be in range [0, 1], but found 2")

	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: []uint32{0, 0},
	})
	assert.ErrorContains(t, err, "custom_quorum_numbers must not contain duplicates")

}

func TestGetBlobStatus(t *testing.T) {
	data := make([]byte, 1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)

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

func TestGetBlobDispersingStatus(t *testing.T) {
	data := make([]byte, 1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)

	status, _, requestID := disperseBlob(t, dispersalServer, data)
	assert.Equal(t, status, pb.BlobStatus_PROCESSING)
	assert.NotNil(t, requestID)
	blobKey, err := disperser.ParseBlobKey(string(requestID))
	assert.NoError(t, err)
	err = queue.MarkBlobDispersing(context.Background(), blobKey)
	assert.NoError(t, err)
	meta, err := queue.GetBlobMetadata(context.Background(), blobKey)
	assert.NoError(t, err)
	assert.Equal(t, meta.BlobStatus, disperser.Dispersing)

	reply, err := dispersalServer.GetBlobStatus(context.Background(), &pb.BlobStatusRequest{
		RequestId: requestID,
	})
	assert.NoError(t, err)
	assert.Equal(t, reply.GetStatus(), pb.BlobStatus_PROCESSING)
}

func TestRetrieveBlob(t *testing.T) {

	for i := 0; i < 3; i++ {
		// Create random data
		data := make([]byte, 1024)
		_, err := rand.Read(data)
		assert.NoError(t, err)

		data = codec.ConvertByPaddingEmptyByte(data)

		// Disperse the random data
		status, blobSize, requestID := disperseBlob(t, dispersalServer, data)
		assert.Equal(t, status, pb.BlobStatus_PROCESSING)
		assert.NotNil(t, requestID)

		reply, err := dispersalServer.GetBlobStatus(context.Background(), &pb.BlobStatusRequest{
			RequestId: requestID,
		})
		assert.NoError(t, err)
		assert.Equal(t, reply.GetStatus(), pb.BlobStatus_PROCESSING)

		fmt.Println("requestID", requestID)

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
		retrieveData, err := retrieveBlob(dispersalServer, requestID, 1)
		assert.NoError(t, err)

		assert.Equal(t, data, retrieveData)
	}

}

func TestRetrieveBlobFailsWhenBlobNotConfirmed(t *testing.T) {
	// Create random data
	data := make([]byte, 1024)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)

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
	_, err = retrieveBlob(dispersalServer, requestID, 2)
	assert.NotNil(t, err)
	assert.Equal(t, "rpc error: code = NotFound desc = no metadata found for the given batch header hash and blob index", err.Error())

}

func TestDisperseBlobWithExceedSizeLimit(t *testing.T) {
	data := make([]byte, 2*1024*1024+10)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: []uint32{0, 1},
	})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "rpc error: code = InvalidArgument desc = blob size cannot exceed 2 MiB")
}

func TestParseAllowlist(t *testing.T) {
	fs := flag.NewFlagSet("disperser", flag.ContinueOnError)
	allowlistFileFlag := apiserver.AllowlistFileFlag("disperser")
	allowlistFileFlag.Apply(fs)

	overwriteFile(t, allowlistFile, `
[
  {
    "name": "eigenlabs",
    "account": "0.1.2.3",
    "quorumID": 0,
    "blobRate": 0.01,
    "byteRate": 1024
  },
  {
    "name": "eigenlabs",
    "account": "0.1.2.3",
    "quorumID": 1,
    "blobRate": 1,
    "byteRate": 1048576
  },
  {
    "name": "foo",
    "account": "5.5.5.5",
    "quorumID": 1,
    "blobRate": 0.1,
    "byteRate": 4092
  }
]
	`)
	err := fs.Parse([]string{"--auth.allowlist-file", allowlistFile.Name()})
	assert.NoError(t, err)
	c := cli.NewContext(nil, fs, nil)
	rateConfig, err := apiserver.ReadCLIConfig(c)
	assert.NoError(t, err)

	assert.Contains(t, rateConfig.Allowlist, "0.1.2.3")
	assert.Contains(t, rateConfig.Allowlist, "5.5.5.5")
	assert.Contains(t, rateConfig.Allowlist["0.1.2.3"], uint8(0))
	assert.Contains(t, rateConfig.Allowlist["0.1.2.3"], uint8(1))
	assert.Contains(t, rateConfig.Allowlist["5.5.5.5"], uint8(1))
	assert.NotContains(t, rateConfig.Allowlist["5.5.5.5"], uint8(0))
	assert.Equal(t, rateConfig.Allowlist["0.1.2.3"][0].Name, "eigenlabs")
	assert.Equal(t, rateConfig.Allowlist["0.1.2.3"][0].BlobRate, uint32(0.01*1e6))
	assert.Equal(t, rateConfig.Allowlist["0.1.2.3"][0].Throughput, uint32(1024))
	assert.Equal(t, rateConfig.Allowlist["0.1.2.3"][1].Name, "eigenlabs")
	assert.Equal(t, rateConfig.Allowlist["0.1.2.3"][1].BlobRate, uint32(1e6))
	assert.Equal(t, rateConfig.Allowlist["0.1.2.3"][1].Throughput, uint32(1048576))
	assert.Equal(t, rateConfig.Allowlist["5.5.5.5"][1].Name, "foo")
	assert.Equal(t, rateConfig.Allowlist["5.5.5.5"][1].BlobRate, uint32(0.1*1e6))
	assert.Equal(t, rateConfig.Allowlist["5.5.5.5"][1].Throughput, uint32(4092))
}

func TestLoadAllowlistFromFile(t *testing.T) {
	overwriteFile(t, allowlistFile, `
[
  {
    "name": "eigenlabs",
    "account": "0.1.2.3",
    "quorumID": 0,
    "blobRate": 0.01,
    "byteRate": 1024
  },
  {
    "name": "eigenlabs",
    "account": "0.1.2.3",
    "quorumID": 1,
    "blobRate": 1,
    "byteRate": 1048576
  },
  {
    "name": "foo",
    "account": "5.5.5.5",
    "quorumID": 1,
    "blobRate": 0.1,
    "byteRate": 4092
  }
]
	`)
	dispersalServer.LoadAllowlist()
	al := dispersalServer.GetRateConfig().Allowlist
	assert.Contains(t, al, "0.1.2.3")
	assert.Contains(t, al, "5.5.5.5")
	assert.Contains(t, al["0.1.2.3"], uint8(0))
	assert.Contains(t, al["0.1.2.3"], uint8(1))
	assert.Contains(t, al["5.5.5.5"], uint8(1))
	assert.NotContains(t, al["5.5.5.5"], uint8(0))
	assert.Equal(t, al["0.1.2.3"][0].Name, "eigenlabs")
	assert.Equal(t, al["0.1.2.3"][0].BlobRate, uint32(0.01*1e6))
	assert.Equal(t, al["0.1.2.3"][0].Throughput, uint32(1024))
	assert.Equal(t, al["0.1.2.3"][1].Name, "eigenlabs")
	assert.Equal(t, al["0.1.2.3"][1].BlobRate, uint32(1e6))
	assert.Equal(t, al["0.1.2.3"][1].Throughput, uint32(1048576))
	assert.Equal(t, al["5.5.5.5"][1].Name, "foo")
	assert.Equal(t, al["5.5.5.5"][1].BlobRate, uint32(0.1*1e6))
	assert.Equal(t, al["5.5.5.5"][1].Throughput, uint32(4092))

	overwriteFile(t, allowlistFile, `
[
  {
    "name": "hello",
    "account": "0.0.0.0",
    "quorumID": 0,
    "blobRate": 0.1,
    "byteRate": 100
  },
  {
    "name": "world",
    "account": "7.7.7.7",
    "quorumID": 1,
    "blobRate": 1,
    "byteRate": 1234
  }
]
	`)
	dispersalServer.LoadAllowlist()
	al = dispersalServer.GetRateConfig().Allowlist
	assert.NotContains(t, al, "0.1.2.3")
	assert.NotContains(t, al, "5.5.5.5")
	assert.Contains(t, al, "0.0.0.0")
	assert.Contains(t, al, "7.7.7.7")
	assert.Contains(t, al["0.0.0.0"], uint8(0))
	assert.Equal(t, al["0.0.0.0"][0].Name, "hello")
	assert.Equal(t, al["0.0.0.0"][0].BlobRate, uint32(0.1*1e6))
	assert.Equal(t, al["0.0.0.0"][0].Throughput, uint32(100))

	assert.Contains(t, al["7.7.7.7"], uint8(1))
	assert.Equal(t, al["7.7.7.7"][1].Name, "world")
	assert.Equal(t, al["7.7.7.7"][1].BlobRate, uint32(1*1e6))
	assert.Equal(t, al["7.7.7.7"][1].Throughput, uint32(1234))
}

func overwriteFile(t *testing.T, f *os.File, content string) {
	err := f.Truncate(0)
	assert.NoError(t, err)
	_, err = f.Seek(0, 0)
	assert.NoError(t, err)
	_, err = f.WriteString(content)
	assert.NoError(t, err)
}

func setup() {
	var err error
	allowlistFile, err = os.CreateTemp("", "allowlist.*.json")
	if err != nil {
		panic("failed to create allowlist file")
	}

	deployLocalStack = !(os.Getenv("DEPLOY_LOCALSTACK") == "false")
	if !deployLocalStack {
		localStackPort = os.Getenv("LOCALSTACK_PORT")
	}

	if deployLocalStack {
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localStackPort)
		if err != nil {
			teardown()
			panic("failed to start localstack container")
		}

	}

	err = deploy.DeployResources(dockertestPool, localStackPort, metadataTableName, bucketTableName)
	if err != nil {
		teardown()
		panic("failed to deploy AWS resources")
	}

	transactor := &mock.MockTransactor{}
	transactor.On("GetCurrentBlockNumber").Return(uint32(100), nil)
	transactor.On("GetQuorumCount").Return(uint8(2), nil)
	quorumParams := []core.SecurityParam{
		{QuorumID: 0, AdversaryThreshold: 80, ConfirmationThreshold: 100},
		{QuorumID: 1, AdversaryThreshold: 80, ConfirmationThreshold: 100},
	}
	transactor.On("GetQuorumSecurityParams", tmock.Anything).Return(quorumParams, nil)
	transactor.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{}, nil)

	dispersalServer = newTestServer(transactor)
}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
	if allowlistFile != nil {
		_ = os.Remove(allowlistFile.Name())
	}
}

func newTestServer(transactor core.Transactor) *apiserver.DispersalServer {
	logger := logging.NewNoopLogger()

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
	ratelimiter := ratelimit.NewRateLimiter(prometheus.NewRegistry(), globalParams, bucketStore, logger)

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
					Name:       "eigenlabs",
					Throughput: 100 * 1024,
					BlobRate:   5 * 1e6,
				},
				1: {
					Name:       "eigenlabs",
					Throughput: 1024 * 1024,
					BlobRate:   5 * 1e6,
				},
			},
			"0x1aa8226f6d354380dDE75eE6B634875c4203e522": map[uint8]apiserver.PerUserRateInfo{
				0: {
					Name:       "eigenlabs",
					Throughput: 100 * 1024,
					BlobRate:   5 * 1e6,
				},
				1: {
					Name:       "eigenlabs",
					Throughput: 1024 * 1024,
					BlobRate:   5 * 1e6,
				},
			},
		},
		RetrievalBlobRate:   3 * 1e6,
		RetrievalThroughput: 20 * 1024,

		AllowlistFile:            allowlistFile.Name(),
		AllowlistRefreshInterval: 10 * time.Minute,
	}

	queue = blobstore.NewSharedStorage(bucketName, s3Client, blobMetadataStore, logger)

	return apiserver.NewDispersalServer(disperser.ServerConfig{
		GrpcPort:    "51001",
		GrpcTimeout: 1 * time.Second,
	}, queue, transactor, logger, disperser.NewMetrics(prometheus.NewRegistry(), "9001", logger), ratelimiter, rateConfig)
}

func disperseBlob(t *testing.T, server *apiserver.DispersalServer, data []byte) (pb.BlobStatus, uint, []byte) {
	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	reply, err := server.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: []uint32{0, 1},
	})
	assert.NoError(t, err)
	return reply.GetResult(), uint(len(data)), reply.GetRequestId()
}

func retrieveBlob(server *apiserver.DispersalServer, requestID []byte, blobIndex uint32) ([]byte, error) {
	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	batchHeaderHash := crypto.Keccak256(requestID)
	reply, err := server.RetrieveBlob(ctx, &pb.RetrieveBlobRequest{
		BatchHeaderHash: batchHeaderHash,
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
	batchHeaderHash := crypto.Keccak256Hash(requestID)

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
