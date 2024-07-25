package batchstore_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/batcher/batchstore"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var (
	logger = logging.NewNoopLogger()

	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource

	deployLocalStack bool
	localStackPort   = "4566"

	dynamoClient   *dynamodb.Client
	minibatchStore *batchstore.MinibatchStore

	UUID               = uuid.New()
	minibatchTableName = fmt.Sprintf("test-MinibatchStore-%v", UUID)
)

func setup(m *testing.M) {
	deployLocalStack = !(os.Getenv("DEPLOY_LOCALSTACK") == "false")
	fmt.Printf("deployLocalStack: %v\n", deployLocalStack)
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

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}

	_, err := test_utils.CreateTable(context.Background(), cfg, minibatchTableName, batchstore.GenerateTableSchema(minibatchTableName, 10, 10))
	if err != nil {
		teardown()
		panic("failed to create dynamodb table: " + err.Error())
	}

	dynamoClient, err = dynamodb.NewClient(cfg, logger)
	if err != nil {
		teardown()
		panic("failed to create dynamodb client: " + err.Error())
	}

	minibatchStore = batchstore.NewMinibatchStore(dynamoClient, logger, minibatchTableName, time.Hour)
}

func teardown() {
	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func TestPutBatch(t *testing.T) {
	ctx := context.Background()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	ts := time.Now().Truncate(time.Second).UTC()
	batch := &batcher.BatchRecord{
		ID:                   id,
		CreatedAt:            ts,
		ReferenceBlockNumber: 1,
		Status:               batcher.BatchStatusPending,
		HeaderHash:           [32]byte{1},
		AggregatePubKey:      nil,
		AggregateSignature:   nil,
	}
	err = minibatchStore.PutBatch(ctx, batch)
	assert.NoError(t, err)
	err = minibatchStore.PutBatch(ctx, batch)
	assert.Error(t, err)
	b, err := minibatchStore.GetBatch(ctx, batch.ID)
	assert.NoError(t, err)
	assert.Equal(t, batch, b)
	err = minibatchStore.UpdateBatchStatus(ctx, batch.ID, batcher.BatchStatusFormed)
	assert.NoError(t, err)
	u, err := minibatchStore.GetBatch(ctx, batch.ID)
	assert.NoError(t, err)
	assert.Equal(t, batcher.BatchStatusFormed, u.Status)
	err = minibatchStore.UpdateBatchStatus(ctx, batch.ID, batcher.BatchStatusAttested)
	assert.NoError(t, err)
	u, err = minibatchStore.GetBatch(ctx, batch.ID)
	assert.NoError(t, err)
	assert.Equal(t, batcher.BatchStatusAttested, u.Status)
	err = minibatchStore.UpdateBatchStatus(ctx, batch.ID, batcher.BatchStatusFailed)
	assert.NoError(t, err)
	u, err = minibatchStore.GetBatch(ctx, batch.ID)
	assert.NoError(t, err)
	assert.Equal(t, batcher.BatchStatusFailed, u.Status)
	err = minibatchStore.UpdateBatchStatus(ctx, batch.ID, 4)
	assert.Error(t, err)
}

func TestGetBatchesByStatus(t *testing.T) {
	ctx := context.Background()
	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()
	id3, _ := uuid.NewV7()
	ts := time.Now().Truncate(time.Second).UTC()
	batch1 := &batcher.BatchRecord{
		ID:                   id1,
		CreatedAt:            ts,
		ReferenceBlockNumber: 1,
		Status:               batcher.BatchStatusFormed,
		HeaderHash:           [32]byte{1},
		AggregatePubKey:      nil,
		AggregateSignature:   nil,
	}
	_ = minibatchStore.PutBatch(ctx, batch1)
	batch2 := &batcher.BatchRecord{
		ID:                   id2,
		CreatedAt:            ts,
		ReferenceBlockNumber: 1,
		Status:               batcher.BatchStatusFormed,
		HeaderHash:           [32]byte{1},
		AggregatePubKey:      nil,
		AggregateSignature:   nil,
	}
	err := minibatchStore.PutBatch(ctx, batch2)
	assert.NoError(t, err)
	batch3 := &batcher.BatchRecord{
		ID:                   id3,
		CreatedAt:            ts,
		ReferenceBlockNumber: 1,
		Status:               batcher.BatchStatusFormed,
		HeaderHash:           [32]byte{1},
		AggregatePubKey:      nil,
		AggregateSignature:   nil,
	}
	err = minibatchStore.PutBatch(ctx, batch3)
	assert.NoError(t, err)

	attested, err := minibatchStore.GetBatchesByStatus(ctx, batcher.BatchStatusAttested)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(attested))

	formed, err := minibatchStore.GetBatchesByStatus(ctx, batcher.BatchStatusFormed)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(formed))

	err = minibatchStore.UpdateBatchStatus(ctx, id1, batcher.BatchStatusAttested)
	assert.NoError(t, err)

	formed, err = minibatchStore.GetBatchesByStatus(ctx, batcher.BatchStatusFormed)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(formed))
}

func TestPutMinibatch(t *testing.T) {
	ctx := context.Background()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	minibatch := &batcher.MinibatchRecord{
		BatchID:              id,
		MinibatchIndex:       12,
		BlobHeaderHashes:     [][32]byte{{1}},
		BatchSize:            1,
		ReferenceBlockNumber: 1,
	}
	err = minibatchStore.PutMinibatch(ctx, minibatch)
	assert.NoError(t, err)
	m, err := minibatchStore.GetMinibatch(ctx, minibatch.BatchID, minibatch.MinibatchIndex)
	assert.NoError(t, err)
	assert.Equal(t, minibatch, m)
}

func TestGetLatestFormedBatch(t *testing.T) {
	ctx := context.Background()
	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()
	ts := time.Now().Truncate(time.Second).UTC()
	ts2 := ts.Add(10 * time.Second)
	batch1 := &batcher.BatchRecord{
		ID:                   id1,
		CreatedAt:            ts,
		ReferenceBlockNumber: 1,
		Status:               batcher.BatchStatusFormed,
		HeaderHash:           [32]byte{1},
		AggregatePubKey:      nil,
		AggregateSignature:   nil,
	}
	minibatch1 := &batcher.MinibatchRecord{
		BatchID:              id1,
		MinibatchIndex:       1,
		BlobHeaderHashes:     [][32]byte{{1}},
		BatchSize:            1,
		ReferenceBlockNumber: 1,
	}
	batch2 := &batcher.BatchRecord{
		ID:                   id2,
		CreatedAt:            ts2,
		ReferenceBlockNumber: 1,
		Status:               batcher.BatchStatusFormed,
		HeaderHash:           [32]byte{1},
		AggregatePubKey:      nil,
		AggregateSignature:   nil,
	}
	minibatch2 := &batcher.MinibatchRecord{
		BatchID:              id2,
		MinibatchIndex:       1,
		BlobHeaderHashes:     [][32]byte{{1}},
		BatchSize:            1,
		ReferenceBlockNumber: 1,
	}
	minibatch3 := &batcher.MinibatchRecord{
		BatchID:              id2,
		MinibatchIndex:       2,
		BlobHeaderHashes:     [][32]byte{{1}},
		BatchSize:            1,
		ReferenceBlockNumber: 1,
	}
	err := minibatchStore.PutBatch(ctx, batch1)
	assert.NoError(t, err)
	err = minibatchStore.PutMinibatch(ctx, minibatch1)
	assert.NoError(t, err)
	err = minibatchStore.PutBatch(ctx, batch2)
	assert.NoError(t, err)
	err = minibatchStore.PutMinibatch(ctx, minibatch2)
	assert.NoError(t, err)
	err = minibatchStore.PutMinibatch(ctx, minibatch3)
	assert.NoError(t, err)

	batch, minibatches, err := minibatchStore.GetLatestFormedBatch(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(minibatches))
	assert.Equal(t, batch.ID, batch2.ID)
}
func TestPutDispersalRequest(t *testing.T) {
	ctx := context.Background()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	ts := time.Now().Truncate(time.Second).UTC()
	opID := core.OperatorID([32]byte{123})
	request := &batcher.DispersalRequest{
		BatchID:         id,
		MinibatchIndex:  0,
		OperatorID:      opID,
		OperatorAddress: gcommon.HexToAddress("0x0"),
		NumBlobs:        1,
		RequestedAt:     ts,
		BlobHash:        "blobHash",
		MetadataHash:    "metadataHash",
	}
	err = minibatchStore.PutDispersalRequest(ctx, request)
	assert.NoError(t, err)
	r, err := minibatchStore.GetDispersalRequest(ctx, request.BatchID, request.MinibatchIndex, opID)
	assert.NoError(t, err)
	assert.Equal(t, request, r)
}

func TestPutDispersalResponse(t *testing.T) {
	ctx := context.Background()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	ts := time.Now().Truncate(time.Second).UTC()
	opID := core.OperatorID([32]byte{123})
	blobHash := "blobHash"
	metadataHash := "metadataHash"
	response := &batcher.DispersalResponse{
		DispersalRequest: batcher.DispersalRequest{
			BatchID:         id,
			MinibatchIndex:  0,
			OperatorID:      opID,
			OperatorAddress: gcommon.HexToAddress("0x0"),
			NumBlobs:        1,
			RequestedAt:     ts,
			BlobHash:        blobHash,
			MetadataHash:    metadataHash,
		},
		Signatures:  nil,
		RespondedAt: ts,
		Error:       nil,
	}
	err = minibatchStore.PutDispersalResponse(ctx, response)
	assert.NoError(t, err)
	r, err := minibatchStore.GetDispersalResponse(ctx, response.BatchID, response.MinibatchIndex, opID)
	assert.NoError(t, err)
	assert.Equal(t, response, r)
}

func TestDispersalStatus(t *testing.T) {
	ctx := context.Background()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	ts := time.Now().Truncate(time.Second).UTC()
	opID := core.OperatorID([32]byte{123})
	blobHash := "blobHash"
	metadataHash := "metadataHash"

	// no dispersals
	dispersed, err := minibatchStore.BatchDispersed(ctx, id)
	assert.NoError(t, err)
	assert.False(t, dispersed)

	request := &batcher.DispersalRequest{
		BatchID:         id,
		MinibatchIndex:  0,
		OperatorID:      opID,
		OperatorAddress: gcommon.HexToAddress("0x0"),
		NumBlobs:        1,
		RequestedAt:     ts,
		BlobHash:        blobHash,
		MetadataHash:    metadataHash,
	}
	err = minibatchStore.PutDispersalRequest(ctx, request)
	assert.NoError(t, err)

	// dispersal request but no response
	dispersed, err = minibatchStore.BatchDispersed(ctx, id)
	assert.NoError(t, err)
	assert.False(t, dispersed)

	response := &batcher.DispersalResponse{
		DispersalRequest: batcher.DispersalRequest{
			BatchID:         id,
			MinibatchIndex:  0,
			OperatorID:      opID,
			OperatorAddress: gcommon.HexToAddress("0x0"),
			NumBlobs:        1,
			RequestedAt:     ts,
			BlobHash:        blobHash,
			MetadataHash:    metadataHash,
		},
		Signatures:  nil,
		RespondedAt: ts,
		Error:       nil,
	}
	err = minibatchStore.PutDispersalResponse(ctx, response)
	assert.NoError(t, err)

	// dispersal request and response
	dispersed, err = minibatchStore.BatchDispersed(ctx, id)
	assert.NoError(t, err)
	assert.True(t, dispersed)
}
