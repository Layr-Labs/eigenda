package clients

import (
	"context"
	"fmt"
	relaygrpc "github.com/Layr-Labs/eigenda/api/grpc/relay"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"google.golang.org/grpc"
)

var _ relaygrpc.RelayClient = (*limitedRelayClient)(nil)

// TODO (cody-littley / litt3): when randomly selecting a relay client to use, we should avoid selecting a
// client that currently has exhausted its concurrency permits.

// limitedRelayClient encapsulates a gRPC client and a channel for limiting concurrent requests to that client.
type limitedRelayClient struct {
	// client is the underlying gRPC client.
	client relaygrpc.RelayClient

	// relayID is the ID of the relay that this client is connected to.
	relayKey corev2.RelayKey

	// permits is a channel for limiting the number of concurrent requests to the client.
	// when a request is initiated, a value is sent to the channel, and when the request is completed,
	// the value is received. The channel has a buffer size of `MaxConcurrentRequests`, and will block
	// if the number of concurrent requests exceeds this limit.
	permits chan struct{}
}

// newLimitedRelayClient creates a new limitedRelayClient.
func newLimitedRelayClient(
	client relaygrpc.RelayClient,
	relayKey corev2.RelayKey,
	maxConcurrentRequests uint) (*limitedRelayClient, error) {

	if maxConcurrentRequests == 0 {
		return nil, fmt.Errorf("maxConcurrentRequests must be greater than 0")
	}

	return &limitedRelayClient{
		client:   client,
		relayKey: relayKey,
		permits:  make(chan struct{}, maxConcurrentRequests),
	}, nil
}

func (l *limitedRelayClient) GetBlob(
	ctx context.Context,
	in *relaygrpc.GetBlobRequest,
	opts ...grpc.CallOption) (*relaygrpc.GetBlobReply, error) {

	select {
	case l.permits <- struct{}{}:
		// permit acquired
	case <-ctx.Done():
		return nil,
			fmt.Errorf("context cancelled while waiting for permit to get blob from relay %d", l.relayKey)
	}
	defer func() {
		<-l.permits
	}()
	return l.client.GetBlob(ctx, in, opts...)
}

func (l *limitedRelayClient) GetChunks(
	ctx context.Context,
	in *relaygrpc.GetChunksRequest,
	opts ...grpc.CallOption) (*relaygrpc.GetChunksReply, error) {

	select {
	case l.permits <- struct{}{}:
		// permit acquired
	case <-ctx.Done():
		return nil,
			fmt.Errorf("context cancelled while waiting for permit to get chunks from relay %d", l.relayKey)
	}
	defer func() {
		<-l.permits
	}()
	return l.client.GetChunks(ctx, in, opts...)
}
