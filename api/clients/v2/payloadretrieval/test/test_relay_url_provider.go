package test

import (
	"context"

	"github.com/Layr-Labs/eigenda/api/clients/v2/relay"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
)

// TestRelayUrlProvider implements RelayUrlProvider, for test cases
//
// NOT SAFE for concurrent use
type TestRelayUrlProvider struct {
	urlMap map[v2.RelayKey]string
}

var _ relay.RelayUrlProvider = &TestRelayUrlProvider{}

func NewTestRelayUrlProvider() *TestRelayUrlProvider {
	return &TestRelayUrlProvider{
		urlMap: make(map[v2.RelayKey]string),
	}
}

func (rup *TestRelayUrlProvider) GetRelayUrl(_ context.Context, relayKey v2.RelayKey) (string, error) {
	return rup.urlMap[relayKey], nil
}

func (rup *TestRelayUrlProvider) GetRelayCount(_ context.Context) (uint32, error) {
	return uint32(len(rup.urlMap)), nil
}

func (rup *TestRelayUrlProvider) StoreRelayUrl(relayKey v2.RelayKey, url string) {
	rup.urlMap[relayKey] = url
}
