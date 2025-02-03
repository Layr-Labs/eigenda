package prover

import (
	enc "github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
)

func (g *ParametrizedProver) Decode(frames []enc.Frame, indices []uint64, maxInputSize uint64) ([]byte, error) {
	rsFrames := make([]rs.FrameCoeffs, len(frames))
	for ind, frame := range frames {
		rsFrames[ind] = frame.Coeffs
	}

	return g.Encoder.Decode(rsFrames, indices, maxInputSize, g.EncodingParams)
}
