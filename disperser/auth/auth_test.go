package auth

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"testing"
)

var _ v2.DispersalServer = (*mockDispersalServer)(nil)

type mockDispersalServer struct {
	v2.DispersalServer
}

func (s *mockDispersalServer) StoreChunks(context.Context, *v2.StoreChunksRequest) (*v2.StoreChunksReply, error) {
	fmt.Printf("called StoreChunks\n")
	return nil, nil
}

func (s *mockDispersalServer) NodeInfo(context.Context, *v2.NodeInfoRequest) (*v2.NodeInfoReply, error) {
	fmt.Printf("called NodeInfo\n")
	return nil, nil
}

func buildClient(t *testing.T) v2.DispersalClient {
	addr := "0.0.0.0:50051"

	options := make([]grpc.DialOption, 0)
	//options = append(options, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(addr, options...)
	require.NoError(t, err)

	return v2.NewDispersalClient(conn)
}

func buildServer(t *testing.T) (v2.DispersalServer, *grpc.Server) {
	dispersalServer := &mockDispersalServer{}

	options := make([]grpc.ServerOption, 0)
	server := grpc.NewServer(options...)
	v2.RegisterDispersalServer(server, dispersalServer)

	addr := "0.0.0.0:50051"
	listener, err := net.Listen("tcp", addr)
	require.NoError(t, err)

	go func() {
		err = server.Serve(listener)
		require.NoError(t, err)
	}()

	return dispersalServer, server
}

func TestServerWithTLS(t *testing.T) {
	dispersalServer, server := buildServer(t)
	defer server.Stop()
	require.NotNil(t, dispersalServer) // TODO remove

	client := buildClient(t)

	response, err := client.NodeInfo(context.Background(), &v2.NodeInfoRequest{})
	require.NoError(t, err)
	require.NotNil(t, response)
}
