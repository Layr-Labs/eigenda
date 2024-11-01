package clients_test

import (
	"context"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/core/auth"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestPutBlobNoopSigner(t *testing.T) {
	config := clients.NewConfig("nohost", "noport", time.Second, false)
	disperserClient, err := clients.NewDisperserClient(config, auth.NewLocalNoopSigner())
	assert.NoError(t, err)

	test := []byte("test")
	test[0] = 0x00 // make sure the first byte of the requst is always 0
	quorums := []uint8{0}
	_, _, err = disperserClient.DisperseBlobAuthenticated(context.Background(), test, quorums)
	st, isGRPCError := status.FromError(err)
	assert.True(t, isGRPCError)
	assert.Equal(t, codes.InvalidArgument.String(), st.Code().String())
	assert.Equal(t, "please configure signer key if you want to use authenticated endpoint noop signer cannot get accountID", st.Message())
}
