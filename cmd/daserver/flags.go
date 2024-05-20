package main

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
	opservice "github.com/ethereum-optimism/optimism/op-service"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
)

const (
	ListenAddrFlagName = "addr"
	PortFlagName       = "port"
)

const EnvVarPrefix = "OP_PLASMA_DA_SERVER"

func prefixEnvVars(name string) []string {
	return opservice.PrefixEnvVar(EnvVarPrefix, name)
}

var (
	ListenAddrFlag = &cli.StringFlag{
		Name:    ListenAddrFlagName,
		Usage:   "server listening address",
		Value:   "127.0.0.1",
		EnvVars: prefixEnvVars("ADDR"),
	}
	PortFlag = &cli.IntFlag{
		Name:    PortFlagName,
		Usage:   "server listening port",
		Value:   3100,
		EnvVars: prefixEnvVars("PORT"),
	}
)

var requiredFlags = []cli.Flag{
	ListenAddrFlag,
	PortFlag,
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

type CLIConfig struct {
	FileStoreDirPath string
	S3Bucket         string
	EigenDAConfig    eigenda.Config
	MetricsCfg       opmetrics.CLIConfig
}

func ReadCLIConfig(ctx *cli.Context) CLIConfig {
	return CLIConfig{
		EigenDAConfig: eigenda.ReadConfig(ctx),
	}
}

func (c CLIConfig) Check() error {

	err := c.EigenDAConfig.Check()
	if err != nil {
		return err
	}
	return nil
}

func (c CLIConfig) EigenDAEnabled() bool {
	return c.EigenDAConfig.RPC != ""
}

func CheckRequired(ctx *cli.Context) error {
	for _, f := range requiredFlags {
		if !ctx.IsSet(f.Names()[0]) {
			return fmt.Errorf("flag %s is required", f.Names()[0])
		}
	}
	return nil
}
