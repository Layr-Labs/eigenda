//go:build icicle

package gpu_utils

import (
	"math"
	"sync"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/ingonyama-zk/icicle/v3/wrappers/golang/core"
	bn254_icicle "github.com/ingonyama-zk/icicle/v3/wrappers/golang/curves/bn254"
)

func ConvertFrToScalarFieldsBytes(data []fr.Element) []bn254_icicle.ScalarField {
	scalars := make([]bn254_icicle.ScalarField, len(data))

	for i := 0; i < len(data); i++ {
		src := data[i] // 4 uint64
		var littleEndian [32]byte

		fr.LittleEndian.PutElement(&littleEndian, src)
		scalars[i].FromBytesLittleEndian(littleEndian[:])
	}
	return scalars
}

func ConvertScalarFieldsToFrBytes(scalars []bn254_icicle.ScalarField) []fr.Element {
	frElements := make([]fr.Element, len(scalars))

	for i := 0; i < len(frElements); i++ {
		v := scalars[i]
		slice64, _ := fr.LittleEndian.Element((*[fr.Bytes]byte)(v.ToBytesLittleEndian()))
		frElements[i] = slice64
	}
	return frElements
}

func BatchConvertGnarkAffineToIcicleAffine(gAffineList []bn254.G1Affine) []bn254_icicle.Affine {
	icicleAffineList := make([]bn254_icicle.Affine, len(gAffineList))
	for i := 0; i < len(gAffineList); i++ {
		GnarkAffineToIcicleAffine(&gAffineList[i], &icicleAffineList[i])
	}
	return icicleAffineList
}

func GnarkAffineToIcicleAffine(g1 *bn254.G1Affine, iciAffine *bn254_icicle.Affine) {
	var littleEndBytesX, littleEndBytesY [32]byte
	fp.LittleEndian.PutElement(&littleEndBytesX, g1.X)
	fp.LittleEndian.PutElement(&littleEndBytesY, g1.Y)

	iciAffine.X.FromBytesLittleEndian(littleEndBytesX[:])
	iciAffine.Y.FromBytesLittleEndian(littleEndBytesY[:])
}

func HostSliceIcicleProjectiveToGnarkAffine(ps core.HostSlice[bn254_icicle.Projective], numWorker int) []bn254.G1Affine {
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

func IcicleProjectiveToGnarkAffine(p bn254_icicle.Projective) bn254.G1Affine {
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
