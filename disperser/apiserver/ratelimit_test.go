package apiserver_test

import (
	"context"
	"crypto/rand"
	"net"
	"strings"
	"testing"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/api/grpc/mock"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/peer"

	tmock "github.com/stretchr/testify/mock"
)

func TestRatelimit(t *testing.T) {
	data50KiB := make([]byte, 50*1024)
	_, err := rand.Read(data50KiB)
	assert.NoError(t, err)
	data1KiB := make([]byte, 1024)
	_, err = rand.Read(data1KiB)
	assert.NoError(t, err)

	// Try with a non-allowlisted IP
	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("1.1.1.1"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	quorumParams := []*core.SecurityParam{
		{QuorumID: 0, AdversaryThreshold: 50, ConfirmationThreshold: 100},
		{QuorumID: 1, AdversaryThreshold: 50, ConfirmationThreshold: 100},
	}
	transactor.On("GetQuorumSecurityParams", tmock.Anything).Return(quorumParams, nil)
	transactor.On("GetRequiredQuorumNumbers", tmock.Anything).Return([]uint8{}, nil)

	// Try with non-allowlisted IP
	// Should fail with account throughput limit because unauth throughput limit is 20 KiB/s for quorum 0
	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:          data50KiB,
		QuorumNumbers: []uint32{0},
	})
	assert.ErrorContains(t, err, "account throughput limit")

	// Try with non-allowlisted IP. Should fail with account blob limit because blob rate (3 blobs/s) X bucket size (3s) is smaller than 20 blobs.
	numLimited := 0
	for i := 0; i < 20; i++ {
		_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
			Data:          data1KiB,
			QuorumNumbers: []uint32{1},
		})
		if err != nil && strings.Contains(err.Error(), "account blob limit") {
			numLimited++
		}
	}
	assert.Greater(t, numLimited, 0)

	// Now try with an allowlisted IP
	// This should succeed because the account throughput limit is 100 KiB/s for quorum 0
	p = &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("1.2.3.4"),
			Port: 51001,
		},
	}
	ctx = peer.NewContext(context.Background(), p)

	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:          data50KiB,
		QuorumNumbers: []uint32{0},
	})
	assert.NoError(t, err)

	// This should succeed because the account blob limit (5 blobs/s) X bucket size (3s) is larger than 10 blobs.
	for i := 0; i < 10; i++ {
		_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
			Data:          data1KiB,
			QuorumNumbers: []uint32{1},
		})
		assert.NoError(t, err)
	}
}

func TestAuthRatelimit(t *testing.T) {

	data50KiB := make([]byte, 50*1024)
	_, err := rand.Read(data50KiB)
	assert.NoError(t, err)
	data1KiB := make([]byte, 1024)
	_, err = rand.Read(data1KiB)
	assert.NoError(t, err)

	// Use an unauthenticated signer
	privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdeb"
	signer := auth.NewSigner(privateKeyHex)

	errorChan := make(chan error, 10)

	// Should fail with account throughput limit because unauth throughput limit is 20 KiB/s for quorum 0
	simulateClient(t, signer, "2.2.2.2", data50KiB, []uint32{0}, errorChan, false)

	err = <-errorChan
	assert.ErrorContains(t, err, "account throughput limit")

	// Should fail with account blob limit because blob rate (3 blobs/s) X bucket size (3s) is smaller than 10 blobs.
	for i := 0; i < 20; i++ {
		simulateClient(t, signer, "3.3.3.3", data1KiB, []uint32{1}, errorChan, false)
	}
	numLimited := 0
	for i := 0; i < 20; i++ {
		err = <-errorChan
		if err != nil && strings.Contains(err.Error(), "account blob limit") {
			numLimited++
		}
	}
	assert.Greater(t, numLimited, 0)

	// Use an authenticated signer
	privateKeyHex = "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
	signer = auth.NewSigner(privateKeyHex)

	// This should succeed because the account throughput limit is 100 KiB/s for quorum 0
	simulateClient(t, signer, "4.4.4.4", data50KiB, []uint32{0}, errorChan, false)

	err = <-errorChan
	assert.NoError(t, err)

	// This should succeed because the account blob limit (5 blobs/s) X bucket size (3s) is larger than 10 blobs.
	for i := 0; i < 10; i++ {
		simulateClient(t, signer, "5.5.5.5", data1KiB, []uint32{1}, errorChan, false)
	}
	numLimited = 0
	for i := 0; i < 10; i++ {
		err = <-errorChan
		if err != nil && strings.Contains(err.Error(), "account blob limit") {
			numLimited++
		}
	}
	assert.Equal(t, numLimited, 0)

}

func simulateClient(t *testing.T, signer core.BlobRequestSigner, origin string, data []byte, quorums []uint32, errorChan chan error, shouldSucceed bool) {

	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP(origin),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)
	stream := mock.MakeStreamMock(ctx)

	go func() {
		err := dispersalServer.DisperseBlobAuthenticated(stream)
		errorChan <- err
		stream.Close()
	}()

	err := stream.SendFromClient(&pb.AuthenticatedRequest{
		Payload: &pb.AuthenticatedRequest_DisperseRequest{
			DisperseRequest: &pb.DisperseBlobRequest{
				Data:          data,
				QuorumNumbers: quorums,
				AccountId:     signer.GetAccountID(),
			},
		},
	})
	assert.NoError(t, err)

	reply, err := stream.RecvToClient()
	assert.NoError(t, err)

	authHeaderReply, ok := reply.Payload.(*pb.AuthenticatedReply_BlobAuthHeader)
	assert.True(t, ok)

	authHeader := core.BlobAuthHeader{
		BlobCommitments: encoding.BlobCommitments{},
		AccountID:       "",
		Nonce:           authHeaderReply.BlobAuthHeader.ChallengeParameter,
	}

	authData, err := signer.SignBlobRequest(authHeader)
	assert.NoError(t, err)

	// Process challenge and send back challenge_reply
	err = stream.SendFromClient(&pb.AuthenticatedRequest{Payload: &pb.AuthenticatedRequest_AuthenticationData{
		AuthenticationData: &pb.AuthenticationData{
			AuthenticationData: authData,
		},
	}})
	assert.NoError(t, err)

	if shouldSucceed {

		reply, err = stream.RecvToClient()
		assert.NoError(t, err)

		disperseReply, ok := reply.Payload.(*pb.AuthenticatedReply_DisperseReply)
		assert.True(t, ok)

		assert.Equal(t, disperseReply.DisperseReply.Result, pb.BlobStatus_PROCESSING)

	}

}
