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
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigenda/disperser/batcher/batchstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
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
		NumMinibatches:       0,
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
		NumMinibatches:       0,
	}
	_ = minibatchStore.PutBatch(ctx, batch1)
	batch2 := &batcher.BatchRecord{
		ID:                   id2,
		CreatedAt:            ts,
		ReferenceBlockNumber: 1,
		Status:               batcher.BatchStatusFormed,
		NumMinibatches:       0,
	}
	err := minibatchStore.PutBatch(ctx, batch2)
	assert.NoError(t, err)
	batch3 := &batcher.BatchRecord{
		ID:                   id3,
		CreatedAt:            ts,
		ReferenceBlockNumber: 1,
		Status:               batcher.BatchStatusFormed,
		NumMinibatches:       0,
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
		NumMinibatches:       2,
	}
	batch2 := &batcher.BatchRecord{
		ID:                   id2,
		CreatedAt:            ts2,
		ReferenceBlockNumber: 1,
		Status:               batcher.BatchStatusFormed,
		NumMinibatches:       2,
	}
	err := minibatchStore.PutBatch(ctx, batch1)
	assert.NoError(t, err)
	err = minibatchStore.PutBatch(ctx, batch2)
	assert.NoError(t, err)

	batch, err := minibatchStore.GetLatestFormedBatch(ctx)
	assert.NoError(t, err)
	assert.Equal(t, batch2.ID, batch.ID)
}

func TestPutDispersal(t *testing.T) {
	ctx := context.Background()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	ts := time.Now().Truncate(time.Second).UTC()
	opID := core.OperatorID([32]byte{123})
	response := &batcher.MinibatchDispersal{
		BatchID:         id,
		MinibatchIndex:  0,
		OperatorID:      opID,
		OperatorAddress: gcommon.HexToAddress("0x0"),
		Socket:          "socket",
		NumBlobs:        1,
		RequestedAt:     ts,
		DispersalResponse: batcher.DispersalResponse{
			Signatures:  nil,
			RespondedAt: ts,
			Error:       nil,
		},
	}
	err = minibatchStore.PutDispersal(ctx, response)
	assert.NoError(t, err)
	r, err := minibatchStore.GetDispersal(ctx, response.BatchID, response.MinibatchIndex, opID)
	assert.NoError(t, err)
	assert.Equal(t, response, r)
}

func TestDispersalStatus(t *testing.T) {
	ctx := context.Background()
	id, err := uuid.NewV7()
	assert.NoError(t, err)
	ts := time.Now().Truncate(time.Second).UTC()
	opID := core.OperatorID([32]byte{123})

	// no dispersals
	dispersed, err := minibatchStore.BatchDispersed(ctx, id, 0)
	assert.NoError(t, err)
	assert.False(t, dispersed)

	request := &batcher.MinibatchDispersal{
		BatchID:         id,
		MinibatchIndex:  0,
		OperatorID:      opID,
		OperatorAddress: gcommon.HexToAddress("0x0"),
		NumBlobs:        1,
		RequestedAt:     ts,
	}
	err = minibatchStore.PutDispersal(ctx, request)
	assert.NoError(t, err)

	// dispersal request but no response
	dispersed, err = minibatchStore.BatchDispersed(ctx, id, 1)
	assert.NoError(t, err)
	assert.False(t, dispersed)

	response := &batcher.DispersalResponse{
		Signatures:  nil,
		RespondedAt: ts,
		Error:       nil,
	}
	err = minibatchStore.UpdateDispersalResponse(ctx, request, response)
	assert.NoError(t, err)

	// dispersal request and response
	dispersed, err = minibatchStore.BatchDispersed(ctx, id, 1)
	assert.NoError(t, err)
	assert.True(t, dispersed)
	// test with different number of minibatches
	dispersed, err = minibatchStore.BatchDispersed(ctx, id, 2)
	assert.NoError(t, err)
	assert.False(t, dispersed)

}

func TestGetBlobMinibatchMappings(t *testing.T) {
	ctx := context.Background()
	batchID, err := uuid.NewV7()
	assert.NoError(t, err)
	blobKey := disperser.BlobKey{
		BlobHash:     "blob-hash",
		MetadataHash: "metadata-hash",
	}
	var commitX, commitY, lengthX, lengthY fp.Element
	_, err = commitX.SetString("21661178944771197726808973281966770251114553549453983978976194544185382599016")
	assert.NoError(t, err)
	_, err = commitY.SetString("9207254729396071334325696286939045899948985698134704137261649190717970615186")
	assert.NoError(t, err)
	commitment := &encoding.G1Commitment{
		X: commitX,
		Y: commitY,
	}
	_, err = lengthX.SetString("18730744272503541936633286178165146673834730535090946570310418711896464442549")
	assert.NoError(t, err)
	_, err = lengthY.SetString("15356431458378126778840641829778151778222945686256112821552210070627093656047")
	assert.NoError(t, err)
	var lengthXA0, lengthXA1, lengthYA0, lengthYA1 fp.Element
	_, err = lengthXA0.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	assert.NoError(t, err)
	_, err = lengthXA1.SetString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	assert.NoError(t, err)
	_, err = lengthYA0.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	assert.NoError(t, err)
	_, err = lengthYA1.SetString("4082367875863433681332203403145435568316851327593401208105741076214120093531")
	assert.NoError(t, err)

	var lengthProof, lengthCommitment bn254.G2Affine
	lengthProof.X.A0 = lengthXA0
	lengthProof.X.A1 = lengthXA1
	lengthProof.Y.A0 = lengthYA0
	lengthProof.Y.A1 = lengthYA1

	lengthCommitment = lengthProof
	expectedDataLength := 111
	expectedChunkLength := uint(222)
	err = minibatchStore.PutBlobMinibatchMappings(ctx, []*batcher.BlobMinibatchMapping{
		{
			BlobKey:        &blobKey,
			BatchID:        batchID,
			MinibatchIndex: 11,
			BlobIndex:      22,
			BlobHeader: core.BlobHeader{
				BlobCommitments: encoding.BlobCommitments{
					Commitment:       commitment,
					LengthCommitment: (*encoding.G2Commitment)(&lengthCommitment),
					Length:           uint(expectedDataLength),
					LengthProof:      (*encoding.LengthProof)(&lengthProof),
				},
				QuorumInfos: []*core.BlobQuorumInfo{
					{
						ChunkLength: expectedChunkLength,
						SecurityParam: core.SecurityParam{
							QuorumID:              1,
							ConfirmationThreshold: 55,
							AdversaryThreshold:    33,
							QuorumRate:            123,
						},
					},
				},
				AccountID: "account-id",
			},
		},
	})
	assert.NoError(t, err)

	mapping, err := minibatchStore.GetBlobMinibatchMappings(ctx, blobKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(mapping))
	assert.Equal(t, &blobKey, mapping[0].BlobKey)
	assert.Equal(t, batchID, mapping[0].BatchID)
	assert.Equal(t, uint(11), mapping[0].MinibatchIndex)
	assert.Equal(t, uint(22), mapping[0].BlobIndex)
	assert.Equal(t, commitment, mapping[0].BlobCommitments.Commitment)
	lengthCommitmentBytes, err := mapping[0].BlobCommitments.LengthCommitment.Serialize()
	assert.NoError(t, err)
	expectedLengthCommitmentBytes := lengthCommitment.Bytes()
	assert.Equal(t, expectedLengthCommitmentBytes[:], lengthCommitmentBytes[:])
	assert.Equal(t, expectedDataLength, int(mapping[0].BlobCommitments.Length))
	lengthProofBytes, err := mapping[0].BlobCommitments.LengthProof.Serialize()
	assert.NoError(t, err)
	expectedLengthProofBytes := lengthProof.Bytes()
	assert.Equal(t, expectedLengthProofBytes[:], lengthProofBytes[:])
	assert.Len(t, mapping[0].QuorumInfos, 1)
	assert.Equal(t, expectedChunkLength, mapping[0].QuorumInfos[0].ChunkLength)
	assert.Equal(t, core.SecurityParam{
		QuorumID:              1,
		ConfirmationThreshold: 55,
		AdversaryThreshold:    33,
		QuorumRate:            123,
	}, mapping[0].QuorumInfos[0].SecurityParam)
}
