package ondemand_test

import (
	"fmt"
	"os"
	"testing"

	commonaws "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/random"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/require"
)

// TestMain sets up Localstack/Dynamo for all tests in the ondemand package and tears down after.
func TestMain(m *testing.M) {
	_, cleanup := test.GetOrDeployLocalstack()
	defer cleanup()

	code := m.Run()
	os.Exit(code)
}

// createPaymentTable creates a DynamoDB table for on-demand payment testing
// Uses the existing CreateOnDemandTable function from meterer package to ensure
// our test table schema exactly matches the production schema.
// Appends a random suffix to the table name to prevent collisions between tests.
func createPaymentTable(t *testing.T, tableName string) string {
	t.Helper()
	testRandom := random.NewTestRandom()
	randomSuffix := testRandom.Intn(999999)
	fullTableName := fmt.Sprintf("%s_%d", tableName, randomSuffix)

	// Create local client config for table creation
	localstackPort := test.DefaultLocalstackPort
	if os.Getenv("DEPLOY_LOCALSTACK") == "false" {
		localstackPort = os.Getenv("LOCALSTACK_PORT")
	}

	clientConfig := commonaws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localstackPort),
	}

	err := meterer.CreateOnDemandTable(clientConfig, fullTableName)
	require.NoError(t, err, "failed to create on-demand table")

	return fullTableName
}

// deleteTable deletes a DynamoDB table used in testing
func deleteTable(t *testing.T, tableName string) {
	t.Helper()
	ctx := t.Context()

	dynamoClient, cleanup := test.GetOrDeployLocalstack()
	defer cleanup()

	_, err := dynamoClient.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	require.NoError(t, err, "failed to delete table")
}
