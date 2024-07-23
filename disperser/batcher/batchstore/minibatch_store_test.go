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

	fmt.Printf("m: %v\n", m)
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
		HeaderHash:           [32]byte{1},
		AggregatePubKey:      nil,
		AggregateSignature:   nil,
	}
	err = minibatchStore.PutBatch(ctx, batch)
	assert.NoError(t, err)
	b, err := minibatchStore.GetBatch(ctx, batch.ID)
	assert.NoError(t, err)
	assert.Equal(t, batch, b)
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

func TestPutDispersalRequest(t *testing.T) {
	ctx := context.Background()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	ts := time.Now().Truncate(time.Second).UTC()
	request := &batcher.DispersalRequest{
		BatchID:         id,
		MinibatchIndex:  0,
		OperatorID:      core.OperatorID([32]byte{123}),
		OperatorAddress: gcommon.HexToAddress("0x0"),
		NumBlobs:        1,
		RequestedAt:     ts,
	}
	err = minibatchStore.PutDispersalRequest(ctx, request)
	assert.NoError(t, err)
	r, err := minibatchStore.GetDispersalRequest(ctx, request.BatchID, request.MinibatchIndex)
	assert.NoError(t, err)
	assert.Equal(t, request, r)
}

func TestPutDispersalResponse(t *testing.T) {
	ctx := context.Background()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	ts := time.Now().Truncate(time.Second).UTC()
	response := &batcher.DispersalResponse{
		DispersalRequest: batcher.DispersalRequest{
			BatchID:         id,
			MinibatchIndex:  0,
			OperatorID:      core.OperatorID([32]byte{1}),
			OperatorAddress: gcommon.HexToAddress("0x0"),
			NumBlobs:        1,
			RequestedAt:     ts,
		},
		Signatures:  nil,
		RespondedAt: ts,
		Error:       nil,
	}
	err = minibatchStore.PutDispersalResponse(ctx, response)
	assert.NoError(t, err)
	r, err := minibatchStore.GetDispersalResponse(ctx, response.BatchID, response.MinibatchIndex)
	assert.NoError(t, err)
	assert.Equal(t, response, r)
}
