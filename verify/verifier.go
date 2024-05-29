package verify

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/encoding"

	"github.com/Layr-Labs/eigenda-proxy/eigenda"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

type Verifier struct {
	cfg    *kzg.KzgConfig
	prover *prover.Prover
}

func NewVerifier(cfg *kzg.KzgConfig) (*Verifier, error) {

	prover, err := prover.NewProver(cfg, false) // don't load G2 points
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
	// ChunkLength and TotalChunks aren't relevant for computing data
	// commitment which is why they're currently set arbitrarily
	encoder, err := v.prover.GetKzgEncoder(
		encoding.ParamsFromSysPar(420, 69, uint64(len(blob))),
	)
	if err != nil {
		return err
	}

	inputFr, err := rs.ToFrArray(blob)
	if err != nil {
		return fmt.Errorf("cannot convert bytes to field elements, %w", err)
	}

	poly, _, _, err := encoder.Encoder.Encode(inputFr)
	if err != nil {
		return err
	}

	commit, err := encoder.Commit(poly.Coeffs)
	if err != nil {
		return err
	}

	x, y := cert.BlobCommitmentFields()

	errMsg := ""
	if !commit.X.Equal(x) || !commit.Y.Equal(y) {
		errMsg += fmt.Sprintf("field elements do not match, x generated commit: %x, x initial commit: %x, ", commit.X.Marshal(), (*x).Marshal())
		errMsg += fmt.Sprintf("y generated commit: %x, y initial commit: %x", commit.Y.Marshal(), (*y).Marshal())
		return fmt.Errorf(errMsg)
	}

	return nil
}
