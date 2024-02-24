package verifier

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"math"

	"github.com/Layr-Labs/eigenda/encoding"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

// Sample is the basic unit for a verification
// A blob may contain multiple Samples
type Sample struct {
	Commitment bn254.G1Affine
	Proof      bn254.G1Affine
	RowIndex   int // corresponds to a row in the verification matrix
	Coeffs     []fr.Element
	X          uint // X is the evaluating index which corresponds to the leading coset
}

// generate a random value using Fiat Shamir transform
// we can also pseudo randomness generated locally, but we have to ensure no adversary can manipulate it
// Hashing everything takes about 1ms, so Fiat Shamir transform does not incur much cost
func GenRandomFactor(samples []Sample) (fr.Element, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	for _, sample := range samples {
		err := enc.Encode(sample.Commitment)
		if err != nil {
			return fr.Element{}, err
		}
	}

	var randomFr fr.Element

	err := kzg.HashToSingleField(&randomFr, buffer.Bytes())
	if err != nil {
		return fr.Element{}, err
	}

	return randomFr, nil
}

// Every sample has its own randomness, even though multiple samples can come from identical blob
// Randomnesss for each sample is computed by repeatedly raising the power of the root randomness
func GenRandomnessVector(samples []Sample) ([]fr.Element, error) {
	// root randomness
	r, err := GenRandomFactor(samples)
	if err != nil {
		return nil, err
	}

	n := len(samples)

	randomsFr := make([]fr.Element, n)
	randomsFr[0].Set(&r)

	// power of r
	for j := 0; j < n-1; j++ {
		randomsFr[j+1].Mul(&randomsFr[j], &r)
	}
	return randomsFr, nil
}

// the rhsG1 comprises of three terms, see https://ethresear.ch/t/a-universal-verification-equation-for-data-availability-sampling/13240/1
func genRhsG1(samples []Sample, randomsFr []fr.Element, m int, params encoding.EncodingParams, ks *kzg.KZGSettings, proofs []bn254.G1Affine) (*bn254.G1Affine, error) {
	n := len(samples)
	commits := make([]bn254.G1Affine, m)
	D := params.ChunkLength

	var tmp fr.Element

	// first term
	// get coeffs to compute the aggregated commitment
	// note the coeff is affected by how many chunks are validated per blob
	// if x chunks are sampled from one blob, we need to compute the sum of all x random field element corresponding to each sample
	aggCommitCoeffs := make([]fr.Element, m)
	setCommit := make([]bool, m)
	for k := 0; k < n; k++ {
		s := samples[k]
		row := s.RowIndex
		//bls.AddModFr(&aggCommitCoeffs[row], &aggCommitCoeffs[row], &randomsFr[k])
		aggCommitCoeffs[row].Add(&aggCommitCoeffs[row], &randomsFr[k])

		if !setCommit[row] {
			commits[row].Set(&s.Commitment)
			//bls.CopyG1(&commits[row], &s.Commitment)
			setCommit[row] = true
		} else {
			//bls.EqualG1(&commits[row], &s.Commitment)
			if !commits[row].Equal(&s.Commitment) {
				return nil, errors.New("Samples of the same row has different commitments")
			}
		}
	}

	var aggCommit bn254.G1Affine
	_, err := aggCommit.MultiExp(commits, aggCommitCoeffs, ecc.MultiExpConfig{})
	if err != nil {
		return nil, err
	}

	// second term
	// compute the aggregated interpolation polynomial
	aggPolyCoeffs := make([]fr.Element, D)

	// we sum over the weighted coefficients (by the random field element) over all D monomial in all n samples
	for k := 0; k < n; k++ {
		coeffs := samples[k].Coeffs

		rk := randomsFr[k]
		// for each monomial in a given polynomial, multiply its coefficient with the corresponding random field,
		// then sum it with others. Given ChunkLen (D) is identical for all samples in a subBatch.
		// The operation is always valid.
		for j := uint64(0); j < D; j++ {
			tmp.Mul(&coeffs[j], &rk)
			//bls.MulModFr(&tmp, &coeffs[j], &rk)
			//bls.AddModFr(&aggPolyCoeffs[j], &aggPolyCoeffs[j], &tmp)
			aggPolyCoeffs[j].Add(&aggPolyCoeffs[j], &tmp)
		}
	}

	// All samples in a subBatch has identical chunkLen
	var aggPolyG1 bn254.G1Affine
	_, err = aggPolyG1.MultiExp(ks.Srs.G1[:D], aggPolyCoeffs, ecc.MultiExpConfig{})
	if err != nil {
		return nil, err
	}
	//aggPolyG1 := bls.LinCombG1(ks.Srs.G1[:D], aggPolyCoeffs)

	// third term
	// leading coset is an evaluation index, here we compute the weighted leading coset evaluation by random fields
	lcCoeffs := make([]fr.Element, n)

	// get leading coset powers
	leadingDs := make([]fr.Element, n)

	for k := 0; k < n; k++ {

		// got the leading coset field element
		h := ks.ExpandedRootsOfUnity[samples[k].X]
		var hPow fr.Element
		hPow.SetOne()
		//bls.CopyFr(&hPow, &bls.ONE)

		// raising the power for each leading coset
		for j := uint64(0); j < D; j++ {
			hPow.Mul(&hPow, &h)
			//bls.MulModFr(&tmp, &hPow, &h)
			//bls.CopyFr(&hPow, &tmp)
		}
		//bls.CopyFr(&leadingDs[k], &hPow)
		leadingDs[k].Set(&hPow)
	}

	// applying the random weights to leading coset elements
	for k := 0; k < n; k++ {
		rk := randomsFr[k]
		//bls.MulModFr(&lcCoeffs[k], &rk, &leadingDs[k])
		lcCoeffs[k].Mul(&rk, &leadingDs[k])
	}

	var offsetG1 bn254.G1Affine
	_, err = offsetG1.MultiExp(proofs, lcCoeffs, ecc.MultiExpConfig{})
	if err != nil {
		return nil, err
	}

	//offsetG1 := bls.LinCombG1(proofs, lcCoeffs)

	var rhsG1 bn254.G1Affine
	//bls.SubG1(&rhsG1, aggCommit, aggPolyG1)
	rhsG1.Sub(&aggCommit, &aggPolyG1)
	//bls.AddG1(&rhsG1, &rhsG1, offsetG1)
	rhsG1.Add(&rhsG1, &offsetG1)
	return &rhsG1, nil
}

// TODO(mooselumph): Cleanup this function
func (v *Verifier) UniversalVerifySubBatch(params encoding.EncodingParams, samplesCore []encoding.Sample, numBlobs int) error {

	samples := make([]Sample, len(samplesCore))

	for i, sc := range samplesCore {
		x, err := rs.GetLeadingCosetIndex(
			uint64(sc.AssignmentIndex),
			params.NumChunks,
		)
		if err != nil {
			return err
		}

		sample := Sample{
			Commitment: (bn254.G1Affine)(*sc.Commitment),
			Proof:      sc.Chunk.Proof,
			RowIndex:   sc.BlobIndex,
			Coeffs:     sc.Chunk.Coeffs,
			X:          uint(x),
		}
		samples[i] = sample
	}

	return v.UniversalVerify(params, samples, numBlobs)
}

// UniversalVerify implements batch verification on a set of chunks given the same chunk dimension (chunkLen, numChunk).
// The details is given in Ethereum Research post whose authors are George Kadianakis, Ansgar Dietrichs, Dankrad Feist
// https://ethresear.ch/t/a-universal-verification-equation-for-data-availability-sampling/13240
//
// m is number of blob, samples is a list of chunks
//
// The order of samples do not matter.
// Each sample need not have unique row, it is possible that multiple chunks of the same blob are validated altogether
func (group *Verifier) UniversalVerify(params encoding.EncodingParams, samples []Sample, m int) error {
	// precheck
	for i, s := range samples {
		if s.RowIndex >= m {
			fmt.Printf("sample %v has %v Row, but there are only %v blobs\n", i, s.RowIndex, m)
			return errors.New("sample.RowIndex and numBlob are inconsistent")
		}
	}

	verifier, err := group.GetKzgVerifier(params)
	if err != nil {
		return err
	}
	ks := verifier.Ks

	D := params.ChunkLength

	if D > group.SRSNumberToLoad {
		return fmt.Errorf("requested chunkLen %v is larger than Loaded SRS points %v.", D, group.SRSNumberToLoad)
	}

	n := len(samples)
	fmt.Printf("Batch verify %v frames of %v symbols out of %v blobs \n", n, params.ChunkLength, m)

	// generate random field elements to aggregate equality check
	randomsFr, err := GenRandomnessVector(samples)
	if err != nil {
		return err
	}

	// array of proofs
	proofs := make([]bn254.G1Affine, n)
	for i := 0; i < n; i++ {
		//bls.CopyG1(&proofs[i], &samples[i].Proof)
		proofs[i].Set(&samples[i].Proof)
	}

	// lhs g1
	//lhsG1 := bls.LinCombG1(proofs, randomsFr)
	var lhsG1 bn254.G1Affine
	_, err = lhsG1.MultiExp(proofs, randomsFr, ecc.MultiExpConfig{})
	if err != nil {
		return err
	}
	// lhs g2
	exponent := uint64(math.Log2(float64(D)))
	G2atD, err := kzg.ReadG2PointOnPowerOf2(exponent, group.KzgConfig)

	if err != nil {
		// then try to access if there is a full list of g2 srs
		G2atD, err = kzg.ReadG2Point(D, group.KzgConfig)
		if err != nil {
			return err
		}
		fmt.Println("Acessed the entire G2")
	}

	lhsG2 := &G2atD

	// rhs g2
	rhsG2 := &kzg.GenG2

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

	return PairingsVerify(&lhsG1, lhsG2, rhsG1, rhsG2)
}
