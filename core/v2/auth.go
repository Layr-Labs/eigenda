package v2

import pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"

type BlobRequestAuthenticator interface {
	AuthenticateBlobRequest(header *BlobHeader) error
	AuthenticatePaymentStateRequest(request *pb.GetPaymentStateRequest) error
}

type BlobRequestSigner interface {
	SignBlobRequest(header *BlobHeader) ([]byte, error)
	SignPaymentStateRequest() ([]byte, error)
	GetAccountID() (string, error)
}
