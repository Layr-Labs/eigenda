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
	"testing"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFFTRoundtrip(t *testing.T) {
	fs := NewFFTSettings(4)
	data := make([]bls.Fr, fs.MaxWidth)
	for i := uint64(0); i < fs.MaxWidth; i++ {
		bls.AsFr(&data[i], i)
	}
	coeffs, err := fs.FFT(data, false)
	require.Nil(t, err)
	require.NotNil(t, coeffs)

	res, err := fs.FFT(coeffs, true)
	require.Nil(t, err)
	require.NotNil(t, coeffs)

	for i := range res {
		assert.True(t, bls.EqualFr(&res[i], &data[i]))
	}

	t.Log("zero", bls.FrStr(&bls.ZERO))
	t.Log("zero", bls.FrStr(&bls.ONE))
}

func TestInvFFT(t *testing.T) {
	fs := NewFFTSettings(4)
	data := make([]bls.Fr, fs.MaxWidth)
	for i := uint64(0); i < fs.MaxWidth; i++ {
		bls.AsFr(&data[i], i)
	}

	res, err := fs.FFT(data, true)
	require.Nil(t, err)
	require.NotNil(t, res)

	expected := []bls.Fr{
		bls.ToFr("10944121435919637611123202872628637544274182200208017171849102093287904247816"),
		bls.ToFr("1936030771851033959223912058450265953781825736913396623629635806885115007405"),
		bls.ToFr("16407567355707715082381689537916387329395994555403796510305004205827931381005"),
		bls.ToFr("10191068092603585790326358584923261075982428954421092317052884890230353083980"),
		bls.ToFr("21888242871839275220042445260109153167277707414472061641729655619866599103259"),
		bls.ToFr("21152419124866706061239949059012548909204540700669677175965090584889269743773"),
		bls.ToFr("16407567355707715086789610508212631171937308527291741914242101339246350165720"),
		bls.ToFr("12897381804114154238953344473132041472086565426937872290416035768380869236628"),
		bls.ToFr("10944121435919637611123202872628637544274182200208017171849102093287904247808"),
		bls.ToFr("8990861067725120983293061272125233616461798973478162053282168418194939258988"),
		bls.ToFr("5480675516131560135456795237044643916611055873124292429456102847329458329896"),
		bls.ToFr("735823746972569161006456686244726179343823699746357167733113601686538751843"),
		bls.ToFr("2203960485148121921270656985943972701968548566709209392357"),
		bls.ToFr("11697174779235689431920047160334014012565935445994942026645319296345455411636"),
		bls.ToFr("5480675516131560139864716207340887759152369845012237833393199980747877114611"),
		bls.ToFr("19952212099988241263022493686807009134766538663502637720068568379690693488211"),
	}

	for i := range res {
		assert.True(t, bls.EqualFr(&res[i], &expected[i]))
	}
}
