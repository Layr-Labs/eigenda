package testbed

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	GraphNodeImage = "graphprotocol/graph-node:v0.35.0"
	PostgresImage  = "postgres:13"
	IPFSImage      = "ipfs/kubo:v0.24.0"

	GraphNodeHTTPPort    = "8000/tcp"
	GraphNodeWSPort      = "8001/tcp"
	GraphNodeJSONPort    = "8020/tcp"
	GraphNodeIndexPort   = "8030/tcp"
	GraphNodeMetricsPort = "8040/tcp"
	PostgresPort         = "5432/tcp"
	IPFSAPIPort          = "5001/tcp"
	IPFSGatewayPort      = "8080/tcp"
)

// GraphNodeOptions configures The Graph node container
type GraphNodeOptions struct {
	PostgresDB     string                        // Database name for Graph Node
	PostgresUser   string                        // Database user
	PostgresPass   string                        // Database password
	EthereumRPC    string                        // Ethereum RPC endpoint (will be set to Anvil RPC if Anvil is enabled)
	IPFSEndpoint   string                        // Optional external IPFS endpoint
	ExposeHostPort bool                          // If true, binds container ports to host
	HostHTTPPort   string                        // Custom host port for HTTP (defaults to "8000" if empty and ExposeHostPort is true)
	HostWSPort     string                        // Custom host port for WebSocket (defaults to "8001" if empty and ExposeHostPort is true)
	HostAdminPort  string                        // Custom host port for Admin (defaults to "8020" if empty and ExposeHostPort is true)
	HostIPFSPort   string                        // Custom host port for IPFS (defaults to "5001" if empty and ExposeHostPort is true)
	Logger         logging.Logger                // Logger for container operations (required)
	Network        *testcontainers.DockerNetwork // Optional Docker network to use
}

// GraphNodeContainer manages a Graph Node cluster with PostgreSQL and IPFS
type GraphNodeContainer struct {
	graphNode testcontainers.Container
	postgres  testcontainers.Container
	ipfs      testcontainers.Container
	network   *testcontainers.DockerNetwork
	httpURL   string
	wsURL     string
	adminURL  string
	logger    logging.Logger
}

// NewGraphNodeContainerWithOptions creates and starts a complete Graph Node setup with custom options
func NewGraphNodeContainerWithOptions(ctx context.Context, opts GraphNodeOptions) (*GraphNodeContainer, error) {
	if opts.Logger == nil {
		return nil, fmt.Errorf("logger is required in GraphNodeOptions")
	}
	logger := opts.Logger
	logger.Info("Starting Graph Node cluster")

	// Set defaults
	if opts.PostgresDB == "" {
		opts.PostgresDB = "graph-node"
	}
	if opts.PostgresUser == "" {
		opts.PostgresUser = "graph-node"
	}
	if opts.PostgresPass == "" {
		opts.PostgresPass = "let-me-in"
	}

	// Create network if not provided
	var nw *testcontainers.DockerNetwork
	var err error
	if opts.Network != nil {
		nw = opts.Network
		logger.Debug("Using provided Docker network")
	} else {
		// Create a new network for the Graph Node cluster
		networkName := fmt.Sprintf("graph-node-network-%d", time.Now().UnixNano())
		nw, err = network.New(ctx,
			network.WithAttachable(),
			network.WithLabels(map[string]string{"testbed": "graph-node"}),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create network: %w", err)
		}
		logger.Debug("Created Docker network", "name", networkName)
	}

	// Generate unique names for all containers to avoid conflicts
	timestamp := time.Now().UnixNano()
	postgresName := fmt.Sprintf("postgres-graph-test-%d", timestamp)
	ipfsName := fmt.Sprintf("ipfs-graph-test-%d", timestamp)
	graphNodeName := fmt.Sprintf("graph-node-test-%d", timestamp)

	// Start PostgreSQL first
	logger.Debug("Starting PostgreSQL container", "name", postgresName)
	postgres, err := startPostgres(ctx, opts, nw, postgresName, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres: %w", err)
	}

	// Start IPFS (optional, Graph Node can use external IPFS)
	var ipfs testcontainers.Container
	ipfsEndpoint := opts.IPFSEndpoint
	if ipfsEndpoint == "" {
		logger.Debug("Starting IPFS container", "name", ipfsName)
		ipfs, err = startIPFS(ctx, opts, nw, ipfsName, logger)
		if err != nil {
			return nil, fmt.Errorf("failed to start ipfs: %w", err)
		}
		// Use container name for internal network communication
		ipfsEndpoint = fmt.Sprintf("http://%s:5001", ipfsName)
	}

	// Start Graph Node
	logger.Debug("Starting Graph Node container", "name", graphNodeName)
	graphNode, err := startGraphNode(ctx, opts, nw, ipfsEndpoint, opts.EthereumRPC, graphNodeName, postgresName, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to start graph node: %w", err)
	}

	// Get Graph Node URLs
	host, err := graphNode.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get graph node host: %w", err)
	}

	httpPort, err := graphNode.MappedPort(ctx, "8000")
	if err != nil {
		return nil, fmt.Errorf("failed to get graph node http port: %w", err)
	}

	wsPort, err := graphNode.MappedPort(ctx, "8001")
	if err != nil {
		return nil, fmt.Errorf("failed to get graph node ws port: %w", err)
	}

	adminPort, err := graphNode.MappedPort(ctx, "8020")
	if err != nil {
		return nil, fmt.Errorf("failed to get graph node admin port: %w", err)
	}

	httpURL := fmt.Sprintf("http://%s:%s", host, httpPort.Port())
	wsURL := fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
	adminURL := fmt.Sprintf("http://%s:%s", host, adminPort.Port())

	logger.Info("Graph Node cluster started successfully", "httpURL", httpURL, "wsURL", wsURL, "adminURL", adminURL)

	return &GraphNodeContainer{
		graphNode: graphNode,
		postgres:  postgres,
		ipfs:      ipfs,
		network:   nw,
		httpURL:   httpURL,
		wsURL:     wsURL,
		adminURL:  adminURL,
		logger:    logger,
	}, nil
}

// HTTPURL returns the Graph Node HTTP endpoint
func (g *GraphNodeContainer) HTTPURL() string {
	return g.httpURL
}

// WebSocketURL returns the Graph Node WebSocket endpoint
func (g *GraphNodeContainer) WebSocketURL() string {
	return g.wsURL
}

// AdminURL returns the Graph Node admin endpoint for deployments
func (g *GraphNodeContainer) AdminURL() string {
	return g.adminURL
}

// Terminate stops and removes all containers
func (g *GraphNodeContainer) Terminate(ctx context.Context) error {
	if g == nil {
		return nil
	}

	g.logger.Info("Terminating Graph Node cluster")
	var errs []error

	if g.graphNode != nil {
		g.logger.Debug("Terminating Graph Node container")
		if err := g.graphNode.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate graph node: %w", err))
		}
	}

	if g.ipfs != nil {
		g.logger.Debug("Terminating IPFS container")
		if err := g.ipfs.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate ipfs: %w", err))
		}
	}

	if g.postgres != nil {
		g.logger.Debug("Terminating PostgreSQL container")
		if err := g.postgres.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate postgres: %w", err))
		}
	}

	// Only remove network if we created it (not provided externally)
	if g.network != nil {
		g.logger.Debug("Removing Docker network")
		if err := g.network.Remove(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to remove network: %w", err))
		}
	}

	if len(errs) > 0 {
		g.logger.Error("Errors terminating Graph Node cluster", "errors", errs)
		return fmt.Errorf("errors terminating containers: %v", errs)
	}

	g.logger.Debug("Graph Node cluster terminated successfully")
	return nil
}

// GetPostgres returns the PostgreSQL container for external access
func (g *GraphNodeContainer) GetPostgres() testcontainers.Container {
	return g.postgres
}

// GetIPFS returns the IPFS container for external access
func (g *GraphNodeContainer) GetIPFS() testcontainers.Container {
	return g.ipfs
}

// GetGraphNode returns the Graph Node container for external access
func (g *GraphNodeContainer) GetGraphNode() testcontainers.Container {
	return g.graphNode
}

// IPFSURL returns the IPFS API endpoint URL
func (g *GraphNodeContainer) IPFSURL(ctx context.Context) (string, error) {
	if g.ipfs == nil {
		return "", fmt.Errorf("IPFS container not available")
	}

	host, err := g.ipfs.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get IPFS host: %w", err)
	}

	port, err := g.ipfs.MappedPort(ctx, "5001")
	if err != nil {
		return "", fmt.Errorf("failed to get IPFS port: %w", err)
	}

	return fmt.Sprintf("http://%s:%s", host, port.Port()), nil
}

// PostgresURL returns the PostgreSQL connection URL
func (g *GraphNodeContainer) PostgresURL(ctx context.Context) (string, error) {
	if g.postgres == nil {
		return "", fmt.Errorf("PostgreSQL container not available")
	}

	host, err := g.postgres.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get PostgreSQL host: %w", err)
	}

	port, err := g.postgres.MappedPort(ctx, "5432")
	if err != nil {
		return "", fmt.Errorf("failed to get PostgreSQL port: %w", err)
	}

	// Return connection string format
	return fmt.Sprintf("postgresql://graph-node:let-me-in@%s:%s/graph-node", host, port.Port()), nil
}

// startPostgres creates and starts a PostgreSQL container
func startPostgres(
	ctx context.Context, opts GraphNodeOptions, nw *testcontainers.DockerNetwork, containerName string, logger logging.Logger,
) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        PostgresImage,
		ExposedPorts: []string{PostgresPort},
		Env: map[string]string{
			"POSTGRES_DB":          opts.PostgresDB,
			"POSTGRES_USER":        opts.PostgresUser,
			"POSTGRES_PASSWORD":    opts.PostgresPass,
			"POSTGRES_INITDB_ARGS": "--locale=C --encoding=UTF8",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(60 * time.Second),
		Name:       containerName,
		Networks:   []string{nw.Name},
		NetworkAliases: map[string][]string{
			nw.Name: {containerName, "postgres"},
		},
	}

	genericReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           newTestcontainersLogger(logger),
	}

	return testcontainers.GenericContainer(ctx, genericReq)
}

// startIPFS creates and starts an IPFS container
func startIPFS(ctx context.Context, opts GraphNodeOptions, nw *testcontainers.DockerNetwork, containerName string, logger logging.Logger) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        IPFSImage,
		ExposedPorts: []string{IPFSAPIPort, IPFSGatewayPort},
		WaitingFor:   wait.ForListeningPort("5001/tcp").WithStartupTimeout(60 * time.Second),
		Name:         containerName,
		Networks:     []string{nw.Name},
		NetworkAliases: map[string][]string{
			nw.Name: {containerName, "ipfs"},
		},
	}

	// Add host port bindings if requested
	if opts.ExposeHostPort {
		ipfsPort := opts.HostIPFSPort
		if ipfsPort == "" {
			ipfsPort = "5001"
		}

		req.HostConfigModifier = func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				"5001/tcp": []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: ipfsPort},
				},
			}
		}
	}

	genericReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           newTestcontainersLogger(logger),
	}

	return testcontainers.GenericContainer(ctx, genericReq)
}

// startGraphNode creates and starts a Graph Node container
func startGraphNode(
	ctx context.Context,
	opts GraphNodeOptions,
	nw *testcontainers.DockerNetwork,
	ipfsEndpoint, ethereumRPC, containerName, postgresName string,
	logger logging.Logger,
) (testcontainers.Container, error) {
	// Construct postgres connection string
	postgresURL := fmt.Sprintf("postgresql://%s:%s@%s:5432/%s", opts.PostgresUser, opts.PostgresPass, postgresName, opts.PostgresDB)

	req := testcontainers.ContainerRequest{
		Image: GraphNodeImage,
		ExposedPorts: []string{
			GraphNodeHTTPPort,
			GraphNodeWSPort,
			GraphNodeJSONPort,
			GraphNodeIndexPort,
			GraphNodeMetricsPort,
		},
		Env: map[string]string{
			"postgres_host": postgresName,
			"postgres_user": opts.PostgresUser,
			"postgres_pass": opts.PostgresPass,
			"postgres_db":   opts.PostgresDB,
			"postgres_port": "5432",
			"ipfs":          ipfsEndpoint,
			"ethereum":      "devnet:" + ethereumRPC,
			"GRAPH_LOG":     "debug",
			"RUST_LOG":      "info",
			// Alternative postgres configuration method
			"DATABASE_URL": postgresURL,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("8000/tcp"),
			wait.ForLog("Starting GraphQL HTTP server").WithStartupTimeout(90*time.Second),
		),
		Name:     containerName,
		Networks: []string{nw.Name},
		NetworkAliases: map[string][]string{
			nw.Name: {containerName, "graph-node"},
		},
	}

	// Add host port bindings if requested
	if opts.ExposeHostPort {
		httpPort := opts.HostHTTPPort
		if httpPort == "" {
			httpPort = "8000"
		}
		wsPort := opts.HostWSPort
		if wsPort == "" {
			wsPort = "8001"
		}
		adminPort := opts.HostAdminPort
		if adminPort == "" {
			adminPort = "8020"
		}

		req.HostConfigModifier = func(hc *container.HostConfig) {
			hc.PortBindings = nat.PortMap{
				"8000/tcp": []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: httpPort},
				},
				"8001/tcp": []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: wsPort},
				},
				"8020/tcp": []nat.PortBinding{
					{HostIP: "0.0.0.0", HostPort: adminPort},
				},
			}
		}
	}

	genericReq := testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           newTestcontainersLogger(logger),
	}

	return testcontainers.GenericContainer(ctx, genericReq)
}

