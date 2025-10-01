package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
	gcommon "github.com/ethereum/go-ethereum/common"
	"github.com/testcontainers/testcontainers-go/network"
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

	_, err = startChainInfra(ctx, config)
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

	logger.Info("Generating disperser keypair")
	err = config.GenerateDisperserKeypair()
	if err != nil {
		logger.Errorf("could not generate disperser keypair: %v", err)
		panic(err)
	}

	// Create eth client
	ethClient, err := geth.NewMultiHomingClient(geth.EthClientConfig{
		RPCURLs:          []string{config.Deployers[0].RPC},
		PrivateKeyString: config.Pks.EcdsaMap[config.EigenDA.Deployer].PrivateKey[2:],
		NumConfirmations: 0,
		NumRetries:       3,
	}, gcommon.Address{}, logger)
	if err != nil {
		logger.Errorf("could not create eth client for registration: %v", err)
		panic(err)
	}

	logger.Info("Registering disperser keypair on-chain")
	config.PerformDisperserRegistrations(ethClient)

	// Register blob versions
	config.RegisterBlobVersions(ethClient)

	// Register relay URLs
	relayURLs := []string{
		"localhost:32035",
		"localhost:32037",
		"localhost:32039",
		"localhost:32041",
	}
	config.RegisterRelays(ethClient, relayURLs, ethClient.GetAccountAddress())

	logger.Info("Generating variables")
	err = config.GenerateAllVariables("0.0.0.0:34001")
	if err != nil {
		logger.Errorf("could not generate environment variables: %v", err)
		panic(err)
	}

	logger.Info("Deployment complete. You can now run `make start-services` to start the services.")
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
func startChainInfra(ctx *cli.Context, config *deploy.Config) (*testbed.AnvilContainer, error) {
	// Create a shared Docker network for all containers
	// TODO(samlaf): seems like there's no way with testcontainers-go@v0.38 to give this network a name...
	// https://pkg.go.dev/github.com/testcontainers/testcontainers-go@v0.38.0/network#WithNetworkName
	// only returns an option to be passed to container requests... so we would have to use it on the first container
	// we create, which would require changing our testbed package.
	dockerNetwork, err := network.New(ctx.Context,
		network.WithDriver("bridge"),
		network.WithAttachable(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker network: %w", err)
	}
	logger.Info("Created Docker network", "name", dockerNetwork.Name)

	anvilC, err := testbed.NewAnvilContainerWithOptions(ctx.Context, testbed.AnvilOptions{
		ExposeHostPort: true,
		HostPort:       "8545",
		Logger:         logger,
		Network:        dockerNetwork,
		BlockTime:      1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start anvil container: %w", err)
	}

	if deployer, ok := config.GetDeployer(config.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
		fmt.Println("Starting graph node")
		_, err := testbed.NewGraphNodeContainerWithOptions(ctx.Context, testbed.GraphNodeOptions{
			PostgresDB:     "graph-node",
			PostgresUser:   "graph-node",
			PostgresPass:   "let-me-in",
			ExposeHostPort: true,
			HostHTTPPort:   "8000",
			HostWSPort:     "8001",
			HostAdminPort:  "8020",
			HostIPFSPort:   "5001",
			Logger:         logger,
			Network:        dockerNetwork,
			// internal endpoint will work because they are in the same dockerNetwork
			EthereumRPC: anvilC.InternalEndpoint(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to start graph node: %w", err)
		}
	}

	return anvilC, nil

}

func startLocalstack(ctx *cli.Context, config *deploy.Config) error {
	context, cancel := context.WithTimeout(ctx.Context, 30*time.Second)
	defer cancel()

	localstackContainer, err := testbed.NewLocalStackContainerWithOptions(context, testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       ctx.String(localstackPortFlagName),
		Services:       []string{"s3", "dynamodb", "kms"},
		Logger:         logger,
	})
	if err != nil {
		return fmt.Errorf("failed to start localstack container: %w", err)
	}

	deployConfig := testbed.DeployResourcesConfig{
		LocalStackEndpoint:  localstackContainer.Endpoint(),
		MetadataTableName:   metadataTableName,
		BucketTableName:     bucketTableName,
		V2MetadataTableName: metadataTableNameV2,
		AWSConfig:           localstackContainer.GetAWSClientConfig(),
	}
	if err := testbed.DeployResources(context, deployConfig); err != nil {
		return fmt.Errorf("failed to deploy resources: %w", err)
	}

	return nil
}
