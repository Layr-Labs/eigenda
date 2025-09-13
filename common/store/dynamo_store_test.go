package store_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/testbed"
	"github.com/stretchr/testify/require"
)

var (
	logger = test.GetLogger()

	localStackContainer *testbed.LocalStackContainer

	deployLocalStack bool
	localStackPort   = "4572"

	dynamoClient     dynamodb.Client
	dynamoParamStore common.KVStore[common.RateBucketParams]
	bucketTableName  = "BucketStore"
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(_ *testing.M) {
	deployLocalStack = (os.Getenv("DEPLOY_LOCALSTACK") != "false")
	if !deployLocalStack {
		localStackPort = os.Getenv("LOCALSTACK_PORT")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if deployLocalStack {
		// Start LocalStack container
		var err error
		localStackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       localStackPort,
			Services:       []string{"dynamodb"},
			Logger:         logger,
		})
		if err != nil {
			teardown()
			logger.Fatal("Failed to start localstack container:", err)
		}

		// Extract port from the endpoint for compatibility with existing code
		// The endpoint is in format "http://host:port", we need just the port
		endpoint := localStackContainer.Endpoint()
		if idx := strings.LastIndex(endpoint, ":"); idx != -1 {
			localStackPort = endpoint[idx+1:]
		}
	}

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}

	_, err := test_utils.CreateTable(ctx, cfg, bucketTableName, store.GenerateTableSchema(10, 10, bucketTableName))
	if err != nil {
		teardown()
		logger.Fatal("Failed to create dynamodb table:", err)
	}

	dynamoClient, err = dynamodb.NewClient(cfg, logger)
	if err != nil {
		teardown()
		logger.Fatal("Failed to create dynamodb client:", err)
	}

	dynamoParamStore = store.NewDynamoParamStore[common.RateBucketParams](dynamoClient, bucketTableName)
}

func teardown() {
	if deployLocalStack && localStackContainer != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := localStackContainer.Terminate(ctx); err != nil {
			logger.Error("failed to terminate LocalStack container", "error", err)
		}
	}
}

func TestDynamoBucketStore(t *testing.T) {
	ctx := t.Context()

	p := &common.RateBucketParams{
		BucketLevels:    []time.Duration{time.Second, time.Minute},
		LastRequestTime: time.Now().UTC(),
	}

	t.Run("get_nonexistent_item", func(t *testing.T) {
		p2, err := dynamoParamStore.GetItem(ctx, "testRetriever")
		require.Error(t, err, "should error when item doesn't exist")
		require.Nil(t, p2, "should return nil when item doesn't exist")
	})

	t.Run("update_and_get_item", func(t *testing.T) {
		err := dynamoParamStore.UpdateItem(ctx, "testRetriever", p)
		require.NoError(t, err, "failed to update item in store")

		p2, err := dynamoParamStore.GetItem(ctx, "testRetriever")
		require.NoError(t, err, "failed to get item from store")
		require.Equal(t, p, p2, "retrieved item should match stored item")
	})
}
