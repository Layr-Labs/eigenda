package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

// Describes the log level.
type LogLevel string

const (
	// Log all levels
	LogLevelDebug LogLevel = "debug"
	// Log info, warn, error
	LogLevelInfo LogLevel = "info"
	// Log warn, error
	LogLevelWarn LogLevel = "warn"
	// Log only errors
	LogLevelError LogLevel = "error"
)

// Describes the log format.
type LogFormat string

const (
	// Log in JSON format.
	JSONLogFormat LogFormat = "json"
	// Log in human-readable text format.
	TextLogFormat LogFormat = "text"
)

var _ VerifiableConfig = &SimpleLoggerConfig{}

// Roughly equivalent to common.LoggerConfig, but without complex types that trip up the config parser. This
// struct should be used when embedding logger configuration in other config structs.
type SimpleLoggerConfig struct {
	// Format of the log output. Valid options are "json" and "text".
	Format LogFormat

	// Enable source code location
	AddSource bool

	// Minimum level to log. Valid options are "debug", "info", "warn", and "error".
	Level LogLevel

	// Time format, only supported with text handler
	TimeFormat string

	// Disable color, only supported with text handler (i.e. no color in json).
	NoColor bool
}

// Create a SimpleLoggerConfig with default values. These defaults are appropriate for production deployments.
func DefaultSimpleLoggerConfig() *SimpleLoggerConfig {
	return &SimpleLoggerConfig{
		Format:     JSONLogFormat,
		AddSource:  true,
		Level:      LogLevelDebug,
		TimeFormat: "",
		NoColor:    false,
	}
}

func (s *SimpleLoggerConfig) Verify() error {
	if s.Format != JSONLogFormat && s.Format != TextLogFormat {
		return fmt.Errorf("invalid log format: %s", s.Format)
	}

	if s.Level != LogLevelDebug && s.Level != LogLevelInfo && s.Level != LogLevelWarn && s.Level != LogLevelError {
		return fmt.Errorf("invalid log level: %s", s.Level)
	}

	return nil
}

// TODO(cody.littley): once all configurations are migrated to use SimpleLoggerConfig,
//  consider removing LoggerConfig entirely.

// Convert this SimpleLoggerConfig to a full LoggerConfig (i.e. the config the logger framework consumes).
func (s *SimpleLoggerConfig) ToLoggerConfig() (*common.LoggerConfig, error) {
	var level slog.Leveler
	switch s.Level {
	case LogLevelDebug:
		level = slog.LevelDebug
	case LogLevelInfo:
		level = slog.LevelInfo
	case LogLevelWarn:
		level = slog.LevelWarn
	case LogLevelError:
		level = slog.LevelError
	default:
		return nil, fmt.Errorf("invalid log level: %s", s.Level)
	}

	return &common.LoggerConfig{
		Format:       common.LogFormat(s.Format),
		OutputWriter: os.Stdout,
		HandlerOpts: logging.SLoggerOptions{
			AddSource:  s.AddSource,
			Level:      level,
			TimeFormat: s.TimeFormat,
			NoColor:    s.NoColor,
		},
	}, nil
}

// Build a logger from this SimpleLoggerConfig.
func (s *SimpleLoggerConfig) BuildLogger() (logging.Logger, error) {
	loggerConfig, err := s.ToLoggerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to convert SimpleLoggerConfig to LoggerConfig: %w", err)
	}

	logger, err := common.NewLogger(loggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	return logger, nil
}
