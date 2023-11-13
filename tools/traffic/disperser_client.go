package traffic

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/disperser"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type DisperserClient interface {
	DisperseBlob(ctx context.Context, data []byte, quorumID, quorumThreshold, adversityThreshold uint8) (*disperser.BlobStatus, []byte, error)
	GetBlobStatus(ctx context.Context, key []byte) (*disperser_rpc.BlobStatusReply, error)
}

type client struct {
	config *Config
}

func NewDisperserClient(config *Config) DisperserClient {
	return &client{
		config: config,
	}
}

func (c *client) getDialOptions() []grpc.DialOption {
	if c.config.UseSecureGrpcFlag {
		config := &tls.Config{}
		credential := credentials.NewTLS(config)
		return []grpc.DialOption{grpc.WithTransportCredentials(credential)}
	} else {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}
}

func (c *client) DisperseBlob(ctx context.Context, data []byte, quorumID, quorumThreshold, adversityThreshold uint8) (*disperser.BlobStatus, []byte, error) {
	addr := fmt.Sprintf("%v:%v", c.config.Hostname, c.config.GrpcPort)

	dialOptions := c.getDialOptions()
	conn, err := grpc.Dial(addr, dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = conn.Close() }()

	disperserClient := disperser_rpc.NewDisperserClient(conn)
	ctxTimeout, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	request := &disperser_rpc.DisperseBlobRequest{
		Data: data,
		SecurityParams: []*disperser_rpc.SecurityParams{
			{
				QuorumId:           uint32(quorumID),
				AdversaryThreshold: uint32(adversityThreshold),
				QuorumThreshold:    uint32(quorumThreshold),
			},
		},
	}

	reply, err := disperserClient.DisperseBlob(ctxTimeout, request)
	if err != nil {
		return nil, nil, err
	}

	blobStatus, err := disperser.FromBlobStatusProto(reply.GetResult())
	if err != nil {
		return nil, nil, err
	}

	return blobStatus, reply.GetRequestId(), nil
}

func (c *client) GetBlobStatus(ctx context.Context, requestID []byte) (*disperser_rpc.BlobStatusReply, error) {
	addr := fmt.Sprintf("%v:%v", c.config.Hostname, c.config.GrpcPort)
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	disperserClient := disperser_rpc.NewDisperserClient(conn)
	ctxTimeout, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()

	request := &disperser_rpc.BlobStatusRequest{
		RequestId: requestID,
	}

	reply, err := disperserClient.GetBlobStatus(ctxTimeout, request)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
