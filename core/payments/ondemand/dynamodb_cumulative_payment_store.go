package ondemand

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
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
// This implementation stores cumulative payment values in a DynamoDB table, providing
// persistent storage for on-demand payment tracking. The table uses AccountID as the
// partition key and stores the CumulativePayment value as a number.
type DynamoDBCumulativePaymentStore struct {
	logger       logging.Logger
	dynamoClient dynamodb.Client
	tableName    string
	accountID    gethcommon.Address
}

var _ CumulativePaymentStore = (*DynamoDBCumulativePaymentStore)(nil)

// Creates a new DynamoDB-backed cumulative payment store
func NewDynamoDBCumulativePaymentStore(
	logger logging.Logger,
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
		logger:       logger,
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
		s.logger.Debug("No payment record found for account, returning 0", "accountID", s.accountID.Hex())
		return big.NewInt(0), nil
	}

	// Extract CumulativePayment attribute
	paymentAttr, ok := result[attributeCumulativePayment]
	if !ok {
		s.logger.Debugf("%s attribute not found, returning 0 for accountID %s", attributeCumulativePayment, s.accountID.Hex())
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

	s.logger.Debug("Retrieved cumulative payment",
		"accountID", s.accountID.Hex(),
		"payment", payment.String())

	return payment, nil
}

// SetCumulativePayment stores the new cumulative payment value in DynamoDB
//
// Returns an error if:
// - The new payment is not greater than the existing payment
// - There's a failure writing to DynamoDB
func (s *DynamoDBCumulativePaymentStore) SetCumulativePayment(ctx context.Context, newCumulativePayment *big.Int) error {
	if newCumulativePayment == nil {
		return fmt.Errorf("newCumulativePayment cannot be nil")
	}

	// Create the item to store
	item := dynamodb.Item{
		attributeAccountID:         &types.AttributeValueMemberS{Value: s.accountID.Hex()},
		attributeCumulativePayment: &types.AttributeValueMemberN{Value: newCumulativePayment.String()},
	}

	exprValueNewPayment := ":newPayment"

	// Use conditional expression to ensure:
	// 1. If no record exists, accept the payment (first payment for this account)
	// 2. If record exists, only accept if new payment is greater than existing
	// This ensures cumulative payments only increase, preventing concurrent update issues
	conditionExpression := "attribute_not_exists(" + attributeCumulativePayment + ") OR " +
		attributeCumulativePayment + " < " + exprValueNewPayment

	expressionValues := map[string]types.AttributeValue{
		exprValueNewPayment: &types.AttributeValueMemberN{Value: newCumulativePayment.String()},
	}

	err := s.dynamoClient.PutItemWithCondition(
		ctx,
		s.tableName,
		item,
		conditionExpression,
		nil, // No expression attribute names needed
		expressionValues,
	)

	if err != nil {
		if errors.Is(err, dynamodb.ErrConditionFailed) {
			return fmt.Errorf("new value (%s) must be greater than existing value", newCumulativePayment.String())
		}
		return fmt.Errorf("set payment for account %s: %w", s.accountID.Hex(), err)
	}

	s.logger.Debug("Set cumulative payment",
		"accountID", s.accountID.Hex(),
		"payment", newCumulativePayment.String())

	return nil
}
