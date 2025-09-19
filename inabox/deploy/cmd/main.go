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
			return chainInfra(ctx, config)
		case localstackCmdName:
			return localstack(ctx, config)
		case expCmdName:
			if err := config.DeployExperiment(); err != nil {
				return fmt.Errorf("failed to deploy experiment: %w", err)
			}
		case generateEnvCmdName:
			config.GenerateAllVariables()
		case allCmdName:
			return all(ctx, config)
		}

		return nil

	}

}

func chainInfra(ctx *cli.Context, config *deploy.Config) error {
	// Disable Ryuk since we likely want to run the test for a long time
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}

	_, err := testbed.NewAnvilContainerWithOptions(context.Background(), testbed.AnvilOptions{
		ExposeHostPort: true,
		HostPort:       "8545",
		Logger:         logger,
	})
	if err != nil {
		return fmt.Errorf("failed to start anvil container: %w", err)
	}

	if deployer, ok := config.GetDeployer(config.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
		fmt.Println("Starting graph node")
		_, err := testbed.NewGraphNodeContainerWithOptions(context.Background(), testbed.GraphNodeOptions{
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
		// Wait for Graph Node to be ready
		fmt.Println("Waiting for Graph Node to be ready...")
		time.Sleep(10 * time.Second)
	}

	return nil

}

func localstack(ctx *cli.Context, config *deploy.Config) error {
	// Disable Ryuk since we likely want to run the test for a long time
	if err := os.Setenv("TESTCONTAINERS_RYUK_DISABLED", "true"); err != nil {
		return fmt.Errorf("failed to set environment variable: %w", err)
	}

	context, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_, err := testbed.NewLocalStackContainerWithOptions(context, testbed.LocalStackOptions{
		ExposeHostPort: true,
		HostPort:       ctx.String(localstackFlagName),
		Services:       []string{"s3", "dynamodb", "kms"},
		Logger:         logger,
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

func all(ctx *cli.Context, config *deploy.Config) error {

	err := chainInfra(ctx, config)
	if err != nil {
		return err
	}

	err = localstack(ctx, config)
	if err != nil {
		return err
	}

	config.DeployExperiment()

	return nil

}
