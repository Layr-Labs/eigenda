package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/tools/traffic"
	"github.com/Layr-Labs/eigenda/tools/traffic/config"
	"github.com/urfave/cli"
)

var (
	version   = ""
	gitCommit = ""
	gitDate   = ""
)

func main() {
	app := cli.NewApp()
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)
	app.Name = "da-traffic-generator"
	app.Usage = "EigenDA Traffic Generator"
	app.Description = "Service for generating traffic to EigenDA disperser"
	app.Flags = config.Flags
	app.Action = trafficGeneratorMain
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}
}

func trafficGeneratorMain(ctx *cli.Context) error {
	generatorConfig, err := config.NewConfig(ctx)
	if err != nil {
		return err
	}

	generator, err := traffic.NewTrafficGeneratorV2(generatorConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create new traffic generator\n%s", err))
	}

	return generator.Start()
}
