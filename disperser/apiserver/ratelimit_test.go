package apiserver_test

import (
	"context"
	"crypto/rand"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	pb "github.com/Layr-Labs/eigenda/api/grpc/disperser"
	"github.com/Layr-Labs/eigenda/api/grpc/mock"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/peer"
)

func TestRatelimit(t *testing.T) {
	data50KiB := make([]byte, 49600) // 50*1024/32*31
	_, err := rand.Read(data50KiB)

	data50KiB = codec.ConvertByPaddingEmptyByte(data50KiB)

	assert.NoError(t, err)
	data1KiB := make([]byte, 1024) // 1024/32*31
	_, err = rand.Read(data1KiB)
	assert.NoError(t, err)

	data1KiB = codec.ConvertByPaddingEmptyByte(data1KiB)

	// Try with a non-allowlisted IP
	p := &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("1.1.1.1"),
			Port: 51001,
		},
	}
	ctx := peer.NewContext(context.Background(), p)

	// Try with non-allowlisted IP
	// Should fail with account throughput limit because unauth throughput limit is 20 KiB/s for quorum 0
	_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
		Data:                data50KiB,
		CustomQuorumNumbers: []uint32{0},
	})
	assert.ErrorContains(t, err, "Account throughput rate limit")

	// Try with non-allowlisted IP. Should fail with account blob limit because blob rate (3 blobs/s) X bucket size (3s) is smaller than 20 blobs.
	numLimited := 0
	for i := 0; i < 20; i++ {

		_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
			Data:                data1KiB,
			CustomQuorumNumbers: []uint32{1},
		})
		if err != nil && strings.Contains(err.Error(), "Account blob rate limit") {
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
		Data:                data50KiB,
		CustomQuorumNumbers: []uint32{0},
	})
	assert.NoError(t, err)

	// This should succeed because the account blob limit (5 blobs/s) X bucket size (3s) is larger than 10 blobs.
	for i := 0; i < 10; i++ {
		_, err = dispersalServer.DisperseBlob(ctx, &pb.DisperseBlobRequest{
			Data:                data1KiB,
			CustomQuorumNumbers: []uint32{1},
		})
		assert.NoError(t, err)
	}
}

func TestAuthRatelimit(t *testing.T) {

	data50KiB := make([]byte, 49600)
	_, err := rand.Read(data50KiB)
	assert.NoError(t, err)

	data50KiB = codec.ConvertByPaddingEmptyByte(data50KiB)

	data1KiB := make([]byte, 992)
	_, err = rand.Read(data1KiB)
	assert.NoError(t, err)

	data1KiB = codec.ConvertByPaddingEmptyByte(data1KiB)

	// Use an unauthenticated signer
	privateKeyHex := "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdeb"
	signer := auth.NewLocalBlobRequestSigner(privateKeyHex)

	errorChan := make(chan error, 10)

	// Should fail with account throughput limit because unauth throughput limit is 20 KiB/s for quorum 0
	simulateClient(t, signer, "2.2.2.2", data50KiB, []uint32{0}, 0, errorChan, false)

	err = <-errorChan
	assert.ErrorContains(t, err, "Account throughput rate limit")

	// Should fail with account blob limit because blob rate (3 blobs/s) X bucket size (3s) is smaller than 10 blobs.
	for i := 0; i < 20; i++ {
		simulateClient(t, signer, "3.3.3.3", data1KiB, []uint32{1}, 0, errorChan, false)
	}
	numLimited := 0
	for i := 0; i < 20; i++ {
		err = <-errorChan
		if err != nil && strings.Contains(err.Error(), "Account blob rate limit") {
			numLimited++
		}
	}
	assert.Greater(t, numLimited, 0)

	// Use an authenticated signer
	privateKeyHex = "0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcded"
	signer = auth.NewLocalBlobRequestSigner(privateKeyHex)

	// This should succeed because the account throughput limit is 100 KiB/s for quorum 0
	simulateClient(t, signer, "4.4.4.4", data50KiB, []uint32{0}, 0, errorChan, false)

	err = <-errorChan
	assert.NoError(t, err)

	// This should succeed because the account blob limit (5 blobs/s) X bucket size (3s) is larger than 10 blobs.
	for i := 0; i < 10; i++ {
		simulateClient(t, signer, "5.5.5.5", data1KiB, []uint32{1}, 0, errorChan, false)
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

func TestRetrievalRateLimit(t *testing.T) {

	// Create random data
	data := make([]byte, 992)
	_, err := rand.Read(data)
	assert.NoError(t, err)

	data = codec.ConvertByPaddingEmptyByte(data)

	// Disperse the random data
	status, blobSize, requestID := disperseBlob(t, dispersalServer, data)
	assert.Equal(t, status, pb.BlobStatus_PROCESSING)
	assert.NotNil(t, requestID)

	reply, err := dispersalServer.GetBlobStatus(context.Background(), &pb.BlobStatusRequest{
		RequestId: requestID,
	})
	assert.NoError(t, err)
	assert.Equal(t, reply.GetStatus(), pb.BlobStatus_PROCESSING)

	// Simulate blob confirmation so that we can retrieve the blob
	securityParams := []*core.SecurityParam{
		{
			QuorumID:              0,
			AdversaryThreshold:    80,
			ConfirmationThreshold: 100,
		},
		{
			QuorumID:              1,
			AdversaryThreshold:    80,
			ConfirmationThreshold: 100,
		},
	}
	_ = simulateBlobConfirmation(t, requestID, blobSize, securityParams, 1)

	reply, err = dispersalServer.GetBlobStatus(context.Background(), &pb.BlobStatusRequest{
		RequestId: requestID,
	})
	assert.NoError(t, err)
	assert.Equal(t, reply.GetStatus(), pb.BlobStatus_CONFIRMED)

	// Retrieve the blob and compare it with the original data
	numLimited := 0
	tt := time.Now()
	for i := 0; i < 15; i++ {
		_, err = retrieveBlob(dispersalServer, requestID, 1)
		fmt.Println(time.Since(tt))
		tt = time.Now()
		if err != nil && strings.Contains(err.Error(), "request ratelimited: Retrieval blob rate limit") {
			numLimited++
		}
	}
	assert.Greater(t, numLimited, 0)
}
func simulateClient(t *testing.T, signer core.BlobRequestSigner, origin string, data []byte, quorums []uint32, delay time.Duration, errorChan chan error, shouldSucceed bool) {

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

	accountId, err := signer.GetAccountID()
	assert.NoError(t, err)

	err = stream.SendFromClient(&pb.AuthenticatedRequest{
		Payload: &pb.AuthenticatedRequest_DisperseRequest{
			DisperseRequest: &pb.DisperseBlobRequest{
				Data:                data,
				CustomQuorumNumbers: quorums,
				AccountId:           accountId,
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

	time.Sleep(delay)

	// Process challenge and send back challenge_reply
	err = stream.SendFromClient(&pb.AuthenticatedRequest{Payload: &pb.AuthenticatedRequest_AuthenticationData{
		AuthenticationData: &pb.AuthenticationData{
			AuthenticationData: authData,
		},
	}})

	if shouldSucceed {

		assert.NoError(t, err)

		reply, err = stream.RecvToClient()
		assert.NoError(t, err)

		disperseReply, ok := reply.Payload.(*pb.AuthenticatedReply_DisperseReply)
		assert.True(t, ok)

		assert.Equal(t, disperseReply.DisperseReply.Result, pb.BlobStatus_PROCESSING)

	}

}
