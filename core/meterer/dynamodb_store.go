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
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// DynamoDBOffchainStore implements OffchainStore using DynamoDB
type DynamoDBOffchainStore struct {
	dynamoClient         commondynamodb.Client
	reservationTableName string
	onDemandTableName    string
	globalBinTableName   string
	logger               logging.Logger
}

// NewDynamoDBOffchainStore creates a new DynamoDBOffchainStore
func NewDynamoDBOffchainStore(
	cfg commonaws.ClientConfig,
	reservationTableName string,
	onDemandTableName string,
	globalBinTableName string,
	logger logging.Logger,
) (OffchainStore, error) {
	dynamoClient, err := commondynamodb.NewClient(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamodb client: %w", err)
	}

	err = dynamoClient.TableExists(context.Background(), reservationTableName)
	if err != nil {
		return nil, fmt.Errorf("reservation table does not exist: %w", err)
	}
	err = dynamoClient.TableExists(context.Background(), onDemandTableName)
	if err != nil {
		return nil, fmt.Errorf("on-demand table does not exist: %w", err)
	}
	err = dynamoClient.TableExists(context.Background(), globalBinTableName)
	if err != nil {
		return nil, fmt.Errorf("global bin table does not exist: %w", err)
	}

	return &DynamoDBOffchainStore{
		dynamoClient:         dynamoClient,
		reservationTableName: reservationTableName,
		onDemandTableName:    onDemandTableName,
		globalBinTableName:   globalBinTableName,
		logger:               logger,
	}, nil
}

func (s *DynamoDBOffchainStore) UpdateReservationBin(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64, size uint64) (uint64, error) {
	key := map[string]types.AttributeValue{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
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

func (s *DynamoDBOffchainStore) UpdateGlobalBin(ctx context.Context, reservationPeriod uint64, size uint64) (uint64, error) {
	key := map[string]types.AttributeValue{
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
	}

	res, err := s.dynamoClient.IncrementBy(ctx, s.globalBinTableName, key, "BinUsage", size)
	if err != nil {
		return 0, fmt.Errorf("failed to increment global bin usage: %w", err)
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
		return 0, fmt.Errorf("failed to parse BinUsage: %w", err)
	}

	return binUsageValue, nil
}

func (s *DynamoDBOffchainStore) AddOnDemandPayment(ctx context.Context, paymentMetadata core.PaymentMetadata, paymentCharged *big.Int) (*big.Int, error) {
	// Create new item with only AccountID and CumulativePayment
	item := commondynamodb.Item{
		"AccountID":         &types.AttributeValueMemberS{Value: paymentMetadata.AccountID.Hex()},
		"CumulativePayment": &types.AttributeValueMemberN{Value: paymentMetadata.CumulativePayment.String()},
	}

	// Use conditional expression to ensure:
	// 1. If no record exists, accept the payment
	// 2. If record exists, the increment must be at least the payment charged
	//    (which also ensures the new payment is larger than the existing one since paymentCharged > 0)
	paymentCheckpoint := big.NewInt(0).Sub(paymentMetadata.CumulativePayment, paymentCharged)
	if paymentCheckpoint.Sign() < 0 {
		return nil, fmt.Errorf("payment validation failed: payment charged is greater than cumulative payment")
	}
	conditionExpression := "attribute_not_exists(CumulativePayment) OR " +
		"CumulativePayment <= :payment"

	expressionValues := map[string]types.AttributeValue{
		":payment": &types.AttributeValueMemberN{Value: paymentCheckpoint.String()},
	}

	oldItem, err := s.dynamoClient.PutItemWithConditionAndReturn(ctx, s.onDemandTableName, item, conditionExpression, nil, expressionValues)
	if err != nil {
		if errors.Is(err, commondynamodb.ErrConditionFailed) {
			return nil, fmt.Errorf("insufficient cumulative payment increment: %w", err)
		}
		return nil, fmt.Errorf("failed to add on-demand payment: %w", err)
	}

	// If there was no previous item, return zero
	if len(oldItem) == 0 {
		return big.NewInt(0), nil
	}

	// Extract the old CumulativePayment value
	oldPaymentAttr, ok := oldItem["CumulativePayment"]
	if !ok {
		return big.NewInt(0), nil
	}

	// Type assertion with check
	oldPaymentNum, ok := oldPaymentAttr.(*types.AttributeValueMemberN)
	if !ok {
		return big.NewInt(0), fmt.Errorf("CumulativePayment has invalid type: %T", oldPaymentAttr)
	}

	oldPayment := new(big.Int)
	if _, success := oldPayment.SetString(oldPaymentNum.Value, 10); !success {
		return big.NewInt(0), fmt.Errorf("failed to parse old payment value: %s", oldPaymentNum.Value)
	}

	return oldPayment, nil
}

func (s *DynamoDBOffchainStore) RollbackOnDemandPayment(ctx context.Context, accountID gethcommon.Address, newPayment, oldPayment *big.Int) error {
	// Initialize oldPayment to zero if it's nil
	if oldPayment == nil {
		oldPayment = big.NewInt(0)
	}

	// Create the item with the old payment value (which might be zero)
	item := commondynamodb.Item{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"CumulativePayment": &types.AttributeValueMemberN{Value: oldPayment.String()},
	}

	// Construct a condition expression as a string
	conditionExpression := "attribute_not_exists(CumulativePayment) OR CumulativePayment = :expectedPayment"

	// Create the expression attribute values map
	expressionValues := map[string]types.AttributeValue{
		":expectedPayment": &types.AttributeValueMemberN{Value: newPayment.String()},
	}

	err := s.dynamoClient.PutItemWithCondition(
		ctx,
		s.onDemandTableName,
		item,
		conditionExpression,
		nil, // No expression attribute names needed
		expressionValues,
	)

	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		if s.logger != nil {
			s.logger.Debug("Skipping rollback as current payment doesn't match the expected value",
				"accountID", accountID.Hex(),
				"expectedPayment", newPayment.String())
		}
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to rollback payment: %w", err)
	}

	if s.logger != nil {
		s.logger.Debug("Successfully rolled back payment to previous value",
			"accountID", accountID.Hex(),
			"rolledBackFrom", newPayment.String(),
			"rolledBackTo", oldPayment.String())
	}

	return nil
}

func (s *DynamoDBOffchainStore) GetPeriodRecords(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64) ([MinNumBins]*pb.PeriodRecord, error) {
	// Fetch the 3 bins start from the current bin
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(s.reservationTableName),
		KeyConditionExpression: aws.String("AccountID = :account AND ReservationPeriod >= :reservationPeriod"),
		ExpressionAttributeValues: commondynamodb.ExpressionValues{
			":account":           &types.AttributeValueMemberS{Value: accountID.Hex()},
			":reservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
		},
		ScanIndexForward: aws.Bool(true),
		Limit:            aws.Int32(MinNumBins),
	}
	bins, err := s.dynamoClient.QueryWithInput(ctx, queryInput)
	if err != nil {
		return [MinNumBins]*pb.PeriodRecord{}, fmt.Errorf("failed to query payments for account: %w", err)
	}

	records := [MinNumBins]*pb.PeriodRecord{}
	for i := 0; i < len(bins) && i < int(MinNumBins); i++ {
		periodRecord, err := parsePeriodRecord(bins[i])
		if err != nil {
			return [MinNumBins]*pb.PeriodRecord{}, fmt.Errorf("failed to parse bin %d record: %w", i, err)
		}
		records[i] = periodRecord
	}

	return records, nil
}

func (s *DynamoDBOffchainStore) GetLargestCumulativePayment(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	// Get the single record for this account
	key := commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: accountID.Hex()},
	}

	result, err := s.dynamoClient.GetItem(ctx, s.onDemandTableName, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment for account: %w", err)
	}

	// If no item found, return zero
	if len(result) == 0 {
		return big.NewInt(0), nil
	}

	// Extract CumulativePayment
	largestPaymentAttr, ok := result["CumulativePayment"]
	if !ok {
		return big.NewInt(0), nil
	}

	// Type assertion with check
	largestPaymentNum, ok := largestPaymentAttr.(*types.AttributeValueMemberN)
	if !ok {
		return nil, fmt.Errorf("CumulativePayment has invalid type: %T", largestPaymentAttr)
	}

	payment := new(big.Int)
	if _, success := payment.SetString(largestPaymentNum.Value, 10); !success {
		return nil, fmt.Errorf("failed to parse payment value: %s", largestPaymentNum.Value)
	}

	return payment, nil
}

func (s *DynamoDBOffchainStore) GetGlobalBinUsage(ctx context.Context, reservationPeriod uint64) (uint64, error) {
	key := map[string]types.AttributeValue{
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
	}

	result, err := s.dynamoClient.GetItem(ctx, s.globalBinTableName, key)
	if err != nil {
		return 0, fmt.Errorf("failed to get global bin usage: %w", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	binUsage, ok := result["BinUsage"]
	if !ok {
		return 0, nil
	}

	binUsageAttr, ok := binUsage.(*types.AttributeValueMemberN)
	if !ok {
		return 0, nil
	}

	binUsageValue, err := strconv.ParseUint(binUsageAttr.Value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse BinUsage: %w", err)
	}

	return binUsageValue, nil
}

func (s *DynamoDBOffchainStore) GetReservationBin(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64) (uint64, error) {
	key := map[string]types.AttributeValue{
		"AccountID":         &types.AttributeValueMemberS{Value: accountID.Hex()},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
	}

	result, err := s.dynamoClient.GetItem(ctx, s.reservationTableName, key)
	if err != nil {
		return 0, fmt.Errorf("failed to get reservation bin: %w", err)
	}

	if len(result) == 0 {
		return 0, nil
	}

	binUsage, ok := result["BinUsage"]
	if !ok {
		return 0, nil
	}

	binUsageAttr, ok := binUsage.(*types.AttributeValueMemberN)
	if !ok {
		return 0, nil
	}

	binUsageValue, err := strconv.ParseUint(binUsageAttr.Value, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("failed to parse BinUsage: %w", err)
	}

	return binUsageValue, nil
}

func (s *DynamoDBOffchainStore) GetOnDemandPayment(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	key := commondynamodb.Key{
		"AccountID": &types.AttributeValueMemberS{Value: accountID.Hex()},
	}

	result, err := s.dynamoClient.GetItem(ctx, s.onDemandTableName, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get on-demand payment: %w", err)
	}

	if len(result) == 0 {
		return big.NewInt(0), nil
	}

	paymentAttr, ok := result["CumulativePayment"]
	if !ok {
		return big.NewInt(0), nil
	}

	paymentNum, ok := paymentAttr.(*types.AttributeValueMemberN)
	if !ok {
		return nil, fmt.Errorf("CumulativePayment has invalid type: %T", paymentAttr)
	}

	payment := new(big.Int)
	if _, success := payment.SetString(paymentNum.Value, 10); !success {
		return nil, fmt.Errorf("failed to parse payment value: %s", paymentNum.Value)
	}

	return payment, nil
}

func (s *DynamoDBOffchainStore) GetGlobalBin(ctx context.Context, reservationPeriod uint64) (uint64, error) {
	return s.GetGlobalBinUsage(ctx, reservationPeriod)
}

// Destroy shuts down and permanently deletes all data in the store
func (s *DynamoDBOffchainStore) Destroy() error {
	// DynamoDB doesn't have a destroy method, so we just return nil
	return nil
}

func parsePeriodRecord(bin map[string]types.AttributeValue) (*pb.PeriodRecord, error) {
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

	binUsage, ok := bin["BinUsage"]
	if !ok {
		return nil, errors.New("BinUsage is not present in the response")
	}

	binUsageAttr, ok := binUsage.(*types.AttributeValueMemberN)
	if !ok {
		return nil, fmt.Errorf("unexpected type for BinUsage: %T", binUsage)
	}

	binUsageValue, err := strconv.ParseUint(binUsageAttr.Value, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse BinUsage: %w", err)
	}

	return &pb.PeriodRecord{
		Index: uint32(reservationPeriodValue),
		Usage: uint64(binUsageValue),
	}, nil
}
