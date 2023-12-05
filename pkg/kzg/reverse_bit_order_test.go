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
package kzg

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReverseBitOrder(t *testing.T) {
	for s := 2; s < 2048; s *= 2 {
		t.Run(fmt.Sprintf("size_%d", s), func(t *testing.T) {
			data := make([]uint32, s)
			for i := 0; i < s; i++ {
				data[i] = uint32(i)
			}

			reverseBitOrder(uint32(s), func(i, j uint32) {
				data[i], data[j] = data[j], data[i]
			})

			for i := 0; i < s; i++ {
				assert.Equal(t, reverseBitsLimited(uint32(s), uint32(i)), data[i], "bad reversal at %d", i)
				expected := fmt.Sprintf("%0"+fmt.Sprintf("%d", s)+"b", i)
				got := fmt.Sprintf("%0"+fmt.Sprintf("%d", s)+"b", data[i])
				assert.Equal(t, len(expected), len(got), "bad length: %d, expected %d", len(got), len(expected))

				for j := 0; j < len(expected); j++ {
					// TODO: add check
				}
			}
		})
	}
}

func TestRevBitorderBitIndex(t *testing.T) {
	for i := 0; i < 32; i++ {
		got := bitIndex(uint32(1 << i))
		assert.Equal(t, got, uint8(i), "bit index %d is wrong: %d", i, got)
	}
}

func TestReverseBits(t *testing.T) {
	rng := rand.New(rand.NewSource(1234))
	for i := 0; i < 10000; i++ {
		v := rng.Uint32()
		expected := revStr(fmt.Sprintf("%032b", v))
		out := reverseBits(v)
		got := fmt.Sprintf("%032b", out)
		assert.Equal(t, expected, got, "bit mismatch: expected: %s, got: %s ", expected, got)
	}
}

func revStr(v string) string {
	out := make([]byte, len(v))
	for i := 0; i < len(v); i++ {
		out[i] = v[len(v)-1-i]
	}
	return string(out)
}
