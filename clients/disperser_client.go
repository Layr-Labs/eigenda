package clients

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Hostname          string
	Port              string
	Timeout           time.Duration
	UseSecureGrpcFlag bool
}

func NewConfig(hostname, port string, timeout time.Duration, useSecureGrpcFlag bool) *Config {
	return &Config{
		Hostname:          hostname,
		Port:              port,
		Timeout:           timeout,
		UseSecureGrpcFlag: useSecureGrpcFlag,
	}
}

type DisperserClient interface {
	DisperseBlob(ctx context.Context, data []byte, customQuorums []uint8) (*disperser.BlobStatus, []byte, error)
	DisperseBlobAuthenticated(ctx context.Context, data []byte, customQuorums []uint8) (*disperser.BlobStatus, []byte, error)
	GetBlobStatus(ctx context.Context, key []byte) (*disperser_rpc.BlobStatusReply, error)
}

type disperserClient struct {
	config *Config
	signer core.BlobRequestSigner
}

var _ DisperserClient = &disperserClient{}

func NewDisperserClient(config *Config, signer core.BlobRequestSigner) DisperserClient {
	return &disperserClient{
		config: config,
		signer: signer,
	}
}

func (c *disperserClient) getDialOptions() []grpc.DialOption {
	if c.config.UseSecureGrpcFlag {
		config := &tls.Config{}
		credential := credentials.NewTLS(config)
		return []grpc.DialOption{grpc.WithTransportCredentials(credential)}
	} else {
		return []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	}
}

func (c *disperserClient) DisperseBlob(ctx context.Context, data []byte, quorums []uint8) (*disperser.BlobStatus, []byte, error) {
	addr := fmt.Sprintf("%v:%v", c.config.Hostname, c.config.Port)

	dialOptions := c.getDialOptions()
	conn, err := grpc.Dial(addr, dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = conn.Close() }()

	disperserClient := disperser_rpc.NewDisperserClient(conn)
	ctxTimeout, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	quorumNumbers := make([]uint32, len(quorums))
	for i, q := range quorums {
		quorumNumbers[i] = uint32(q)
	}

	request := &disperser_rpc.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: quorumNumbers,
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

func (c *disperserClient) DisperseBlobAuthenticated(ctx context.Context, data []byte, quorums []uint8) (*disperser.BlobStatus, []byte, error) {

	addr := fmt.Sprintf("%v:%v", c.config.Hostname, c.config.Port)

	dialOptions := c.getDialOptions()
	conn, err := grpc.Dial(addr, dialOptions...)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = conn.Close() }()

	disperserClient := disperser_rpc.NewDisperserClient(conn)
	ctxTimeout, cancel := context.WithTimeout(ctx, c.config.Timeout)

	defer cancel()

	stream, err := disperserClient.DisperseBlobAuthenticated(ctxTimeout)
	if err != nil {
		return nil, nil, fmt.Errorf("frror while calling DisperseBlobAuthenticated: %v", err)
	}

	quorumNumbers := make([]uint32, len(quorums))
	for i, q := range quorums {
		quorumNumbers[i] = uint32(q)
	}

	request := &disperser_rpc.DisperseBlobRequest{
		Data:                data,
		CustomQuorumNumbers: quorumNumbers,
		AccountId:           c.signer.GetAccountID(),
	}

	// Send the initial request
	err = stream.Send(&disperser_rpc.AuthenticatedRequest{Payload: &disperser_rpc.AuthenticatedRequest_DisperseRequest{
		DisperseRequest: request,
	}})

	if err != nil {
		return nil, nil, fmt.Errorf("failed to send request: %v", err)
	}

	// Get the Challenge
	reply, err := stream.Recv()
	if err != nil {
		return nil, nil, fmt.Errorf("error while receiving: %v", err)
	}
	authHeaderReply, ok := reply.Payload.(*disperser_rpc.AuthenticatedReply_BlobAuthHeader)
	if !ok {
		return nil, nil, fmt.Errorf("expected challenge")
	}

	authHeader := core.BlobAuthHeader{
		BlobCommitments: encoding.BlobCommitments{},
		AccountID:       "",
		Nonce:           authHeaderReply.BlobAuthHeader.ChallengeParameter,
	}

	authData, err := c.signer.SignBlobRequest(authHeader)
	if err != nil {
		return nil, nil, fmt.Errorf("error signing blob request")
	}

	// Process challenge and send back challenge_reply
	err = stream.Send(&disperser_rpc.AuthenticatedRequest{Payload: &disperser_rpc.AuthenticatedRequest_AuthenticationData{
		AuthenticationData: &disperser_rpc.AuthenticationData{
			AuthenticationData: authData,
		},
	}})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to send challenge reply: %v", err)
	}

	reply, err = stream.Recv()
	if err != nil {
		return nil, nil, fmt.Errorf("error while receiving final reply: %v", err)
	}
	disperseReply, ok := reply.Payload.(*disperser_rpc.AuthenticatedReply_DisperseReply) // Process the final disperse_reply
	if !ok {
		return nil, nil, fmt.Errorf("expected DisperseReply")
	}

	blobStatus, err := disperser.FromBlobStatusProto(disperseReply.DisperseReply.GetResult())
	if err != nil {
		return nil, nil, err
	}

	return blobStatus, disperseReply.DisperseReply.GetRequestId(), nil
}

func (c *disperserClient) GetBlobStatus(ctx context.Context, requestID []byte) (*disperser_rpc.BlobStatusReply, error) {
	addr := fmt.Sprintf("%v:%v", c.config.Hostname, c.config.Port)
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
