package mock

import (
	"context"

	disperser_rpc "github.com/Layr-Labs/eigenda/api/grpc/disperser"
)

// Currently only implements the RetrieveBlob RPC
type DisperserServer struct {
	disperser_rpc.UnimplementedDisperserServer
}

// RetrieveBlob returns a ~5MiB(+header_size) blob. It is used to test that the client correctly sets the max message size,
// to be able to support large blobs (default grpc max message size is 4MiB).
func (m *DisperserServer) RetrieveBlob(ctx context.Context, req *disperser_rpc.RetrieveBlobRequest) (*disperser_rpc.RetrieveBlobReply, error) {
	// Create a blob larger than default max size (4MiB)
	largeBlob := make([]byte, 5*1024*1024) // 5MiB
	for i := range largeBlob {
		largeBlob[i] = byte(i % 256)
	}

	return &disperser_rpc.RetrieveBlobReply{
		Data: largeBlob,
	}, nil
}
