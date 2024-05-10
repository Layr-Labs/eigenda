package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
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
	config, err := traffic.NewConfig(ctx)
	if err != nil {
		return err
	}

	var signer core.BlobRequestSigner
	if config.SignerPrivateKey != "" {
		log.Println("Using signer private key")
		signer = auth.NewSigner(config.SignerPrivateKey)
	}

	generator, err := traffic.NewTrafficGenerator(config, signer)
	if err != nil {
		panic("failed to create new traffic generator")
	}

	return generator.Run()
}
