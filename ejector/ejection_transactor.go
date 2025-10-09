package ejector

import (
	"context"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// EjectionTransactor executes transactions related to ejections. This layer of abstraction allows for easier
// unit testing of the ejector logic.
type EjectionTransactor interface {

	// Begin ejection proceedings against the operator with the given address.
	StartEjection(ctx context.Context, addressToEject gethcommon.Address) error

	// Checks to see if an ejection is currently in progress for the operator with the given address.
	IsEjectionInProgress(ctx context.Context, addressToCheck gethcommon.Address) (bool, error)

	// Checks to see if the validator with the given address is present in any quorum.
	IsValidatorPresentInAnyQuorum(ctx context.Context, addressToCheck gethcommon.Address) (bool, error)

	// Complete the ejection proceedings against the operator with the given address.
	CompleteEjection(ctx context.Context, addressToEject gethcommon.Address) error
}
