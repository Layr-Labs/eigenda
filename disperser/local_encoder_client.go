package disperser

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
)

type LocalEncoderClient struct {
	mu sync.Mutex

	prover *prover.Prover
}

var _ EncoderClient = (*LocalEncoderClient)(nil)

func NewLocalEncoderClient(prover *prover.Prover) *LocalEncoderClient {
	return &LocalEncoderClient{
		prover: prover,
	}
}

func (m *LocalEncoderClient) EncodeBlob(ctx context.Context, data []byte, encodingParams encoding.EncodingParams) (*encoding.BlobCommitments, *core.ChunksData, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	commits, chunks, err := m.prover.EncodeAndProve(data, encodingParams)
	if err != nil {
		return nil, nil, fmt.Errorf("prover.EncodeAndProve: %w", err)
	}

	bytes := make([][]byte, 0, len(chunks))
	for _, c := range chunks {
		serialized, err := c.SerializeGob()
		if err != nil {
			return nil, nil, fmt.Errorf("serialize chunk: %w", err)
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
