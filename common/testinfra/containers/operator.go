package containers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// OperatorConfig defines configuration for an EigenDA operator node
type OperatorConfig struct {
	// Operator identification
	ID       int    // Operator index (0-based)
	Hostname string // Hostname for the operator

	// Root directory paths
	EigenDADirectory string // Path to eigenda root directory (for locating resources)

	// Network ports
	DispersalPort           string
	RetrievalPort           string
	InternalDispersalPort   string
	InternalRetrievalPort   string
	V2DispersalPort         string
	V2RetrievalPort         string
	InternalV2DispersalPort string
	InternalV2RetrievalPort string
	NodeAPIPort             string
	MetricsPort             string

	// Chain configuration
	ChainRPC              string
	PrivateKey            string // BLS private key
	BlsKeyFile            string // Path to BLS key file
	BlsKeyPassword        string // Password for BLS key file
	EcdsaPrivateKey       string // ECDSA private key for transactions
	EcdsaKeyFile          string // Path to ECDSA key file
	EcdsaKeyPassword      string // Password for ECDSA key file
	PublicIPCheckInterval string // How often to check public IP (e.g., "10s")
	NumConfirmations      int

	// Contract addresses
	BLSOperatorStateRetriever string
	EigenDAServiceManager     string

	// Storage configuration
	DBPath           string
	LogPath          string
	IdleDuration     string
	DbSizePollPeriod string
	MinFreeSpace     uint64

	// Encoding configuration
	EncodingVersion                      string
	G1Path                               string
	G2Path                               string
	CachePath                            string
	SRSOrder                             uint64
	SRSLoad                              uint64
	NumBatchDeserializationWorkers       int
	NumBatchHeaderDeserializationWorkers int
	EnableGnarkBundleEncoding            bool

	// Feature flags
	EnableNodeAPI  bool
	EnableMetrics  bool
	EnableTestMode bool
	NodeMode       string // v1-only, v2-only, v1-and-v2

	// Rate limiting
	RetrievalRateLimit            int
	RetrievalBucketSize           int
	AttestationProtocolLimit      int
	AttestationProtocolBucketSize int

	// Logging
	LogFormat string
	LogLevel  string

	// Container image
	Image string
}

// OperatorContainer represents a running operator container
type OperatorContainer struct {
	testcontainers.Container
	config  OperatorConfig
	logPath string
}

// DefaultOperatorConfig returns a default operator configuration
func DefaultOperatorConfig(id int) OperatorConfig {
	// Calculate base ports for this operator
	// Each operator gets a range of 10 ports for different services
	// Starting at 32010 to avoid conflicts with disperser (32003, 32005) and other services
	basePort := 32010 + (id * 10)

	return OperatorConfig{
		ID:                      id,
		Hostname:                "0.0.0.0",
		DispersalPort:           fmt.Sprintf("%d", basePort+1),
		RetrievalPort:           fmt.Sprintf("%d", basePort+2),
		InternalDispersalPort:   fmt.Sprintf("%d", basePort+1),
		InternalRetrievalPort:   fmt.Sprintf("%d", basePort+2),
		V2DispersalPort:         fmt.Sprintf("%d", basePort+3),
		V2RetrievalPort:         fmt.Sprintf("%d", basePort+4),
		InternalV2DispersalPort: fmt.Sprintf("%d", basePort+3),
		InternalV2RetrievalPort: fmt.Sprintf("%d", basePort+4),
		NodeAPIPort:             fmt.Sprintf("%d", basePort+5),
		MetricsPort:             fmt.Sprintf("%d", basePort+6),

		NumConfirmations: 0,

		DBPath:           "/data/db",
		LogPath:          "/logs",
		IdleDuration:     "10s",
		DbSizePollPeriod: "10s",
		MinFreeSpace:     10 * 1024 * 1024 * 1024, // 10 GB

		EncodingVersion:                      "v0",
		G1Path:                               "",
		G2Path:                               "",
		CachePath:                            "",
		SRSOrder:                             10000,
		SRSLoad:                              10000,
		NumBatchDeserializationWorkers:       4,
		NumBatchHeaderDeserializationWorkers: 4,
		EnableGnarkBundleEncoding:            true,

		EnableNodeAPI:  true,
		EnableMetrics:  true,
		EnableTestMode: false,
		NodeMode:       "v1-and-v2",

		RetrievalRateLimit:            100,
		RetrievalBucketSize:           100,
		AttestationProtocolLimit:      100,
		AttestationProtocolBucketSize: 100,

		LogFormat: "text",
		LogLevel:  "debug",

		PublicIPCheckInterval: "10s",

		Image: "ghcr.io/layr-labs/eigenda/node:dev",
	}
}

// NewOperatorContainerWithNetwork creates and starts a new operator container with a custom network
func NewOperatorContainerWithNetwork(ctx context.Context, config OperatorConfig, nw *testcontainers.DockerNetwork) (*OperatorContainer, error) {
	// Create a temporary directory for logs
	logDir, err := os.MkdirTemp("", fmt.Sprintf("operator-%d-logs-*", config.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}
	logPath := filepath.Join(logDir, "operator.log")

	// Determine the secrets directory path (must be absolute)
	secretsDir, err := filepath.Abs(filepath.Join(config.EigenDADirectory, "testinfra", "secrets"))
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for secrets: %w", err)
	}
	fmt.Printf("DEBUG: Mounting secrets from: %s to /app/secrets\n", secretsDir)

	// Determine the resources directory path for KZG params (must be absolute)
	resourcesDir, err := filepath.Abs(filepath.Join(config.EigenDADirectory, "resources"))
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path for resources: %w", err)
	}
	fmt.Printf("DEBUG: Mounting resources from: %s to /app/resources\n", resourcesDir)

	// Create a temporary directory for database
	dbDir, err := os.MkdirTemp("", fmt.Sprintf("operator-%d-db-*", config.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to create db directory: %w", err)
	}

	// Build environment variables
	env := buildOperatorEnv(config)

	// Debug log critical environment variables
	fmt.Printf("DEBUG: Operator %d env - NODE_BLS_OPERATOR_STATE_RETRIVER=%s, NODE_EIGENDA_SERVICE_MANAGER=%s\n",
		config.ID, env["NODE_BLS_OPERATOR_STATE_RETRIVER"], env["NODE_EIGENDA_SERVICE_MANAGER"])

	// Configure container request (no command args needed, only env vars)
	req := testcontainers.ContainerRequest{
		Image: config.Image,
		Env:   env,
		ExposedPorts: []string{
			config.InternalDispersalPort + "/tcp",
			config.InternalRetrievalPort + "/tcp",
			config.InternalV2DispersalPort + "/tcp",
			config.InternalV2RetrievalPort + "/tcp",
			config.NodeAPIPort + "/tcp",
			config.MetricsPort + "/tcp",
		},
		Networks:       []string{},
		NetworkAliases: map[string][]string{},
		Mounts: testcontainers.ContainerMounts{
			{
				Source: testcontainers.GenericBindMountSource{
					HostPath: logDir,
				},
				Target: "/logs",
			},
			{
				Source: testcontainers.GenericBindMountSource{
					HostPath: dbDir,
				},
				Target: "/data",
			},
			{
				Source: testcontainers.GenericBindMountSource{
					HostPath: secretsDir,
				},
				Target: "/app/secrets",
			},
			{
				Source: testcontainers.GenericBindMountSource{
					HostPath: resourcesDir,
				},
				Target: "/app/resources",
			},
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(config.InternalDispersalPort+"/tcp")).WithStartupTimeout(60*time.Second),
			wait.ForLog("v2 dispersal enabled").WithStartupTimeout(90*time.Second),
		),
		Name:            fmt.Sprintf("eigenda-operator-%d", config.ID),
		AlwaysPullImage: false, // Use local image if available
	}

	// Add port bindings when hostname is 0.0.0.0 or using localhost domain
	// This ensures the operator ports are accessible from the host at the expected ports
	if config.Hostname == "0.0.0.0" || strings.HasSuffix(config.Hostname, ".localtest.me") {
		req.HostConfigModifier = func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				nat.Port(config.InternalDispersalPort + "/tcp"):   []nat.PortBinding{{HostPort: config.DispersalPort}},
				nat.Port(config.InternalRetrievalPort + "/tcp"):   []nat.PortBinding{{HostPort: config.RetrievalPort}},
				nat.Port(config.InternalV2DispersalPort + "/tcp"): []nat.PortBinding{{HostPort: config.V2DispersalPort}},
				nat.Port(config.InternalV2RetrievalPort + "/tcp"): []nat.PortBinding{{HostPort: config.V2RetrievalPort}},
				nat.Port(config.NodeAPIPort + "/tcp"):             []nat.PortBinding{{HostPort: config.NodeAPIPort}},
				nat.Port(config.MetricsPort + "/tcp"):             []nat.PortBinding{{HostPort: config.MetricsPort}},
			}
		}
	}

	// Add network configuration if provided
	if nw != nil {
		req.Networks = []string{nw.Name}
		// Use localhost domain which resolves to 127.0.0.1
		// This allows the batcher to use host-gateway to reach operators
		operatorHostname := fmt.Sprintf("operator-%d.localtest.me", config.ID)
		req.NetworkAliases = map[string][]string{
			nw.Name: {operatorHostname, config.Hostname, fmt.Sprintf("operator-%d", config.ID)},
		}
	}

	// Create and start container
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           testcontainers.Logger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start operator container %d: %w", config.ID, err)
	}

	// Option A: Host IP (how the container can be reached from the host)
	host, _ := container.Host(ctx)
	fmt.Println("Host:", host)

	// Option B: Container’s internal IP (bridge network)
	insp, _ := container.Inspect(ctx)
	ip := insp.NetworkSettings.IPAddress
	fmt.Println("Container internal IP:", ip)

	// Option C: Mapped port (host port bound to container’s port)
	port, _ := container.MappedPort(ctx, "80")
	fmt.Println("Mapped port:", port.Int())

	fmt.Printf("Operator %d logs will be available at: %s\n", config.ID, logPath)

	return &OperatorContainer{
		Container: container,
		config:    config,
		logPath:   logPath,
	}, nil
}

// Config returns the operator configuration
func (c *OperatorContainer) Config() OperatorConfig {
	return c.config
}

// LogPath returns the path to the operator log file on the host
func (c *OperatorContainer) LogPath() string {
	return c.logPath
}

// GetDispersalAddress returns the full dispersal address for this operator (hostname:ports)
func (c *OperatorContainer) GetDispersalAddress() string {
	// Format: hostname:dispersal;retrieval;v2dispersal;v2retrieval
	return fmt.Sprintf("%s:%s;%s;%s;%s",
		c.config.Hostname,
		c.config.DispersalPort,
		c.config.RetrievalPort,
		c.config.V2DispersalPort,
		c.config.V2RetrievalPort,
	)
}

// GetInternalDispersalAddress returns the internal dispersal address for Docker network communication
func (c *OperatorContainer) GetInternalDispersalAddress() string {
	// Format: hostname:dispersal;retrieval;v2dispersal;v2retrieval
	return fmt.Sprintf("%s:%s;%s;%s;%s",
		c.config.Hostname,
		c.config.InternalDispersalPort,
		c.config.InternalRetrievalPort,
		c.config.InternalV2DispersalPort,
		c.config.InternalV2RetrievalPort,
	)
}

// buildOperatorEnv builds environment variables for the operator container
// All configuration is now done through environment variables only
func buildOperatorEnv(config OperatorConfig) map[string]string {
	// Strip 0x prefix from private keys if present
	privateKey := strings.TrimPrefix(config.PrivateKey, "0x")
	privateKey = strings.TrimPrefix(privateKey, "0X")

	ecdsaKey := strings.TrimPrefix(config.EcdsaPrivateKey, "0x")
	ecdsaKey = strings.TrimPrefix(ecdsaKey, "0X")

	env := map[string]string{
		// Identification
		"NODE_HOSTNAME": config.Hostname,

		// Ports
		"NODE_DISPERSAL_PORT":             config.DispersalPort,
		"NODE_RETRIEVAL_PORT":             config.RetrievalPort,
		"NODE_INTERNAL_DISPERSAL_PORT":    config.InternalDispersalPort,
		"NODE_INTERNAL_RETRIEVAL_PORT":    config.InternalRetrievalPort,
		"NODE_V2_DISPERSAL_PORT":          config.V2DispersalPort,
		"NODE_V2_RETRIEVAL_PORT":          config.V2RetrievalPort,
		"NODE_INTERNAL_V2_DISPERSAL_PORT": config.InternalV2DispersalPort,
		"NODE_INTERNAL_V2_RETRIEVAL_PORT": config.InternalV2RetrievalPort,
		"NODE_API_PORT":                   config.NodeAPIPort,
		"NODE_METRICS_PORT":               config.MetricsPort,

		// Features
		"NODE_ENABLE_NODE_API":  fmt.Sprintf("%t", config.EnableNodeAPI),
		"NODE_ENABLE_METRICS":   fmt.Sprintf("%t", config.EnableMetrics),
		"NODE_ENABLE_TEST_MODE": fmt.Sprintf("%t", config.EnableTestMode),
		"NODE_RUNTIME_MODE":     config.NodeMode,

		// Chain
		"NODE_CHAIN_RPC":                   config.ChainRPC,
		"NODE_TEST_PRIVATE_BLS":            privateKey,
		"NODE_NUM_CONFIRMATIONS":           fmt.Sprintf("%d", config.NumConfirmations),
		"NODE_BLS_OPERATOR_STATE_RETRIVER": config.BLSOperatorStateRetriever,
		"NODE_EIGENDA_SERVICE_MANAGER":     config.EigenDAServiceManager,

		// Storage
		"NODE_DB_PATH":             config.DBPath,
		"NODE_LOG_PATH":            filepath.Join(config.LogPath, "operator.log"),
		"NODE_IDLE_DURATION":       config.IdleDuration,
		"NODE_DB_SIZE_POLL_PERIOD": config.DbSizePollPeriod,
		"NODE_MIN_FREE_SPACE":      fmt.Sprintf("%d", config.MinFreeSpace),

		// Encoding
		"NODE_ENCODING_VERSION":                         config.EncodingVersion,
		"NODE_SRS_ORDER":                                fmt.Sprintf("%d", config.SRSOrder),
		"NODE_NUM_BATCH_DESERIALIZATION_WORKERS":        fmt.Sprintf("%d", config.NumBatchDeserializationWorkers),
		"NODE_NUM_BATCH_HEADER_DESERIALIZATION_WORKERS": fmt.Sprintf("%d", config.NumBatchHeaderDeserializationWorkers),
		"NODE_ENABLE_GNARK_BUNDLE_ENCODING":             fmt.Sprintf("%t", config.EnableGnarkBundleEncoding),

		// Rate limiting
		"NODE_RETRIEVAL_RATE_LIMIT":             fmt.Sprintf("%d", config.RetrievalRateLimit),
		"NODE_RETRIEVAL_BUCKET_SIZE":            fmt.Sprintf("%d", config.RetrievalBucketSize),
		"NODE_ATTESTATION_PROTOCOL_LIMIT":       fmt.Sprintf("%d", config.AttestationProtocolLimit),
		"NODE_ATTESTATION_PROTOCOL_BUCKET_SIZE": fmt.Sprintf("%d", config.AttestationProtocolBucketSize),

		// Logging
		"NODE_LOG_FORMAT": config.LogFormat,
		"NODE_LOG_LEVEL":  config.LogLevel,

		// Additional configuration that was previously in command flags
		"NODE_TIMEOUT":                  "10s",
		"NODE_QUORUM_ID_LIST":           "0,1",
		"NODE_PUBLIC_IP_PROVIDER":       "mockip",
		"NODE_CHURNER_URL":              "churner:32001",
		"NODE_CHURNER_USE_SECURE_GRPC":  "false",
		"NODE_PUBLIC_IP_CHECK_INTERVAL": "0", // Disable IP update checks in containerized environment
		"NODE_DISABLE_DISPERSAL_AUTH":   "false",
		"NODE_REGISTER_AT_NODE_START":   "true",
		"NODE_RELAY_USE_SECURE_GRPC":    "false", // Disable TLS for relay connections in test environment

		// KZG configuration
		"NODE_G1_PATH":    "/app/resources/srs/g1.point",
		"NODE_CACHE_PATH": "/app/resources/srs/SRSTables",
		"NODE_SRS_LOAD":   fmt.Sprintf("%d", config.SRSLoad),
	}

	// Add optional configurations if provided
	if config.G1Path != "" {
		env["NODE_G1_PATH"] = config.G1Path
	}
	if config.G2Path != "" {
		env["NODE_G2_PATH"] = config.G2Path
	} else {
		env["NODE_G2_POWER_OF_2_PATH"] = "/app/resources/srs/g2.point.powerOf2"
	}
	if config.CachePath != "" {
		env["NODE_CACHE_PATH"] = config.CachePath
	}

	// BLS key configuration
	// Set BLS key file if provided
	if config.BlsKeyFile != "" {
		env["NODE_BLS_KEY_FILE"] = config.BlsKeyFile
		env["NODE_BLS_KEY_PASSWORD"] = config.BlsKeyPassword
	}

	// ECDSA key configuration
	// Always set NODE_PRIVATE_KEY if we have an ECDSA private key
	if ecdsaKey != "" {
		env["NODE_PRIVATE_KEY"] = ecdsaKey
	}

	// Additionally set key file configuration if provided
	if config.EcdsaKeyFile != "" {
		env["NODE_ECDSA_KEY_FILE"] = config.EcdsaKeyFile
		env["NODE_ECDSA_KEY_PASSWORD"] = config.EcdsaKeyPassword
	}

	return env
}
