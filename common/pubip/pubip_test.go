package pubip

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProviderOrDefault(t *testing.T) {
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	provider := ProviderOrDefault(logger, SeepIPProvider)
	require.Equal(t, SeepIPProvider, provider.Name())
	seeIPProvider, ok := provider.(*simpleProvider)
	require.True(t, ok)
	require.Equal(t, SeeIPURL, seeIPProvider.URL)

	provider = ProviderOrDefault(logger, IpifyProvider)
	require.Equal(t, IpifyProvider, provider.Name())
	ipifyProvider, ok := provider.(*simpleProvider)
	require.True(t, ok)
	require.Equal(t, IpifyURL, ipifyProvider.URL)

	provider = ProviderOrDefault(logger, MockIpProvider)
	require.Equal(t, MockIpProvider, provider.Name())
	_, ok = provider.(*mockProvider)
	require.True(t, ok)

	// invalid provider, should yield default
	provider = ProviderOrDefault(logger, "this is not a supported provider")
	require.Equal(t, fmt.Sprintf("multiProvider(%s, %s)", SeepIPProvider, IpifyProvider), provider.Name())
	multi, ok := provider.(*multiProvider)
	require.True(t, ok)
	require.Equal(t, 2, len(multi.providers))
	require.Equal(t, SeepIPProvider, multi.providers[0].Name())
	require.Equal(t, IpifyProvider, multi.providers[1].Name())

	provider = providerOrDefault(logger, fmt.Sprintf("%s,%s", SeepIPProvider, IpifyProvider))
	require.Equal(t, fmt.Sprintf("multiProvider(%s, %s)", SeepIPProvider, IpifyProvider), provider.Name())
	multi, ok = provider.(*multiProvider)
	require.True(t, ok)
	require.Equal(t, 2, len(multi.providers))
	require.Equal(t, SeepIPProvider, multi.providers[0].Name())
	require.Equal(t, IpifyProvider, multi.providers[1].Name())

	provider = providerOrDefault(logger, fmt.Sprintf("%s,%s,%s", IpifyProvider, SeepIPProvider, MockIpProvider))
	require.Equal(t, fmt.Sprintf("multiProvider(%s, %s, %s)",
		IpifyProvider, SeepIPProvider, MockIpProvider), provider.Name())
	multi, ok = provider.(*multiProvider)
	require.True(t, ok)
	require.Equal(t, 3, len(multi.providers))
	require.Equal(t, IpifyProvider, multi.providers[0].Name())
	require.Equal(t, SeepIPProvider, multi.providers[1].Name())
	require.Equal(t, MockIpProvider, multi.providers[2].Name())

	// invalid provider, should yield default
	provider = providerOrDefault(logger, fmt.Sprintf("%s,not a real provider,%s", IpifyProvider, MockIpProvider))
	require.Equal(t, fmt.Sprintf("multiProvider(%s, %s)", SeepIPProvider, IpifyProvider), provider.Name())
	multi, ok = provider.(*multiProvider)
	require.True(t, ok)
	require.Equal(t, 2, len(multi.providers))
	require.Equal(t, SeepIPProvider, multi.providers[0].Name())
	require.Equal(t, IpifyProvider, multi.providers[1].Name())
}

var _ Provider = (*testProvider)(nil)

type testProvider struct {
	// if true then this PublicIPAddress will return an error
	returnErr bool

	// number of times PublicIPAddress was called
	count int

	// ip address to return when PublicIPAddress is called
	ip string
}

func (t *testProvider) Name() string {
	return "test"
}

func (t *testProvider) PublicIPAddress(ctx context.Context) (string, error) {
	t.count++
	if t.returnErr {
		return "", fmt.Errorf("intentional error")
	}
	return t.ip, nil
}

func TestMultiProvider(t *testing.T) {
	rand := random.NewTestRandom(t)
	logger, err := common.NewLogger(common.DefaultLoggerConfig())
	require.NoError(t, err)

	provider1 := &testProvider{
		ip: rand.String(10),
	}
	provider2 := &testProvider{
		ip: rand.String(10),
	}
	provider3 := &testProvider{
		ip: rand.String(10),
	}
	provider := NewMultiProvider(logger, provider1, provider2, provider3)

	ip, err := provider.PublicIPAddress(context.Background())
	require.NoError(t, err)
	require.Equal(t, 1, provider1.count)
	require.Equal(t, 0, provider2.count)
	require.Equal(t, 0, provider3.count)
	require.Equal(t, provider1.ip, ip)

	provider1.returnErr = true
	ip, err = provider.PublicIPAddress(context.Background())
	require.NoError(t, err)
	require.Equal(t, 2, provider1.count)
	require.Equal(t, 1, provider2.count)
	require.Equal(t, 0, provider3.count)
	require.Equal(t, provider2.ip, ip)

	provider2.returnErr = true
	ip, err = provider.PublicIPAddress(context.Background())
	require.NoError(t, err)
	require.Equal(t, 3, provider1.count)
	require.Equal(t, 2, provider2.count)
	require.Equal(t, 1, provider3.count)
	require.Equal(t, provider3.ip, ip)

	provider3.returnErr = true
	ip, err = provider.PublicIPAddress(context.Background())
	require.Error(t, err)
	require.Equal(t, 4, provider1.count)
	require.Equal(t, 3, provider2.count)
	require.Equal(t, 2, provider3.count)
	require.Equal(t, "", ip)
}

func TestSimpleProvider_PublicIPAddress(t *testing.T) {
	tests := []struct {
		name        string
		requestDoer RequestDoerFunc
		expectErr   bool
		expected    string
	}{
		{
			name: "success",
			requestDoer: func(req *http.Request) (*http.Response, error) {
				w := httptest.NewRecorder()
				_, _ = w.WriteString("\n\n8.8.8.8\n\n")
				return w.Result(), nil
			},
			expectErr: false,
			expected:  "8.8.8.8",
		},
		{
			name: "http error status",
			requestDoer: func(req *http.Request) (*http.Response, error) {
				w := httptest.NewRecorder()
				w.WriteHeader(http.StatusInternalServerError)
				return w.Result(), nil
			},
			expectErr: true,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := CustomProvider(
				tt.requestDoer,
				"test",
				"https://api.seeip.org")

			ip, err := p.PublicIPAddress(context.Background())
			assert.Equal(t, tt.expected, ip)

			if tt.expectErr {
				assert.NotNil(t, err)
			}
		})
	}
}
