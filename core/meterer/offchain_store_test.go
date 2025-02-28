package meterer_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

func TestReservationBinsBasicOperations(t *testing.T) {
	tableName := "reservations_test_basic"
	err := meterer.CreateReservationTable(clientConfig, tableName)
	assert.NoError(t, err)

	ctx := context.Background()
	err = dynamoClient.PutItem(ctx, tableName,
		commondynamodb.Item{
			"AccountID":         &types.AttributeValueMemberS{Value: "account1"},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: "1"},
			"BinUsage":          &types.AttributeValueMemberN{Value: "1000"},
			"UpdatedAt":         &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
	)
	assert.NoError(t, err)

	item, err := dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: "account1"},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: "1"},
	})
	assert.NoError(t, err)

	assert.Equal(t, "account1", item["AccountID"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "1", item["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
	assert.Equal(t, "1000", item["BinUsage"].(*types.AttributeValueMemberN).Value)

	items, err := dynamoClient.Query(ctx, tableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{Value: "account1"},
	})
	assert.NoError(t, err)
	assert.Len(t, items, 1)

	_, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: "account2"},
	})
	assert.Error(t, err)

	_, err = dynamoClient.UpdateItem(ctx, tableName, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: "account1"},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: "1"},
	}, commondynamodb.Item{
		"BinUsage": &types.AttributeValueMemberN{Value: "2000"},
	})
	assert.NoError(t, err)
	err = dynamoClient.PutItem(ctx, tableName,
		commondynamodb.Item{
			"AccountID":         &types.AttributeValueMemberS{Value: "account2"},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: "1"},
			"BinUsage":          &types.AttributeValueMemberN{Value: "3000"},
			"UpdatedAt":         &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
	)
	assert.NoError(t, err)

	item, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: "account1"},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: "1"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "2000", item["BinUsage"].(*types.AttributeValueMemberN).Value)

	items, err = dynamoClient.Query(ctx, tableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{Value: "account1"},
	})
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "2000", items[0]["BinUsage"].(*types.AttributeValueMemberN).Value)

	item, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"AccountID":         &types.AttributeValueMemberS{Value: "account2"},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: "1"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "3000", item["BinUsage"].(*types.AttributeValueMemberN).Value)

	err = dynamoClient.DeleteTable(ctx, tableName)
	assert.NoError(t, err)
}

func TestGlobalBinsBasicOperations(t *testing.T) {
	tableName := "global_test_basic"
	err := meterer.CreateGlobalReservationTable(clientConfig, tableName)
	assert.NoError(t, err)

	ctx := context.Background()
	numItems := 30
	items := make([]commondynamodb.Item, numItems)
	for i := 0; i < numItems; i += 1 {
		items[i] = commondynamodb.Item{
			"ReservationPeriod": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", i)},
			"BinUsage":          &types.AttributeValueMemberN{Value: "1000"},
			"UpdatedAt":         &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		}
	}
	unprocessed, err := dynamoClient.PutItems(ctx, tableName, items)
	assert.NoError(t, err)
	assert.Len(t, unprocessed, 0)

	queryResult, err := dynamoClient.Query(ctx, tableName, "ReservationPeriod = :index", commondynamodb.ExpressionValues{
		":index": &types.AttributeValueMemberN{
			Value: "1",
		}})
	assert.NoError(t, err)
	assert.Len(t, queryResult, 1)
	assert.Equal(t, "1", queryResult[0]["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
	assert.Equal(t, "1000", queryResult[0]["BinUsage"].(*types.AttributeValueMemberN).Value)

	queryResult, err = dynamoClient.Query(ctx, tableName, "ReservationPeriod = :index", commondynamodb.ExpressionValues{
		":index": &types.AttributeValueMemberN{
			Value: "1",
		}})
	assert.NoError(t, err)
	assert.Len(t, queryResult, 1)
	assert.Equal(t, "1", queryResult[0]["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
	assert.Equal(t, "1000", queryResult[0]["BinUsage"].(*types.AttributeValueMemberN).Value)

	queryResult, err = dynamoClient.Query(ctx, tableName, "ReservationPeriod = :index", commondynamodb.ExpressionValues{
		":index": &types.AttributeValueMemberN{
			Value: "32",
		}})
	assert.NoError(t, err)
	assert.Len(t, queryResult, 0)

	_, err = dynamoClient.UpdateItem(ctx, tableName, commondynamodb.Key{
		"ReservationPeriod": &types.AttributeValueMemberN{Value: "1"},
	}, commondynamodb.Item{
		"BinUsage": &types.AttributeValueMemberN{Value: "2000"},
	})
	assert.NoError(t, err)

	err = dynamoClient.PutItem(ctx, tableName,
		commondynamodb.Item{
			"ReservationPeriod": &types.AttributeValueMemberN{Value: "2"},
			"BinUsage":          &types.AttributeValueMemberN{Value: "3000"},
			"UpdatedAt":         &types.AttributeValueMemberS{Value: time.Now().Format(time.RFC3339)},
		},
	)
	assert.NoError(t, err)

	queryResult, err = dynamoClient.Query(ctx, tableName, "ReservationPeriod = :index", commondynamodb.ExpressionValues{
		":index": &types.AttributeValueMemberN{
			Value: "1",
		}})
	assert.NoError(t, err)
	assert.Len(t, queryResult, 1)
	assert.Equal(t, "1", queryResult[0]["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
	assert.Equal(t, "2000", queryResult[0]["BinUsage"].(*types.AttributeValueMemberN).Value)

	queryResult, err = dynamoClient.Query(ctx, tableName, "ReservationPeriod = :index", commondynamodb.ExpressionValues{
		":index": &types.AttributeValueMemberN{
			Value: "2",
		}})
	assert.NoError(t, err)
	assert.Len(t, queryResult, 1)
	assert.Equal(t, "2", queryResult[0]["ReservationPeriod"].(*types.AttributeValueMemberN).Value)
	assert.Equal(t, "3000", queryResult[0]["BinUsage"].(*types.AttributeValueMemberN).Value)
}

func TestOnDemandUsageBasicOperations(t *testing.T) {
	tableName := "ondemand_test_basic"
	err := meterer.CreateOnDemandTable(clientConfig, tableName)
	assert.NoError(t, err)

	ctx := context.Background()

	err = dynamoClient.PutItem(ctx, tableName,
		commondynamodb.Item{
			"AccountID":          &types.AttributeValueMemberS{Value: "account1"},
			"CumulativePayments": &types.AttributeValueMemberN{Value: "1"},
			"DataLength":         &types.AttributeValueMemberN{Value: "1000"},
		},
	)
	assert.NoError(t, err)

	numItems := 30
	repetitions := 5
	items := make([]commondynamodb.Item, numItems)
	for i := 0; i < numItems; i += 1 {
		items[i] = commondynamodb.Item{
			"AccountID":          &types.AttributeValueMemberS{Value: fmt.Sprintf("account%d", i%repetitions)},
			"CumulativePayments": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", i)},
			"DataLength":         &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", i*1000)},
		}
	}
	unprocessed, err := dynamoClient.PutItems(ctx, tableName, items)
	assert.NoError(t, err)
	assert.Len(t, unprocessed, 0)

	// get item
	item, err := dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"AccountID":          &types.AttributeValueMemberS{Value: "account1"},
		"CumulativePayments": &types.AttributeValueMemberN{Value: "1"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "1", item["CumulativePayments"].(*types.AttributeValueMemberN).Value)
	assert.Equal(t, "1000", item["DataLength"].(*types.AttributeValueMemberN).Value)

	queryResult, err := dynamoClient.Query(ctx, tableName, "AccountID = :account", commondynamodb.ExpressionValues{
		":account": &types.AttributeValueMemberS{
			Value: "account1",
		}})
	assert.NoError(t, err)
	assert.Len(t, queryResult, numItems/repetitions)
	for _, item := range queryResult {
		cumulativePayments, _ := strconv.Atoi(item["CumulativePayments"].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, fmt.Sprintf("%d", cumulativePayments*1000), item["DataLength"].(*types.AttributeValueMemberN).Value)
	}
	queryResult, err = dynamoClient.Query(ctx, tableName, "AccountID = :account_id", commondynamodb.ExpressionValues{
		":account_id": &types.AttributeValueMemberS{
			Value: fmt.Sprintf("account%d", numItems/repetitions+1),
		}})
	assert.NoError(t, err)
	assert.Len(t, queryResult, 0)

	updatedItem, err := dynamoClient.UpdateItem(ctx, tableName, commondynamodb.Key{
		"AccountID":          &types.AttributeValueMemberS{Value: "account1"},
		"CumulativePayments": &types.AttributeValueMemberN{Value: "1"},
		// "BinUsage": &types.AttributeValueMemberN{Value: "1000"},
	}, commondynamodb.Item{
		"AccountID":          &types.AttributeValueMemberS{Value: "account1"},
		"CumulativePayments": &types.AttributeValueMemberN{Value: "3"},
		"DataLength":         &types.AttributeValueMemberN{Value: "3000"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "3000", updatedItem["DataLength"].(*types.AttributeValueMemberN).Value)

	item, err = dynamoClient.GetItem(ctx, tableName, commondynamodb.Key{
		"AccountID":          &types.AttributeValueMemberS{Value: "account1"},
		"CumulativePayments": &types.AttributeValueMemberN{Value: "1"},
	})
	assert.NoError(t, err)
	assert.Equal(t, "3000", item["DataLength"].(*types.AttributeValueMemberN).Value)
}
