package common_test

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestGetClientAddress(t *testing.T) {

	// Make test context
	// Four proxies. The last proxy's IP address will be in the connection, not in the header
	md := metadata.Pairs("x-forwarded-for", "dummyheader, clientip", "x-forwarded-for", "proxy1, proxy2", "x-forwarded-for", "proxy3")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		t.Fatal("failed to get metadata from context")
	}
	assert.Equal(t, []string{"dummyheader, clientip", "proxy1, proxy2", "proxy3"}, md.Get("x-forwarded-for"))

	ip, err := common.GetClientAddress(ctx, "x-forwarded-for", 4, false)
	assert.NoError(t, err)

	assert.Equal(t, "clientip", ip)

}
