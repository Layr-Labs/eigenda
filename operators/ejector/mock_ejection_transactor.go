package ejector

import (
	"context"
	"fmt"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

var _ EjectionTransactor = &mockEjectionTransactor{}

// mockEjectionTransactor is a mock implementation of the EjectionTransactor interface for testing purposes.
type mockEjectionTransactor struct {

	// A set of addresses for which ejection is currently in progress.
	inProgressEjections map[gethcommon.Address]struct{}

	// A set of addresses for which ejection has been completed.
	completedEjections map[gethcommon.Address]struct{}

	// The values to return for IsValidatorPresentInAnyQuorum calls.
	isValidatorPresentInAnyQuorumResponses map[gethcommon.Address]bool

	// A map of addresses to errors to return for IsEjectionInProgress calls.
	isEjectionInProgressErrors map[gethcommon.Address]error

	// A map of addresses to errors to return for IsValidatorPresentInAnyQuorum calls.
	isValidatorPresentInAnyQuorumErrors map[gethcommon.Address]error

	// A map of addresses to errors to return for CompleteEjection calls.
	completeEjectionErrors map[gethcommon.Address]error
}

func newMockEjectionTransactor() *mockEjectionTransactor {
	return &mockEjectionTransactor{
		inProgressEjections:                    make(map[gethcommon.Address]struct{}),
		completedEjections:                     make(map[gethcommon.Address]struct{}),
		isValidatorPresentInAnyQuorumResponses: make(map[gethcommon.Address]bool),
		isEjectionInProgressErrors:             make(map[gethcommon.Address]error),
	}
}

func (m mockEjectionTransactor) StartEjection(
	_ context.Context,
	addressToEject gethcommon.Address,
) error {

	if _, ok := m.inProgressEjections[addressToEject]; ok {
		return fmt.Errorf("ejection already in progress")
	}

	m.inProgressEjections[addressToEject] = struct{}{}
	return nil
}

func (m mockEjectionTransactor) IsEjectionInProgress(
	_ context.Context,
	addressToCheck gethcommon.Address,
) (bool, error) {

	_, inProgress := m.inProgressEjections[addressToCheck]
	return inProgress, m.isEjectionInProgressErrors[addressToCheck]
}

func (m mockEjectionTransactor) IsValidatorPresentInAnyQuorum(
	_ context.Context,
	addressToCheck gethcommon.Address,
) (bool, error) {
	return m.isValidatorPresentInAnyQuorumResponses[addressToCheck],
		m.isValidatorPresentInAnyQuorumErrors[addressToCheck]
}

func (m mockEjectionTransactor) CompleteEjection(
	_ context.Context,
	addressToEject gethcommon.Address,
) error {

	if err, ok := m.completeEjectionErrors[addressToEject]; ok {
		return err
	}

	if _, ok := m.inProgressEjections[addressToEject]; !ok {
		return fmt.Errorf("no ejection in progress for address %s", addressToEject.Hex())
	}

	if _, ok := m.completedEjections[addressToEject]; ok {
		return fmt.Errorf("ejection already completed for address %s", addressToEject.Hex())
	}

	delete(m.inProgressEjections, addressToEject)
	m.completedEjections[addressToEject] = struct{}{}

	// Once ejected, the validator should no longer be present in any quorum.
	m.isValidatorPresentInAnyQuorumResponses[addressToEject] = true

	return nil
}
