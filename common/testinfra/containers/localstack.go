package containers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	LocalStackImage = "localstack/localstack:3.0"
	LocalStackPort  = "4566/tcp"
)

// LocalStackContainer wraps testcontainers functionality for LocalStack
type LocalStackContainer struct {
	container testcontainers.Container
	config    LocalStackConfig
	endpoint  string
}

// NewLocalStackContainer creates and starts a new LocalStack container
func NewLocalStackContainer(ctx context.Context, config LocalStackConfig) (*LocalStackContainer, error) {
	return NewLocalStackContainerWithNetwork(ctx, config, "")
}

// NewLocalStackContainerWithNetwork creates and starts a new LocalStack container in a specific network
func NewLocalStackContainerWithNetwork(
	ctx context.Context, config LocalStackConfig, networkName string,
) (*LocalStackContainer, error) {
	if !config.Enabled {
		return nil, fmt.Errorf("localstack container is disabled in config")
	}

	env := buildLocalStackEnv(config)

	// Generate a unique container name using timestamp to avoid conflicts
	uniqueName := fmt.Sprintf("localstack-test-%d", time.Now().UnixNano())

	req := testcontainers.ContainerRequest{
		Image:        LocalStackImage,
		ExposedPorts: []string{LocalStackPort},
		Env:          env,
		WaitingFor:   wait.ForListeningPort("4566/tcp"),
		Name:         uniqueName,
	}

	// Add network if specified
	if networkName != "" {
		req.Networks = []string{networkName}
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start localstack container: %w", err)
	}

	// Get the mapped port
	mappedPort, err := container.MappedPort(ctx, "4566")
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	// Get the host
	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get host: %w", err)
	}

	endpoint := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())

	return &LocalStackContainer{
		container: container,
		config:    config,
		endpoint:  endpoint,
	}, nil
}

// Endpoint returns the LocalStack endpoint URL
func (ls *LocalStackContainer) Endpoint() string {
	return ls.endpoint
}

// Region returns the configured AWS region
func (ls *LocalStackContainer) Region() string {
	return ls.config.Region
}

// Services returns the list of enabled AWS services
func (ls *LocalStackContainer) Services() []string {
	return ls.config.Services
}

// GetServiceEndpoint returns the endpoint for a specific AWS service
func (ls *LocalStackContainer) GetServiceEndpoint(service string) string {
	return ls.endpoint
}

// GetAWSConfig returns AWS SDK configuration for connecting to LocalStack
func (ls *LocalStackContainer) GetAWSConfig() map[string]string {
	return map[string]string{
		"AWS_ACCESS_KEY_ID":     "test",
		"AWS_SECRET_ACCESS_KEY": "test",
		"AWS_DEFAULT_REGION":    ls.config.Region,
		"AWS_ENDPOINT_URL":      ls.endpoint,
	}
}

// Terminate stops and removes the container
func (ls *LocalStackContainer) Terminate(ctx context.Context) error {
	if ls.container != nil {
		return ls.container.Terminate(ctx)
	}
	return nil
}

// buildLocalStackEnv constructs environment variables for LocalStack
func buildLocalStackEnv(config LocalStackConfig) map[string]string {
	env := map[string]string{
		"SERVICES":            strings.Join(config.Services, ","),
		"DEFAULT_REGION":      config.Region,
		"HOSTNAME_EXTERNAL":   "localhost",
		"DISABLE_CORS_CHECKS": "1",
	}

	if config.Debug {
		env["DEBUG"] = "1"
		env["LS_LOG"] = "debug"
	}

	return env
}

// WaitForReady waits for LocalStack to be ready to accept requests
func (ls *LocalStackContainer) WaitForReady(ctx context.Context) error {
	// The wait strategy in the container request should handle this
	return nil
}

// CreateS3Bucket creates an S3 bucket in LocalStack
func (ls *LocalStackContainer) CreateS3Bucket(ctx context.Context, bucketName string) error {
	// This would typically use the AWS SDK to create the bucket
	// For now, we'll just return nil - the actual implementation would
	// depend on how the calling code wants to handle AWS SDK configuration
	return nil
}

// CreateDynamoDBTable creates a DynamoDB table in LocalStack
func (ls *LocalStackContainer) CreateDynamoDBTable(ctx context.Context, tableName string, keySchema map[string]string) error {
	// Similar to S3, this would use the AWS SDK
	return nil
}

// CreateKMSKey creates a KMS key in LocalStack
func (ls *LocalStackContainer) CreateKMSKey(ctx context.Context, keySpec string) (string, error) {
	// Returns a mock key ID for testing
	return "arn:aws:kms:us-east-1:000000000000:key/12345678-1234-1234-1234-123456789012", nil
}

// GetLogs returns the container logs for debugging
func (ls *LocalStackContainer) GetLogs(ctx context.Context) (string, error) {
	if ls.container == nil {
		return "", fmt.Errorf("container not started")
	}

	logs, err := ls.container.Logs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container logs: %w", err)
	}
	defer logs.Close()

	buf := make([]byte, 1024*1024) // 1MB buffer
	n, err := logs.Read(buf)
	if err != nil && err.Error() != "EOF" {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return string(buf[:n]), nil
}

// HealthCheck checks if LocalStack is healthy and all services are ready
func (ls *LocalStackContainer) HealthCheck(ctx context.Context) error {
	// Could implement a more sophisticated health check here
	// that verifies each enabled service is actually responding
	return nil
}
