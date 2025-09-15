package signingrate

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/proto"
)

const (
	// DynamoDB attribute names
	attrStartTimestamp = "StartTimestamp"
	attrBucketType     = "BucketType"
	attrEndTimestamp   = "EndTimestamp"
	attrBucket         = "Bucket"

	endTimestampIndex = "EndTimestampIndex"

	// DynamoDB expression placeholders
	placeholderBucket       = ":bucket"
	placeholderEndTimestamp = ":endTimestamp"
	placeholderBucketType   = ":bucketType"
	placeholderBT           = ":bt"
	placeholderStart        = ":start"
)

var _ SigningRateStorage = (*dynamoSigningRateStorage)(nil)

// A DynamoDB implementation of the SigningRateStorage interface.
type dynamoSigningRateStorage struct {
	dynamoClient *dynamodb.Client
	tableName    *string
}

// Create a new DynamoDB-backed SigningRateStorage.
func NewDynamoSigningRateStorage(
	ctx context.Context,
	dynamoClient *dynamodb.Client,
	tableName string,
) (SigningRateStorage, error) {

	if dynamoClient == nil {
		return nil, fmt.Errorf("dynamoClient cannot be nil")
	}
	if tableName == "" {
		return nil, fmt.Errorf("tableName cannot be empty")
	}

	s := &dynamoSigningRateStorage{
		dynamoClient: dynamoClient,
		tableName:    aws.String(tableName),
	}

	err := s.ensureTableExists(ctx)
	if err != nil {
		return nil, fmt.Errorf("error ensuring table exists: %w", err)
	}

	return s, nil
}

func (d *dynamoSigningRateStorage) StoreBuckets(ctx context.Context, buckets []*SigningRateBucket) error {
	for _, bucket := range buckets {
		if err := d.storeBucket(ctx, bucket); err != nil {
			return fmt.Errorf("error storing bucket: %w", err)
		}
	}
	return nil
}

func (d *dynamoSigningRateStorage) storeBucket(ctx context.Context, bucket *SigningRateBucket) error {

	key := getDynamoBucketKey(bucket)
	value, err := proto.Marshal(bucket.ToProtobuf())
	if err != nil {
		return fmt.Errorf("proto marshal failed: %w", err)
	}

	// Note: the "BucketType" attribute is due a quirk in dynamo. It won't let us do certain queries unless we have
	// a partition key. So we create a dummy partition key that is always the same value.

	_, err = d.dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: d.tableName,
		Key:       key,
		UpdateExpression: aws.String(fmt.Sprintf("SET %s = %s, %s = %s, %s = %s",
			attrBucket, placeholderBucket,
			attrEndTimestamp, placeholderEndTimestamp,
			attrBucketType, placeholderBucketType)),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			placeholderBucket:       &types.AttributeValueMemberB{Value: value},
			placeholderEndTimestamp: &types.AttributeValueMemberS{Value: timestampToString(bucket.EndTimestamp())},
			placeholderBucketType:   &types.AttributeValueMemberS{Value: attrBucket},
		},
	})

	if err != nil {
		return fmt.Errorf("dynamo update failed: %w", err)
	}
	return nil
}

// Get the DynamoDB key for a given bucket. The primary key for a bucket is its starting timestamp.
func getDynamoBucketKey(bucket *SigningRateBucket) map[string]types.AttributeValue {
	timestamp := bucket.StartTimestamp()

	return map[string]types.AttributeValue{
		attrStartTimestamp: &types.AttributeValueMemberS{Value: timestampToString(timestamp)},
	}
}

// Convert a timestamp to the string format used in DynamoDB. String is padded with zeros on the left to ensure
// lexicographical ordering based on string comparison. This method assumes that timestamps are non-negative and
// represent seconds since the Unix epoch (i.e. sub-second precision is not supported).
func timestampToString(t time.Time) string {
	// 20 digits can hold a maximally sized uint64, so ensure that's how much we always use.
	return fmt.Sprintf("%020d", t.Unix())
}

func (d *dynamoSigningRateStorage) LoadBuckets(
	ctx context.Context,
	startTimestamp time.Time,
) ([]*SigningRateBucket, error) {

	input := &dynamodb.QueryInput{
		TableName: d.tableName,
		IndexName: aws.String(endTimestampIndex),
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = %s AND %s > %s",
			attrBucketType, placeholderBT,
			attrEndTimestamp, placeholderStart)),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			placeholderBT:    &types.AttributeValueMemberS{Value: attrBucket},
			placeholderStart: &types.AttributeValueMemberS{Value: timestampToString(startTimestamp)},
		},
		ProjectionExpression: aws.String(attrBucket),
	}

	var out []*SigningRateBucket
	for {
		resp, err := d.dynamoClient.Query(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("dynamo query failed: %w", err)
		}

		for _, item := range resp.Items {
			bin, ok := item[attrBucket].(*types.AttributeValueMemberB)
			if !ok {
				// Row missing payload; skip
				continue
			}

			pb := &validator.SigningRateBucket{}
			if err := proto.Unmarshal(bin.Value, pb); err != nil {
				return nil, fmt.Errorf("unmarshal bucket proto: %w", err)
			}

			b := NewBucketFromProto(pb)
			out = append(out, b)
		}

		if len(resp.LastEvaluatedKey) == 0 {
			break
		}
		input.ExclusiveStartKey = resp.LastEvaluatedKey
	}

	// Index returns rows ordered by EndTimestamp which may not be unique. Sort by StartTimestamp, which are unique.
	sort.Slice(out, func(i, j int) bool {
		return out[i].StartTimestamp().Before(out[j].StartTimestamp())
	})

	return out, nil
}

func (d *dynamoSigningRateStorage) ensureTableExists(ctx context.Context) error {
	_, err := d.dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: d.tableName,
	})
	if err == nil {
		// Table exists, wait until ACTIVE
		return d.waitForTableActive(ctx)
	}

	var rnfe *types.ResourceNotFoundException
	if !errors.As(err, &rnfe) {
		return fmt.Errorf("describe table: %w", err)
	}

	_, err = d.dynamoClient.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: d.tableName,
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String(attrStartTimestamp), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String(attrBucketType), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String(attrEndTimestamp), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String(attrStartTimestamp), KeyType: types.KeyTypeHash},
		},
		BillingMode: types.BillingModePayPerRequest,
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(endTimestampIndex),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String(attrBucketType), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String(attrEndTimestamp), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
				// No ProvisionedThroughput because we're PAY_PER_REQUEST
			},
		},
	})
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	// Wait for table ACTIVE
	return d.waitForTableActive(ctx)
}

func (d *dynamoSigningRateStorage) waitForTableActive(ctx context.Context) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	timeout := time.After(10 * time.Minute)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for table to become ACTIVE")
		case <-ticker.C:
			out, err := d.dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: d.tableName,
			})
			if err != nil {
				return fmt.Errorf("describe table while waiting: %w", err)
			}
			if out.Table != nil && out.Table.TableStatus == types.TableStatusActive {
				// Also verify the GSI is ACTIVE (created at table creation)
				ok := true
				for _, g := range out.Table.GlobalSecondaryIndexes {
					if g.IndexName != nil && *g.IndexName == endTimestampIndex &&
						g.IndexStatus != types.IndexStatusActive {
						ok = false
						break
					}
				}
				if ok {
					return nil
				}
			}
		}
	}
}
