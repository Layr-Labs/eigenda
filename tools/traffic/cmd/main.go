package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/tools/traffic"
	"github.com/Layr-Labs/eigenda/tools/traffic/flags"
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
	app.Flags = flags.Flags
	app.Action = trafficGeneratorMain
	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}
}

func trafficGeneratorMain(ctx *cli.Context) error {
	config := traffic.NewConfig(ctx)
	generator, err := traffic.NewTrafficGenerator(config)
	if err != nil {
		panic("failed to create new traffic generator")
	}

	return generator.Run()
}
