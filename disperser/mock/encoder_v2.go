package mock

import (
	"context"

	corev2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigenda/disperser"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/stretchr/testify/mock"
)

type MockEncoderClientV2 struct {
	mock.Mock
}

var _ disperser.EncoderClientV2 = (*MockEncoderClientV2)(nil)

func NewMockEncoderClientV2() *MockEncoderClientV2 {
	return &MockEncoderClientV2{}
}

func (m *MockEncoderClientV2) EncodeBlob(ctx context.Context, blobKey corev2.BlobKey, encodingParams encoding.EncodingParams) (*encoding.FragmentInfo, error) {
	args := m.Called()
	var fragmentInfo *encoding.FragmentInfo
	if args.Get(0) != nil {
		fragmentInfo = args.Get(0).(*encoding.FragmentInfo)
	}
	return fragmentInfo, args.Error(1)
}
