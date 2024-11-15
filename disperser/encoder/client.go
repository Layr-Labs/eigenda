package encoder

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	pb "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	addr    string
	timeout time.Duration
}

func NewEncoderClient(addr string, timeout time.Duration) (disperser.EncoderClient, error) {
	return client{
		addr:    addr,
		timeout: timeout,
	}, nil
}

func (c client) EncodeBlob(ctx context.Context, data []byte, encodingParams encoding.EncodingParams) (*encoding.BlobCommitments, *core.ChunksData, error) {
	conn, err := grpc.NewClient(
		c.addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*1024)), // 1 GiB
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial encoder: %w", err)
	}
	defer conn.Close()

	encoder := pb.NewEncoderClient(conn)
	reply, err := encoder.EncodeBlob(ctx, &pb.EncodeBlobRequest{
		Data: data,
		EncodingParams: &pb.EncodingParams{
			ChunkLength: uint32(encodingParams.ChunkLength),
			NumChunks:   uint32(encodingParams.NumChunks),
		},
	})
	if err != nil {
		return nil, nil, err
	}

	commitment, err := new(encoding.G1Commitment).Deserialize(reply.GetCommitment().GetCommitment())
	if err != nil {
		return nil, nil, err
	}
	lengthCommitment, err := new(encoding.G2Commitment).Deserialize(reply.GetCommitment().GetLengthCommitment())
	if err != nil {
		return nil, nil, err
	}
	lengthProof, err := new(encoding.LengthProof).Deserialize(reply.GetCommitment().GetLengthProof())
	if err != nil {
		return nil, nil, err
	}
	var format core.ChunkEncodingFormat
	switch reply.GetChunkEncodingFormat() {
	case pb.ChunkEncodingFormat_GNARK:
		format = core.GnarkChunkEncodingFormat
	case pb.ChunkEncodingFormat_GOB:
		format = core.GobChunkEncodingFormat
	case pb.ChunkEncodingFormat_UNKNOWN:
		format = core.GobChunkEncodingFormat
	}
	chunksData := &core.ChunksData{
		Chunks:   reply.GetChunks(),
		Format:   format,
		ChunkLen: int(encodingParams.ChunkLength),
	}
	return &encoding.BlobCommitments{
		Commitment:       commitment,
		LengthCommitment: lengthCommitment,
		LengthProof:      lengthProof,
		Length:           uint(reply.GetCommitment().GetLength()),
	}, chunksData, nil
}
