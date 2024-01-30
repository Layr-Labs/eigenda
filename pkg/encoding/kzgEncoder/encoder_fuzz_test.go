package kzgEncoder_test

import (
	"testing"

	rs "github.com/Layr-Labs/eigenda/pkg/encoding/encoder"
	kzgRs "github.com/Layr-Labs/eigenda/pkg/encoding/kzgEncoder"
	"github.com/stretchr/testify/assert"
)

func FuzzOnlySystematic(f *testing.F) {

	f.Add(GETTYSBURG_ADDRESS_BYTES)
	f.Fuzz(func(t *testing.T, input []byte) {

		group, _ := kzgRs.NewKzgEncoderGroup(kzgConfig)

		params := rs.GetEncodingParams(10, 3, uint64(len(input)))
		enc, err := group.NewKzgEncoder(params)
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
