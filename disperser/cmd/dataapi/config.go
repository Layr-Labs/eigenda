package main

import (
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/disperser/cmd/dataapi/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus"
	"github.com/urfave/cli"
)

type Config struct {
	AwsClientConfig  aws.ClientConfig
	BlobstoreConfig  blobstore.Config
	EthClientConfig  geth.EthClientConfig
	LoggerConfig     common.LoggerConfig
	PrometheusConfig prometheus.Config
	MetricsConfig    dataapi.MetricsConfig
	ChainStateConfig thegraph.Config

	SocketAddr                   string
	PrometheusApiAddr            string
	SubgraphApiBatchMetadataAddr string
	SubgraphApiOperatorStateAddr string
	ServerMode                   string
	AllowOrigins                 []string

	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string

	DisperserHostname  string
	ChurnerHostname    string
	BatcherHealthEndpt string
}

func NewConfig(ctx *cli.Context) (Config, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}
	ethClientConfig := geth.ReadEthClientConfig(ctx)
	config := Config{
		BlobstoreConfig: blobstore.Config{
			BucketName: ctx.GlobalString(flags.S3BucketNameFlag.Name),
			TableName:  ctx.GlobalString(flags.DynamoTableNameFlag.Name),
		},
		AwsClientConfig:               aws.ReadClientConfig(ctx, flags.FlagPrefix),
		EthClientConfig:               ethClientConfig,
		LoggerConfig:                  *loggerConfig,
		SocketAddr:                    ctx.GlobalString(flags.SocketAddrFlag.Name),
		SubgraphApiBatchMetadataAddr:  ctx.GlobalString(flags.SubgraphApiBatchMetadataAddrFlag.Name),
		SubgraphApiOperatorStateAddr:  ctx.GlobalString(flags.SubgraphApiOperatorStateAddrFlag.Name),
		BLSOperatorStateRetrieverAddr: ctx.GlobalString(flags.BlsOperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:     ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
		ServerMode:                    ctx.GlobalString(flags.ServerModeFlag.Name),
		PrometheusConfig: prometheus.Config{
			ServerURL: ctx.GlobalString(flags.PrometheusServerURLFlag.Name),
			Username:  ctx.GlobalString(flags.PrometheusServerUsernameFlag.Name),
			Secret:    ctx.GlobalString(flags.PrometheusServerSecretFlag.Name),
			Cluster:   ctx.GlobalString(flags.PrometheusMetricsClusterLabelFlag.Name),
		},
		AllowOrigins: ctx.GlobalStringSlice(flags.AllowOriginsFlag.Name),

		MetricsConfig: dataapi.MetricsConfig{
			HTTPPort:      ctx.GlobalString(flags.MetricsHTTPPort.Name),
			EnableMetrics: ctx.GlobalBool(flags.EnableMetricsFlag.Name),
		},
		DisperserHostname:  ctx.GlobalString(flags.DisperserHostnameFlag.Name),
		ChurnerHostname:    ctx.GlobalString(flags.ChurnerHostnameFlag.Name),
		BatcherHealthEndpt: ctx.GlobalString(flags.BatcherHealthEndptFlag.Name),
		ChainStateConfig:   thegraph.ReadCLIConfig(ctx),
	}
	return config, nil
}
