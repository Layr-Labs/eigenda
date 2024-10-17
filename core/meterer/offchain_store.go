package meterer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type OffchainStore struct {
	dynamoClient         *commondynamodb.Client
	reservationTableName string
	onDemandTableName    string
	globalBinTableName   string
	logger               logging.Logger
	// TODO: add maximum storage for both tables
}

func NewOffchainStore(
	cfg commonaws.ClientConfig,
	reservationTableName string,
	onDemandTableName string,
	globalBinTableName string,
	logger logging.Logger,
) (OffchainStore, error) {

	dynamoClient, err := commondynamodb.NewClient(cfg, logger)
	if err != nil {
		return OffchainStore{}, err
	}

	err = dynamoClient.TableExists(context.Background(), reservationTableName)
	if err != nil {
		return OffchainStore{}, err
	}
	err = dynamoClient.TableExists(context.Background(), onDemandTableName)
	if err != nil {
		return OffchainStore{}, err
	}
	err = dynamoClient.TableExists(context.Background(), globalBinTableName)
	if err != nil {
		return OffchainStore{}, err
	}
	//TODO: add a separate thread to periodically clean up the tables
	// delete expired reservation bins (<i-1) and old on-demand payments (retain max N payments)
	return OffchainStore{
		dynamoClient:         dynamoClient,
		reservationTableName: reservationTableName,
		onDemandTableName:    onDemandTableName,
		globalBinTableName:   globalBinTableName,
		logger:               logger,
	}, nil
}

type ReservationBin struct {
	AccountID string
	BinIndex  uint32
	BinUsage  uint32
	UpdatedAt time.Time
}

type PaymentTuple struct {
	CumulativePayment uint64
	DataLength        uint32
}

type GlobalBin struct {
	BinIndex  uint32
	BinUsage  uint64
	UpdatedAt time.Time
}

func (s *OffchainStore) UpdateReservationBin(ctx context.Context, accountID string, binIndex uint64, size uint64) (uint64, error) {
	key := map[string]types.AttributeValue{
		"AccountID": &types.AttributeValueMemberS{Value: accountID},
		"BinIndex":  &types.AttributeValueMemberN{Value: strconv.FormatUint(binIndex, 10)},
	}

	res, err := s.dynamoClient.IncrementBy(ctx, s.reservationTableName, key, "BinUsage", size)
	if err != nil {
		return 0, fmt.Errorf("failed to increment bin usage: %w", err)
	}

	binUsage, ok := res["BinUsage"]
	if !ok {
		return 0, errors.New("BinUsage is not present in the response")
	}

	binUsageAttr, ok := binUsage.(*types.AttributeValueMemberN)
	if !ok {
		return 0, fmt.Errorf("unexpected type for BinUsage: %T", binUsage)
	}

	binUsageValue, err := strconv.ParseUint(binUsageAttr.Value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse BinUsage: %w", err)
	}

	return binUsageValue, nil
}

func (s *OffchainStore) UpdateGlobalBin(ctx context.Context, binIndex uint64, size uint64) (uint64, error) {
	key := map[string]types.AttributeValue{
		"BinIndex": &types.AttributeValueMemberN{Value: strconv.FormatUint(binIndex, 10)},
	}

	res, err := s.dynamoClient.IncrementBy(ctx, s.globalBinTableName, key, "BinUsage", size)
	if err != nil {
		return 0, err
	}

	binUsage, ok := res["BinUsage"]
	if !ok {
		return 0, nil
	}

	binUsageAttr, ok := binUsage.(*types.AttributeValueMemberN)
	if !ok {
		return 0, nil
	}

	binUsageValue, err := strconv.ParseUint(binUsageAttr.Value, 10, 32)
	if err != nil {
		return 0, err
	}

	return binUsageValue, nil
}

func (s *OffchainStore) AddOnDemandPayment(ctx context.Context, paymentMetadata core.PaymentMetadata, symbolsCharged uint32) error {
	result, err := s.dynamoClient.GetItem(ctx, s.onDemandTableName,
		commondynamodb.Item{
			"AccountID":          &types.AttributeValueMemberS{Value: paymentMetadata.AccountID},
			"CumulativePayments": &types.AttributeValueMemberN{Value: paymentMetadata.CumulativePayment.String()},
		},
	)
	if err != nil {
		fmt.Println("new payment record: %w", err)
	}
	if result != nil {
		return fmt.Errorf("exact payment already exists")
	}
	err = s.dynamoClient.PutItem(ctx, s.onDemandTableName,
		commondynamodb.Item{
			"AccountID":          &types.AttributeValueMemberS{Value: paymentMetadata.AccountID},
			"CumulativePayments": &types.AttributeValueMemberN{Value: paymentMetadata.CumulativePayment.String()},
			"DataLength":         &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(symbolsCharged), 10)},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to add payment: %w", err)
	}
	return nil
}

// RemoveOnDemandPayment removes a specific payment from the list for a specific account
func (s *OffchainStore) RemoveOnDemandPayment(ctx context.Context, accountID string, payment *big.Int) error {
	err := s.dynamoClient.DeleteItem(ctx, s.onDemandTableName,
		commondynamodb.Key{
			"AccountID":          &types.AttributeValueMemberS{Value: accountID},
			"CumulativePayments": &types.AttributeValueMemberN{Value: payment.String()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to remove payment: %w", err)
	}

	return nil
}

// GetRelevantOnDemandRecords gets previous cumulative payment, next cumulative payment, blob size of next payment
// The queries are done sequentially instead of one-go for efficient querying and would not cause race condition errors for honest requests
func (s *OffchainStore) GetRelevantOnDemandRecords(ctx context.Context, accountID string, cumulativePayment *big.Int) (uint64, uint64, uint32, error) {
	// Fetch the largest entry smaller than the given cumulativePayment
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(s.onDemandTableName),
		KeyConditionExpression: aws.String("AccountID = :account AND CumulativePayments < :cumulativePayment"),
		ExpressionAttributeValues: commondynamodb.ExpressionValues{
			":account":           &types.AttributeValueMemberS{Value: accountID},
			":cumulativePayment": &types.AttributeValueMemberN{Value: cumulativePayment.String()},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(1),
	}
	smallerResult, err := s.dynamoClient.QueryWithInput(ctx, queryInput)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to query smaller payments for account: %w", err)
	}
	var prevPayment uint64
	if len(smallerResult) > 0 {
		prevPayment, err = strconv.ParseUint(smallerResult[0]["CumulativePayments"].(*types.AttributeValueMemberN).Value, 10, 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("failed to parse previous payment: %w", err)
		}
	}

	// Fetch the smallest entry larger than the given cumulativePayment
	queryInput = &dynamodb.QueryInput{
		TableName:              aws.String(s.onDemandTableName),
		KeyConditionExpression: aws.String("AccountID = :account AND CumulativePayments > :cumulativePayment"),
		ExpressionAttributeValues: commondynamodb.ExpressionValues{
			":account":           &types.AttributeValueMemberS{Value: accountID},
			":cumulativePayment": &types.AttributeValueMemberN{Value: cumulativePayment.String()},
		},
		ScanIndexForward: aws.Bool(true),
		Limit:            aws.Int32(1),
	}
	largerResult, err := s.dynamoClient.QueryWithInput(ctx, queryInput)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to query the next payment for account: %w", err)
	}
	var nextPayment uint64
	var nextDataLength uint32
	if len(largerResult) > 0 {
		nextPayment, err = strconv.ParseUint(largerResult[0]["CumulativePayments"].(*types.AttributeValueMemberN).Value, 10, 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("failed to parse next payment: %w", err)
		}
		dataLength, err := strconv.ParseUint(largerResult[0]["DataLength"].(*types.AttributeValueMemberN).Value, 10, 32)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("failed to parse blob size: %w", err)
		}
		nextDataLength = uint32(dataLength)
	}

	return prevPayment, nextPayment, nextDataLength, nil
}
