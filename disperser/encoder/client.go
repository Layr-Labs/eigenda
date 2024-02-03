package encoder

import (
	"context"
	"fmt"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/disperser"
	pb "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder"
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

func (c client) EncodeBlob(ctx context.Context, data []byte, encodingParams core.EncodingParams) (*core.BlobCommitments, []*core.Chunk, error) {
	conn, err := grpc.Dial(
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

	commitment, err := new(core.Commitment).Deserialize(reply.GetCommitment().GetCommitment())
	if err != nil {
		return nil, nil, err
	}
	lengthCommitment, err := new(core.LengthCommitment).Deserialize(reply.GetCommitment().GetLengthCommitment())
	if err != nil {
		return nil, nil, err
	}
	lengthProof, err := new(core.LengthProof).Deserialize(reply.GetCommitment().GetLengthProof())
	if err != nil {
		return nil, nil, err
	}
	chunks := make([]*core.Chunk, len(reply.GetChunks()))
	for i, chunk := range reply.GetChunks() {
		deserialized, err := new(core.Chunk).Deserialize(chunk)
		if err != nil {
			return nil, nil, err
		}
		chunks[i] = deserialized
	}
	return &core.BlobCommitments{
		Commitment:       commitment,
		LengthCommitment: lengthCommitment,
		LengthProof:      lengthProof,
		Length:           uint(reply.GetCommitment().GetLength()),
	}, chunks, nil
}
