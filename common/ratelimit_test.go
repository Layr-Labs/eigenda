package common_test

import (
	"context"
	"net"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func TestGetClientAddress(t *testing.T) {

	// Make test context
	// Four proxies. The last proxy's IP address will be in the connection, not in the header
	md := metadata.Pairs("x-forwarded-for", "dummyheader, clientip", "x-forwarded-for", "proxy1, proxy2", "x-forwarded-for", "proxy3")

	ctx := peer.NewContext(context.Background(), &peer.Peer{
		Addr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 1234,
		},
	})

	ctx = metadata.NewIncomingContext(ctx, md)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		t.Fatal("failed to get metadata from context")
	}
	assert.Equal(t, []string{"dummyheader, clientip", "proxy1, proxy2", "proxy3"}, md.Get("x-forwarded-for"))

	ip, err := common.GetClientAddress(ctx, "x-forwarded-for", 4, false)
	assert.NoError(t, err)
	assert.Equal(t, "clientip", ip)

	ip, err = common.GetClientAddress(ctx, "x-forwarded-for", 7, false)
	assert.Error(t, err)
	assert.Equal(t, "", ip)

	ip, err = common.GetClientAddress(ctx, "x-forwarded-for", 7, true)
	assert.NoError(t, err)
	assert.Equal(t, "0.0.0.0", ip)

	ip, err = common.GetClientAddress(ctx, "", 0, true)
	assert.NoError(t, err)
	assert.Equal(t, "0.0.0.0", ip)

	ip, err = common.GetClientAddress(ctx, "", 0, false)
	assert.NoError(t, err)
	assert.Equal(t, "0.0.0.0", ip)

}
