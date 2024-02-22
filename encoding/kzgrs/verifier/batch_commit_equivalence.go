package verifier

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/Layr-Labs/eigenda/pkg/kzg/bn254"
)

type CommitmentPair struct {
	Commitment       bn254.G1Point
	LengthCommitment bn254.G2Point
}

// generate a random value using Fiat Shamir transform
// we can also pseudo randomness generated locally, but we have to ensure no adversary can manipulate it
// Hashing everything takes about 1ms, so Fiat Shamir transform does not incur much cost
func GenRandomFactorForEquivalence(g1commits []bn254.G1Point, g2commits []bn254.G2Point) (bn254.Fr, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	for _, commit := range g1commits {
		err := enc.Encode(commit)
		if err != nil {
			return bn254.ZERO, err
		}
	}

	for _, commit := range g2commits {
		err := enc.Encode(commit)
		if err != nil {
			return bn254.ZERO, err
		}
	}

	var randomFr bn254.Fr

	err := bn254.HashToSingleField(&randomFr, buffer.Bytes())
	if err != nil {
		return bn254.ZERO, err
	}

	return randomFr, nil
}

func CreateRandomnessVector(g1commits []bn254.G1Point, g2commits []bn254.G2Point) ([]bn254.Fr, error) {
	r, err := GenRandomFactorForEquivalence(g1commits, g2commits)
	if err != nil {
		return nil, err
	}
	n := len(g1commits)

	if len(g1commits) != len(g2commits) {
		return nil, errors.New("Inconsistent number of blobs for g1 and g2")
	}

	randomsFr := make([]bn254.Fr, n)
	bn254.CopyFr(&randomsFr[0], &r)

	// power of r
	for j := 0; j < n-1; j++ {
		bn254.MulModFr(&randomsFr[j+1], &randomsFr[j], &r)
	}

	return randomsFr, nil
}

func (group *Verifier) BatchVerifyCommitEquivalence(commitmentsPair []CommitmentPair) error {

	g1commits := make([]bn254.G1Point, len(commitmentsPair))
	g2commits := make([]bn254.G2Point, len(commitmentsPair))
	for i := 0; i < len(commitmentsPair); i++ {
		g1commits[i] = commitmentsPair[i].Commitment
		g2commits[i] = commitmentsPair[i].LengthCommitment
	}

	randomsFr, err := CreateRandomnessVector(g1commits, g2commits)
	if err != nil {
		return err
	}

	lhsG1 := bn254.LinCombG1(g1commits, randomsFr)
	lhsG2 := &bn254.GenG2

	rhsG2 := bn254.LinCombG2(g2commits, randomsFr)
	rhsG1 := &bn254.GenG1

	if bn254.PairingsVerify(lhsG1, lhsG2, rhsG1, rhsG2) {
		return nil
	} else {
		return errors.New("Universal Verify Incorrect paring")
	}
}
