package pubip

import "context"

var _ Provider = (*mockProvider)(nil)

// mockProvider is a mock implementation of the Provider interface.
type mockProvider struct {
}

func (m mockProvider) Name() string {
	return "mockip"
}

func (m mockProvider) PublicIPAddress(ctx context.Context) (string, error) {
	return "localhost", nil
}
