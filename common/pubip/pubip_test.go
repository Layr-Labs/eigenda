package pubip

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProviderOrDefault(t *testing.T) {
	p := ProviderOrDefault(SeepIPProvider)
	assert.Equal(t, SeeIP, p)
	p = ProviderOrDefault(IpifyProvider)
	assert.Equal(t, Ipify, p)
	p = ProviderOrDefault("test")
	assert.Equal(t, SeeIP, p)
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
			p := SimpleProvider{
				RequestDoer: tt.requestDoer,
				Name:        "test",
				URL:         "https://api.seeip.org",
			}

			ip, err := p.PublicIPAddress(context.Background())
			assert.Equal(t, tt.expected, ip)

			if tt.expectErr {
				assert.NotNil(t, err)
			}
		})
	}
}
