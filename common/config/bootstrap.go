package config

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/pprof"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/urfave/cli/v2"
)

var (
	pprofFlag = &cli.IntFlag{
		Name:    "pprof-port",
		Aliases: []string{"p"},
		Usage:   "Port for the pprof server.",
		Value:   6060,
	}
	pprofPortFlag = &cli.IntFlag{
		Name:    "pprof-port",
		Aliases: []string{"o"},
		Usage:   "Port for the pprof server.",
		Value:   6060,
	}
	debugFlag = &cli.BoolFlag{
		Name:    "debug",
		Aliases: []string{"d"},
		Usage:   "Enable debug mode. Program will pause for a debugger to attach.",
	}
	disableEnvVarsFlag = &cli.BoolFlag{
		Name:    "disable-env-vars",
		Aliases: []string{"e"},
		Usage:   "Disable loading configuration from environment variables.",
	}
	overrideEnvPrefixFlag = &cli.StringFlag{
		Name:    "env-prefix",
		Aliases: []string{"r"},
		Usage:   "If set, overrides the environment variable prefix used to load configuration from env vars.",
	}
	configFileFlag = &cli.StringSliceFlag{
		Name:    "config",
		Aliases: []string{"c"},
		Usage:   "Path to a configuration file. Can be specified multiple times to load multiple files.",
	}
)

// Reads command line arguments, loads configuration from files and environment variables as specified.
func Bootstrap[T VerifiableConfig](
	constructor func() T,
	defaultEnvVarPrefix string,
) (T, error) {

	// We need a logger before we have a logger config. Once we parse config, we can initialize the real logger.
	bootstrapLogger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to create bootstrap logger: %w", err)
	}

	action, cfgChan := buildHandler(bootstrapLogger, constructor, defaultEnvVarPrefix)

	app := &cli.App{
		Flags: []cli.Flag{
			pprofPortFlag,
			debugFlag,
			disableEnvVarsFlag,
			overrideEnvPrefixFlag,
			configFileFlag,
		},
		Action: action,
	}

	err = app.Run(os.Args)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error parsing command line arguments: %w", err)
	}

	cfg := <-cfgChan
	return cfg, nil
}

func buildHandler[T VerifiableConfig](
	logger logging.Logger,
	constructor func() T,
	defaultEnvVarPrefix string,
) (cli.ActionFunc, chan T) {

	cfgChan := make(chan T, 1)
	action := func(cliCTX *cli.Context) error {
		pprofEnabled := cliCTX.Bool(pprofFlag.Name)
		pprofPort := cliCTX.Int(pprofPortFlag.Name)
		debug := cliCTX.Bool(debugFlag.Name)
		disableEnvVars := cliCTX.Bool(disableEnvVarsFlag.Name)
		overrideEnvPrefix := cliCTX.String(overrideEnvPrefixFlag.Name)
		configFiles := cliCTX.StringSlice(configFileFlag.Name)

		prefix := defaultEnvVarPrefix
		if disableEnvVars {
			prefix = ""
		} else if overrideEnvPrefix != "" {
			prefix = overrideEnvPrefix
		}

		if debug {
			waitForDebugger(logger)
		}

		if pprofEnabled {
			startPprofServer(logger, pprofPort)
		}

		cfg, err := ParseConfig(constructor, prefix, configFiles...)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		cfgChan <- cfg
		return nil
	}
	return action, cfgChan
}

// waitForDebugger pauses execution to allow a human time to attach a debugger to the process.
func waitForDebugger(logger logging.Logger) {
	pid := os.Getpid()
	logger.Infof("Waiting for debugger to attach (pid: %d).\n", pid)

	logger.Infof("Press Enter to continue...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n') // block until newline is read
}

func startPprofServer(logger logging.Logger, port int) {
	logger.Infof("pprof enabled on port %d", port)
	profiler := pprof.NewPprofProfiler(fmt.Sprintf("%d", port), logger)
	go profiler.Start()
}
