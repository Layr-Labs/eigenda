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

// ═══════════════════════════════════════════════════════════════════════════════════════════
// DynamoDB Storage Structure Documentation
// ═══════════════════════════════════════════════════════════════════════════════════════════
//
// ## What We're Storing
//
// This storage layer persists signing rate buckets (SigningRateBucket objects) to DynamoDB.
// Each bucket represents a time window containing signing rate data. We need to:
//   1. Store new buckets as they're created
//   2. Retrieve all buckets that ended after a specific time (for loading historical data)
//
// ## DynamoDB Basics
//
// DynamoDB is a NoSQL key-value database. Unlike SQL databases with tables and flexible queries,
// DynamoDB has strict requirements about how you access data:
//
// **Primary Key (Partition Key)**
//   - Every table MUST have a primary key that uniquely identifies each item (row)
//   - You can retrieve items directly by their primary key (very fast, single-digit millisecond)
//   - You CANNOT query by other attributes without creating indexes
//
// **Global Secondary Index (GSI)**
//   - A GSI is an alternative "view" of your table with a different key structure
//   - Lets you query the table using different attributes than the primary key
//   - GSIs MUST have a partition key, and optionally a sort key for range queries
//   - GSIs duplicate your data (managed automatically by DynamoDB)
//
// **Important Constraint**: All DynamoDB queries MUST specify a partition key value.
//   - You cannot do a "scan all items where X > Y" without a partition key
//   - This is a fundamental DynamoDB limitation/design choice for "performance"
//   - Since we don't have a natural partition key for our query pattern, this code
//     uses a hacky workaround (explained below).
//
// ## Our Table Structure
//
// **Main Table:**
//   - Primary Key: StartTimestamp (when the bucket started)
//     * This makes sense because each bucket has a unique start time
//     * Allows us to store/retrieve specific buckets efficiently
//   - Other Attributes:
//     * EndTimestamp: When the bucket ended
//     * Payload: The serialized protobuf data (the actual bucket contents)
//     * PayloadType: A dummy constant value (used as a dummy partition key, explained below)
//
// **Global Secondary Index (EndTimestampIndex):**
//   - Partition Key: PayloadType (always set to "Payload" - a constant dummy value)
//   - Sort Key: EndTimestamp (allows range queries like "EndTimestamp > X")
//
// ## Why This Design?
//
// **Problem**: We need to query "all buckets where EndTimestamp > X" to load historical data.
//   - We can't use the main table because its key is StartTimestamp
//   - DynamoDB won't let us query by EndTimestamp without an index
//
// **Solution**: Create a Global Secondary Index with EndTimestamp as the sort key.
//   - But GSIs require a partition key (DynamoDB rule)
//   - We don't have a natural partition key for this query pattern
//
// **The Dummy Partition Key Trick**:
//   - We create an artificial attribute called PayloadType that's always "Payload"
//   - Every item gets the same PayloadType value
//   - This puts all items in the same partition for the GSI
//   - Now we can query: "PayloadType = 'Payload' AND EndTimestamp > X"
//
// **Why Zero-Pad Timestamps?**
//   - DynamoDB sorts strings lexicographically (like dictionary order)
//   - "9" > "10" in string comparison, but 9 < 10 numerically
//   - Zero-padding ensures string sort order matches numerical order
//   - "0009" < "0010" (correct!)
//   - We pad to 20 digits to handle the full uint64 range
//
// ## Example Query Flow
//
// To load all buckets ending after time T:
//   1. Query the EndTimestampIndex GSI
//   2. Condition: PayloadType = "Payload" AND EndTimestamp > timestampToString(T)
//   3. DynamoDB returns matching items ordered by EndTimestamp
//   4. We deserialize the Payload attribute to get the bucket objects
//   5. Sort by StartTimestamp for deterministic ordering (EndTimestamp may not be unique)
//
// ## Data Format
//
// Each item in the table looks like:
//   {
//     "StartTimestamp": "00000000001234567890",  // Primary key (zero-padded)
//     "EndTimestamp":   "00000000001234568890",  // Used for range queries (zero-padded)
//     "PayloadType":    "Payload",               // Dummy partition key (always "Payload")
//     "Payload":        <binary protobuf data>   // Serialized SigningRateBucket
//   }
//
// ═══════════════════════════════════════════════════════════════════════════════════════════

const (
	// DynamoDB attribute names - these define the column names in our table
	attrStartTimestamp = "StartTimestamp" // Primary key: when the bucket started (unique identifier)
	attrPayloadType    = "PayloadType"    // Artificial partition key for Global Secondary Index queries (always "Payload")
	attrEndTimestamp   = "EndTimestamp"   // When the bucket ended (used for range queries)
	attrPayload        = "Payload"        // The serialized protobuf data

	// Global Secondary Index name - allows us to query by EndTimestamp ranges
	// DynamoDB requires a partition key for all queries, so we use PayloadType as a dummy partition
	endTimestampIndex = "EndTimestampIndex"
	payloadTypeValue  = "Payload" // Constant value for the dummy partition key

	// DynamoDB expression placeholders - these are security tokens that prevent injection attacks
	// DynamoDB requires all values in expressions to be parameterized using these placeholders
	placeholderPayload      = ":payload"
	placeholderEndTimestamp = ":endTimestamp"
	placeholderPayloadType  = ":payloadType"
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

func (d *dynamoSigningRateStorage) StoreBuckets(ctx context.Context, buckets []*validator.SigningRateBucket) error {
	for _, bucket := range buckets {
		if err := d.storeBucket(ctx, bucket); err != nil {
			return fmt.Errorf("error storing bucket: %w", err)
		}
	}
	return nil
}

func (d *dynamoSigningRateStorage) storeBucket(ctx context.Context, bucket *validator.SigningRateBucket) error {

	// Create the primary key for this bucket (StartTimestamp)
	key := getDynamoBucketKey(bucket)

	// Serialize the bucket data as protobuf bytes for storage
	value, err := proto.Marshal(bucket)
	if err != nil {
		return fmt.Errorf("proto marshal failed: %w", err)
	}

	// Note: PayloadType is a dummy partition key required for Global Secondary Index queries.
	// DynamoDB requires all queries to specify a partition key, but we want to query by EndTimestamp ranges.
	// So we create an artificial partition key that's always the same value ("Payload").

	// Use UpdateItem instead of PutItem because it creates the item if it doesn't exist,
	// or updates it if it does exist (upsert behavior).
	_, err = d.dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: d.tableName,
		Key:       key, // Primary key: {StartTimestamp: "00000001234567890"}

		// SET expression updates/creates the specified attributes
		// This is DynamoDB's expression language for atomic updates
		UpdateExpression: aws.String(fmt.Sprintf("SET %s = %s, %s = %s, %s = %s",
			attrPayload, placeholderPayload,           // Store serialized bucket data
			attrEndTimestamp, placeholderEndTimestamp, // Store end timestamp for range queries
			attrPayloadType, placeholderPayloadType)), // Store dummy partition key for Global Secondary Index

		// Map placeholder tokens to actual values - this prevents injection attacks
		// and allows DynamoDB to optimize the query
		ExpressionAttributeValues: map[string]types.AttributeValue{
			// Binary data type for protobuf bytes
			placeholderPayload: &types.AttributeValueMemberB{Value: value},
			// String data type for timestamp (zero-padded for lexicographical sorting)
			placeholderEndTimestamp: &types.AttributeValueMemberS{Value: timestampToString(bucket.GetEndTimestamp())},
			// String data type for dummy partition key
			placeholderPayloadType: &types.AttributeValueMemberS{Value: payloadTypeValue},
		},
	})

	if err != nil {
		return fmt.Errorf("dynamo update failed: %w", err)
	}
	return nil
}

// Get the DynamoDB key for a given bucket. The primary key for a bucket is its starting timestamp.
// getDynamoBucketKey creates the primary key for a bucket in DynamoDB.
// DynamoDB keys must be a map of attribute names to AttributeValue objects.
// We use StartTimestamp as the primary key because it's unique for each bucket.
func getDynamoBucketKey(bucket *validator.SigningRateBucket) map[string]types.AttributeValue {
	timestamp := bucket.GetStartTimestamp()

	// Return a composite key map - in this case just the primary key
	// AttributeValueMemberS indicates this is a String type in DynamoDB
	return map[string]types.AttributeValue{
		attrStartTimestamp: &types.AttributeValueMemberS{Value: timestampToString(timestamp)},
	}
}

// Convert a timestamp to the string format used in DynamoDB. String is padded with zeros on the left to ensure
// lexicographical ordering based on string comparison. This method assumes that timestamps are non-negative and
// represent seconds since the Unix epoch (i.e. sub-second precision is not supported).
// timestampToString converts a Unix timestamp to a zero-padded string for DynamoDB storage.
// DynamoDB stores everything as strings/numbers/binary, and string comparison is lexicographical.
// By zero-padding to 20 digits, we ensure that string sorting matches numerical sorting:
// "00000000000000000001" < "00000000000000000010" (correct)
// vs "1" > "10" (incorrect if not padded)
// 20 digits can hold the maximum uint64 value (18,446,744,073,709,551,615).
func timestampToString(unixTime uint64) string {
	return fmt.Sprintf("%020d", unixTime)
}

// LoadBuckets retrieves all signing rate buckets that ended after the given start time.
// This method demonstrates DynamoDB's pagination and Global Secondary Index usage.
func (d *dynamoSigningRateStorage) LoadBuckets(
	ctx context.Context,
	startTimestamp time.Time,
) ([]*validator.SigningRateBucket, error) {

	// Query the Global Secondary Index instead of the main table
	// Global Secondary Index allows us to query by EndTimestamp ranges, which isn't possible with the main table
	// that only has StartTimestamp as the key
	input := &dynamodb.QueryInput{
		TableName: d.tableName,
		// Use the Global Secondary Index that has PayloadType as partition key and EndTimestamp as sort key
		IndexName: aws.String(endTimestampIndex),

		// KeyConditionExpression defines the query conditions
		// Format: "partition_key = value AND sort_key > value"
		// We must specify the partition key (PayloadType) and can add range conditions on sort key (EndTimestamp)
		KeyConditionExpression: aws.String(fmt.Sprintf("%s = %s AND %s > %s",
			attrPayloadType, placeholderPayloadType,   // PayloadType = "Payload" (dummy partition)
			attrEndTimestamp, placeholderStart)),       // EndTimestamp > startTimestamp

		// Parameterized values for the query conditions (prevents injection attacks)
		ExpressionAttributeValues: map[string]types.AttributeValue{
			// All items have this same dummy partition key value
			placeholderPayloadType: &types.AttributeValueMemberS{Value: payloadTypeValue},
			// Convert the start time to our zero-padded string format for comparison
			placeholderStart: &types.AttributeValueMemberS{
				Value: timestampToString(uint64(startTimestamp.Unix())),
			},
		},
		// ProjectionExpression limits which attributes are returned (saves bandwidth/cost)
		// We only need the Payload attribute since that contains the full serialized bucket
		ProjectionExpression: aws.String(attrPayload),
	}

	var out []*validator.SigningRateBucket

	// DynamoDB paginates results automatically. We need to loop to get all pages.
	// Each Query call returns at most 1MB of data or 1000 items, whichever comes first.
	for {
		resp, err := d.dynamoClient.Query(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("dynamo query failed: %w", err)
		}

		// Process each item in this page of results
		for _, item := range resp.Items {
			// Extract the binary payload attribute and verify it's the right type
			// DynamoDB returns a generic AttributeValue interface, we need to cast to the specific type
			bin, ok := item[attrPayload].(*types.AttributeValueMemberB)
			if !ok {
				// This shouldn't happen unless the data is corrupted, but skip gracefully
				continue
			}

			// Deserialize the protobuf data back into a bucket object
			pb := &validator.SigningRateBucket{}
			if err := proto.Unmarshal(bin.Value, pb); err != nil {
				return nil, fmt.Errorf("unmarshal bucket proto: %w", err)
			}

			out = append(out, pb)
		}

		// Check if there are more pages to fetch
		// LastEvaluatedKey contains the primary key of the last item processed
		// If it's empty, we've reached the end of the results
		if len(resp.LastEvaluatedKey) == 0 {
			break
		}

		// Set the starting point for the next page of results
		// DynamoDB will continue from after this key
		input.ExclusiveStartKey = resp.LastEvaluatedKey
	}

	// The Global Secondary Index returns rows ordered by EndTimestamp, but EndTimestamp values may not be unique.
	// Sort by StartTimestamp to ensure deterministic ordering, since StartTimestamp is unique.
	sort.Slice(out, func(i, j int) bool {
		return out[i].GetStartTimestamp() < out[j].GetStartTimestamp()
	})

	return out, nil
}

// ensureTableExists checks if the DynamoDB table exists and creates it if necessary.
// This method demonstrates DynamoDB table creation with Global Secondary Indexes.
func (d *dynamoSigningRateStorage) ensureTableExists(ctx context.Context) error {
	// First, try to describe the table to see if it exists
	_, err := d.dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: d.tableName,
	})
	if err == nil {
		// Table exists, but it might still be in CREATING status, so wait until ACTIVE
		return d.waitForTableActive(ctx)
	}

	// Check if the error is specifically "table not found"
	var rnfe *types.ResourceNotFoundException
	if !errors.As(err, &rnfe) {
		// Some other error occurred (permissions, network, etc.)
		return fmt.Errorf("describe table: %w", err)
	}

	// Table doesn't exist, so create it
	_, err = d.dynamoClient.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: d.tableName,

		// AttributeDefinitions specify the data types for attributes used in keys
		// You only need to define attributes that are used in primary keys or index keys
		// Other attributes are schemaless and don't need to be defined here
		AttributeDefinitions: []types.AttributeDefinition{
			// Primary key attribute for the main table
			{AttributeName: aws.String(attrStartTimestamp), AttributeType: types.ScalarAttributeTypeS},
			// Global Secondary Index partition key (dummy key that's always the same value)
			{AttributeName: aws.String(attrPayloadType), AttributeType: types.ScalarAttributeTypeS},
			// Global Secondary Index sort key (allows range queries on EndTimestamp)
			{AttributeName: aws.String(attrEndTimestamp), AttributeType: types.ScalarAttributeTypeS},
		},

		// KeySchema defines the primary key structure for the main table
		KeySchema: []types.KeySchemaElement{
			// HASH key is the partition key - determines which physical partition stores the item
			// We use StartTimestamp as our primary key since each bucket has a unique start time
			{AttributeName: aws.String(attrStartTimestamp), KeyType: types.KeyTypeHash},
			// No RANGE key needed for the main table since StartTimestamp alone is unique
		},

		// Use pay-per-request billing instead of provisioned capacity
		// This automatically scales and we only pay for actual usage
		BillingMode: types.BillingModePayPerRequest,

		// Global Secondary Indexes allow alternative access patterns
		// Global Secondary Indexes have their own key structure and can be queried independently of the main table
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(endTimestampIndex),

				// Global Secondary Index key structure: PayloadType (partition) + EndTimestamp (sort)
				// This allows us to query "all items with PayloadType='Payload' where EndTimestamp > X"
				KeySchema: []types.KeySchemaElement{
					// Partition key for the Global Secondary Index - we use a dummy constant value
					// This puts all items in the same partition, which is fine for our use case
					{AttributeName: aws.String(attrPayloadType), KeyType: types.KeyTypeHash},
					// Sort key for the Global Secondary Index - allows range queries on EndTimestamp
					{AttributeName: aws.String(attrEndTimestamp), KeyType: types.KeyTypeRange},
				},

				// ProjectionType determines what attributes are copied to the Global Secondary Index
				// ALL means all attributes are available, so we can query the Global Secondary Index without
				// needing to look up the main table (avoids additional read costs)
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},

				// No ProvisionedThroughput specified because we're using PAY_PER_REQUEST billing
				// The Global Secondary Index inherits the billing mode from the main table
			},
		},
	})
	if err != nil {
		return fmt.Errorf("create table: %w", err)
	}

	// Table creation is asynchronous - wait for it to become ACTIVE before using it
	return d.waitForTableActive(ctx)
}

// waitForTableActive polls DynamoDB until the table and all its indexes are ready for use.
// DynamoDB table/index creation is asynchronous and can take several minutes.
func (d *dynamoSigningRateStorage) waitForTableActive(ctx context.Context) error {
	// Poll every 2 seconds to check table status
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Give up after 10 minutes - table creation shouldn't take this long
	timeout := time.After(10 * time.Minute)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for table to become ACTIVE")
		case <-ticker.C:
			// Query the table's current status
			out, err := d.dynamoClient.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: d.tableName,
			})
			if err != nil {
				return fmt.Errorf("describe table while waiting: %w", err)
			}

			// Check if the main table is ACTIVE
			// Possible statuses: CREATING, ACTIVE, DELETING, UPDATING
			if out.Table != nil && out.Table.TableStatus == types.TableStatusActive {
				// Table is ACTIVE, but we also need to check that all Global Secondary Indexes are ACTIVE
				// Global Secondary Indexes can have their own status independent of the main table
				ok := true
				for _, globalSecondaryIndex := range out.Table.GlobalSecondaryIndexes {
					// Find our EndTimestampIndex and check its status
					if globalSecondaryIndex.IndexName != nil && *globalSecondaryIndex.IndexName == endTimestampIndex {
						// Global Secondary Index possible statuses: CREATING, ACTIVE, DELETING, UPDATING
						if globalSecondaryIndex.IndexStatus != types.IndexStatusActive {
							ok = false
							break
						}
					}
				}

				// Both table and all Global Secondary Indexes are ACTIVE - ready to use
				if ok {
					return nil
				}
			}
			// If we get here, either table or Global Secondary Index is not ACTIVE yet, continue polling
		}
	}
}
