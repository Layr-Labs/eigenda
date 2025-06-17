package main

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/disperser/cmd/apiserver/flags"
	"github.com/Layr-Labs/eigenda/disperser/cmd/apiserver/lib"
	"github.com/urfave/cli"
	"log"
	"os"
)

var (
	// version is the version of the binary.
	version   string
	gitCommit string
	gitDate   string
)

func main() {
	app := cli.NewApp()
	app.Flags = flags.Flags
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)
	app.Name = "disperser"
	app.Usage = "EigenDA Disperser Server"
	app.Description = "Service for accepting blobs for dispersal"

	app.Action = lib.RunDisperserServer
	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}
