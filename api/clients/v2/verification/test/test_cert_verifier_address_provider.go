package test

import (
	"sync/atomic"

	"github.com/Layr-Labs/eigenda/api/clients/v2"
)

// TestCertVerifierAddressProvider is an implementation of CertVerifierAddressProvider which allows the value of the
// cert verifier address to be set arbitrarily
//
// This struct is safe for concurrent use
type TestCertVerifierAddressProvider struct {
	certVerifierAddress atomic.Value
}

var _ clients.CertVerifierAddressProvider = &TestCertVerifierAddressProvider{}

func (s *TestCertVerifierAddressProvider) GetCertVerifierAddress(_ uint64) (string, error) {
	return s.certVerifierAddress.Load().(string), nil
}

func (s *TestCertVerifierAddressProvider) SetCertVerifierAddress(inputCertVerifierAddress string) {
	s.certVerifierAddress.Store(inputCertVerifierAddress)
}
