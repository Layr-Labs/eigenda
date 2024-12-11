package clients

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/Layr-Labs/eigenda/disperser/auth"
	"sync"

	commonpb "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	nodegrpc "github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/Layr-Labs/eigenda/core"
	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"google.golang.org/grpc"
)

type NodeClientV2Config struct {
	Hostname          string
	Port              string
	UseSecureGrpcFlag bool
	// The .pem file containing the private key used to sign StoreChunks() requests. If "" then no signing is done.
	PrivateKeyFile string
}

type NodeClientV2 interface {
	StoreChunks(ctx context.Context, certs *corev2.Batch) (*core.Signature, error)
	Close() error
}

type nodeClientV2 struct {
	config   *NodeClientV2Config
	initOnce sync.Once
	conn     *grpc.ClientConn
	key      *ecdsa.PrivateKey

	dispersalClient nodegrpc.DispersalClient
}

var _ NodeClientV2 = (*nodeClientV2)(nil)

func NewNodeClientV2(config *NodeClientV2Config) (*nodeClientV2, error) {
	if config == nil || config.Hostname == "" || config.Port == "" {
		return nil, fmt.Errorf("invalid config: %v", config)
	}

	var key *ecdsa.PrivateKey // TODO update flags
	if config.PrivateKeyFile != "" {
		var err error
		key, err = auth.ReadPrivateECDSAKeyFile(config.PrivateKeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %v", err)
		}
	}

	return &nodeClientV2{
		config: config,
		key:    key,
	}, nil
}

func (c *nodeClientV2) StoreChunks(ctx context.Context, batch *corev2.Batch) (*core.Signature, error) {
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
	}

	if c.key != nil {
		signature, err := auth.SignStoreChunksRequest(c.key, request) // TODO
		if err != nil {
			return nil, fmt.Errorf("failed to sign request: %v", err)
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
func (c *nodeClientV2) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		c.dispersalClient = nil
		return err
	}
	return nil
}

func (c *nodeClientV2) initOnceGrpcConnection() error {
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
