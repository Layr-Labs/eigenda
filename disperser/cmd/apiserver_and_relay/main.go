package main

import (
	"fmt"
	"log"
	"os"

	apiserverFlags "github.com/Layr-Labs/eigenda/disperser/cmd/apiserver/flags"
	apiserverLib "github.com/Layr-Labs/eigenda/disperser/cmd/apiserver/lib"
	relayFlags "github.com/Layr-Labs/eigenda/relay/cmd/flags"
	relayLib "github.com/Layr-Labs/eigenda/relay/cmd/lib"
	"github.com/urfave/cli"
)

var (
	// version, gitCommit, gitDate are populated at build time (via -ldflags)
	version   string
	gitCommit string
	gitDate   string
)

func main() {
	app := cli.NewApp()
	app.Description = "EigenDA Disperser API Server (accepts blobs for dispersal) and Relay (serves blobs and chunks data)"
	app.Name = "API Server and Relay"
	app.Usage = "EigenDA Disperser API Server and Relay"
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)

	app.Commands = []cli.Command{
		{
			Name:   "apiserver",
			Usage:  "Run the EigenDA Disperser API server",
			Flags:  apiserverFlags.Flags,
			Action: apiserverLib.RunDisperserServer,
		},
		{
			Name:   "relay",
			Usage:  "Run the EigenDA Relay",
			Flags:  relayFlags.Flags,
			Action: relayLib.RunRelay,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}
