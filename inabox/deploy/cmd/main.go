package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/urfave/cli/v2"
)

var (
	testNameFlagName       = "testname"
	rootPathFlagName       = "root-path"
	localstackPortFlagName = "localstack-port"

	metadataTableName   = "test-BlobMetadata"
	bucketTableName     = "test-BucketStore"
	metadataTableNameV2 = "test-BlobMetadata-v2"

	logger = test.GetLogger()
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    testNameFlagName,
				Usage:   "name of the test to run (in `inabox/testdata`)",
				EnvVars: []string{"EIGENDA_TESTDATA_PATH"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:  rootPathFlagName,
				Usage: "path to the root of repo",
				Value: "../",
			},
			&cli.StringFlag{
				Name:  localstackPortFlagName,
				Value: "",
				Usage: "path to the config file",
			},
		},
		Action:      DeployAll,
		Description: "Deploys all infra, resources, and contracts needed to spin up a local EigenDA inabox devnet.",
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func DeployAll(ctx *cli.Context) error {
	config, err := readTestConfig(ctx)
	if err != nil {
		return fmt.Errorf("get test config: %w", err)
	}

	// Disable Ryuk since we likely want to run the test for a long time
	// This will prevent testcontainer's GC container from starting,
	// and will hence let the containers run indefinitely.
	// They can be stopped manually using `make stop-infra`.
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}

	err = startChainInfra(ctx, config)
	if err != nil {
		return fmt.Errorf("start chain infra: %w", err)
	}

	err = startLocalstack(ctx, config)
	if err != nil {
		return fmt.Errorf("start localstack: %w", err)
	}

	err = config.DeployExperiment()
	if err != nil {
		return fmt.Errorf("deploy experiment: %w", err)
	}

	logger.Info("Deployment complete. You can now run `make run-e2e` to run the e2e tests.")

	return nil

}

func readTestConfig(ctx *cli.Context) (*deploy.Config, error) {
	rootPath, err := filepath.Abs(ctx.String(rootPathFlagName))
	if err != nil {
		return nil, fmt.Errorf("get absolute root path: %w", err)
	}
	testname := ctx.String(testNameFlagName)
	if testname == "" {
		testname, err = deploy.GetLatestTestDirectory(rootPath)
		if err != nil {
			return nil, fmt.Errorf("get latest test directory: %w", err)
		}
	}
	config := deploy.ReadTestConfig(testname, rootPath)
	return config, nil
}

// Spins up an anvil chain and a graph node (if DeploySubgraphs=true)
func startChainInfra(ctx *cli.Context, config *deploy.Config) error {

	_, err := testbed.NewAnvilContainerWithOptions(ctx.Context, testbed.AnvilOptions{
		ExposeHostPort: true,
		HostPort:       "8545",
		Logger:         logger,
	})
	if err != nil {
		return fmt.Errorf("failed to start anvil container: %w", err)
	}

	if deployer, ok := config.GetDeployer(config.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
		fmt.Println("Starting graph node")
		_, err := testbed.NewGraphNodeContainerWithOptions(ctx.Context, testbed.GraphNodeOptions{
			PostgresDB:     "graph-node",
			PostgresUser:   "graph-node",
			PostgresPass:   "let-me-in",
			EthereumRPC:    "http://localhost:8545",
			ExposeHostPort: true,
			HostHTTPPort:   "8000",
			HostWSPort:     "8001",
			HostAdminPort:  "8020",
			HostIPFSPort:   "5001",
			Logger:         logger,
		})
		if err != nil {
			return fmt.Errorf("failed to start graph node: %w", err)
		}
	}

	return nil

}

func startLocalstack(ctx *cli.Context, config *deploy.Config) error {
	context, cancel := context.WithTimeout(ctx.Context, 30*time.Second)
	defer cancel()

	_, err := testbed.NewLocalStackContainerWithOptions(context, testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       ctx.String(localstackPortFlagName),
		Services:       []string{"s3", "dynamodb", "kms"},
		Logger:         logger,
	})
	if err != nil {
		return fmt.Errorf("failed to start localstack container: %w", err)
	}

	deployConfig := testbed.DeployResourcesConfig{
		LocalStackEndpoint:  fmt.Sprintf("http://%s:%s", "0.0.0.0", ctx.String(localstackPortFlagName)),
		MetadataTableName:   metadataTableName,
		BucketTableName:     bucketTableName,
		V2MetadataTableName: metadataTableNameV2,
	}
	if err := testbed.DeployResources(context, deployConfig); err != nil {
		return fmt.Errorf("failed to deploy resources: %w", err)
	}

	return nil
}
