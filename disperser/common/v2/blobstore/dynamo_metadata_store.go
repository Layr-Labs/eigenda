package blobstore

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	commondynamodb "github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common"
	v2 "github.com/Layr-Labs/eigenda/disperser/common/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	StatusIndexName            = "StatusIndex"
	OperatorDispersalIndexName = "OperatorDispersalIndex"
	OperatorResponseIndexName  = "OperatorResponseIndex"

	blobKeyPrefix             = "BlobKey#"
	dispersalKeyPrefix        = "Dispersal#"
	blobMetadataSK            = "BlobMetadata"
	blobCertSK                = "BlobCertificate"
	dispersalRequestSKPrefix  = "DispersalRequest#"
	dispersalResponseSKPrefix = "DispersalResponse#"
)

var (
	statusUpdatePrecondition = map[v2.BlobStatus][]v2.BlobStatus{
		v2.Queued:    {},
		v2.Encoded:   {v2.Queued},
		v2.Certified: {v2.Encoded},
		v2.Failed:    {v2.Queued, v2.Encoded},
	}
	ErrInvalidStateTransition = errors.New("invalid state transition")
)

// BlobMetadataStore is a blob metadata storage backed by DynamoDB
type BlobMetadataStore struct {
	dynamoDBClient *commondynamodb.Client
	logger         logging.Logger
	tableName      string
}

func NewBlobMetadataStore(dynamoDBClient *commondynamodb.Client, logger logging.Logger, tableName string) *BlobMetadataStore {
	logger.Debugf("creating blob metadata store v2 with table %s", tableName)
	return &BlobMetadataStore{
		dynamoDBClient: dynamoDBClient,
		logger:         logger.With("component", "blobMetadataStoreV2"),
		tableName:      tableName,
	}
}

func (s *BlobMetadataStore) PutBlobMetadata(ctx context.Context, blobMetadata *v2.BlobMetadata) error {
	item, err := MarshalBlobMetadata(blobMetadata)
	if err != nil {
		return err
	}

	err = s.dynamoDBClient.PutItemWithCondition(ctx, s.tableName, item, "attribute_not_exists(PK) AND attribute_not_exists(SK)", nil, nil)
	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		return common.ErrAlreadyExists
	}

	return err
}

func (s *BlobMetadataStore) UpdateBlobStatus(ctx context.Context, blobKey corev2.BlobKey, status v2.BlobStatus) error {
	validStatuses := statusUpdatePrecondition[status]
	if len(validStatuses) == 0 {
		return fmt.Errorf("%w: invalid status transition to %s", ErrInvalidStateTransition, status.String())
	}

	expValues := make([]expression.OperandBuilder, len(validStatuses))
	for i, validStatus := range validStatuses {
		expValues[i] = expression.Value(int(validStatus))
	}
	condition := expression.Name("BlobStatus").In(expValues[0], expValues[1:]...)
	_, err := s.dynamoDBClient.UpdateItemWithCondition(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: blobKeyPrefix + blobKey.Hex(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: blobMetadataSK,
		},
	}, map[string]types.AttributeValue{
		"BlobStatus": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
		"UpdatedAt": &types.AttributeValueMemberN{
			Value: strconv.FormatInt(time.Now().UnixNano(), 10),
		},
	}, condition)

	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		blob, err := s.GetBlobMetadata(ctx, blobKey)
		if err != nil {
			s.logger.Errorf("failed to get blob metadata for key %s: %v", blobKey.Hex(), err)
		}

		if blob.BlobStatus == status {
			return fmt.Errorf("%w: blob already in status %s", common.ErrAlreadyExists, status.String())
		}

		return fmt.Errorf("%w: invalid status transition to %s", ErrInvalidStateTransition, status.String())
	}

	return err
}

func (s *BlobMetadataStore) GetBlobMetadata(ctx context.Context, blobKey corev2.BlobKey) (*v2.BlobMetadata, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: blobKeyPrefix + blobKey.Hex(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: blobMetadataSK,
		},
	})

	if item == nil {
		return nil, fmt.Errorf("%w: metadata not found for key %s", common.ErrMetadataNotFound, blobKey.Hex())
	}

	if err != nil {
		return nil, err
	}

	metadata, err := UnmarshalBlobMetadata(item)
	if err != nil {
		return nil, err
	}

	return metadata, nil
}

// GetBlobMetadataByStatus returns all the metadata with the given status that were updated after lastUpdatedAt
// Because this function scans the entire index, it should only be used for status with a limited number of items.
func (s *BlobMetadataStore) GetBlobMetadataByStatus(ctx context.Context, status v2.BlobStatus, lastUpdatedAt uint64) ([]*v2.BlobMetadata, error) {
	items, err := s.dynamoDBClient.QueryIndex(ctx, s.tableName, StatusIndexName, "BlobStatus = :status AND UpdatedAt > :updatedAt", commondynamodb.ExpressionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
		":updatedAt": &types.AttributeValueMemberN{
			Value: strconv.FormatInt(time.Now().Unix(), 10),
		}})
	if err != nil {
		return nil, err
	}

	metadata := make([]*v2.BlobMetadata, len(items))
	for i, item := range items {
		metadata[i], err = UnmarshalBlobMetadata(item)
		if err != nil {
			return nil, err
		}
	}

	return metadata, nil
}

// GetBlobMetadataCountByStatus returns the count of all the metadata with the given status
// Because this function scans the entire index, it should only be used for status with a limited number of items.
func (s *BlobMetadataStore) GetBlobMetadataCountByStatus(ctx context.Context, status v2.BlobStatus) (int32, error) {
	count, err := s.dynamoDBClient.QueryIndexCount(ctx, s.tableName, StatusIndexName, "BlobStatus = :status", commondynamodb.ExpressionValues{
		":status": &types.AttributeValueMemberN{
			Value: strconv.Itoa(int(status)),
		},
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (s *BlobMetadataStore) PutBlobCertificate(ctx context.Context, blobCert *corev2.BlobCertificate, fragmentInfo *encoding.FragmentInfo) error {
	item, err := MarshalBlobCertificate(blobCert, fragmentInfo)
	if err != nil {
		return err
	}

	err = s.dynamoDBClient.PutItemWithCondition(ctx, s.tableName, item, "attribute_not_exists(PK) AND attribute_not_exists(SK)", nil, nil)
	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		return common.ErrAlreadyExists
	}

	return err
}

func (s *BlobMetadataStore) GetBlobCertificate(ctx context.Context, blobKey corev2.BlobKey) (*corev2.BlobCertificate, *encoding.FragmentInfo, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: blobKeyPrefix + blobKey.Hex(),
		},
		"SK": &types.AttributeValueMemberS{
			Value: blobCertSK,
		},
	})

	if err != nil {
		return nil, nil, err
	}

	if item == nil {
		return nil, nil, fmt.Errorf("%w: certificate not found for key %s", common.ErrMetadataNotFound, blobKey.Hex())
	}

	cert, fragmentInfo, err := UnmarshalBlobCertificate(item)
	if err != nil {
		return nil, nil, err
	}

	return cert, fragmentInfo, nil
}

// GetBlobCertificates returns the certificates for the given blob keys
// Note: the returned certificates are NOT necessarily ordered by the order of the input blob keys
func (s *BlobMetadataStore) GetBlobCertificates(ctx context.Context, blobKeys []corev2.BlobKey) ([]*corev2.BlobCertificate, []*encoding.FragmentInfo, error) {
	keys := make([]map[string]types.AttributeValue, len(blobKeys))
	for i, blobKey := range blobKeys {
		keys[i] = map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: blobKeyPrefix + blobKey.Hex(),
			},
			"SK": &types.AttributeValueMemberS{
				Value: blobCertSK,
			},
		}
	}

	items, err := s.dynamoDBClient.GetItems(ctx, s.tableName, keys)
	if err != nil {
		return nil, nil, err
	}

	certs := make([]*corev2.BlobCertificate, len(items))
	fragmentInfos := make([]*encoding.FragmentInfo, len(items))
	for i, item := range items {
		cert, fragmentInfo, err := UnmarshalBlobCertificate(item)
		if err != nil {
			return nil, nil, err
		}
		certs[i] = cert
		fragmentInfos[i] = fragmentInfo
	}

	return certs, fragmentInfos, nil
}

func (s *BlobMetadataStore) PutDispersalRequest(ctx context.Context, req *corev2.DispersalRequest) error {
	item, err := MarshalDispersalRequest(req)
	if err != nil {
		return err
	}

	err = s.dynamoDBClient.PutItemWithCondition(ctx, s.tableName, item, "attribute_not_exists(PK) AND attribute_not_exists(SK)", nil, nil)
	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		return common.ErrAlreadyExists
	}

	return err
}

func (s *BlobMetadataStore) GetDispersalRequest(ctx context.Context, batchHeaderHash [32]byte, operatorID core.OperatorID) (*corev2.DispersalRequest, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: dispersalKeyPrefix + hex.EncodeToString(batchHeaderHash[:]),
		},
		"SK": &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%s%s", dispersalRequestSKPrefix, operatorID.Hex()),
		},
	})

	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, fmt.Errorf("%w: dispersal request not found for batch header hash %x and operator %s", common.ErrMetadataNotFound, batchHeaderHash, operatorID.Hex())
	}

	req, err := UnmarshalDispersalRequest(item)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (s *BlobMetadataStore) PutDispersalResponse(ctx context.Context, res *corev2.DispersalResponse) error {
	item, err := MarshalDispersalResponse(res)
	if err != nil {
		return err
	}

	err = s.dynamoDBClient.PutItemWithCondition(ctx, s.tableName, item, "attribute_not_exists(PK) AND attribute_not_exists(SK)", nil, nil)
	if errors.Is(err, commondynamodb.ErrConditionFailed) {
		return common.ErrAlreadyExists
	}

	return err
}

func (s *BlobMetadataStore) GetDispersalResponse(ctx context.Context, batchHeaderHash [32]byte, operatorID core.OperatorID) (*corev2.DispersalResponse, error) {
	item, err := s.dynamoDBClient.GetItem(ctx, s.tableName, map[string]types.AttributeValue{
		"PK": &types.AttributeValueMemberS{
			Value: dispersalKeyPrefix + hex.EncodeToString(batchHeaderHash[:]),
		},
		"SK": &types.AttributeValueMemberS{
			Value: fmt.Sprintf("%s%s", dispersalResponseSKPrefix, operatorID.Hex()),
		},
	})

	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, fmt.Errorf("%w: dispersal response not found for batch header hash %x and operator %s", common.ErrMetadataNotFound, batchHeaderHash, operatorID.Hex())
	}

	res, err := UnmarshalDispersalResponse(item)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GenerateTableSchema(tableName string, readCapacityUnits int64, writeCapacityUnits int64) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			// PK is the composite partition key
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			// SK is the composite sort key
			{
				AttributeName: aws.String("SK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("BlobStatus"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("UpdatedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("OperatorID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("DispersedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("RespondedAt"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("PK"),
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
				IndexName: aws.String(StatusIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("BlobStatus"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("UpdatedAt"),
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
			{
				IndexName: aws.String(OperatorDispersalIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("OperatorID"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("DispersedAt"),
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
			{
				IndexName: aws.String(OperatorResponseIndexName),
				KeySchema: []types.KeySchemaElement{
					{
						AttributeName: aws.String("OperatorID"),
						KeyType:       types.KeyTypeHash,
					},
					{
						AttributeName: aws.String("RespondedAt"),
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

func MarshalBlobMetadata(metadata *v2.BlobMetadata) (commondynamodb.Item, error) {
	fields, err := attributevalue.MarshalMap(metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal blob metadata: %w", err)
	}

	// Add PK and SK fields
	blobKey, err := metadata.BlobHeader.BlobKey()
	if err != nil {
		return nil, err
	}
	fields["PK"] = &types.AttributeValueMemberS{Value: blobKeyPrefix + blobKey.Hex()}
	fields["SK"] = &types.AttributeValueMemberS{Value: blobMetadataSK}

	return fields, nil
}

func UnmarshalBlobKey(item commondynamodb.Item) (corev2.BlobKey, error) {
	type Blob struct {
		PK string
	}

	blob := Blob{}
	err := attributevalue.UnmarshalMap(item, &blob)
	if err != nil {
		return corev2.BlobKey{}, err
	}

	bk := strings.TrimPrefix(blob.PK, blobKeyPrefix)
	return corev2.HexToBlobKey(bk)
}

func UnmarshalBlobMetadata(item commondynamodb.Item) (*v2.BlobMetadata, error) {
	metadata := v2.BlobMetadata{}
	err := attributevalue.UnmarshalMap(item, &metadata)
	if err != nil {
		return nil, err
	}
	return &metadata, nil
}

func MarshalBlobCertificate(blobCert *corev2.BlobCertificate, fragmentInfo *encoding.FragmentInfo) (commondynamodb.Item, error) {
	fields, err := attributevalue.MarshalMap(blobCert)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal blob certificate: %w", err)
	}

	// merge fragment info
	fragmentInfoFields, err := attributevalue.MarshalMap(fragmentInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal fragment info: %w", err)
	}
	for k, v := range fragmentInfoFields {
		fields[k] = v
	}

	// Add PK and SK fields
	blobKey, err := blobCert.BlobHeader.BlobKey()
	if err != nil {
		return nil, err
	}
	fields["PK"] = &types.AttributeValueMemberS{Value: blobKeyPrefix + blobKey.Hex()}
	fields["SK"] = &types.AttributeValueMemberS{Value: blobCertSK}

	return fields, nil
}

func UnmarshalBlobCertificate(item commondynamodb.Item) (*corev2.BlobCertificate, *encoding.FragmentInfo, error) {
	cert := corev2.BlobCertificate{}
	err := attributevalue.UnmarshalMap(item, &cert)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal blob certificate: %w", err)
	}
	fragmentInfo := encoding.FragmentInfo{}
	err = attributevalue.UnmarshalMap(item, &fragmentInfo)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal fragment info: %w", err)
	}
	return &cert, &fragmentInfo, nil
}

func UnmarshalBatchHeaderHash(item commondynamodb.Item) ([32]byte, error) {
	type Object struct {
		PK string
	}

	obj := Object{}
	err := attributevalue.UnmarshalMap(item, &obj)
	if err != nil {
		return [32]byte{}, err
	}

	root := strings.TrimPrefix(obj.PK, dispersalKeyPrefix)
	return hexToHash(root)
}

func UnmarshalOperatorID(item commondynamodb.Item) (*core.OperatorID, error) {
	type Object struct {
		OperatorID string
	}

	obj := Object{}
	err := attributevalue.UnmarshalMap(item, &obj)
	if err != nil {
		return nil, err
	}

	operatorID, err := core.OperatorIDFromHex(obj.OperatorID)
	if err != nil {
		return nil, err
	}

	return &operatorID, nil
}

func MarshalDispersalRequest(req *corev2.DispersalRequest) (commondynamodb.Item, error) {
	fields, err := attributevalue.MarshalMap(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dispersal request: %w", err)
	}

	batchHeaderHash, err := req.BatchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}
	hashstr := hex.EncodeToString(batchHeaderHash[:])

	fields["PK"] = &types.AttributeValueMemberS{Value: dispersalKeyPrefix + hashstr}
	fields["SK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", dispersalRequestSKPrefix, req.OperatorID.Hex())}
	fields["OperatorID"] = &types.AttributeValueMemberS{Value: req.OperatorID.Hex()}

	return fields, nil
}

func UnmarshalDispersalRequest(item commondynamodb.Item) (*corev2.DispersalRequest, error) {
	req := corev2.DispersalRequest{}
	err := attributevalue.UnmarshalMap(item, &req)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dispersal request: %w", err)
	}

	operatorID, err := UnmarshalOperatorID(item)
	if err != nil {
		return nil, err
	}
	req.OperatorID = *operatorID

	return &req, nil
}

func MarshalDispersalResponse(res *corev2.DispersalResponse) (commondynamodb.Item, error) {
	fields, err := attributevalue.MarshalMap(res)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal dispersal response: %w", err)
	}

	batchHeaderHash, err := res.BatchHeader.Hash()
	if err != nil {
		return nil, fmt.Errorf("failed to hash batch header: %w", err)
	}
	hashstr := hex.EncodeToString(batchHeaderHash[:])

	fields["PK"] = &types.AttributeValueMemberS{Value: dispersalKeyPrefix + hashstr}
	fields["SK"] = &types.AttributeValueMemberS{Value: fmt.Sprintf("%s%s", dispersalResponseSKPrefix, res.OperatorID.Hex())}
	fields["OperatorID"] = &types.AttributeValueMemberS{Value: res.OperatorID.Hex()}

	return fields, nil
}

func UnmarshalDispersalResponse(item commondynamodb.Item) (*corev2.DispersalResponse, error) {
	res := corev2.DispersalResponse{}
	err := attributevalue.UnmarshalMap(item, &res)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dispersal response: %w", err)
	}

	operatorID, err := UnmarshalOperatorID(item)
	if err != nil {
		return nil, err
	}
	res.OperatorID = *operatorID

	return &res, nil
}

func hexToHash(h string) ([32]byte, error) {
	s := strings.TrimPrefix(h, "0x")
	s = strings.TrimPrefix(s, "0X")
	b, err := hex.DecodeString(s)
	if err != nil {
		return [32]byte{}, err
	}
	return [32]byte(b), nil
}
