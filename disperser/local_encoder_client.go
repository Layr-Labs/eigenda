package disperser

import (
	"context"
	"sync"

	"github.com/Layr-Labs/eigenda/encoding"
)

type LocalEncoderClient struct {
	mu sync.Mutex

	prover encoding.Prover
}

var _ EncoderClient = (*LocalEncoderClient)(nil)

func NewLocalEncoderClient(prover encoding.Prover) *LocalEncoderClient {
	return &LocalEncoderClient{
		prover: prover,
	}
}

func (m *LocalEncoderClient) EncodeBlob(ctx context.Context, data []byte, encodingParams encoding.EncodingParams) (*encoding.BlobCommitments, []*encoding.Frame, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	commits, chunks, err := m.prover.EncodeAndProveDataAsEvals(data, encodingParams)
	if err != nil {
		return nil, nil, err
	}

	return &commits, chunks, nil
}
