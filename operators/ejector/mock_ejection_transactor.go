package ejector

import (
	"context"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

var _ EjectionTransactor = &mockEjectionTransactor{}

// mockEjectionTransactor is a mock implementation of the EjectionTransactor interface for testing purposes.
type mockEjectionTransactor struct{}

func NewMockEjectionTransactor() EjectionTransactor {
	return &mockEjectionTransactor{}
}

func (m mockEjectionTransactor) StartEjection(
	ctx context.Context,
	addressToEject gethcommon.Address,
) error {
	//TODO implement me
	panic("implement me")
}

func (m mockEjectionTransactor) IsEjectionInProgress(
	ctx context.Context,
	addressToCheck gethcommon.Address,
) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEjectionTransactor) IsValidatorPresentInAnyQuorum(
	ctx context.Context,
	addressToCheck gethcommon.Address,
) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockEjectionTransactor) CompleteEjection(
	ctx context.Context,
	addressToEject gethcommon.Address,
) error {
	//TODO implement me
	panic("implement me")
}
