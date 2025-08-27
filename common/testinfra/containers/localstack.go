package containers

import (
	"context"
	"fmt"
	"strings"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
	"github.com/testcontainers/testcontainers-go/network"
)

const (
	LocalStackImage = "localstack/localstack:latest"
)

// LocalStackConfig configures the LocalStack AWS simulation container
type LocalStackConfig struct {
	Enabled  bool     `json:"enabled"`
	Services []string `json:"services"` // AWS services to enable: s3, dynamodb, kms, secretsmanager
	Region   string   `json:"region"`
	Debug    bool     `json:"debug"`
}

// LocalStackContainer wraps the official LocalStack testcontainers module
type LocalStackContainer struct {
	container *localstack.LocalStackContainer
	config    LocalStackConfig
	endpoint  string
}

// NewLocalStackContainerWithNetwork creates and starts a new LocalStack container in a specific network
func NewLocalStackContainerWithNetwork(
	ctx context.Context, config LocalStackConfig, nw *testcontainers.DockerNetwork,
) (*LocalStackContainer, error) {
	if !config.Enabled {
		return nil, fmt.Errorf("localstack container is disabled in config")
	}

	if nw == nil {
		return nil, fmt.Errorf("network is required - LocalStack containers must use a shared network")
	}

	// Build container customizers
	var opts []testcontainers.ContainerCustomizer

	// Add environment variables
	env := buildLocalStackEnv(config)
	opts = append(opts, testcontainers.WithEnv(env))
	
	// Add network configuration using the recommended approach
	opts = append(opts, network.WithNetwork([]string{"localstack"}, nw))

	// Start the container using the official module
	container, err := localstack.Run(ctx, LocalStackImage, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to start localstack container: %w", err)
	}

	// Get the endpoint immediately after container starts
	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "4566")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
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

// InternalEndpoint returns the LocalStack endpoint URL for internal Docker network communication
func (ls *LocalStackContainer) InternalEndpoint() string {
	return "http://localstack:4566"
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
	// All services use the same endpoint in LocalStack v2+
	return ls.Endpoint()
}

// GetAWSConfig returns AWS SDK configuration for connecting to LocalStack
func (ls *LocalStackContainer) GetAWSConfig() map[string]string {
	return map[string]string{
		"AWS_ACCESS_KEY_ID":     "test",
		"AWS_SECRET_ACCESS_KEY": "test",
		"AWS_DEFAULT_REGION":    ls.config.Region,
		"AWS_ENDPOINT_URL":      ls.Endpoint(),
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
	// The official module handles waiting automatically during Run()
	return nil
}
