package thegraph_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Layr-Labs/eigenda/core/thegraph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGraphQLQuerier struct {
	mock.Mock
}

func (m *MockGraphQLQuerier) Query(ctx context.Context, query any, variables map[string]any) error {
	args := m.Called(ctx, query, variables)
	return args.Error(0)
}

func TestRetryQuerier_Query(t *testing.T) {
	ctx := context.Background()
	query := "query"
	variables := map[string]any{"key": "value"}

	mockQuerier := new(MockGraphQLQuerier)
	mockQuerier.On("Query", ctx, query, variables).Return(errors.New("query error")).Once()
	mockQuerier.On("Query", ctx, query, variables).Return(errors.New("query error")).Once()
	mockQuerier.On("Query", ctx, query, variables).Return(nil)

	retryQuerier := thegraph.NewRetryQuerier(mockQuerier, time.Millisecond, 2)

	err := retryQuerier.Query(ctx, query, variables)
	assert.NoError(t, err)

	mockQuerier.AssertExpectations(t)
}

func TestRetryQuerier_ExceedMaxRetries(t *testing.T) {
	ctx := context.Background()
	query := "query"
	variables := map[string]any{"key": "value"}

	mockQuerier := new(MockGraphQLQuerier)
	mockQuerier.On("Query", ctx, query, variables).Return(errors.New("query error")).Once()
	mockQuerier.On("Query", ctx, query, variables).Return(errors.New("query error")).Once()
	mockQuerier.On("Query", ctx, query, variables).Return(errors.New("query error")).Once()

	retryQuerier := thegraph.NewRetryQuerier(mockQuerier, time.Millisecond, 2)

	err := retryQuerier.Query(ctx, query, variables)
	assert.ErrorContains(t, err, "max retries exceeded")

	mockQuerier.AssertExpectations(t)
}

func TestRetryQuerier_Timeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	query := "query"
	variables := map[string]any{"key": "value"}

	mockQuerier := new(MockGraphQLQuerier)
	mockQuerier.On("Query", ctx, query, variables).Return(errors.New("query error")).Once()
	mockQuerier.On("Query", ctx, query, variables).Return(errors.New("query error")).Once()
	mockQuerier.On("Query", ctx, query, variables).Return(nil)

	retryQuerier := thegraph.NewRetryQuerier(mockQuerier, 100*time.Millisecond, 2)

	err := retryQuerier.Query(ctx, query, variables)
	assert.ErrorContains(t, err, "context deadline exceeded")

}
