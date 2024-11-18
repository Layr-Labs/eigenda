package main

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/cmd/controller/flags"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/urfave/cli"
)

const MaxUint16 = ^uint16(0)

type Config struct {
	EncodingManagerConfig          controller.EncodingManagerConfig
	DispatcherConfig               controller.DispatcherConfig
	NumConcurrentEncodingRequests  int
	NumConcurrentDispersalRequests int
	NodeClientCacheSize            int

	DynamoDBTableName string

	EthClientConfig  geth.EthClientConfig
	AwsClientConfig  aws.ClientConfig
	LoggerConfig     common.LoggerConfig
	IndexerConfig    indexer.Config
	ChainStateConfig thegraph.Config
	UseGraph         bool
	IndexerDataDir   string

	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
}

func NewConfig(ctx *cli.Context) (Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}
	ethClientConfig := geth.ReadEthClientConfigRPCOnly(ctx)
	numRelayAssignments := ctx.GlobalInt(flags.NumRelayAssignmentFlag.Name)
	if numRelayAssignments < 1 || numRelayAssignments > int(MaxUint16) {
		return Config{}, fmt.Errorf("invalid number of relay assignments: %d", numRelayAssignments)
	}
	availableRelays := ctx.GlobalIntSlice(flags.AvailableRelaysFlag.Name)
	if len(availableRelays) == 0 {
		return Config{}, fmt.Errorf("no available relays specified")
	}
	relays := make([]corev2.RelayKey, len(availableRelays))
	for i, relay := range availableRelays {
		if relay < 0 || relay > 65_535 {
			return Config{}, fmt.Errorf("invalid relay: %d", relay)
		}
		relays[i] = corev2.RelayKey(relay)
	}
	config := Config{
		DynamoDBTableName: ctx.GlobalString(flags.DynamoDBTableNameFlag.Name),
		EthClientConfig:   ethClientConfig,
		AwsClientConfig:   aws.ReadClientConfig(ctx, flags.FlagPrefix),
		LoggerConfig:      *loggerConfig,
		EncodingManagerConfig: controller.EncodingManagerConfig{
			PullInterval:           ctx.GlobalDuration(flags.EncodingPullIntervalFlag.Name),
			EncodingRequestTimeout: ctx.GlobalDuration(flags.EncodingRequestTimeoutFlag.Name),
			StoreTimeout:           ctx.GlobalDuration(flags.EncodingStoreTimeoutFlag.Name),
			NumEncodingRetries:     ctx.GlobalInt(flags.NumEncodingRetriesFlag.Name),
			NumRelayAssignment:     uint16(numRelayAssignments),
			AvailableRelays:        relays,
			EncoderAddress:         ctx.GlobalString(flags.EncoderAddressFlag.Name),
		},
		DispatcherConfig: controller.DispatcherConfig{
			PullInterval:           ctx.GlobalDuration(flags.DispatcherPullIntervalFlag.Name),
			FinalizationBlockDelay: ctx.GlobalUint64(flags.FinalizationBlockDelayFlag.Name),
			NodeRequestTimeout:     ctx.GlobalDuration(flags.NodeRequestTimeoutFlag.Name),
			NumRequestRetries:      ctx.GlobalInt(flags.NumRequestRetriesFlag.Name),
		},
		NumConcurrentEncodingRequests:  ctx.GlobalInt(flags.NumConcurrentEncodingRequestsFlag.Name),
		NumConcurrentDispersalRequests: ctx.GlobalInt(flags.NumConcurrentDispersalRequestsFlag.Name),
		NodeClientCacheSize:            ctx.GlobalInt(flags.NodeClientCacheNumEntriesFlag.Name),
		IndexerConfig:                  indexer.ReadIndexerConfig(ctx),
		ChainStateConfig:               thegraph.ReadCLIConfig(ctx),
		UseGraph:                       ctx.GlobalBool(flags.UseGraphFlag.Name),
		IndexerDataDir:                 ctx.GlobalString(flags.IndexerDataDirFlag.Name),

		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
	}
	return config, nil
}
