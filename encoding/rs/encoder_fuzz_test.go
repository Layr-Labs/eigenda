package rs_test

import (
	"math"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/fft"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	rs_cpu "github.com/Layr-Labs/eigenda/encoding/rs/cpu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func FuzzOnlySystematic(f *testing.F) {

	f.Add(GETTYSBURG_ADDRESS_BYTES)
	f.Fuzz(func(t *testing.T, input []byte) {

		params := encoding.ParamsFromSysPar(10, 3, uint64(len(input)))
		enc, err := rs.NewEncoder(params, true)
		if err != nil {
			t.Errorf("Error making rs: %q", err)
		}

		n := uint8(math.Log2(float64(enc.NumEvaluations())))
		if enc.ChunkLength == 1 {
			n = uint8(math.Log2(float64(2 * enc.NumChunks)))
		}
		fs := fft.NewFFTSettings(n)

		RsComputeDevice := &rs_cpu.RsCpuComputeDevice{
			Fs:             fs,
			EncodingParams: params,
		}

		enc.Computer = RsComputeDevice
		require.NotNil(t, enc)

		//encode the data
		frames, _, err := enc.EncodeBytes(input)
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
