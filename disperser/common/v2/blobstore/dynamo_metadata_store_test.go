package blobstore_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestBlobMetadataStoreOperations(t *testing.T) {
	ctx := context.Background()
	blobKey1, blobHeader1 := newBlob(t)
	blobKey2, blobHeader2 := newBlob(t)
	now := time.Now()
	metadata1 := &v2.BlobMetadata{
		BlobHeader: blobHeader1,
		BlobStatus: v2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	metadata2 := &v2.BlobMetadata{
		BlobHeader: blobHeader2,
		BlobStatus: v2.Certified,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err := blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	assert.NoError(t, err)
	err = blobMetadataStore.PutBlobMetadata(ctx, metadata2)
	assert.NoError(t, err)

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey1)
	assert.NoError(t, err)
	assert.Equal(t, metadata1, fetchedMetadata)
	fetchedMetadata, err = blobMetadataStore.GetBlobMetadata(ctx, blobKey2)
	assert.NoError(t, err)
	assert.Equal(t, metadata2, fetchedMetadata)

	queued, err := blobMetadataStore.GetBlobMetadataByStatus(ctx, v2.Queued, 0)
	assert.NoError(t, err)
	assert.Len(t, queued, 1)
	assert.Equal(t, metadata1, queued[0])
	// query to get newer blobs should result in 0 results
	queued, err = blobMetadataStore.GetBlobMetadataByStatus(ctx, v2.Queued, metadata1.UpdatedAt+100)
	assert.NoError(t, err)
	assert.Len(t, queued, 0)

	certified, err := blobMetadataStore.GetBlobMetadataByStatus(ctx, v2.Certified, 0)
	assert.NoError(t, err)
	assert.Len(t, certified, 1)
	assert.Equal(t, metadata2, certified[0])

	queuedCount, err := blobMetadataStore.GetBlobMetadataCountByStatus(ctx, v2.Queued)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), queuedCount)

	// attempt to put metadata with the same key should fail
	err = blobMetadataStore.PutBlobMetadata(ctx, metadata1)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey1.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey2.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		},
	})
}

func TestBlobMetadataStoreGetBlobMetadataByStatusPaginated(t *testing.T) {
	ctx := context.Background()
	numBlobs := 103
	pageSize := 10
	keys := make([]corev2.BlobKey, numBlobs)
	headers := make([]*corev2.BlobHeader, numBlobs)
	metadataList := make([]*v2.BlobMetadata, numBlobs)
	dynamoKeys := make([]commondynamodb.Key, numBlobs)
	expectedCursors := make([]*blobstore.StatusIndexCursor, 0)
	for i := 0; i < numBlobs; i++ {
		blobKey, blobHeader := newBlob(t)
		now := time.Now()
		metadata := &v2.BlobMetadata{
			BlobHeader: blobHeader,
			BlobStatus: v2.Encoded,
			Expiry:     uint64(now.Add(time.Hour).Unix()),
			NumRetries: 0,
			UpdatedAt:  uint64(now.UnixNano()),
		}

		err := blobMetadataStore.PutBlobMetadata(ctx, metadata)
		require.NoError(t, err)
		keys[i] = blobKey
		headers[i] = blobHeader
		dynamoKeys[i] = commondynamodb.Key{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		}
		metadataList[i] = metadata
		if (i+1)%pageSize == 0 {
			expectedCursors = append(expectedCursors, &blobstore.StatusIndexCursor{
				BlobKey:   &blobKey,
				UpdatedAt: metadata.UpdatedAt,
			})
		}
	}

	// Querying blobs in Queued status should return 0 results
	cursor := &blobstore.StatusIndexCursor{
		BlobKey:   nil,
		UpdatedAt: 0,
	}
	metadata, newCursor, err := blobMetadataStore.GetBlobMetadataByStatusPaginated(ctx, v2.Queued, cursor, 10)
	require.NoError(t, err)
	require.Len(t, metadata, 0)
	require.Equal(t, cursor, newCursor)

	// Querying blobs in Encoded status should return results
	cursor = &blobstore.StatusIndexCursor{
		BlobKey:   nil,
		UpdatedAt: 0,
	}
	i := 0
	numIterations := (numBlobs + pageSize - 1) / pageSize
	for i < numIterations {
		metadata, cursor, err = blobMetadataStore.GetBlobMetadataByStatusPaginated(ctx, v2.Encoded, cursor, int32(pageSize))
		require.NoError(t, err)
		if i < len(expectedCursors) {
			require.Len(t, metadata, pageSize)
			require.NotNil(t, cursor)
			require.Equal(t, cursor.BlobKey, expectedCursors[i].BlobKey)
			require.Equal(t, cursor.UpdatedAt, expectedCursors[i].UpdatedAt)
		} else {
			require.Len(t, metadata, numBlobs%pageSize)
			require.Equal(t, cursor.BlobKey, &keys[numBlobs-1])
			require.Equal(t, cursor.UpdatedAt, metadataList[numBlobs-1].UpdatedAt)
		}
		i++
	}
	lastCursor := cursor
	metadata, cursor, err = blobMetadataStore.GetBlobMetadataByStatusPaginated(ctx, v2.Encoded, cursor, int32(pageSize))
	require.NoError(t, err)
	require.Len(t, metadata, 0)
	require.Equal(t, cursor, lastCursor)

	deleteItems(t, dynamoKeys)
}

func TestBlobMetadataStoreCerts(t *testing.T) {
	ctx := context.Background()
	blobKey, blobHeader := newBlob(t)
	blobCert := &corev2.BlobCertificate{
		BlobHeader: blobHeader,
		RelayKeys:  []corev2.RelayKey{0, 2, 4},
	}
	fragmentInfo := &encoding.FragmentInfo{
		TotalChunkSizeBytes: 100,
		FragmentSizeBytes:   1024 * 1024 * 4,
	}
	err := blobMetadataStore.PutBlobCertificate(ctx, blobCert, fragmentInfo)
	assert.NoError(t, err)

	fetchedCert, fetchedFragmentInfo, err := blobMetadataStore.GetBlobCertificate(ctx, blobKey)
	assert.NoError(t, err)
	assert.Equal(t, blobCert, fetchedCert)
	assert.Equal(t, fragmentInfo, fetchedFragmentInfo)

	// blob cert with the same key should fail
	blobCert1 := &corev2.BlobCertificate{
		BlobHeader: blobHeader,
		RelayKeys:  []corev2.RelayKey{0},
	}
	err = blobMetadataStore.PutBlobCertificate(ctx, blobCert1, fragmentInfo)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	// get multiple certs
	numCerts := 100
	keys := make([]corev2.BlobKey, numCerts)
	for i := 0; i < numCerts; i++ {
		blobCert := &corev2.BlobCertificate{
			BlobHeader: &corev2.BlobHeader{
				BlobVersion:     0,
				QuorumNumbers:   []core.QuorumID{0},
				BlobCommitments: mockCommitment,
				PaymentMetadata: core.PaymentMetadata{
					AccountID:         "0x123",
					ReservationPeriod: uint32(i),
					CumulativePayment: big.NewInt(321),
				},
				Signature: []byte("signature"),
			},
			RelayKeys: []corev2.RelayKey{0},
		}
		blobKey, err := blobCert.BlobHeader.BlobKey()
		assert.NoError(t, err)
		keys[i] = blobKey
		err = blobMetadataStore.PutBlobCertificate(ctx, blobCert, fragmentInfo)
		assert.NoError(t, err)
	}

	certs, fragmentInfos, err := blobMetadataStore.GetBlobCertificates(ctx, keys)
	assert.NoError(t, err)
	assert.Len(t, certs, numCerts)
	assert.Len(t, fragmentInfos, numCerts)
	reservationPeriodes := make(map[uint32]struct{})
	for i := 0; i < numCerts; i++ {
		assert.Equal(t, fragmentInfos[i], fragmentInfo)
		reservationPeriodes[certs[i].BlobHeader.PaymentMetadata.ReservationPeriod] = struct{}{}
	}
	assert.Len(t, reservationPeriodes, numCerts)
	for i := 0; i < numCerts; i++ {
		assert.Contains(t, reservationPeriodes, uint32(i))
	}

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobCertificate"},
		},
	})
}

func TestBlobMetadataStoreUpdateBlobStatus(t *testing.T) {
	ctx := context.Background()
	blobKey, blobHeader := newBlob(t)

	now := time.Now()
	metadata := &v2.BlobMetadata{
		BlobHeader: blobHeader,
		BlobStatus: v2.Queued,
		Expiry:     uint64(now.Add(time.Hour).Unix()),
		NumRetries: 0,
		UpdatedAt:  uint64(now.UnixNano()),
	}
	err := blobMetadataStore.PutBlobMetadata(ctx, metadata)
	assert.NoError(t, err)

	// Update the blob status to invalid status
	err = blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Certified)
	assert.ErrorIs(t, err, blobstore.ErrInvalidStateTransition)

	// Update the blob status to a valid status
	err = blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Encoded)
	assert.NoError(t, err)

	// Update the blob status to same status
	err = blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Encoded)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	fetchedMetadata, err := blobMetadataStore.GetBlobMetadata(ctx, blobKey)
	assert.NoError(t, err)
	assert.Equal(t, fetchedMetadata.BlobStatus, v2.Encoded)
	assert.Greater(t, fetchedMetadata.UpdatedAt, metadata.UpdatedAt)

	// Update the blob status to a valid status
	err = blobMetadataStore.UpdateBlobStatus(ctx, blobKey, v2.Failed)
	assert.NoError(t, err)

	fetchedMetadata, err = blobMetadataStore.GetBlobMetadata(ctx, blobKey)
	assert.NoError(t, err)
	assert.Equal(t, fetchedMetadata.BlobStatus, v2.Failed)

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BlobMetadata"},
		},
	})
}

func TestBlobMetadataStoreDispersals(t *testing.T) {
	ctx := context.Background()
	opID := core.OperatorID{0, 1}
	dispersalRequest := &corev2.DispersalRequest{
		OperatorID:      opID,
		OperatorAddress: gethcommon.HexToAddress("0x1234567"),
		Socket:          "socket",
		DispersedAt:     uint64(time.Now().UnixNano()),

		BatchHeader: corev2.BatchHeader{
			BatchRoot:            [32]byte{1, 2, 3},
			ReferenceBlockNumber: 100,
		},
	}

	err := blobMetadataStore.PutDispersalRequest(ctx, dispersalRequest)
	assert.NoError(t, err)

	bhh, err := dispersalRequest.BatchHeader.Hash()
	assert.NoError(t, err)

	fetchedRequest, err := blobMetadataStore.GetDispersalRequest(ctx, bhh, dispersalRequest.OperatorID)
	assert.NoError(t, err)
	assert.Equal(t, dispersalRequest, fetchedRequest)

	// attempt to put dispersal request with the same key should fail
	err = blobMetadataStore.PutDispersalRequest(ctx, dispersalRequest)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	dispersalResponse := &corev2.DispersalResponse{
		DispersalRequest: dispersalRequest,
		RespondedAt:      uint64(time.Now().UnixNano()),
		Signature:        [32]byte{1, 1, 1},
		Error:            "error",
	}

	err = blobMetadataStore.PutDispersalResponse(ctx, dispersalResponse)
	assert.NoError(t, err)

	fetchedResponse, err := blobMetadataStore.GetDispersalResponse(ctx, bhh, dispersalRequest.OperatorID)
	assert.NoError(t, err)
	assert.Equal(t, dispersalResponse, fetchedResponse)

	// attempt to put dispersal response with the same key should fail
	err = blobMetadataStore.PutDispersalResponse(ctx, dispersalResponse)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "DispersalRequest#" + opID.Hex()},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "DispersalResponse#" + opID.Hex()},
		},
	})
}

func TestBlobMetadataStoreVerificationInfo(t *testing.T) {
	ctx := context.Background()
	blobKey := corev2.BlobKey{1, 1, 1}
	batchHeader := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 100,
	}
	bhh, err := batchHeader.Hash()
	assert.NoError(t, err)
	verificationInfo := &corev2.BlobVerificationInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey,
		BlobIndex:      10,
		InclusionProof: []byte("proof"),
	}

	err = blobMetadataStore.PutBlobVerificationInfo(ctx, verificationInfo)
	assert.NoError(t, err)

	fetchedInfo, err := blobMetadataStore.GetBlobVerificationInfo(ctx, blobKey, bhh)
	assert.NoError(t, err)
	assert.Equal(t, verificationInfo, fetchedInfo)

	// attempt to put verification info with the same key should fail
	err = blobMetadataStore.PutBlobVerificationInfo(ctx, verificationInfo)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	// put multiple verification infos
	blobKey1 := corev2.BlobKey{2, 2, 2}
	verificationInfo1 := &corev2.BlobVerificationInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey1,
		BlobIndex:      12,
		InclusionProof: []byte("proof 1"),
	}
	blobKey2 := corev2.BlobKey{3, 3, 3}
	verificationInfo2 := &corev2.BlobVerificationInfo{
		BatchHeader:    batchHeader,
		BlobKey:        blobKey2,
		BlobIndex:      14,
		InclusionProof: []byte("proof 2"),
	}
	err = blobMetadataStore.PutBlobVerificationInfos(ctx, []*corev2.BlobVerificationInfo{verificationInfo1, verificationInfo2})
	assert.NoError(t, err)

	// test retries
	nonTransientError := errors.New("non transient error")
	mockDynamoClient.On("PutItems", mock.Anything, mock.Anything, mock.Anything).Return(nil, nonTransientError).Once()
	err = mockedBlobMetadataStore.PutBlobVerificationInfos(ctx, []*corev2.BlobVerificationInfo{verificationInfo1, verificationInfo2})
	assert.ErrorIs(t, err, nonTransientError)

	mockDynamoClient.On("PutItems", mock.Anything, mock.Anything, mock.Anything).Return([]dynamodb.Item{
		{
			"PK": &types.AttributeValueMemberS{Value: "BlobKey#" + blobKey1.Hex()},
			"SK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
		},
	}, nil).Run(func(args mock.Arguments) {
		items := args.Get(2).([]dynamodb.Item)
		assert.Len(t, items, 2)
	}).Once()
	mockDynamoClient.On("PutItems", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil).
		Run(func(args mock.Arguments) {
			items := args.Get(2).([]dynamodb.Item)
			assert.Len(t, items, 1)
		}).
		Once()
	err = mockedBlobMetadataStore.PutBlobVerificationInfos(ctx, []*corev2.BlobVerificationInfo{verificationInfo1, verificationInfo2})
	assert.NoError(t, err)
	mockDynamoClient.AssertNumberOfCalls(t, "PutItems", 3)
}

func TestBlobMetadataStoreBatchAttestation(t *testing.T) {
	ctx := context.Background()
	h := &corev2.BatchHeader{
		BatchRoot:            [32]byte{1, 2, 3},
		ReferenceBlockNumber: 100,
	}
	bhh, err := h.Hash()
	assert.NoError(t, err)

	err = blobMetadataStore.PutBatchHeader(ctx, h)
	assert.NoError(t, err)

	fetchedHeader, err := blobMetadataStore.GetBatchHeader(ctx, bhh)
	assert.NoError(t, err)
	assert.Equal(t, h, fetchedHeader)

	// attempt to put batch header with the same key should fail
	err = blobMetadataStore.PutBatchHeader(ctx, h)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	keyPair, err := core.GenRandomBlsKeys()
	assert.NoError(t, err)

	apk := keyPair.GetPubKeyG2()
	attestation := &corev2.Attestation{
		BatchHeader: h,
		AttestedAt:  uint64(time.Now().UnixNano()),
		NonSignerPubKeys: []*core.G1Point{
			core.NewG1Point(big.NewInt(1), big.NewInt(2)),
			core.NewG1Point(big.NewInt(3), big.NewInt(4)),
		},
		APKG2: apk,
		QuorumAPKs: map[uint8]*core.G1Point{
			0: core.NewG1Point(big.NewInt(5), big.NewInt(6)),
			1: core.NewG1Point(big.NewInt(7), big.NewInt(8)),
		},
		Sigma: &core.Signature{
			G1Point: core.NewG1Point(big.NewInt(9), big.NewInt(10)),
		},
		QuorumNumbers: []core.QuorumID{0, 1},
		QuorumResults: map[uint8]uint8{
			0: 100,
			1: 80,
		},
	}

	err = blobMetadataStore.PutAttestation(ctx, attestation)
	assert.NoError(t, err)

	fetchedAttestation, err := blobMetadataStore.GetAttestation(ctx, bhh)
	assert.NoError(t, err)
	assert.Equal(t, attestation, fetchedAttestation)

	// attempt to put attestation with the same key should fail
	err = blobMetadataStore.PutAttestation(ctx, attestation)
	assert.ErrorIs(t, err, common.ErrAlreadyExists)

	// attempt to retrieve batch header and attestation at the same time
	fetchedHeader, fetchedAttestation, err = blobMetadataStore.GetSignedBatch(ctx, bhh)
	assert.NoError(t, err)
	assert.Equal(t, h, fetchedHeader)
	assert.Equal(t, attestation, fetchedAttestation)

	deleteItems(t, []commondynamodb.Key{
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "BatchHeader"},
		},
		{
			"PK": &types.AttributeValueMemberS{Value: "BatchHeader#" + hex.EncodeToString(bhh[:])},
			"SK": &types.AttributeValueMemberS{Value: "Attestation"},
		},
	})
}

func deleteItems(t *testing.T, keys []commondynamodb.Key) {
	failed, err := dynamoClient.DeleteItems(context.Background(), metadataTableName, keys)
	assert.NoError(t, err)
	assert.Len(t, failed, 0)
}

func newBlob(t *testing.T) (corev2.BlobKey, *corev2.BlobHeader) {
	accountBytes := make([]byte, 32)
	_, err := rand.Read(accountBytes)
	require.NoError(t, err)
	accountID := hex.EncodeToString(accountBytes)
	reservationPeriod, err := rand.Int(rand.Reader, big.NewInt(256))
	require.NoError(t, err)
	cumulativePayment, err := rand.Int(rand.Reader, big.NewInt(1024))
	require.NoError(t, err)
	sig := make([]byte, 32)
	_, err = rand.Read(sig)
	require.NoError(t, err)
	bh := &corev2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   []core.QuorumID{0},
		BlobCommitments: mockCommitment,
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         accountID,
			ReservationPeriod: uint32(reservationPeriod.Int64()),
			CumulativePayment: cumulativePayment,
		},
		Signature: sig,
	}
	bk, err := bh.BlobKey()
	require.NoError(t, err)
	return bk, bh
}
