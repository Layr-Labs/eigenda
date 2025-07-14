package deploy

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func StartDockertestWithAnvilContainer(anvilPort string) (*dockertest.Pool, *dockertest.Resource, error) {
	fmt.Println("Starting Anvil container")
	pool, err := dockertest.NewPool("")
	if err != nil {
		fmt.Println("Could not construct pool: %w", err)
		return nil, nil, err
	}

	err = pool.Client.Ping()
	if err != nil {
		fmt.Println("Could not connect to Docker: %w", err)
		return nil, nil, err
	}

	runOpts := &dockertest.RunOptions{
		Repository:   "ghcr.io/foundry-rs/foundry",
		Tag:          "latest",
		Name:         "anvil-inabox",
		ExposedPorts: []string{anvilPort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(anvilPort): {
				{HostIP: "0.0.0.0", HostPort: anvilPort},
			},
		},
		Cmd: []string{"anvil --host 0.0.0.0"},
	}

	fmt.Printf("Running container with cmd: %v\n", runOpts.Cmd)
	
	resource, err := pool.RunWithOptions(runOpts, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		fmt.Println("Could not start resource: %w", err)
		return nil, nil, err
	}

	// Debug: inspect the container configuration
	containerInfo, err := pool.Client.InspectContainer(resource.Container.ID)
	if err == nil {
		fmt.Printf("Container CMD: %v\n", containerInfo.Config.Cmd)
		fmt.Printf("Container Entrypoint: %v\n", containerInfo.Config.Entrypoint)
		fmt.Printf("Container ENV: %v\n", containerInfo.Config.Env)
	}

	pool.MaxWait = 10 * time.Second
	if err := pool.Retry(func() error {
		fmt.Println("Waiting for anvil to start")
		
		// Try to connect to Anvil RPC endpoint
		client := &http.Client{Timeout: 2 * time.Second}
		req, err := http.NewRequest("POST", fmt.Sprintf("http://0.0.0.0:%s", anvilPort), nil)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Anvil is not running:", err)
			return err
		}
		defer core.CloseLogOnError(resp.Body, "anvil response body", nil)

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Anvil returned non-OK status: %s\n", resp.Status)
			return errors.New("non-ok status")
		}

		fmt.Println("Anvil is running and responding!")
		return nil

	}); err != nil {
		fmt.Println("Could not connect to anvil:", err)
		return nil, nil, err
	}

	log.Printf("Anvil started successfully! URL: http://0.0.0.0:%s", anvilPort)

	return pool, resource, nil
}

func PurgeDockertestAnvilResources(pool *dockertest.Pool, resource *dockertest.Resource) {
	fmt.Println("Stopping Anvil Dockertest resources")
	if resource != nil {
		fmt.Println("Expiring anvil docker resource")
		if err := resource.Expire(1); err != nil {
			log.Fatalf("Could not expire resource: %s", err)
		}
	}

	if resource != nil && pool != nil {
		fmt.Println("Purging anvil docker resource")
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}