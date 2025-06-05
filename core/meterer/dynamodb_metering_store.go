package meterer

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

// DynamoDBMeteringStore implements the MeteringStore interface using DynamoDB
type DynamoDBMeteringStore struct {
	dynamoClient         commondynamodb.Client
	reservationTableName string
	onDemandTableName    string
	globalBinTableName   string
	logger               logging.Logger
	// TODO: add maximum storage for both tables
}

// NewDynamoDBMeteringStore creates a new DynamoDB-backed metering store
func NewDynamoDBMeteringStore(
	cfg commonaws.ClientConfig,
	reservationTableName string,
	onDemandTableName string,
	globalBinTableName string,
	logger logging.Logger,
) (*DynamoDBMeteringStore, error) {
	dynamoClient, err := commondynamodb.NewClient(cfg, logger)
	if err != nil {
		return nil, err
	}

	err = dynamoClient.TableExists(context.Background(), reservationTableName)
	if err != nil {
		return nil, err
	}
	err = dynamoClient.TableExists(context.Background(), onDemandTableName)
	if err != nil {
		return nil, err
	}
	err = dynamoClient.TableExists(context.Background(), globalBinTableName)
	if err != nil {
		return nil, err
	}
	//TODO: add a separate thread to periodically clean up the tables
	// delete expired reservation bins (<i-1) and old on-demand payments (retain max N payments)
	return &DynamoDBMeteringStore{
		dynamoClient:         dynamoClient,
		reservationTableName: reservationTableName,
		onDemandTableName:    onDemandTableName,
		globalBinTableName:   globalBinTableName,
		logger:               logger,
	}, nil
}

// IncrementBinUsages updates the bin usage for each quorum in quorumNumbers for a specific account and reservation period.
// The key AccountID is formatted as {AccountID}:{quorumNumber}.
func (s *DynamoDBMeteringStore) IncrementBinUsages(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID, reservationPeriods map[core.QuorumID]uint64, sizes map[core.QuorumID]uint64) (map[core.QuorumID]uint64, error) {
	binUsages := make(map[core.QuorumID]uint64)

	// Build ops for atomic batch increment
	ops := make([]commondynamodb.TransactAddOp, len(quorumNumbers))
	for i, quorumNumber := range quorumNumbers {
		accountIDAndQuorum := accountID.Hex() + ":" + strconv.FormatUint(uint64(quorumNumber), 10)
		key := map[string]types.AttributeValue{
			"AccountID":         &types.AttributeValueMemberS{Value: accountIDAndQuorum},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriods[quorumNumber], 10)},
		}
		ops[i] = commondynamodb.TransactAddOp{
			Key:   key,
			Attr:  "BinUsage",
			Value: float64(sizes[quorumNumber]), // positive for increment
		}
	}

	err := s.dynamoClient.TransactAddBy(ctx, s.reservationTableName, ops)
	if err != nil {
		return nil, err
	}

	// Fetch new values for each key
	for _, quorumNumber := range quorumNumbers {
		accountIDAndQuorum := accountID.Hex() + ":" + strconv.FormatUint(uint64(quorumNumber), 10)
		key := map[string]types.AttributeValue{
			"AccountID":         &types.AttributeValueMemberS{Value: accountIDAndQuorum},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriods[quorumNumber], 10)},
		}
		item, getErr := s.dynamoClient.GetItem(ctx, s.reservationTableName, key)
		if getErr != nil {
			return nil, getErr
		}
		binUsage, ok := item["BinUsage"]
		if !ok {
			return nil, fmt.Errorf("BinUsage is not present in the response")
		}
		binUsageAttr, ok := binUsage.(*types.AttributeValueMemberN)
		if !ok {
			return nil, fmt.Errorf("unexpected type for BinUsage: %T", binUsage)
		}
		binUsageValue, parseErr := strconv.ParseUint(binUsageAttr.Value, 10, 64)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse BinUsage: %w", parseErr)
		}
		binUsages[quorumNumber] = binUsageValue
	}

	return binUsages, nil
}

func (s *DynamoDBMeteringStore) UpdateGlobalBin(ctx context.Context, reservationPeriod uint64, size uint64) (uint64, error) {
	key := map[string]types.AttributeValue{
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
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

func (s *DynamoDBMeteringStore) AddOnDemandPayment(ctx context.Context, paymentMetadata core.PaymentMetadata, paymentCharged *big.Int) (*big.Int, error) {
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

// RollbackOnDemandPayment rolls back a payment to the previous value
// If oldPayment is 0, it writes a zero value instead of deleting the record
// This method uses a conditional expression to ensure we only roll back if the current value matches newPayment
func (s *DynamoDBMeteringStore) RollbackOnDemandPayment(ctx context.Context, accountID gethcommon.Address, newPayment, oldPayment *big.Int) error {
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

// GetPeriodRecordsMultiQuorum retrieves period records for multiple quorums efficiently.
// This function is optimized for retrieving period records for all quorums in a single database operation.
// Returns an array of PeriodRecords up to numBins in length, with records for each requested quorum.
func (s *DynamoDBMeteringStore) GetPeriodRecordsMultiQuorum(
	ctx context.Context,
	accountID gethcommon.Address,
	reservationPeriod uint64,
	quorumNumbers []core.QuorumID,
	numBins uint32,
) (map[core.QuorumID]*pb.PeriodRecords, error) {
	if len(quorumNumbers) == 0 {
		return nil, nil
	}

	// Prepare all keys for batch get
	var keys []map[string]types.AttributeValue
	for _, quorum := range quorumNumbers {
		accountIDAndQuorum := accountID.Hex() + ":" + strconv.FormatUint(uint64(quorum), 10)
		for i := 0; i < int(numBins); i++ {
			key := map[string]types.AttributeValue{
				"AccountID":         &types.AttributeValueMemberS{Value: accountIDAndQuorum},
				"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod+uint64(i), 10)},
			}
			keys = append(keys, key)
		}
	}

	items, err := s.dynamoClient.GetItems(ctx, s.reservationTableName, keys, true)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get period records for account: %w", err)
	}

	records := make(map[core.QuorumID]*pb.PeriodRecords)
	for _, item := range items {
		quorumNumber, periodRecord, err := parsePeriodRecord(item)
		if err != nil {
			s.logger.Debug("Failed to parse period record", "err", err)
			continue
		}
		records[quorumNumber] = &pb.PeriodRecords{
			Records: []*pb.PeriodRecord{periodRecord},
		}
	}

	return records, nil
}

func (s *DynamoDBMeteringStore) GetLargestCumulativePayment(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
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

func parsePeriodRecord(bin map[string]types.AttributeValue) (core.QuorumID, *pb.PeriodRecord, error) {
	reservationPeriod, ok := bin["ReservationPeriod"]
	if !ok {
		return 0, nil, errors.New("ReservationPeriod is not present in the response")
	}

	reservationPeriodAttr, ok := reservationPeriod.(*types.AttributeValueMemberN)
	if !ok {
		return 0, nil, fmt.Errorf("unexpected type for ReservationPeriod: %T", reservationPeriod)
	}

	reservationPeriodValue, err := strconv.ParseUint(reservationPeriodAttr.Value, 10, 32)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse ReservationPeriod: %w", err)
	}

	binUsage, ok := bin["BinUsage"]
	if !ok {
		return 0, nil, errors.New("BinUsage is not present in the response")
	}

	binUsageAttr, ok := binUsage.(*types.AttributeValueMemberN)
	if !ok {
		return 0, nil, fmt.Errorf("unexpected type for BinUsage: %T", binUsage)
	}

	binUsageValue, err := strconv.ParseUint(binUsageAttr.Value, 10, 32)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse BinUsage: %w", err)
	}
	accountIDAndQuorum, ok := bin["AccountIDAndQuorum"]
	if !ok {
		return 0, nil, errors.New("AccountIDAndQuorum is not present in the response")
	}

	accountIDAndQuorumAttr, ok := accountIDAndQuorum.(*types.AttributeValueMemberS)
	if !ok {
		return 0, nil, fmt.Errorf("unexpected type for AccountIDAndQuorum: %T", accountIDAndQuorum)
	}

	parts := strings.Split(accountIDAndQuorumAttr.Value, ":")
	if len(parts) != 2 {
		return 0, nil, fmt.Errorf("invalid AccountIDAndQuorum format: %s", accountIDAndQuorumAttr.Value)
	}

	quorumNumber, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to parse QuorumNumber: %w", err)
	}
	if quorumNumber > math.MaxUint8 {
		return 0, nil, fmt.Errorf("QuorumNumber exceeds maximum value for uint8: %d", quorumNumber)
	}

	return core.QuorumID(quorumNumber), &pb.PeriodRecord{
		Index: uint32(reservationPeriodValue),
		Usage: uint64(binUsageValue),
	}, nil
}

// DecrementBinUsages atomically decrements the bin usage for each quorum in quorumNumbers for a specific account and reservation period.
// The key is AccountIDAndQuorum, formatted as {AccountID}:{quorumNumber}.
func (s *DynamoDBMeteringStore) DecrementBinUsages(ctx context.Context, accountID gethcommon.Address, quorumNumbers []core.QuorumID, reservationPeriods map[core.QuorumID]uint64, sizes map[core.QuorumID]uint64) error {
	// Build ops for atomic batch decrement
	ops := make([]commondynamodb.TransactAddOp, len(quorumNumbers))
	for i, quorumNumber := range quorumNumbers {
		accountIDAndQuorum := accountID.Hex() + ":" + strconv.FormatUint(uint64(quorumNumber), 10)
		key := map[string]types.AttributeValue{
			"AccountID":         &types.AttributeValueMemberS{Value: accountIDAndQuorum},
			"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriods[quorumNumber], 10)},
		}
		ops[i] = commondynamodb.TransactAddOp{
			Key:   key,
			Attr:  "BinUsage",
			Value: -float64(sizes[quorumNumber]), // negative for decrement
		}
	}

	return s.dynamoClient.TransactAddBy(ctx, s.reservationTableName, ops)
}
