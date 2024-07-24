package batchstore

import (
	"context"
	"fmt"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser/batcher"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

const (
	batchSK             = "BATCH#"
	minibatchSK         = "MINIBATCH#"
	dispersalRequestSK  = "DISPERSAL_REQUEST#"
	dispersalResponseSK = "DISPERSAL_RESPONSE#"
)

type DynamoBatchRecord struct {
	BatchID              string
	CreatedAt            time.Time
	ReferenceBlockNumber uint
	HeaderHash           [32]byte
	AggregatePubKey      *core.G2Point
	AggregateSignature   *core.Signature
}

type DynamoMinibatchRecord struct {
	BatchID              string
	MinibatchIndex       uint
	BlobHeaderHashes     [][32]byte
	BatchSize            uint64
	ReferenceBlockNumber uint
}

type DynamoDispersalRequest struct {
	BatchID         string
	MinibatchIndex  uint
	OperatorID      string
	OperatorAddress string
	NumBlobs        uint
	RequestedAt     time.Time
}

type DynamoDispersalResponse struct {
	DynamoDispersalRequest
	Signatures  []*core.Signature
	RespondedAt time.Time
	Error       error
}
type MinibatchStore struct {
	dynamoDBClient *commondynamodb.Client
	tableName      string
	logger         logging.Logger
	ttl            time.Duration
}

var _ batcher.MinibatchStore = (*MinibatchStore)(nil)

func NewMinibatchStore(dynamoDBClient *commondynamodb.Client, logger logging.Logger, tableName string, ttl time.Duration) *MinibatchStore {
	logger.Debugf("creating minibatch store with table %s with TTL: %s", tableName, ttl)
	return &MinibatchStore{
		dynamoDBClient: dynamoDBClient,
		logger:         logger.With("component", "MinibatchStore"),
		tableName:      tableName,
		ttl:            ttl,
	}
}

func GenerateTableSchema(tableName string, readCapacityUnits int64, writeCapacityUnits int64) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("BatchID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("OperatorID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("RequestedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("BatchID"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("SK"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName: aws.String(tableName),
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("OperatorID_RequestedAt_Index"),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("OperatorID"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("RequestedAt"),
						KeyType:       types.KeyTypeRange,
					},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeAll,
				},
				ProvisionedThroughput: &types.ProvisionedThroughput{
					ReadCapacityUnits:  aws.Int64(readCapacityUnits),
					WriteCapacityUnits: aws.Int64(writeCapacityUnits),
				},
			},
		},
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(readCapacityUnits),
			WriteCapacityUnits: aws.Int64(writeCapacityUnits),
		},
	}
}

func ToDynamoBatchRecord(br batcher.BatchRecord) DynamoBatchRecord {
	return DynamoBatchRecord{
		BatchID:              br.ID.String(),
		CreatedAt:            br.CreatedAt,
		ReferenceBlockNumber: br.ReferenceBlockNumber,
		HeaderHash:           br.HeaderHash,
		AggregatePubKey:      br.AggregatePubKey,
		AggregateSignature:   br.AggregateSignature,
	}
}

func ToDynamoMinibatchRecord(br batcher.MinibatchRecord) DynamoMinibatchRecord {
	return DynamoMinibatchRecord{
		BatchID:              br.BatchID.String(),
		MinibatchIndex:       br.MinibatchIndex,
		BlobHeaderHashes:     br.BlobHeaderHashes,
		BatchSize:            br.BatchSize,
		ReferenceBlockNumber: br.ReferenceBlockNumber,
	}
}

func ToDynamoDispersalRequest(dr batcher.DispersalRequest) DynamoDispersalRequest {
	return DynamoDispersalRequest{
		BatchID:         dr.BatchID.String(),
		MinibatchIndex:  dr.MinibatchIndex,
		OperatorID:      dr.OperatorID.Hex(),
		OperatorAddress: dr.OperatorAddress.Hex(),
		NumBlobs:        dr.NumBlobs,
		RequestedAt:     dr.RequestedAt,
	}
}

func ToDynamoDispersalResponse(dr batcher.DispersalResponse) DynamoDispersalResponse {
	return DynamoDispersalResponse{
		DynamoDispersalRequest: ToDynamoDispersalRequest(dr.DispersalRequest),
		Signatures:             dr.Signatures,
		RespondedAt:            dr.RespondedAt,
		Error:                  dr.Error,
	}
}

func FromDynamoBatchRecord(dbr DynamoBatchRecord) (batcher.BatchRecord, error) {
	batchID, err := uuid.Parse(dbr.BatchID)
	if err != nil {
		return batcher.BatchRecord{}, fmt.Errorf("failed to convert dynamo batch record batch ID %v from string: %v", dbr.BatchID, err)
	}

	return batcher.BatchRecord{
		ID:                   batchID,
		CreatedAt:            dbr.CreatedAt,
		ReferenceBlockNumber: dbr.ReferenceBlockNumber,
		HeaderHash:           dbr.HeaderHash,
		AggregatePubKey:      dbr.AggregatePubKey,
		AggregateSignature:   dbr.AggregateSignature,
	}, nil
}

func FromDynamoMinibatchRecord(dbr DynamoMinibatchRecord) (batcher.MinibatchRecord, error) {
	batchID, err := uuid.Parse(dbr.BatchID)
	if err != nil {
		return batcher.MinibatchRecord{}, fmt.Errorf("failed to convert dynamo minibatch record batch ID %v from string: %v", dbr.BatchID, err)
	}

	return batcher.MinibatchRecord{
		BatchID:              batchID,
		MinibatchIndex:       dbr.MinibatchIndex,
		BlobHeaderHashes:     dbr.BlobHeaderHashes,
		BatchSize:            dbr.BatchSize,
		ReferenceBlockNumber: dbr.ReferenceBlockNumber,
	}, nil
}

func FromDynamoDispersalRequest(ddr DynamoDispersalRequest) (batcher.DispersalRequest, error) {
	batchID, err := uuid.Parse(ddr.BatchID)
	if err != nil {
		return batcher.DispersalRequest{}, fmt.Errorf("failed to convert dynamo dispersal request batch ID %v from string: %v", ddr.BatchID, err)
	}
	operatorID, err := core.OperatorIDFromHex(ddr.OperatorID)
	if err != nil {
		return batcher.DispersalRequest{}, fmt.Errorf("failed to convert dynamo dispersal request operator ID %v from hex: %v", ddr.OperatorID, err)
	}

	return batcher.DispersalRequest{
		BatchID:         batchID,
		MinibatchIndex:  ddr.MinibatchIndex,
		OperatorID:      operatorID,
		OperatorAddress: gcommon.HexToAddress(ddr.OperatorAddress),
		NumBlobs:        ddr.NumBlobs,
		RequestedAt:     ddr.RequestedAt,
	}, nil
}

func FromDynamoDispersalResponse(ddr DynamoDispersalResponse) (batcher.DispersalResponse, error) {
	request, err := FromDynamoDispersalRequest(ddr.DynamoDispersalRequest)
	if err != nil {
		return batcher.DispersalResponse{}, err
	}

	return batcher.DispersalResponse{
		DispersalRequest: request,
		Signatures:       ddr.Signatures,
		RespondedAt:      ddr.RespondedAt,
		Error:            ddr.Error,
	}, nil
}

func MarshalBatchRecord(batch *batcher.BatchRecord) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(ToDynamoBatchRecord(*batch))
	if err != nil {
		return nil, err
	}
	fields["SK"] = &types.AttributeValueMemberS{Value: batchSK + batch.ID.String()}
	fields["CreatedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", batch.CreatedAt.UTC().Unix())}
	return fields, nil
}

func MarshalMinibatchRecord(minibatch *batcher.MinibatchRecord) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(ToDynamoMinibatchRecord(*minibatch))
	if err != nil {
		return nil, err
	}
	fields["SK"] = &types.AttributeValueMemberS{Value: minibatchSK + fmt.Sprintf("%d", minibatch.MinibatchIndex)}
	return fields, nil
}

func MarshalDispersalRequest(request *batcher.DispersalRequest) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(ToDynamoDispersalRequest(*request))
	if err != nil {
		return nil, err
	}
	fields["SK"] = &types.AttributeValueMemberS{Value: dispersalRequestSK + fmt.Sprintf("%d#%s", request.MinibatchIndex, request.OperatorID.Hex())}
	fields["RequestedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", request.RequestedAt.UTC().Unix())}
	return fields, nil
}

func MarshalDispersalResponse(response *batcher.DispersalResponse) (map[string]types.AttributeValue, error) {
	fields, err := attributevalue.MarshalMap(ToDynamoDispersalResponse(*response))
	if err != nil {
		return nil, err
	}
	fields["SK"] = &types.AttributeValueMemberS{Value: dispersalResponseSK + fmt.Sprintf("%d#%s", response.MinibatchIndex, response.OperatorID.Hex())}
	fields["RespondedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", response.RespondedAt.UTC().Unix())}
	fields["RequestedAt"] = &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", response.DispersalRequest.RequestedAt.UTC().Unix())}
	return fields, nil
}

func UnmarshalBatchRecord(item commondynamodb.Item) (*batcher.BatchRecord, error) {
	dbr := DynamoBatchRecord{}
	err := attributevalue.UnmarshalMap(item, &dbr)
	if err != nil {
		return nil, err
	}

	batch, err := FromDynamoBatchRecord(dbr)
	if err != nil {
		return nil, err
	}

	batch.CreatedAt = batch.CreatedAt.UTC()
	return &batch, nil
}

func UnmarshalMinibatchRecord(item commondynamodb.Item) (*batcher.MinibatchRecord, error) {
	dbr := DynamoMinibatchRecord{}
	err := attributevalue.UnmarshalMap(item, &dbr)
	if err != nil {
		return nil, err
	}

	minibatch, err := FromDynamoMinibatchRecord(dbr)
	if err != nil {
		return nil, err
	}
	return &minibatch, nil
}

func UnmarshalDispersalRequest(item commondynamodb.Item) (*batcher.DispersalRequest, error) {
	ddr := DynamoDispersalRequest{}
	err := attributevalue.UnmarshalMap(item, &ddr)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dispersal request from DynamoDB: %v", err)
	}

	request, err := FromDynamoDispersalRequest(ddr)
	if err != nil {
		return nil, err
	}

	request.RequestedAt = request.RequestedAt.UTC()
	return &request, nil
}

func UnmarshalDispersalResponse(item commondynamodb.Item) (*batcher.DispersalResponse, error) {
	ddr := DynamoDispersalResponse{}
	err := attributevalue.UnmarshalMap(item, &ddr)
	if err != nil {
		return nil, err
	}

	response, err := FromDynamoDispersalResponse(ddr)
	if err != nil {
		return nil, err
	}
	response.RespondedAt = response.RespondedAt.UTC()
	response.DispersalRequest.RequestedAt = response.DispersalRequest.RequestedAt.UTC()
	return &response, nil
}

func (m *MinibatchStore) PutBatch(ctx context.Context, batch *batcher.BatchRecord) error {
	item, err := MarshalBatchRecord(batch)
	if err != nil {
		return err
	}

	return m.dynamoDBClient.PutItem(ctx, m.tableName, item)
}

func (m *MinibatchStore) PutMinibatch(ctx context.Context, minibatch *batcher.MinibatchRecord) error {
	item, err := MarshalMinibatchRecord(minibatch)
	if err != nil {
		return err
	}

	return m.dynamoDBClient.PutItem(ctx, m.tableName, item)
}

func (m *MinibatchStore) PutMiniBatch(ctx context.Context, minibatch *batcher.MinibatchRecord) error {
	return m.PutMinibatch(ctx, minibatch)
}

func (m *MinibatchStore) PutDispersalRequest(ctx context.Context, request *batcher.DispersalRequest) error {
	item, err := MarshalDispersalRequest(request)
	if err != nil {
		return err
	}

	fmt.Printf("%v", item)
	return m.dynamoDBClient.PutItem(ctx, m.tableName, item)
}

func (m *MinibatchStore) PutDispersalResponse(ctx context.Context, response *batcher.DispersalResponse) error {
	item, err := MarshalDispersalResponse(response)
	if err != nil {
		return err
	}

	return m.dynamoDBClient.PutItem(ctx, m.tableName, item)
}

func (m *MinibatchStore) GetBatch(ctx context.Context, batchID uuid.UUID) (*batcher.BatchRecord, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"BatchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: batchSK + batchID.String(),
		},
	})
	if err != nil {
		m.logger.Errorf("failed to get batch from DynamoDB: %v", err)
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	batch, err := UnmarshalBatchRecord(item)
	if err != nil {
		m.logger.Errorf("failed to unmarshal batch record from DynamoDB: %v", err)
		return nil, err
	}
	return batch, nil
}

// GetPendingBatch implements batcher.MinibatchStore.
func (m *MinibatchStore) GetPendingBatch(ctx context.Context) (*batcher.BatchRecord, error) {
	panic("unimplemented")
}

func (m *MinibatchStore) GetMinibatch(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) (*batcher.MinibatchRecord, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"BatchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: minibatchSK + fmt.Sprintf("%d", minibatchIndex),
		},
	})
	if err != nil {
		m.logger.Errorf("failed to get minibatch from DynamoDB: %v", err)
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	minibatch, err := UnmarshalMinibatchRecord(item)
	if err != nil {
		m.logger.Errorf("failed to unmarshal minibatch record from DynamoDB: %v", err)
		return nil, err
	}
	return minibatch, nil
}

func (m *MinibatchStore) GetMiniBatch(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) (*batcher.MinibatchRecord, error) {
	return m.GetMinibatch(ctx, batchID, minibatchIndex)
}

func (m *MinibatchStore) GetDispersalRequest(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*batcher.DispersalRequest, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"BatchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: dispersalRequestSK + fmt.Sprintf("%d#%s", minibatchIndex, opID.Hex()),
		},
	})
	if err != nil {
		m.logger.Errorf("failed to get dispersal request from DynamoDB: %v", err)
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	request, err := UnmarshalDispersalRequest(item)
	if err != nil {
		m.logger.Errorf("failed to unmarshal dispersal request from DynamoDB: %v", err)
		return nil, err
	}
	return request, nil
}

func (m *MinibatchStore) GetDispersalRequests(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*batcher.DispersalRequest, error) {
	items, err := m.dynamoDBClient.Query(ctx, m.tableName, "BatchID = :batchID AND SK = :sk", commondynamodb.ExpresseionValues{
		":batchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		":sk": &types.AttributeValueMemberS{
			Value: dispersalRequestSK + fmt.Sprintf("%s#%d", batchID.String(), minibatchIndex),
		},
	})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no dispersal requests found for BatchID %s MinibatchIndex %d", batchID, minibatchIndex)
	}

	requests := make([]*batcher.DispersalRequest, len(items))
	for i, item := range items {
		requests[i], err = UnmarshalDispersalRequest(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal dispersal requests at index %d: %v", i, err)
			return nil, err
		}
	}

	return requests, nil
}

func (m *MinibatchStore) GetDispersalResponse(ctx context.Context, batchID uuid.UUID, minibatchIndex uint, opID core.OperatorID) (*batcher.DispersalResponse, error) {
	item, err := m.dynamoDBClient.GetItem(ctx, m.tableName, map[string]types.AttributeValue{
		"BatchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: dispersalResponseSK + fmt.Sprintf("%d#%s", minibatchIndex, opID.Hex()),
		},
	})
	if err != nil {
		m.logger.Errorf("failed to get dispersal response from DynamoDB: %v", err)
		return nil, err
	}

	if item == nil {
		return nil, nil
	}

	response, err := UnmarshalDispersalResponse(item)
	if err != nil {
		m.logger.Errorf("failed to unmarshal dispersal response from DynamoDB: %v", err)
		return nil, err
	}
	return response, nil
}

func (m *MinibatchStore) GetDispersalResponses(ctx context.Context, batchID uuid.UUID, minibatchIndex uint) ([]*batcher.DispersalResponse, error) {
	items, err := m.dynamoDBClient.Query(ctx, m.tableName, "BatchID = :batchID AND SK = :sk", commondynamodb.ExpresseionValues{
		":batchID": &types.AttributeValueMemberS{
			Value: batchID.String(),
		},
		":sk": &types.AttributeValueMemberS{
			Value: dispersalResponseSK + fmt.Sprintf("%s#%d", batchID.String(), minibatchIndex),
		},
	})
	if err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no dispersal responses found for BatchID %s MinibatchIndex %d", batchID, minibatchIndex)
	}

	responses := make([]*batcher.DispersalResponse, len(items))
	for i, item := range items {
		responses[i], err = UnmarshalDispersalResponse(item)
		if err != nil {
			m.logger.Errorf("failed to unmarshal dispersal response at index %d: %v", i, err)
			return nil, err
		}
	}

	return responses, nil
}
