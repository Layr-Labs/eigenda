package telemetry

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/urfave/cli/v2"
)

const (
	OTelEnabledFlagName          = "otel.enabled"
	OTelServiceNameFlagName      = "otel.service-name"
	OTelExporterEndpointFlagName = "otel.exporter.otlp.endpoint"
	OTelExporterInsecureFlagName = "otel.exporter.otlp.insecure"
	OTelTraceSampleRateFlagName  = "otel.trace.sample-rate"
)

// Config holds the configuration for OpenTelemetry tracing
type Config struct {
	// Enabled determines if OpenTelemetry tracing is enabled
	Enabled bool
	// ServiceName is the name of the service for tracing
	ServiceName string
	// ExporterEndpoint is the OTLP exporter endpoint (e.g., "localhost:4317" for gRPC)
	ExporterEndpoint string
	// ExporterInsecure determines if the exporter should use an insecure connection
	ExporterInsecure bool
	// TraceSampleRate is the sampling rate for traces (0.0 to 1.0)
	TraceSampleRate float64
}

var _ config.VerifiableConfig = (*Config)(nil)

func CLIFlags(envPrefix, category string) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:     OTelEnabledFlagName,
			Usage:    "Enable OpenTelemetry tracing",
			EnvVars:  []string{envPrefix + "_OTEL_ENABLED"},
			Category: category,
			Value:    false,
		},
		&cli.StringFlag{
			Name:     OTelServiceNameFlagName,
			Usage:    "Service name for OpenTelemetry tracing",
			EnvVars:  []string{envPrefix + "_OTEL_SERVICE_NAME"},
			Category: category,
			Value:    "eigenda-proxy",
		},
		&cli.StringFlag{
			Name:     OTelExporterEndpointFlagName,
			Usage:    "OpenTelemetry OTLP exporter endpoint (e.g., 'localhost:4318' for HTTP)",
			EnvVars:  []string{envPrefix + "_OTEL_EXPORTER_OTLP_ENDPOINT"},
			Category: category,
			Value:    "localhost:4318",
		},
		&cli.BoolFlag{
			Name:     OTelExporterInsecureFlagName,
			Usage:    "Use insecure connection for OTLP exporter",
			EnvVars:  []string{envPrefix + "_OTEL_EXPORTER_OTLP_INSECURE"},
			Category: category,
			Value:    true,
		},
		&cli.Float64Flag{
			Name:     OTelTraceSampleRateFlagName,
			Usage:    "Trace sampling rate (0.0 to 1.0, where 1.0 means sample all traces)",
			EnvVars:  []string{envPrefix + "_OTEL_TRACE_SAMPLE_RATE"},
			Category: category,
			Value:    1.0,
		},
	}
}

func ReadConfig(ctx *cli.Context) Config {
	return Config{
		Enabled:          ctx.Bool(OTelEnabledFlagName),
		ServiceName:      ctx.String(OTelServiceNameFlagName),
		ExporterEndpoint: ctx.String(OTelExporterEndpointFlagName),
		ExporterInsecure: ctx.Bool(OTelExporterInsecureFlagName),
		TraceSampleRate:  ctx.Float64(OTelTraceSampleRateFlagName),
	}
}

// Verify implements config.VerifiableConfig
func (c *Config) Verify() error {
	if !c.Enabled {
		return nil
	}

	if c.ServiceName == "" {
		return fmt.Errorf("service name must be provided when OpenTelemetry is enabled")
	}

	if c.ExporterEndpoint == "" {
		return fmt.Errorf("exporter endpoint must be provided when OpenTelemetry is enabled")
	}

	if c.TraceSampleRate < 0.0 || c.TraceSampleRate > 1.0 {
		return fmt.Errorf("trace sample rate must be between 0.0 and 1.0, got: %f", c.TraceSampleRate)
	}

	return nil
}
