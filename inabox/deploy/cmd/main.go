package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Layr-Labs/eigenda/inabox/deploy"
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
			return localstack(ctx)
		case expCmdName:
			config.DeployExperiment()
		case generateEnvCmdName:
			config.GenerateAllVariables()
		case allCmdName:
			return all(ctx, config)
		}

		return nil

	}

}

func chainInfra(ctx *cli.Context, config *deploy.Config) error {

	config.StartAnvil()

	if deployer, ok := config.GetDeployer(config.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
		fmt.Println("Starting graph node")
		config.StartGraphNode()
	}

	return nil

}

func localstack(ctx *cli.Context) error {

	pool, _, err := deploy.StartDockertestWithLocalstackContainer(ctx.String(localstackFlagName))
	if err != nil {
		return err
	}

	if ctx.Bool(deployResourcesFlagName) {
		return deploy.DeployResources(pool, ctx.String(localstackFlagName), metadataTableName, bucketTableName, metadataTableNameV2)
	}

	return nil
}

func all(ctx *cli.Context, config *deploy.Config) error {

	err := chainInfra(ctx, config)
	if err != nil {
		return err
	}

	err = localstack(ctx)
	if err != nil {
		return err
	}

	config.DeployExperiment()

	return nil

}
