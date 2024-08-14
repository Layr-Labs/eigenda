package disperser

import (
	"context"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
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

func (m *LocalEncoderClient) EncodeBlob(ctx context.Context, data []byte, encodingParams encoding.EncodingParams) (*encoding.BlobCommitments, *core.ChunksData, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	commits, chunks, err := m.prover.EncodeAndProve(data, encodingParams)
	if err != nil {
		return nil, nil, err
	}

	bytes := make([][]byte, 0, len(chunks))
	for _, c := range chunks {
		serialized, err := c.Serialize()
		if err != nil {
			return nil, nil, err
		}
		bytes = append(bytes, serialized)
	}
	chunksData := &core.ChunksData{
		Chunks:   bytes,
		Format:   core.GobChunkEncodingFormat,
		ChunkLen: int(encodingParams.ChunkLength),
	}

	return &commits, chunksData, nil
}
