package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/common"
)

type NoopRatelimiter struct {
}

func (r *NoopRatelimiter) AllowRequest(ctx context.Context, retrieverID string, blobSize uint, rate common.RateParam) (bool, error) {
	return true, nil
}
