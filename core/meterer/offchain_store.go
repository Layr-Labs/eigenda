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

func (s *OffchainStore) UpdateReservationBin(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64, size uint64) (uint64, error) {
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

func (s *OffchainStore) AddOnDemandPayment(ctx context.Context, paymentMetadata core.PaymentMetadata, paymentCharged *big.Int) error {
	result, err := s.dynamoClient.GetItem(ctx, s.onDemandTableName,
		commondynamodb.Item{
			"AccountID":          &types.AttributeValueMemberS{Value: paymentMetadata.AccountID.Hex()},
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
			"AccountID":          &types.AttributeValueMemberS{Value: paymentMetadata.AccountID.Hex()},
			"CumulativePayments": &types.AttributeValueMemberN{Value: paymentMetadata.CumulativePayment.String()},
			"PaymentCharged":     &types.AttributeValueMemberN{Value: paymentCharged.String()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to add payment: %w", err)
	}
	return nil
}

// RemoveOnDemandPayment removes a specific payment from the list for a specific account
func (s *OffchainStore) RemoveOnDemandPayment(ctx context.Context, accountID gethcommon.Address, payment *big.Int) error {
	err := s.dynamoClient.DeleteItem(ctx, s.onDemandTableName,
		commondynamodb.Key{
			"AccountID":          &types.AttributeValueMemberS{Value: accountID.Hex()},
			"CumulativePayments": &types.AttributeValueMemberN{Value: payment.String()},
		},
	)

	if err != nil {
		return fmt.Errorf("failed to remove payment: %w", err)
	}

	return nil
}

// GetRelevantOnDemandRecords gets previous cumulative payment, next cumulative payment, and paymentCharged
// The queries are done sequentially instead of one-go for efficient querying and would not cause race condition errors for honest requests
func (s *OffchainStore) GetRelevantOnDemandRecords(ctx context.Context, accountID gethcommon.Address, cumulativePayment *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	// Fetch the largest entry smaller than the given cumulativePayment
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(s.onDemandTableName),
		KeyConditionExpression: aws.String("AccountID = :account AND CumulativePayments < :cumulativePayment"),
		ExpressionAttributeValues: commondynamodb.ExpressionValues{
			":account":           &types.AttributeValueMemberS{Value: accountID.Hex()},
			":cumulativePayment": &types.AttributeValueMemberN{Value: cumulativePayment.String()},
		},
		ScanIndexForward: aws.Bool(false),
		Limit:            aws.Int32(1),
	}
	smallerResult, err := s.dynamoClient.QueryWithInput(ctx, queryInput)
	if err != nil {
		return nil, nil, big.NewInt(0), fmt.Errorf("failed to query smaller payments for account: %w", err)
	}
	prevPayment := big.NewInt(0)
	if len(smallerResult) > 0 {
		cumulativePaymentsAttr, ok := smallerResult[0]["CumulativePayments"]
		if !ok {
			return nil, nil, big.NewInt(0), fmt.Errorf("CumulativePayments field not found in result")
		}
		cumulativePaymentsNum, ok := cumulativePaymentsAttr.(*types.AttributeValueMemberN)
		if !ok {
			return nil, nil, big.NewInt(0), fmt.Errorf("CumulativePayments has invalid type")
		}
		setPrevPayment, success := prevPayment.SetString(cumulativePaymentsNum.Value, 10)
		if !success {
			return nil, nil, big.NewInt(0), fmt.Errorf("failed to parse previous payment: %w", err)
		}
		prevPayment = setPrevPayment
	}

	// Fetch the smallest entry larger than the given cumulativePayment
	queryInput = &dynamodb.QueryInput{
		TableName:              aws.String(s.onDemandTableName),
		KeyConditionExpression: aws.String("AccountID = :account AND CumulativePayments > :cumulativePayment"),
		ExpressionAttributeValues: commondynamodb.ExpressionValues{
			":account":           &types.AttributeValueMemberS{Value: accountID.Hex()},
			":cumulativePayment": &types.AttributeValueMemberN{Value: cumulativePayment.String()},
		},
		ScanIndexForward: aws.Bool(true),
		Limit:            aws.Int32(1),
	}
	largerResult, err := s.dynamoClient.QueryWithInput(ctx, queryInput)
	if err != nil {
		return nil, nil, big.NewInt(0), fmt.Errorf("failed to query the next payment for account: %w", err)
	}
	nextPayment := big.NewInt(0)
	paymentCharged := big.NewInt(0)
	if len(largerResult) > 0 {
		cumulativePaymentsAttr, ok := largerResult[0]["CumulativePayments"]
		if !ok {
			return nil, nil, big.NewInt(0), fmt.Errorf("CumulativePayments field not found in result")
		}
		cumulativePaymentsNum, ok := cumulativePaymentsAttr.(*types.AttributeValueMemberN)
		if !ok {
			return nil, nil, big.NewInt(0), fmt.Errorf("CumulativePayments has invalid type")
		}
		setNextPayment, success := nextPayment.SetString(cumulativePaymentsNum.Value, 10)
		if !success {
			return nil, nil, big.NewInt(0), fmt.Errorf("failed to parse previous payment: %w", err)
		}
		nextPayment = setNextPayment

		chargeAttr, ok := largerResult[0]["PaymentCharged"]
		if !ok {
			return nil, nil, big.NewInt(0), fmt.Errorf("PaymentCharged field not found in result")
		}
		chargeNum, ok := chargeAttr.(*types.AttributeValueMemberN)
		if !ok {
			return nil, nil, big.NewInt(0), fmt.Errorf("PaymentCharged has invalid type")
		}
		if _, success := paymentCharged.SetString(chargeNum.Value, 10); !success {
			return nil, nil, big.NewInt(0), fmt.Errorf("failed to parse paymentCharged value: %w", err)
		}
	}

	return prevPayment, nextPayment, paymentCharged, nil
}

func (s *OffchainStore) GetPeriodRecords(ctx context.Context, accountID gethcommon.Address, reservationPeriod uint64) ([MinNumBins]*pb.PeriodRecord, error) {
	// Fetch the 3 bins start from the current bin
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(s.reservationTableName),
		KeyConditionExpression: aws.String("AccountID = :account AND ReservationPeriod > :reservationPeriod"),
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

func (s *OffchainStore) GetLargestCumulativePayment(ctx context.Context, accountID gethcommon.Address) (*big.Int, error) {
	// Fetch the largest cumulative payment
	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(s.onDemandTableName),
		KeyConditionExpression: aws.String("AccountID = :account"),
		ExpressionAttributeValues: commondynamodb.ExpressionValues{
			":account": &types.AttributeValueMemberS{Value: accountID.Hex()},
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

	// Safely extract CumulativePayments
	cumulativePaymentsAttr, ok := payments[0]["CumulativePayments"]
	if !ok {
		return nil, fmt.Errorf("CumulativePayments field not found in result")
	}

	// Type assertion with check
	cumulativePaymentsNum, ok := cumulativePaymentsAttr.(*types.AttributeValueMemberN)
	if !ok {
		return nil, fmt.Errorf("CumulativePayments has invalid type: %T", cumulativePaymentsAttr)
	}

	payment := new(big.Int)
	if _, success := payment.SetString(cumulativePaymentsNum.Value, 10); !success {
		return nil, fmt.Errorf("failed to parse payment value: %s", cumulativePaymentsNum.Value)
	}

	return payment, nil
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
