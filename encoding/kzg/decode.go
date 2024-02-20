package kzgEncoder

import (
	rs "github.com/Layr-Labs/eigenda/encoding/rs"
)

func (g *KzgEncoder) Decode(frames []Frame, indices []uint64, maxInputSize uint64) ([]byte, error) {
	rsFrames := make([]rs.Frame, len(frames))
	for ind, frame := range frames {
		rsFrames[ind] = rs.Frame{Coeffs: frame.Coeffs}
	}

	return g.Encoder.Decode(rsFrames, indices, maxInputSize)
}
