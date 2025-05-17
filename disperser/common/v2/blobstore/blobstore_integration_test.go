package blobstore_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	cryptorand "crypto/rand"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
	"github.com/Layr-Labs/eigenda/common/aws/mock"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
)

// testEnv holds all the test dependencies
type testEnv struct {
	logger                    logging.Logger
	dockertestPool            *dockertest.Pool
	dockertestResource        *dockertest.Resource
	postgresContainer         *dockertest.Resource
	deployLocalStack          bool
	localStackPort            string
	usePostgres               bool
	s3Client                  s3.Client
	dynamoClient              dynamodb.Client
	mockDynamoClient          *mock.MockDynamoDBClient
	blobStore                 *blobstore.BlobStore
	blobMetadataStore         blobstore.MetadataStore
	mockedBlobMetadataStore   *blobstore.BlobMetadataStore
	postgresBlobMetadataStore *blobstore.PostgresBlobMetadataStore
	UUID                      uuid.UUID
	s3BucketName              string
	metadataTableName         string
	pgUser                    string
	pgPassword                string
	pgDB                      string
	mockCommitment            encoding.BlobCommitments
	rng                       *rand.Rand // Deterministic random source
}

// newTestEnv creates a new test environment
func newTestEnv() *testEnv {
	testUUID := uuid.New()

	// Create a deterministic random source with seed=1
	// This ensures tests are reproducible
	rng := rand.New(rand.NewSource(1))

	return &testEnv{
		logger:            testutils.GetLogger(),
		localStackPort:    "4571",
		usePostgres:       os.Getenv("USE_POSTGRES") == "true",
		deployLocalStack:  !(os.Getenv("DEPLOY_LOCALSTACK") == "false"),
		UUID:              testUUID,
		s3BucketName:      "test-eigenda-blobstore",
		metadataTableName: fmt.Sprintf("test-BlobMetadata-%v", testUUID),
		pgUser:            "postgres",
		pgPassword:        "postgres",
		pgDB:              "testdb",
		mockCommitment:    encoding.BlobCommitments{},
		rng:               rng,
	}
}

// Global test environment for TestMain
var globalEnv *testEnv

// Creates a new test environment for each test
func setupForTest(t *testing.T) *testEnv {
	t.Helper()

	// Create a copy of the global environment's configuration
	env := &testEnv{
		logger:            testutils.GetLogger().With("test", t.Name()),
		localStackPort:    globalEnv.localStackPort,
		usePostgres:       globalEnv.usePostgres,
		deployLocalStack:  globalEnv.deployLocalStack,
		UUID:              globalEnv.UUID,
		s3BucketName:      globalEnv.s3BucketName,
		metadataTableName: globalEnv.metadataTableName,
		pgUser:            globalEnv.pgUser,
		pgPassword:        globalEnv.pgPassword,
		pgDB:              globalEnv.pgDB,
	}

	// We reuse the clients and stores from the global environment
	env.s3Client = globalEnv.s3Client
	env.dynamoClient = globalEnv.dynamoClient
	env.mockDynamoClient = globalEnv.mockDynamoClient
	env.blobStore = globalEnv.blobStore
	env.blobMetadataStore = globalEnv.blobMetadataStore
	env.mockedBlobMetadataStore = globalEnv.mockedBlobMetadataStore
	env.postgresContainer = globalEnv.postgresContainer
	env.postgresBlobMetadataStore = globalEnv.postgresBlobMetadataStore
	env.mockCommitment = globalEnv.mockCommitment

	// Create a deterministic random source for this test
	// Using test name as part of the seed for better isolation
	seed := int64(1)
	for _, c := range t.Name() {
		seed += int64(c)
	}
	env.rng = rand.New(rand.NewSource(seed))

	// Reset database state for this test
	if err := resetDatabaseState(env); err != nil {
		t.Fatalf("Failed to reset database state: %v", err)
	}

	// Register cleanup
	t.Cleanup(func() {
		// Do nothing.
	})

	return env
}

// resetDatabaseState resets the database state for a new test
func resetDatabaseState(env *testEnv) error {
	if env.usePostgres {
		// For PostgreSQL, we drop and recreate all tables
		if env.postgresBlobMetadataStore != nil {
			return env.postgresBlobMetadataStore.ResetTables()
		}
	} else {
		// For DynamoDB, we recreate the table
		ctx := context.Background()

		// Delete the table if it exists
		// We ignore errors here since the table might not exist
		err := test_utils.DeleteTable(ctx, aws.ClientConfig{
			Region:          "us-east-1",
			AccessKey:       "localstack",
			SecretAccessKey: "localstack",
			EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", env.localStackPort),
		}, env.metadataTableName)

		// Sleep briefly to allow table deletion to complete
		time.Sleep(500 * time.Millisecond)

		// Recreate the table
		_, err = test_utils.CreateTable(ctx, aws.ClientConfig{
			Region:          "us-east-1",
			AccessKey:       "localstack",
			SecretAccessKey: "localstack",
			EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", env.localStackPort),
		}, env.metadataTableName, blobstore.GenerateTableSchema(env.metadataTableName, 10, 10))
		if err != nil {
			return fmt.Errorf("failed to recreate DynamoDB table: %w", err)
		}

		// Sleep briefly to allow table creation to complete
		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func TestMain(m *testing.M) {
	// Initialize global test environment
	globalEnv = newTestEnv()

	// Setup resources before tests
	if err := setupResources(globalEnv); err != nil {
		fmt.Printf("Failed to set up test resources: %v\n", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()
	if r := recover(); r != nil {
		// Log the panic so it surfaces in CI logs
		fmt.Printf("Recovered in TestMain: %v\n", r)
	}

	// Clean up containers on exit
	teardownResources(globalEnv)
	os.Exit(code)
}

func setupResources(env *testEnv) error {
	// Start LocalStack if needed - we need this for S3 regardless of database type
	if env.deployLocalStack {
		var err error
		env.dockertestPool, env.dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(env.localStackPort)
		if err != nil {
			return fmt.Errorf("failed to start localstack container: %w", err)
		}
	}

	// Set up either PostgreSQL or DynamoDB based on flag
	if env.usePostgres {
		if err := setupPostgres(env); err != nil {
			return fmt.Errorf("failed to setup postgres: %w", err)
		}
	} else {
		if err := setupDynamoDB(env); err != nil {
			return fmt.Errorf("failed to setup dynamodb: %w", err)
		}
	}

	// Set up S3 for blob storage (used for all metadata store types)
	if err := setupS3(env); err != nil {
		return fmt.Errorf("failed to setup S3: %w", err)
	}

	// Set up mock data
	if err := setupMockData(env); err != nil {
		return fmt.Errorf("failed to setup mock data: %w", err)
	}

	return nil
}

func setupPostgres(env *testEnv) error {
	var err error
	// Create a new pool
	if env.dockertestPool == nil {
		env.dockertestPool, err = dockertest.NewPool("")
		if err != nil {
			return fmt.Errorf("could not connect to docker: %w", err)
		}
	}

	// Pull and run a PostgreSQL container
	env.postgresContainer, err = env.dockertestPool.Run("postgres", "13", []string{
		"POSTGRES_USER=" + env.pgUser,
		"POSTGRES_PASSWORD=" + env.pgPassword,
		"POSTGRES_DB=" + env.pgDB,
	})
	if err != nil {
		return fmt.Errorf("could not start PostgreSQL container: %w", err)
	}

	// Get the PostgreSQL port
	pgPortStr := env.postgresContainer.GetPort("5432/tcp")
	pgPort, err := strconv.Atoi(pgPortStr)
	if err != nil {
		teardownResources(env)
		return fmt.Errorf("failed to parse PostgreSQL port: %w", err)
	}

	// Wait for PostgreSQL to become ready
	if err := env.dockertestPool.Retry(func() error {
		// Create a test connection to PostgreSQL
		pgConfig := blobstore.PostgreSQLConfig{
			Host:     "localhost",
			Port:     pgPort,
			Username: env.pgUser,
			Password: env.pgPassword,
			Database: env.pgDB,
			SSLMode:  "disable",
		}

		// Try to create a metadata store (this will initialize the tables)
		var err error
		env.postgresBlobMetadataStore, err = blobstore.NewPostgresBlobMetadataStore(pgConfig, env.logger)
		return err
	}); err != nil {
		teardownResources(env)
		return fmt.Errorf("could not connect to PostgreSQL: %w", err)
	}

	// Use the PostgreSQL metadata store
	env.blobMetadataStore = env.postgresBlobMetadataStore
	return nil
}

func setupDynamoDB(env *testEnv) error {
	if !env.deployLocalStack {
		env.localStackPort = os.Getenv("LOCALSTACK_PORT")
	}

	// NOTE: The LocalStack container startup has been moved to setupResources
	// to ensure it runs regardless of which database we're using

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", env.localStackPort),
	}

	_, err := test_utils.CreateTable(context.Background(), cfg, env.metadataTableName, blobstore.GenerateTableSchema(env.metadataTableName, 10, 10))
	if err != nil {
		return fmt.Errorf("failed to create dynamodb table: %w", err)
	}

	env.dynamoClient, err = dynamodb.NewClient(cfg, env.logger)
	if err != nil {
		return fmt.Errorf("failed to create dynamodb client: %w", err)
	}
	env.mockDynamoClient = &mock.MockDynamoDBClient{}

	// Create DynamoDB metadata store
	dynamoMetadataStore := blobstore.NewBlobMetadataStore(env.dynamoClient, env.logger, env.metadataTableName)
	env.mockedBlobMetadataStore = blobstore.NewBlobMetadataStore(env.mockDynamoClient, env.logger, env.metadataTableName)

	// Use the DynamoDB metadata store
	env.blobMetadataStore = dynamoMetadataStore
	return nil
}

func setupS3(env *testEnv) error {
	// Set up AWS client config for S3
	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", env.localStackPort),
	}

	var err error
	env.s3Client, err = s3.NewClient(context.Background(), cfg, env.logger)
	if err != nil {
		return fmt.Errorf("failed to create s3 client: %w", err)
	}
	err = env.s3Client.CreateBucket(context.Background(), env.s3BucketName)
	if err != nil {
		return fmt.Errorf("failed to create s3 bucket: %w", err)
	}
	env.blobStore = blobstore.NewBlobStore(env.s3BucketName, env.s3Client, env.logger)
	return nil
}

// Helper function to create a fp.Element from a decimal string
func mustSetElementFromString(s string) fp.Element {
	var element fp.Element
	_, err := element.SetString(s)
	if err != nil {
		// We're using panic in this test helper (not in actual test code)
		// since this is initialization code that should never fail
		panic(fmt.Sprintf("failed to parse element string %s: %v", s, err))
	}
	return element
}

func setupMockData(env *testEnv) error {
	var X1, Y1 fp.Element
	X1 = *X1.SetBigInt(big.NewInt(1))
	Y1 = *Y1.SetBigInt(big.NewInt(2))

	// Use the helper to make this more readable and less error-prone
	lengthXA0 := mustSetElementFromString("10857046999023057135944570762232829481370756359578518086990519993285655852781")
	lengthXA1 := mustSetElementFromString("11559732032986387107991004021392285783925812861821192530917403151452391805634")
	lengthYA0 := mustSetElementFromString("8495653923123431417604973247489272438418190587263600148770280649306958101930")
	lengthYA1 := mustSetElementFromString("4082367875863433681332203403145435568316851327593401208105741076214120093531")

	var lengthProof, lengthCommitment bn254.G2Affine
	lengthProof.X.A0 = lengthXA0
	lengthProof.X.A1 = lengthXA1
	lengthProof.Y.A0 = lengthYA0
	lengthProof.Y.A1 = lengthYA1

	lengthCommitment = lengthProof

	env.mockCommitment = encoding.BlobCommitments{
		Commitment: &encoding.G1Commitment{
			X: X1,
			Y: Y1,
		},
		LengthCommitment: (*encoding.G2Commitment)(&lengthCommitment),
		LengthProof:      (*encoding.G2Commitment)(&lengthProof),
		Length:           10,
	}
	return nil
}

func teardownResources(env *testEnv) {
	// Cleanup PostgreSQL if it was used
	if env.usePostgres && env.postgresContainer != nil {
		// Close the PostgreSQL metadata store connection
		if env.postgresBlobMetadataStore != nil {
			env.postgresBlobMetadataStore.Close()
		}
		// Purge the PostgreSQL container
		if env.dockertestPool != nil {
			// Add nil check and error handling
			if err := env.dockertestPool.Purge(env.postgresContainer); err != nil {
				// Just log the error but don't fail the test
				fmt.Printf("Warning: Could not purge PostgreSQL container: %s\n", err)
			}
		}
	}

	// Cleanup LocalStack resources if they were deployed
	if env.deployLocalStack && env.dockertestResource != nil {
		// Wrap in a defer/recover to prevent test failures if PurgeDockertestResources panics
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Warning: Recovered from panic in PurgeDockertestResources: %v\n", r)
				}
			}()
			deploy.PurgeDockertestResources(env.dockertestPool, env.dockertestResource)
		}()
	}
}

// getStoreType returns a string indicating whether we're using PostgreSQL or DynamoDB
func (env *testEnv) getStoreType() string {
	if env.usePostgres {
		return "PostgreSQL"
	}
	return "DynamoDB"
}

func (env *testEnv) newBlob(t *testing.T) (corev2.BlobKey, *corev2.BlobHeader) {
	accountBytes := make([]byte, 32)
	_, err := rand.Read(accountBytes)
	require.NoError(t, err)
	accountID := gethcommon.HexToAddress(hex.EncodeToString(accountBytes))
	timestamp := time.Now().UnixNano()
	cumulativePayment, err := cryptorand.Int(cryptorand.Reader, big.NewInt(1024))
	require.NoError(t, err)
	sig := make([]byte, 32)
	_, err = cryptorand.Read(sig)
	require.NoError(t, err)
	bh := &corev2.BlobHeader{
		BlobVersion:     0,
		QuorumNumbers:   []core.QuorumID{0},
		BlobCommitments: env.mockCommitment,
		PaymentMetadata: core.PaymentMetadata{
			AccountID:         accountID,
			Timestamp:         timestamp,
			CumulativePayment: cumulativePayment,
		},
	}
	bk, err := bh.BlobKey()
	require.NoError(t, err)
	return bk, bh
}

// // newBlob creates a new blob for testing with deterministic random values
// func (env *testEnv) newBlob(t *testing.T) (corev2.BlobKey, *corev2.BlobHeader) {
// 	t.Helper()

// 	// Use our deterministic random source
// 	accountBytes := make([]byte, 32)
// 	_, err := env.rng.Read(accountBytes)
// 	require.NoError(t, err)

// 	accountID := gethcommon.HexToAddress(hex.EncodeToString(accountBytes))

// 	// Use a fixed timestamp for deterministic tests
// 	// The fixed timestamp ensures that tests are reproducible
// 	timestamp := int64(1633027200000000000) // 2021-10-01 00:00:00 UTC in nanoseconds

// 	// Generate deterministic random payment value
// 	cumulativePayment := big.NewInt(int64(env.rng.Intn(1024)))

// 	bh := &corev2.BlobHeader{
// 		BlobVersion:     0,
// 		QuorumNumbers:   []core.QuorumID{0},
// 		BlobCommitments: env.mockCommitment,
// 		PaymentMetadata: core.PaymentMetadata{
// 			AccountID:         accountID,
// 			Timestamp:         timestamp,
// 			CumulativePayment: cumulativePayment,
// 		},
// 	}

// 	bk, err := bh.BlobKey()
// 	require.NoError(t, err)
// 	return bk, bh
// }

// // TestSimpleMetadataStore tests basic operations on the metadata store
// func TestSimpleMetadataStore(t *testing.T) {
// 	// Create a test-specific environment
// 	env := setupForTest(t)

// 	defer func() {
// 		if r := recover(); r != nil {
// 			// Handle the panic, log it, or perform assertions on the recovered value
// 			t.Logf("Recovered from panic: %v", r)
// 			t.Fail()
// 		}
// 	}()

// 	// Create a test context with a timeout
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Create a test blob
// 	blobKey, blobHeader := env.newBlob(t)

// 	// Create a BlobMetadata with fixed deterministic timestamps
// 	// Fixed timestamps ensure that tests are reproducible
// 	requestedAt := uint64(1633027200000000000) // 2021-10-01 00:00:00 UTC in nanoseconds
// 	updatedAt := requestedAt + 1000000000      // 1 second later

// 	blobMetadata := &commonv2.BlobMetadata{
// 		BlobHeader:  blobHeader,
// 		RequestedAt: requestedAt,
// 		BlobStatus:  commonv2.Queued,
// 		UpdatedAt:   updatedAt,
// 	}

// 	// Put the metadata
// 	err := env.blobMetadataStore.PutBlobMetadata(ctx, blobMetadata)
// 	require.NoError(t, err, "Failed to put blob metadata")

// 	// Get the metadata
// 	retrievedMetadata, err := env.blobMetadataStore.GetBlobMetadata(ctx, blobKey)
// 	require.NoError(t, err, "Failed to get blob metadata")

// 	// Verify the metadata
// 	t.Logf("Store type: %s", env.getStoreType())
// 	t.Logf("Retrieved blob status: %v", retrievedMetadata.BlobStatus)
// 	require.Equal(t, commonv2.Queued, retrievedMetadata.BlobStatus, "Unexpected blob status")

// 	// Update the status
// 	err = env.blobMetadataStore.UpdateBlobStatus(ctx, blobKey, commonv2.Encoded)
// 	require.NoError(t, err, "Failed to update blob status")

// 	// Get the metadata again
// 	retrievedMetadata, err = env.blobMetadataStore.GetBlobMetadata(ctx, blobKey)
// 	require.NoError(t, err, "Failed to get blob metadata after update")

// 	// Verify the updated status
// 	require.Equal(t, commonv2.Encoded, retrievedMetadata.BlobStatus, "Unexpected updated blob status")
// }
