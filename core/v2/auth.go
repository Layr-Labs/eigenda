package v2

type BlobRequestAuthenticator interface {
	AuthenticateBlobRequest(header *BlobHeader, signature []byte) error
	AuthenticatePaymentStateRequest(signature []byte, accountId string) error
}

type BlobRequestSigner interface {
	SignBlobRequest(header *BlobHeader) ([]byte, error)
	SignPaymentStateRequest() ([]byte, error)
	GetAccountID() (string, error)
}
