package node

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/stats"
)

// StreamTracker implements the stats.Handler interface to track gRPC stream metrics
type StreamTracker struct {
	// Maps a connection (transport) to the number of active streams
	streamsPerConnection     map[string]int
	streamsPerConnectionLock sync.Mutex
	metrics                  *Metrics
	logger                   logging.Logger
	connectionID             string
}

// NewStreamTracker creates a new StreamTracker for monitoring gRPC stream usage
func NewStreamTracker(metrics *Metrics, logger logging.Logger, connectionID string) *StreamTracker {
	return &StreamTracker{
		streamsPerConnection: make(map[string]int),
		metrics:              metrics,
		logger:               logger.With("component", "StreamTracker"),
		connectionID:         connectionID,
	}
}

// TagConn implements stats.Handler
func (s *StreamTracker) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

// TagRPC implements stats.Handler
func (s *StreamTracker) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

// HandleRPC implements stats.Handler
func (s *StreamTracker) HandleRPC(ctx context.Context, stat stats.RPCStats) {
	switch stat.(type) {
	case *stats.Begin:
		s.streamsPerConnectionLock.Lock()
		s.streamsPerConnection[s.connectionID]++
		count := s.streamsPerConnection[s.connectionID]
		s.streamsPerConnectionLock.Unlock()
		s.metrics.RecordStreamsInFlight(s.connectionID, count)
	case *stats.End:
		s.streamsPerConnectionLock.Lock()
		s.streamsPerConnection[s.connectionID]--
		count := s.streamsPerConnection[s.connectionID]
		s.streamsPerConnectionLock.Unlock()
		s.metrics.RecordStreamsInFlight(s.connectionID, count)
	}
}

// HandleConn implements stats.Handler
func (s *StreamTracker) HandleConn(ctx context.Context, stat stats.ConnStats) {
	switch stat.(type) {
	case *stats.ConnBegin:
		// Initialize the connection counter at 0
		s.streamsPerConnectionLock.Lock()
		s.streamsPerConnection[s.connectionID] = 0
		s.streamsPerConnectionLock.Unlock()
		s.metrics.RecordStreamsInFlight(s.connectionID, 0)
	case *stats.ConnEnd:
		// Clean up when connection is closed
		s.streamsPerConnectionLock.Lock()
		delete(s.streamsPerConnection, s.connectionID)
		s.streamsPerConnectionLock.Unlock()
	}
}

// UnaryClientInterceptor returns a gRPC client-side interceptor that measures queue time and latency
func UnaryClientInterceptor(metrics *Metrics, connectionID string) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		queueStart := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		queueLatency := time.Since(queueStart)
		metrics.RecordQueueLatency(connectionID, method, queueLatency)
		return err
	}
}

// StreamClientInterceptor returns a gRPC client-side stream interceptor that measures queue time
func StreamClientInterceptor(metrics *Metrics, connectionID string) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		queueStart := time.Now()
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		queueLatency := time.Since(queueStart)
		metrics.RecordQueueLatency(connectionID, method, queueLatency)
		return clientStream, err
	}
}

// GetDialOptions returns the gRPC dial options configured with metric interceptors and stream trackers
func GetDialOptions(metrics *Metrics, logger logging.Logger, useSecure bool, maxMsgSize uint, connectionID string) []grpc.DialOption {
	// Create a unique connection ID if none provided
	if connectionID == "" {
		connectionID = fmt.Sprintf("conn-%d", time.Now().UnixNano())
	}

	// Base options
	opts := []grpc.DialOption{
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(int(maxMsgSize))),
	}

	// Add secure or insecure transport credentials
	if useSecure {
		// Not implemented - would add TLS credentials
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Add interceptors for metrics
	streamTracker := NewStreamTracker(metrics, logger, connectionID)
	opts = append(opts,
		grpc.WithStatsHandler(streamTracker),
		grpc.WithUnaryInterceptor(UnaryClientInterceptor(metrics, connectionID)),
		grpc.WithStreamInterceptor(StreamClientInterceptor(metrics, connectionID)),
	)

	return opts
}
