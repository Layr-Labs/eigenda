/*
EigenDA Proxy provides a simple REST API to facilitate interacting with the EigenDA Network.
*/
package main

import (
	"context"
	"os"

	"github.com/Layr-Labs/eigenda-proxy/config"
	"github.com/ethereum/go-ethereum/log"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	"github.com/ethereum-optimism/optimism/op-service/ctxinterrupt"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/metrics/doc"
)

var (
	Version = "unknown"
	Commit  = "unknown"
	Date    = "unknown"
)

func main() {
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(config.Flags)
	app.Version = Version
	app.Name = "eigenda-proxy"
	app.Usage = "EigenDA Proxy Sidecar Service"
	app.Description = "Service for more trustless and secure interactions with EigenDA"
	app.Action = StartProxySvr
	app.Commands = []*cli.Command{
		{
			Name:        "doc",
			Subcommands: doc.NewSubcommands(metrics.NewMetrics("default")),
		},
	}

	// load env file (if applicable)
	if p := os.Getenv("ENV_PATH"); p != "" {
		if err := godotenv.Load(p); err != nil {
			panic(err)
		}
	}

	ctx := ctxinterrupt.WithSignalWaiterMain(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}
}
