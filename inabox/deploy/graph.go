package deploy

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type GraphResources struct {
	Pool             *dockertest.Pool
	PostgresResource *dockertest.Resource
	IpfsResource     *dockertest.Resource
	GraphResource    *dockertest.Resource
	NetworkResource  *docker.Network
}

// Shared configuration constants for postgres
const (
	postgresUser     = "graph-node"
	postgresPassword = "let-me-in"
	postgresDB       = "graph-node"
)

// This function starts the necessary Docker containers for Graph Node services:
// Postgres, IPFS, and Graph Node itself.
// It mimics https://github.com/graphprotocol/graph-node/blob/97992cb852c2e63c7013bbe95f98dc99c1beddce/docker/docker-compose.yml
func StartDockertestWithGraphServices() (*GraphResources, error) {
	fmt.Println("Starting Graph Node services")
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("could not construct pool: %w", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		return nil, fmt.Errorf("could not connect to Docker: %w", err)
	}

	// Create a docker network so containers can communicate by hostname (like docker-compose does)
	network, err := pool.Client.CreateNetwork(docker.CreateNetworkOptions{
		Name: "graph-network",
	})
	if err != nil {
		return nil, fmt.Errorf("could not create docker network: %w", err)
	}

	resources := &GraphResources{Pool: pool, NetworkResource: network}

	// Start Postgres first
	fmt.Println("Starting Postgres container")
	postgresResource, err := startPostgresContainer(pool, network.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to start postgres: %w", err)
	}
	resources.PostgresResource = postgresResource

	// Start IPFS second
	fmt.Println("Starting IPFS container")
	ipfsResource, err := startIpfsContainer(pool, network.ID)
	if err != nil {
		PurgeDockertestGraphResources(resources)
		return nil, fmt.Errorf("failed to start ipfs: %w", err)
	}
	resources.IpfsResource = ipfsResource

	// Start Graph Node last (depends on postgres and ipfs)
	fmt.Println("Starting Graph Node container")
	graphResource, err := startGraphNodeContainer(pool, network.ID, "graph-postgres", "graph-ipfs:5001")
	if err != nil {
		PurgeDockertestGraphResources(resources)
		return nil, fmt.Errorf("failed to start graph node: %w", err)
	}
	resources.GraphResource = graphResource

	log.Printf("Graph Node services started successfully! GraphQL endpoint: http://0.0.0.0:8000")
	return resources, nil
}

func startPostgresContainer(pool *dockertest.Pool, networkID string) (*dockertest.Resource, error) {
	runOpts := &dockertest.RunOptions{
		Repository:   "postgres",
		Tag:          "14",
		Name:         "graph-postgres",
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port("5432"): {
				{HostIP: "0.0.0.0", HostPort: "5432"},
			},
		},
		Networks: []*dockertest.Network{{Network: &docker.Network{ID: networkID}}},
		Env: []string{
			"POSTGRES_USER=" + postgresUser,
			"POSTGRES_PASSWORD=" + postgresPassword,
			"POSTGRES_DB=" + postgresDB,
			"POSTGRES_INITDB_ARGS=-E UTF8 --locale=C",
		},
		Cmd: []string{
			"postgres",
			"-cshared_preload_libraries=pg_stat_statements",
			"-cmax_connections=200",
		},
	}

	resource, err := pool.RunWithOptions(runOpts, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, err
	}

	// Health check for Postgres
	pool.MaxWait = 30 * time.Second
	if err := pool.Retry(func() error {
		fmt.Println("Waiting for postgres to start")

		// Use pg_isready to check if postgres is ready
		exec, err := pool.Client.CreateExec(docker.CreateExecOptions{
			Container: resource.Container.ID,
			Cmd:       []string{"pg_isready", "-d", postgresDB, "-U", postgresUser},
		})
		if err != nil {
			return err
		}

		err = pool.Client.StartExec(exec.ID, docker.StartExecOptions{})
		if err != nil {
			return fmt.Errorf("postgres is not ready: %w", err)
		}

		fmt.Println("Postgres is running and responding!")
		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not connect to postgres: %w", err)
	}

	return resource, nil
}

func startIpfsContainer(pool *dockertest.Pool, networkID string) (*dockertest.Resource, error) {
	runOpts := &dockertest.RunOptions{
		Repository:   "ipfs/kubo",
		Tag:          "v0.14.0",
		Name:         "graph-ipfs",
		ExposedPorts: []string{"5001"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port("5001"): {
				{HostIP: "0.0.0.0", HostPort: "5001"},
			},
		},
		Networks: []*dockertest.Network{{Network: &docker.Network{ID: networkID}}},
	}

	resource, err := pool.RunWithOptions(runOpts, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		return nil, err
	}

	// Simple readiness check for IPFS - ensure it's ready to accept connections
	// This mimics the docker-compose "service_started" condition
	pool.MaxWait = 30 * time.Second
	if err := pool.Retry(func() error {
		fmt.Println("Waiting for IPFS to be ready to accept connections")

		client := &http.Client{Timeout: 2 * time.Second}
		// Use /api/v0/id endpoint which should be available and not return 405
		resp, err := client.Get("http://0.0.0.0:5001/api/v0/id")
		if err != nil {
			return fmt.Errorf("IPFS is not accepting connections: %w", err)
		}
		defer core.CloseLogOnError(resp.Body, "ipfs response body", nil)

		// Just ensure IPFS is responding (don't check specific status codes)
		fmt.Println("IPFS is ready and accepting connections!")
		return nil
	}); err != nil {
		return nil, fmt.Errorf("could not verify IPFS is ready: %w", err)
	}

	return resource, nil
}

func startGraphNodeContainer(pool *dockertest.Pool, networkID, postgresHost, ipfsHost string) (*dockertest.Resource, error) {
	runOpts := &dockertest.RunOptions{
		Repository:   "graphprotocol/graph-node",
		Tag:          "latest",
		Name:         "graph-node",
		ExposedPorts: []string{"8000", "8001", "8020", "8030", "8040"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port("8000"): {{HostIP: "0.0.0.0", HostPort: "8000"}},
			docker.Port("8001"): {{HostIP: "0.0.0.0", HostPort: "8001"}},
			docker.Port("8020"): {{HostIP: "0.0.0.0", HostPort: "8020"}},
			docker.Port("8030"): {{HostIP: "0.0.0.0", HostPort: "8030"}},
			docker.Port("8040"): {{HostIP: "0.0.0.0", HostPort: "8040"}},
		},
		Networks: []*dockertest.Network{{Network: &docker.Network{ID: networkID}}},
		Env: []string{
			"postgres_host=" + postgresHost,
			"postgres_user=" + postgresUser,
			"postgres_pass=" + postgresPassword,
			"postgres_db=" + postgresDB,
			"ipfs=" + ipfsHost,
			// TODO: we should run anvil in the same network and use its hostname instead of host.docker.internal.
			"ethereum=devnet:http://host.docker.internal:8545",
			"GRAPH_LOG=info",
		},
	}

	resource, err := pool.RunWithOptions(runOpts, func(config *docker.HostConfig) {
		// config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		// Add extra hosts for host.docker.internal
		config.ExtraHosts = []string{"host.docker.internal:host-gateway"}
	})
	if err != nil {
		return nil, err
	}

	// TODO: add a health check for graph-node.
	// Currently no endpoint is available, see https://github.com/graphprotocol/graph-node/issues/5484
	fmt.Println("Graph Node container started successfully")

	return resource, nil
}

func PurgeDockertestGraphResources(resources *GraphResources) {
	fmt.Println("Stopping Graph Node Dockertest resources")

	if resources == nil {
		return
	}

	// Stop in reverse order of startup
	if resources.GraphResource != nil {
		fmt.Println("Expiring graph node docker resource")
		if err := resources.GraphResource.Expire(1); err != nil {
			log.Printf("Could not expire graph node resource: %s", err)
		}
		if resources.Pool != nil {
			fmt.Println("Purging graph node docker resource")
			if err := resources.Pool.Purge(resources.GraphResource); err != nil {
				log.Printf("Could not purge graph node resource: %s", err)
			}
		}
	}

	if resources.IpfsResource != nil {
		fmt.Println("Expiring ipfs docker resource")
		if err := resources.IpfsResource.Expire(1); err != nil {
			log.Printf("Could not expire ipfs resource: %s", err)
		}
		if resources.Pool != nil {
			fmt.Println("Purging ipfs docker resource")
			if err := resources.Pool.Purge(resources.IpfsResource); err != nil {
				log.Printf("Could not purge ipfs resource: %s", err)
			}
		}
	}

	if resources.PostgresResource != nil {
		fmt.Println("Expiring postgres docker resource")
		if err := resources.PostgresResource.Expire(1); err != nil {
			log.Printf("Could not expire postgres resource: %s", err)
		}
		if resources.Pool != nil {
			fmt.Println("Purging postgres docker resource")
			if err := resources.Pool.Purge(resources.PostgresResource); err != nil {
				log.Printf("Could not purge postgres resource: %s", err)
			}
		}
	}

	// Remove the network last, after all containers are cleaned up
	if resources.NetworkResource != nil && resources.Pool != nil {
		fmt.Println("Removing docker network")
		if err := resources.Pool.Client.RemoveNetwork(resources.NetworkResource.ID); err != nil {
			log.Printf("Could not remove docker network: %s", err)
		}
	}
}
