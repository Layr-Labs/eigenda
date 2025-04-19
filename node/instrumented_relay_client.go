package node

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
)

// RelayConnectionManager wraps the RelayClient to provide instrumented connections
type RelayConnectionManager struct {
	relayClient clients.RelayClient
	metrics     *Metrics
	logger      logging.Logger
	useSecure   bool
	maxMsgSize  uint
}

// NewInstrumentedRelayClient creates a new RelayClient with instrumented connections
func NewInstrumentedRelayClient(
	config *clients.RelayClientConfig,
	logger logging.Logger,
	relayUrlProvider relay.RelayUrlProvider,
	metrics *Metrics,
) (clients.RelayClient, error) {
	// Create the base relay client
	client, err := clients.NewRelayClient(config, logger, relayUrlProvider)
	if err != nil {
		return nil, err
	}

	// Create and return the wrapped client
	return &RelayConnectionManager{
		relayClient: client,
		metrics:     metrics,
		logger:      logger.With("component", "InstrumentedRelayClient"),
		useSecure:   config.UseSecureGrpcFlag,
		maxMsgSize:  config.MaxGRPCMessageSize,
	}, nil
}

// GetBlob retrieves a blob from a relay
func (r *RelayConnectionManager) GetBlob(
	ctx context.Context,
	relayKey corev2.RelayKey,
	blobKey corev2.BlobKey,
) ([]byte, error) {
	connectionID := fmt.Sprintf("relay-%d", relayKey)
	queueStart := time.Now()
	result, err := r.relayClient.GetBlob(ctx, relayKey, blobKey)
	if r.metrics != nil {
		queueLatency := time.Since(queueStart)
		r.metrics.RecordQueueLatency(connectionID, "GetBlob", queueLatency)
	}
	return result, err
}

// GetChunksByRange retrieves blob chunks from a relay by chunk index range
func (r *RelayConnectionManager) GetChunksByRange(
	ctx context.Context,
	relayKey corev2.RelayKey,
	requests []*clients.ChunkRequestByRange,
) ([][]byte, error) {
	// This method is now directly instrumented in the DownloadBundles function
	// for more granular measurements of each specific download operation
	return r.relayClient.GetChunksByRange(ctx, relayKey, requests)
}

// GetChunksByIndex retrieves blob chunks from a relay by index
func (r *RelayConnectionManager) GetChunksByIndex(
	ctx context.Context,
	relayKey corev2.RelayKey,
	requests []*clients.ChunkRequestByIndex,
) ([][]byte, error) {
	connectionID := fmt.Sprintf("relay-%d", relayKey)
	queueStart := time.Now()
	result, err := r.relayClient.GetChunksByIndex(ctx, relayKey, requests)
	if r.metrics != nil {
		queueLatency := time.Since(queueStart)
		r.metrics.RecordQueueLatency(connectionID, "GetChunksByIndex", queueLatency)
	}
	return result, err
}

// Close closes the relay client
func (r *RelayConnectionManager) Close() error {
	return r.relayClient.Close()
}

// getDialOptions returns connection-specific dial options with instrumentation
func (r *RelayConnectionManager) getDialOptions(relayKey corev2.RelayKey) []grpc.DialOption {
	connectionID := fmt.Sprintf("relay-%d", relayKey)
	return GetDialOptions(r.metrics, r.logger, r.useSecure, r.maxMsgSize, connectionID)
}
