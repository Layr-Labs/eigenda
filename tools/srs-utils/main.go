package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Layr-Labs/eigenda/tools/srs-utils/parser"
	"github.com/Layr-Labs/eigenda/tools/srs-utils/verifier"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Commands: []cli.Command{
			{
				Name:    "verify",
				Aliases: []string{"v"},
				Usage:   "verify if the parsed SRS are consistent",
				Action: func(cCtx *cli.Context) error {
					config := verifier.ReadCLIConfig(cCtx)
					verifier.VerifySRS(config)
					return nil
				},
				Flags: verifier.Flags,
			},
			{
				Name:    "parse",
				Aliases: []string{"p"},
				Usage:   "parse data from ptau challenge file into EigenDA SRS format",
				Flags:   parser.Flags,
				Action: func(cCtx *cli.Context) error {
					config := parser.ReadCLIConfig(cCtx)
					fmt.Printf("config %v\n", config.PtauPath)
					parser.ParsePtauChallenge(config)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
