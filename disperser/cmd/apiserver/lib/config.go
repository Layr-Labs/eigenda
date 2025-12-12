package lib

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/Layr-Labs/eigenda/disperser/cmd/apiserver/flags"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	"github.com/urfave/cli"
)

type DisperserVersion uint

const (
	V1 DisperserVersion = 1
	V2 DisperserVersion = 2
)

type Config struct {
	DisperserVersion  DisperserVersion
	AwsClientConfig   aws.ClientConfig
	BlobstoreConfig   blobstore.Config
	ServerConfig      disperser.ServerConfig
	LoggerConfig      common.LoggerConfig
	MetricsConfig     disperser.MetricsConfig
	RatelimiterConfig ratelimit.Config
	RateConfig        apiserver.RateConfig
	// KzgCommitterConfig is only needed when DisperserVersion is V2.
	// It's used by the grpc endpoint we expose to compute client commitments.
	KzgCommitterConfig            committer.Config
	EnableRatelimiter             bool
	EnablePaymentMeterer          bool
	ReservedOnly                  bool
	ChainReadTimeout              time.Duration
	ReservationsTableName         string
	OnDemandTableName             string
	GlobalRateTableName           string
	BucketTableName               string
	BucketStoreSize               int
	EthClientConfig               geth.EthClientConfig
	MaxBlobSize                   int
	MaxNumSymbolsPerBlob          uint32
	OnchainStateRefreshInterval   time.Duration
	ControllerAddress             string
	UseControllerMediatedPayments bool

	EigenDADirectory                string
	OperatorStateRetrieverAddr      string
	EigenDAServiceManagerAddr       string
	AuthPmtStateRequestMaxPastAge   time.Duration
	AuthPmtStateRequestMaxFutureAge time.Duration
	MaxDispersalAge                 time.Duration
	MaxFutureDispersalTime          time.Duration
}

func NewConfig(ctx *cli.Context) (Config, error) {
	version := ctx.GlobalUint(flags.DisperserVersionFlag.Name)
	if version != uint(V1) && version != uint(V2) {
		return Config{}, fmt.Errorf("unknown disperser version %d", version)
	}

	ratelimiterConfig, err := ratelimit.ReadCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}

	rateConfig, err := apiserver.ReadCLIConfig(ctx)
	if err != nil {
		return Config{}, err
	}

	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return Config{}, err
	}

	kzgCommitterConfig := committer.ReadCLIConfig(ctx)
	if version == uint(V2) {
		if err := kzgCommitterConfig.Verify(); err != nil {
			return Config{}, fmt.Errorf("disperser version 2: kzg committer config verify: %w", err)
		}
	}

	config := Config{
		DisperserVersion: DisperserVersion(version),
		AwsClientConfig:  aws.ReadClientConfig(ctx, flags.FlagPrefix),
		ServerConfig: disperser.ServerConfig{
			GrpcPort:                       ctx.GlobalString(flags.GrpcPortFlag.Name),
			GrpcTimeout:                    ctx.GlobalDuration(flags.GrpcTimeoutFlag.Name),
			MaxConnectionAge:               ctx.GlobalDuration(flags.MaxConnectionAgeFlag.Name),
			MaxConnectionAgeGrace:          ctx.GlobalDuration(flags.MaxConnectionAgeGraceFlag.Name),
			MaxIdleConnectionAge:           ctx.GlobalDuration(flags.MaxIdleConnectionAgeFlag.Name),
			PprofHttpPort:                  ctx.GlobalString(flags.PprofHttpPort.Name),
			EnablePprof:                    ctx.GlobalBool(flags.EnablePprof.Name),
			DisableGetBlobCommitment:       ctx.GlobalBool(flags.DisableGetBlobCommitment.Name),
			SigningRateRetentionPeriod:     ctx.GlobalDuration(flags.SigningRateRetentionPeriodFlag.Name),
			SigningRatePollInterval:        ctx.GlobalDuration(flags.SigningRatePollIntervalFlag.Name),
			DisperserId:                    uint32(ctx.GlobalUint64(flags.DisperserIdFlag.Name)),
			TolerateMissingAnchorSignature: ctx.GlobalBool(flags.TolerateMissingAnchorSignatureFlag.Name),
		},
		BlobstoreConfig: blobstore.Config{
			BucketName:       ctx.GlobalString(flags.S3BucketNameFlag.Name),
			TableName:        ctx.GlobalString(flags.DynamoDBTableNameFlag.Name),
			Backend:          blobstore.ObjectStorageBackend(ctx.GlobalString(flags.ObjectStorageBackendFlag.Name)),
			OCIRegion:        ctx.GlobalString(flags.OCIRegionFlag.Name),
			OCICompartmentID: ctx.GlobalString(flags.OCICompartmentIDFlag.Name),
			OCINamespace:     ctx.GlobalString(flags.OCINamespaceFlag.Name),
		},
		LoggerConfig: *loggerConfig,
		MetricsConfig: disperser.MetricsConfig{
			HTTPPort:                 ctx.GlobalString(flags.MetricsHTTPPort.Name),
			EnableMetrics:            ctx.GlobalBool(flags.EnableMetrics.Name),
			DisablePerAccountMetrics: ctx.GlobalBool(flags.DisablePerAccountMetricsFlag.Name),
		},
		RatelimiterConfig:             ratelimiterConfig,
		RateConfig:                    rateConfig,
		KzgCommitterConfig:            kzgCommitterConfig,
		EnableRatelimiter:             ctx.GlobalBool(flags.EnableRatelimiter.Name),
		EnablePaymentMeterer:          ctx.GlobalBool(flags.EnablePaymentMeterer.Name),
		ReservedOnly:                  ctx.GlobalBoolT(flags.ReservedOnly.Name),
		ControllerAddress:             ctx.GlobalString(flags.ControllerAddressFlag.Name),
		UseControllerMediatedPayments: ctx.GlobalBool(flags.UseControllerMediatedPayments.Name),
		ReservationsTableName:         ctx.GlobalString(flags.ReservationsTableName.Name),
		OnDemandTableName:             ctx.GlobalString(flags.OnDemandTableName.Name),
		GlobalRateTableName:           ctx.GlobalString(flags.GlobalRateTableName.Name),
		BucketTableName:               ctx.GlobalString(flags.BucketTableName.Name),
		BucketStoreSize:               ctx.GlobalInt(flags.BucketStoreSize.Name),
		ChainReadTimeout:              ctx.GlobalDuration(flags.ChainReadTimeout.Name),
		EthClientConfig:               geth.ReadEthClientConfigRPCOnly(ctx),
		MaxBlobSize:                   ctx.GlobalInt(flags.MaxBlobSize.Name),
		MaxNumSymbolsPerBlob:          uint32(ctx.GlobalUint(flags.MaxNumSymbolsPerBlob.Name)),
		OnchainStateRefreshInterval:   ctx.GlobalDuration(flags.OnchainStateRefreshInterval.Name),

		EigenDADirectory:                ctx.GlobalString(flags.EigenDADirectoryFlag.Name),
		OperatorStateRetrieverAddr:      ctx.GlobalString(flags.OperatorStateRetrieverFlag.Name),
		EigenDAServiceManagerAddr:       ctx.GlobalString(flags.EigenDAServiceManagerFlag.Name),
		AuthPmtStateRequestMaxPastAge:   ctx.GlobalDuration(flags.AuthPmtStateRequestMaxPastAge.Name),
		AuthPmtStateRequestMaxFutureAge: ctx.GlobalDuration(flags.AuthPmtStateRequestMaxFutureAge.Name),
		MaxDispersalAge:                 ctx.GlobalDuration(flags.MaxDispersalAgeFlag.Name),
		MaxFutureDispersalTime:          ctx.GlobalDuration(flags.MaxFutureDispersalTimeFlag.Name),
	}
	return config, nil
}
