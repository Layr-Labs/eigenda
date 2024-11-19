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
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"google.golang.org/grpc/peer"

	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common"
	pbcommonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	pbv2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/stretchr/testify/assert"
	tmock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
	commitmentProto, err := commitments.ToProtobuf()
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
	blobHeader, err := corev2.BlobHeaderFromProtobuf(blobHeaderProto)
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
	commitmentProto, err := commitments.ToProtobuf()
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
	ctx := peer.NewContext(context.Background(), c.Peer)

	blobHeader := &corev2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: mockCommitment,
		QuorumNumbers:   []core.QuorumID{0},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         "0x1234",
			BinIndex:          0,
			CumulativePayment: big.NewInt(532),
		},
	}
	blobKey, err := blobHeader.BlobKey()
	require.NoError(t, err)
	now := time.Now()
	metadata := &dispv2.BlobMetadata{
		BlobHeader: blobHeader,
		BlobStatus: dispv2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err = c.BlobMetadataStore.PutBlobMetadata(ctx, metadata)
	require.NoError(t, err)
	blobCert := &corev2.BlobCertificate{
		BlobHeader: blobHeader,
		RelayKeys:  []corev2.RelayKey{0, 1, 2},
	}
	err = c.BlobMetadataStore.PutBlobCertificate(ctx, blobCert, nil)
	require.NoError(t, err)

	// Non ceritified blob status
	status, err := c.DispersalServerV2.GetBlobStatus(ctx, &pbv2.BlobStatusRequest{
		BlobKey: blobKey[:],
	})
	require.NoError(t, err)
	require.Equal(t, pbv2.BlobStatus_QUEUED, status.Status)
	err = c.BlobMetadataStore.UpdateBlobStatus(ctx, blobKey, dispv2.Encoded)
	require.NoError(t, err)
	status, err = c.DispersalServerV2.GetBlobStatus(ctx, &pbv2.BlobStatusRequest{
		BlobKey: blobKey[:],
	})
	require.NoError(t, err)
	require.Equal(t, pbv2.BlobStatus_ENCODED, status.Status)

	// Certified blob status
	err = c.BlobMetadataStore.UpdateBlobStatus(ctx, blobKey, dispv2.Certified)
	require.NoError(t, err)
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 100,
	}
	err = c.BlobMetadataStore.PutBatchHeader(ctx, batchHeader)
	require.NoError(t, err)
	verificationInfo0 := &corev2.BlobVerificationInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey,
		BlobIndex:      123,
		InclusionProof: []byte("inclusion proof"),
	}
	err = c.BlobMetadataStore.PutBlobVerificationInfo(ctx, verificationInfo0)
	require.NoError(t, err)

	attestation := &corev2.Attestation{
		BatchHeader: batchHeader,
		NonSignerPubKeys: []*core.G1Point{
			core.NewG1Point(big.NewInt(1), big.NewInt(2)),
			core.NewG1Point(big.NewInt(3), big.NewInt(4)),
		},
		APKG2: &core.G2Point{
			G2Affine: &bn254.G2Affine{
				X: mockCommitment.LengthCommitment.X,
				Y: mockCommitment.LengthCommitment.Y,
			},
		},
		Sigma: &core.Signature{
			G1Point: core.NewG1Point(big.NewInt(5), big.NewInt(6)),
		},
	}
	err = c.BlobMetadataStore.PutAttestation(ctx, attestation)
	require.NoError(t, err)

	reply, err := c.DispersalServerV2.GetBlobStatus(context.Background(), &pbv2.BlobStatusRequest{
		BlobKey: blobKey[:],
	})
	require.NoError(t, err)
	require.Equal(t, pbv2.BlobStatus_CERTIFIED, reply.GetStatus())
	blobHeaderProto, err := blobHeader.ToProtobuf()
	require.NoError(t, err)
	blobCertProto, err := blobCert.ToProtobuf()
	require.NoError(t, err)
	require.Equal(t, blobHeaderProto, reply.GetBlobVerificationInfo().GetBlobCertificate().GetBlobHeader())
	require.Equal(t, blobCertProto.Relays, reply.GetBlobVerificationInfo().GetBlobCertificate().GetRelays())
	require.Equal(t, verificationInfo0.BlobIndex, reply.GetBlobVerificationInfo().GetBlobIndex())
	require.Equal(t, verificationInfo0.InclusionProof, reply.GetBlobVerificationInfo().GetInclusionProof())
	require.Equal(t, batchHeader.BatchRoot[:], reply.GetSignedBatch().GetHeader().BatchRoot)
	require.Equal(t, batchHeader.ReferenceBlockNumber, reply.GetSignedBatch().GetHeader().ReferenceBlockNumber)
	attestationProto := attestation.ToProtobuf()
	require.Equal(t, attestationProto, reply.GetSignedBatch().GetAttestation())
}

func TestV2GetBlobCommitment(t *testing.T) {
	c := newTestServerV2(t)
	data := make([]byte, 50)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)
	commit, err := prover.GetCommitments(data)
	require.NoError(t, err)
	reply, err := c.DispersalServerV2.GetBlobCommitment(context.Background(), &pbv2.BlobCommitmentRequest{
		Data: data,
	})
	require.NoError(t, err)
	commitment, err := new(encoding.G1Commitment).Deserialize(reply.BlobCommitment.Commitment)
	require.NoError(t, err)
	assert.Equal(t, commit.Commitment, commitment)
	lengthCommitment, err := new(encoding.G2Commitment).Deserialize(reply.BlobCommitment.LengthCommitment)
	require.NoError(t, err)
	assert.Equal(t, commit.LengthCommitment, lengthCommitment)
	lengthProof, err := new(encoding.G2Commitment).Deserialize(reply.BlobCommitment.LengthProof)
	require.NoError(t, err)
	assert.Equal(t, commit.LengthProof, lengthProof)
	assert.Equal(t, uint32(commit.Length), reply.BlobCommitment.Length)
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
	}, rateConfig, blobStore, blobMetadataStore, chainReader, nil, auth.NewAuthenticator(), prover, 100, time.Hour, logger)

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
