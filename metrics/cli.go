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

func DefaultCLIConfig() CLIConfig {
	return CLIConfig{
		Enabled:    false,
		ListenAddr: defaultListenAddr,
		ListenPort: defaultListenPort,
	}
}

func CLIFlags(envPrefix string, category string) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:     EnabledFlagName,
			Usage:    "Enable the metrics server",
			Category: category,
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

type CLIConfig struct {
	Enabled    bool
	ListenAddr string
	ListenPort int
}

func (m CLIConfig) Check() error {
	if !m.Enabled {
		return nil
	}

	if m.ListenPort < 0 || m.ListenPort > math.MaxUint16 {
		return ErrInvalidPort
	}

	return nil
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	return CLIConfig{
		Enabled:    ctx.Bool(EnabledFlagName),
		ListenAddr: ctx.String(ListenAddrFlagName),
		ListenPort: ctx.Int(PortFlagName),
	}
}
