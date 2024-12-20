package clients

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api"
	"sync"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	nodegrpc "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"google.golang.org/grpc"
)

type NodeClientConfig struct {
	Hostname          string
	Port              string
	UseSecureGrpcFlag bool
}

type NodeClient interface {
	StoreChunks(ctx context.Context, certs *corev2.Batch) (*core.Signature, error)
	Close() error
}

type nodeClient struct {
	config        *NodeClientConfig
	initOnce      sync.Once
	conn          *grpc.ClientConn
	requestSigner DispersalRequestSigner

	dispersalClient nodegrpc.DispersalClient
}

var _ NodeClient = (*nodeClient)(nil)

func NewNodeClient(config *NodeClientConfig, requestSigner DispersalRequestSigner) (NodeClient, error) {
	if config == nil || config.Hostname == "" || config.Port == "" {
		return nil, fmt.Errorf("invalid config: %v", config)
	}
	return &nodeClient{
		config:        config,
		requestSigner: requestSigner,
	}, nil
}

func (c *nodeClient) StoreChunks(ctx context.Context, batch *corev2.Batch) (*core.Signature, error) {
	if len(batch.BlobCertificates) == 0 {
		return nil, fmt.Errorf("no blob certificates in the batch")
	}

	if err := c.initOnceGrpcConnection(); err != nil {
		return nil, err
	}

	blobCerts := make([]*commonpb.BlobCertificate, len(batch.BlobCertificates))
	for i, cert := range batch.BlobCertificates {
		var err error
		blobCerts[i], err = cert.ToProtobuf()
		if err != nil {
			return nil, fmt.Errorf("failed to convert blob certificate to protobuf: %v", err)
		}
	}

	request := &nodegrpc.StoreChunksRequest{
		Batch: &commonpb.Batch{
			Header: &commonpb.BatchHeader{
				BatchRoot:            batch.BatchHeader.BatchRoot[:],
				ReferenceBlockNumber: batch.BatchHeader.ReferenceBlockNumber,
			},
			BlobCertificates: blobCerts,
		},
		DisperserID: api.EigenLabsDisperserID, // this will need to be updated when dispersers are decentralized
	}

	if c.requestSigner != nil {
		// Sign the request to store chunks
		signature, err := c.requestSigner.SignStoreChunksRequest(ctx, request)
		if err != nil {
			return nil, fmt.Errorf("failed to sign store chunks request: %v", err)
		}
		request.Signature = signature
	}

	// Call the gRPC method to store chunks
	response, err := c.dispersalClient.StoreChunks(ctx, request)
	if err != nil {
		return nil, err
	}

	// Extract signatures from the response
	if response == nil {
		return nil, fmt.Errorf("received nil response from StoreChunks")
	}

	sigBytes := response.GetSignature()
	point, err := new(core.Signature).Deserialize(sigBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize signature: %v", err)
	}
	return &core.Signature{G1Point: point}, nil
}

// Close closes the grpc connection to the disperser server.
// It is thread safe and can be called multiple times.
func (c *nodeClient) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		c.dispersalClient = nil
		return err
	}
	return nil
}

func (c *nodeClient) initOnceGrpcConnection() error {
	var initErr error
	c.initOnce.Do(func() {
		addr := fmt.Sprintf("%v:%v", c.config.Hostname, c.config.Port)
		dialOptions := getGrpcDialOptions(c.config.UseSecureGrpcFlag)
		conn, err := grpc.NewClient(addr, dialOptions...)
		if err != nil {
			initErr = err
			return
		}
		c.conn = conn
		c.dispersalClient = nodegrpc.NewDispersalClient(conn)
	})
	return initErr
}
