package clients

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
)

// CertVerifierAddressProvider defines an object which can translate block number to cert verifier address
//
// This provider uses reference block number as a key, since updates to a cert verifier address in a running system are
// coordinated by defining the reference block number at which a new cert verifier address takes effect. Specifically,
// a blob shall be verified by the latest defined cert verifier contract with a reference block number key that doesn't
// exceed the reference block number of the blob's batch.
type CertVerifierAddressProvider interface {
	// GetCertVerifierAddress returns the EigenDACertVerifierAddress that is active at the input reference block number
	GetCertVerifierAddress(ctx context.Context, referenceBlockNumber uint64) (common.Address, error)
}
