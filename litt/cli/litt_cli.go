package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli/v2"
)

// buildCliParser creates a command line parser for the LittDB CLI tool.
func buildCLIParser() *cli.App {
	app := &cli.App{
		Name:  "litt",
		Usage: "LittDB command line interface",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Enable debug mode. Program will pause for a debugger to attach.",
			},
		},
		Before: handleDebugMode,
		Commands: []*cli.Command{
			{
				Name:      "ls",
				Usage:     "List tables in a LittDB instance",
				ArgsUsage: "--src <path1> ... --src <pathN>",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "src",
						Aliases:  []string{"s"},
						Usage:    "Source paths where the DB data is found, at least one is required.",
						Required: true,
					},
				},
				Action: lsCommand,
			},
			{
				Name: "table-info",
				Usage: "Get information about a LittDB table. " +
					"If the DB is spread across multiple paths, all paths must be provided.",
				ArgsUsage: "--src <path1> ... --src <pathN> <table-name>",
				Args:      true,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "src",
						Aliases:  []string{"s"},
						Usage:    "Source paths where the DB data is found, at least one is required.",
						Required: true,
					},
				},
				Action: tableInfoCommand,
			},
			{
				Name:  "rebase",
				Usage: "Restructure LittDB file system layout.",
				ArgsUsage: "--src <source-path1> ... --src <source-pathN> " +
					"--dest <destination-path1> ... --dest <destination-pathN>",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "src",
						Aliases:  []string{"s"},
						Usage:    "Source paths where the data is found, at least one is required.",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:     "dst",
						Aliases:  []string{"d"},
						Usage:    "Destination paths for the rebased LittDB, at least one is required.",
						Required: true,
					},
					&cli.BoolFlag{
						Name:    "shallow",
						Aliases: []string{"S"},
						Usage: "If true, then copies are made shallowly " +
							"(e.g. with symlinks and hardlinks, where possible). ",
						Required: false,
					},
					&cli.BoolFlag{
						Name:     "preserve",
						Aliases:  []string{"p"},
						Usage:    "If enabled, then the old files are not removed.",
						Required: false,
					},
					&cli.BoolFlag{
						Name:     "quiet",
						Aliases:  []string{"q"},
						Usage:    "Reduces the verbosity of the output.",
						Required: false,
					},
				},
				Action: rebaseCommand,
			},
			{
				Name:      "benchmark",
				Usage:     "Run a LittDB benchmark.",
				ArgsUsage: "<path/to/benchmark/config.json>",
				Args:      true,
				Action:    benchmarkCommand,
			},
			{
				Name:      "prune",
				Usage:     "Delete data from a LittDB database/snapshot.",
				ArgsUsage: "--src <path1> ... --src <pathN> --max-age <duration in seconds>",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "src",
						Aliases:  []string{"s"},
						Usage:    "Source paths where the DB data is found, at least one is required.",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:     "table",
						Aliases:  []string{"t"},
						Usage:    "Prune this table. If not specified, all tables will be pruned.",
						Required: false,
					},
					&cli.Uint64Flag{
						Name:    "max-age",
						Aliases: []string{"a"},
						Usage: "Maximum age of segments to keep, in seconds. " +
							"Segments older than this will be deleted.",
						Value:    0, // Default to 0, meaning no age limit
						Required: true,
					},
				},
				Action: pruneCommand,
			},
			{
				Name:  "push",
				Usage: "Push data to a remote location using ssh and rsync.",
				ArgsUsage: "--src <source-path1> ... --src <source-pathN> " +
					"--dst <remote-path1> ... --dst <remote-pathN> " +
					"[-i path/to/key] [-p port] [--no-gc] [--quiet] [--threads 42] [--throttle 100]" +
					"<user>@<host>",
				Args: true,
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:     "src",
						Aliases:  []string{"s"},
						Usage:    "Source paths where the data is found, at least one is required.",
						Required: true,
					},
					&cli.StringSliceFlag{
						Name:     "dst",
						Aliases:  []string{"d"},
						Usage:    "Remote destination paths, at least one is required.",
						Required: true,
					},
					&cli.Uint64Flag{
						Name:    "port",
						Aliases: []string{"p"},
						Usage:   "SSH port to connect to the remote host.",
						Value:   22,
					},
					&cli.StringFlag{
						Name:    "key",
						Aliases: []string{"i"},
						Usage:   "Path to the SSH private key file for authentication.",
						Value:   "~/.ssh/id_rsa",
					},
					&cli.BoolFlag{
						Name:    "no-gc",
						Aliases: []string{"n"},
						Usage:   "If true, do not delete files pushed to the remote host.",
					},
					&cli.BoolFlag{
						Name:     "quiet",
						Aliases:  []string{"q"},
						Usage:    "Reduces the verbosity of the output.",
						Required: false,
					},
					&cli.Uint64Flag{
						Name:    "threads",
						Aliases: []string{"t"},
						Usage:   "Number of parallel rsync operations.",
						Value:   8,
					},
					&cli.Float64Flag{
						Name:    "throttle",
						Aliases: []string{"T"},
						Usage:   "Max network utilization, in mb/s",
						Value:   0,
					},
				},
				Action: pushCommand,
			},
		},
	}
	return app
}

// If the --debug flag is set, this function will block until SIGUSR1 is received to allow a debugger to attach.
func handleDebugMode(ctx *cli.Context) error {
	debugModeEnabled := ctx.Bool("debug")

	if !debugModeEnabled {
		return nil
	}

	logger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	if err != nil {
		return fmt.Errorf("failed to create logger: %v", err)
	}

	pid := os.Getpid()
	logger.Infof("Waiting for debugger to attach (pid: %d).\n", pid)

	logger.Infof("Press Enter to continue...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n') // block until newline is read

	return nil
}
