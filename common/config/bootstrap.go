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

// TODO(cody.littley): we should migrate this from urfave to cobra, since we already use cobra for the config
// framework. This would let us drop the urfave dependency.

var (
	pprofFlag = &cli.BoolFlag{
		Name:    "pprof",
		Aliases: []string{"p"},
		Usage:   "If set, starts a pprof server.",
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
	verifyConfigFlag = &cli.BoolFlag{
		Name:    "verify-config",
		Aliases: []string{"v"},
		Usage:   "If set, verifies configuration then exits.",
	}
)

// Reads command line arguments, loads configuration from files and environment variables as specified.
func Bootstrap[T DocumentedConfig](
	// A function that returns a new instance of the config struct with default values set.
	constructor func() T,
	// A list of environment variables that should be ignored when sanity checking environment variables.
	// Useful for situations where external systems set environment variables that would otherwise cause problems.
	ignoredEnvVars ...string,
) (T, error) {

	// We need a logger before we have a logger config. Once we parse config, we can initialize the real logger.
	bootstrapLogger, err := common.NewLogger(common.DefaultConsoleLoggerConfig())
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to create bootstrap logger: %w", err)
	}

	action, cfgChan := buildHandler(bootstrapLogger, constructor, ignoredEnvVars)

	app := &cli.App{
		Flags: []cli.Flag{
			pprofFlag,
			pprofPortFlag,
			debugFlag,
			disableEnvVarsFlag,
			overrideEnvPrefixFlag,
			configFileFlag,
			verifyConfigFlag,
		},
		Action: action,
	}

	err = app.Run(os.Args)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("error parsing command line arguments: %w", err)
	}

	// If the help flag was set, the action never runs and cfgChan is never written to.
	// Check if we have a config; if not, the help was shown and we should exit.
	select {
	case cfg := <-cfgChan:
		return cfg, nil
	default:
		// Help was shown, return zero value
		var zero T
		return zero, nil
	}
}

func buildHandler[T DocumentedConfig](
	logger logging.Logger,
	constructor func() T,
	ignoredEnvVars []string,
) (cli.ActionFunc, chan T) {

	cfgChan := make(chan T, 1)
	action := func(cliCTX *cli.Context) error {
		pprofEnabled := cliCTX.Bool(pprofFlag.Name)
		pprofPort := cliCTX.Int(pprofPortFlag.Name)
		debug := cliCTX.Bool(debugFlag.Name)
		disableEnvVars := cliCTX.Bool(disableEnvVarsFlag.Name)
		overrideEnvPrefix := cliCTX.String(overrideEnvPrefixFlag.Name)
		configFiles := cliCTX.StringSlice(configFileFlag.Name)
		verifyConfig := cliCTX.Bool(verifyConfigFlag.Name)

		if debug {
			waitForDebugger(logger)
		}

		if pprofEnabled {
			startPprofServer(logger, pprofPort)
		}

		defaultConfig := constructor()

		prefix := defaultConfig.GetEnvVarPrefix()
		if disableEnvVars {
			prefix = ""
		} else if overrideEnvPrefix != "" {
			prefix = overrideEnvPrefix
		}

		cfg, err := ParseConfig(logger, defaultConfig, prefix, ignoredEnvVars, configFiles...)
		if err != nil {
			return fmt.Errorf("failed to load configuration: %w", err)
		}

		if verifyConfig {
			logger.Info("Configuration is valid. Exiting.")
			os.Exit(0)
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
