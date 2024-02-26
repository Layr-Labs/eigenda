package verifier

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

type CommitmentPair struct {
	Commitment       bn254.G1Affine
	LengthCommitment bn254.G2Affine
}

// generate a random value using Fiat Shamir transform
// we can also pseudo randomness generated locally, but we have to ensure no adversary can manipulate it
// Hashing everything takes about 1ms, so Fiat Shamir transform does not incur much cost
func GenRandomFactorForEquivalence(g1commits []bn254.G1Affine, g2commits []bn254.G2Affine) (fr.Element, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	for _, commit := range g1commits {
		err := enc.Encode(commit)
		if err != nil {
			return fr.Element{}, err
		}
	}

	for _, commit := range g2commits {
		err := enc.Encode(commit)
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

func CreateRandomnessVector(g1commits []bn254.G1Affine, g2commits []bn254.G2Affine) ([]fr.Element, error) {
	r, err := GenRandomFactorForEquivalence(g1commits, g2commits)
	if err != nil {
		return nil, err
	}
	n := len(g1commits)

	if len(g1commits) != len(g2commits) {
		return nil, errors.New("inconsistent number of blobs for g1 and g2")
	}

	randomsFr := make([]fr.Element, n)

	randomsFr[0].Set(&r)

	// power of r
	for j := 0; j < n-1; j++ {
		randomsFr[j+1].Mul(&randomsFr[j], &r)
	}

	return randomsFr, nil
}

func (v *Verifier) VerifyCommitEquivalenceBatch(commitments []encoding.BlobCommitments) error {
	commitmentsPair := make([]CommitmentPair, len(commitments))

	for i, c := range commitments {
		commitmentsPair[i] = CommitmentPair{
			Commitment:       (bn254.G1Affine)(*c.Commitment),
			LengthCommitment: (bn254.G2Affine)(*c.LengthCommitment),
		}
	}
	return v.BatchVerifyCommitEquivalence(commitmentsPair)
}

func (group *Verifier) BatchVerifyCommitEquivalence(commitmentsPair []CommitmentPair) error {

	g1commits := make([]bn254.G1Affine, len(commitmentsPair))
	g2commits := make([]bn254.G2Affine, len(commitmentsPair))
	for i := 0; i < len(commitmentsPair); i++ {
		g1commits[i] = commitmentsPair[i].Commitment
		g2commits[i] = commitmentsPair[i].LengthCommitment
	}

	randomsFr, err := CreateRandomnessVector(g1commits, g2commits)
	if err != nil {
		return err
	}

	var lhsG1 bn254.G1Affine
	_, err = lhsG1.MultiExp(g1commits, randomsFr, ecc.MultiExpConfig{})
	if err != nil {
		return err
	}

	lhsG2 := &kzg.GenG2

	var rhsG2 bn254.G2Affine
	_, err = rhsG2.MultiExp(g2commits, randomsFr, ecc.MultiExpConfig{})
	if err != nil {
		return err
	}
	rhsG1 := &kzg.GenG1

	err = PairingsVerify(&lhsG1, lhsG2, rhsG1, &rhsG2)
	if err == nil {
		return nil
	} else {
		return errors.New("Universal Verify Incorrect paring")
	}
}
