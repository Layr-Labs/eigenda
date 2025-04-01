package apiserver_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/testutils/random"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	auth "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/mock"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	dispv2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"google.golang.org/grpc/peer"

	pbcommonv2 "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	pbv2 "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	gethcommon "github.com/ethereum/go-ethereum/common"
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
	commitments, err := prover.GetCommitmentsForPaddedLength(data)
	assert.NoError(t, err)
	accountID, err := c.Signer.GetAccountID()
	assert.NoError(t, err)
	commitmentProto, err := commitments.ToProtobuf()
	assert.NoError(t, err)
	blobHeaderProto := &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	blobHeader, err := corev2.BlobHeaderFromProtobuf(blobHeaderProto)
	fmt.Println("blobHeader", blobHeader)
	assert.NoError(t, err)
	signer, err := auth.NewLocalBlobRequestSigner(privateKeyHex)
	assert.NoError(t, err)
	sig, err := signer.SignBlobRequest(blobHeader)
	assert.NoError(t, err)

	now := time.Now()
	reply, err := c.DispersalServerV2.DisperseBlob(ctx, &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  sig,
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

	// Try dispersing the same blob; blob key check will fail if the blob is already stored
	reply, err = c.DispersalServerV2.DisperseBlob(ctx, &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  sig,
		BlobHeader: blobHeaderProto,
	})
	assert.Nil(t, reply)
	assert.ErrorContains(t, err, "blob already exists")

	data2 := make([]byte, 50)
	_, err = rand.Read(data)
	assert.NoError(t, err)

	data2 = codec.ConvertByPaddingEmptyByte(data2)
	commitments, err = prover.GetCommitmentsForPaddedLength(data2)
	assert.NoError(t, err)
	commitmentProto, err = commitments.ToProtobuf()
	assert.NoError(t, err)
	blobHeaderProto2 := &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	blobHeader2, err := corev2.BlobHeaderFromProtobuf(blobHeaderProto2)
	assert.NoError(t, err)
	sig2, err := signer.SignBlobRequest(blobHeader2)
	assert.NoError(t, err)

	reply, err = c.DispersalServerV2.DisperseBlob(ctx, &pbv2.DisperseBlobRequest{
		Blob:       data2,
		Signature:  sig2,
		BlobHeader: blobHeaderProto2,
	})
	assert.Nil(t, reply)
	assert.ErrorContains(t, err, "failed to update cumulative payment: insufficient cumulative payment increment")
}

func TestV2DisperseBlobRequestValidation(t *testing.T) {
	c := newTestServerV2(t)
	data := make([]byte, 50)
	_, err := rand.Read(data)
	assert.NoError(t, err)
	signer, err := auth.NewLocalBlobRequestSigner(privateKeyHex)
	assert.NoError(t, err)
	data = codec.ConvertByPaddingEmptyByte(data)
	commitments, err := prover.GetCommitmentsForPaddedLength(data)
	assert.NoError(t, err)
	accountID, err := c.Signer.GetAccountID()
	assert.NoError(t, err)
	// request with no blob commitments
	invalidReqProto := &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65},
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
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65},
		BlobHeader: invalidReqProto,
	})
	assert.ErrorContains(t, err, "too many quorum numbers specified")

	// request with invalid quorum
	invalidReqProto = &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{2, 54},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65},
		BlobHeader: invalidReqProto,
	})
	assert.ErrorContains(t, err, "invalid quorum")

	// request with invalid blob version
	invalidReqProto = &pbcommonv2.BlobHeader{
		Version:       2,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65},
		BlobHeader: invalidReqProto,
	})
	assert.ErrorContains(t, err, "invalid blob version 2")

	invalidReqProto = &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	// request with invalid signature
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65},
		BlobHeader: invalidReqProto,
	})
	assert.ErrorContains(t, err, "authentication failed")

	// request with invalid payment metadata
	invalidReqProto = &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         0,
			CumulativePayment: big.NewInt(0).Bytes(),
		},
	}
	blobHeader, err := corev2.BlobHeaderFromProtobuf(invalidReqProto)
	assert.NoError(t, err)
	sig, err := signer.SignBlobRequest(blobHeader)
	assert.NoError(t, err)

	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  sig,
		BlobHeader: invalidReqProto,
	})
	assert.ErrorContains(t, err, "invalid payment metadata")

	// request with invalid commitment
	invalidCommitment := commitmentProto
	invalidCommitment.Length = commitmentProto.Length - 1
	invalidReqProto = &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    invalidCommitment,
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	blobHeader, err = corev2.BlobHeaderFromProtobuf(invalidReqProto)
	assert.NoError(t, err)
	sig, err = signer.SignBlobRequest(blobHeader)
	assert.NoError(t, err)
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  sig,
		BlobHeader: invalidReqProto,
	})
	assert.ErrorContains(t, err, "is less than blob length")

	// request with blob size exceeding the limit
	data = make([]byte, 321)
	_, err = rand.Read(data)
	assert.NoError(t, err)
	data = codec.ConvertByPaddingEmptyByte(data)
	commitments, err = prover.GetCommitmentsForPaddedLength(data)
	assert.NoError(t, err)
	commitmentProto, err = commitments.ToProtobuf()
	assert.NoError(t, err)
	validHeader := &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	blobHeader, err = corev2.BlobHeaderFromProtobuf(validHeader)
	assert.NoError(t, err)
	sig, err = signer.SignBlobRequest(blobHeader)
	assert.NoError(t, err)
	_, err = c.DispersalServerV2.DisperseBlob(context.Background(), &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  sig,
		BlobHeader: validHeader,
	})
	assert.ErrorContains(t, err, "blob size too big")
}

func TestV2GetBlobStatus(t *testing.T) {
	c := newTestServerV2(t)
	ctx := peer.NewContext(context.Background(), c.Peer)

	blobHeader := &corev2.BlobHeader{
		BlobVersion:     0,
		BlobCommitments: mockCommitment,
		QuorumNumbers:   []core.QuorumID{0},
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         gethcommon.HexToAddress("0x1234"),
			Timestamp:         0,
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

	// Queued/Encoded blob status
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

	// Complete blob status
	err = c.BlobMetadataStore.UpdateBlobStatus(ctx, blobKey, dispv2.Complete)
	require.NoError(t, err)
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 100,
	}
	err = c.BlobMetadataStore.PutBatchHeader(ctx, batchHeader)
	require.NoError(t, err)
	inclusionInfo0 := &corev2.BlobInclusionInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey,
		BlobIndex:      123,
		InclusionProof: []byte("inclusion proof"),
	}
	err = c.BlobMetadataStore.PutBlobInclusionInfo(ctx, inclusionInfo0)
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
	require.Equal(t, pbv2.BlobStatus_COMPLETE, reply.GetStatus())
	blobHeaderProto, err := blobHeader.ToProtobuf()
	require.NoError(t, err)
	blobCertProto, err := blobCert.ToProtobuf()
	require.NoError(t, err)
	require.Equal(t, blobHeaderProto, reply.GetBlobInclusionInfo().GetBlobCertificate().GetBlobHeader())
	require.Equal(t, blobCertProto.RelayKeys, reply.GetBlobInclusionInfo().GetBlobCertificate().GetRelayKeys())
	require.Equal(t, inclusionInfo0.BlobIndex, reply.GetBlobInclusionInfo().GetBlobIndex())
	require.Equal(t, inclusionInfo0.InclusionProof, reply.GetBlobInclusionInfo().GetInclusionProof())
	require.Equal(t, batchHeader.BatchRoot[:], reply.GetSignedBatch().GetHeader().BatchRoot)
	require.Equal(t, batchHeader.ReferenceBlockNumber, reply.GetSignedBatch().GetHeader().ReferenceBlockNumber)
	attestationProto, err := attestation.ToProtobuf()
	require.NoError(t, err)
	require.Equal(t, attestationProto, reply.GetSignedBatch().GetAttestation())
}

func TestV2GetBlobCommitment(t *testing.T) {
	c := newTestServerV2(t)
	data := make([]byte, 50)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)
	commit, err := prover.GetCommitmentsForPaddedLength(data)
	require.NoError(t, err)
	reply, err := c.DispersalServerV2.GetBlobCommitment(context.Background(), &pbv2.BlobCommitmentRequest{
		Blob: data,
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
	logger := testutils.GetLogger()
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

	// append test name to each table name for an unique store
	mockState := &mock.MockOnchainPaymentState{}
	mockState.On("RefreshOnchainPaymentState", tmock.Anything).Return(nil).Maybe()
	mockState.On("GetReservationWindow", tmock.Anything).Return(uint64(1), nil)
	mockState.On("GetPricePerSymbol", tmock.Anything).Return(uint64(2), nil)
	mockState.On("GetGlobalSymbolsPerSecond", tmock.Anything).Return(uint64(1009), nil)
	mockState.On("GetGlobalRatePeriodInterval", tmock.Anything).Return(uint64(1), nil)
	mockState.On("GetMinNumSymbols", tmock.Anything).Return(uint64(3), nil)

	now := uint64(time.Now().Unix())
	mockState.On("GetReservedPaymentByAccount", tmock.Anything, tmock.Anything).Return(&core.ReservedPayment{SymbolsPerSecond: 100, StartTimestamp: now + 1200, EndTimestamp: now + 1800, QuorumSplits: []byte{50, 50}, QuorumNumbers: []uint8{0, 1}}, nil)
	mockState.On("GetOnDemandPaymentByAccount", tmock.Anything, tmock.Anything).Return(&core.OnDemandPayment{CumulativePayment: big.NewInt(3864)}, nil)
	mockState.On("GetOnDemandQuorumNumbers", tmock.Anything).Return([]uint8{0, 1}, nil)

	if err := mockState.RefreshOnchainPaymentState(context.Background()); err != nil {
		panic("failed to make initial query to the on-chain state")
	}
	table_names := []string{"reservations_server_" + t.Name(), "ondemand_server_" + t.Name(), "global_server_" + t.Name()}
	err = meterer.CreateReservationTable(awsConfig, table_names[0])
	if err != nil {
		teardown()
		panic("failed to create reservation table")
	}
	err = meterer.CreateOnDemandTable(awsConfig, table_names[1])
	if err != nil {
		teardown()
		panic("failed to create ondemand table")
	}
	err = meterer.CreateGlobalReservationTable(awsConfig, table_names[2])
	if err != nil {
		teardown()
		panic("failed to create global reservation table")
	}

	store, err := meterer.NewOffchainStore(
		awsConfig,
		table_names[0],
		table_names[1],
		table_names[2],
		logger,
	)
	if err != nil {
		teardown()
		panic("failed to create offchain store")
	}
	meterer := meterer.NewMeterer(meterer.Config{}, mockState, store, logger)

	chainReader.On("GetCurrentBlockNumber").Return(uint32(100), nil)
	chainReader.On("GetQuorumCount").Return(uint8(2), nil)
	chainReader.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{0, 1}, nil)
	chainReader.On("GetBlockStaleMeasure", tmock.Anything).Return(uint32(10), nil)
	chainReader.On("GetStoreDurationBlocks", tmock.Anything).Return(uint32(100), nil)
	chainReader.On("GetAllVersionedBlobParams", tmock.Anything).Return(map[v2.BlobVersion]*core.BlobVersionParameters{
		0: {
			NumChunks:       8192,
			CodingRate:      8,
			MaxNumOperators: 3537,
		},
	}, nil)

	s, err := apiserver.NewDispersalServerV2(
		disperser.ServerConfig{
			GrpcPort:    "51002",
			GrpcTimeout: 1 * time.Second,
		},
		blobStore,
		blobMetadataStore,
		chainReader,
		meterer,
		auth.NewBlobRequestAuthenticator(),
		prover,
		10,
		time.Hour,
		logger,
		prometheus.NewRegistry(),
		disperser.MetricsConfig{
			HTTPPort:      "9094",
			EnableMetrics: false,
		},
	)
	assert.NoError(t, err)

	err = s.RefreshOnchainState(context.Background())
	assert.NoError(t, err)
	signer, err := auth.NewLocalBlobRequestSigner(privateKeyHex)
	assert.NoError(t, err)
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

func TestInvalidLength(t *testing.T) {
	c := newTestServerV2(t)
	ctx := peer.NewContext(context.Background(), c.Peer)
	data := make([]byte, 50)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)
	commitments, err := prover.GetCommitmentsForPaddedLength(data)
	assert.NoError(t, err)

	// Length we are commiting to should be a power of 2.
	require.Equal(t, commitments.Length, encoding.NextPowerOf2(commitments.Length))

	// Changing the number of commitments should cause an error before a validity check of the commitments
	commitments.Length += 1

	accountID, err := c.Signer.GetAccountID()
	assert.NoError(t, err)
	commitmentProto, err := commitments.ToProtobuf()
	assert.NoError(t, err)
	blobHeaderProto := &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	blobHeader, err := corev2.BlobHeaderFromProtobuf(blobHeaderProto)
	assert.NoError(t, err)
	signer, err := auth.NewLocalBlobRequestSigner(privateKeyHex)
	assert.NoError(t, err)
	sig, err := signer.SignBlobRequest(blobHeader)
	assert.NoError(t, err)

	_, err = c.DispersalServerV2.DisperseBlob(ctx, &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  sig,
		BlobHeader: blobHeaderProto,
	})

	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid commitment length, must be a power of 2")
}

func TestTooShortCommitment(t *testing.T) {
	rand := random.NewTestRandom()

	c := newTestServerV2(t)
	ctx := peer.NewContext(context.Background(), c.Peer)
	data := rand.VariableBytes(2, 100)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)
	commitments, err := prover.GetCommitmentsForPaddedLength(data)
	assert.NoError(t, err)

	// Length we are commiting to should be a power of 2.
	require.Equal(t, commitments.Length, encoding.NextPowerOf2(commitments.Length))

	// Choose a smaller commitment length than is legal. Make sure it's a power of 2 so that it doesn't
	// fail prior to the commitment length check.
	commitments.Length /= 2

	accountID, err := c.Signer.GetAccountID()
	assert.NoError(t, err)
	commitmentProto, err := commitments.ToProtobuf()
	assert.NoError(t, err)
	blobHeaderProto := &pbcommonv2.BlobHeader{
		Version:       0,
		QuorumNumbers: []uint32{0, 1},
		Commitment:    commitmentProto,
		PaymentHeader: &pbcommonv2.PaymentHeader{
			AccountId:         accountID.Hex(),
			Timestamp:         5,
			CumulativePayment: big.NewInt(100).Bytes(),
		},
	}
	blobHeader, err := corev2.BlobHeaderFromProtobuf(blobHeaderProto)
	assert.NoError(t, err)
	signer, err := auth.NewLocalBlobRequestSigner(privateKeyHex)
	assert.NoError(t, err)
	sig, err := signer.SignBlobRequest(blobHeader)
	assert.NoError(t, err)

	_, err = c.DispersalServerV2.DisperseBlob(ctx, &pbv2.DisperseBlobRequest{
		Blob:       data,
		Signature:  sig,
		BlobHeader: blobHeaderProto,
	})

	require.Error(t, err)
	require.True(t, strings.Contains(err.Error(), "invalid commitment length") ||
		strings.Contains(err.Error(), "is less than blob length"))
}
