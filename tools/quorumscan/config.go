package quorumscan

import (
	"strconv"
	"strings"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/tools/quorumscan/flags"
	"github.com/urfave/cli"
)

type Config struct {
	LoggerConfig       common.LoggerConfig
	BlockNumber        uint64
	Workers            int
	Timeout            time.Duration
	UseRetrievalClient bool
	QuorumIDs          []core.QuorumID
	TopN               uint

	ChainStateConfig thegraph.Config
	EthClientConfig  geth.EthClientConfig

	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
}

func ReadConfig(ctx *cli.Context) *Config {
	quorumIDsStr := ctx.String(flags.QuorumIDsFlag.Name)
	quorumIDs := []core.QuorumID{}

	// Parse comma-separated quorum IDs
	if quorumIDsStr != "" {
		for _, idStr := range strings.Split(quorumIDsStr, ",") {
			if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
				quorumIDs = append(quorumIDs, core.QuorumID(id))
			}
		}
	}

	return &Config{
		ChainStateConfig:              thegraph.ReadCLIConfig(ctx),
		EthClientConfig:               geth.ReadEthClientConfig(ctx),
		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
		QuorumIDs:                     quorumIDs,
		BlockNumber:                   ctx.Uint64(flags.BlockNumberFlag.Name),
		TopN:                          ctx.Uint(flags.TopNFlag.Name),
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
