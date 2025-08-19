package containers

import (
	"context"
	"fmt"
	"time"

	"github.com/testcontainers/testcontainers-go"
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

// GraphNodeContainer manages a Graph Node cluster with PostgreSQL and IPFS
type GraphNodeContainer struct {
	graphNode testcontainers.Container
	postgres  testcontainers.Container
	ipfs      testcontainers.Container
	network   *testcontainers.DockerNetwork
	config   GraphNodeConfig
	httpURL  string
	wsURL    string
	adminURL string
}

// NewGraphNodeContainer creates and starts a complete Graph Node setup
// Note: This function assumes a shared network will be provided
func NewGraphNodeContainer(
	ctx context.Context, config GraphNodeConfig, ethereumRPC string, networkName string,
) (*GraphNodeContainer, error) {
	if networkName == "" {
		return nil, fmt.Errorf("networkName is required - GraphNode containers must use a shared network")
	}
	return NewGraphNodeContainerWithNetwork(ctx, config, ethereumRPC, networkName)
}

// NewGraphNodeContainerWithNetwork creates and starts a complete Graph Node setup in a specific network
func NewGraphNodeContainerWithNetwork(
	ctx context.Context, config GraphNodeConfig, ethereumRPC string, networkName string,
) (*GraphNodeContainer, error) {
	if !config.Enabled {
		return nil, fmt.Errorf("graph node container is disabled in config")
	}

	if networkName == "" {
		return nil, fmt.Errorf("networkName is required - GraphNode containers must use a shared network")
	}

	// Generate unique names for all containers to avoid conflicts
	timestamp := time.Now().UnixNano()
	postgresName := fmt.Sprintf("postgres-graph-test-%d", timestamp)
	ipfsName := fmt.Sprintf("ipfs-graph-test-%d", timestamp)
	graphNodeName := fmt.Sprintf("graph-node-test-%d", timestamp)

	// Always use the provided shared network - never create our own
	net := &testcontainers.DockerNetwork{Name: networkName}

	// Start PostgreSQL first
	postgres, err := startPostgres(ctx, config, net.Name, postgresName)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres: %w", err)
	}

	// PostgreSQL will be accessible via container name within the network

	// Start IPFS (optional, Graph Node can use external IPFS)
	var ipfs testcontainers.Container
	ipfsEndpoint := config.IPFSEndpoint
	if ipfsEndpoint == "" {
		ipfs, err = startIPFS(ctx, net.Name, ipfsName)
		if err != nil {
			_ = postgres.Terminate(ctx) // Best effort cleanup
			return nil, fmt.Errorf("failed to start ipfs: %w", err)
		}

		// Use container name for internal network communication
		ipfsEndpoint = ipfsName + ":5001"
	}

	// Start Graph Node
	graphNode, err := startGraphNode(ctx, config, net.Name, ipfsEndpoint, ethereumRPC, graphNodeName, postgresName)
	if err != nil {
		_ = postgres.Terminate(ctx) // Best effort cleanup
		if ipfs != nil {
			_ = ipfs.Terminate(ctx) // Best effort cleanup
		}
		return nil, fmt.Errorf("failed to start graph node: %w", err)
	}

	// Get Graph Node URLs
	host, err := graphNode.Host(ctx)
	if err != nil {
		_ = postgres.Terminate(ctx) // Best effort cleanup
		if ipfs != nil {
			_ = ipfs.Terminate(ctx) // Best effort cleanup
		}
		_ = graphNode.Terminate(ctx) // Best effort cleanup
		return nil, fmt.Errorf("failed to get graph node host: %w", err)
	}

	httpPort, err := graphNode.MappedPort(ctx, "8000")
	if err != nil {
		_ = postgres.Terminate(ctx) // Best effort cleanup
		if ipfs != nil {
			_ = ipfs.Terminate(ctx) // Best effort cleanup
		}
		_ = graphNode.Terminate(ctx) // Best effort cleanup
		return nil, fmt.Errorf("failed to get graph node http port: %w", err)
	}

	wsPort, err := graphNode.MappedPort(ctx, "8001")
	if err != nil {
		_ = postgres.Terminate(ctx) // Best effort cleanup
		if ipfs != nil {
			_ = ipfs.Terminate(ctx) // Best effort cleanup
		}
		_ = graphNode.Terminate(ctx) // Best effort cleanup
		return nil, fmt.Errorf("failed to get graph node ws port: %w", err)
	}

	adminPort, err := graphNode.MappedPort(ctx, "8020")
	if err != nil {
		_ = postgres.Terminate(ctx) // Best effort cleanup
		if ipfs != nil {
			_ = ipfs.Terminate(ctx) // Best effort cleanup
		}
		_ = graphNode.Terminate(ctx) // Best effort cleanup
		return nil, fmt.Errorf("failed to get graph node admin port: %w", err)
	}

	httpURL := fmt.Sprintf("http://%s:%s", host, httpPort.Port())
	wsURL := fmt.Sprintf("ws://%s:%s", host, wsPort.Port())
	adminURL := fmt.Sprintf("http://%s:%s", host, adminPort.Port())

	return &GraphNodeContainer{
		graphNode: graphNode,
		postgres:  postgres,
		ipfs:      ipfs,
		network:   net,
		config:    config,
		httpURL:   httpURL,
		wsURL:     wsURL,
		adminURL:  adminURL,
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
	var errs []error

	if g.graphNode != nil {
		if err := g.graphNode.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate graph node: %w", err))
		}
	}

	if g.ipfs != nil {
		if err := g.ipfs.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate ipfs: %w", err))
		}
	}

	if g.postgres != nil {
		if err := g.postgres.Terminate(ctx); err != nil {
			errs = append(errs, fmt.Errorf("failed to terminate postgres: %w", err))
		}
	}

	// Note: We don't remove the network since it's a shared network managed by InfraManager

	if len(errs) > 0 {
		return fmt.Errorf("errors terminating containers: %v", errs)
	}

	return nil
}

// startPostgres creates and starts a PostgreSQL container
func startPostgres(ctx context.Context, config GraphNodeConfig, networkName, containerName string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        PostgresImage,
		ExposedPorts: []string{PostgresPort},
		Env: map[string]string{
			"POSTGRES_DB":          config.PostgresDB,
			"POSTGRES_USER":        config.PostgresUser,
			"POSTGRES_PASSWORD":    config.PostgresPass,
			"POSTGRES_INITDB_ARGS": "--locale=C --encoding=UTF8",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		Name:       containerName,
		Networks:   []string{networkName},
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

// startIPFS creates and starts an IPFS container
func startIPFS(ctx context.Context, networkName, containerName string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        IPFSImage,
		ExposedPorts: []string{IPFSAPIPort, IPFSGatewayPort},
		WaitingFor:   wait.ForListeningPort("5001/tcp"),
		Name:         containerName,
		Networks:     []string{networkName},
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

// startGraphNode creates and starts a Graph Node container
func startGraphNode(
	ctx context.Context,
	config GraphNodeConfig,
	networkName, ipfsEndpoint, ethereumRPC, containerName, postgresName string,
) (testcontainers.Container, error) {
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
			"postgres_user": config.PostgresUser,
			"postgres_pass": config.PostgresPass,
			"postgres_db":   config.PostgresDB,
			"postgres_port": "5432",
			"ipfs":          ipfsEndpoint,
			"ethereum":      "devnet:" + ethereumRPC,
			"GRAPH_LOG":     "debug",
			"RUST_LOG":      "info",
		},
		WaitingFor: wait.ForListeningPort("8000/tcp"),
		Name:       containerName,
		Networks:   []string{networkName},
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

// DeploySubgraph deploys a subgraph to the Graph Node
func (g *GraphNodeContainer) DeploySubgraph(ctx context.Context, subgraphName, manifestPath string) error {
	// This would typically use the Graph CLI or direct API calls
	// Implementation depends on how subgraph deployment should be handled
	return fmt.Errorf("subgraph deployment not yet implemented")
}

// WaitForReady waits for the Graph Node to be ready
func (g *GraphNodeContainer) WaitForReady(ctx context.Context) error {
	// Additional readiness checks beyond the wait strategy
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
