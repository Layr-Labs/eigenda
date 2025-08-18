package commitments

import (
	"github.com/Layr-Labs/eigenda/api/proxy/common/types/certs"
)

// StandardCommitment is the default commitment used by arbitrum nitro stack, AVSs,
// and any stack that doesn't need any specific bytes prefix.
// Its encoding simply returns the serialized versionedCert.
type StandardCommitment struct {
	versionedCert certs.VersionedCert
}

func NewStandardCommitment(versionedCert certs.VersionedCert) StandardCommitment {
	return StandardCommitment{versionedCert}
}
func (c StandardCommitment) Encode() []byte {
	return c.versionedCert.Encode()
}
