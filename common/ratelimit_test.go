package common_test

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestGetClientAddressWithTrustedProxies(t *testing.T) {

	// Make test context
	md := metadata.Pairs("x-forwarded-for", "proxy1", "x-forwarded-for", "proxy2", "x-forwarded-for", "proxy3")
	ctx := metadata.NewIncomingContext(context.Background(), md)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		t.Fatal("failed to get metadata from context")
	}
	assert.Equal(t, []string{"proxy1", "proxy2", "proxy3"}, md.Get("x-forwarded-for"))

	trustedProxies := map[string]struct{}{"proxy1": {}, "proxy2": {}}
	ip, err := common.GetClientAddressWithTrustedProxies(ctx, "x-forwarded-for", trustedProxies)
	assert.NoError(t, err)

	assert.Equal(t, "proxy3", ip)

}
