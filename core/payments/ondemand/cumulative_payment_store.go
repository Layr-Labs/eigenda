package ondemand

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	attributeAccountID         = "AccountID"
	attributeCumulativePayment = "CumulativePayment"
)

// CumulativePaymentStore provides persistent storage for cumulative payment values using DynamoDB.
//
// The table uses AccountID as the partition key and stores the CumulativePayment value as a number.
//
// This store represents a subset of the logic implemented in [meterer.DynamoDBMeteringStore]. It maintains the same
// table structure for the sake of backwards compatibility, but otherwise is intended to replace the old class, as
// part of the ongoing payments refactor.
//
// TODO(litt3): there are some potential avenues for optimization of this store:
// 1. Use something other than DynamoDB. DynamoDB is being used for historical reasons, but there is only a single
// writer now, which doesn't need any of the distributed DB properties provided by DynamoDB.
// 2. Implement a write queue, so that the caller doesn't need to wait for the write to complete. The callers of the
// CumulativePaymentStore just need *eventual* persistence of the cumulative payment, so using a queue would be
// sufficient, and would free the caller from blocking on I/O. Note that this optimization would make undercharging
// a possibility, if a crash happens before a piece of usage data has been persisted. This is an acceptable
// tradeoff for simplified architecture and improved performance.
type CumulativePaymentStore struct {
	// The DynamoDB client to use for storage operations
	dynamoClient *dynamodb.Client
	// The name of the DynamoDB table to store payments in, stored as *string for use in DynamoDB operations
	tableName *string
	// The account address, pre-built as a key for DynamoDB operations
	accountKey map[string]types.AttributeValue
}

// Creates a new DynamoDB-backed cumulative payment store
func NewCumulativePaymentStore(
	dynamoClient *dynamodb.Client,
	tableName string,
	// The account ID this store is tracking payments for
	accountID gethcommon.Address,
) (*CumulativePaymentStore, error) {
	if dynamoClient == nil {
		return nil, fmt.Errorf("dynamoClient cannot be nil")
	}
	if tableName == "" {
		return nil, fmt.Errorf("tableName cannot be empty")
	}
	if accountID == (gethcommon.Address{}) {
		return nil, fmt.Errorf("accountID cannot be the zero address")
	}

	return &CumulativePaymentStore{
		dynamoClient: dynamoClient,
		tableName:    aws.String(tableName),
		accountKey: map[string]types.AttributeValue{
			attributeAccountID: &types.AttributeValueMemberS{Value: accountID.Hex()},
		},
	}, nil
}

// Stores a new cumulative payment value in DynamoDB
func (s *CumulativePaymentStore) StoreCumulativePayment(
	ctx context.Context,
	newCumulativePayment *big.Int,
) error {
	if s == nil {
		// sane no-op behavior, since using a payment store is optional
		return nil
	}

	if newCumulativePayment == nil {
		return errors.New("newCumulativePayment cannot be nil")
	}
	if newCumulativePayment.Sign() < 0 {
		return fmt.Errorf("cumulative payment cannot be negative: received %s", newCumulativePayment.String())
	}

	_, err := s.dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        s.tableName,
		Key:              s.accountKey,
		UpdateExpression: aws.String("SET #cp = :new"),
		ExpressionAttributeNames: map[string]string{
			"#cp": attributeCumulativePayment,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":new": &types.AttributeValueMemberN{Value: newCumulativePayment.String()},
		},
	})
	if err != nil {
		return fmt.Errorf("update cumulative payment: %w", err)
	}

	return nil
}

// Retrieves the current cumulative payment value from DynamoDB
func (s *CumulativePaymentStore) GetCumulativePayment(ctx context.Context) (*big.Int, error) {
	resp, err := s.dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:            s.tableName,
		Key:                  s.accountKey,
		ConsistentRead:       aws.Bool(true),
		ProjectionExpression: aws.String(attributeCumulativePayment),
	})
	if err != nil {
		return nil, fmt.Errorf("get item: %w", err)
	}

	if len(resp.Item) == 0 {
		return big.NewInt(0), nil
	}

	attributeValue, ok := resp.Item[attributeCumulativePayment]
	if !ok {
		return big.NewInt(0), nil
	}

	attributeNumber, ok := attributeValue.(*types.AttributeValueMemberN)
	if !ok {
		return nil, fmt.Errorf("%s has invalid type: %T", attributeCumulativePayment, attributeValue)
	}

	cumulativePayment := new(big.Int)
	if _, success := cumulativePayment.SetString(attributeNumber.Value, 10); !success {
		return nil, fmt.Errorf("parse cumulative payment value: %s", attributeNumber.Value)
	}

	return cumulativePayment, nil
}
