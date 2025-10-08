package ejector

import (
	"context"

	geth "github.com/ethereum/go-ethereum/common"
)

// EjectionTransactor executes transactions related to ejections. This layer of abstraction allows for easier
// unit testing of the ejector logic.
type EjectionTransactor interface {

	// Begin ejection proceedings against the operator with the given address.
	StartEjection(ctx context.Context, addressToEject geth.Address) error

	// Checks to see if an ejection is currently in progress for the operator with the given address.
	IsEjectionInProgress(ctx context.Context, addressToCheck geth.Address) (bool, error)

	// Checks to see if the validator with the given address is present in any quorum.
	IsValidatorPresentInAnyQuorum(ctx context.Context, addressToCheck geth.Address) (bool, error)

	// Complete the ejection proceedings against the operator with the given address.
	CompleteEjection(ctx context.Context, addressToEject geth.Address) error
}

var _ EjectionTransactor = &ejectionTransactor{}

// ejectionTransactor is the production implementation of the EjectionTransactor interface.
type ejectionTransactor struct {
}

// Create a new EjectionTransactor.
func NewEjectionTransactor() EjectionTransactor {
	return &ejectionTransactor{}
}

// CompleteEjection implements EjectionTransactor.
func (e *ejectionTransactor) CompleteEjection(
	ctx context.Context,
	addressToEject geth.Address,
) error {
	panic("unimplemented")
}

// IsEjectionInProgress implements EjectionTransactor.
func (e *ejectionTransactor) IsEjectionInProgress(
	ctx context.Context,
	addressToCheck geth.Address,
) (bool, error) {
	panic("unimplemented")
}

// IsValidatorPresentInAnyQuorum implements EjectionTransactor.
func (e *ejectionTransactor) IsValidatorPresentInAnyQuorum(
	ctx context.Context,
	addressToCheck geth.Address,
) (bool, error) {
	panic("unimplemented")
}

// StartEjection implements EjectionTransactor.
func (e *ejectionTransactor) StartEjection(
	ctx context.Context,
	addressToEject geth.Address) error {
	panic("unimplemented")
}
