package verify

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding/rs"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/op-plasma-eigenda/eigenda"
)

type Verifier struct {
	cfg    *kzg.KzgConfig
	prover *prover.Prover
}

func NewVerifier(cfg *kzg.KzgConfig) (*Verifier, error) {

	prover, err := prover.NewProver(cfg, true)
	if err != nil {
		return nil, err
	}

	return &Verifier{
		cfg:    cfg,
		prover: prover,
	}, nil
}

// Verify regenerates a commitment from the blob and asserts equivalence
// to the commitment in the certificate
// TODO: Optimize implementation by opening a point on the commitment instead
func (v *Verifier) Verify(cert eigenda.Cert, blob []byte) error {
	encoder, err := v.prover.GetKzgEncoder(
		encoding.ParamsFromSysPar(6, 69, uint64(len(blob))),
	)
	if err != nil {
		return err
	}
	inputFr, err := rs.ToFrArray(blob)
	if err != nil {
		return err
	}

	polyEvals, _, err := encoder.ExtendPolyEval(inputFr)
	if err != nil {
		return err
	}
	commit, err := encoder.Commit(polyEvals)
	if err != nil {
		return err
	}

	x, y := cert.BlobCommitmentFields()

	if commit.X.NotEqual(x) == 0 {
		return fmt.Errorf("x element mismatch %s:%s %s:%s", "gen_commit", x.String(), "initial_commit", x.String())
	}

	if commit.Y.NotEqual(y) == 0 {
		return fmt.Errorf("x element mismatch %s:%s %s:%s", "gen_commit", y.String(), "initial_commit", y.String())
	}

	return nil
}
