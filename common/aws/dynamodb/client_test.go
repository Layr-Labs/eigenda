package dynamodb_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
	dynamoClient       *commondynamodb.Client
	clientConfig       commonaws.ClientConfig

	localStackPort = "4567"
)

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	teardown()
	os.Exit(code)
}

func setup(m *testing.M) {

	deployLocalstack := !(os.Getenv("DEPLOY_LOCALSTACK") == "false")
	if !deployLocalstack {
		localStackPort = os.Getenv("LOCALSTACK_PORT")
	}

	if deployLocalstack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localStackPort)
		if err != nil {
			teardown()
			panic("failed to start localstack container")
		}
	}

	logger, err := logging.GetLogger(logging.DefaultCLIConfig())
	if err != nil {
		teardown()
		panic("failed to get logger")
	}

	clientConfig = commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}
	dynamoClient, err = commondynamodb.NewClient(clientConfig, logger)
	if err != nil {
		teardown()
		panic("failed to create dynamodb client")
	}
}

func teardown() {
	deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
}

func createTable(t *testing.T, tableName string) {
	ctx := context.Background()
	tableDescription, err := test_utils.CreateTable(ctx, clientConfig, tableName, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("MetadataKey"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("MetadataKey"),
			KeyType:       types.KeyTypeHash,
		}},
		TableName: aws.String(tableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, tableDescription)
}

func TestBasicOperations(t *testing.T) {
	tableName := "Processing"
	createTable(t, tableName)

	ctx := context.Background()
	err := dynamoClient.PutItem(ctx, tableName,
		commondynamodb.Item{
			"MetadataKey": &types.AttributeValueMemberS{Value: "key"},
			"RequestedAt": &types.AttributeValueMemberN{Value: "123"},
			"SecurityParams": &types.AttributeValueMemberL{
				Value: []types.AttributeValue{
					&types.AttributeValueMemberM{
						Value: map[string]types.AttributeValue{
							"QuorumID":           &types.AttributeValueMemberN{Value: "0"},
							"AdversaryThreshold": &types.AttributeValueMemberN{Value: "80"},
						},
					},
					&types.AttributeValueMemberM{
						Value: map[string]types.AttributeValue{
							"QuorumID":           &types.AttributeValueMemberN{Value: "1"},
							"AdversaryThreshold": &types.AttributeValueMemberN{Value: "70"},
						},
					},
				},
			},
			"BlobSize": &types.AttributeValueMemberN{Value: "123"},
			"BlobKey":  &types.AttributeValueMemberS{Value: "blob1"},
			"Status":   &types.AttributeValueMemberS{Value: "Processing"},
		},
	)
	assert.NoError(t, err)

	item, err := dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"MetadataKey": &types.AttributeValueMemberS{Value: "key"},
	})
	assert.NoError(t, err)

	assert.Equal(t, "key", item["MetadataKey"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "123", item["RequestedAt"].(*types.AttributeValueMemberN).Value)
	assert.Equal(t, "Processing", item["Status"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "blob1", item["BlobKey"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "123", item["BlobSize"].(*types.AttributeValueMemberN).Value)
	assert.Equal(t, []types.AttributeValue{
		&types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"QuorumID":           &types.AttributeValueMemberN{Value: "0"},
				"AdversaryThreshold": &types.AttributeValueMemberN{Value: "80"},
			},
		},
		&types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"QuorumID":           &types.AttributeValueMemberN{Value: "1"},
				"AdversaryThreshold": &types.AttributeValueMemberN{Value: "70"},
			},
		},
	}, item["SecurityParams"].(*types.AttributeValueMemberL).Value)

	_, err = dynamoClient.UpdateItem(ctx, tableName, commondynamodb.Key{
		"MetadataKey": &types.AttributeValueMemberS{Value: "key"},
	}, commondynamodb.Item{
		"Status": &types.AttributeValueMemberS{Value: "Confirmed"},
		"BatchHeaderHash": &types.AttributeValueMemberS{
			Value: "0x123",
		},
		"BlobIndex": &types.AttributeValueMemberN{
			Value: "0",
		},
	})
	assert.NoError(t, err)

	item, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"MetadataKey": &types.AttributeValueMemberS{Value: "key"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "key", item["MetadataKey"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "Confirmed", item["Status"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "0x123", item["BatchHeaderHash"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "0", item["BlobIndex"].(*types.AttributeValueMemberN).Value)

	err = dynamoClient.DeleteTable(ctx, tableName)
	assert.NoError(t, err)
}

func TestBatchOperations(t *testing.T) {
	tableName := "Processing"
	createTable(t, tableName)

	ctx := context.Background()
	numItems := 30
	items := make([]commondynamodb.Item, numItems)
	for i := 0; i < numItems; i += 1 {
		items[i] = commondynamodb.Item{
			"MetadataKey": &types.AttributeValueMemberS{Value: fmt.Sprintf("key%d", i)},
			"BlobKey":     &types.AttributeValueMemberS{Value: fmt.Sprintf("blob%d", i)},
		}
	}
	unprocessed, err := dynamoClient.PutItems(ctx, tableName, items)
	assert.NoError(t, err)
	assert.Len(t, unprocessed, 0)

	fetchedItem, err := dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"MetadataKey": &types.AttributeValueMemberS{Value: "key0"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, fetchedItem)
	assert.Equal(t, fetchedItem["BlobKey"].(*types.AttributeValueMemberS).Value, "blob0")

	fetchedItem, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"MetadataKey": &types.AttributeValueMemberS{Value: "key1"},
	})
	assert.NoError(t, err)
	assert.NotNil(t, fetchedItem)
	assert.Equal(t, fetchedItem["BlobKey"].(*types.AttributeValueMemberS).Value, "blob1")

	keys := make([]commondynamodb.Key, numItems)
	for i := 0; i < numItems; i += 1 {
		keys[i] = commondynamodb.Key{
			"MetadataKey": &types.AttributeValueMemberS{Value: fmt.Sprintf("key%d", i)},
		}
	}

	unprocessedKeys, err := dynamoClient.DeleteItems(ctx, tableName, keys)
	assert.NoError(t, err)
	assert.Len(t, unprocessedKeys, 0)

	fetchedItem, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"MetadataKey": &types.AttributeValueMemberS{Value: "key0"},
	})
	assert.NoError(t, err)
	assert.Len(t, fetchedItem, 0)

	fetchedItem, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"MetadataKey": &types.AttributeValueMemberS{Value: "key1"},
	})
	assert.NoError(t, err)
	assert.Len(t, fetchedItem, 0)
}
