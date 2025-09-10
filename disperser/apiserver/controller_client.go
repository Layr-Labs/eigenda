package apiserver

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	pbcommon "github.com/Layr-Labs/eigenda/api/grpc/common/v2"
	"github.com/Layr-Labs/eigenda/api/grpc/controller"
	"github.com/Layr-Labs/eigenda/api/hashing"
	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ControllerClient wraps the controller service client and handles signing of payment authorization requests
type ControllerClient struct {
	controllerAddress   string
	disperserKMSKeyID   string
	disperserKeyManager *kms.Client
	disperserPublicKey  *ecdsa.PublicKey

	clientConnection *grpc.ClientConn
	serviceClient    controller.ControllerServiceClient
}

// Creates a client for communicating with the controller GRPC server
func NewControllerClient(
	ctx context.Context,
	controllerAddress string,
	kmsRegion string,
	kmsEndpoint string,
	disperserKMSKeyID string,
) (*ControllerClient, error) {
	if controllerAddress == "" {
		return nil, fmt.Errorf("controller address is required")
	}
	if disperserKMSKeyID == "" {
		return nil, fmt.Errorf("disperser KMS key ID is required for controller client")
	}
	if kmsRegion == "" {
		return nil, fmt.Errorf("KMS region is required for controller client")
	}

	awsConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(kmsRegion))
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}

	var keyManager *kms.Client
	if kmsEndpoint != "" {
		keyManager = kms.New(kms.Options{
			Region:       kmsRegion,
			BaseEndpoint: aws.String(kmsEndpoint),
		})
	} else {
		keyManager = kms.NewFromConfig(awsConfig)
	}

	publicKey, err := aws2.LoadPublicKeyKMS(ctx, keyManager, disperserKMSKeyID)
	if err != nil {
		return nil, fmt.Errorf("load public key from KMS: %w", err)
	}

	clientConnection, err := grpc.NewClient(
		controllerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("new grpc client: %w", err)
	}

	client := controller.NewControllerServiceClient(clientConnection)

	return &ControllerClient{
		controllerAddress:   controllerAddress,
		disperserKMSKeyID:   disperserKMSKeyID,
		disperserKeyManager: keyManager,
		disperserPublicKey:  publicKey,
		clientConnection:    clientConnection,
		serviceClient:       client,
	}, nil
}

// Sends a signed payment authorization request to the controller
func (c *ControllerClient) AuthorizePayment(
	ctx context.Context,
	blobHeader *pbcommon.BlobHeader,
) error {
	authorizePaymentRequest := &controller.AuthorizePaymentRequest{BlobHeader: blobHeader}

	hash, err := hashing.HashAuthorizePaymentRequest(authorizePaymentRequest)
	if err != nil {
		return fmt.Errorf("hash authorize payment request: %w", err)
	}

	signature, err := aws2.SignKMS(ctx, c.disperserKeyManager, c.disperserKMSKeyID, c.disperserPublicKey, hash)
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
func (c *ControllerClient) Close() error {
	if c.clientConnection != nil {
		err := c.clientConnection.Close()
		c.clientConnection = nil
		c.serviceClient = nil
		return err
	}
	return nil
}
