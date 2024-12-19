package verification

import (
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/verifier"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"runtime"
	"testing"
)

var (
	gettysburgAddressBytes = codec.ConvertByPaddingEmptyByte(
		[]byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth."))

	kzgConfig       *kzg.KzgConfig
	srsNumberToLoad uint64
)

func setup() {
	log.Println("Setting up suite")

	srsNumberToLoad = 2900

	kzgConfig = &kzg.KzgConfig{
		G1Path:          "../../../../inabox/resources/kzg/g1.point",
		G2Path:          "../../../../inabox/resources/kzg/g2.point",
		G2PowerOf2Path:  "../../../../inabox/resources/kzg/g2.point.powerOf2",
		CacheDir:        "../../../../inabox/resources/kzg/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: srsNumberToLoad,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    false,
	}
}

// randomlyModifyBytes picks a random byte from the input array, and increments it
func randomlyModifyBytes(testRandom *random.TestRandom, inputBytes []byte) {
	indexToModify := testRandom.Intn(len(inputBytes))
	inputBytes[indexToModify] = inputBytes[indexToModify] + 1
}

func TestMain(m *testing.M) {
	setup()
	result := m.Run()
	teardown()
	os.Exit(result)
}

func teardown() {
	log.Println("Tearing down")
}

func TestComputeAndCompareKzgCommitmentSuccess(t *testing.T) {
	kzgVerifier, err := verifier.NewVerifier(kzgConfig, nil)
	assert.NotNil(t, kzgVerifier)
	assert.Nil(t, err)

	commitment, err := GenerateBlobCommitment(kzgVerifier, gettysburgAddressBytes)
	assert.NotNil(t, commitment)
	assert.Nil(t, err)

	// make sure the commitment verifies correctly
	err = GenerateAndCompareBlobCommitment(
		kzgVerifier,
		commitment,
		gettysburgAddressBytes)
	assert.Nil(t, err)
}

func TestComputeAndCompareKzgCommitmentFailure(t *testing.T) {
	kzgVerifier, err := verifier.NewVerifier(kzgConfig, nil)
	assert.NotNil(t, kzgVerifier)
	assert.Nil(t, err)

	commitment, err := GenerateBlobCommitment(kzgVerifier, gettysburgAddressBytes)
	assert.NotNil(t, commitment)
	assert.Nil(t, err)

	// randomly modify the bytes, and make sure the commitment verification fails
	testRandom := random.NewTestRandom(t)
	randomlyModifyBytes(testRandom, gettysburgAddressBytes)
	err = GenerateAndCompareBlobCommitment(
		kzgVerifier,
		commitment,
		gettysburgAddressBytes)
	assert.NotNil(t, err)
}

func TestGenerateBlobCommitmentEquality(t *testing.T) {
	kzgVerifier, err := verifier.NewVerifier(kzgConfig, nil)
	assert.NotNil(t, kzgVerifier)
	assert.Nil(t, err)

	// generate two identical commitments
	commitment1, err := GenerateBlobCommitment(kzgVerifier, gettysburgAddressBytes)
	assert.NotNil(t, commitment1)
	assert.Nil(t, err)
	commitment2, err := GenerateBlobCommitment(kzgVerifier, gettysburgAddressBytes)
	assert.NotNil(t, commitment2)
	assert.Nil(t, err)

	// commitments to identical bytes should be equal
	assert.Equal(t, commitment1, commitment2)

	// randomly modify a byte
	testRandom := random.NewTestRandom(t)
	randomlyModifyBytes(testRandom, gettysburgAddressBytes)
	commitmentA, err := GenerateBlobCommitment(kzgVerifier, gettysburgAddressBytes)
	assert.NotNil(t, commitmentA)
	assert.Nil(t, err)

	// commitments to non-identical bytes should not be equal
	assert.NotEqual(t, commitment1, commitmentA)
}

func TestGenerateBlobCommitmentTooLong(t *testing.T) {
	kzgVerifier, err := verifier.NewVerifier(kzgConfig, nil)
	assert.NotNil(t, kzgVerifier)
	assert.Nil(t, err)

	// this is the absolute maximum number of bytes we can handle, given how the verifier was configured
	almostTooLongByteCount := srsNumberToLoad * 32

	// an array of exactly this size should be fine
	almostTooLongBytes := make([]byte, almostTooLongByteCount)
	commitment1, err := GenerateBlobCommitment(kzgVerifier, almostTooLongBytes)
	assert.NotNil(t, commitment1)
	assert.Nil(t, err)

	// but 1 more byte is more than we can handle
	tooLongBytes := make([]byte, almostTooLongByteCount+1)
	commitment2, err := GenerateBlobCommitment(kzgVerifier, tooLongBytes)
	assert.Nil(t, commitment2)
	assert.NotNil(t, err)
}
