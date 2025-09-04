package testbed

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

const (
	LocalStackImage = "localstack/localstack:latest"
	LocalStackPort  = "4566/tcp"
)

// LocalStackOptions configures the LocalStack AWS simulation container
type LocalStackOptions struct {
	ExposeHostPort bool           // If true, binds container port 4566 to host port (default: 4570)
	HostPort       string         // Custom host port to bind to (defaults to "4570" if empty and ExposeHostPort is true)
	Services       []string       // AWS services to enable (defaults to s3, dynamodb, kms)
	Region         string         // AWS region (defaults to us-east-1)
	Debug          bool           // Enable debug logging
	Logger         logging.Logger // Logger for container operations (required)
}

// LocalStackContainer wraps the official LocalStack testcontainers module
type LocalStackContainer struct {
	container *localstack.LocalStackContainer
	options   LocalStackOptions
	endpoint  string
	logger    logging.Logger
}

// NewLocalStackContainer creates and starts a new LocalStack container with a default noop logger
func NewLocalStackContainer(ctx context.Context) (*LocalStackContainer, error) {
	// Create a silent logger that discards all output
	noopLogger := logging.NewTextSLogger(io.Discard, &logging.SLoggerOptions{})
	return NewLocalStackContainerWithOptions(ctx, LocalStackOptions{
		Logger: noopLogger,
	})
}

// NewLocalStackContainerWithOptions creates and starts a new LocalStack container with custom options
func NewLocalStackContainerWithOptions(ctx context.Context, opts LocalStackOptions) (*LocalStackContainer, error) {
	if opts.Logger == nil {
		return nil, fmt.Errorf("logger is required in LocalStackOptions")
	}
	
	// Set defaults
	if len(opts.Services) == 0 {
		opts.Services = []string{"s3", "dynamodb", "kms"}
	}
	if opts.Region == "" {
		opts.Region = "us-east-1"
	}

	logger := opts.Logger
	logger.Info("Starting LocalStack container", "services", opts.Services, "region", opts.Region)

	var customizers []testcontainers.ContainerCustomizer
	
	// Add logger
	customizers = append(customizers, testcontainers.WithLogger(newTestcontainersLogger(logger)))
	
	env := buildLocalStackEnv(opts)
	customizers = append(customizers, testcontainers.WithEnv(env))

	// Add host port binding if requested
	if opts.ExposeHostPort {
		hostPort := opts.HostPort
		if hostPort == "" {
			hostPort = "4570" // Default to 4570 for LocalStack (similar to Anvil using 8545)
		}
		customizers = append(customizers, testcontainers.WithHostConfigModifier(func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				LocalStackPort: []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: hostPort,
					},
				},
			}
		}))
	}

	// Start the container using the official module
	logger.Debug("Creating LocalStack container with image", "image", LocalStackImage)
	container, err := localstack.Run(ctx, LocalStackImage, customizers...)
	if err != nil {
		logger.Error("Failed to start LocalStack container", "error", err)
		return nil, fmt.Errorf("failed to start localstack container: %w", err)
	}

	// Get the endpoint
	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "4566")
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	endpoint := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())

	logger.Info("LocalStack container started successfully", "endpoint", endpoint)

	return &LocalStackContainer{
		container: container,
		options:   opts,
		endpoint:  endpoint,
		logger:    logger,
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
	return ls.options.Region
}

// Services returns the list of enabled AWS services
func (ls *LocalStackContainer) Services() []string {
	return ls.options.Services
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
		"AWS_DEFAULT_REGION":    ls.options.Region,
		"AWS_ENDPOINT_URL":      ls.Endpoint(),
	}
}

// Terminate stops and removes the container
func (ls *LocalStackContainer) Terminate(ctx context.Context) error {
	if ls == nil || ls.container == nil {
		return nil
	}
	ls.logger.Info("Terminating LocalStack container")
	if err := ls.container.Terminate(ctx); err != nil {
		ls.logger.Error("Failed to terminate LocalStack container", "error", err)
		return fmt.Errorf("failed to terminate LocalStack container: %w", err)
	}
	ls.logger.Debug("LocalStack container terminated successfully")
	return nil
}

// buildLocalStackEnv constructs environment variables for LocalStack
func buildLocalStackEnv(opts LocalStackOptions) map[string]string {
	env := map[string]string{
		"SERVICES":            strings.Join(opts.Services, ","),
		"DEFAULT_REGION":      opts.Region,
		"HOSTNAME_EXTERNAL":   "localhost",
		"DISABLE_CORS_CHECKS": "1",
	}

	if opts.Debug {
		env["DEBUG"] = "1"
		env["LS_LOG"] = "debug"
	}

	return env
}

// Deprecated: Use LocalStackOptions instead
type LocalStackConfig = LocalStackOptions

// Deprecated: Use NewLocalStackContainerWithOptions instead
func DefaultLocalStackConfig() LocalStackOptions {
	return LocalStackOptions{
		Services: []string{"s3", "dynamodb", "kms"},
		Region:   "us-east-1",
		Debug:    false,
	}
}
