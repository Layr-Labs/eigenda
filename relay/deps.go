package relay

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/aws/dynamodb"
	"github.com/Layr-Labs/eigenda/common/aws/s3"
	"github.com/Layr-Labs/eigenda/common/geth"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/eth"
	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/Layr-Labs/eigenda/disperser/common/v2/blobstore"
	"github.com/Layr-Labs/eigenda/relay/chunkstore"
	"github.com/Layr-Labs/eigensdk-go/logging"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
)

// ServerDependencies holds all the dependencies needed to create a relay Server.
type ServerDependencies struct {
	MetricsRegistry *prometheus.Registry
	MetadataStore   blobstore.MetadataStore
	BlobStore       *blobstore.BlobStore
	ChunkReader     chunkstore.ChunkReader
	ChainReader     core.Reader
	ChainState      core.IndexedChainState
	Logger          logging.Logger
}

// ServerDependenciesConfig contains only the configuration values for setting up relay server dependencies.
type ServerDependenciesConfig struct {
	AWSConfig                  aws.ClientConfig
	MetadataTableName          string
	BucketName                 string
	OperatorStateRetrieverAddr string
	ServiceManagerAddr         string
	ChainStateConfig           thegraph.Config
	EthClientConfig            geth.EthClientConfig
	LoggerConfig               common.LoggerConfig
}

// NewServerDependencies creates all the dependencies needed for a relay server.
// This function is shared between the CLI entrypoint and test harnesses.
func NewServerDependencies(
	ctx context.Context,
	config ServerDependenciesConfig,
) (*ServerDependencies, error) {
	logger, err := common.NewLogger(&config.LoggerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	ethClient, err := geth.NewMultiHomingClient(config.EthClientConfig, gethcommon.Address{}, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create eth client: %w", err)
	}

	// Create DynamoDB client
	dynamoClient, err := dynamodb.NewClient(config.AWSConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamodb client: %w", err)
	}

	// Create S3 client
	s3Client, err := s3.NewClient(ctx, config.AWSConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 client: %w", err)
	}

	// Create metrics registry
	metricsRegistry := prometheus.NewRegistry()

	// Create metadata store
	baseMetadataStore := blobstore.NewBlobMetadataStore(dynamoClient, logger, config.MetadataTableName)
	metadataStore := blobstore.NewInstrumentedMetadataStore(baseMetadataStore, blobstore.InstrumentedMetadataStoreConfig{
		ServiceName: "relay",
		Registry:    metricsRegistry,
		Backend:     blobstore.BackendDynamoDB,
	})

	// Create blob store and chunk reader
	blobStore := blobstore.NewBlobStore(config.BucketName, s3Client, logger)
	chunkReader := chunkstore.NewChunkReader(logger, s3Client, config.BucketName)

	// Create eth writer
	tx, err := eth.NewWriter(logger, ethClient, config.OperatorStateRetrieverAddr, config.ServiceManagerAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to create eth writer: %w", err)
	}

	// Create chain state
	cs := eth.NewChainState(tx, ethClient)
	ics := thegraph.MakeIndexedChainState(config.ChainStateConfig, cs, logger)

	return &ServerDependencies{
		MetricsRegistry: metricsRegistry,
		MetadataStore:   metadataStore,
		BlobStore:       blobStore,
		ChunkReader:     chunkReader,
		ChainReader:     tx,
		ChainState:      ics,
		Logger:          logger,
	}, nil
}
