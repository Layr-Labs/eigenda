package verification

import "github.com/Layr-Labs/eigenda/api/clients/v2"

// StaticCertVerifierAddressProvider implements the CertVerifierAddressProvider, and simply returns the configured
// address every time the GetCertVerifierAddress method is called
type StaticCertVerifierAddressProvider struct {
	certVerifierAddress string
}

// NewStaticCertVerifierAddressProvider creates a CertVerifierAddressProvider which always returns the input address
// when GetCertVerifierAddress is called
func NewStaticCertVerifierAddressProvider(certVerifierAddress string) *StaticCertVerifierAddressProvider {
	return &StaticCertVerifierAddressProvider{certVerifierAddress: certVerifierAddress}
}

var _ clients.CertVerifierAddressProvider = &StaticCertVerifierAddressProvider{}

func (s *StaticCertVerifierAddressProvider) GetCertVerifierAddress(_ uint64) (string, error) {
	return s.certVerifierAddress, nil
}
