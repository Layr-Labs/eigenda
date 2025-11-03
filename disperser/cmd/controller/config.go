package main

import (
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/healthcheck"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/core/payments/ondemand/ondemandvalidation"
	"github.com/Layr-Labs/eigenda/core/payments/reservation/reservationvalidation"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser/cmd/controller/flags"
	"github.com/Layr-Labs/eigenda/disperser/controller"
	"github.com/Layr-Labs/eigenda/disperser/controller/server"
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/urfave/cli"
)

const MaxUint16 = ^uint16(0)

type Config struct {
	EncodingManagerConfig controller.EncodingManagerConfig
	DispatcherConfig      controller.ControllerConfig

	DynamoDBTableName string

	EthClientConfig                     geth.EthClientConfig
	AwsClientConfig                     aws.ClientConfig
	DisperserStoreChunksSigningDisabled bool
	DispersalRequestSignerConfig        clients.DispersalRequestSignerConfig
	LoggerConfig                        common.LoggerConfig
	IndexerConfig                       indexer.Config
	ChainStateConfig                    thegraph.Config
	UseGraph                            bool

	EigenDAContractDirectoryAddress string

	MetricsPort                  int
	ControllerReadinessProbePath string
	ServerConfig                 server.Config
	HeartbeatMonitorConfig       healthcheck.HeartbeatMonitorConfig

	PaymentAuthorizationConfig controller.PaymentAuthorizationConfig
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

	grpcServerConfig, err := common.NewGRPCServerConfig(
		ctx.GlobalBool(flags.GrpcServerEnableFlag.Name),
		uint16(ctx.GlobalUint64(flags.GrpcPortFlag.Name)),
		ctx.GlobalInt(flags.GrpcMaxMessageSizeFlag.Name),
		ctx.GlobalDuration(flags.GrpcMaxIdleConnectionAgeFlag.Name),
		ctx.GlobalDuration(flags.GrpcAuthorizationRequestMaxPastAgeFlag.Name),
		ctx.GlobalDuration(flags.GrpcAuthorizationRequestMaxFutureAgeFlag.Name),
	)
	if err != nil {
		return Config{}, fmt.Errorf("invalid gRPC server config: %w", err)
	}

	serverConfig, err := server.NewConfig(
		grpcServerConfig,
		ctx.GlobalBool(flags.GrpcPaymentAuthenticationFlag.Name),
	)
	if err != nil {
		return Config{}, fmt.Errorf("invalid controller service config: %w", err)
	}

	paymentVaultUpdateInterval := ctx.GlobalDuration(flags.PaymentVaultUpdateIntervalFlag.Name)

	onDemandConfig, err := ondemandvalidation.NewOnDemandLedgerCacheConfig(
		ctx.GlobalInt(flags.OnDemandPaymentsLedgerCacheSizeFlag.Name),
		ctx.GlobalString(flags.OnDemandPaymentsTableNameFlag.Name),
		paymentVaultUpdateInterval,
	)
	if err != nil {
		return Config{}, fmt.Errorf("create on-demand config: %w", err)
	}

	reservationConfig, err := reservationvalidation.NewReservationLedgerCacheConfig(
		ctx.GlobalInt(flags.ReservationPaymentsLedgerCacheSizeFlag.Name),
		// TODO(litt3): once the checkpointed onchain config registry is ready, that should be used
		// instead of hardcoding. At that point, this field will be removed from the config struct
		// entirely, and the value will be fetched dynamically at runtime.
		75*time.Second,
		// this doesn't need to be configurable. there are no plans to ever use a different value
		ratelimit.OverfillOncePermitted,
		paymentVaultUpdateInterval,
	)
	if err != nil {
		return Config{}, fmt.Errorf("create reservation config: %w", err)
	}

	paymentAuthorizationConfig := controller.PaymentAuthorizationConfig{
		OnDemandConfig:    onDemandConfig,
		ReservationConfig: reservationConfig,
	}

	heartbeatMonitorConfig := healthcheck.HeartbeatMonitorConfig{
		FilePath:         ctx.GlobalString(flags.ControllerHealthProbePathFlag.Name),
		MaxStallDuration: ctx.GlobalDuration(flags.ControllerHeartbeatMaxStallDurationFlag.Name),
	}
	if err := heartbeatMonitorConfig.Verify(); err != nil {
		return Config{}, fmt.Errorf("invalid heartbeat monitor config: %w", err)
	}

	awsClientConfig := aws.ReadClientConfig(ctx, flags.FlagPrefix)
	config := Config{
		DynamoDBTableName:                   ctx.GlobalString(flags.DynamoDBTableNameFlag.Name),
		EthClientConfig:                     ethClientConfig,
		AwsClientConfig:                     aws.ReadClientConfig(ctx, flags.FlagPrefix),
		DisperserStoreChunksSigningDisabled: ctx.GlobalBool(flags.DisperserStoreChunksSigningDisabledFlag.Name),
		LoggerConfig:                        *loggerConfig,
		DispersalRequestSignerConfig: clients.DispersalRequestSignerConfig{
			KeyID:      ctx.GlobalString(flags.DisperserKMSKeyIDFlag.Name),
			PrivateKey: ctx.GlobalString(flags.DisperserPrivateKeyFlag.Name),
			Region:     awsClientConfig.Region,
			Endpoint:   awsClientConfig.EndpointURL,
		},
		EncodingManagerConfig: controller.EncodingManagerConfig{
			PullInterval:                ctx.GlobalDuration(flags.EncodingPullIntervalFlag.Name),
			EncodingRequestTimeout:      ctx.GlobalDuration(flags.EncodingRequestTimeoutFlag.Name),
			StoreTimeout:                ctx.GlobalDuration(flags.EncodingStoreTimeoutFlag.Name),
			NumEncodingRetries:          ctx.GlobalInt(flags.NumEncodingRetriesFlag.Name),
			NumRelayAssignment:          uint16(numRelayAssignments),
			AvailableRelays:             relays,
			EncoderAddress:              ctx.GlobalString(flags.EncoderAddressFlag.Name),
			MaxNumBlobsPerIteration:     int32(ctx.GlobalInt(flags.MaxNumBlobsPerIterationFlag.Name)),
			OnchainStateRefreshInterval: ctx.GlobalDuration(flags.OnchainStateRefreshIntervalFlag.Name),
			NumConcurrentRequests:       ctx.GlobalInt(flags.NumConcurrentEncodingRequestsFlag.Name),
		},
		DispatcherConfig: controller.ControllerConfig{
			PullInterval:                          ctx.GlobalDuration(flags.DispatcherPullIntervalFlag.Name),
			FinalizationBlockDelay:                ctx.GlobalUint64(flags.FinalizationBlockDelayFlag.Name),
			AttestationTimeout:                    ctx.GlobalDuration(flags.AttestationTimeoutFlag.Name),
			BatchMetadataUpdatePeriod:             ctx.GlobalDuration(flags.BatchMetadataUpdatePeriodFlag.Name),
			BatchAttestationTimeout:               ctx.GlobalDuration(flags.BatchAttestationTimeoutFlag.Name),
			SignatureTickInterval:                 ctx.GlobalDuration(flags.SignatureTickIntervalFlag.Name),
			NumRequestRetries:                     ctx.GlobalInt(flags.NumRequestRetriesFlag.Name),
			MaxBatchSize:                          int32(ctx.GlobalInt(flags.MaxBatchSizeFlag.Name)),
			SignificantSigningThresholdPercentage: uint8(ctx.GlobalUint(flags.SignificantSigningThresholdPercentageFlag.Name)),
			SignificantSigningMetricsThresholds:   ctx.GlobalStringSlice(flags.SignificantSigningMetricsThresholdsFlag.Name),
			NumConcurrentRequests:                 ctx.GlobalInt(flags.NumConcurrentDispersalRequestsFlag.Name),
			NodeClientCacheSize:                   ctx.GlobalInt(flags.NodeClientCacheNumEntriesFlag.Name),
		},
		IndexerConfig:                   indexer.ReadIndexerConfig(ctx),
		ChainStateConfig:                thegraph.ReadCLIConfig(ctx),
		UseGraph:                        ctx.GlobalBool(flags.UseGraphFlag.Name),
		EigenDAContractDirectoryAddress: ctx.GlobalString(flags.EigenDAContractDirectoryAddressFlag.Name),
		MetricsPort:                     ctx.GlobalInt(flags.MetricsPortFlag.Name),
		ControllerReadinessProbePath:    ctx.GlobalString(flags.ControllerReadinessProbePathFlag.Name),
		ServerConfig:                    serverConfig,
		HeartbeatMonitorConfig:          heartbeatMonitorConfig,
		PaymentAuthorizationConfig:      paymentAuthorizationConfig,
	}

	if err := config.DispersalRequestSignerConfig.Verify(); err != nil {
		return Config{}, fmt.Errorf("invalid dispersal request signer config: %w", err)
	}

	if err := config.EncodingManagerConfig.Verify(); err != nil {
		return Config{}, fmt.Errorf("invalid encoding manager config: %w", err)
	}
	if err := config.DispatcherConfig.Verify(); err != nil {
		return Config{}, fmt.Errorf("invalid dispatcher config: %w", err)
	}
	if err := config.PaymentAuthorizationConfig.Verify(); err != nil {
		return Config{}, fmt.Errorf("invalid payment authorization config: %w", err)
	}

	return config, nil
}
