package lib

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/api/grpc/validator"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/common/math"
	"github.com/Layr-Labs/eigenda/common/ratelimit"
	"github.com/Layr-Labs/eigenda/common/store"
	authv2 "github.com/Layr-Labs/eigenda/core/auth/v2"
	"github.com/Layr-Labs/eigenda/core/eth"
	mt "github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/Layr-Labs/eigenda/core/signingrate"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/disperser/apiserver"
	"github.com/Layr-Labs/eigenda/disperser/common/blobstore"
	blobstorev2 "github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/encoding/v2/kzg/committer"
	gethcommon "github.com/ethereum/go-ethereum/common"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunDisperserServer(ctx *cli.Context) error {
	config, err := NewConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := common.NewLogger(&config.LoggerConfig)
	if err != nil {
		return err
	}

	client, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		logger.Error("Cannot create chain.Client", "err", err)
		return err
	}

	chainId, err := client.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("get chain ID: %w", err)
	}
	config.ServerConfig.ChainId = chainId

	transactor, err := eth.NewReader(
		logger, client, config.OperatorStateRetrieverAddr, config.EigenDAServiceManagerAddr)
	if err != nil {
		return err
	}
	blockStaleMeasure, err := transactor.GetBlockStaleMeasure(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get BLOCK_STALE_MEASURE: %w", err)
	}
	storeDurationBlocks, err := transactor.GetStoreDurationBlocks(context.Background())
	if err != nil || storeDurationBlocks == 0 {
		return fmt.Errorf("failed to get STORE_DURATION_BLOCKS: %w", err)
	}

	objectStorageClient, err := blobstore.CreateObjectStorageClient(
		context.Background(), config.BlobstoreConfig, config.AwsClientConfig, logger)
	if err != nil {
		return err
	}

	dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
	if err != nil {
		return err
	}

	reg := prometheus.NewRegistry()

	var meterer *mt.Meterer
	if config.EnablePaymentMeterer {
		mtConfig := mt.Config{
			ChainReadTimeout: config.ChainReadTimeout,
			UpdateInterval:   config.OnchainStateRefreshInterval,
		}

		paymentChainState, err := mt.NewOnchainPaymentState(context.Background(), transactor, logger)
		if err != nil {
			return fmt.Errorf("failed to create onchain payment state: %w", err)
		}
		if err := paymentChainState.RefreshOnchainPaymentState(context.Background()); err != nil {
			return fmt.Errorf("failed to make initial query to the on-chain state: %w", err)
		}

		meteringStore, err := mt.NewDynamoDBMeteringStore(
			config.AwsClientConfig,
			config.ReservationsTableName,
			config.OnDemandTableName,
			config.GlobalRateTableName,
			logger,
		)
		if err != nil {
			return fmt.Errorf("failed to create offchain store: %w", err)
		}
		// add some default sensible configs
		meterer = mt.NewMeterer(
			mtConfig,
			paymentChainState,
			meteringStore,
			logger,
			// metrics.NewNoopMetrics(),
		)
		meterer.Start(context.Background())
	}

	var ratelimiter common.RateLimiter
	if config.EnableRatelimiter {
		globalParams := config.RatelimiterConfig.GlobalRateParams

		var bucketStore common.KVStore[common.RateBucketParams]
		if config.BucketTableName != "" {
			dynamoClient, err := dynamodb.NewClient(config.AwsClientConfig, logger)
			if err != nil {
				return err
			}
			bucketStore = store.NewDynamoParamStore[common.RateBucketParams](dynamoClient, config.BucketTableName)
		} else {
			bucketStore, err = store.NewLocalParamStore[common.RateBucketParams](config.BucketStoreSize)
			if err != nil {
				return err
			}
		}
		ratelimiter = ratelimit.NewRateLimiter(reg, globalParams, bucketStore, logger)
	}

	if config.MaxBlobSize <= 0 || config.MaxBlobSize > 32*1024*1024 {
		return fmt.Errorf("configured max blob size is invalid %v", config.MaxBlobSize)
	}

	if !math.IsPowerOfTwo(uint64(config.MaxBlobSize)) {
		return fmt.Errorf("configured max blob size must be power of 2 %v", config.MaxBlobSize)
	}

	bucketName := config.BlobstoreConfig.BucketName
	logger.Info("Blob store", "bucket", bucketName)
	if config.DisperserVersion == V2 {
		committer, err := committer.NewFromConfig(config.KzgCommitterConfig)
		if err != nil {
			return fmt.Errorf("new committer: %w", err)
		}
		baseBlobMetadataStore := blobstorev2.NewBlobMetadataStore(
			dynamoClient,
			logger,
			config.BlobstoreConfig.TableName)
		blobMetadataStore := blobstorev2.NewInstrumentedMetadataStore(
			baseBlobMetadataStore,
			blobstorev2.InstrumentedMetadataStoreConfig{
				ServiceName: "apiserver",
				Registry:    reg,
				Backend:     blobstorev2.BackendDynamoDB,
			})
		blobStore := blobstorev2.NewBlobStore(bucketName, objectStorageClient, logger)

		var controllerConnection *grpc.ClientConn
		var controllerClient controller.ControllerServiceClient
		if config.UseControllerMediatedPayments {
			if config.ControllerAddress == "" {
				return fmt.Errorf("controller address is required when using controller-mediated payments")
			}
			connection, err := grpc.NewClient(
				config.ControllerAddress,
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			)
			if err != nil {
				return fmt.Errorf("create controller connection: %w", err)
			}
			controllerConnection = connection
			controllerClient = controller.NewControllerServiceClient(connection)
		}

		// Create listener for the gRPC server
		addr := fmt.Sprintf("%s:%s", "0.0.0.0", config.ServerConfig.GrpcPort)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			return fmt.Errorf("failed to create listener: %w", err)
		}

		signingRateTracker, err := signingrate.NewSigningRateTracker(
			logger,
			config.ServerConfig.SigningRateRetentionPeriod,
			time.Second, // bucket size is unimportant, since it is unused when mirroring from controller
			time.Now,
		)
		if err != nil {
			return fmt.Errorf("failed to create signing rate tracker: %w", err)
		}
		signingRateTracker = signingrate.NewThreadsafeSigningRateTracker(context.Background(), signingRateTracker)

		// A function that can fetch signing rate data from the controller.
		scraper := func(ctx context.Context, startTime time.Time) ([]*validator.SigningRateBucket, error) {
			data, err := controllerClient.GetValidatorSigningRateDump(
				ctx,
				&controller.GetValidatorSigningRateDumpRequest{
					StartTimestamp: uint64(startTime.Unix()),
				})
			if err != nil {
				return nil, fmt.Errorf("GetValidatorSigningRateDump RPC failed: %w", err)
			}
			return data.GetSigningRateBuckets(), nil
		}

		// Clone signing rate data from controller. This is blocking, so that when we start the server we have
		// data to serve right away.
		err = signingrate.DoInitialScrape(
			context.Background(),
			logger,
			scraper,
			signingRateTracker,
			config.ServerConfig.SigningRateRetentionPeriod)
		if err != nil {
			return fmt.Errorf("do initial scrape: %w", err)
		}

		// In the background, periodically refresh signing rate data from controller.
		go signingrate.MirrorSigningRate(
			context.Background(),
			logger,
			scraper,
			signingRateTracker,
			config.ServerConfig.SigningRatePollInterval,
			config.ServerConfig.SigningRateRetentionPeriod,
		)

		server, err := apiserver.NewDispersalServerV2(
			config.ServerConfig,
			time.Now,
			blobStore,
			blobMetadataStore,
			transactor,
			meterer,
			authv2.NewPaymentStateAuthenticator(
				config.AuthPmtStateRequestMaxPastAge,
				config.AuthPmtStateRequestMaxFutureAge),
			committer,
			config.MaxNumSymbolsPerBlob,
			config.OnchainStateRefreshInterval,
			config.MaxDispersalAge,
			config.MaxFutureDispersalTime,
			logger,
			reg,
			config.MetricsConfig,
			config.ReservedOnly,
			config.UseControllerMediatedPayments,
			controllerConnection,
			controllerClient,
			listener,
			signingRateTracker,
		)
		if err != nil {
			return err
		}
		return server.Start(context.Background())
	}

	blobMetadataStore := blobstore.NewBlobMetadataStore(
		dynamoClient,
		logger,
		config.BlobstoreConfig.TableName,
		time.Duration((storeDurationBlocks+blockStaleMeasure)*12)*time.Second)
	blobStore := blobstore.NewSharedStorage(bucketName, objectStorageClient, blobMetadataStore, logger)

	grpcMetrics := grpcprom.NewServerMetrics()
	metrics := disperser.NewMetrics(reg, config.MetricsConfig.HTTPPort, logger)
	server := apiserver.NewDispersalServer(
		config.ServerConfig,
		blobStore,
		transactor,
		logger,
		metrics,
		grpcMetrics,
		meterer,
		ratelimiter,
		config.RateConfig,
		config.MaxBlobSize,
	)

	reg.MustRegister(grpcMetrics)

	// Enable Metrics Block
	if config.MetricsConfig.EnableMetrics {
		httpSocket := fmt.Sprintf(":%s", config.MetricsConfig.HTTPPort)
		metrics.Start(context.Background())
		logger.Info("Enabled metrics for Disperser", "socket", httpSocket)
	}

	return server.Start(context.Background())
}
