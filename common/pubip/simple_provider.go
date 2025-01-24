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

var _ Provider = (*simpleProvider)(nil)

type RequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type RequestDoerFunc func(req *http.Request) (*http.Response, error)

var _ RequestDoer = (RequestDoerFunc)(nil)

func (f RequestDoerFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

// simpleProvider is a simple implementation of the Provider interface that checks with a single endpoint.
type simpleProvider struct {
	RequestDoer RequestDoer
	name        string
	URL         string
}

// CustomProvider creates a new simpleProvider with the given request doer, name, and URL.
func CustomProvider(requestDoer RequestDoer, name, url string) Provider {
	return &simpleProvider{
		RequestDoer: requestDoer,
		name:        name,
		URL:         url,
	}
}

// NewSimpleProvider creates a new simpleProvider with the given name and URL.
func NewSimpleProvider(name, url string) Provider {
	return &simpleProvider{
		name: name,
		URL:  url,
	}
}

func (s *simpleProvider) Name() string {
	return s.name
}

func (s *simpleProvider) PublicIPAddress(ctx context.Context) (string, error) {
	ip, err := s.doRequest(ctx, s.URL)
	if err != nil {
		return "", fmt.Errorf("%s: failed to retrieve public ip address: %w", s.name, err)
	}
	return ip, nil
}

func (s *simpleProvider) doRequest(ctx context.Context, url string) (string, error) {
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
