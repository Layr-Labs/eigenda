package ondemand

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	awsdynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	attributeAccountID         = "AccountID"
	attributeCumulativePayment = "CumulativePayment"
)

// DynamoDBCumulativePaymentStore implements the CumulativePaymentStore interface using DynamoDB
//
// This struct does NOT support decrementing cumulative payments. It is designed to exist on the disperser, where
// it doesn't make sense for cumulative payment to ever decrease. Therefore, for extra safety, decreasing cumulative
// payment is forbidden.
//
// This implementation provides persistent storage for on-demand payment tracking. The table uses AccountID as the
// partition key and stores the CumulativePayment value as a number.
//
// This store represents a subset of the logic implemented in [meterer.DynamoDBMeteringStore]. It maintains the same
// table structure for the sake of backwards compatibility, but otherwise is intended to replace the old class, as
// part of the ongoing payments refactor.
type DynamoDBCumulativePaymentStore struct {
	dynamoClient dynamodb.Client
	tableName    string
	accountID    gethcommon.Address
}

var _ CumulativePaymentStore = (*DynamoDBCumulativePaymentStore)(nil)

// Creates a new DynamoDB-backed cumulative payment store
func NewDynamoDBCumulativePaymentStore(
	// The DynamoDB client to use for storage operations
	dynamoClient dynamodb.Client,
	// The name of the DynamoDB table to store payments in
	tableName string,
	// The account ID this store is tracking payments for
	accountID gethcommon.Address,
) (*DynamoDBCumulativePaymentStore, error) {
	if dynamoClient == nil {
		return nil, fmt.Errorf("dynamoClient cannot be nil")
	}
	if tableName == "" {
		return nil, fmt.Errorf("tableName cannot be empty")
	}
	if accountID == (gethcommon.Address{}) {
		return nil, fmt.Errorf("accountID cannot be the zero address")
	}

	return &DynamoDBCumulativePaymentStore{
		dynamoClient: dynamoClient,
		tableName:    tableName,
		accountID:    accountID,
	}, nil
}

// GetCumulativePayment retrieves the stored cumulative payment for the account from DynamoDB
//
// If no payment record exists for the account, returns 0.
// Returns an error if there's a failure accessing DynamoDB or parsing the stored value.
func (s *DynamoDBCumulativePaymentStore) GetCumulativePayment(ctx context.Context) (*big.Int, error) {
	input := &awsdynamodb.GetItemInput{
		TableName: aws.String(s.tableName),
		Key: map[string]types.AttributeValue{
			attributeAccountID: &types.AttributeValueMemberS{Value: s.accountID.Hex()},
		},
		// Use strongly consistent read to ensure we get the latest value
		ConsistentRead: aws.Bool(true),
	}

	result, err := s.dynamoClient.GetItemWithInput(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment for account %s: %w", s.accountID.Hex(), err)
	}

	// If no item found, return zero (new account)
	if len(result) == 0 {
		return big.NewInt(0), nil
	}

	// Extract CumulativePayment attribute
	paymentAttr, ok := result[attributeCumulativePayment]
	if !ok {
		return big.NewInt(0), nil
	}

	// Type assertion to ensure it's a number
	paymentNumber, ok := paymentAttr.(*types.AttributeValueMemberN)
	if !ok {
		return nil, fmt.Errorf("%s has invalid type: %T", attributeCumulativePayment, paymentAttr)
	}

	// Parse the string value to big.Int
	payment := new(big.Int)
	if _, success := payment.SetString(paymentNumber.Value, 10); !success {
		return nil, fmt.Errorf("parse payment value: %s", paymentNumber.Value)
	}

	return payment, nil
}

// SetCumulativePayment stores the new cumulative payment value in DynamoDB
//
// The operation is idempotent.
//
// Returns an error if:
// - The new payment is less than the existing payment (decrements are not allowed)
// - There's a failure writing to DynamoDB
func (s *DynamoDBCumulativePaymentStore) SetCumulativePayment(
	ctx context.Context,
	newCumulativePayment *big.Int,
) error {
	if newCumulativePayment == nil {
		return fmt.Errorf("newCumulativePayment cannot be nil")
	}
	if newCumulativePayment.Sign() < 0 {
		return fmt.Errorf("newCumulativePayment cannot be negative: %s", newCumulativePayment.String())
	}

	key := dynamodb.Key{
		attributeAccountID: &types.AttributeValueMemberS{Value: s.accountID.Hex()},
	}

	item := dynamodb.Item{
		attributeCumulativePayment: &types.AttributeValueMemberN{Value: newCumulativePayment.String()},
	}

	// 1. If no record exists, accept the payment (first payment for this account)
	// 2. If record exists, only accept if new payment is greater than or equal to existing
	//    (allowing idempotent retries with the same value)
	conditionBuilder := expression.Or(
		expression.AttributeNotExists(expression.Name(attributeCumulativePayment)),
		expression.LessThanEqual(
			expression.Name(attributeCumulativePayment),
			expression.Value(&types.AttributeValueMemberN{Value: newCumulativePayment.String()}),
		),
	)

	_, err := s.dynamoClient.UpdateItemWithCondition(
		ctx,
		s.tableName,
		key,
		item,
		conditionBuilder,
	)

	if err != nil {
		if errors.Is(err, dynamodb.ErrConditionFailed) {
			return fmt.Errorf(
				"new value (%s) must be greater than or equal to existing value", newCumulativePayment.String())
		}
		return fmt.Errorf("set payment for account %s: %w", s.accountID.Hex(), err)
	}

	return nil
}
