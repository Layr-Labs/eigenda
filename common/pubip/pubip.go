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
	SeeIP  = &SimpleProvider{Name: "seeip", URL: "https://api.seeip.org"}
	Ipify  = &SimpleProvider{Name: "ipify", URL: "https://api.ipify.org"}
	MockIp = &SimpleProvider{Name: "mockip", URL: ""}
)

type RequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestDoerFunc func(req *http.Request) (*http.Response, error)

var _ RequestDoer = (RequestDoerFunc)(nil)

func (f RequestDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

type Provider interface {
	PublicIPAddress(ctx context.Context) (string, error)
}

type SimpleProvider struct {
	RequestDoer RequestDoer
	Name        string
	URL         string
}

var _ Provider = (*SimpleProvider)(nil)

func (s *SimpleProvider) PublicIPAddress(ctx context.Context) (string, error) {
	if s.Name == MockIpProvider {
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
