package testbed

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ChurnerConfig defines configuration for the churner container
type ChurnerConfig struct {
	// Enable churner service
	Enabled bool `env:"-"` // Skip in env mapping

	// Log configuration
	LogFormat string `env:"CHURNER_LOG_FORMAT"`
	LogLevel  string `env:"CHURNER_LOG_LEVEL"`

	// Network configuration
	Hostname string `env:"CHURNER_HOSTNAME"`
	GRPCPort string `env:"CHURNER_GRPC_PORT"`

	// EigenDA contract addresses
	EigenDADirectory       string `env:"CHURNER_EIGENDA_DIRECTORY"`
	OperatorStateRetriever string `env:"CHURNER_BLS_OPERATOR_STATE_RETRIVER"` // Note: typo in original
	ServiceManager         string `env:"CHURNER_EIGENDA_SERVICE_MANAGER"`

	// Chain configuration
	ChainRPC   string `env:"CHURNER_CHAIN_RPC"`
	PrivateKey string `env:"CHURNER_PRIVATE_KEY"`

	// Graph configuration
	GraphURL string `env:"CHURNER_GRAPH_URL"`

	// Metrics configuration
	EnableMetrics   bool   `env:"CHURNER_ENABLE_METRICS"`
	MetricsHTTPPort string `env:"CHURNER_METRICS_HTTP_PORT"`

	// Churner specific configuration
	PerPublicKeyRateLimit time.Duration `env:"CHURNER_PER_PUBLIC_KEY_RATE_LIMIT"`
	ChurnApprovalInterval time.Duration `env:"CHURNER_CHURN_APPROVAL_INTERVAL"`

	// Container configuration (not exposed as env vars)
	Image          string        `env:"-"`
	StartupTimeout time.Duration `env:"-"`
	ExposeHostPort bool          `env:"-"` // If true, binds container port to host port
	HostPort       string        `env:"-"` // Custom host port to bind to (defaults to GRPCPort if empty and ExposeHostPort is true)

	// Additional env vars that don't have direct struct fields
	LogPath string `env:"CHURNER_LOG_PATH"`
}

// ChurnerContainer represents a running churner container
type ChurnerContainer struct {
	testcontainers.Container
	config     ChurnerConfig
	url        string
	logPath    string
	logDir     string
	network    *testcontainers.DockerNetwork
	internalIP string
	cancelLog  context.CancelFunc
}

// DefaultChurnerConfig returns a default churner configuration suitable for testing
func DefaultChurnerConfig() ChurnerConfig {
	return ChurnerConfig{
		Enabled:                true,
		LogFormat:              "text",
		LogLevel:               "debug",
		Hostname:               "0.0.0.0",
		GRPCPort:               "32002",
		EigenDADirectory:       "", // Will be populated from contract deployment
		OperatorStateRetriever: "", // Will be populated from contract deployment
		ServiceManager:         "", // Will be populated from contract deployment
		ChainRPC:               "", // Will be populated from Anvil
		PrivateKey:             "", // Will be populated from deployer key
		GraphURL:               "", // Will be populated from GraphNode if enabled
		EnableMetrics:          true,
		MetricsHTTPPort:        "9095",
		PerPublicKeyRateLimit:  1 * time.Second, // Fast for testing
		ChurnApprovalInterval:  900 * time.Second,
		Image:                  "ghcr.io/layr-labs/eigenda/churner:dev",
		StartupTimeout:         30 * time.Second,
		ExposeHostPort:         false,
		HostPort:               "",
	}
}

// ToEnvMap converts the ChurnerConfig to environment variables
// directly from struct tags
func (c ChurnerConfig) ToEnvMap() (map[string]string, error) {
	env := make(map[string]string)
	v := reflect.ValueOf(c)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env")

		// Skip fields with no env tag or "-" tag
		if tag == "" || tag == "-" {
			continue
		}

		value := v.Field(i)
		var strValue string

		switch value.Kind() {
		case reflect.String:
			strValue = value.String()
		case reflect.Bool:
			strValue = strconv.FormatBool(value.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// Special handling for time.Duration which is int64
			if field.Type == reflect.TypeOf(time.Duration(0)) {
				strValue = value.Interface().(time.Duration).String()
			} else {
				strValue = strconv.FormatInt(value.Int(), 10)
			}
		default:
			// Skip unsupported types
			continue
		}

		// Only add non-empty values to the environment
		// The churner binary will handle its own defaults and validation
		if strValue != "" {
			env[tag] = strValue
		}
	}

	return env, nil
}

// NewChurnerContainerWithNetwork creates and starts a new churner container with a custom network
func NewChurnerContainerWithNetwork(ctx context.Context, config ChurnerConfig, network *testcontainers.DockerNetwork) (*ChurnerContainer, error) {
	if !config.Enabled {
		return nil, nil
	}

	env, err := config.ToEnvMap()
	if err != nil {
		return nil, fmt.Errorf("failed to build environment variables: %w", err)
	}

	// Configure container request with network
	req := testcontainers.ContainerRequest{
		Image:        config.Image,
		Env:          env,
		ExposedPorts: []string{config.GRPCPort + "/tcp"},
		Networks:     []string{network.Name},
		NetworkAliases: map[string][]string{
			network.Name: {"churner"},
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(config.GRPCPort+"/tcp")).WithStartupTimeout(config.StartupTimeout),
			wait.ForLog("churner server listening at").WithStartupTimeout(config.StartupTimeout),
		),
		Name:            "eigenda-churner",
		AlwaysPullImage: false,
	}

	// Add metrics port if enabled
	if config.EnableMetrics {
		req.ExposedPorts = append(req.ExposedPorts, config.MetricsHTTPPort+"/tcp")
	}

	// Add host port binding if requested
	if config.ExposeHostPort {
		hostPort := config.HostPort
		if hostPort == "" {
			hostPort = config.GRPCPort
		}
		req.HostConfigModifier = func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				nat.Port(config.GRPCPort + "/tcp"): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: hostPort,
					},
				},
			}
			// Also bind metrics port if enabled
			if config.EnableMetrics {
				metricsHostPort := config.MetricsHTTPPort
				hc.PortBindings[nat.Port(config.MetricsHTTPPort+"/tcp")] = []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: metricsHostPort,
					},
				}
			}
		}
	}

	// Create and start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start churner container: %w", err)
	}

	// Get the mapped port
	mappedPort, err := container.MappedPort(ctx, nat.Port(config.GRPCPort))
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	// Get the container host
	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	// Get internal IP address
	containerJSON, err := container.Inspect(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to inspect container: %w", err)
	}

	internalIP := ""
	if containerJSON.NetworkSettings != nil && containerJSON.NetworkSettings.Networks != nil {
		if networkInfo, ok := containerJSON.NetworkSettings.Networks[network.Name]; ok {
			internalIP = networkInfo.IPAddress
		}
	}

	churnerURL := fmt.Sprintf("%s:%s", host, mappedPort.Port())

	// Create a timestamped directory for logs
	timestamp := time.Now().Format("2006-01-02-15-04-05")
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	logDir := filepath.Join(cwd, "logs", fmt.Sprintf("churner-%s", timestamp))
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	hostLogPath := filepath.Join(logDir, "churner.log")

	logFile, err := os.Create(hostLogPath)
	if err != nil {
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to create log file: %w", err)
	}

	// Start streaming container logs to file
	logCtx, cancelLog := context.WithCancel(context.Background())
	logReader, err := container.Logs(logCtx)
	if err != nil {
		logFile.Close()
		cancelLog()
		container.Terminate(ctx)
		return nil, fmt.Errorf("failed to get log reader: %w", err)
	}

	// Stream logs in background
	go func() {
		defer logReader.Close()
		defer logFile.Close()
		_, _ = io.Copy(logFile, logReader)
	}()

	churner := &ChurnerContainer{
		Container:  container,
		config:     config,
		url:        churnerURL,
		logPath:    hostLogPath,
		logDir:     logDir,
		network:    network,
		internalIP: internalIP,
		cancelLog:  cancelLog,
	}

	// Dump config to .env file for debugging
	if err := churner.DumpConfigToEnv(); err != nil {
		fmt.Printf("Warning: failed to dump config to .env: %v\n", err)
	}

	fmt.Printf("Churner started successfully\n")
	fmt.Printf("  - External URL: %s\n", churnerURL)
	fmt.Printf("  - Internal URL: churner:%s\n", config.GRPCPort)
	fmt.Printf("  - Internal IP: %s\n", internalIP)
	fmt.Printf("  - Logs: %s\n", hostLogPath)
	if config.EnableMetrics {
		metricsPort, _ := container.MappedPort(ctx, nat.Port(config.MetricsHTTPPort))
		fmt.Printf("  - Metrics: http://%s:%s/metrics\n", host, metricsPort.Port())
	}

	return churner, nil
}

// URL returns the churner service URL accessible from the host
func (c *ChurnerContainer) URL() string {
	return c.url
}

// InternalURL returns the churner service URL for internal network communication
func (c *ChurnerContainer) InternalURL() string {
	return fmt.Sprintf("churner:%s", c.config.GRPCPort)
}

// InternalIP returns the internal IP address of the container
func (c *ChurnerContainer) InternalIP() string {
	return c.internalIP
}

// Config returns the churner configuration
func (c *ChurnerContainer) Config() ChurnerConfig {
	return c.config
}

// LogPath returns the path to the churner log file on the host
func (c *ChurnerContainer) LogPath() string {
	return c.logPath
}

// Network returns the Docker network this container is connected to
func (c *ChurnerContainer) Network() *testcontainers.DockerNetwork {
	return c.network
}

// DumpConfigToEnv writes the churner configuration to a .env file in the logs directory
func (c *ChurnerContainer) DumpConfigToEnv() error {
	if c.logDir == "" {
		return fmt.Errorf("log directory not set")
	}

	envPath := filepath.Join(c.logDir, "churner-config.env")

	// Get env map from config
	envMap, err := c.config.ToEnvMap()
	if err != nil {
		return fmt.Errorf("failed to convert config to env map: %w", err)
	}

	// Create the env file
	file, err := os.Create(envPath)
	if err != nil {
		return fmt.Errorf("failed to create env file: %w", err)
	}
	defer file.Close()

	// Write header
	_, err = fmt.Fprintf(file, "# Churner Configuration Dump\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "# Generated at: %s\n", time.Now().Format(time.RFC3339))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "# Container URL: %s\n", c.url)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "# Internal IP: %s\n", c.internalIP)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "\n")
	if err != nil {
		return err
	}

	// Write env vars in sorted order for consistency
	keys := make([]string, 0, len(envMap))
	for key := range envMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		_, err = fmt.Fprintf(file, "%s=%s\n", key, envMap[key])
		if err != nil {
			return err
		}
	}

	// Add additional runtime information as comments
	_, err = fmt.Fprintf(file, "\n# Runtime Information (not env vars)\n")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "# Image: %s\n", c.config.Image)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "# Startup Timeout: %s\n", c.config.StartupTimeout)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(file, "# Log Path: %s\n", c.logPath)
	if err != nil {
		return err
	}

	fmt.Printf("  - Config dump: %s\n", envPath)
	return nil
}

// Stop gracefully stops the churner container
func (c *ChurnerContainer) Stop(ctx context.Context) error {
	fmt.Printf("Stopping churner container...\n")

	// Cancel log streaming if it's running
	if c.cancelLog != nil {
		c.cancelLog()
	}

	if err := c.Container.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to stop churner container: %w", err)
	}

	// Log directory is kept for debugging purposes
	if c.logPath != "" {
		fmt.Printf("Churner logs preserved at: %s\n", c.logPath)
	}

	return nil
}
