package common

import (
	"context"
	"errors"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// Requester ID is the ID of the party making the request. In the case of a rollup making a dispersal request, the Requester
// ID is the authenticated Account ID. For retrieval requests, the requester ID will be the requester's IP address.
type RequesterID = string

// RequesterName is the friendly name of the party making the request. In the case
// of a rollup making a dispersal request, the RequesterName is the name of the rollup.
type RequesterName = string

type RequestParams struct {
	RequesterID   RequesterID
	RequesterName RequesterName
	BlobSize      uint
	Rate          RateParam
	Info          interface{}
}

type RateLimiter interface {
	// AllowRequest checks whether the request should be allowed. If the request is allowed, the function returns true.
	// If the request is not allowed, the function returns false and the RequestParams of the request that was not allowed.
	// In order for the request to be allowed, all of the requests represented by the RequestParams slice must be allowed.
	// Each RequestParams object represents a single request. Each request is subjected to the same GlobalRateParams, but the
	// individual parameters of the request can differ.
	//
	// If CountFailed is set to true in the GlobalRateParams, AllowRequest will count failed requests towards the rate limit.
	// If CountFailed is set to false, the rate limiter will stop processing requests as soon as it encounters a request that
	// is not allowed.
	AllowRequest(ctx context.Context, params []RequestParams) (bool, *RequestParams, error)
}

type GlobalRateParams struct {
	// BucketSizes are the time scales at which the rate limit is enforced.
	// For each time scale, the rate limiter will make sure that the given rate (possibly subject to a relaxation given
	// by one of the Multipliers) is observed when the request bandwidth is averaged at this time scale.
	// In terms of implementation, the rate limiter uses a set of "time buckets". A time bucket, i, is filled to a maximum of
	// `BucketSizes[i]` at a rate of 1, and emptied by an amount equal to `(size of request)/RateParam` each time a
	// request is processed.
	BucketSizes []time.Duration
	// Multipliers specify how much the supplied rate limit should be relaxed for each time scale.
	// For i'th BuckeSize, the RateParam*Multiplier[i] will be applied.
	Multipliers []float32
	// CountFailed indicates whether failed requests should be counted towards the rate limit.
	CountFailed bool
}

// RateParam is the type used for expressing a bandwidth based rate limit in units of Bytes/second
type RateParam = uint32

type RateBucketParams struct {
	// BucketLevels stores the amount of time contained in each bucket. For instance, if the bucket contains 1 minute, this means
	// that the requester can consume one minute worth of bandwidth (in terms of amount of data, this equals RateParam * one minute)
	// before the rate limiter will throttle them
	BucketLevels []time.Duration
	// LastRequestTime stores the time of the last request received from a given requester. All times are stored in UTC.
	LastRequestTime time.Time
}

// GetClientAddress returns the client address from the context. If the header is not empty, it will
// take the ip address located at the `numProxies` position from the end of the header. If the ip address cannot be
// found in the header, it will use the connection ip if `allowDirectConnectionFallback` is true. Otherwise, it will return
// an error.
func GetClientAddress(ctx context.Context, header string, numProxies int, allowDirectConnectionFallback bool) (string, error) {

	if header != "" && numProxies > 0 {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok && len(md.Get(header)) > 0 {
			parts := splitHeader(md.Get(header))
			if len(parts) >= numProxies {
				return parts[len(parts)-numProxies], nil
			}
		}
	}

	if header == "" || allowDirectConnectionFallback {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return "", errors.New("failed to get peer from request")
		}
		addr := p.Addr.String()
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			return "", err
		}
		return host, nil
	}

	return "", errors.New("failed to get ip")
}

func splitHeader(header []string) []string {
	var result []string
	for _, h := range header {
		for _, p := range strings.Split(h, ",") {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
	}
	return result
}
