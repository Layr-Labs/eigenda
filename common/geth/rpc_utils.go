package geth

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"

	"github.com/ethereum/go-ethereum/ethclient"
)

// Removes sensitive information from an RPC URL for safe logging.
// Returns scheme://hostname:port (e.g., "https://rpc.example.com:8545").
// Strips credentials, paths, and query parameters that might contain secrets.
func SanitizeRpcUrl(rawUrl string) string {
	parsed, err := url.Parse(rawUrl)
	if err != nil {
		return "[invalid-url]"
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return "[malformed-url]"
	}
	return parsed.Scheme + "://" + parsed.Host
}

// Categorizes connection errors without exposing sensitive details.
func ClassifyDialError(err error) string {
	if err == nil {
		return "unknown"
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return "timeout"
	}
	if errors.Is(err, context.Canceled) {
		return "canceled"
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		if dnsErr.IsTimeout {
			return "dns_timeout"
		}
		if dnsErr.IsNotFound {
			return "dns_not_found"
		}
		return "dns_error"
	}

	var opErr *net.OpError
	if errors.As(err, &opErr) {
		if opErr.Timeout() {
			return "timeout"
		}
		switch opErr.Op {
		case "dial":
			return "connection_refused"
		case "read":
			return "read_error"
		case "write":
			return "write_error"
		default:
			return "network_error:" + opErr.Op
		}
	}

	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return "invalid_url"
	}

	return "unknown"
}

// Wraps ethclient.DialContext and ensures errors never leak URL credentials.
// Always use this instead of calling ethclient.DialContext directly.
func SafeDial(ctx context.Context, rawUrl string) (*ethclient.Client, error) {
	client, err := ethclient.DialContext(ctx, rawUrl)
	if err != nil {
		return nil, fmt.Errorf("dial RPC endpoint %s (%s)", SanitizeRpcUrl(rawUrl), ClassifyDialError(err))
	}
	return client, nil
}
