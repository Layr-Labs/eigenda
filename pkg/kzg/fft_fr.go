package kzg

import (
	"fmt"

	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

func (fs *FFTSettings) simpleFT(vals []bls.Fr, valsOffset uint64, valsStride uint64, rootsOfUnity []bls.Fr, rootsOfUnityStride uint64, out []bls.Fr) {
	l := uint64(len(out))
	var v bls.Fr
	var tmp bls.Fr
	var last bls.Fr
	for i := uint64(0); i < l; i++ {
		jv := &vals[valsOffset]
		r := &rootsOfUnity[0]
		bls.MulModFr(&v, jv, r)
		bls.CopyFr(&last, &v)

		for j := uint64(1); j < l; j++ {
			jv := &vals[valsOffset+j*valsStride]
			r := &rootsOfUnity[((i*j)%l)*rootsOfUnityStride]
			bls.MulModFr(&v, jv, r)
			bls.CopyFr(&tmp, &last)
			bls.AddModFr(&last, &tmp, &v)
		}
		bls.CopyFr(&out[i], &last)
	}
}

func (fs *FFTSettings) _fft(vals []bls.Fr, valsOffset uint64, valsStride uint64, rootsOfUnity []bls.Fr, rootsOfUnityStride uint64, out []bls.Fr) {
	if len(out) <= 4 { // if the value count is small, run the unoptimized version instead. // TODO tune threshold.
		fs.simpleFT(vals, valsOffset, valsStride, rootsOfUnity, rootsOfUnityStride, out)
		return
	}

	half := uint64(len(out)) >> 1
	// L will be the left half of out
	fs._fft(vals, valsOffset, valsStride<<1, rootsOfUnity, rootsOfUnityStride<<1, out[:half])
	// R will be the right half of out
	fs._fft(vals, valsOffset+valsStride, valsStride<<1, rootsOfUnity, rootsOfUnityStride<<1, out[half:]) // just take even again

	var yTimesRoot bls.Fr
	var x, y bls.Fr
	for i := uint64(0); i < half; i++ {
		// temporary copies, so that writing to output doesn't conflict with input
		bls.CopyFr(&x, &out[i])
		bls.CopyFr(&y, &out[i+half])
		root := &rootsOfUnity[i*rootsOfUnityStride]
		bls.MulModFr(&yTimesRoot, &y, root)
		bls.AddModFr(&out[i], &x, &yTimesRoot)
		bls.SubModFr(&out[i+half], &x, &yTimesRoot)
	}
}

func (fs *FFTSettings) FFT(vals []bls.Fr, inv bool) ([]bls.Fr, error) {
	n := uint64(len(vals))
	if n > fs.MaxWidth {
		return nil, fmt.Errorf("got %d values but only have %d roots of unity", n, fs.MaxWidth)
	}
	n = nextPowOf2(n)
	// We make a copy so we can mutate it during the work.
	valsCopy := make([]bls.Fr, n)
	for i := 0; i < len(vals); i++ {
		bls.CopyFr(&valsCopy[i], &vals[i])
	}
	for i := uint64(len(vals)); i < n; i++ {
		bls.CopyFr(&valsCopy[i], &bls.ZERO)
	}
	out := make([]bls.Fr, n)
	if err := fs.InplaceFFT(valsCopy, out, inv); err != nil {
		return nil, err
	}
	return out, nil
}

func (fs *FFTSettings) InplaceFFT(vals []bls.Fr, out []bls.Fr, inv bool) error {
	n := uint64(len(vals))
	if n > fs.MaxWidth {
		return fmt.Errorf("got %d values but only have %d roots of unity", n, fs.MaxWidth)
	}
	if !bls.IsPowerOfTwo(n) {
		return fmt.Errorf("got %d values but not a power of two", n)
	}
	if inv {
		var invLen bls.Fr
		bls.AsFr(&invLen, n)
		bls.InvModFr(&invLen, &invLen)
		rootz := fs.ReverseRootsOfUnity[:fs.MaxWidth]
		stride := fs.MaxWidth / n

		fs._fft(vals, 0, 1, rootz, stride, out)
		var tmp bls.Fr
		for i := 0; i < len(out); i++ {
			bls.MulModFr(&tmp, &out[i], &invLen)
			bls.CopyFr(&out[i], &tmp) // TODO: depending on Fr implementation, allow to directly write back to an input
		}
		return nil
	} else {
		rootz := fs.ExpandedRootsOfUnity[:fs.MaxWidth]
		stride := fs.MaxWidth / n
		// Regular FFT
		fs._fft(vals, 0, 1, rootz, stride, out)
		return nil
	}
}

// rearrange Fr elements in reverse bit order. Supports 2**31 max element count.
func reverseBitOrderFr(values []bls.Fr) error {
	if len(values) > (1 << 31) {
		return ErrFrListTooLarge
	}
	var tmp bls.Fr
	reverseBitOrder(uint32(len(values)), func(i, j uint32) {
		bls.CopyFr(&tmp, &values[i])
		bls.CopyFr(&values[i], &values[j])
		bls.CopyFr(&values[j], &tmp)
	})
	return nil
}

// rearrange Fr ptr elements in reverse bit order. Supports 2**31 max element count.
// func reverseBitOrderFrPtr(values []*bls.Fr) error {
// 	if len(values) > (1 << 31) {
// 		return ErrFrListTooLarge
// 	}
// 	reverseBitOrder(uint32(len(values)), func(i, j uint32) {
// 		values[i], values[j] = values[j], values[i]
// 	})
// 	return nil
// }
