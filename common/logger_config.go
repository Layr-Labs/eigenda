package common

import (
	"io"
	"log/slog"
	"os"

	"github.com/urfave/cli"
)

const (
	PathFlagName  = "log.path"
	LevelFlagName = "log.level"
)

type LoggerConfig struct {
	OutputWriter io.Writer
	HandlerOpts  slog.HandlerOptions
}

func LoggerCLIFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   PrefixFlag(flagPrefix, LevelFlagName),
			Usage:  `The lowest log level that will be output. Accepted options are "debug", "info", "warn", "error"`,
			Value:  "info",
			EnvVar: PrefixEnvVar(envPrefix, "LOG_LEVEL"),
		},
		cli.StringFlag{
			Name:   PrefixFlag(flagPrefix, PathFlagName),
			Usage:  "Path to file where logs will be written",
			Value:  "",
			EnvVar: PrefixEnvVar(envPrefix, "LOG_PATH"),
		},
	}
}

func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		OutputWriter: os.Stdout,
		HandlerOpts: slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
	}
}

func ReadLoggerCLIConfig(ctx *cli.Context, flagPrefix string) (*LoggerConfig, error) {
	cfg := DefaultLoggerConfig()
	path := ctx.GlobalString(PrefixFlag(flagPrefix, PathFlagName))
	if path != "" {
		f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}
		cfg.OutputWriter = io.MultiWriter(os.Stdout, f)
	}
	logLevel := ctx.GlobalString(PrefixFlag(flagPrefix, LevelFlagName))
	var level slog.Level
	err := level.UnmarshalText([]byte(logLevel))
	if err != nil {
		panic("failed to parse log level " + logLevel)
	}
	cfg.HandlerOpts.Level = level

	return &cfg, nil
}
