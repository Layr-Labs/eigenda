package rs

import (
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

// Consider input data as the polynomial Coefficients, c
// This functions computes the evaluations of the such the interpolation polynomial
// Passing through input data, evaluated at series of root of unity.
// Consider the following points (w, d[0]), (wφ, d[1]), (wφ^2, d[2]), (wφ^3, d[3])
// Suppose F be the fft matrix, then the systamtic equation that going through those points is
// d = W F c, where each row corresponds to equation being evaluated at [1, φ, φ^2, φ^3]
// where W is a diagonal matrix with diagonal [1 w w^2 w^3] for shifting the evaluation points

// The index is transformed using FFT, for example 001 => 100, 110 => 011
// The reason behind is because Reed Solomon extension using FFT insert evaluation within original
// Data. i.e. [o_1, o_2, o_3..] with coding ratio 0.5 becomes [o_1, p_1, o_2, p_2...]

func (g *Encoder) GetInterpolationPolyEval(
	interpolationPoly []bls.Fr,
	j uint32,
) ([]bls.Fr, error) {
	evals := make([]bls.Fr, g.ChunkLen)
	w := g.Fs.ExpandedRootsOfUnity[uint64(j)]
	shiftedInterpolationPoly := make([]bls.Fr, len(interpolationPoly))

	//multiply each term of the polynomial by x^i so the fourier transform results in the desired evaluations
	//The fourier matrix looks like
	// ___                    ___
	// | 1  1   1    1  . . . . |
	// | 1  φ   φ^2 φ^3         |
	// | 1  φ^2 φ^4 φ^6         |
	// | 1  φ^3 φ^6 φ^9         |  = F
	// | .   .          .       |
	// | .   .            .     |
	// | .   .              .   |
	// |__                    __|

	//
	// F * p = [p(1), p(φ), p(φ^2), ...]
	//
	// but we want
	//
	// [p(w), p(wφ), p(wφ^2), ...]
	//
	// we can do this by computing shiftedInterpolationPoly = q = p(wx) and then doing
	//
	// F * q = [p(w), p(wφ), p(wφ^2), ...]
	//
	// to get our desired evaluations
	// cool idea protolambda :)
	var wPow bls.Fr
	bls.CopyFr(&wPow, &bls.ONE)
	var tmp, tmp2 bls.Fr
	for i := 0; i < len(interpolationPoly); i++ {
		bls.MulModFr(&tmp2, &interpolationPoly[i], &wPow)
		bls.CopyFr(&shiftedInterpolationPoly[i], &tmp2)
		bls.MulModFr(&tmp, &wPow, &w)
		bls.CopyFr(&wPow, &tmp)
	}

	err := g.Fs.InplaceFFT(shiftedInterpolationPoly, evals, false)
	return evals, err
}

// Since both F W are invertible, c = W^-1 F^-1 d, convert it back. F W W^-1 F^-1 d = c
func (g *Encoder) GetInterpolationPolyCoeff(chunk []bls.Fr, k uint32) ([]bls.Fr, error) {
	coeffs := make([]bls.Fr, g.ChunkLen)
	w := g.Fs.ExpandedRootsOfUnity[uint64(k)]
	shiftedInterpolationPoly := make([]bls.Fr, len(chunk))
	err := g.Fs.InplaceFFT(chunk, shiftedInterpolationPoly, true)
	if err != nil {
		return coeffs, err
	}
	var wPow bls.Fr
	bls.CopyFr(&wPow, &bls.ONE)
	var tmp, tmp2 bls.Fr
	for i := 0; i < len(chunk); i++ {
		bls.InvModFr(&tmp, &wPow)
		bls.MulModFr(&tmp2, &shiftedInterpolationPoly[i], &tmp)
		bls.CopyFr(&coeffs[i], &tmp2)
		bls.MulModFr(&tmp, &wPow, &w)
		bls.CopyFr(&wPow, &tmp)
	}
	return coeffs, nil
}
