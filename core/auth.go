package core

type BlobRequestAuthenticator interface {
	AuthenticateBlobRequest(header BlobAuthHeader) error
}

type BlobRequestSigner interface {
	SignBlobRequest(header BlobAuthHeader) ([]byte, error)
	GetAccountID() (string, error)
}
