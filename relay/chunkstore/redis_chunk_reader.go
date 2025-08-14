package chunkstore

import (
	"context"
	"crypto/tls"
	"fmt"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/valkey-io/valkey-go"
)

var _ ChunkReader = (*redisChunkReader)(nil)

// A redis implementation of ChunkReader.
type redisChunkReader struct {
	client valkey.Client
}

// NewRedisChunkReader creates a new RedisChunkReader.
func NewRedisChunkReader(host string, username string, password string) (ChunkReader, error) {

	opts := valkey.ClientOption{
		InitAddress:  []string{host},
		Username:     username,
		Password:     password,
		TLSConfig:    &tls.Config{MinVersion: tls.VersionTLS12},
		DisableCache: true,
	}
	client, err := valkey.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %w", err)
	}

	return &redisChunkReader{
		client: client,
	}, nil
}

func (r *redisChunkReader) GetBinaryChunkProofs(
	ctx context.Context,
	blobKey corev2.BlobKey,
) ([][]byte, error) {

	key := frameProofKey(blobKey)

	result := r.client.Do(ctx, r.client.B().Get().Key(key).Build())
	if result.Error() != nil {
		return nil, fmt.Errorf("failed to get frame proofs: %w", result.Error())
	}

	value, err := result.AsBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to decode frame proofs: %w", err) // TODO handle case where not found
	}

	proofs, err := rs.SplitSerializedFrameProofs(value)
	if err != nil {
		return nil, fmt.Errorf("failed to split serialized frame proofs: %w", err)
	}

	return proofs, nil
}

func (r *redisChunkReader) GetBinaryChunkCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
	fragmentInfo *encoding.FragmentInfo,
) (uint32, [][]byte, error) {
	key := coefficientsKey(blobKey)

	result := r.client.Do(ctx, r.client.B().Get().Key(key).Build())
	if result.Error() != nil {
		return 0, nil, fmt.Errorf("failed to get coefficients: %w", result.Error())
	}

	value, err := result.AsBytes()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to decode coefficients: %w", err) // TODO handle case where not found
	}

	elementCount, frames, err := rs.SplitSerializedFrameCoeffs(value)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to split serialized frame coefficients: %w", err)
	}

	return elementCount, frames, nil
}
