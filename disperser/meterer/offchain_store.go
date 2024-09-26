package meterer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
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
	AccountID string    `dynamodbav:"AccountID"`
	BinIndex  uint64    `dynamodbav:"BinIndex"`
	BinUsage  uint32    `dynamodbav:"BinUsage"`
	UpdatedAt time.Time `dynamodbav:"UpdatedAt"`
}

type PaymentTuple struct {
	CumulativePayment uint64
	BlobSize          uint32
}

type OnDemandCollection struct {
	AccountID string `dynamodbav:"AccountID"`
	// Payments  []PaymentTuple `dynamodbav:"payments"`
	CumulativePayments map[uint64]uint32 `dynamodbav:"CumulativePayments"` // Payment is the key, BlobSize is the value
}

type GlobalBin struct {
	BinIndex  uint64    `dynamodbav:"BinIndex"`
	BinUsage  uint64    `dynamodbav:"BinUsage"`
	UpdatedAt time.Time `dynamodbav:"UpdatedAt"`
}

// func (s *OffchainStore) UpdateReservationBin(ctx context.Context, accountID string, binIndex uint64, size uint32) (uint32, error) {
func (s *OffchainStore) UpdateReservationBin(ctx context.Context, accountID string, binIndex uint64, size uint32) (uint32, error) {
	key := map[string]types.AttributeValue{
		"AccountID": &types.AttributeValueMemberS{Value: accountID},
		"BinIndex":  &types.AttributeValueMemberN{Value: strconv.FormatUint(binIndex, 10)},
	}

	update := map[string]types.AttributeValue{
		"BinUsage": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(size), 10)},
	}

	fmt.Printf("updating item %s %+v\n", "BinUsage", update["BinUsage"])

	// fmt.Println("updating reservation bin", accountID, binIndex, size, s.reservationTableName)
	res, err := s.dynamoClient.UpdateItemIncrement(ctx, s.reservationTableName, key, update)
	// res, err := s.dynamoClient.UpdateItem(ctx, s.reservationTableName, commondynamodb.Key{
	// 	"AccountID": &types.AttributeValueMemberS{Value: accountID},
	// 	"BinIndex":  &types.AttributeValueMemberN{Value: strconv.FormatUint(binIndex, 10)},
	// }, commondynamodb.Item{
	// 	"BinUsage": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(size), 10)},
	// })
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

// func (s *OffchainStore) UpdateReservationBin(ctx context.Context, accountID string, binIndex uint64, size uint32) (uint32, error) {
func (s *OffchainStore) UpdateGlobalBin(ctx context.Context, binIndex uint32, size uint32) (uint64, error) {
	key := map[string]types.AttributeValue{
		"BinIndex": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(binIndex), 10)},
	}

	update := map[string]types.AttributeValue{
		":inc": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(size), 10)},
		":now": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)},
	}

	// updateExpression := "ADD BinUsage :inc SET UpdatedAt = :now"

	res, err := s.dynamoClient.UpdateItem(ctx, s.globalBinTableName, key, update)
	if err != nil {
		return 0, err
	}

	binUsage, ok := res["BinUsage"]
	if !ok {
		return 0, nil // Return 0 if BinUsage is not present in the response
	}

	binUsageAttr, ok := binUsage.(*types.AttributeValueMemberN)
	if !ok {
		return 0, nil // Return 0 if BinUsage is not of the expected type
	}

	binUsageValue, err := strconv.ParseUint(binUsageAttr.Value, 10, 32)
	if err != nil {
		return 0, err
	}

	return binUsageValue, nil
}

// func (s *OffchainStore) RemoveOutdatedBins(ctx context.Context, currentBinIndex uint64) error {
// 	var processError error
// 	err := s.dynamoClient.ScanPaginator(ctx, s.reservationTableName, nil,
// 		func(page *dynamodb.ScanOutput) bool {
// 			for _, item := range page.Items {
// 				var usage ReservationBin
// 				err := attributevalue.UnmarshalMap(item, &usage)
// 				if err != nil {
// 					processError = err
// 					return false
// 				}

// 				if usage.BinIndex < currentBinIndex-2 {
// 					key := map[string]types.AttributeValue{
// 						"AccountID": &types.AttributeValueMemberS{Value: usage.AccountID},
// 						"BinIndex":  &types.AttributeValueMemberN{Value: strconv.FormatUint(usage.BinIndex, 10)},
// 					}
// 					err := s.dynamoClient.DeleteItem(ctx, s.reservationTableName, key)
// 					if err != nil {
// 						processError = err
// 						return false
// 					}
// 				}
// 			}
// 			return true
// 		})

//		if err != nil {
//			return err
//		}
//		return processError
//	}

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

// TODO: could we let the return type be a pointer to a slice of ReservationBins?
func (s *OffchainStore) FindReservationBins(ctx context.Context, accountID string) ([]ReservationBin, error) {
	// key := .Key{
	// 	"AccountID": &types.AttributeValueMemberS{Value: accountID},
	// }

	// keys := make([]commondynamodb.Key, numItems)
	// for i := 0; i < numItems; i += 1 {
	// 	keys[i] = commondynamodb.Key{
	// 		"MetadataKey": &types.AttributeValueMemberS{Value: fmt.Sprintf("key%d", i)},
	// 	}
	// }
	// items, err := s.dynamoDBClient.QueryIndex(ctx, s.tableName, statusIndexName, "BlobStatus = :status", commondynamodb.ExpresseionValues{
	// 	":status": &types.AttributeValueMemberN{
	// 		Value: strconv.Itoa(int(status)),
	// 	}})

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

func (s *OffchainStore) AddOnDemandPayment(ctx context.Context, blobHeader BlobHeader) (*big.Int, error) {
	key := map[string]types.AttributeValue{
		"account_id": &types.AttributeValueMemberS{Value: blobHeader.AccountID},
	}
	// // Update expression to add the payment with its associated blob size
	// updateExpression := "SET payment_blobs.#payment = :blobSize ADD CumulativePayments :payment"
	// expressionAttributeNames := map[string]string{
	// 	"#payment": fmt.Sprintf("%d", blobHeader.CumulativePayment),
	// }
	// expressionAttributeValues := map[string]types.AttributeValue{
	// 	":blobSize": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", blobHeader.BlobSize)},
	// 	":payment":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", blobHeader.CumulativePayment)},
	// }
	update := map[string]types.AttributeValue{
		":blobSize": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", blobHeader.BlobSize)},
		":payment":  &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", blobHeader.CumulativePayment)},
	}

	// Perform the update
	res, err := s.dynamoClient.UpdateItem(ctx, s.onDemandTableName, key, update)

	if err != nil {
		return nil, fmt.Errorf("failed to add payment: %w", err)
	}

	cumulativePayments, ok := res["CumulativePayments"]
	if !ok {
		return nil, fmt.Errorf("CumulativePayments attribute not found in the response")
	}

	cumulativePaymentsN, ok := cumulativePayments.(*types.AttributeValueMemberN)
	if !ok {
		return nil, fmt.Errorf("CumulativePayments attribute is not of type AttributeValueMemberN")
	}

	cumulativePaymentsValue, err := strconv.ParseUint(cumulativePaymentsN.Value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CumulativePayments value: %w", err)
	}

	return big.NewInt(0).SetUint64(cumulativePaymentsValue), nil

	// _, err := s.dynamoClient.UpdateItem(ctx, s.onDemandTableName, key, update)

	// if err != nil {
	// 	return err
	// }

	// return nil
}

// RemoveOnDemandPayment removes a specific payment from the list for a specific account
func (s *OffchainStore) RemoveOnDemandPayment(ctx context.Context, accountID string, payment uint64) error {
	key := map[string]types.AttributeValue{
		"account_id": &types.AttributeValueMemberS{Value: accountID},
	}
	item := map[string]types.AttributeValue{
		"payments": &types.AttributeValueMemberNS{Value: []string{fmt.Sprintf("%d", payment)}},
	}

	_, err := s.dynamoClient.UpdateItem(ctx, s.onDemandTableName, key, item)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

// Add this function to get relevant on-demand payment records
func (s *OffchainStore) GetRelevantOnDemandRecords(ctx context.Context, accountID string, cumulativePayment uint64) (uint64, uint64, uint32, error) {
	// Query the DynamoDB table for all payments of the account
	key := map[string]types.AttributeValue{
		"account_id": &types.AttributeValueMemberS{Value: accountID},
	}

	result, err := s.dynamoClient.GetItem(ctx, s.onDemandTableName, key)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to get item: %w", err)
	}

	// Extract the payments from the result
	var payments []OnDemandCollection
	if paymentsAttr, ok := result["payments"]; ok {
		if paymentsSet, ok := paymentsAttr.(*types.AttributeValueMemberNS); ok {
			for _, paymentStr := range paymentsSet.Value {
				payment, err := strconv.ParseUint(paymentStr, 10, 64)
				if err != nil {
					return 0, 0, 0, fmt.Errorf("failed to parse payment: %w", err)
				}
				payments = append(payments, OnDemandCollection{CumulativePayments: map[uint64]uint32{payment: 0}})
			}
		}
	}

	// Binary search to find the insertion point
	index := sort.Search(len(payments), func(i int) bool {
		return payments[i].CumulativePayments[cumulativePayment] == 0
	})

	var prevPayment, nextPayment uint64
	var nextBlobSize uint32

	if index > 0 {
		for payment, _ := range payments[index-1].CumulativePayments {
			prevPayment = payment
			break
		}
	}

	if index < len(payments)-1 {
		for payment, blobSize := range payments[index+1].CumulativePayments {
			nextPayment = payment
			nextBlobSize = blobSize
			break
		}
	}

	return prevPayment, nextPayment, nextBlobSize, nil
}
