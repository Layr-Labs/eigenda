package v2

import gethcommon "github.com/ethereum/go-ethereum/common"

type BlobRequestAuthenticator interface {
	AuthenticateBlobRequest(header *BlobHeader, signature []byte) error
	AuthenticatePaymentStateRequest(signature []byte, accountId gethcommon.Address) error
}

type BlobRequestSigner interface {
	SignBlobRequest(header *BlobHeader) ([]byte, error)
	SignPaymentStateRequest() ([]byte, error)
	GetAccountID() (gethcommon.Address, error)
}
