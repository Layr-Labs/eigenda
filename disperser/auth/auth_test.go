package auth

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/grpc/node/v2"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"testing"
)

// The purpose of these tests are to verify that TLS key generation works as expected.
// TODO recreate keys each time a test is run

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

	cert, err := tls.LoadX509KeyPair("./test-disperser.crt", "./test-disperser.key")
	require.NoError(t, err)

	nodeCert, err := os.ReadFile("./test-node.crt")
	require.NoError(t, err)
	certPool := x509.NewCertPool()
	ok := certPool.AppendCertsFromPEM(nodeCert)
	require.True(t, ok)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		//ServerName:   "0.0.0.0",
		//InsecureSkipVerify: true,
	})

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(creds))
	require.NoError(t, err)

	return v2.NewDispersalClient(conn)
}

func buildServer(t *testing.T) (v2.DispersalServer, *grpc.Server) {
	dispersalServer := &mockDispersalServer{}

	cert, err := tls.LoadX509KeyPair("./test-node.crt", "./test-node.key")
	require.NoError(t, err)

	disperserCert, err := os.ReadFile("./test-disperser.crt")
	require.NoError(t, err)
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(disperserCert)
	ok := certPool.AppendCertsFromPEM(disperserCert)
	require.True(t, ok)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert, // TODO commenting this makes things pass
		//ServerName:   "0.0.0.0",
	})

	server := grpc.NewServer(grpc.Creds(creds))
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
