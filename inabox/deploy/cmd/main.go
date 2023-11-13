package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/ory/dockertest/v3"
	"github.com/urfave/cli/v2"
)

var (
	testNameFlagName   = "testname"
	rootPathFlagName   = "root-path"
	localstackFlagName = "localstack-port"

	metadataTableName = "test-BlobMetadata"
	bucketTableName   = "test-BucketStore"

	infraCmdName     = "infra"
	resourcesCmdName = "resources"
	expCmdName       = "exp"
	allCmdName       = "all"
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
		},
		Commands: []*cli.Command{
			{
				Name:   infraCmdName,
				Usage:  "deploy the infrastructure (anvil, graph, localstack) for the inabox test",
				Action: getRunner(infraCmdName),
			},
			{
				Name:   resourcesCmdName,
				Usage:  "deploy the AWS resources needed for the inabox test",
				Action: getRunner(resourcesCmdName),
			},
			{
				Name:   expCmdName,
				Usage:  "deploy the contracts and create configurations for all EigenDA components",
				Action: getRunner(expCmdName),
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

		config := deploy.NewTestConfig(testname, rootPath)

		switch command {
		case infraCmdName:
			_, _, err = infra(ctx, config)
			return err
		case resourcesCmdName:
			return resources(ctx)
		case expCmdName:
			config.DeployExperiment()
		case allCmdName:
			return all(ctx, config)
		}

		return nil

	}

}

func infra(ctx *cli.Context, config *deploy.Config) (*dockertest.Pool, *dockertest.Resource, error) {

	pool, resources, err := deploy.StartDockertestWithLocalstackContainer(ctx.String(localstackFlagName))
	config.StartAnvil()

	if deployer, ok := config.GetDeployer(config.EigenDA.Deployer); ok && deployer.DeploySubgraphs {
		fmt.Println("Starting graph node")
		config.StartGraphNode()
	}

	return pool, resources, err

}

func resources(ctx *cli.Context) error {
	return deploy.DeployResources(nil, ctx.String(localstackFlagName), metadataTableName, bucketTableName)
}

func all(ctx *cli.Context, config *deploy.Config) error {

	pool, _, err := infra(ctx, config)
	if err != nil {
		return err
	}

	err = deploy.DeployResources(pool, ctx.String(localstackFlagName), metadataTableName, bucketTableName)
	if err != nil {
		return err
	}

	config.DeployExperiment()

	return nil

}
