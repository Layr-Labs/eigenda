//go:build icicle

package icicle

import (
	"math"
	"sync"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	iciclebn254 "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
)

func ConvertFrToScalarFieldsBytes(data []fr.Element) []iciclebn254.ScalarField {
	scalars := make([]iciclebn254.ScalarField, len(data))

	for i := 0; i < len(data); i++ {
		src := data[i] // 4 uint64
		var littleEndian [32]byte

		fr.LittleEndian.PutElement(&littleEndian, src)
		scalars[i].FromBytesLittleEndian(littleEndian[:])
	}
	return scalars
}

func ConvertScalarFieldsToFrBytes(scalars []iciclebn254.ScalarField) []fr.Element {
	frElements := make([]fr.Element, len(scalars))

	for i := 0; i < len(frElements); i++ {
		v := scalars[i]
		slice64, _ := fr.LittleEndian.Element((*[fr.Bytes]byte)(v.ToBytesLittleEndian()))
		frElements[i] = slice64
	}
	return frElements
}

func BatchConvertGnarkAffineToIcicleAffine(gAffineList []bn254.G1Affine) []iciclebn254.Affine {
	icicleAffineList := make([]iciclebn254.Affine, len(gAffineList))
	for i := 0; i < len(gAffineList); i++ {
		GnarkAffineToIcicleAffine(&gAffineList[i], &icicleAffineList[i])
	}
	return icicleAffineList
}

func GnarkAffineToIcicleAffine(g1 *bn254.G1Affine, iciAffine *iciclebn254.Affine) {
	var littleEndBytesX, littleEndBytesY [32]byte
	fp.LittleEndian.PutElement(&littleEndBytesX, g1.X)
	fp.LittleEndian.PutElement(&littleEndBytesY, g1.Y)

	iciAffine.X.FromBytesLittleEndian(littleEndBytesX[:])
	iciAffine.Y.FromBytesLittleEndian(littleEndBytesY[:])
}

func HostSliceIcicleProjectiveToGnarkAffine(ps core.HostSlice[iciclebn254.Projective], numWorker int) []bn254.G1Affine {
	output := make([]bn254.G1Affine, len(ps))

	if len(ps) < numWorker {
		numWorker = len(ps)
	}

	var wg sync.WaitGroup

	interval := int(math.Ceil(float64(len(ps)) / float64(numWorker)))

	for w := 0; w < numWorker; w++ {
		wg.Add(1)
		start := w * interval
		end := (w + 1) * interval
		if len(ps) < end {
			end = len(ps)
		}

		go func(workerStart, workerEnd int) {
			defer wg.Done()
			for i := workerStart; i < workerEnd; i++ {
				output[i] = IcicleProjectiveToGnarkAffine(ps[i])
			}

		}(start, end)
	}
	wg.Wait()
	return output
}

func IcicleProjectiveToGnarkAffine(p iciclebn254.Projective) bn254.G1Affine {
	px, _ := fp.LittleEndian.Element((*[fp.Bytes]byte)((&p.X).ToBytesLittleEndian()))
	py, _ := fp.LittleEndian.Element((*[fp.Bytes]byte)((&p.Y).ToBytesLittleEndian()))
	pz, _ := fp.LittleEndian.Element((*[fp.Bytes]byte)((&p.Z).ToBytesLittleEndian()))

	zInv := new(fp.Element)
	x := new(fp.Element)
	y := new(fp.Element)

	zInv.Inverse(&pz)

	x.Mul(&px, zInv)
	y.Mul(&py, zInv)

	return bn254.G1Affine{X: *x, Y: *y}
}
