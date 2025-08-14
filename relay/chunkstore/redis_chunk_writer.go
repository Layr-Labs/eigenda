package chunkstore

import (
	"context"
	"crypto/tls"
	"fmt"
	"unsafe"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/valkey-io/valkey-go"
)

var _ ChunkWriter = (*redisChunkWriter)(nil)

// TODO this doesn't belong here, but it won't compile if we try to import it
func UnsafeBytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(&b[0], len(b))
}

// A redis based ChunkWriter.
type redisChunkWriter struct {
	client valkey.Client
}

// Create a new RedisChunkWriter.
func NewRedisChunkWriter(host string, username string, password string) (ChunkWriter, error) {

	opts := valkey.ClientOption{
		InitAddress: []string{host},
		Username:    username,
		Password:    password,
		TLSConfig:   &tls.Config{MinVersion: tls.VersionTLS12},
	}
	client, err := valkey.NewClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %w", err)
	}

	return &redisChunkWriter{
		client: client,
	}, nil
}

// Get a key for storing frame proofs for a particular blob.
func frameProofKey(blobKey corev2.BlobKey) string {
	return "p-" + blobKey.Hex()
}

// Get a key for storing frame coefficients for a particular blob.
func coefficientsKey(blobKey corev2.BlobKey) string {
	return "c-" + blobKey.Hex()
}

func (r *redisChunkWriter) PutFrameProofs(
	ctx context.Context,
	blobKey corev2.BlobKey,
	proofs []*encoding.Proof,
) error {

	key := frameProofKey(blobKey)

	value, err := rs.SerializeFrameProofs(proofs)
	if err != nil {
		return fmt.Errorf("failed to encode proofs: %w", err)
	}

	result := r.client.Do(ctx,
		r.client.B().Set().
			Key(key).
			Value(UnsafeBytesToString(value)).
			ExSeconds(600). // TODO make this configurable
			Build())

	if result.Error() != nil {
		return fmt.Errorf("failed to put frame proofs: %w", result.Error())
	}

	return nil
}

func (r *redisChunkWriter) PutFrameCoefficients(
	ctx context.Context,
	blobKey corev2.BlobKey,
	frames []rs.FrameCoeffs,
) (*encoding.FragmentInfo, error) {

	key := coefficientsKey(blobKey)

	value, err := rs.SerializeFrameCoeffsSlice(frames)
	if err != nil {
		return nil, fmt.Errorf("failed to encode frames: %w", err)
	}

	result := r.client.Do(ctx,
		r.client.B().Set().
			Key(key).
			Value(UnsafeBytesToString(value)).
			ExSeconds(600). // TODO make this configurable
			Build())

	if result.Error() != nil {
		return nil, fmt.Errorf("failed to put frame coefficients: %w", result.Error())
	}

	// fragment info is really only needed for the S3 implementation, so we can just return a default value.

	return &encoding.FragmentInfo{}, nil
}

func (r *redisChunkWriter) ProofExists(
	ctx context.Context,
	blobKey corev2.BlobKey,
) (bool, error) {

	key := frameProofKey(blobKey)

	result := r.client.Do(ctx, r.client.B().Exists().Key(key).Build())
	if result.Error() != nil {
		return false, fmt.Errorf("failed to check if proofs exist: %w", result.Error())
	}

	count, err := result.ToInt64()
	if err != nil {
		return false, fmt.Errorf("failed to convert result to int64: %w", err)
	}

	return count > 0, nil
}

func (r *redisChunkWriter) CoefficientsExists(
	ctx context.Context,
	blobKey corev2.BlobKey,
) (bool, *encoding.FragmentInfo, error) {

	key := coefficientsKey(blobKey)

	result := r.client.Do(ctx, r.client.B().Exists().Key(key).Build())
	if result.Error() != nil {
		return false, nil, fmt.Errorf("failed to check if coefficients exist: %w", result.Error())
	}

	count, err := result.ToInt64()
	if err != nil {
		return false, nil, fmt.Errorf("failed to convert result to int64: %w", err)
	}

	// We don't have fragment info in Redis, so we return a default value.
	return count > 0, &encoding.FragmentInfo{}, nil
}
