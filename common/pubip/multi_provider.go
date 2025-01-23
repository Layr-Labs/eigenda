package pubip

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"strings"
)

var _ Provider = (*multiProvider)(nil)

// An implementation of Provider that uses multiple providers. It attempts each provider in order until one succeeds.
type multiProvider struct {
	logger    logging.Logger
	providers []Provider
}

func (m *multiProvider) Name() string {
	sb := strings.Builder{}
	sb.WriteString("multiProvider(")
	for i, provider := range m.providers {
		sb.WriteString(provider.Name())
		if i < len(m.providers)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(")")
	return sb.String()
}

// NewMultiProvider creates a new multiProvider with the given providers.
func NewMultiProvider(
	logger logging.Logger,
	providers ...Provider) Provider {

	return &multiProvider{
		logger:    logger,
		providers: providers,
	}
}

func (m *multiProvider) PublicIPAddress(ctx context.Context) (string, error) {
	for _, provider := range m.providers {
		ip, err := provider.PublicIPAddress(ctx)
		if err == nil {
			return ip, nil
		}
		m.logger.Errorf("failed to get public IP address from %s: %v", provider, err)
	}

	return "", fmt.Errorf("failed to get public IP address from any provider")
}
