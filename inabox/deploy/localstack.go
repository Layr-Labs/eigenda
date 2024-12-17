package deploy

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/store"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	test_utils "github.com/Layr-Labs/eigenda/common/aws/dynamodb/utils"
)

func StartDockertestWithLocalstackContainer(localStackPort string) (*dockertest.Pool, *dockertest.Resource, error) {
	fmt.Println("Starting Localstack container")
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

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository:   "localstack/localstack",
		Tag:          "latest",
		Name:         "localstack-test",
		ExposedPorts: []string{localStackPort},
		PortBindings: map[docker.Port][]docker.PortBinding{
			docker.Port(localStackPort): {
				{HostIP: "0.0.0.0", HostPort: localStackPort},
			},
		},
		Env: []string{
			fmt.Sprintf("GATEWAY_LISTEN=0.0.0.0:%s", localStackPort),
			fmt.Sprintf("LOCALSTACK_HOST=localhost.localstack.cloud:%s", localStackPort),
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		fmt.Println("Could not start resource: %w", err)
		return nil, nil, err
	}

	pool.MaxWait = 10 * time.Second
	if err := pool.Retry(func() error {
		fmt.Println("Waiting for localstack to start")
		resp, err := http.Get(fmt.Sprintf("http://0.0.0.0:%s", localStackPort))
		if err != nil {
			fmt.Println("Server is not running:", err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Server returned non-OK status: %s\n", resp.Status)
			return errors.New("non-ok status")

		}

		fmt.Println("Server is running and responding!")
		return nil

	}); err != nil {
		fmt.Println("Could not connect to localstack:", err)
		return nil, nil, err
	}

	log.Printf("Localstack started successfully! URL: http://0.0.0.0:%s", localStackPort)

	return pool, resource, nil
}

func DeployResources(
	pool *dockertest.Pool,
	localStackPort,
	metadataTableName,
	bucketTableName,
	v2MetadataTableName string,
) error {

	if pool == nil {
		var err error
		pool, err = dockertest.NewPool("")
		if err != nil {
			fmt.Println("Could not construct pool: %w", err)
			return err
		}
	}

	// exponential backoff-retry, because the application in
	// the container might not be ready to accept connections yet
	pool.MaxWait = 10 * time.Second
	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../..")
	changeDirectory(filepath.Join(rootPath, "inabox"))
	if err := pool.Retry(func() error {
		fmt.Println("Creating S3 bucket")
		return execCmd("./create-s3-bucket.sh", []string{}, []string{fmt.Sprintf("AWS_URL=http://0.0.0.0:%s", localStackPort)})
	}); err != nil {
		fmt.Println("Could not connect to docker:", err)
		return err
	}

	cfg := aws.ClientConfig{
		Region:          "us-east-1",
		AccessKey:       "localstack",
		SecretAccessKey: "localstack",
		EndpointURL:     fmt.Sprintf("http://0.0.0.0:%s", localStackPort),
	}
	_, err := test_utils.CreateTable(context.Background(), cfg, metadataTableName, blobstore.GenerateTableSchema(metadataTableName, 10, 10))
	if err != nil {
		return err
	}

	_, err = test_utils.CreateTable(context.Background(), cfg, bucketTableName, store.GenerateTableSchema(10, 10, bucketTableName))
	if err != nil {
		return err
	}

	fmt.Println("Creating v2 tables")
	if v2MetadataTableName != "" {
		// Create v2 metadata table
		_, err = test_utils.CreateTable(context.Background(), cfg, v2MetadataTableName, blobstorev2.GenerateTableSchema(v2MetadataTableName, 10, 10))
		if err != nil {
			return err
		}

		v2PaymentName := "e2e_v2_"
		// create payment related tables
		err = meterer.CreateReservationTable(cfg, v2PaymentName+"reservation")
		if err != nil {
			fmt.Println("err", err)
			return err
		}
		err = meterer.CreateOnDemandTable(cfg, v2PaymentName+"ondemand")
		if err != nil {
			fmt.Println("err", err)
			return err
		}
		err = meterer.CreateGlobalReservationTable(cfg, v2PaymentName+"global_reservation")
		if err != nil {
			fmt.Println("err", err)
			return err
		}
	}

	return err

}

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
