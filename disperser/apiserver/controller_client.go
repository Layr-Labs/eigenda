package apiserver

import (
	"context"
	"fmt"

	pb "github.com/Layr-Labs/eigenda/api/grpc/controller"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ControllerClient struct {
	clientConnection *grpc.ClientConn
	client           pb.ControllerServiceClient
}

func NewControllerClient(address string) (*ControllerClient, error) {
	if address == "" {
		return nil, fmt.Errorf("controller address is empty")
	}

	clientConnection, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial controller: %w", err)
	}

	return &ControllerClient{
		clientConnection: clientConnection,
		client:           pb.NewControllerServiceClient(clientConnection),
	}, nil
}

func (cc *ControllerClient) AuthorizePayment(
	ctx context.Context,
	authorizePaymentRequest *pb.AuthorizePaymentRequest,
) (*pb.AuthorizePaymentReply, error) {
	return cc.client.AuthorizePayment(ctx, authorizePaymentRequest)
}

func (cc *ControllerClient) Close() error {
	if cc.clientConnection != nil {
		return cc.clientConnection.Close()
	}
	return nil
}
