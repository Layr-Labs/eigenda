package genenv

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/Layr-Labs/eigenda/inabox/strategies/containers/config"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

func GenAnvilState(lock *config.ConfigLock) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Panicf("Failed creating dockertest pool: %v", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Panicf("Could not connect to Docker: %v", err)
	}

	// Build Docker image from Dockerfile
	resource, err := pool.BuildAndRunWithBuildOptions(&dockertest.BuildOptions{
		Dockerfile: filepath.Join(lock.RootPath, "inabox/strategies/containers/gen-anvil-state/Dockerfile"), // Name of the Dockerfile
		ContextDir: lock.RootPath,                                                                           // Directory containing the Dockerfile
	}, &dockertest.RunOptions{
		Mounts: []string{
			fmt.Sprintf("%v/contracts:/contracts", lock.RootPath),
			fmt.Sprintf("%v:/data", lock.Path),
		}, // Set your volume mounts here
		// Env:    []string{"ENV_VAR=value"},               // Set environment variables if needed
	})
	if err != nil {
		log.Panicf("Could not connect to Docker: %v", err)
	}

	// Clean up after
	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Printf("Could not purge resource: %s", err)
		}
	}()

	// Wait for the container to exit
	statusCh, err := pool.Client.AttachToContainerNonBlocking(docker.AttachToContainerOptions{
		Container:    resource.Container.ID,
		OutputStream: nil,
		ErrorStream:  nil,
		Stderr:       true,
		Stdout:       true,
		Stream:       true,
	})
	if err != nil {
		log.Panicf("Could not connect to Docker: %w", err)
	}

	statusCh.Wait()
}
