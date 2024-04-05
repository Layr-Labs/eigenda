package encoder

import (
	"context"
	"fmt"

	"github.com/Layr-Labs/eigenda/disperser"
	pb "github.com/Layr-Labs/eigenda/disperser/api/grpc/encoder"
	"github.com/Layr-Labs/eigenda/encoding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type client struct {
	addr string
}

func NewEncoderClient(addr string) (disperser.EncoderClient, error) {
	return client{
		addr: addr,
	}, nil
}

func (c client) EncodeBlob(ctx context.Context, data []byte, encodingParams encoding.EncodingParams) (*encoding.BlobCommitments, []*encoding.Frame, error) {
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
	chunks := make([]*encoding.Frame, len(reply.GetChunks()))
	for i, chunk := range reply.GetChunks() {
		deserialized, err := new(encoding.Frame).Deserialize(chunk)
		if err != nil {
			return nil, nil, err
		}
		chunks[i] = deserialized
	}
	return &encoding.BlobCommitments{
		Commitment:       commitment,
		LengthCommitment: lengthCommitment,
		LengthProof:      lengthProof,
		Length:           uint(reply.GetCommitment().GetLength()),
	}, chunks, nil
}
