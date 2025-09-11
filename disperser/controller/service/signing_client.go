package service

import (
	"context"
	"fmt"

	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	controller "github.com/Layr-Labs/eigenda/api/grpc/controller/v1"
	"github.com/Layr-Labs/eigenda/api/hashing"
	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Wraps the controller service client and handles signing of requests
type SigningClient struct {
	controllerAddress string
	kmsSigner         *aws2.KMSSigner

	clientConnection *grpc.ClientConn
	serviceClient    controller.ControllerServiceClient
}

// Creates a client for communicating with the controller GRPC server
func NewSigningClient(
	ctx context.Context,
	controllerAddress string,
	kmsSigner *aws2.KMSSigner,
) (*SigningClient, error) {
	if controllerAddress == "" {
		return nil, fmt.Errorf("controller address is required")
	}
	if kmsSigner == nil {
		return nil, fmt.Errorf("KMS signer is required")
	}

	clientConnection, err := grpc.NewClient(
		controllerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("new grpc client: %w", err)
	}

	serviceClient := controller.NewControllerServiceClient(clientConnection)

	return &SigningClient{
		controllerAddress: controllerAddress,
		kmsSigner:         kmsSigner,
		clientConnection:  clientConnection,
		serviceClient:     serviceClient,
	}, nil
}

// Sends a signed payment authorization request to the controller
func (c *SigningClient) AuthorizePayment(ctx context.Context, blobHeader *pbcommon.BlobHeader) error {
	authorizePaymentRequest := &controller.AuthorizePaymentRequest{BlobHeader: blobHeader}

	hash, err := hashing.HashAuthorizePaymentRequest(authorizePaymentRequest)
	if err != nil {
		return fmt.Errorf("hash authorize payment request: %w", err)
	}

	signature, err := c.kmsSigner.Sign(ctx, hash)
	if err != nil {
		return fmt.Errorf("sign authorization request: %w", err)
	}

	authorizePaymentRequest.DisperserSignature = signature

	_, err = c.serviceClient.AuthorizePayment(ctx, authorizePaymentRequest)
	if err != nil {
		return fmt.Errorf("authorize payment: %w", err)
	}

	return nil
}

// Closes the grpc connection to the controller server
func (c *SigningClient) Close() error {
	if c.clientConnection != nil {
		err := c.clientConnection.Close()
		if err != nil {
			return fmt.Errorf("close connection: %w", err)
		}

		c.clientConnection = nil
		c.serviceClient = nil
	}
	return nil
}
