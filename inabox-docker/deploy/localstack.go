package deploy

import (
	"fmt"
	"log"

	"github.com/ory/dockertest/v3"
)

// func StartDockertestWithLocalstackContainer(localStackPort string) (*dockertest.Pool, *dockertest.Resource, error) {
// 	fmt.Println("Starting Localstack container")
// 	pool, err := dockertest.NewPool("")
// 	if err != nil {
// 		fmt.Println("Could not construct pool: %w", err)
// 		return nil, nil, err
// 	}

// 	err = pool.Client.Ping()
// 	if err != nil {
// 		fmt.Println("Could not connect to Docker: %w", err)
// 		return nil, nil, err
// 	}

// 	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
// 		Repository:   "localstack/localstack",
// 		Tag:          "latest",
// 		ExposedPorts: []string{localStackPort},
// 		PortBindings: map[docker.Port][]docker.PortBinding{
// 			docker.Port(localStackPort): {
// 				{HostIP: "0.0.0.0", HostPort: localStackPort},
// 			},
// 		},
// 		Env: []string{
// 			fmt.Sprintf("GATEWAY_LISTEN=0.0.0.0:%s", localStackPort),
// 			fmt.Sprintf("LOCALSTACK_HOST=localhost.localstack.cloud:%s", localStackPort),
// 		},
// 	}, func(config *docker.HostConfig) {
// 		// set AutoRemove to true so that stopped container goes away by itself
// 		config.AutoRemove = true
// 		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
// 	})
// 	if err != nil {
// 		fmt.Println("Could not start resource: %w", err)
// 		return nil, nil, err
// 	}

// 	pool.MaxWait = 10 * time.Second
// 	if err := pool.Retry(func() error {
// 		fmt.Println("Waiting for localstack to start")
// 		resp, err := http.Get(fmt.Sprintf("http://0.0.0.0:%s", localStackPort))
// 		if err != nil {
// 			fmt.Println("Server is not running:", err)
// 			return err
// 		}
// 		defer resp.Body.Close()

// 		if resp.StatusCode != http.StatusOK {
// 			fmt.Printf("Server returned non-OK status: %s\n", resp.Status)
// 			return errors.New("non-ok status")

// 		}

// 		fmt.Println("Server is running and responding!")
// 		return nil

// 	}); err != nil {
// 		fmt.Println("Could not connect to localstack:", err)
// 		return nil, nil, err
// 	}

// 	log.Printf("Localstack started successfully! URL: http://0.0.0.0:%s", localStackPort)

// 	return pool, resource, nil
// }

func PurgeDockertestResources(pool *dockertest.Pool, resource *dockertest.Resource) {
	fmt.Println("Stopping Dockertest resources")
	if resource != nil {
		fmt.Println("Expiring docker resource")
		if err := resource.Expire(1); err != nil {
			log.Fatalf("Could not expire resource: %s", err)
		}
	}

	if resource != nil && pool != nil {
		fmt.Println("Purging docker resource")
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}
}
