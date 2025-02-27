package controller

import (
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/stretchr/testify/mock"
)

type MockBlobSet struct {
	mock.Mock
}

func (q *MockBlobSet) AddBlob(blobKey v2.BlobKey) {
	_ = q.Called(blobKey)
}

func (q *MockBlobSet) RemoveBlob(blobKey v2.BlobKey) {
	_ = q.Called(blobKey)
}

func (q *MockBlobSet) Size() int {
	args := q.Called()
	return args.Int(0)
}

func (q *MockBlobSet) Contains(blobKey v2.BlobKey) bool {
	args := q.Called(blobKey)
	return args.Bool(0)
}
