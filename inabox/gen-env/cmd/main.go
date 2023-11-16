package main

import (
	"log"
	"os"
	"path/filepath"

	genenv "github.com/Layr-Labs/eigenda/inabox/gen-env"
	"github.com/urfave/cli/v2"
)

var (
	testNameFlagName = "testname"
	rootPathFlagName = "root-path"
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

	testName := ctx.String(testNameFlagName)
	if testName == "" {
		testName = genenv.GetLatestTestDirectory(rootPath)
	}

	lock := genenv.GenerateConfigLock(rootPath, testName)
	genenv.GenerateDockerCompose(lock)
	genenv.CompileDockerCompose(rootPath, testName)

	return nil
}
