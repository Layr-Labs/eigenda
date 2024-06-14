package commitments

import (
	"fmt"
	"log"
)

// Define the parent and child types
type CertEncodingVersion byte

const (
	CertEncodingV0 CertEncodingVersion = 0
)

type EigenDACommitment struct {
	certV0 []byte
}

var _ Commitment = (*EigenDACommitment)(nil)

func EigenDACertV0(value []byte) EigenDACommitment {
	return EigenDACommitment{certV0: value}
}

func (e EigenDACommitment) IsCertV0() bool {
	return e.certV0 != nil
}

func (e EigenDACommitment) MustCertV0Value() []byte {
	if e.certV0 != nil {
		return e.certV0
	}
	log.Panic("CommitmentEither does not contain a Keccak256Commitment value")
	return nil // This will never be reached, but is required for compilation.
}

func (e EigenDACommitment) Marshal() ([]byte, error) {
	if e.IsCertV0() {
		return append([]byte{byte(CertEncodingV0)}, e.certV0...), nil
	} else {
		return nil, fmt.Errorf("EigenDADAServiceOPCommitment is of unknown type")
	}
}

func (e *EigenDACommitment) Unmarshal(bz []byte) error {
	if len(bz) < 1 {
		return fmt.Errorf("OP commitment does not contain eigenda commitment encoding version prefix byte")
	}
	head := CertEncodingVersion(bz[0])
	tail := bz[1:]
	switch head {
	case CertEncodingV0:
		e.certV0 = tail
	default:
		return fmt.Errorf("unrecognized EigenDA commitment encoding type byte: %x", bz[0])
	}
	return nil
}
