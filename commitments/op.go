package commitments

import (
	"fmt"
	"log"
)

type OPCommitmentType byte

const (
	// Keccak256CommitmentTypeByte represents a commitment using Keccak256 hashing.
	Keccak256CommitmentTypeByte OPCommitmentType = 0
	// GenericCommitmentTypeByte represents a commitment using a DA service.
	GenericCommitmentTypeByte OPCommitmentType = 1
)

type OPCommitment struct {
	keccak256Commitment []byte
	genericCommitment   *DAServiceOPCommitment
}

var _ Commitment = (*OPCommitment)(nil)

func Keccak256Commitment(value []byte) OPCommitment {
	return OPCommitment{keccak256Commitment: value}
}

func GenericCommitment(value DAServiceOPCommitment) OPCommitment {
	return OPCommitment{genericCommitment: &value}
}

func (e OPCommitment) IsKeccak256Commitment() bool {
	return e.keccak256Commitment != nil
}

func (e OPCommitment) IsGenericCommitment() bool {
	return e.genericCommitment != nil
}

func (e OPCommitment) MustKeccak256CommitmentValue() []byte {
	if e.keccak256Commitment != nil {
		return e.keccak256Commitment
	}
	log.Panic("OPCommitment does not contain a Keccak256Commitment value")
	return nil // This will never be reached, but is required for compilation.
}

func (e OPCommitment) MustGenericCommitmentValue() DAServiceOPCommitment {
	if e.genericCommitment != nil {
		return *e.genericCommitment
	}
	log.Panic("OPCommitment does not contain a DAServiceCommitment value")
	return DAServiceOPCommitment{} // This will never be reached, but is required for compilation.
}

func (e OPCommitment) Marshal() ([]byte, error) {
	if e.IsGenericCommitment() {
		bytes, err := e.MustGenericCommitmentValue().Marshal()
		if err != nil {
			return nil, err
		}
		return append([]byte{byte(GenericCommitmentTypeByte)}, bytes...), nil
	} else if e.IsKeccak256Commitment() {
		return append([]byte{byte(Keccak256CommitmentTypeByte)}, e.MustKeccak256CommitmentValue()...), nil
	} else {
		return nil, fmt.Errorf("OPCommitment is neither a Keccak256 commitment nor a DA service commitment")
	}
}

func (e *OPCommitment) Unmarshal(bz []byte) error {
	if len(bz) < 1 {
		return fmt.Errorf("OPCommitment does not contain a commitment type prefix byte")
	}
	head := OPCommitmentType(bz[0])
	tail := bz[1:]
	switch head {
	case Keccak256CommitmentTypeByte:
		e.keccak256Commitment = tail
	case GenericCommitmentTypeByte:
		daServiceCommitment := DAServiceOPCommitment{}
		err := daServiceCommitment.Unmarshal(tail)
		if err != nil {
			return err
		}
		e.genericCommitment = &daServiceCommitment
	default:
		return fmt.Errorf("unrecognized commitment type byte: %x", bz[0])
	}
	return nil
}
