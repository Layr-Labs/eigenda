package main

import (
	"fmt"
	"math"
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
	"github.com/Layr-Labs/eigenda/indexer"
	"github.com/urfave/cli"
)

func NewConfig(ctx *cli.Context) (*controller.ControllerConfig, error) {
	loggerConfig, err := common.ReadLoggerCLIConfig(ctx, flags.FlagPrefix)
	if err != nil {
		return nil, fmt.Errorf("read logger config: %w", err)
	}
	ethClientConfig := geth.ReadEthClientConfigRPCOnly(ctx)
	numRelayAssignments := ctx.GlobalInt(flags.NumRelayAssignmentFlag.Name)
	if numRelayAssignments < 1 || numRelayAssignments > math.MaxUint16 {
		return nil, fmt.Errorf("invalid number of relay assignments: %d", numRelayAssignments)
	}
	availableRelays := ctx.GlobalIntSlice(flags.AvailableRelaysFlag.Name)
	if len(availableRelays) == 0 {
		return nil, fmt.Errorf("no available relays specified")
	}
	relays := make([]corev2.RelayKey, len(availableRelays))
	for i, relay := range availableRelays {
		if relay < 0 || relay > 65_535 {
			return nil, fmt.Errorf("invalid relay: %d", relay)
		}
		relays[i] = corev2.RelayKey(relay)
	}

	grpcServerConfig, err := common.NewGRPCServerConfig(
		uint16(ctx.GlobalUint64(flags.GrpcPortFlag.Name)),
		ctx.GlobalInt(flags.GrpcMaxMessageSizeFlag.Name),
		ctx.GlobalDuration(flags.GrpcMaxIdleConnectionAgeFlag.Name),
		ctx.GlobalDuration(flags.GrpcAuthorizationRequestMaxPastAgeFlag.Name),
		ctx.GlobalDuration(flags.GrpcAuthorizationRequestMaxFutureAgeFlag.Name),
	)
	if err != nil {
		return nil, fmt.Errorf("invalid gRPC server config: %w", err)
	}

	paymentVaultUpdateInterval := ctx.GlobalDuration(flags.PaymentVaultUpdateIntervalFlag.Name)

	onDemandConfig, err := ondemandvalidation.NewOnDemandLedgerCacheConfig(
		ctx.GlobalInt(flags.OnDemandPaymentsLedgerCacheSizeFlag.Name),
		ctx.GlobalString(flags.OnDemandPaymentsTableNameFlag.Name),
		paymentVaultUpdateInterval,
	)
	if err != nil {
		return nil, fmt.Errorf("create on-demand config: %w", err)
	}

	reservationConfig, err := reservationvalidation.NewReservationLedgerCacheConfig(
		ctx.GlobalInt(flags.ReservationPaymentsLedgerCacheSizeFlag.Name),
		// TODO(litt3): once the checkpointed onchain config registry is ready, that should be used
		// instead of hardcoding. At that point, this field will be removed from the config struct
		// entirely, and the value will be fetched dynamically at runtime.
		90*time.Second,
		// this doesn't need to be configurable. there are no plans to ever use a different value
		ratelimit.OverfillOncePermitted,
		paymentVaultUpdateInterval,
	)
	if err != nil {
		return nil, fmt.Errorf("create reservation config: %w", err)
	}

	paymentAuthorizationConfig := controller.PaymentAuthorizationConfig{
		OnDemandConfig:                 onDemandConfig,
		ReservationConfig:              reservationConfig,
		EnablePerAccountPaymentMetrics: ctx.GlobalBool(flags.EnablePerAccountPaymentMetricsFlag.Name),
	}

	heartbeatMonitorConfig := healthcheck.HeartbeatMonitorConfig{
		FilePath:         ctx.GlobalString(flags.ControllerHealthProbePathFlag.Name),
		MaxStallDuration: ctx.GlobalDuration(flags.ControllerHeartbeatMaxStallDurationFlag.Name),
	}
	if err := heartbeatMonitorConfig.Verify(); err != nil {
		return nil, fmt.Errorf("invalid heartbeat monitor config: %w", err)
	}

	awsClientConfig := aws.ReadClientConfig(ctx, flags.FlagPrefix)
	disperserID := uint32(ctx.GlobalUint64(flags.DisperserIDFlag.Name))
	config := &controller.ControllerConfig{
		DynamoDBTableName:                   ctx.GlobalString(flags.DynamoDBTableNameFlag.Name),
		DisperserID:                         disperserID,
		EthClientConfig:                     ethClientConfig,
		AwsClient:                           aws.ReadClientConfig(ctx, flags.FlagPrefix),
		DisperserStoreChunksSigningDisabled: ctx.GlobalBool(flags.DisperserStoreChunksSigningDisabledFlag.Name),
		Logger:                              *loggerConfig,
		DispersalRequestSigner: clients.DispersalRequestSignerConfig{
			KeyID:      ctx.GlobalString(flags.DisperserKMSKeyIDFlag.Name),
			PrivateKey: ctx.GlobalString(flags.DisperserPrivateKeyFlag.Name),
			Region:     awsClientConfig.Region,
			Endpoint:   awsClientConfig.EndpointURL,
		},
		EncodingManager: controller.EncodingManagerConfig{
			PullInterval:                      ctx.GlobalDuration(flags.EncodingPullIntervalFlag.Name),
			EncodingRequestTimeout:            ctx.GlobalDuration(flags.EncodingRequestTimeoutFlag.Name),
			StoreTimeout:                      ctx.GlobalDuration(flags.EncodingStoreTimeoutFlag.Name),
			NumEncodingRetries:                ctx.GlobalInt(flags.NumEncodingRetriesFlag.Name),
			NumRelayAssignment:                uint16(numRelayAssignments),
			AvailableRelays:                   relays,
			EncoderAddress:                    ctx.GlobalString(flags.EncoderAddressFlag.Name),
			MaxNumBlobsPerIteration:           int32(ctx.GlobalInt(flags.MaxNumBlobsPerIterationFlag.Name)),
			OnchainStateRefreshInterval:       ctx.GlobalDuration(flags.OnchainStateRefreshIntervalFlag.Name),
			NumConcurrentRequests:             ctx.GlobalInt(flags.NumConcurrentEncodingRequestsFlag.Name),
			EnablePerAccountBlobStatusMetrics: ctx.GlobalBool(flags.EnablePerAccountBlobStatusMetricsFlag.Name),
		},
		PullInterval:                           ctx.GlobalDuration(flags.DispatcherPullIntervalFlag.Name),
		FinalizationBlockDelay:                 ctx.GlobalUint64(flags.FinalizationBlockDelayFlag.Name),
		AttestationTimeout:                     ctx.GlobalDuration(flags.AttestationTimeoutFlag.Name),
		BatchMetadataUpdatePeriod:              ctx.GlobalDuration(flags.BatchMetadataUpdatePeriodFlag.Name),
		BatchAttestationTimeout:                ctx.GlobalDuration(flags.BatchAttestationTimeoutFlag.Name),
		SignatureTickInterval:                  ctx.GlobalDuration(flags.SignatureTickIntervalFlag.Name),
		MaxBatchSize:                           int32(ctx.GlobalInt(flags.MaxBatchSizeFlag.Name)),
		SignificantSigningThresholdFraction:    ctx.GlobalFloat64(flags.SignificantSigningThresholdFractionFlag.Name),
		NumConcurrentRequests:                  ctx.GlobalInt(flags.NumConcurrentDispersalRequestsFlag.Name),
		NodeClientCacheSize:                    ctx.GlobalInt(flags.NodeClientCacheNumEntriesFlag.Name),
		CollectDetailedValidatorSigningMetrics: ctx.GlobalBool(flags.DetailedValidatorMetricsFlag.Name),
		EnablePerAccountBlobStatusMetrics:      ctx.GlobalBool(flags.EnablePerAccountBlobStatusMetricsFlag.Name),
		MaxDispersalAge:                        ctx.GlobalDuration(flags.MaxDispersalAgeFlag.Name),
		MaxDispersalFutureAge:                  ctx.GlobalDuration(flags.MaxDispersalFutureAgeFlag.Name),
		SigningRateRetentionPeriod:             ctx.GlobalDuration(flags.SigningRateRetentionPeriodFlag.Name),
		SigningRateBucketSpan:                  ctx.GlobalDuration(flags.SigningRateBucketSpanFlag.Name),
		BlobDispersalQueueSize:                 uint32(ctx.GlobalUint64(flags.BlobDispersalQueueSizeFlag.Name)),
		BlobDispersalRequestBatchSize:          uint32(ctx.GlobalUint64(flags.BlobDispersalRequestBatchSizeFlag.Name)),
		BlobDispersalRequestBackoffPeriod:      ctx.GlobalDuration(flags.BlobDispersalRequestBackoffPeriodFlag.Name),
		SigningRateFlushPeriod:                 ctx.GlobalDuration(flags.SigningRateFlushPeriodFlag.Name),
		SigningRateDynamoDbTableName:           ctx.GlobalString(flags.SigningRateDynamoDbTableNameFlag.Name),
		Indexer:                                indexer.ReadIndexerConfig(ctx),
		ChainState:                             thegraph.ReadCLIConfig(ctx),
		UseGraph:                               ctx.GlobalBool(flags.UseGraphFlag.Name),
		EigenDAContractDirectoryAddress:        ctx.GlobalString(flags.EigenDAContractDirectoryAddressFlag.Name),
		MetricsPort:                            ctx.GlobalInt(flags.MetricsPortFlag.Name),
		ControllerReadinessProbePath:           ctx.GlobalString(flags.ControllerReadinessProbePathFlag.Name),
		Server:                                 grpcServerConfig,
		HeartbeatMonitor:                       heartbeatMonitorConfig,
		PaymentAuthorization:                   paymentAuthorizationConfig,
		UserAccountRemappingFilePath:           ctx.GlobalString(flags.UserAccountRemappingFileFlag.Name),
		ValidatorIdRemappingFilePath:           ctx.GlobalString(flags.ValidatorIdRemappingFileFlag.Name),
	}

	err = config.Verify()
	if err != nil {
		return nil, fmt.Errorf("verify controller config: %w", err)
	}

	return config, nil
}
