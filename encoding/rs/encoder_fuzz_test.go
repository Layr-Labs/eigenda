package rs_test

import (
	"testing"

	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/assert"
)

func FuzzOnlySystematic(f *testing.F) {

	f.Add(GETTYSBURG_ADDRESS_BYTES)
	f.Fuzz(func(t *testing.T, input []byte) {

		params := rs.ParamsFromSysPar(10, 3, uint64(len(input)))
		enc, err := rs.NewEncoder(params, true)
		if err != nil {
			t.Errorf("Error making rs: %q", err)
		}

		//encode the data
		_, frames, _, err := enc.EncodeBytes(input)
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
