package apiserver

import (
	"context"
	"fmt"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ControllerClient interface {
	AuthorizePayment(ctx context.Context, req *pb.AuthorizePaymentRequest) (*pb.AuthorizePaymentReply, error)
	Close() error
}

type controllerClient struct {
	conn   *grpc.ClientConn
	client pb.ControllerClient
}

func NewControllerClient(address string) (ControllerClient, error) {
	if address == "" {
		return nil, fmt.Errorf("controller address is empty")
	}

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial controller: %w", err)
	}

	return &controllerClient{
		conn:   conn,
		client: pb.NewControllerClient(conn),
	}, nil
}

func (c *controllerClient) AuthorizePayment(ctx context.Context, req *pb.AuthorizePaymentRequest) (*pb.AuthorizePaymentReply, error) {
	return c.client.AuthorizePayment(ctx, req)
}

func (c *controllerClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}