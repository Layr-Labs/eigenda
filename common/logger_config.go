package common

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/Layr-Labs/eigensdk-go/logging"
	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/urfave/cli"
)

const (
	PathFlagName   = "log.path"
	LevelFlagName  = "log.level"
	FormatFlagName = "log.format"
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
		cli.StringFlag{
			Name:   PrefixFlag(flagPrefix, FormatFlagName),
			Usage:  "The format of the log file. Accepted options are 'json' and 'text'",
			Value:  "json",
			EnvVar: PrefixEnvVar(envPrefix, "LOG_FORMAT"),
		},
	}
}

// DefaultLoggerConfig returns a LoggerConfig with the default settings for a JSON logger.
// In general, this should be the baseline config for most services running in production.
func DefaultLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
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
func DefaultTextLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
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
func DefaultConsoleLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Format:       TextLogFormat,
		OutputWriter: os.Stdout,
		HandlerOpts: logging.SLoggerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
			NoColor:   false,
		},
	}
}

func ReadLoggerCLIConfig(ctx *cli.Context, flagPrefix string) (*LoggerConfig, error) {
	cfg := DefaultLoggerConfig()
	format := ctx.GlobalString(PrefixFlag(flagPrefix, FormatFlagName))
	if format == "json" {
		cfg.Format = JSONLogFormat
	} else if format == "text" {
		cfg.Format = TextLogFormat
	} else {
		return nil, fmt.Errorf("invalid log file format %s", format)
	}

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

	return cfg, nil
}

func NewLogger(cfg *LoggerConfig) (logging.Logger, error) {
	if cfg.Format == JSONLogFormat {
		return logging.NewJsonSLogger(cfg.OutputWriter, &cfg.HandlerOpts), nil
	}
	if cfg.Format == TextLogFormat {
		return logging.NewTextSLogger(cfg.OutputWriter, &cfg.HandlerOpts), nil
	}
	return nil, fmt.Errorf("unknown log format: %s", cfg.Format)
}

// InterceptorLogger returns a grpclogging.Logger that uses the provided logging.Logger.
// grpclogging.Logger is an interface that allows logging gRPC interceptor messages.
// Ref: https://github.com/grpc-ecosystem/go-grpc-middleware/blob/main/interceptors/logging/examples/slog/example_test.go
func InterceptorLogger(logger logging.Logger) grpclogging.Logger {
	return grpclogging.LoggerFunc(func(ctx context.Context, lvl grpclogging.Level, msg string, fields ...any) {
		switch lvl {
		case grpclogging.LevelDebug:
			logger.Debug(msg, fields...)
		case grpclogging.LevelInfo:
			logger.Info(msg, fields...)
		case grpclogging.LevelWarn:
			logger.Warn(msg, fields...)
		case grpclogging.LevelError:
			logger.Error(msg, fields...)
		default:
			logger.Info(msg, fields...)
		}
	})
}
