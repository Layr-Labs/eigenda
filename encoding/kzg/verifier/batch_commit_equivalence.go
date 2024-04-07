package verifier

import (
	"crypto/rand"
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

// Create a random number with crypto/rand.
// Gnark provides SetRandom() function, but the implementation below is for explicity
func GetRandomFr() (fr.Element, error) {
	r, err := rand.Int(rand.Reader, fr.Modulus())
	if err != nil {
		return fr.Element{}, err
	}
	var rElement fr.Element
	rElement.SetBigInt(r)
	return rElement, nil
}

func CreateRandomnessVector(n int) ([]fr.Element, error) {
	if n <= 0 {
		return nil, errors.New("the length of vector must be positive")
	}
	r, err := GetRandomFr()
	if err != nil {
		return nil, err
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

	randomsFr, err := CreateRandomnessVector(len(g1commits))
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
		return errors.New("incorrect universal batch verification")
	}
}
