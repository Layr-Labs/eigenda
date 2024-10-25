package v2

type BlobRequestAuthenticator interface {
	AuthenticateBlobRequest(header *BlobHeader) error
}

type BlobRequestSigner interface {
	SignBlobRequest(header *BlobHeader) ([]byte, error)
	GetAccountID() (string, error)
}
