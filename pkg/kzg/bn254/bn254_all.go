// This code is sourced from the go-kzg Repository by protolambda.
// Original code: https://github.com/protolambda/go-kzg
// MIT License
//
// Copyright (c) 2020 @protolambda
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
