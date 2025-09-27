package rs_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	GETTYSBURG_ADDRESS_BYTES = codec.ConvertByPaddingEmptyByte([]byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth."))
	numNode                  = uint64(4)
	numSys                   = uint64(3)
	numPar                   = numNode - numSys
)

func TestEncodeDecode_InvertsWhenSamplingAllFrames(t *testing.T) {
	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	inputFr, err := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	assert.Nil(t, err)
	frames, _, err := enc.Encode(inputFr, params)
	assert.Nil(t, err)

	// sample some Frames
	samples, indices := sampleFrames(frames, uint64(len(frames)))
	data, err := enc.Decode(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)), params)

	require.Nil(t, err)
	require.NotNil(t, data)

	assert.Equal(t, data, GETTYSBURG_ADDRESS_BYTES)
}

func TestEncodeDecode_InvertsWhenSamplingMissingFrame(t *testing.T) {
	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	inputFr, err := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	assert.Nil(t, err)
	frames, _, err := enc.Encode(inputFr, params)
	assert.Nil(t, err)

	// sample some Frames
	samples, indices := sampleFrames(frames, uint64(len(frames)-1))
	data, err := enc.Decode(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)), params)

	require.Nil(t, err)
	require.NotNil(t, data)

	assert.Equal(t, data, GETTYSBURG_ADDRESS_BYTES)
}

func TestEncodeDecode_InvertsWithMissingAndDuplicateFrames(t *testing.T) {
	numSys := uint64(3)
	numPar := uint64(5)
	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))

	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	inputFr, err := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	assert.Nil(t, err)
	frames, _, err := enc.Encode(inputFr, params)
	assert.Nil(t, err)

	assert.EqualValues(t, len(frames), numSys+numPar)

	// sample some Frames
	samples, indices := sampleFrames(frames, uint64(len(frames))-numPar)

	// duplicate two of the frames
	samples = append(samples, samples[0:2]...)
	indices = append(indices, indices[0:2]...)

	data, err := enc.Decode(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)), params)

	require.Nil(t, err)
	require.NotNil(t, data)

	assert.Equal(t, data, GETTYSBURG_ADDRESS_BYTES)
}

func TestEncodeDecode_ErrorsWhenNotEnoughSampledFrames(t *testing.T) {
	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	fmt.Println("Num Chunks: ", params.NumChunks)

	inputFr, err := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	assert.Nil(t, err)
	frames, _, err := enc.Encode(inputFr, params)
	assert.Nil(t, err)

	// sample some Frames
	samples, indices := sampleFrames(frames, uint64(len(frames)-2))
	data, err := enc.Decode(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)), params)

	require.Nil(t, data)
	require.NotNil(t, err)

	assert.EqualError(t, err, "number of frame must be sufficient")
}

func TestEncodeDecode_ErrorsWhenNotEnoughSampledFramesWithDuplicates(t *testing.T) {
	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(GETTYSBURG_ADDRESS_BYTES)))
	cfg := encoding.DefaultConfig()
	enc, err := rs.NewEncoder(cfg)
	assert.Nil(t, err)

	fmt.Println("Num Chunks: ", params.NumChunks)

	inputFr, err := rs.ToFrArray(GETTYSBURG_ADDRESS_BYTES)
	assert.Nil(t, err)
	frames, _, err := enc.Encode(inputFr, params)
	assert.Nil(t, err)

	// sample some Frames
	samples, indices := sampleFrames(frames, uint64(len(frames)-2))

	// duplicate two of the frames
	samples = append(samples, samples[0:2]...)
	indices = append(indices, indices[0:2]...)

	data, err := enc.Decode(samples, indices, uint64(len(GETTYSBURG_ADDRESS_BYTES)), params)

	require.Nil(t, data)
	require.NotNil(t, err)

	assert.EqualError(t, err, "number of frame must be sufficient")
}

func sampleFrames(frames []rs.FrameCoeffs, num uint64) ([]rs.FrameCoeffs, []uint64) {
	samples := make([]rs.FrameCoeffs, num)
	indices := rand.Perm(len(frames))
	indices = indices[:num]

	frameIndices := make([]uint64, num)
	for i, j := range indices {
		samples[i] = frames[j]
		frameIndices[i] = uint64(j)
	}
	return samples, frameIndices
}

func FuzzOnlySystematic(f *testing.F) {

	f.Add(GETTYSBURG_ADDRESS_BYTES)
	f.Fuzz(func(t *testing.T, input []byte) {

		params := encoding.ParamsFromSysPar(10, 3, uint64(len(input)))
		cfg := encoding.DefaultConfig()
		enc, err := rs.NewEncoder(cfg)
		if err != nil {
			t.Errorf("Error making rs: %q", err)
		}

		//encode the data
		frames, _, err := enc.EncodeBytes(input, params)
		if err != nil {
			t.Errorf("Error Encoding:\n Data:\n %q \n Err: %q", input, err)
		}

		//sample the correct systematic Frames
		samples, indices := sampleFrames(frames, uint64(len(frames)))

		data, err := enc.Decode(samples, indices, uint64(len(input)), params)
		if err != nil {
			t.Errorf("Error Decoding:\n Data:\n %q \n Err: %q", input, err)
		}
		assert.Equal(t, input, data, "Input data was not equal to the decoded data")
	})
}
