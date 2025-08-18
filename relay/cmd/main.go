package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/relay/cmd/flags"
	"github.com/Layr-Labs/eigenda/relay/cmd/lib"
	"github.com/urfave/cli"
)

var (
	version   string
	gitCommit string
	gitDate   string
)

func main() {
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)
	app.Name = "relay"
	app.Usage = "EigenDA Relay"
	app.Description = "EigenDA relay for serving blobs and chunks data"

	app.Action = lib.RunRelay
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}
	select {}
}
