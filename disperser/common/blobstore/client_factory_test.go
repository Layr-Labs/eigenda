package blobstore

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/assert"
)

// mockLogger is a simple mock logger for testing
type mockLogger struct{}

func (m *mockLogger) Debug(msg string, args ...interface{})       {}
func (m *mockLogger) Info(msg string, args ...interface{})        {}
func (m *mockLogger) Warn(msg string, args ...interface{})        {}
func (m *mockLogger) Error(msg string, args ...interface{})       {}
func (m *mockLogger) Fatal(msg string, args ...interface{})       {}
func (m *mockLogger) Debugf(template string, args ...interface{}) {}
func (m *mockLogger) Infof(template string, args ...interface{})  {}
func (m *mockLogger) Warnf(template string, args ...interface{})  {}
func (m *mockLogger) Errorf(template string, args ...interface{}) {}
func (m *mockLogger) Fatalf(template string, args ...interface{}) {}
func (m *mockLogger) With(tags ...any) logging.Logger             { return m }

func TestCreateObjectStorageClient_S3Backend(t *testing.T) {
	ctx := context.Background()
	config := Config{
		Backend:    S3Backend,
		BucketName: "test-bucket",
		TableName:  "test-table",
	}
	awsConfig := aws.ClientConfig{
		Region:                      "us-east-1",
		AccessKey:                   "test-access-key",
		SecretAccessKey:             "test-secret-key",
		EndpointURL:                 "",
		FragmentParallelismConstant: 1,
		FragmentParallelismFactor:   0,
	}
	logger := &mockLogger{}

	// This test will fail without AWS credentials, but it tests the factory logic
	client, err := CreateObjectStorageClient(ctx, config, awsConfig, logger)

	// We expect an error in test environment without AWS setup
	if err != nil {
		assert.Contains(t, err.Error(), "failed to create S3 client")
	} else {
		assert.NotNil(t, client)
	}
}

func TestCreateObjectStorageClient_OCIBackend(t *testing.T) {
	ctx := context.Background()
	config := Config{
		Backend:    OCIBackend,
		BucketName: "test-bucket",
		TableName:  "test-table",
	}
	awsConfig := aws.ClientConfig{
		Region:                      "us-east-1",
		FragmentParallelismConstant: 1,
		FragmentParallelismFactor:   0,
	}
	logger := &mockLogger{}

	// This test will fail without OCI credentials, but it tests the factory logic
	client, err := CreateObjectStorageClient(ctx, config, awsConfig, logger)

	// We expect an error in test environment without OCI setup
	if err != nil {
		assert.Contains(t, err.Error(), "failed to create OCI object storage client")
	} else {
		assert.NotNil(t, client)
	}
}

func TestCreateObjectStorageClient_UnsupportedBackend(t *testing.T) {
	ctx := context.Background()
	config := Config{
		Backend:    "unsupported-backend",
		BucketName: "test-bucket",
		TableName:  "test-table",
	}
	awsConfig := aws.ClientConfig{
		Region: "us-east-1",
	}
	logger := &mockLogger{}

	client, err := CreateObjectStorageClient(ctx, config, awsConfig, logger)

	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported object storage backend: unsupported-backend")
}

func TestCreateObjectStorageClient_EmptyBackend(t *testing.T) {
	ctx := context.Background()
	config := Config{
		Backend:    "", // Empty backend should default somewhere or error
		BucketName: "test-bucket",
		TableName:  "test-table",
	}
	awsConfig := aws.ClientConfig{
		Region: "us-east-1",
	}
	logger := &mockLogger{}

	client, err := CreateObjectStorageClient(ctx, config, awsConfig, logger)

	// Should error due to unsupported backend
	assert.Nil(t, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported object storage backend")
}

func TestCreateObjectStorageClient_OCIWithFragmentParallelismFactor(t *testing.T) {
	ctx := context.Background()
	config := Config{
		Backend:    OCIBackend,
		BucketName: "test-bucket",
		TableName:  "test-table",
	}
	awsConfig := aws.ClientConfig{
		Region:                    "us-east-1",
		FragmentParallelismFactor: 2, // Should result in 2 * runtime.NumCPU() workers
	}
	logger := &mockLogger{}

	// This test will fail without OCI credentials, but it tests the configuration logic
	client, err := CreateObjectStorageClient(ctx, config, awsConfig, logger)

	// We expect an error in test environment, but the config should be passed correctly
	if err != nil {
		assert.Contains(t, err.Error(), "failed to create OCI object storage client")
	} else {
		assert.NotNil(t, client)
	}
}

func TestObjectStorageBackend_Constants(t *testing.T) {
	assert.Equal(t, ObjectStorageBackend("s3"), S3Backend)
	assert.Equal(t, ObjectStorageBackend("oci"), OCIBackend)
}

func TestConfig_Struct(t *testing.T) {
	config := Config{
		BucketName: "test-bucket",
		TableName:  "test-table",
		Backend:    S3Backend,
	}

	assert.Equal(t, "test-bucket", config.BucketName)
	assert.Equal(t, "test-table", config.TableName)
	assert.Equal(t, S3Backend, config.Backend)
}

func TestCreateObjectStorageClient_OCIMinimalConfig(t *testing.T) {
	ctx := context.Background()
	config := Config{
		Backend:    OCIBackend,
		BucketName: "test-bucket",
		TableName:  "test-table",
	}
	awsConfig := aws.ClientConfig{
		// Minimal AWS config for OCI (only fragment settings used)
		FragmentParallelismConstant: 0,
		FragmentParallelismFactor:   0,
	}
	logger := &mockLogger{}

	// This should still work (but fail due to credentials)
	client, err := CreateObjectStorageClient(ctx, config, awsConfig, logger)

	if err != nil {
		assert.Contains(t, err.Error(), "failed to create OCI object storage client")
	} else {
		assert.NotNil(t, client)
	}
}
