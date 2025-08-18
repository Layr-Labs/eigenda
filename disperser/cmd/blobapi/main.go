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
	app.Flags = mergeFlags(apiserverFlags.Flags, relayFlags.Flags)
	app.Description = "EigenDA Disperser API Server (accepts blobs for dispersal) and Relay (serves blobs and chunks data)"
	app.Name = "BlobAPI"
	app.Usage = "EigenDA Disperser API Server and Relay"
	app.Version = fmt.Sprintf("%s-%s-%s", version, gitCommit, gitDate)

	app.Action = func(ctx *cli.Context) error {
		// exactly the same code you had in the subcommand:
		apiserverDone := make(chan error, 1)
		relayDone := make(chan error, 1)

		go func() { apiserverDone <- apiserverLib.RunDisperserServer(ctx) }()
		go func() { relayDone <- relayLib.RunRelay(ctx) }()

		select {
		case err := <-apiserverDone:
			return fmt.Errorf("apiserver exited: %w", err)
		case err := <-relayDone:
			return fmt.Errorf("relay exited: %w", err)
		}
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("application failed: %v", err)
	}

	select {}
}

// mergeFlags merges two slices of cli.Flag, dropping any with the same primary name.
func mergeFlags(a, b []cli.Flag) []cli.Flag {
	seen := make(map[string]bool, len(a)+len(b))
	out := make([]cli.Flag, 0, len(a)+len(b))

	// First add all of “a”
	for _, f := range a {
		name := f.GetName()
		seen[name] = true
		out = append(out, f)
	}
	// Then add only those in “b” whose primary name we haven’t seen
	for _, f := range b {
		if !seen[f.GetName()] {
			seen[f.GetName()] = true
			out = append(out, f)
		}
	}
	return out
}
