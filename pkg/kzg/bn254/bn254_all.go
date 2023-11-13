package bn254

import (
	"encoding/hex"
	"errors"
)

func (p *G1Point) MarshalText() []byte {
	return []byte(hex.EncodeToString(ToCompressedG1(p)))
}

// UnmarshalText decodes hex formatted text (no 0x prefix) into a G1Point
func (p *G1Point) UnmarshalText(text []byte) error {
	if p == nil {
		return errors.New("cannot decode into nil G1Point")
	}
	data, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	d, err := FromCompressedG1(data)
	if err != nil {
		return err
	}
	*p = *d
	return nil
}

// MarshalText encodes G2Point into hex formatted text (no 0x prefix)
func (p *G2Point) MarshalText() []byte {
	return []byte(hex.EncodeToString(ToCompressedG2(p)))
}

// UnmarshalText decodes hex formatted text (no 0x prefix) into a G2Point
func (p *G2Point) UnmarshalText(text []byte) error {
	if p == nil {
		return errors.New("cannot decode into nil G2Point")
	}
	data, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	d, err := FromCompressedG2(data)
	if err != nil {
		return err
	}
	*p = *d
	return nil
}
