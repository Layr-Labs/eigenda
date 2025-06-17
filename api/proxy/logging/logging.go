package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/urfave/cli/v2"
)

/*
	TODO: https://github.com/Layr-Labs/eigenda-proxy/issues/268

	This CLI logic is already defined in the eigenda monorepo:
	 https://github.com/Layr-Labs/eigenda/blob/0d293cc031987c43f653535732c6e1f1fa65a0b2/common/logger_config.go
	This regression is due to the fact the proxy leverage urfave/cli/v2 whereas
	core eigenda predominantly uses urfave/cli (i.e, v1).

*/

const (
	PathFlagName   = "path"
	LevelFlagName  = "level"
	FormatFlagName = "format"
	// deprecated
	PidFlagName   = "pid"
	ColorFlagName = "color"

	// Flag
	FlagPrefix = "log"
)

type LogFormat string

const (
	JSONLogFormat LogFormat = "json"
	TextLogFormat LogFormat = "text"
)

type LoggerConfig struct {
	Format       LogFormat
	OutputWriter io.Writer
	HandlerOpts  logging.SLoggerOptions
}

func withEnvPrefix(envPrefix, s string) []string {
	return []string{envPrefix + "_LOG_" + s}
}

func CLIFlags(envPrefix string, category string) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:     common.PrefixFlag(FlagPrefix, LevelFlagName),
			Category: category,
			Usage:    `The lowest log level that will be output. Accepted options are "debug", "info", "warn", "error"`,
			Value:    "info",
			EnvVars:  withEnvPrefix(envPrefix, "LEVEL"),
		},
		&cli.StringFlag{
			Name:     common.PrefixFlag(FlagPrefix, PathFlagName),
			Category: category,
			Usage:    "Path to file where logs will be written",
			Value:    "",
			EnvVars:  withEnvPrefix(envPrefix, "PATH"),
		},
		&cli.StringFlag{
			Name:     common.PrefixFlag(FlagPrefix, FormatFlagName),
			Category: category,
			Usage:    "The format of the log file. Accepted options are 'json' and 'text'",
			Value:    "text",
			EnvVars:  withEnvPrefix(envPrefix, "FORMAT"),
		},
		// Deprecated since used by op-service logging which has been replaced
		// by eigengo-sdk logger
		&cli.BoolFlag{
			Name:     common.PrefixFlag(FlagPrefix, PidFlagName),
			Category: category,
			Usage:    "Show pid in the log",
			EnvVars:  withEnvPrefix(envPrefix, "PID"),
			Hidden:   true,
			Action: func(_ *cli.Context, _ bool) error {
				return fmt.Errorf("flag --%s is deprecated", PidFlagName)
			},
		},
		&cli.BoolFlag{
			Name:     common.PrefixFlag(FlagPrefix, ColorFlagName),
			Category: category,
			Usage:    "Color the log output if in terminal mode",
			EnvVars:  []string{common.PrefixEnvVar(envPrefix, "LOG_COLOR")},
			Hidden:   true,
			Action: func(_ *cli.Context, _ bool) error {
				return fmt.Errorf("flag --%s is deprecated", ColorFlagName)
			},
		},
	}
}

// DefaultLoggerConfig returns a LoggerConfig with the default settings for a JSON logger.
// In general, this should be the baseline config for most services running in production.
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Format:       JSONLogFormat,
		OutputWriter: os.Stdout,
		HandlerOpts: logging.SLoggerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
			NoColor:   true,
		},
	}
}

// DefaultTextLoggerConfig returns a LoggerConfig with the default settings for a text logger.
// For use in tests or other scenarios where the logs are consumed by humans.
func DefaultTextLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Format:       TextLogFormat,
		OutputWriter: os.Stdout,
		HandlerOpts: logging.SLoggerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
			NoColor:   true, // color is nice in the console, but not nice when written to a file
		},
	}
}

// DefaultConsoleLoggerConfig returns a LoggerConfig with the default settings
// for logging to a console (i.e. with human eyeballs). Adds color, and so should
// not be used when logs are captured in a file.
func DefaultConsoleLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Format:       TextLogFormat,
		OutputWriter: os.Stdout,
		HandlerOpts: logging.SLoggerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
			NoColor:   false,
		},
	}
}

func ReadLoggerCLIConfig(ctx *cli.Context) (*LoggerConfig, error) {
	cfg := DefaultLoggerConfig()
	format := ctx.String(common.PrefixFlag(FlagPrefix, FormatFlagName))
	switch format {
	case "json":
		cfg.Format = JSONLogFormat

	case "text":
		cfg.Format = TextLogFormat

	default:
		return nil, fmt.Errorf("invalid log file format %s", format)
	}

	path := ctx.String(common.PrefixFlag(FlagPrefix, PathFlagName))
	if path != "" {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		cfg.OutputWriter = io.MultiWriter(os.Stdout, f)
	}
	logLevel := ctx.String(common.PrefixFlag(FlagPrefix, LevelFlagName))
	var level slog.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		panic("failed to parse log level " + logLevel)
	}
	cfg.HandlerOpts.Level = level

	return &cfg, nil
}

func NewLogger(cfg LoggerConfig) (logging.Logger, error) {
	if cfg.Format == JSONLogFormat {
		return logging.NewJsonSLogger(cfg.OutputWriter, &cfg.HandlerOpts), nil
	}
	if cfg.Format == TextLogFormat {
		return logging.NewTextSLogger(cfg.OutputWriter, &cfg.HandlerOpts), nil
	}
	return nil, fmt.Errorf("unknown log format: %s", cfg.Format)
}
