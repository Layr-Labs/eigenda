package metrics

import (
	"errors"
	"math"

	"github.com/urfave/cli/v2"
)

const (
	EnabledFlagName    = "metrics.enabled"
	ListenAddrFlagName = "metrics.addr"
	PortFlagName       = "metrics.port"
	defaultListenAddr  = "0.0.0.0"
	defaultListenPort  = 7300

	EnvPrefix = "metrics"
)

var ErrInvalidPort = errors.New("invalid metrics port")

func withEnvPrefix(envPrefix, s string) []string {
	return []string{envPrefix + "_METRICS_" + s}
}

func DefaultConfig() Config {
	return Config{
		Enabled: false,
		Host:    defaultListenAddr,
		Port:    defaultListenPort,
	}
}

func CLIFlags(envPrefix string, category string) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:     EnabledFlagName,
			Usage:    "Enable the metrics server. On by default, so use --metrics.enabled=false to disable.",
			Category: category,
			Value:    true,
			EnvVars:  withEnvPrefix(envPrefix, "ENABLED"),
		},
		&cli.StringFlag{
			Name:     ListenAddrFlagName,
			Usage:    "Metrics listening address",
			Category: category,
			Value:    defaultListenAddr,
			EnvVars:  withEnvPrefix(envPrefix, "ADDR"),
		},
		&cli.IntFlag{
			Name:     PortFlagName,
			Usage:    "Metrics listening port",
			Category: category,
			Value:    defaultListenPort,
			EnvVars:  withEnvPrefix(envPrefix, "PORT"),
		},
	}
}

func (m Config) Check() error {
	if !m.Enabled {
		return nil
	}

	if m.Port < 0 || m.Port > math.MaxUint16 {
		return ErrInvalidPort
	}

	return nil
}

func ReadConfig(ctx *cli.Context) Config {
	return Config{
		Enabled: ctx.Bool(EnabledFlagName),
		Host:    ctx.String(ListenAddrFlagName),
		Port:    ctx.Int(PortFlagName),
	}
}
