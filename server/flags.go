package server

import (
	"github.com/Layr-Labs/eigenda-proxy/store"
	"github.com/urfave/cli/v2"

	opservice "github.com/ethereum-optimism/optimism/op-service"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
)

const (
	ListenAddrFlagName = "addr"
	PortFlagName       = "port"
)

const EnvVarPrefix = "EIGENDA_PROXY"

func prefixEnvVars(name string) []string {
	return opservice.PrefixEnvVar(EnvVarPrefix, name)
}

// Flags contains the list of configuration options available to the binary.
var Flags = []cli.Flag{
	&cli.StringFlag{
		Name:    ListenAddrFlagName,
		Usage:   "server listening address",
		Value:   "0.0.0.0",
		EnvVars: prefixEnvVars("ADDR"),
	},
	&cli.IntFlag{
		Name:    PortFlagName,
		Usage:   "server listening port",
		Value:   3100,
		EnvVars: prefixEnvVars("PORT"),
	},
}

func init() {
	Flags = append(Flags, oplog.CLIFlags(EnvVarPrefix)...)
	Flags = append(Flags, CLIFlags()...)
	Flags = append(Flags, opmetrics.CLIFlags(EnvVarPrefix)...)
}

type CLIConfig struct {
	S3Config      store.S3Config
	EigenDAConfig Config
	MetricsCfg    opmetrics.CLIConfig
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	config := ReadConfig(ctx)
	return CLIConfig{
		EigenDAConfig: config,
		MetricsCfg:    opmetrics.ReadCLIConfig(ctx),
		S3Config:      config.S3Config,
	}
}

func (c CLIConfig) Check() error {
	err := c.EigenDAConfig.Check()
	if err != nil {
		return err
	}
	return nil
}
