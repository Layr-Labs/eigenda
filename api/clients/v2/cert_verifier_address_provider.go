package clients

import "github.com/ethereum/go-ethereum/common"

// CertVerifierAddressProvider defines an object which can translate block number to cert verifier address
//
// This provider uses block number as a key, since updates to a cert verifier address in a running system are
// coordinated by defining the block number at which a new cert verifier address takes effect.
type CertVerifierAddressProvider interface {
	// GetCertVerifierAddress returns the EigenDACertVerifierAddress that is active at the input block number
	GetCertVerifierAddress(blockNumber uint64) (common.Address, error)
}
