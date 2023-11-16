package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Layr-Labs/eigenda/inabox/strategies/containers/config"
	genenv "github.com/Layr-Labs/eigenda/inabox/strategies/containers/gen-env"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
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
				Usage:   "name of the test to run (in `inabox/strategies/containers/testdata`)",
				EnvVars: []string{"EIGENDA_TESTDATA_PATH"},
				Value:   "",
			},
			&cli.StringFlag{
				Name:  rootPathFlagName,
				Usage: "path to the root of repo",
				Value: "../../../",
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

	configsDir := filepath.Join(rootPath, "inabox", "strategies", "containers", "testdata")
	if testname == "" {
		files, err := os.ReadDir(configsDir)
		if err != nil {
			panic(err)
			// return err
		}
		if len(files) == 0 {
			return errors.New("no default experiment available")
		}
		testname = files[len(files)-1].Name()
	}

	lock := config.NewConfigLock(testname, rootPath)

	// Create a new experiment and deploy the contracts
	if err != nil {
		log.Panicf("could not load private keys: %v", err)
	}

	fmt.Println("Generating service config variables")
	genenv.GenerateDockerCompose(lock)

	// Write config.lock
	bz, err := yaml.Marshal(lock)
	if err != nil {
		log.Panicf("Yaml serialization of config.lock failed: %v", err)
	}
	config.WriteFile(filepath.Join(configsDir, testname, "config.lock"), bz)

	return nil
}
