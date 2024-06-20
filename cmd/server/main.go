package main

import (
	"context"
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"github.com/Layr-Labs/eigenda-proxy/metrics"
	"github.com/Layr-Labs/eigenda-proxy/server"
	"github.com/ethereum-optimism/optimism/op-service/cliapp"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/metrics/doc"
	"github.com/ethereum-optimism/optimism/op-service/opio"
)

var Version = "v1.1.0"

func main() {
	oplog.SetupDefaults()

	app := cli.NewApp()
	app.Flags = cliapp.ProtectFlags(server.Flags)
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

	ctx := opio.WithInterruptBlocker(context.Background())
	err := app.RunContext(ctx, os.Args)
	if err != nil {
		log.Crit("Application failed", "message", err)
	}

}
