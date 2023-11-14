package common

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// Requester ID is the ID of the party making the request. In the case of a rollup making a dispersal request, the Requester
// ID is the authenticated Account ID. For retrieval requests, the requester ID will be the requester's IP address.
type RequesterID = string

type RateLimiter interface {
	AllowRequest(ctx context.Context, requesterID RequesterID, blobSize uint, rate RateParam) (bool, error)
}

type GlobalRateParams struct {
	// BucketSizes are the time scales at which the rate limit is enforced.
	// For each time scale, the rate limiter will make sure that the give rate (possibly subject to a relaxation given
	// by one of the Multipliers) is observed when the request bandwidth is averaged at this time scale.
	// In terms of implementation, the rate limiter uses a set of "time buckets". A time bucket, i, is filled to a maximum of
	// `BucketSizes[i]` at a rate of 1, and emptied by an amount equal to `(size of request)/RateParam` each time a
	// request is processed.
	BucketSizes []time.Duration
	// Multipliers speicify how much the supplied rate limit should be relaxed for each time scale.
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

func GetClientAddress(ctx context.Context, header string) (string, error) {
	if header != "" {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok || len(md.Get(header)) == 0 {
			return "", fmt.Errorf("failed to get ip from header")
		}
		return md.Get(header)[len(md.Get(header))-1], nil
	} else {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return "", fmt.Errorf("failed to get peer from request")
		}
		addr := p.Addr.String()
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			return "", err
		}
		return host, nil
	}
}

// GetClientAddressWithTrustedProxies is a modified version of GetClientAddress that takes into account the possibility
// that the request has been proxied through multiple trusted proxies. The function takes in a list of trusted proxies
// and will pass through the request headers until it finds the first non-trusted proxy. If the request has not been
// proxied, the function will return the same result as GetClientAddress.
func GetClientAddressWithTrustedProxies(ctx context.Context, headerName string, trustedProxies map[string]struct{}) (string, error) {

	// If no headerName is specified, return the remote address
	if headerName == "" {
		p, ok := peer.FromContext(ctx)
		if !ok {
			return "", fmt.Errorf("failed to get peer from request")
		}
		addr := p.Addr.String()
		host, _, err := net.SplitHostPort(addr)
		if err != nil {
			return "", err
		}
		return host, nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md.Get(headerName)) == 0 {
		return "", fmt.Errorf("the header %s is not present in the request", headerName)
	}

	ips := md.Get(headerName)

	// Iterate over the IPs from right to left
	for i := len(ips) - 1; i >= 0; i-- {
		ip := ips[i]

		// If the IP is not in the trusted proxies list, return it
		if _, ok := trustedProxies[ip]; !ok {
			return ip, nil
		}
	}

	return "", fmt.Errorf("all IPs in header are trusted proxies")
}
