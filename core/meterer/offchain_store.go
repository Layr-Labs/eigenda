package meterer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const MinNumPeriods int32 = 3

type OffchainStore struct {
	dynamoClient         commondynamodb.Client
	reservationTableName string
	onDemandTableName    string
	globalBinTableName   string
	logger               logging.Logger
	// TODO: add maximum storage for both tables
	MaxOnDemandStorage uint64
}

func NewOffchainStore(
	cfg commonaws.ClientConfig,
	reservationTableName string,
	onDemandTableName string,
	globalBinTableName string,
	maxOnDemandStorage uint64,
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
	// delete expired reservation periods (<i-1) and old on-demand payments (retain max N payments)
	return OffchainStore{
		dynamoClient:         dynamoClient,
		reservationTableName: reservationTableName,
		onDemandTableName:    onDemandTableName,
		globalBinTableName:   globalBinTableName,
		logger:               logger,
		MaxOnDemandStorage:   maxOnDemandStorage,
	}, nil
}

func (s *OffchainStore) UpdateReservationPeriod(ctx context.Context, accountID string, reservationPeriod uint64, size uint64) (uint64, error) {
	key := map[string]types.AttributeValue{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
	}

	res, err := s.dynamoClient.IncrementBy(ctx, s.reservationTableName, key, "PeriodUsage", size)
	if err != nil {
		return 0, fmt.Errorf("failed to increment bin usage: %w", err)
	}

	periodUsage, ok := res["PeriodUsage"]
	if !ok {
		return 0, errors.New("PeriodUsage is not present in the response")
	}

	periodUsageAttr, ok := periodUsage.(*types.AttributeValueMemberN)
	if !ok {
		return 0, fmt.Errorf("unexpected type for PeriodUsage: %T", periodUsage)
	}

	periodUsageValue, err := strconv.ParseUint(periodUsageAttr.Value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse PeriodUsage: %w", err)
	}

	return periodUsageValue, nil
}

func (s *OffchainStore) UpdateGlobalPeriod(ctx context.Context, reservationPeriod uint32, size uint64) (uint64, error) {
	key := map[string]types.AttributeValue{
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(reservationPeriod), 10)},
	}

	res, err := s.dynamoClient.IncrementBy(ctx, s.globalBinTableName, key, "PeriodUsage", size)
	if err != nil {
		return 0, err
	}

	periodUsage, ok := res["PeriodUsage"]
	if !ok {
		return 0, nil
	}

	periodUsageAttr, ok := periodUsage.(*types.AttributeValueMemberN)
	if !ok {
		return 0, nil
	}

	periodUsageValue, err := strconv.ParseUint(periodUsageAttr.Value, 10, 32)
	if err != nil {
		return 0, err
	}

	return periodUsageValue, nil
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

	if err := s.PruneOnDemandPayments(ctx, paymentMetadata.AccountID); err != nil {
		// Don't fail the request if pruning fails, just log a warning
		s.logger.Warn("failed to prune on-demand payments", "accountID", paymentMetadata.AccountID, "error", err)
	}

	return nil
}

func (s *OffchainStore) PruneOnDemandPayments(ctx context.Context, accountID string) error {
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(s.onDemandTableName),
		KeyConditionExpression: aws.String("AccountID = :account"),
		ExpressionAttributeValues: commondynamodb.ExpressionValues{
			":account": &types.AttributeValueMemberS{Value: accountID},
		},
		ScanIndexForward: aws.Bool(true), // ascending order
	}

	payments, err := s.dynamoClient.QueryWithInput(ctx, queryInput)
	if err != nil {
		return fmt.Errorf("failed to query existing payments: %w", err)
	}

	if len(payments) >= int(s.MaxOnDemandStorage) {
		numToDelete := len(payments) - int(s.MaxOnDemandStorage) + 1
		// Create keys for all payments to delete (taking the smallest cumulative payments)
		keysToDelete := make([]commondynamodb.Key, numToDelete)
		for i := 0; i < numToDelete; i++ {
			keysToDelete[i] = commondynamodb.Key{
				"AccountID":          payments[i]["AccountID"],
				"CumulativePayments": payments[i]["CumulativePayments"],
			}
		}

		// Delete the items in batches
		failedKeys, err := s.dynamoClient.DeleteItems(ctx, s.onDemandTableName, keysToDelete)
		if err != nil {
			return fmt.Errorf("failed to delete oldest payments: %w", err)
		}
		if len(failedKeys) > 0 {
			return fmt.Errorf("failed to delete %d payments", len(failedKeys))
		}
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
func (s *OffchainStore) GetRelevantOnDemandRecords(ctx context.Context, accountID string, cumulativePayment *big.Int) (*big.Int, *big.Int, uint32, error) {
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
		return nil, nil, 0, fmt.Errorf("failed to query smaller payments for account: %w", err)
	}
	prevPayment := big.NewInt(0)
	if len(smallerResult) > 0 {
		cumulativePaymentsAttr, ok := smallerResult[0]["CumulativePayments"]
		if !ok {
			return nil, nil, 0, fmt.Errorf("CumulativePayments field not found in result")
		}
		cumulativePaymentsNum, ok := cumulativePaymentsAttr.(*types.AttributeValueMemberN)
		if !ok {
			return nil, nil, 0, fmt.Errorf("CumulativePayments has invalid type")
		}
		setPrevPayment, success := prevPayment.SetString(cumulativePaymentsNum.Value, 10)
		if !success {
			return nil, nil, 0, fmt.Errorf("failed to parse previous payment: %w", err)
		}
		prevPayment = setPrevPayment
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
		return nil, nil, 0, fmt.Errorf("failed to query the next payment for account: %w", err)
	}
	nextPayment := big.NewInt(0)
	nextDataLength := uint32(0)
	if len(largerResult) > 0 {
		cumulativePaymentsAttr, ok := largerResult[0]["CumulativePayments"]
		if !ok {
			return nil, nil, 0, fmt.Errorf("CumulativePayments field not found in result")
		}
		cumulativePaymentsNum, ok := cumulativePaymentsAttr.(*types.AttributeValueMemberN)
		if !ok {
			return nil, nil, 0, fmt.Errorf("CumulativePayments has invalid type")
		}
		setNextPayment, success := nextPayment.SetString(cumulativePaymentsNum.Value, 10)
		if !success {
			return nil, nil, 0, fmt.Errorf("failed to parse previous payment: %w", err)
		}
		nextPayment = setNextPayment

		dataLengthAttr, ok := largerResult[0]["DataLength"]
		if !ok {
			return nil, nil, 0, fmt.Errorf("DataLength field not found in result")
		}
		dataLengthNum, ok := dataLengthAttr.(*types.AttributeValueMemberN)
		if !ok {
			return nil, nil, 0, fmt.Errorf("DataLength has invalid type")
		}
		dataLength, err := strconv.ParseUint(dataLengthNum.Value, 10, 32)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed to parse data length: %w", err)
		}
		nextDataLength = uint32(dataLength)
	}

	return prevPayment, nextPayment, nextDataLength, nil
}

func (s *OffchainStore) GetReservationPeriodRecords(ctx context.Context, accountID string, reservationPeriod uint32) ([MinNumPeriods]*pb.ReservationPeriodRecord, error) {
	// Fetch the 3 periods start from the current bin
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(s.reservationTableName),
		KeyConditionExpression: aws.String("AccountID = :account AND ReservationPeriod > :reservationPeriod"),
		ExpressionAttributeValues: commondynamodb.ExpressionValues{
			":account":           &types.AttributeValueMemberS{Value: accountID},
			":reservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(reservationPeriod), 10)},
		},
		ScanIndexForward: aws.Bool(true),
		Limit:            aws.Int32(MinNumPeriods),
	}
	periods, err := s.dynamoClient.QueryWithInput(ctx, queryInput)
	if err != nil {
		return [MinNumPeriods]*pb.ReservationPeriodRecord{}, fmt.Errorf("failed to query payments for account: %w", err)
	}

	records := [MinNumPeriods]*pb.ReservationPeriodRecord{}
	for i := 0; i < len(periods) && i < int(MinNumPeriods); i++ {
		reservationPeriodRecord, err := parseReservationPeriodRecord(periods[i])
		if err != nil {
			return [MinNumPeriods]*pb.ReservationPeriodRecord{}, fmt.Errorf("failed to parse bin %d record: %w", i, err)
		}
		records[i] = reservationPeriodRecord
	}

	return records, nil
}

func (s *OffchainStore) GetLargestCumulativePayment(ctx context.Context, accountID string) (*big.Int, error) {
	// Fetch the largest cumulative payment
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(s.onDemandTableName),
		KeyConditionExpression: aws.String("AccountID = :account"),
		ExpressionAttributeValues: commondynamodb.ExpressionValues{
			":account": &types.AttributeValueMemberS{Value: accountID},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(1),
	}
	payments, err := s.dynamoClient.QueryWithInput(ctx, queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to query payments for account: %w", err)
	}

	if len(payments) == 0 {
		return big.NewInt(0), nil
	}

	var payment *big.Int
	_, success := payment.SetString(payments[0]["CumulativePayments"].(*types.AttributeValueMemberN).Value, 10)
	if !success {
		return nil, fmt.Errorf("failed to parse payment: %w", err)
	}

	return payment, nil
}

// DeleteOldPeriods removes all reservation bin entries with indices strictly less than the provided reservationPeriod
func (s *OffchainStore) DeleteOldPeriods(ctx context.Context, reservationPeriod uint32) error {
	// get all keys that need to be deleted
	queryInput := &dynamodb.QueryInput{
		TableName:        aws.String(s.reservationTableName),
		FilterExpression: aws.String("ReservationPeriod < :reservationPeriod"),
		ExpressionAttributeValues: commondynamodb.ExpressionValues{
			":reservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(reservationPeriod), 10)},
		},
	}

	items, err := s.dynamoClient.QueryWithInput(ctx, queryInput)
	if err != nil {
		return fmt.Errorf("failed to query old periods: %w", err)
	}

	keys := make([]commondynamodb.Key, len(items))
	for i, item := range items {
		keys[i] = commondynamodb.Key{
			"AccountID":         item["AccountID"],
			"ReservationPeriod": item["ReservationPeriod"],
		}
	}

	// Delete the items in batches
	if len(keys) > 0 {
		failedKeys, err := s.dynamoClient.DeleteItems(ctx, s.reservationTableName, keys)
		if err != nil {
			return fmt.Errorf("failed to delete old periods: %w", err)
		}
		if len(failedKeys) > 0 {
			return fmt.Errorf("failed to delete %d periods", len(failedKeys))
		}
	}

	return nil
}

func parseReservationPeriodRecord(bin map[string]types.AttributeValue) (*pb.ReservationPeriodRecord, error) {
	reservationPeriod, ok := bin["ReservationPeriod"]
	if !ok {
		return nil, errors.New("ReservationPeriod is not present in the response")
	}

	reservationPeriodAttr, ok := reservationPeriod.(*types.AttributeValueMemberN)
	if !ok {
		return nil, fmt.Errorf("unexpected type for ReservationPeriod: %T", reservationPeriod)
	}

	reservationPeriodValue, err := strconv.ParseUint(reservationPeriodAttr.Value, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ReservationPeriod: %w", err)
	}

	periodUsage, ok := bin["PeriodUsage"]
	if !ok {
		return nil, errors.New("PeriodUsage is not present in the response")
	}

	periodUsageAttr, ok := periodUsage.(*types.AttributeValueMemberN)
	if !ok {
		return nil, fmt.Errorf("unexpected type for PeriodUsage: %T", periodUsage)
	}

	periodUsageValue, err := strconv.ParseUint(periodUsageAttr.Value, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PeriodUsage: %w", err)
	}

	return &pb.ReservationPeriodRecord{
		Index: uint32(reservationPeriodValue),
		Usage: uint64(periodUsageValue),
	}, nil
}
