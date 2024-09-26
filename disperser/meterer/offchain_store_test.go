package meterer_test

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"testing"
// 	"time"

// 	"github.com/Layr-Labs/eigenda/common"
// 	commonaws "github.com/Layr-Labs/eigenda/common/aws"
// 	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
// 	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
// 	"github.com/Layr-Labs/eigenda/inabox/deploy"
// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
// 	"github.com/ory/dockertest/v3"
// 	"github.com/stretchr/testify/assert"
// )

// var (
// 	dockertestPool     *dockertest.Pool
// 	dockertestResource *dockertest.Resource
// 	dynamoClient       *commondynamodb.Client
// 	clientConfig       commonaws.ClientConfig

// 	deployLocalStack bool
// 	localStackPort   = "4567"
// )

// func TestMain(m *testing.M) {
// 	setup(m)
// 	code := m.Run()
// 	teardown()
// 	os.Exit(code)
// }

// func setup(m *testing.M) {

// 	deployLocalStack = !(os.Getenv("DEPLOY_LOCALSTACK") == "false")
// 	if !deployLocalStack {
// 		localStackPort = os.Getenv("LOCALSTACK_PORT")
// 	}

// 	if deployLocalStack {
// 		var err error
// 		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localStackPort)
// 		if err != nil {
// 			teardown()
// 			panic("failed to start localstack container")
// 		}
// 	}

// 	loggerConfig := common.DefaultLoggerConfig()
// 	logger, err := common.NewLogger(loggerConfig)
// 	if err != nil {
// 		teardown()
// 		panic("failed to create logger")
// 	}

// 	clientConfig = commonaws.ClientConfig{
// 		Region:          "us-east-1",
// 		AccessKey:       "localstack",
// 		SecretAccessKey: "localstack",
// 		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
// 	}

// 	dynamoClient, err = commondynamodb.NewClient(clientConfig, logger)
// 	if err != nil {
// 		teardown()
// 		panic("failed to create dynamodb client")
// 	}
// }

// func teardown() {
// 	if deployLocalStack {
// 		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
// 	}
// }

// // func CreateReservationTable(t *testing.T, tableName string) {
// // 	ctx := context.Background()
// // 	tableDescription, err := test_utils.CreateTable(ctx, clientConfig, tableName, &dynamodb.CreateTableInput{
// // 		AttributeDefinitions: []types.AttributeDefinition{
// // 			{
// // 				AttributeName: aws.String("AccountID"),
// // 				AttributeType: types.ScalarAttributeTypeS,
// // 			},
// // 			{
// // 				AttributeName: aws.String("BinIndex"),
// // 				AttributeType: types.ScalarAttributeTypeN,
// // 			},
// // 		},
// // 		KeySchema: []types.KeySchemaElement{
// // 			{
// // 				AttributeName: aws.String("AccountID"),
// // 				KeyType:       types.KeyTypeHash,
// // 			},
// // 			{
// // 				AttributeName: aws.String("BinIndex"),
// // 				KeyType:       types.KeyTypeRange,
// // 			},
// // 		},
// // 		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
// // 			{
// // 				IndexName: aws.String("AccountIDIndex"),
// // 				KeySchema: []types.KeySchemaElement{
// // 					{
// // 						AttributeName: aws.String("AccountID"),
// // 						KeyType:       types.KeyTypeHash,
// // 					},
// // 				},
// // 				Projection: &types.Projection{
// // 					ProjectionType: types.ProjectionTypeAll, // ProjectionTypeAll means all attributes are projected into the index
// // 				},
// // 				ProvisionedThroughput: &types.ProvisionedThroughput{
// // 					ReadCapacityUnits:  aws.Int64(10),
// // 					WriteCapacityUnits: aws.Int64(10),
// // 				},
// // 			},
// // 		},
// // 		TableName: aws.String(tableName),
// // 		ProvisionedThroughput: &types.ProvisionedThroughput{
// // 			ReadCapacityUnits:  aws.Int64(10),
// // 			WriteCapacityUnits: aws.Int64(10),
// // 		},
// // 	})
// // 	assert.NoError(t, err)
// // 	assert.NotNil(t, tableDescription)
// // }

// // func CreateGlobalReservationTable(t *testing.T, tableName string) {
// // 	ctx := context.Background()
// // 	tableDescription, err := test_utils.CreateTable(ctx, clientConfig, tableName, &dynamodb.CreateTableInput{
// // 		AttributeDefinitions: []types.AttributeDefinition{
// // 			{
// // 				AttributeName: aws.String("BinIndex"),
// // 				AttributeType: types.ScalarAttributeTypeN,
// // 			},
// // 		},
// // 		KeySchema: []types.KeySchemaElement{
// // 			{
// // 				AttributeName: aws.String("BinIndex"),
// // 				KeyType:       types.KeyTypeHash,
// // 			},
// // 		},
// // 		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
// // 			{
// // 				IndexName: aws.String("BinIndexIndex"),
// // 				KeySchema: []types.KeySchemaElement{
// // 					{
// // 						AttributeName: aws.String("BinIndex"),
// // 						KeyType:       types.KeyTypeHash,
// // 					},
// // 				},
// // 				Projection: &types.Projection{
// // 					ProjectionType: types.ProjectionTypeAll, // ProjectionTypeAll means all attributes are projected into the index
// // 				},
// // 				ProvisionedThroughput: &types.ProvisionedThroughput{
// // 					ReadCapacityUnits:  aws.Int64(10),
// // 					WriteCapacityUnits: aws.Int64(10),
// // 				},
// // 			},
// // 		},
// // 		TableName: aws.String(tableName),
// // 		ProvisionedThroughput: &types.ProvisionedThroughput{
// // 			ReadCapacityUnits:  aws.Int64(10),
// // 			WriteCapacityUnits: aws.Int64(10),
// // 		},
// // 	})
// // 	assert.NoError(t, err)
// // 	assert.NotNil(t, tableDescription)
// // }

// // func CreateOnDemandTable(t *testing.T, tableName string) {
// // 	ctx := context.Background()
// // 	tableDescription, err := test_utils.CreateTable(ctx, clientConfig, tableName, &dynamodb.CreateTableInput{
// // 		AttributeDefinitions: []types.AttributeDefinition{
// // 			{
// // 				AttributeName: aws.String("AccountID"),
// // 				AttributeType: types.ScalarAttributeTypeS,
// // 			},
// // 			{
// // 				AttributeName: aws.String("CumulativePayments"),
// // 				AttributeType: types.ScalarAttributeTypeS,
// // 			},
// // 		},
// // 		KeySchema: []types.KeySchemaElement{
// // 			{
// // 				AttributeName: aws.String("AccountID"),
// // 				KeyType:       types.KeyTypeHash,
// // 			},
// // 			{
// // 				AttributeName: aws.String("CumulativePayments"),
// // 				KeyType:       types.KeyTypeRange,
// // 			},
// // 		},
// // 		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
// // 			{
// // 				IndexName: aws.String("AccountIDIndex"),
// // 				KeySchema: []types.KeySchemaElement{
// // 					{
// // 						AttributeName: aws.String("AccountID"),
// // 						KeyType:       types.KeyTypeHash,
// // 					},
// // 				},
// // 				Projection: &types.Projection{
// // 					ProjectionType: types.ProjectionTypeAll, // ProjectionTypeAll means all attributes are projected into the index
// // 				},
// // 				ProvisionedThroughput: &types.ProvisionedThroughput{
// // 					ReadCapacityUnits:  aws.Int64(10),
// // 					WriteCapacityUnits: aws.Int64(10),
// // 				},
// // 			},
// // 		},
// // 		TableName: aws.String(tableName),
// // 		ProvisionedThroughput: &types.ProvisionedThroughput{
// // 			ReadCapacityUnits:  aws.Int64(10),
// // 			WriteCapacityUnits: aws.Int64(10),
// // 		},
// // 	})
// // 	assert.NoError(t, err)
// // 	assert.NotNil(t, tableDescription)
// // }

// func TestReservationBins(t *testing.T) {
// 	tableName := "reservations"
// 	CreateReservationTable(t, tableName)
// 	indexName := "AccountIDIndex"

// 	ctx := context.Background()
// 	err := dynamoClient.PutItem(ctx, tableName,
// 		commondynamodb.Item{
// 			"AccountID": &types.AttributeValueMemberS{Value: "account1"},
// 			"BinIndex":  &types.AttributeValueMemberN{Value: "1"},
// 			"BinUsage":  &types.AttributeValueMemberN{Value: "1000"},
// 			"UpdatedAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
// 		},
// 	)
// 	assert.NoError(t, err)

// 	item, err := dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
// 		"AccountID": &types.AttributeValueMemberS{Value: "account1"},
// 		"BinIndex":  &types.AttributeValueMemberN{Value: "1"},
// 	})
// 	assert.NoError(t, err)

// 	assert.Equal(t, "account1", item["AccountID"].(*types.AttributeValueMemberS).Value)
// 	assert.Equal(t, "1", item["BinIndex"].(*types.AttributeValueMemberN).Value)
// 	assert.Equal(t, "1000", item["BinUsage"].(*types.AttributeValueMemberN).Value)

// 	items, err := dynamoClient.QueryIndex(ctx, tableName, indexName, "AccountID = :account", commondynamodb.ExpresseionValues{
// 		":account": &types.AttributeValueMemberS{Value: "account1"},
// 	})
// 	assert.NoError(t, err)
// 	assert.Len(t, items, 1)

// 	item, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
// 		"AccountID": &types.AttributeValueMemberS{Value: "account2"},
// 	})
// 	assert.Error(t, err)

// 	_, err = dynamoClient.UpdateItem(ctx, tableName, commondynamodb.Key{
// 		"AccountID": &types.AttributeValueMemberS{Value: "account1"},
// 		"BinIndex":  &types.AttributeValueMemberN{Value: "1"},
// 	}, commondynamodb.Item{
// 		"BinUsage": &types.AttributeValueMemberN{Value: "2000"},
// 	})
// 	assert.NoError(t, err)
// 	err = dynamoClient.PutItem(ctx, tableName,
// 		commondynamodb.Item{
// 			"AccountID": &types.AttributeValueMemberS{Value: "account2"},
// 			"BinIndex":  &types.AttributeValueMemberN{Value: "1"},
// 			"BinUsage":  &types.AttributeValueMemberN{Value: "3000"},
// 			"UpdatedAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
// 		},
// 	)
// 	assert.NoError(t, err)

// 	item, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
// 		"AccountID": &types.AttributeValueMemberS{Value: "account1"},
// 		"BinIndex":  &types.AttributeValueMemberN{Value: "1"},
// 	})
// 	assert.NoError(t, err)
// 	assert.Equal(t, "2000", item["BinUsage"].(*types.AttributeValueMemberN).Value)

// 	items, err = dynamoClient.Query(ctx, tableName, "AccountID = :account", commondynamodb.ExpresseionValues{
// 		":account": &types.AttributeValueMemberS{Value: "account1"},
// 	})
// 	assert.NoError(t, err)
// 	assert.Len(t, items, 1)
// 	assert.Equal(t, "2000", items[0]["BinUsage"].(*types.AttributeValueMemberN).Value)

// 	item, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
// 		"AccountID": &types.AttributeValueMemberS{Value: "account2"},
// 		"BinIndex":  &types.AttributeValueMemberN{Value: "1"},
// 	})
// 	assert.NoError(t, err)
// 	assert.Equal(t, "3000", item["BinUsage"].(*types.AttributeValueMemberN).Value)

// 	err = dynamoClient.DeleteTable(ctx, tableName)
// 	assert.NoError(t, err)
// }

// func TestGlobalBins(t *testing.T) {
// 	tableName := "global"
// 	CreateGlobalReservationTable(t, tableName)
// 	indexName := "BinIndexIndex"
// 	// expression := "BinUsage + :inc"

// 	ctx := context.Background()
// 	numItems := 30
// 	items := make([]commondynamodb.Item, numItems)
// 	for i := 0; i < numItems; i += 1 {
// 		items[i] = commondynamodb.Item{
// 			"BinIndex":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", i)},
// 			"BinUsage":  &types.AttributeValueMemberN{Value: "1000"},
// 			"UpdatedAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
// 		}
// 	}
// 	unprocessed, err := dynamoClient.PutItems(ctx, tableName, items)
// 	assert.NoError(t, err)
// 	assert.Len(t, unprocessed, 0)

// 	queryResult, err := dynamoClient.QueryIndex(ctx, tableName, indexName, "BinIndex = :index", commondynamodb.ExpresseionValues{
// 		":index": &types.AttributeValueMemberN{
// 			Value: "1",
// 		}})
// 	assert.NoError(t, err)
// 	assert.Len(t, queryResult, 1)
// 	assert.Equal(t, "1", queryResult[0]["BinIndex"].(*types.AttributeValueMemberN).Value)
// 	assert.Equal(t, "1000", queryResult[0]["BinUsage"].(*types.AttributeValueMemberN).Value)

// 	queryResult, err = dynamoClient.QueryIndex(ctx, tableName, indexName, "BinIndex = :index", commondynamodb.ExpresseionValues{
// 		":index": &types.AttributeValueMemberN{
// 			Value: "1",
// 		}})
// 	assert.NoError(t, err)
// 	assert.Len(t, queryResult, 1)
// 	assert.Equal(t, "1", queryResult[0]["BinIndex"].(*types.AttributeValueMemberN).Value)
// 	assert.Equal(t, "1000", queryResult[0]["BinUsage"].(*types.AttributeValueMemberN).Value)

// 	queryResult, err = dynamoClient.QueryIndex(ctx, tableName, indexName, "BinIndex = :index", commondynamodb.ExpresseionValues{
// 		":index": &types.AttributeValueMemberN{
// 			Value: "32",
// 		}})
// 	assert.NoError(t, err)
// 	assert.Len(t, queryResult, 0)

// 	_, err = dynamoClient.UpdateItem(ctx, tableName, commondynamodb.Key{
// 		"BinIndex": &types.AttributeValueMemberN{Value: "1"},
// 		// "BinUsage": &types.AttributeValueMemberN{Value: "1000"},
// 	}, commondynamodb.Item{
// 		"BinUsage": &types.AttributeValueMemberN{Value: "2000"},
// 	})
// 	assert.NoError(t, err)
// 	// assert.Equal(t, "1", updatedItem["BinIndex"].(*types.AttributeValueMemberN).Value)
// 	// assert.Equal(t, "2000", updatedItem["BinUsage"].(*types.AttributeValueMemberN).Value)

// 	err = dynamoClient.PutItem(ctx, tableName,
// 		commondynamodb.Item{
// 			"BinIndex":  &types.AttributeValueMemberN{Value: "2"},
// 			"BinUsage":  &types.AttributeValueMemberN{Value: "3000"},
// 			"UpdatedAt": &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
// 		},
// 	)
// 	assert.NoError(t, err)

// 	queryResult, err = dynamoClient.QueryIndex(ctx, tableName, indexName, "BinIndex = :index", commondynamodb.ExpresseionValues{
// 		":index": &types.AttributeValueMemberN{
// 			Value: "1",
// 		}})
// 	assert.NoError(t, err)
// 	assert.Len(t, queryResult, 1)
// 	assert.Equal(t, "1", queryResult[0]["BinIndex"].(*types.AttributeValueMemberN).Value)
// 	assert.Equal(t, "2000", queryResult[0]["BinUsage"].(*types.AttributeValueMemberN).Value)

// 	queryResult, err = dynamoClient.QueryIndex(ctx, tableName, indexName, "BinIndex = :index", commondynamodb.ExpresseionValues{
// 		":index": &types.AttributeValueMemberN{
// 			Value: "2",
// 		}})
// 	assert.NoError(t, err)
// 	assert.Len(t, queryResult, 1)
// 	assert.Equal(t, "2", queryResult[0]["BinIndex"].(*types.AttributeValueMemberN).Value)
// 	assert.Equal(t, "3000", queryResult[0]["BinUsage"].(*types.AttributeValueMemberN).Value)

// 	// items, err = dynamoClient.QueryIndex(ctx, tableName, indexName, "BinIndex = :index", commondynamodb.ExpresseionValues{
// 	// 	":index": &types.AttributeValueMemberN{
// 	// 		Value: "1",
// 	// 	},
// 	// 	""
// 	// })
// 	// assert.NoError(t, err)
// 	// assert.Equal(t, "2000", items[0]["BinUsage"].(*types.AttributeValueMemberN).Value)
// 	// assert.Equal(t, "3000", items[1]["BinUsage"].(*types.AttributeValueMemberN).Value)
// }

// func TestOnDemandUsage(t *testing.T) {
// 	tableName := "ondemand"
// 	CreateOnDemandTable(t, tableName)
// 	indexName := "AccountIDIndex"

// 	ctx := context.Background()

// 	err := dynamoClient.PutItem(ctx, tableName,
// 		commondynamodb.Item{
// 			"AccountID":          &types.AttributeValueMemberS{Value: "account1"},
// 			"CumulativePayments": &types.AttributeValueMemberS{Value: "1"},
// 			"BlobSize":           &types.AttributeValueMemberN{Value: "1000"},
// 			// "UpdatedAt":          &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
// 		},
// 	)
// 	assert.NoError(t, err)

// 	numItems := 30
// 	repetitions := 5
// 	items := make([]commondynamodb.Item, numItems)
// 	for i := 0; i < numItems; i += 1 {
// 		items[i] = commondynamodb.Item{
// 			"AccountID":          &types.AttributeValueMemberS{Value: fmt.Sprintf("account%d", i%repetitions)},
// 			"CumulativePayments": &types.AttributeValueMemberS{Value: fmt.Sprintf("%d", i)},
// 			"BlobSize":           &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", i*1000)},
// 		}
// 	}
// 	unprocessed, err := dynamoClient.PutItems(ctx, tableName, items)
// 	assert.NoError(t, err)
// 	assert.Len(t, unprocessed, 0)

// 	// get item
// 	item, err := dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
// 		"AccountID":          &types.AttributeValueMemberS{Value: "account1"},
// 		"CumulativePayments": &types.AttributeValueMemberS{Value: "1"},
// 	})
// 	assert.NoError(t, err)
// 	assert.Equal(t, "1", item["CumulativePayments"].(*types.AttributeValueMemberS).Value)
// 	assert.Equal(t, "1000", item["BlobSize"].(*types.AttributeValueMemberN).Value)

// 	queryResult, err := dynamoClient.QueryIndex(ctx, tableName, indexName, "AccountID = :account", commondynamodb.ExpresseionValues{
// 		":account": &types.AttributeValueMemberS{
// 			Value: "account1",
// 		}})
// 	assert.NoError(t, err)
// 	assert.Len(t, queryResult, numItems/repetitions)
// 	for _, item := range queryResult {
// 		cumulativePayments, _ := strconv.Atoi(item["CumulativePayments"].(*types.AttributeValueMemberS).Value)
// 		assert.Equal(t, fmt.Sprintf("%d", cumulativePayments*1000), item["BlobSize"].(*types.AttributeValueMemberN).Value)
// 	}
// 	queryResult, err = dynamoClient.QueryIndex(ctx, tableName, indexName, "AccountID = :account_id", commondynamodb.ExpresseionValues{
// 		":account_id": &types.AttributeValueMemberS{
// 			Value: fmt.Sprintf("account%d", numItems/repetitions+1),
// 		}})
// 	assert.NoError(t, err)
// 	assert.Len(t, queryResult, 0)

// 	updatedItem, err := dynamoClient.UpdateItem(ctx, tableName, commondynamodb.Key{
// 		"AccountID":          &types.AttributeValueMemberS{Value: "account1"},
// 		"CumulativePayments": &types.AttributeValueMemberS{Value: "1"},
// 		// "BinUsage": &types.AttributeValueMemberN{Value: "1000"},
// 	}, commondynamodb.Item{
// 		"AccountID":          &types.AttributeValueMemberS{Value: "account1"},
// 		"CumulativePayments": &types.AttributeValueMemberS{Value: "3"},
// 		"BlobSize":           &types.AttributeValueMemberN{Value: "3000"},
// 	})
// 	assert.NoError(t, err)
// 	assert.Equal(t, "3000", updatedItem["BlobSize"].(*types.AttributeValueMemberN).Value)

// 	item, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
// 		"AccountID":          &types.AttributeValueMemberS{Value: "account1"},
// 		"CumulativePayments": &types.AttributeValueMemberS{Value: "1"},
// 	})
// 	fmt.Println(item)
// 	fmt.Println(item["AccountID"].(*types.AttributeValueMemberS).Value)
// 	fmt.Println(item["CumulativePayments"].(*types.AttributeValueMemberS).Value)
// 	fmt.Println(item["BlobSize"].(*types.AttributeValueMemberN).Value)
// 	assert.NoError(t, err)
// 	assert.Equal(t, "3000", item["BlobSize"].(*types.AttributeValueMemberN).Value)

// 	// item, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
// 	// 	"AccountID":          &types.AttributeValueMemberS{Value: "account1"},
// 	// 	"CumulativePayments": &types.AttributeValueMemberS{Value: "3"},
// 	// })
// 	// fmt.Println(err)
// 	// fmt.Println(item)
// 	// assert.Error(t, err)
// 	// assert.Equal(t, "account1", item["AccountID"].(*types.AttributeValueMemberS).Value)
// 	// assert.Equal(t, "3", item["CumulativePayments"].(*types.AttributeValueMemberS).Value)
// 	// assert.Equal(t, "3000", item["BlobSize"].(*types.AttributeValueMemberN).Value)
// 	// item is nil
// 	// err = dynamoClient.PutItem(ctx, tableName,
// 	// 	commondynamodb.Item{
// 	// 		"AccountID":          &types.AttributeValueMemberS{Value: "account2"},
// 	// 		"CumulativePayments": &types.AttributeValueMemberN{Value: "3000"},
// 	// 		"BlobSize":           &types.AttributeValueMemberN{Value: "3000"},
// 	// 		"UpdatedAt":          &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
// 	// 	},
// 	// )
// 	// assert.NoError(t, err)

// 	// queryResult, err = dynamoClient.QueryIndex(ctx, tableName, indexName, "AccountID = :account_id", commondynamodb.ExpresseionValues{
// 	// 	":account_id": &types.AttributeValueMemberS{
// 	// 		Value: "account1",
// 	// 	}})
// 	// assert.NoError(t, err)
// 	// assert.Len(t, queryResult, 1)
// 	// assert.Equal(t, "1", queryResult[0]["CumulativePayments"].(*types.AttributeValueMemberN).Value)
// 	// assert.Equal(t, "2000", queryResult[0]["BlobSize"].(*types.AttributeValueMemberN).Value)

// 	// queryResult, err = dynamoClient.QueryIndex(ctx, tableName, indexName, "AccountID = :account_id", commondynamodb.ExpresseionValues{
// 	// 	":account_id": &types.AttributeValueMemberS{
// 	// 		Value: "account2",
// 	// 	}})
// 	// assert.NoError(t, err)
// 	// assert.Len(t, queryResult, 1)
// 	// assert.Equal(t, "2", queryResult[0]["CumulativePayments"].(*types.AttributeValueMemberN).Value)
// 	// assert.Equal(t, "3000", queryResult[0]["BlobSize"].(*types.AttributeValueMemberN).Value)

// 	// // items, err = dynamoClient.QueryIndex(ctx, tableName, indexName, "BinIndex = :index", commondynamodb.ExpresseionValues{
// 	// // 	":index": &types.AttributeValueMemberN{
// 	// // 		Value: "1",
// 	// // 	},
// 	// // 	""
// 	// // })
// 	// // assert.NoError(t, err)
// 	// // assert.Equal(t, "2000", items[0]["BinUsage"].(*types.AttributeValueMemberN).Value)
// 	// // assert.Equal(t, "3000", items[1]["BinUsage"].(*types.AttributeValueMemberN).Value)
// }
