package thegraph

import (
	"context"
	"errors"
	"time"
)

type RetryQuerier struct {
	GraphQLQuerier
	PullInterval time.Duration
	MaxRetries   int
}

var _ GraphQLQuerier = (*RetryQuerier)(nil)

func NewRetryQuerier(q GraphQLQuerier, interval time.Duration, maxRetries int) *RetryQuerier {
	return &RetryQuerier{
		GraphQLQuerier: q,
		PullInterval:   interval,
		MaxRetries:     maxRetries,
	}
}

func (q *RetryQuerier) Query(ctx context.Context, query any, variables map[string]any) error {
	for {
		retryCount := 0
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
			// Optionally, add a delay before retrying
			time.Sleep(q.PullInterval)
		}
	}
}
