package relayload

import (
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/tools/relayload/flags"
	"github.com/urfave/cli"
)

type Config struct {
	LoggerConfig common.LoggerConfig
	Timeout      time.Duration
	RelayUrl     string
	OperatorId   string
	OperatorPKey string
	DataApiUrl   string
	NumThreads   int
	RangeSizes   []int
	RequestSizes []int
}

func stringSliceToIntSlice(strSlice []string) []int {
	intSlice := make([]int, len(strSlice))
	for i, s := range strSlice {
		intSlice[i], _ = strconv.Atoi(s)
	}
	return intSlice
}

func ReadConfig(ctx *cli.Context) *Config {
	return &Config{
		RelayUrl:     ctx.GlobalString(flags.RelayUrlFlag.Name),
		OperatorId:   ctx.GlobalString(flags.OperatorIdFlag.Name),
		OperatorPKey: ctx.GlobalString(flags.OperatorPKeyFlag.Name),
		DataApiUrl:   ctx.GlobalString(flags.DataApiUrlFlag.Name),
		NumThreads:   ctx.GlobalInt(flags.NumThreadsFlag.Name),
		RangeSizes:   stringSliceToIntSlice(strings.Split(ctx.GlobalString(flags.RangeSizesFlag.Name), ",")),
		RequestSizes: stringSliceToIntSlice(strings.Split(ctx.GlobalString(flags.RequestSizesFlag.Name), ",")),
	}
}

func NewConfig(ctx *cli.Context) (*Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, err
	}

	config := ReadConfig(ctx)
	config.LoggerConfig = *loggerConfig

	return config, nil
}
