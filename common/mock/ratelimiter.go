package mock

import (
	"context"

	"github.com/Layr-Labs/eigenda/common"
)

type NoopRatelimiter struct {
}

var _ common.RateLimiter = &NoopRatelimiter{}

func (r *NoopRatelimiter) AllowRequest(ctx context.Context, params []common.RequestParams) (bool, *common.RequestParams, error) {
	return true, nil, nil
}
