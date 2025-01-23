package pubip

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	SeepIPProvider = "seeip"
	IpifyProvider  = "ipify"
	MockIpProvider = "mockip"
)

var (
	SeeIP  = &SimpleProvider{name: "seeip", URL: "https://api.seeip.org"}
	Ipify  = &SimpleProvider{name: "ipify", URL: "https://api.ipify.org"}
	MockIp = &SimpleProvider{name: "mockip", URL: ""}
)

type RequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestDoerFunc func(req *http.Request) (*http.Response, error)

var _ RequestDoer = (RequestDoerFunc)(nil)

func (f RequestDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Provider is an interface for getting a machine's public IP address.
type Provider interface {
	// Name returns the name of the provider
	Name() string
	// PublicIPAddress returns the public IP address of the node
	PublicIPAddress(ctx context.Context) (string, error)
}

type SimpleProvider struct {
	RequestDoer RequestDoer
	name        string
	URL         string
}

func (s *SimpleProvider) Name() string {
	return s.name
}

var _ Provider = (*SimpleProvider)(nil)

func (s *SimpleProvider) PublicIPAddress(ctx context.Context) (string, error) {
	if s.name == MockIpProvider {
		return "localhost", nil
	}
	ip, err := s.doRequest(ctx, s.URL)
	if err != nil {
		return "", fmt.Errorf("%s: failed to retrieve public ip address: %w", s.Name, err)
	}
	return ip, nil
}

func (s *SimpleProvider) doRequest(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	if s.RequestDoer == nil {
		s.RequestDoer = http.DefaultClient
	}
	resp, err := s.RequestDoer.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= http.StatusBadRequest {
		return "", errors.New(resp.Status)
	}

	var b bytes.Buffer
	_, err = io.Copy(&b, resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(b.String()), nil
}

func ProviderOrDefault(name string) Provider {
	p := map[string]Provider{
		SeepIPProvider: SeeIP,
		IpifyProvider:  Ipify,
		MockIpProvider: MockIp,
	}[name]
	if p == nil {
		p = SeeIP
	}
	return p
}
