package mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockFinalizer struct {
	mock.Mock
}

func NewFinalizer() *MockFinalizer {
	return &MockFinalizer{}
}

func (b *MockFinalizer) Start(ctx context.Context) {}

func (b *MockFinalizer) FinalizeBlobs(ctx context.Context) error {
	args := b.Called()
	return args.Error(0)
}
