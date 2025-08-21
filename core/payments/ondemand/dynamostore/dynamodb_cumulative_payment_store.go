package dynamostore

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/Layr-Labs/eigenda/core/payments/ondemand"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

const (
	attributeAccountID         = "AccountID"
	attributeCumulativePayment = "CumulativePayment"
)

// DynamoDBCumulativePaymentStore implements the CumulativePaymentStore interface using DynamoDB
//
// This struct intentionally does NOT support decrementing cumulative payments. It is designed to exist on the
// disperser, where it doesn't make sense for cumulative payment to ever decrease. Therefore, for extra safety,
// decreasing cumulative payment is forbidden.
//
// This implementation provides persistent storage for on-demand payment tracking. The table uses AccountID as the
// partition key and stores the CumulativePayment value as a number.
//
// This store represents a subset of the logic implemented in [meterer.DynamoDBMeteringStore]. It maintains the same
// table structure for the sake of backwards compatibility, but otherwise is intended to replace the old class, as
// part of the ongoing payments refactor.
type DynamoDBCumulativePaymentStore struct {
	// The DynamoDB client to use for storage operations
	dynamoClient *dynamodb.Client
	// The name of the DynamoDB table to store payments in, stored as *string for use in DynamoDB operations
	tableName *string
	// The account address, pre-built as a key for DynamoDB operations
	accountKey map[string]types.AttributeValue
}

var _ ondemand.CumulativePaymentStore = (*DynamoDBCumulativePaymentStore)(nil)

// Creates a new DynamoDB-backed cumulative payment store
func NewDynamoDBCumulativePaymentStore(
	dynamoClient *dynamodb.Client,
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
		tableName:    aws.String(tableName),
		accountKey: map[string]types.AttributeValue{
			attributeAccountID: &types.AttributeValueMemberS{Value: accountID.Hex()},
		},
	}, nil
}

// AddCumulativePayment atomically adds a given amount to the cumulative payment
//
// This method uses DynamoDB's conditional update feature to ensure thread-safe operations.
// It performs an atomic read-modify-write operation that:
// 1. Adds the amount to the existing cumulative payment (or initializes to the amount if no record exists)
// 2. Ensures the resulting value doesn't exceed maxCumulativePayment
//
// Returns the new cumulative payment value after the addition.
// Returns an InsufficientFundsError if the addition would exceed maxCumulativePayment.
func (s *DynamoDBCumulativePaymentStore) AddCumulativePayment(
	ctx context.Context,
	amount *big.Int,
	maxCumulativePayment *big.Int,
) (*big.Int, error) {
	if amount == nil {
		return nil, errors.New("amount cannot be nil")
	}
	if amount.Sign() <= 0 {
		return nil, fmt.Errorf("amount must be positive (decrementing not supported): received %s", amount.String())
	}

	if maxCumulativePayment == nil {
		return nil, errors.New("maxCumulativePayment cannot be nil")
	}
	if maxCumulativePayment.Sign() <= 0 {
		return nil, fmt.Errorf("maxCumulativePayment must be positive: received %s", maxCumulativePayment.String())
	}

	// Calculate the maximum allowed current cumulative payment value.
	// If the cumulativePayment value that exists prior to the addition is greater than this, that means the account
	// doesn't have enough funds to pay
	maxAllowedCurrent := new(big.Int).Sub(maxCumulativePayment, amount)

	// if maxAllowedCurrent is negative, that means the cost of this single dispersal alone exceeds total deposits
	if maxAllowedCurrent.Sign() < 0 {
		// we're passing in nil as the CurrentCumulativePayment as a shortcut: the blob cost alone is greater than
		// the maxCumulativePayment, so CurrentCumulativePayment is irrelevant. It's not worth implementing extra
		// fetching logic, just for this strange edge case.
		return nil, &ondemand.InsufficientFundsError{
			CurrentCumulativePayment: nil,
			MaxCumulativePayment:     maxCumulativePayment,
			BlobCost:                 amount,
		}
	}

	// Build the update expression that performs the atomic addition.
	// This expression does: CumulativePayment = (CumulativePayment || 0) + amount
	// - if_not_exists(CumulativePayment, :zero) returns the current value if it exists, or 0 if it doesn't
	// - Then we add :inc (the amount) to that value
	// - The result is stored back in CumulativePayment
	updateExpression := fmt.Sprintf("SET %s = if_not_exists(%s, :zero) + :inc",
		attributeCumulativePayment, attributeCumulativePayment)

	// Build the condition that must be true for the update to succeed.
	// This ensures we don't exceed the maximum allowed cumulative payment.
	// The condition is: (attribute doesn't exist) OR (current value <= maxAllowedCurrent)
	// - For new accounts: attribute_not_exists is true, so condition passes.
	//   This is safe because we've already verified above that amount alone doesn't exceed maxCumulativePayment
	// - For existing accounts: current value must be <= maxAllowedCurrent
	//   This ensures that current + amount <= maxCumulativePayment
	conditionExpression := fmt.Sprintf("attribute_not_exists(%s) OR %s <= :max",
		attributeCumulativePayment, attributeCumulativePayment)

	// Define the placeholder values used in the expressions above
	expressionAttributeValues := map[string]types.AttributeValue{
		// Used when initializing a new account
		":zero": &types.AttributeValueMemberN{Value: "0"},
		// The amount to add to cumulative payment
		":inc": &types.AttributeValueMemberN{Value: amount.String()},
		// Maximum allowed current value
		":max": &types.AttributeValueMemberN{Value: maxAllowedCurrent.String()},
	}

	updateItemInput := &dynamodb.UpdateItemInput{
		TableName:                           s.tableName,
		Key:                                 s.accountKey,
		UpdateExpression:                    aws.String(updateExpression),
		ConditionExpression:                 aws.String(conditionExpression),
		ExpressionAttributeValues:           expressionAttributeValues,
		ReturnValues:                        types.ReturnValueUpdatedNew,
		ReturnValuesOnConditionCheckFailure: types.ReturnValuesOnConditionCheckFailureAllOld,
	}

	updateItemOutput, err := s.dynamoClient.UpdateItem(ctx, updateItemInput)
	if err != nil {
		var conditionCheckFailedException *types.ConditionalCheckFailedException
		if errors.As(err, &conditionCheckFailedException) {
			currentValue, parseErr := extractPaymentValue(conditionCheckFailedException.Item)
			if parseErr != nil {
				return nil, fmt.Errorf(
					"conditional check failed with error %w, and couldn't parse current value: %w", err, parseErr)
			}
			return nil, &ondemand.InsufficientFundsError{
				CurrentCumulativePayment: currentValue,
				MaxCumulativePayment:     maxCumulativePayment,
				BlobCost:                 amount,
			}
		}

		return nil, fmt.Errorf("update cumulative payment: %w", err)
	}

	return extractPaymentValue(updateItemOutput.Attributes)
}

// extractPaymentValue extracts and parses attributeCumulativePayment from a DynamoDB item
func extractPaymentValue(item map[string]types.AttributeValue) (*big.Int, error) {
	if len(item) == 0 {
		return big.NewInt(0), nil
	}

	attributeValue, ok := item[attributeCumulativePayment]
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
