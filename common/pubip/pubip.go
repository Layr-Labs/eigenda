package pubip

import (
	"context"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"strings"
)

const (
	SeepIPProvider = "seeip"
	SeeIPURL       = "https://api.seeip.org"

	IpifyProvider = "ipify"
	IpifyURL      = "https://api.ipify.org"

	MockIpProvider = "mockip"
)

// Provider is an interface for getting a machine's public IP address.
type Provider interface {
	// Name returns the name of the provider
	Name() string
	// PublicIPAddress returns the public IP address of the node
	PublicIPAddress(ctx context.Context) (string, error)
}

// buildSimpleProviderByName returns a simple provider with the given name.
// Returns nil if the name is not recognized.
func buildSimpleProviderByName(name string) Provider {
	if name == SeepIPProvider {
		return NewSimpleProvider(SeepIPProvider, SeeIPURL)
	} else if name == IpifyProvider {
		return NewSimpleProvider(IpifyProvider, IpifyURL)
	} else if name == MockIpProvider {
		return &mockProvider{}
	}
	return nil
}

// buildDefaultProviders returns a default provider.
func buildDefaultProvider(logger logging.Logger) Provider {
	return NewMultiProvider(logger, buildSimpleProviderByName(SeepIPProvider), buildSimpleProviderByName(IpifyProvider))
}

func providerOrDefault(logger logging.Logger, names ...string) Provider {

	for i, name := range names {
		names[i] = strings.ToLower(strings.TrimSpace(name))
	}

	if len(names) == 0 {
		return buildDefaultProvider(logger)
	} else if len(names) == 1 {
		provider := buildSimpleProviderByName(names[0])
		if provider == nil {
			logger.Warnf("Unknown IP provider '%s'", names[0])
			return buildDefaultProvider(logger)
		}
		return provider
	} else {
		providers := make([]Provider, len(names))
		for i, name := range names {
			providers[i] = buildSimpleProviderByName(name)
			if providers[i] == nil {
				logger.Warnf("Unknown IP provider '%s'", name)
				return buildDefaultProvider(logger)
			}
		}

		return NewMultiProvider(logger, providers...)
	}
}

// ProviderOrDefault returns a provider with the provided name, or a default provider if the name is not recognized.
// Provider strings are not case-sensitive.
func ProviderOrDefault(logger logging.Logger, names ...string) Provider {
	provider := providerOrDefault(logger, names...)
	logger.Infof("Using IP provider '%s'", provider.Name())
	return provider
}
