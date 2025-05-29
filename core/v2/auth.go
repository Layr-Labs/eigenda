package v2

import (
	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser/v2"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type BlobRequestAuthenticator interface {
	AuthenticateBlobRequest(header *BlobHeader, signature []byte) error
	AuthenticatePaymentStateRequest(accountId gethcommon.Address, request *pb.GetPaymentStateRequest) error
	AuthenticateQuorumSpecificPaymentStateRequest(accountId gethcommon.Address, request *pb.GetQuorumSpecificPaymentStateRequest) error
}

type BlobRequestSigner interface {
	SignBlobRequest(header *BlobHeader) ([]byte, error)
	SignPaymentStateRequest(timestamp uint64) ([]byte, error)
	GetAccountID() (gethcommon.Address, error)
}
