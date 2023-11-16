package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/Layr-Labs/eigenda/inabox-docker/deploy"
	"github.com/urfave/cli/v2"
)

var (
	testNameFlagName   = "testname"
	rootPathFlagName   = "root-path"
	localstackFlagName = "localstack-port"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    testNameFlagName,
				Usage:   "name of the test to run (in `inabox-docker/testdata`)",
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
		Action: action,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx *cli.Context) error {
	rootPath, err := filepath.Abs(ctx.String(rootPathFlagName))
	if err != nil {
		return err
	}

	testname := ctx.String(testNameFlagName)

	if testname == "" {
		files, err := os.ReadDir(filepath.Join(rootPath, "inabox-docker", "testdata"))
		if err != nil {
			return err
		}
		if len(files) == 0 {
			return errors.New("no default experiment available")
		}
		testname = files[len(files)-1].Name()
	}

	config := deploy.NewTestConfig(testname, rootPath)

	config.GenerateServiceConfig()

	return nil
}
