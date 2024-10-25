package core

import commonpb "github.com/Layr-Labs/eigenda/api/grpc/common"

type BlobRequestAuthenticator interface {
	AuthenticateBlobRequest(header BlobAuthHeader) error
}

type BlobRequestSigner interface {
	SignBlobRequest(header BlobAuthHeader) ([]byte, error)
	GetAccountID() (string, error)
}

type PaymentSigner interface {
	SignBlobPayment(header *commonpb.PaymentHeader) ([]byte, error)
	GetAccountID() string
}
