package logging

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/urfave/cli"
)

const (
	PathFlagName      = "log.path"
	FileLevelFlagName = "log.level-file"
	StdLevelFlagName  = "log.level-std"
)

type Config struct {
	Path      string
	Prefix    string
	FileLevel string
	StdLevel  string
}

func CLIFlags(envPrefix string, flagPrefix string) []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   common.PrefixFlag(flagPrefix, StdLevelFlagName),
			Usage:  `The lowest log level that will be output to stdout. Accepted options are "trace", "debug", "info", "warn", "error"`,
			Value:  "info",
			EnvVar: common.PrefixEnvVar(envPrefix, "STD_LOG_LEVEL"),
		},
		cli.StringFlag{
			Name:   common.PrefixFlag(flagPrefix, FileLevelFlagName),
			Usage:  `The lowest log level that will be output to file logs. Accepted options are "trace", "debug", "info", "warn", "error"`,
			Value:  "info",
			EnvVar: common.PrefixEnvVar(envPrefix, "FILE_LOG_LEVEL"),
		},
		cli.StringFlag{
			Name:   common.PrefixFlag(flagPrefix, PathFlagName),
			Usage:  "Path to file where logs will be written",
			Value:  "",
			EnvVar: common.PrefixEnvVar(envPrefix, "LOG_PATH"),
		},
	}
}

func DefaultCLIConfig() Config {
	return Config{
		Path:      "",
		FileLevel: "debug",
		StdLevel:  "debug",
	}
}

func ReadCLIConfig(ctx *cli.Context, flagPrefix string) Config {
	cfg := DefaultCLIConfig()
	cfg.StdLevel = ctx.GlobalString(common.PrefixFlag(flagPrefix, StdLevelFlagName))
	cfg.FileLevel = ctx.GlobalString(common.PrefixFlag(flagPrefix, FileLevelFlagName))
	cfg.Path = ctx.GlobalString(common.PrefixFlag(flagPrefix, PathFlagName))
	return cfg
}
