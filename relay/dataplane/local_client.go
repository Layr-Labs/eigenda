package dataplane

import "github.com/Layr-Labs/eigenda/common/kvstore"

var _ S3Client = &localClient{}

// localClient implements the S3Client interface, but the data is stored locally using a kvstore.Store.
// This may be useful for testing, but is not intended for production use.
type localClient struct {
	store kvstore.Store
}

// NewLocalClient creates a new S3Client instance that stores data locally in the provided store
func NewLocalClient(store kvstore.Store) S3Client {
	return &localClient{
		store: store,
	}
}

func (l *localClient) Upload(key string, data []byte, fragmentSize int) error {
	return l.store.Put([]byte(key), data)
}

func (l *localClient) Download(key string, fileSize int, fragmentSize int) ([]byte, error) {
	return l.store.Get([]byte(key))
}

func (l *localClient) Close() error {
	return l.store.Shutdown()
}
