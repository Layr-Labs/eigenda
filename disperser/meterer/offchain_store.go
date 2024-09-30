package meterer

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
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
) (*OffchainStore, error) {

	dynamoClient, err := commondynamodb.NewClient(cfg, logger)
	if err != nil {
		return nil, err
	}
	if reservationTableName == "" || onDemandTableName == "" || globalBinTableName == "" {
		return nil, fmt.Errorf("table names cannot be empty")
	}

	return &OffchainStore{
		dynamoClient:         dynamoClient,
		reservationTableName: reservationTableName,
		onDemandTableName:    onDemandTableName,
		globalBinTableName:   globalBinTableName,
		logger:               logger,
	}, nil
}

type ReservationBin struct {
	AccountID string
	BinIndex  uint64
	BinUsage  uint32
	UpdatedAt time.Time
}

type PaymentTuple struct {
	CumulativePayment uint64
	BlobSize          uint32
}

type GlobalBin struct {
	BinIndex  uint64
	BinUsage  uint64
	UpdatedAt time.Time
}

func (s *OffchainStore) UpdateReservationBin(ctx context.Context, accountID string, binIndex uint64, size uint32) (uint32, error) {
	key := map[string]types.AttributeValue{
		"AccountID": &types.AttributeValueMemberS{Value: accountID},
		"BinIndex":  &types.AttributeValueMemberN{Value: strconv.FormatUint(binIndex, 10)},
	}

	update := map[string]types.AttributeValue{
		"BinUsage": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(size), 10)},
	}

	res, err := s.dynamoClient.UpdateItemIncrement(ctx, s.reservationTableName, key, update)
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

	return uint32(binUsageValue), nil
}

func (s *OffchainStore) UpdateGlobalBin(ctx context.Context, binIndex uint64, size uint32) (uint64, error) {
	key := map[string]types.AttributeValue{
		"BinIndex": &types.AttributeValueMemberN{Value: strconv.FormatUint(binIndex, 10)},
	}

	update := map[string]types.AttributeValue{
		"BinUsage": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(size), 10)},
	}
	res, err := s.dynamoClient.UpdateItemIncrement(ctx, s.globalBinTableName, key, update)
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

func (s *OffchainStore) FindReservationBin(ctx context.Context, accountID string, binIndex uint64) (*ReservationBin, error) {
	key := map[string]types.AttributeValue{
		"AccountID": &types.AttributeValueMemberS{Value: accountID},
		"BinIndex":  &types.AttributeValueMemberN{Value: strconv.FormatUint(binIndex, 10)},
	}

	result, err := s.dynamoClient.GetItem(ctx, s.reservationTableName, key)
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("reservation not found")
	}

	var reservation ReservationBin
	err = attributevalue.UnmarshalMap(result, &reservation)
	if err != nil {
		return nil, err
	}

	return &reservation, nil
}

// Find all reservation bins for a given account
func (s *OffchainStore) FindReservationBins(ctx context.Context, accountID string) ([]ReservationBin, error) {
	result, err := s.dynamoClient.QueryIndex(ctx, s.reservationTableName, "AccountIDIndex", "AccountID = :accountID", commondynamodb.ExpresseionValues{
		":accountID": &types.AttributeValueMemberS{Value: accountID},
	})
	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, errors.New("reservation not found")
	}

	var reservations []ReservationBin
	err = attributevalue.UnmarshalListOfMaps(result, &reservations)
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (s *OffchainStore) AddOnDemandPayment(ctx context.Context, blobHeader BlobHeader) error {
	result, err := s.dynamoClient.GetItem(ctx, s.onDemandTableName,
		commondynamodb.Item{
			"AccountID":          &types.AttributeValueMemberS{Value: blobHeader.AccountID},
			"CumulativePayments": &types.AttributeValueMemberS{Value: strconv.FormatUint(blobHeader.CumulativePayment, 10)},
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
			"AccountID":          &types.AttributeValueMemberS{Value: blobHeader.AccountID},
			"CumulativePayments": &types.AttributeValueMemberS{Value: strconv.FormatUint(blobHeader.CumulativePayment, 10)},
			"BlobSize":           &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(blobHeader.BlobSize), 10)},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to add payment: %w", err)
	}
	return nil
}

// RemoveOnDemandPayment removes a specific payment from the list for a specific account
func (s *OffchainStore) RemoveOnDemandPayment(ctx context.Context, accountID string, payment uint64) error {
	err := s.dynamoClient.DeleteItem(ctx, s.onDemandTableName,
		commondynamodb.Key{
			"AccountID":          &types.AttributeValueMemberS{Value: accountID},
			"CumulativePayments": &types.AttributeValueMemberS{Value: strconv.FormatUint(payment, 10)},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to remove payment: %w", err)
	}

	return nil
}

// relevant on-demand payment records: previous cumulative payment, next cumulative payment, blob size of next payment
func (s *OffchainStore) GetRelevantOnDemandRecords(ctx context.Context, accountID string, cumulativePayment uint64) (uint64, uint64, uint32, error) {
	result, err := s.dynamoClient.QueryIndex(ctx, s.onDemandTableName, "AccountIDIndex", "AccountID = :account", commondynamodb.ExpresseionValues{
		":account": &types.AttributeValueMemberS{
			Value: accountID,
		}})
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to query index for account: %w", err)
	}

	// Extract the payments from the result
	var payments []PaymentTuple
	for _, item := range result {
		payment, err := strconv.ParseUint(item["CumulativePayments"].(*types.AttributeValueMemberS).Value, 10, 64)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("failed to parse payment: %w", err)
		}
		blobSize, err := strconv.ParseUint(item["BlobSize"].(*types.AttributeValueMemberN).Value, 10, 32)
		if err != nil {
			return 0, 0, 0, fmt.Errorf("failed to parse blob size: %w", err)
		}
		payments = append(payments, PaymentTuple{
			CumulativePayment: payment,
			BlobSize:          uint32(blobSize),
		})
	}
	// sort payments by cumulative payment to get neighboring payments
	sort.SliceStable(payments, func(i, j int) bool {
		return payments[i].CumulativePayment < payments[j].CumulativePayment
	})
	index := sort.Search(len(payments), func(i int) bool {
		return payments[i].CumulativePayment == cumulativePayment
	})

	var prevPayment, nextPayment uint64
	var nextBlobSize uint32

	if index > 0 {
		prevPayment = payments[index-1].CumulativePayment
	}
	if index < len(payments)-1 {
		nextPayment = payments[index+1].CumulativePayment
		nextBlobSize = payments[index+1].BlobSize
	}

	return prevPayment, nextPayment, nextBlobSize, nil
}
