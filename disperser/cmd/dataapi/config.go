package main

import (
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/logging"
	"github.com/Layr-Labs/eigenda/disperser/blobstore"
	"github.com/Layr-Labs/eigenda/disperser/cmd/dataapi/flags"
	"github.com/Layr-Labs/eigenda/disperser/dataapi"
	"github.com/Layr-Labs/eigenda/disperser/dataapi/prometheus"
	"github.com/urfave/cli"
)

type Config struct {
	AwsClientConfig  aws.ClientConfig
	BlobstoreConfig  blobstore.Config
	EthClientConfig  geth.EthClientConfig
	LoggerConfig     logging.Config
	PrometheusConfig prometheus.Config
	MetricsConfig    dataapi.MetricsConfig

	SocketAddr                   string
	PrometheusApiAddr            string
	SubgraphApiBatchMetadataAddr string
	SubgraphApiOperatorStateAddr string
	ServerMode                   string
	AllowOrigins                 []string

	BLSOperatorStateRetrieverAddr string
	EigenDAServiceManagerAddr     string
}

func NewConfig(ctx *cli.Context) Config {
	config := Config{
		BlobstoreConfig: blobstore.Config{
			BucketName: ctx.GlobalString(flags.S3BucketNameFlag.Name),
			TableName:  ctx.GlobalString(flags.DynamoTableNameFlag.Name),
		},
		AwsClientConfig:               aws.ReadClientConfig(ctx, flags.FlagPrefix),
		EthClientConfig:               geth.ReadEthClientConfig(ctx),
		LoggerConfig:                  logging.ReadCLIConfig(ctx, flags.FlagPrefix),
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
			EnableMetrics: ctx.GlobalBool(flags.EnableMetrics.Name),
		},
	}
	return config
}
