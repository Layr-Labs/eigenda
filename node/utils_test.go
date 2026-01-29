package node

import (
	"context"
	"testing"

	"github.com/Layr-Labs/eigenda/common/pubip"
	"github.com/stretchr/testify/assert"
)

type mockIPProvider struct {
	ip string
}

// Verify mockIPProvider implements pubip.Provider interface
var _ pubip.Provider = (*mockIPProvider)(nil)

func (m *mockIPProvider) Name() string {
	return "mock-provider"
}

func (m *mockIPProvider) PublicIPAddress(ctx context.Context) (string, error) {
	return m.ip, nil
}

func TestSocketAddress(t *testing.T) {
	tests := []struct {
		name             string
		hostname         string
		dispersalPort    string
		retrievalPort    string
		v2DispersalPort  string
		v2RetrievalPort  string
		mockPublicIP     string
		expectedSocket   string
	}{
		{
			name:             "empty hostname uses public IP",
			hostname:         "",
			dispersalPort:    "8080",
			retrievalPort:    "8081",
			v2DispersalPort:  "8082",
			v2RetrievalPort:  "8083",
			mockPublicIP:     "192.168.1.1",
			expectedSocket:   "192.168.1.1:8080;8081;8082;8083",
		},
		{
			name:             "localhost uses public IP",
			hostname:         "localhost",
			dispersalPort:    "8080",
			retrievalPort:    "8081",
			v2DispersalPort:  "8082",
			v2RetrievalPort:  "8083",
			mockPublicIP:     "192.168.1.1",
			expectedSocket:   "192.168.1.1:8080;8081;8082;8083",
		},
		{
			name:             "127.0.0.1 uses public IP",
			hostname:         "127.0.0.1",
			dispersalPort:    "8080",
			retrievalPort:    "8081",
			v2DispersalPort:  "8082",
			v2RetrievalPort:  "8083",
			mockPublicIP:     "192.168.1.1",
			expectedSocket:   "192.168.1.1:8080;8081;8082;8083",
		},
		{
			name:             "regular hostname is used as-is",
			hostname:         "example.com",
			dispersalPort:    "8080",
			retrievalPort:    "8081",
			v2DispersalPort:  "8082",
			v2RetrievalPort:  "8083",
			mockPublicIP:     "192.168.1.1", // Should not be used
			expectedSocket:   "example.com:8080;8081;8082;8083",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := &mockIPProvider{ip: tt.mockPublicIP}
			socket, err := SocketAddress(context.Background(), mockProvider, tt.hostname, tt.dispersalPort, tt.retrievalPort, tt.v2DispersalPort, tt.v2RetrievalPort)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedSocket, socket)
		})
	}
}
