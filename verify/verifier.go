package verify

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fp"

	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

type Verifier struct {
	prover *prover.Prover
}

func NewVerifier(cfg *kzg.KzgConfig) (*Verifier, error) {
	prover, err := prover.NewProver(cfg, false) // don't load G2 points
	if err != nil {
		return nil, err
	}

	return &Verifier{
		prover: prover,
	}, nil
}

func (v *Verifier) Commit(blob []byte) (*bn254.G1Affine, error) {
	// ChunkLength and TotalChunks aren't relevant for computing data
	// commitment which is why they're currently set arbitrarily
	encoder, err := v.prover.GetKzgEncoder(
		encoding.ParamsFromSysPar(420, 69, uint64(len(blob))),
	)
	if err != nil {
		return nil, err
	}

	inputFr, err := rs.ToFrArray(blob)
	if err != nil {
		return nil, fmt.Errorf("cannot convert bytes to field elements, %w", err)
	}

	commit, err := encoder.Commit(inputFr)
	if err != nil {
		return nil, err
	}

	return &commit, nil
}

// Verify regenerates a commitment from the blob and asserts equivalence
// to the commitment in the certificate
// TODO: Optimize implementation by opening a point on the commitment instead
func (v *Verifier) Verify(expectedCommit *common.G1Commitment, blob []byte) error {
	actualCommit, err := v.Commit(blob)
	if err != nil {
		return err
	}

	// convert to field elements
	expectedX := &fp.Element{}
	expectedX.Unmarshal(expectedCommit.X)
	expectedY := &fp.Element{}
	expectedY.Unmarshal(expectedCommit.Y)

	errMsg := ""
	if !actualCommit.X.Equal(expectedX) || !actualCommit.Y.Equal(expectedY) {
		errMsg += fmt.Sprintf("field elements do not match, x actual commit: %x, x expected commit: %x, ", actualCommit.X.Marshal(), (*expectedX).Marshal())
		errMsg += fmt.Sprintf("y actual commit: %x, y expected commit: %x", actualCommit.Y.Marshal(), (*expectedY).Marshal())
		return fmt.Errorf(errMsg)
	}

	return nil
}
