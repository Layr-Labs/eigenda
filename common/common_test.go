package common_test

import (
	"encoding/hex"
	"testing"

	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/stretchr/testify/assert"
)

var (
	gettysburgAddressBytes = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
)

func TestPrefixEnvVar(t *testing.T) {
	assert.Equal(t, "prefix_suffix", common.PrefixEnvVar("prefix", "suffix"))
}

func TestHashBlob(t *testing.T) {
	blob := &core.Blob{
		RequestHeader: core.BlobRequestHeader{
			SecurityParams: []*core.SecurityParam{
				{
					QuorumID:           0,
					AdversaryThreshold: 80,
				},
			},
		},
		Data: gettysburgAddressBytes,
	}
	blobHash, err := common.Hash[*core.Blob](blob)
	blobKey := hex.EncodeToString(blobHash)
	assert.Nil(t, err)
	assert.Len(t, blobKey, 64)
}

func TestHash(t *testing.T) {
	hash, err := common.Hash[string]("test")
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x6f, 0xe3, 0x18, 0xf, 0x70, 0x0, 0x90, 0x69, 0x72, 0x85, 0xac, 0x1e, 0xe, 0x8d, 0xc4, 0x0, 0x25, 0x93, 0x73, 0xd7, 0xbb, 0x94, 0xf0, 0xb1, 0xa9, 0xb0, 0x86, 0xe7, 0xba, 0x22, 0xdc, 0x3d}, hash)
}

func TestEncodeToBytes(t *testing.T) {
	bytes, err := common.EncodeToBytes[string]("test")
	assert.Nil(t, err)
	assert.Equal(t, []byte{0x64, 0x74, 0x65, 0x73, 0x74}, bytes)
}

func TestDecodeFromBytes(t *testing.T) {
	str, err := common.DecodeFromBytes[string]([]byte{0x64, 0x74, 0x65, 0x73, 0x74})
	assert.Nil(t, err)
	assert.Equal(t, "test", str)
}

func TestEncodeDecode(t *testing.T) {
	s := "test"
	bytes, err := common.EncodeToBytes[string](s)
	assert.Nil(t, err)
	str, err := common.DecodeFromBytes[string](bytes)
	assert.Nil(t, err)
	assert.Equal(t, s, str)
}

func TestEncodeDecodeStruct(t *testing.T) {
	type testStruct struct {
		A string
		B int
	}
	s := testStruct{"test", 1}
	bytes, err := common.EncodeToBytes[testStruct](s)
	assert.Nil(t, err)
	str, err := common.DecodeFromBytes[testStruct](bytes)
	assert.Nil(t, err)
	assert.Equal(t, s, str)
}

func TestEncodeDecodeStructWithSlice(t *testing.T) {
	type testStruct struct {
		A []string
		B int
	}
	s := testStruct{[]string{"test", "test2"}, 1}
	bytes, err := common.EncodeToBytes[testStruct](s)
	assert.Nil(t, err)
	str, err := common.DecodeFromBytes[testStruct](bytes)
	assert.Nil(t, err)
	assert.Equal(t, s, str)
}

func TestEncodeDecodeStructWithMap(t *testing.T) {
	type testStruct struct {
		A map[string]string
		B int
	}
	s := testStruct{map[string]string{"test": "test", "test2": "test2"}, 1}
	bytes, err := common.EncodeToBytes[testStruct](s)
	assert.Nil(t, err)
	str, err := common.DecodeFromBytes[testStruct](bytes)
	assert.Nil(t, err)
	assert.Equal(t, s, str)
}

func TestEncodeDecodeStructWithPointer(t *testing.T) {
	type testStruct struct {
		A *string
		B int
	}
	p := "test"
	s := testStruct{&p, 1}
	bytes, err := common.EncodeToBytes[testStruct](s)
	assert.Nil(t, err)
	str, err := common.DecodeFromBytes[testStruct](bytes)
	assert.Nil(t, err)
	assert.Equal(t, s, str)
}
