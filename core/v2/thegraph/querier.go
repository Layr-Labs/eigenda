package thegraph

import (
	"context"
	"errors"
	"time"
)

type RetryQuerier struct {
	GraphQLQuerier
	Backoff    time.Duration
	MaxRetries int
}

var _ GraphQLQuerier = (*RetryQuerier)(nil)

func NewRetryQuerier(q GraphQLQuerier, backoff time.Duration, maxRetries int) *RetryQuerier {
	return &RetryQuerier{
		GraphQLQuerier: q,
		Backoff:        backoff,
		MaxRetries:     maxRetries,
	}
}

func (q *RetryQuerier) Query(ctx context.Context, query any, variables map[string]any) error {

	retryCount := 0
	backoff := q.Backoff
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if retryCount > q.MaxRetries {
				return errors.New("max retries exceeded")
			}
			retryCount++

			err := q.GraphQLQuerier.Query(ctx, query, variables)
			if err == nil {
				return nil
			}

			time.Sleep(backoff)
			backoff *= 2
		}
	}
}
