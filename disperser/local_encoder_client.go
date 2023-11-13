package disperser

import (
	"context"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
)

type LocalEncoderClient struct {
	mu sync.Mutex

	encoder core.Encoder
}

var _ EncoderClient = (*LocalEncoderClient)(nil)

func NewLocalEncoderClient(encoder core.Encoder) *LocalEncoderClient {
	return &LocalEncoderClient{
		encoder: encoder,
	}
}

func (m *LocalEncoderClient) EncodeBlob(ctx context.Context, data []byte, encodingParams core.EncodingParams) (*core.BlobCommitments, []*core.Chunk, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	commits, chunks, err := m.encoder.Encode(data, encodingParams)
	if err != nil {
		return nil, nil, err
	}

	return &commits, chunks, nil
}
