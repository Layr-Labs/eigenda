package signingrate

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/testbed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/require"
)

// setupDynamoClient sets up a DynamoDB client connected to Localstack for testing.
func setupDynamoClient(t *testing.T) (client *dynamodb.Client, cleanup func()) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	localstackPort := 4573

	var localstackContainer *testbed.LocalStackContainer
	var deployLocalStack bool

	ctx := t.Context()

	if os.Getenv("DEPLOY_LOCALSTACK") != "false" {
		deployLocalStack = true
		localstackContainer, err = testbed.NewLocalStackContainerWithOptions(ctx, testbed.LocalStackOptions{
			ExposeHostPort: true,
			HostPort:       fmt.Sprintf("%d", localstackPort),
			Services:       []string{"dynamodb"},
			Logger:         logger,
		})
		require.NoError(t, err)
	} else {
		// localstack is already deployed
		portString := os.Getenv("LOCALSTACK_PORT")
		require.NoError(t, err)
		localstackPort, err = strconv.Atoi(portString)
		require.NoError(t, err)
	}

	clientConfig := commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%d", localstackPort),
	}

	awsConfig := aws.Config{
		Region: clientConfig.Region,
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     clientConfig.AccessKey,
				SecretAccessKey: clientConfig.SecretAccessKey,
			}, nil
		}),
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if clientConfig.EndpointURL != "" {
					return aws.Endpoint{
						PartitionID:   "aws",
						URL:           clientConfig.EndpointURL,
						SigningRegion: clientConfig.Region,
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			}),
	}
	client = dynamodb.NewFromConfig(awsConfig)

	return client, func() {
		if deployLocalStack {
			_ = localstackContainer.Terminate(ctx)
		}
	}
}

// simulateRandomSigningRateActivity simulates random signing activity on the given tracker. Does not attempt to
// advance time.
func simulateRandomSigningRateActivity(
	rand *random.TestRandom,
	tracker SigningRateTracker,
	quorumCount core.QuorumID,
	validatorIDs []core.OperatorID,
	iterations int,
) {

	for i := 0; i < iterations; i++ {
		quorum := core.QuorumID(rand.Intn(int(quorumCount)))
		validator := validatorIDs[rand.Intn(len(validatorIDs))]
		batchSize := uint64(rand.Intn(10) + 1)
		signingLatency := time.Duration(rand.Intn(1000)) * time.Millisecond

		if rand.Bool() {
			tracker.ReportSuccess(quorum, validator, batchSize, signingLatency)
		} else {
			tracker.ReportFailure(quorum, validator, batchSize)
		}
	}
}

// Compare two ValidatorSigningRate slices for equality.
func areValidatorSigningRatesEqual(
	expected []*validator.SigningRateBucket,
	actual []*validator.SigningRateBucket,
) bool {

	if len(expected) != len(actual) {
		return false
	}

	for i := range expected {
		expectedBucket := expected[i]
		actualBucket := actual[i]

		if expectedBucket.GetStartTimestamp() != actualBucket.GetStartTimestamp() ||
			expectedBucket.GetEndTimestamp() != actualBucket.GetEndTimestamp() {
			return false
		}

		if len(expectedBucket.GetQuorumSigningRates()) != len(actualBucket.GetQuorumSigningRates()) {
			return false
		}

		for j := range expectedBucket.GetQuorumSigningRates() {
			expectedQuorum := expectedBucket.GetQuorumSigningRates()[j]
			actualQuorum := actualBucket.GetQuorumSigningRates()[j]

			if expectedQuorum.GetQuorumId() != actualQuorum.GetQuorumId() {
				return false
			}

			if len(expectedQuorum.GetValidatorSigningRates()) != len(actualQuorum.GetValidatorSigningRates()) {
				return false
			}

			for k := range expectedQuorum.GetValidatorSigningRates() {
				expectedValidator := expectedQuorum.GetValidatorSigningRates()[k]
				actualValidator := actualQuorum.GetValidatorSigningRates()[k]

				if !areSigningRatesEqual(expectedValidator, actualValidator) {
					return false
				}
			}
		}
	}

	return true
}

// Setting up local stack is slow. Cram a bunch of test cases into this one test to avoid this cost.
func TestSigningRateStorage(t *testing.T) {
	t.Parallel()

	rand := random.NewTestRandom()

	dynamoClient, cleanup := setupDynamoClient(t)
	defer cleanup()

	tableName := "TestSigningRateStorage"
	storage, err := NewDynamoSigningRateStorage(t.Context(), dynamoClient, tableName)
	require.NoError(t, err)

	validatorCount := rand.Intn(10) + 5
	validatorIDs := make([]core.OperatorID, validatorCount)
	for i := 0; i < validatorCount; i++ {
		validatorIDs[i] = core.OperatorID(rand.Bytes(32))
	}

	quorumCount := core.QuorumID(rand.Intn(5) + 3)

	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	now := rand.Time()
	timePointer := atomic.Pointer[time.Time]{}
	timePointer.Store(&now)
	timeSource := func() time.Time {
		return *timePointer.Load()
	}

	// Use a signing rate tracker as a "source of truth". This data structure is validated in its own unit tests,
	// so trust it here.
	tracker, err := NewSigningRateTracker(logger, time.Hour*100, time.Minute*10, timeSource)
	require.NoError(t, err)

	// Check query behavior when there are no buckets.
	buckets, err := storage.LoadBuckets(t.Context(), time.Unix(0, 0))
	require.NoError(t, err)
	require.Len(t, buckets, 0)

	// Add a single bucket and check it can be retrieved.
	simulateRandomSigningRateActivity(rand, tracker, quorumCount, validatorIDs, 100)

	unflushedBuckets, err := tracker.GetUnflushedBuckets()
	require.NoError(t, err)
	require.Len(t, unflushedBuckets, 1)
	err = storage.StoreBuckets(t.Context(), unflushedBuckets)
	require.NoError(t, err)

	expectedBuckets, err := tracker.GetSigningRateDump(time.Unix(0, 0))
	require.NoError(t, err)
	require.Len(t, expectedBuckets, 1)
	actualBuckets, err := storage.LoadBuckets(t.Context(), time.Unix(0, 0))
	require.NoError(t, err)
	require.True(t, areValidatorSigningRatesEqual(expectedBuckets, actualBuckets))

	// Add several more buckets.
	for i := 0; i < 5; i++ {
		now = now.Add(time.Minute * 10)
		timePointer.Store(&now)
		simulateRandomSigningRateActivity(rand, tracker, quorumCount, validatorIDs, 100)
	}

	unflushedBuckets, err = tracker.GetUnflushedBuckets()
	require.NoError(t, err)
	require.Len(t, unflushedBuckets, 5)
	err = storage.StoreBuckets(t.Context(), unflushedBuckets)
	require.NoError(t, err)

	expectedBuckets, err = tracker.GetSigningRateDump(time.Unix(0, 0))
	require.NoError(t, err)
	require.Len(t, expectedBuckets, 6)
	actualBuckets, err = storage.LoadBuckets(t.Context(), time.Unix(0, 0))
	require.NoError(t, err)
	require.True(t, areValidatorSigningRatesEqual(expectedBuckets, actualBuckets))

	// Query for a subset of the data.

	// Fetch data starting exactly at the start of a bucket.
	targetIndex := len(expectedBuckets) / 2
	startTimestamp := expectedBuckets[targetIndex].GetStartTimestamp()

	actualBuckets, err = storage.LoadBuckets(t.Context(), time.Unix(int64(startTimestamp), 0))
	require.NoError(t, err)
	require.True(t, areValidatorSigningRatesEqual(expectedBuckets[targetIndex:], actualBuckets))

	// If we subtract one second from the starting timestamp, we should snag the previous bucket as well.
	actualBuckets, err = storage.LoadBuckets(t.Context(), time.Unix(int64(startTimestamp)-1, 0))
	require.NoError(t, err)
	require.True(t, areValidatorSigningRatesEqual(expectedBuckets[targetIndex-1:], actualBuckets))

	// Modify the last bucket and ensure it gets overwritten correctly.
	// Note that we are not advancing time, so this activity goes into the last bucket.
	simulateRandomSigningRateActivity(rand, tracker, quorumCount, validatorIDs, 100)

	unflushedBuckets, err = tracker.GetUnflushedBuckets()
	require.NoError(t, err)
	require.Len(t, unflushedBuckets, 1)
	err = storage.StoreBuckets(t.Context(), unflushedBuckets)
	require.NoError(t, err)

	expectedBuckets, err = tracker.GetSigningRateDump(time.Unix(0, 0))
	require.NoError(t, err)
	require.Len(t, expectedBuckets, 6)
	actualBuckets, err = storage.LoadBuckets(t.Context(), time.Unix(0, 0))
	require.NoError(t, err)

	require.True(t, areValidatorSigningRatesEqual(expectedBuckets, actualBuckets))
}
