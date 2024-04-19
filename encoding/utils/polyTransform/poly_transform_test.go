package polyTranform_test

import (
	"math"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	polyTranform "github.com/Layr-Labs/eigenda/encoding/utils/polyTransform"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/require"
)

func TestPolyToEvalsTransform_and_Back(t *testing.T) {
	gettysburgAddressBytes := []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")

	paddedData := codec.ConvertByPaddingEmptyByte(gettysburgAddressBytes)

	paddedFr, err := rs.ToFrArray(paddedData)
	require.Nil(t, err)

	l := uint8(math.Ceil(math.Log2(float64(len(paddedFr)))))

	transformer, err := polyTranform.NewPolyTranform(l)
	require.Nil(t, err)

	coeffsFr, err := transformer.ConvertEvalsToCoeffs(paddedFr)
	require.Nil(t, err)

	evalsFr, err := transformer.ConvertCoeffsToEvals(coeffsFr)
	require.Nil(t, err)

	restoredPaddedData := rs.ToByteArray(evalsFr, uint64(len(paddedData)))

	restored := codec.RemoveEmptyByteFromPaddedBytes(restoredPaddedData)

	require.Equal(t, gettysburgAddressBytes, restored[:len(gettysburgAddressBytes)])
}

func TestPolyToEvalsTransform_and_Oversize(t *testing.T) {
	gettysburgAddressBytes := []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")

	paddedData := codec.ConvertByPaddingEmptyByte(gettysburgAddressBytes)

	paddedFr, err := rs.ToFrArray(paddedData)
	require.Nil(t, err)

	l := uint8(math.Ceil(math.Log2(float64(len(paddedFr)))))

	transformer1, err := polyTranform.NewPolyTranform(l)
	require.Nil(t, err)

	transformer2, err := polyTranform.NewPolyTranform(2 * l)
	require.Nil(t, err)

	coeffsFr1, err := transformer1.ConvertEvalsToCoeffs(paddedFr)
	require.Nil(t, err)

	coeffsFr2, err := transformer2.ConvertEvalsToCoeffs(paddedFr)
	require.Nil(t, err)

	require.Equal(t, coeffsFr1, coeffsFr2)
}

func TestPolyToEvalsTransform_and_IncreaseCapacity(t *testing.T) {
	gettysburgAddressBytes := []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
	paddedData := codec.ConvertByPaddingEmptyByte(gettysburgAddressBytes)

	paddedFr, err := rs.ToFrArray(paddedData)
	require.Nil(t, err)

	l := uint8(math.Ceil(math.Log2(float64(len(paddedFr)))))

	transformer, err := polyTranform.NewPolyTranform(l)
	require.Nil(t, err)

	largerData := make([]fr.Element, 0)
	for i := 0; i < 4; i++ {
		largerData = append(largerData, paddedFr...)
	}

	coeffsFr, err := transformer.ConvertEvalsToCoeffs(largerData)
	require.Nil(t, err)

	evalsFr, err := transformer.ConvertCoeffsToEvals(coeffsFr)
	require.Nil(t, err)

	restoredPaddedData := rs.ToByteArray(evalsFr, uint64(4*len(paddedData)))

	restored := codec.RemoveEmptyByteFromPaddedBytes(restoredPaddedData)

	require.Equal(t, gettysburgAddressBytes, restored[:len(gettysburgAddressBytes)])
}
