package apiserver_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc/peer"

	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common"
	pbcommonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	pbv2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
)

type testComponents struct {
	DispersalServerV2 *apiserver.DispersalServerV2
	BlobStore         *blobstore.BlobStore
	BlobMetadataStore *blobstore.BlobMetadataStore
	ChainReader       *mock.MockWriter
	Signer            *auth.LocalBlobRequestSigner
	Peer              *peer.Peer
}

func TestV2DisperseBlob(t *testing.T) {
	c := newTestServerV2(t)
	ctx := peer.NewContext(context.Background(), c.Peer)
	data := make([]byte, 50)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)
	commitments, err := prover.GetCommitments(data)
	assert.NoError(t, err)
	accountID, err := c.Signer.GetAccountID()
	assert.NoError(t, err)
	commitmentProto, err := commitments.ToProfobuf()
	assert.NoError(t, err)
	blobHeaderProto := &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommon.PaymentHeader{
			AccountId:         accountID,
			BinIndex:          5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	blobHeader, err := corev2.NewBlobHeader(blobHeaderProto)
	assert.NoError(t, err)
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)
	sig, err := signer.SignBlobRequest(blobHeader)
	assert.NoError(t, err)
	blobHeader.Signature = sig
	blobHeaderProto.Signature = sig

	now := time.Now()
	reply, err := c.DispersalServerV2.DisperseBlob(ctx, &pbv2.DisperseBlobRequest{
		Data:       data,
		BlobHeader: blobHeaderProto,
	})
	assert.NoError(t, err)

	blobKey, err := blobHeader.BlobKey()
	assert.NoError(t, err)
	assert.Equal(t, pbv2.BlobStatus_QUEUED, reply.Result)
	assert.Equal(t, blobKey[:], reply.BlobKey)

	// Check if the blob is stored
	storedData, err := c.BlobStore.GetBlob(ctx, blobKey)
	assert.NoError(t, err)
	assert.Equal(t, data, storedData)

	// Check if the blob metadata is stored
	blobMetadata, err := c.BlobMetadataStore.GetBlobMetadata(ctx, blobKey)
	assert.NoError(t, err)
	assert.Equal(t, dispv2.Queued, blobMetadata.BlobStatus)
	assert.Equal(t, blobHeader, blobMetadata.BlobHeader)
	assert.Equal(t, uint64(len(data)), blobMetadata.BlobSize)
	assert.Equal(t, uint(0), blobMetadata.NumRetries)
	assert.Greater(t, blobMetadata.Expiry, uint64(now.Unix()))
	assert.Greater(t, blobMetadata.RequestedAt, uint64(now.UnixNano()))
	assert.Equal(t, blobMetadata.RequestedAt, blobMetadata.UpdatedAt)
}

func TestV2DisperseBlobRequestValidation(t *testing.T) {
	c := newTestServerV2(t)
	data := make([]byte, 50)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)
	commitments, err := prover.GetCommitments(data)
	assert.NoError(t, err)
	accountID, err := c.Signer.GetAccountID()
	assert.NoError(t, err)
	// request with no blob commitments
	invalidReqProto := &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		PaymentHeader: &pbcommon.PaymentHeader{
			AccountId:         accountID,
			BinIndex:          5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Data:       data,
		BlobHeader: invalidReqProto,
	})
	assert.ErrorContains(t, err, "blob header must contain commitments")
	commitmentProto, err := commitments.ToProfobuf()
	assert.NoError(t, err)

	// request with too many quorums
	invalidReqProto = &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1, 2, 3},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommon.PaymentHeader{
			AccountId:         accountID,
			BinIndex:          5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Data:       data,
		BlobHeader: invalidReqProto,
	})
	assert.ErrorContains(t, err, "too many quorum numbers specified")

	// request without required quorums
	invalidReqProto = &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{2},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommon.PaymentHeader{
			AccountId:         accountID,
			BinIndex:          5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Data:       data,
		BlobHeader: invalidReqProto,
	})
	assert.ErrorContains(t, err, "request must contain at least one required quorum")

	// request with invalid blob version
	invalidReqProto = &pbcommonv2.BlobHeader{
		Version:       2,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommon.PaymentHeader{
			AccountId:         accountID,
			BinIndex:          5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Data:       data,
		BlobHeader: invalidReqProto,
	})
	assert.ErrorContains(t, err, "invalid blob version 2")

	// request with invalid signature
	invalidReqProto = &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommon.PaymentHeader{
			AccountId:         accountID,
			BinIndex:          5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
		Signature: []byte{1, 2, 3},
	}
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Data:       data,
		BlobHeader: invalidReqProto,
	})
	assert.ErrorContains(t, err, "authentication failed")
}

func TestV2GetBlobStatus(t *testing.T) {
	c := newTestServerV2(t)
	_, err := c.DispersalServerV2.GetBlobStatus(context.Background(), &pbv2.BlobStatusRequest{
		BlobKey: []byte{1},
	})
	assert.ErrorContains(t, err, "not implemented")
}

func TestV2GetBlobCommitment(t *testing.T) {
	c := newTestServerV2(t)
	_, err := c.DispersalServerV2.GetBlobCommitment(context.Background(), &pbv2.BlobCommitmentRequest{
		Data: []byte{1},
	})
	assert.ErrorContains(t, err, "not implemented")
}

func newTestServerV2(t *testing.T) *testComponents {
	logger := logging.NewNoopLogger()
	// logger, err := common.NewLogger(common.DefaultLoggerConfig())
	// if err != nil {
	// 	panic("failed to create logger")
	// }

	awsConfig := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}
	s3Client, err := s3.NewClient(context.Background(), awsConfig, logger)
	assert.NoError(t, err)
	dynamoClient, err := dynamodb.NewClient(awsConfig, logger)
	assert.NoError(t, err)
	blobMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, v2MetadataTableName)
	blobStore := blobstore.NewBlobStore(s3BucketName, s3Client, logger)
	chainReader := &mock.MockWriter{}
	rateConfig := apiserver.RateConfig{}

	chainReader.On("GetCurrentBlockNumber").Return(uint32(100), nil)
	chainReader.On("GetQuorumCount").Return(uint8(2), nil)
	chainReader.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{0, 1}, nil)
	chainReader.On("GetBlockStaleMeasure", tmock.Anything).Return(uint32(10), nil)
	chainReader.On("GetStoreDurationBlocks", tmock.Anything).Return(uint32(100), nil)

	s := apiserver.NewDispersalServerV2(disperser.ServerConfig{
		GrpcPort:    "51002",
		GrpcTimeout: 1 * time.Second,
	}, rateConfig, blobStore, blobMetadataStore, chainReader, nil, auth.NewAuthenticator(), 100, time.Hour, logger)

	err = s.RefreshOnchainState(context.Background())
	assert.NoError(t, err)
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)
	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 51001,
		},
	}

	return &testComponents{
		DispersalServerV2: s,
		BlobStore:         blobStore,
		BlobMetadataStore: blobMetadataStore,
		ChainReader:       chainReader,
		Signer:            signer,
		Peer:              p,
	}
}
