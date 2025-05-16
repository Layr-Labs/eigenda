package meterer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"sync"

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

const MinNumBins int32 = 3

type OffchainStore struct {
	dynamoClient         commondynamodb.Client
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

// UpdateReservationBin incrementally updates the bin usage for a specific account, period, and quorum.
// Returns the updated bin usage after the increment.
func (s *OffchainStore) UpdateReservationBin(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64, size uint64, quorumNumber uint8) (uint64, error) {
	// Create a composite key combining AccountID and QuorumNumber as the hash key
	accountAndQuorum := fmt.Sprintf("%s_%d", accountID.Hex(), quorumNumber)
	
	key := map[string]types.AttributeValue{
		"AccountAndQuorum":  &types.AttributeValueMemberS{Value: accountAndQuorum},
		"ReservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
	}
	
	// Increment the BinUsage directly (the QuorumNumber is already part of the AccountAndQuorum composite key)
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

// ReservationBinUpdate represents a single update to a reservation bin
type ReservationBinUpdate struct {
	// AccountID is the Ethereum address of the account
	AccountID gethcommon.Address
	// ReservationPeriod is the time period for the reservation
	ReservationPeriod uint64
	// Size is the amount to increment the bin usage by
	Size uint64
	// QuorumNumber is the quorum number for this update
	QuorumNumber uint8
}

// BatchUpdateReservationBins performs multiple bin updates in a single operation for efficiency.
// Returns a map of update indices to their new bin usage values, or errors if any occurred.
func (s *OffchainStore) BatchUpdateReservationBins(ctx context.Context, updates []ReservationBinUpdate) (map[int]uint64, map[int]error) {
	if len(updates) == 0 {
		return map[int]uint64{}, map[int]error{}
	}

	// For small numbers of updates, it's more efficient to use individual operations
	if len(updates) <= 2 {
		results := make(map[int]uint64, len(updates))
		errors := make(map[int]error, len(updates))

		for i, update := range updates {
			newUsage, err := s.UpdateReservationBin(ctx, update.AccountID, update.ReservationPeriod, update.Size, update.QuorumNumber)
			if err != nil {
				errors[i] = err
			} else {
				results[i] = newUsage
			}
		}

		return results, errors
	}

	// With larger update batches, optimize further with concurrent operations
	results := make(map[int]uint64, len(updates))
	errors := make(map[int]error, len(updates))

	// Use a wait group to manage concurrency
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Execute updates concurrently with a reasonable limit on concurrency
	// to avoid overwhelming the DynamoDB connection pool
	concurrencyLimit := 10
	if len(updates) < concurrencyLimit {
		concurrencyLimit = len(updates)
	}
	semaphore := make(chan struct{}, concurrencyLimit)

	for i, update := range updates {
		wg.Add(1)

		// Capture loop variables
		idx := i
		upd := update

		go func() {
			defer wg.Done()

			// Acquire semaphore slot
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// Use the updated UpdateReservationBin method which now works with the composite key
			newUsage, err := s.UpdateReservationBin(ctx, upd.AccountID, upd.ReservationPeriod, upd.Size, upd.QuorumNumber)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				errors[idx] = err
			} else {
				results[idx] = newUsage
			}
		}()
	}

	wg.Wait()
	return results, errors
}

func (s *OffchainStore) UpdateGlobalBin(ctx context.Context, reservationPeriod uint64, size uint64) (uint64, error) {
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

func (s *OffchainStore) AddOnDemandPayment(ctx context.Context, paymentMetadata core.PaymentMetadata, paymentCharged *big.Int) (*big.Int, error) {
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
func (s *OffchainStore) RollbackOnDemandPayment(ctx context.Context, accountID gethcommon.Address, newPayment, oldPayment *big.Int) error {
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

// GetPeriodRecords retrieves up to MinNumBins period records for a specific account and quorum number.
// The records are sorted by reservation period in ascending order, starting from the given period.
func (s *OffchainStore) GetPeriodRecords(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64, quorumNumber uint8) ([MinNumBins]*pb.QuorumPeriodRecord, error) {
	// Create the composite hash key combining AccountID and QuorumNumber
	accountAndQuorum := fmt.Sprintf("%s_%d", accountID.Hex(), quorumNumber)
	
	// Fetch the 3 bins starting from the current bin for the specific account and quorum
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(s.reservationTableName),
		KeyConditionExpression: aws.String("AccountAndQuorum = :accountAndQuorum AND ReservationPeriod >= :reservationPeriod"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":accountAndQuorum":  &types.AttributeValueMemberS{Value: accountAndQuorum},
			":reservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
		},
		ScanIndexForward: aws.Bool(true),
		Limit:            aws.Int32(MinNumBins), // No need for extra filtering now
	}

	bins, err := s.dynamoClient.QueryWithInput(ctx, queryInput)
	if err != nil {
		return [MinNumBins]*pb.QuorumPeriodRecord{}, fmt.Errorf("failed to query payments for account: %w", err)
	}

	records := [MinNumBins]*pb.QuorumPeriodRecord{}
	for i := 0; i < len(bins) && i < int(MinNumBins); i++ {
		periodRecord, err := parseQuorumPeriodRecord(bins[i])
		if err != nil {
			return [MinNumBins]*pb.QuorumPeriodRecord{}, fmt.Errorf("failed to parse bin %d record: %w", i, err)
		}
		records[i] = periodRecord
	}

	return records, nil
}

// GetPeriodRecordsMultiQuorum retrieves period records for multiple quorums efficiently.
// This function is optimized for retrieving period records for all quorums in a single database operation.
// Returns an array of PeriodRecords up to MinNumBins in length, with records for each requested quorum.
func (s *OffchainStore) GetPeriodRecordsMultiQuorum(
	ctx context.Context,
	accountID gethcommon.Address,
	reservationPeriod uint64,
	quorumNumbers []uint8,
) ([]*pb.QuorumPeriodRecord, error) {
	if len(quorumNumbers) == 0 {
		return []*pb.QuorumPeriodRecord{}, nil
	}

	// For multiple quorums, run parallel queries for each AccountAndQuorum composite key
	results := make([]*pb.QuorumPeriodRecord, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex
	errChan := make(chan error, len(quorumNumbers))

	for _, quorumNumber := range quorumNumbers {
		wg.Add(1)
		go func(qNum uint8) {
			defer wg.Done()
			
			// Create composite key for this account and quorum
			accountAndQuorum := fmt.Sprintf("%s_%d", accountID.Hex(), qNum)
			
			// Query for this specific account and quorum
			queryInput := &dynamodb.QueryInput{
				TableName:              aws.String(s.reservationTableName),
				KeyConditionExpression: aws.String("AccountAndQuorum = :accountAndQuorum AND ReservationPeriod >= :reservationPeriod"),
				ExpressionAttributeValues: map[string]types.AttributeValue{
					":accountAndQuorum":  &types.AttributeValueMemberS{Value: accountAndQuorum},
					":reservationPeriod": &types.AttributeValueMemberN{Value: strconv.FormatUint(reservationPeriod, 10)},
				},
				ScanIndexForward: aws.Bool(true),
				Limit:            aws.Int32(MinNumBins),
			}

			bins, err := s.dynamoClient.QueryWithInput(ctx, queryInput)
			if err != nil {
				errChan <- fmt.Errorf("failed to query payments for quorum %d: %w", qNum, err)
				return
			}

			// Parse records and add to results
			quorumRecords := make([]*pb.QuorumPeriodRecord, 0, len(bins))
			for _, bin := range bins {
				periodRecord, err := parseQuorumPeriodRecord(bin)
				if err != nil {
					s.logger.Debug("Failed to parse period record", "err", err, "quorum", qNum)
					continue
				}
				quorumRecords = append(quorumRecords, periodRecord)
			}

			// Add records to the combined result under lock
			mu.Lock()
			results = append(results, quorumRecords...)
			mu.Unlock()
		}(quorumNumber)
	}

	// Wait for all queries to complete
	wg.Wait()
	close(errChan)

	// Check if any errors occurred
	if len(errChan) > 0 {
		err := <-errChan
		return results, fmt.Errorf("error in one or more quorum queries: %w", err)
	}

	return results, nil
}

func (s *OffchainStore) GetLargestCumulativePayment(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
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

func parseQuorumPeriodRecord(bin map[string]types.AttributeValue) (*pb.QuorumPeriodRecord, error) {
	// Parse ReservationPeriod from the response
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

	// Parse BinUsage from the response
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

	// Extract QuorumNumber from the AccountAndQuorum composite key
	var quorumNumber uint32 = 0 // Default to 0 if not extractable
	accountAndQuorum, ok := bin["AccountAndQuorum"]
	if ok {
		accountAndQuorumStr, ok := accountAndQuorum.(*types.AttributeValueMemberS)
		if ok {
			// Parse quorum number from the composite key format "accountID_quorumNumber"
			parts := strings.Split(accountAndQuorumStr.Value, "_")
			if len(parts) == 2 {
				if qNum, err := strconv.ParseUint(parts[1], 10, 32); err == nil {
					quorumNumber = uint32(qNum)
				}
			}
		}
	}

	return &pb.QuorumPeriodRecord{
		Index:        uint32(reservationPeriodValue),
		Usage:        uint64(binUsageValue),
		QuorumNumber: quorumNumber,
	}, nil
}
