package signingrate

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"google.golang.org/protobuf/proto"
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

	err := s.buildEndTimestampIndex(ctx)
	if err != nil {
		return nil, fmt.Errorf("error ensuring EndTimestamp index exists: %w", err)
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
		TableName:        d.tableName,
		Key:              key,
		UpdateExpression: aws.String("SET Bucket = :bucket, EndTimestamp = :endTimestamp, BucketType = :bucketType"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":bucket":       &types.AttributeValueMemberB{Value: value},
			":endTimestamp": &types.AttributeValueMemberS{Value: timestampToString(bucket.EndTimestamp())},
			":bucketType":   &types.AttributeValueMemberS{Value: "Bucket"},
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
		"StartTimestamp": &types.AttributeValueMemberS{Value: timestampToString(timestamp)},
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
		TableName:              d.tableName,
		IndexName:              aws.String("EndTimestampIndex"),
		KeyConditionExpression: aws.String("BucketType = :bt AND EndTimestamp > :start"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":bt":    &types.AttributeValueMemberS{Value: "Bucket"},
			":start": &types.AttributeValueMemberS{Value: timestampToString(startTimestamp)},
		},
		ProjectionExpression: aws.String("Bucket"),
	}

	var out []*SigningRateBucket
	for {
		resp, err := d.dynamoClient.Query(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("dynamo query failed: %w", err)
		}

		for _, item := range resp.Items {
			bin, ok := item["Bucket"].(*types.AttributeValueMemberB)
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

// If it doesn't yet exist, ensure that dynamo table has an index on EndTimestamp.
// This index is needed to efficiently query for all buckets after a certain time.
func (d *dynamoSigningRateStorage) buildEndTimestampIndex(ctx context.Context) error {
	td, err := d.dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: d.tableName,
	})
	if err != nil {
		return fmt.Errorf("describe table: %w", err)
	}

	// Already present?
	for _, g := range td.Table.GlobalSecondaryIndexes {
		if aws.ToString(g.IndexName) == "EndTimestampIndex" {
			if g.IndexStatus == types.IndexStatusActive {
				return nil
			}
			// wait until it's ACTIVE
			err = d.waitForIndexActive(ctx, "EndTimestampIndex")
			if err != nil {
				return fmt.Errorf("wait for index active: %w", err)
			}
			return nil
		}
	}

	// Build the create request
	create := types.CreateGlobalSecondaryIndexAction{
		IndexName: aws.String("EndTimestampIndex"),
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("BucketType"), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String("EndTimestamp"), KeyType: types.KeyTypeRange},
		},
		Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
	}

	// Only needed if the table uses provisioned throughput (not on-demand)
	isProvisioned := td.Table.BillingModeSummary == nil ||
		td.Table.BillingModeSummary.BillingMode == types.BillingModeProvisioned
	if isProvisioned {
		create.ProvisionedThroughput = &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5),
			WriteCapacityUnits: aws.Int64(5),
		}
	}

	_, err = d.dynamoClient.UpdateTable(ctx, &dynamodb.UpdateTableInput{
		TableName: d.tableName,
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("BucketType"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("EndTimestamp"), AttributeType: types.ScalarAttributeTypeS},
		},
		GlobalSecondaryIndexUpdates: []types.GlobalSecondaryIndexUpdate{
			{Create: &create},
		},
	})
	if err != nil {
		return fmt.Errorf("create GSI: %w", err)
	}

	err = d.waitForIndexActive(ctx, "EndTimestampIndex")
	if err != nil {
		return fmt.Errorf("wait for index active: %w", err)
	}

	return nil
}

func (d *dynamoSigningRateStorage) waitForIndexActive(ctx context.Context, indexName string) error {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	timeout := time.After(5 * time.Minute)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for index %s to become ACTIVE", indexName)
		case <-ticker.C:
			td, err := d.dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: d.tableName,
			})
			if err != nil {
				return fmt.Errorf("describe table while waiting: %w", err)
			}
			for _, g := range td.Table.GlobalSecondaryIndexes {
				if aws.ToString(g.IndexName) == indexName && g.IndexStatus == types.IndexStatusActive {
					return nil
				}
			}
		}
	}
}
