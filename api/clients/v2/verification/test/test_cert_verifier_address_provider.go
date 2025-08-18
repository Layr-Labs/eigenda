package test

import (
	"context"
	"sync/atomic"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
	"github.com/ethereum/go-ethereum/common"
)

// TestCertVerifierAddressProvider is an implementation of CertVerifierAddressProvider which allows the value of the
// cert verifier address to be set arbitrarily
//
// This struct is safe for concurrent use
type TestCertVerifierAddressProvider struct {
	certVerifierAddress atomic.Value
}

var _ clients.CertVerifierAddressProvider = &TestCertVerifierAddressProvider{}

func (s *TestCertVerifierAddressProvider) GetCertVerifierAddress(_ context.Context, _ uint64) (common.Address, error) {
	return s.certVerifierAddress.Load().(common.Address), nil
}

func (s *TestCertVerifierAddressProvider) SetCertVerifierAddress(inputCertVerifierAddress common.Address) {
	s.certVerifierAddress.Store(inputCertVerifierAddress)
}
