package kzgEncoder

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"

	rs "github.com/Layr-Labs/eigenda/pkg/encoding/encoder"
	kzg "github.com/Layr-Labs/eigenda/pkg/kzg"
	bls "github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

// Sample is the basic unit for a verification
// A blob may contain multiple Samples
type Sample struct {
	Commitment bls.G1Point
	Proof      bls.G1Point
	RowIndex   int // corresponds to a row in the verification matrix
	Coeffs     []bls.Fr
	X          uint // X is the evaluating index which corresponds to the leading coset
}

// generate a random value using Fiat Shamir transform
// we can also pseudo randomness generated locally, but we have to ensure no adversary can manipulate it
// Hashing everything takes about 1ms, so Fiat Shamir transform does not incur much cost
func GenRandomFactor(samples []Sample) (bls.Fr, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	for _, sample := range samples {
		err := enc.Encode(sample.Commitment)
		if err != nil {
			return bls.ZERO, err
		}
	}

	var randomFr bls.Fr

	err := bls.HashToSingleField(&randomFr, buffer.Bytes())
	if err != nil {
		return bls.ZERO, err
	}

	return randomFr, nil
}

// Every sample has its own randomness, even though multiple samples can come from identical blob
// Randomnesss for each sample is computed by repeatedly raising the power of the root randomness
func GenRandomnessVector(samples []Sample) ([]bls.Fr, error) {
	// root randomness
	r, err := GenRandomFactor(samples)
	if err != nil {
		return nil, err
	}

	n := len(samples)

	randomsFr := make([]bls.Fr, n)
	bls.CopyFr(&randomsFr[0], &r)

	// power of r
	for j := 0; j < n-1; j++ {
		bls.MulModFr(&randomsFr[j+1], &randomsFr[j], &r)
	}
	return randomsFr, nil
}

// the rhsG1 comprises of three terms, see https://ethresear.ch/t/a-universal-verification-equation-for-data-availability-sampling/13240/1
func genRhsG1(samples []Sample, randomsFr []bls.Fr, m int, params rs.EncodingParams, ks *kzg.KZGSettings, proofs []bls.G1Point) (*bls.G1Point, error) {
	n := len(samples)
	commits := make([]bls.G1Point, m)
	D := params.ChunkLen

	var tmp bls.Fr

	// first term
	// get coeffs to compute the aggregated commitment
	// note the coeff is affected by how many chunks are validated per blob
	// if x chunks are sampled from one blob, we need to compute the sum of all x random field element corresponding to each sample
	aggCommitCoeffs := make([]bls.Fr, m)
	setCommit := make([]bool, m)
	for k := 0; k < n; k++ {
		s := samples[k]
		row := s.RowIndex
		bls.AddModFr(&aggCommitCoeffs[row], &aggCommitCoeffs[row], &randomsFr[k])

		if !setCommit[row] {
			bls.CopyG1(&commits[row], &s.Commitment)
			setCommit[row] = true
		} else {
			if !bls.EqualG1(&commits[row], &s.Commitment) {
				return nil, errors.New("Samples of the same row has different commitments")
			}
		}
	}

	aggCommit := bls.LinCombG1(commits, aggCommitCoeffs)

	// second term
	// compute the aggregated interpolation polynomial
	aggPolyCoeffs := make([]bls.Fr, D)

	// we sum over the weighted coefficients (by the random field element) over all D monomial in all n samples
	for k := 0; k < n; k++ {
		coeffs := samples[k].Coeffs

		rk := randomsFr[k]
		// for each monomial in a given polynomial, multiply its coefficient with the corresponding random field,
		// then sum it with others. Given ChunkLen (D) is identical for all samples in a subBatch.
		// The operation is always valid.
		for j := uint64(0); j < D; j++ {
			bls.MulModFr(&tmp, &coeffs[j], &rk)
			bls.AddModFr(&aggPolyCoeffs[j], &aggPolyCoeffs[j], &tmp)
		}
	}

	// All samples in a subBatch has identical chunkLen
	aggPolyG1 := bls.LinCombG1(ks.Srs.G1[:D], aggPolyCoeffs)

	// third term
	// leading coset is an evaluation index, here we compute the weighted leading coset evaluation by random fields
	lcCoeffs := make([]bls.Fr, n)

	// get leading coset powers
	leadingDs := make([]bls.Fr, n)

	for k := 0; k < n; k++ {

		// got the leading coset field element
		h := ks.ExpandedRootsOfUnity[samples[k].X]
		var hPow bls.Fr
		bls.CopyFr(&hPow, &bls.ONE)

		// raising the power for each leading coset
		for j := uint64(0); j < D; j++ {
			bls.MulModFr(&tmp, &hPow, &h)
			bls.CopyFr(&hPow, &tmp)
		}
		bls.CopyFr(&leadingDs[k], &hPow)
	}

	// applying the random weights to leading coset elements
	for k := 0; k < n; k++ {
		rk := randomsFr[k]
		bls.MulModFr(&lcCoeffs[k], &rk, &leadingDs[k])
	}

	offsetG1 := bls.LinCombG1(proofs, lcCoeffs)

	var rhsG1 bls.G1Point
	bls.SubG1(&rhsG1, aggCommit, aggPolyG1)
	bls.AddG1(&rhsG1, &rhsG1, offsetG1)
	return &rhsG1, nil
}

// UniversalVerify implements batch verification on a set of chunks given the same chunk dimension (chunkLen, numChunk).
// The details is given in Ethereum Research post whose authors are George Kadianakis, Ansgar Dietrichs, Dankrad Feist
// https://ethresear.ch/t/a-universal-verification-equation-for-data-availability-sampling/13240
//
// m is number of blob, samples is a list of chunks
//
// The order of samples do not matter.
// Each sample need not have unique row, it is possible that multiple chunks of the same blob are validated altogether
func (group *KzgEncoderGroup) UniversalVerify(params rs.EncodingParams, samples []Sample, m int) error {
	// precheck
	for i, s := range samples {
		if s.RowIndex >= m {
			fmt.Printf("sample %v has %v Row, but there are only %v blobs\n", i, s.RowIndex, m)
			return errors.New("sample.RowIndex and numBlob are inconsistent")
		}
	}

	verifier, _ := group.GetKzgVerifier(params)
	ks := verifier.Ks

	D := params.ChunkLen

	n := len(samples)
	fmt.Printf("Batch verify %v frames of %v symbols out of %v blobs \n", n, params.ChunkLen, m)

	// generate random field elements to aggregate equality check
	randomsFr, err := GenRandomnessVector(samples)
	if err != nil {
		return err
	}

	// array of proofs
	proofs := make([]bls.G1Point, n)
	for i := 0; i < n; i++ {
		bls.CopyG1(&proofs[i], &samples[i].Proof)
	}

	// lhs g1
	lhsG1 := bls.LinCombG1(proofs, randomsFr)

	// lhs g2
	lhsG2 := &ks.Srs.G2[D]

	// rhs g2
	rhsG2 := &bls.GenG2

	// rhs g1
	rhsG1, err := genRhsG1(
		samples,
		randomsFr,
		m,
		params,
		ks,
		proofs,
	)
	if err != nil {
		return err
	}

	if bls.PairingsVerify(lhsG1, lhsG2, rhsG1, rhsG2) {
		return nil
	} else {
		return errors.New("Universal Verify Incorrect paring")
	}
}
