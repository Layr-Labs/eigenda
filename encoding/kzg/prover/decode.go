package prover

import (
	enc "github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

func (g *ParametrizedProver) Decode(frames []enc.Frame, indices []uint64, maxInputSize uint64) ([]byte, error) {
	rsFrames := make([]rs.Frame, len(frames))
	for ind, frame := range frames {
		rsFrames[ind] = rs.Frame{Coeffs: frame.Coeffs}
	}

	return g.Encoder.Decode(rsFrames, indices, maxInputSize, g.EncodingParams)
}
