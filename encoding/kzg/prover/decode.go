package prover

import (
	enc "github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
)

func (g *ParametrizedProver) Decode(frames []enc.Frame, indices []uint64, maxInputSize uint64) ([]fr.Element, error) {
	rsFrames := make([]rs.Frame, len(frames))
	for ind, frame := range frames {
		rsFrames[ind] = rs.Frame{Coeffs: frame.Coeffs}
	}

	evals, err := g.Encoder.DecodeAsEval(rsFrames, indices, maxInputSize)
	if err != nil {
		return nil, err
	}

	return evals, nil
}

func (g *ParametrizedProver) DecodeBytes(frames []enc.Frame, indices []uint64, maxInputSize uint64) ([]byte, error) {
	evals, err := g.Decode(frames, indices, maxInputSize)
	if err != nil {
		return nil, err
	}
	data := rs.ToByteArray(evals, maxInputSize)
	return data, nil
}
