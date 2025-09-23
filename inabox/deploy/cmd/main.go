package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/test"
	"github.com/Layr-Labs/eigenda/test/testbed"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/urfave/cli/v2"
)

var (
	testNameFlagName        = "testname"
	rootPathFlagName        = "root-path"
	localstackFlagName      = "localstack-port"
	deployResourcesFlagName = "deploy-resources"

	metadataTableName   = "test-BlobMetadata"
	bucketTableName     = "test-BucketStore"
	metadataTableNameV2 = "test-BlobMetadata-v2"

	chainCmdName       = "chain"
	localstackCmdName  = "localstack"
	expCmdName         = "exp"
	generateEnvCmdName = "env"
	allCmdName         = "all"
	eigendaCmdName     = "eigenda"

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
				Name:  localstackFlagName,
				Value: "",
				Usage: "path to the config file",
			},
			&cli.StringFlag{
				Name:  deployResourcesFlagName,
				Value: "",
				Usage: "whether to deploy localstack resources",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   chainCmdName,
				Usage:  "deploy the chain infrastructure (anvil, graph) for the inabox test",
				Action: getRunner(chainCmdName),
			},
			{
				Name:   localstackCmdName,
				Usage:  "deploy localstack and create the AWS resources needed for the inabox test",
				Action: getRunner(localstackCmdName),
			},
			{
				Name:   expCmdName,
				Usage:  "deploy the contracts and create configurations for all EigenDA components",
				Action: getRunner(expCmdName),
			},
			{
				Name:   generateEnvCmdName,
				Usage:  "generate the environment variables for the inabox test",
				Action: getRunner(generateEnvCmdName),
			},
			{
				Name:   eigendaCmdName,
				Usage:  "deploy EigenDA infra with churner via testbed and other components via StartBinaries",
				Action: getRunner(eigendaCmdName),
			},
			{
				Name:   allCmdName,
				Usage:  "deploy all infra, resources, contracts",
				Action: getRunner(allCmdName),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getRunner(command string) func(ctx *cli.Context) error {

	return func(ctx *cli.Context) error {

		var config *deploy.Config
		if command != localstackCmdName {
			rootPath, err := filepath.Abs(ctx.String(rootPathFlagName))
			if err != nil {
				return err
			}
			testname := ctx.String(testNameFlagName)
			if testname == "" {
				testname, err = deploy.GetLatestTestDirectory(rootPath)
				if err != nil {
					return err
				}
			}
			config = deploy.NewTestConfig(testname, rootPath)
		}

		switch command {
		case chainCmdName:
			return chainInfra(ctx, config, nil)
		case localstackCmdName:
			return localstack(ctx, config, nil)
		case expCmdName:
			if err := config.DeployExperiment(); err != nil {
				return fmt.Errorf("failed to deploy experiment: %w", err)
			}
		case generateEnvCmdName:
			config.GenerateAllVariables()
		case eigendaCmdName:
			return eigendaInfra(ctx, config, nil)
		case allCmdName:
			return all(ctx, config)
		}

		return nil

	}

}

func chainInfra(ctx *cli.Context, config *deploy.Config, dockerNetwork *testcontainers.DockerNetwork) error {
	// Disable Ryuk since we likely want to run the test for a long time
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}

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
		return fmt.Errorf("failed to create docker network: %w", err)
	}
	logger.Info("Created Docker network", "name", dockerNetwork.Name)

	anvilC, err := testbed.NewAnvilContainerWithOptions(ctx.Context, testbed.AnvilOptions{
		ExposeHostPort: true,
		HostPort:       "8545",
		Logger:         logger,
		Network:        dockerNetwork,
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
			return fmt.Errorf("failed to start graph node: %w", err)
		}
	}

	return nil

}

func localstack(ctx *cli.Context, config *deploy.Config, dockerNetwork *testcontainers.DockerNetwork) error {
	// Disable Ryuk since we likely want to run the test for a long time
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}

	// Create a shared Docker network for all containers if not provided
	var err error
	if dockerNetwork == nil {
		dockerNetwork, err = network.New(context.Background(),
			network.WithDriver("bridge"),
			network.WithAttachable())
		if err != nil {
			return fmt.Errorf("failed to create docker network: %w", err)
		}
		logger.Info("Created Docker network", "name", dockerNetwork.Name)
	}

	context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err = testbed.NewLocalStackContainerWithOptions(context, testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       ctx.String(localstackFlagName),
		Services:       []string{"s3", "dynamodb", "kms"},
		Logger:         logger,
		Network:        dockerNetwork,
	})
	if err != nil {
		return fmt.Errorf("failed to start localstack container: %w", err)
	}

	if ctx.Bool(deployResourcesFlagName) {
		deployConfig := testbed.DeployResourcesConfig{
			LocalStackEndpoint:  fmt.Sprintf("http://%s:%s", "0.0.0.0", ctx.String(localstackFlagName)),
			MetadataTableName:   metadataTableName,
			BucketTableName:     bucketTableName,
			V2MetadataTableName: metadataTableNameV2,
		}
		if err := testbed.DeployResources(context, deployConfig); err != nil {
			return fmt.Errorf("failed to deploy resources: %w", err)
		}
	}

	return nil
}

func eigendaInfra(ctx *cli.Context, config *deploy.Config, dockerNetwork *testcontainers.DockerNetwork) error {
	// Disable Ryuk since we likely want to run the test for a long time
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}

	fmt.Println("Starting EigenDA infrastructure with testbed churner...")

	// Get deployer configuration for churner
	deployer, ok := config.GetDeployer(config.EigenDA.Deployer)
	if !ok {
		return fmt.Errorf("deployer improperly configured")
	}

	// Start churner using testbed
	churnerConfig := testbed.DefaultChurnerConfig()
	churnerConfig.Enabled = true
	churnerConfig.ChainRPC = "http://anvil:8545" // Use the anvil container in the same networks
	churnerConfig.PrivateKey = strings.TrimPrefix(config.Pks.EcdsaMap[deployer.Name].PrivateKey, "0x")
	churnerConfig.ExposeHostPort = true
	churnerConfig.HostPort = "32002"
	churnerConfig.LogLevel = "debug"

	// Set contract addresses if available from deployment
	if config.EigenDA.ServiceManager != "" {
		churnerConfig.ServiceManager = config.EigenDA.ServiceManager
	}
	if config.EigenDA.OperatorStateRetriever != "" {
		churnerConfig.OperatorStateRetriever = config.EigenDA.OperatorStateRetriever
	}

	// Set graph URL if graph node is configured
	if deployer.DeploySubgraphs {
		churnerConfig.GraphURL = "http://graph-node:8000/subgraphs/name/Layr-Labs/eigenda-operator-state"
	}

	churnerCtx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Create a shared Docker network if not provided
	if dockerNetwork == nil {
		var err error
		dockerNetwork, err = network.New(context.Background(),
			network.WithDriver("bridge"),
			network.WithAttachable())
		if err != nil {
			return fmt.Errorf("failed to create docker network: %w", err)
		}
		logger.Info("Created Docker network", "name", dockerNetwork.Name)
	}

	churnerContainer, err := testbed.NewChurnerContainerWithNetwork(churnerCtx, churnerConfig, dockerNetwork)
	if err != nil {
		return fmt.Errorf("failed to start churner container: %w", err)
	}

	if churnerContainer != nil {
		fmt.Printf("Churner started successfully via testbed at %s\n", churnerContainer.URL())
	}

	// Start the rest of EigenDA components using StartBinaries
	fmt.Println("Starting remaining EigenDA components via StartBinaries...")
	config.StartBinaries()

	return nil
}

func all(ctx *cli.Context, config *deploy.Config) error {
	// Disable Ryuk since we likely want to run the test for a long time
	// This needs to run before ANY testcontainer library call is made.
	// Even creating a network will spin up a Ryuk container.
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}

	err := chainInfra(ctx, config)
	if err != nil {
		return fmt.Errorf("deploy chain infra: %w", err)
	}

	err = chainInfra(ctx, config, dockerNetwork)
	if err != nil {
		return fmt.Errorf("deploy localstack: %w", err)
	}

	err = config.DeployExperiment()
	if err != nil {
		return fmt.Errorf("deploy experiment: %w", err)
	}

	return nil

}
