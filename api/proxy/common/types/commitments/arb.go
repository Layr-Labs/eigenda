package commitments

import "github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"

const (
	ArbCustomDAHeaderByte = 0x01
)

// ArbitrumCommitment is the default commitment used by arbitrum nitro stack, AVSs,
// and any stack that doesn't need any specific bytes prefix.
// Its encoding simply returns the serialized versionedCert.
type ArbitrumCommitment struct {
	versionedCert certs.VersionedCert
}

func NewArbCommitment(versionedCert certs.VersionedCert) ArbitrumCommitment {
	return ArbitrumCommitment{versionedCert}
}
func (c ArbitrumCommitment) Encode() []byte {
	return append([]byte{ArbCustomDAHeaderByte}, c.versionedCert.Encode()...)
}
