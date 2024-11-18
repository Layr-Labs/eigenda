package prover_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FuzzOnlySystematic(f *testing.F) {

	f.Add(gettysburgAddressBytes)
	f.Fuzz(func(t *testing.T, input []byte) {
		group, err := prover.NewProver(kzgConfig, nil)
		require.NoError(t, err)

		params := encoding.ParamsFromSysPar(10, 3, uint64(len(input)))
		enc, err := group.GetKzgEncoder(params)
		if err != nil {
			t.Errorf("Error making rs: %q", err)
		}

		//encode the data
		_, _, _, frames, _, err := enc.EncodeBytes(input)

		for _, frame := range frames {
			assert.NotEqual(t, len(frame.Coeffs), 0)
		}

		if err != nil {
			t.Errorf("Error Encoding:\n Data:\n %q \n Err: %q", input, err)
		}

		//sample the correct systematic frames
		samples, indices := sampleFrames(frames, uint64(len(frames)))

		data, err := enc.Decode(samples, indices, uint64(len(input)))
		if err != nil {
			t.Errorf("Error Decoding:\n Data:\n %q \n Err: %q", input, err)
		}
		assert.Equal(t, input, data, "Input data was not equal to the decoded data")
	})
}
