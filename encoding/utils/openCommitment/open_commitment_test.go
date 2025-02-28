package openCommitment_test

import (
	"crypto/rand"
	"log"
	"math/big"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/encoding"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/kzg/prover"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/Layr-Labs/eigenda/encoding/utils/codec"
	oc "github.com/Layr-Labs/eigenda/encoding/utils/openCommitment"

	"github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/stretchr/testify/require"
)

var (
	gettysburgAddressBytes = []byte("Fourscore and seven years ago our fathers brought forth, on this continent, a new nation, conceived in liberty, and dedicated to the proposition that all men are created equal. Now we are engaged in a great civil war, testing whether that nation, or any nation so conceived, and so dedicated, can long endure. We are met on a great battle-field of that war. We have come to dedicate a portion of that field, as a final resting-place for those who here gave their lives, that that nation might live. It is altogether fitting and proper that we should do this. But, in a larger sense, we cannot dedicate, we cannot consecrate—we cannot hallow—this ground. The brave men, living and dead, who struggled here, have consecrated it far above our poor power to add or detract. The world will little note, nor long remember what we say here, but it can never forget what they did here. It is for us the living, rather, to be dedicated here to the unfinished work which they who fought here have thus far so nobly advanced. It is rather for us to be here dedicated to the great task remaining before us—that from these honored dead we take increased devotion to that cause for which they here gave the last full measure of devotion—that we here highly resolve that these dead shall not have died in vain—that this nation, under God, shall have a new birth of freedom, and that government of the people, by the people, for the people, shall not perish from the earth.")
	kzgConfig              *kzg.KzgConfig
	numNode                uint64
	numSys                 uint64
	numPar                 uint64
)

func TestOpenCommitment(t *testing.T) {
	log.Println("Setting up suite")

	kzgConfig = &kzg.KzgConfig{
		G1Path:          "../../../inabox/resources/kzg/g1.point",
		G2Path:          "../../../inabox/resources/kzg/g2.point",
		G2PowerOf2Path:  "../../../inabox/resources/kzg/g2.point.powerOf2",
		CacheDir:        "../../../inabox/resources/kzg/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
		LoadG2Points:    true,
	}

	// input evaluation
	validInput := codec.ConvertByPaddingEmptyByte(gettysburgAddressBytes)
	inputFr, err := rs.ToFrArray(validInput)
	require.Nil(t, err)

	frLen := uint64(len(inputFr))
	paddedInputFr := make([]fr.Element, encoding.NextPowerOf2(frLen))
	// pad input Fr to power of 2 for computing FFT
	for i := 0; i < len(paddedInputFr); i++ {
		if i < len(inputFr) {
			paddedInputFr[i].Set(&inputFr[i])
		} else {
			paddedInputFr[i].SetZero()
		}
	}

	// we need prover only to access kzg SRS, and get kzg commitment of encoding
	group, err := prover.NewProver(kzgConfig, nil)
	require.NoError(t, err)

	// get root of unit for blob
	numNode = 4
	numSys = 4
	numPar = 0
	numOpenChallenge := 10

	params := encoding.ParamsFromSysPar(numSys, numPar, uint64(len(validInput)))
	enc, err := group.GetKzgEncoder(params)
	require.Nil(t, err)

	rs, err := enc.GetRsEncoder(params)
	require.NoError(t, err)

	rootOfUnities := rs.Fs.ExpandedRootsOfUnity[:len(rs.Fs.ExpandedRootsOfUnity)-1]

	// Lagrange basis SRS in normal order, not butterfly
	lagrangeG1SRS, err := rs.Fs.FFTG1(group.Srs.G1[:len(paddedInputFr)], true)
	require.Nil(t, err)

	// commit in lagrange form
	commitLagrange, err := oc.CommitInLagrange(paddedInputFr, lagrangeG1SRS)
	require.Nil(t, err)

	modulo := big.NewInt(int64(len(inputFr)))
	// pick a random place in the blob to open
	for k := 0; k < numOpenChallenge; k++ {

		indexBig, err := rand.Int(rand.Reader, modulo)
		require.Nil(t, err)

		index := int(indexBig.Int64())

		// open at index on the kzg
		proof, valueFr, err := oc.ComputeKzgProof(paddedInputFr, index, lagrangeG1SRS, rootOfUnities)
		require.Nil(t, err)

		_, _, g1Gen, g2Gen := bn254.Generators()

		err = oc.VerifyKzgProof(g1Gen, *commitLagrange, *proof, g2Gen, group.Srs.G2[1], *valueFr, rootOfUnities[index])
		require.Nil(t, err)

		require.Equal(t, *valueFr, inputFr[index])

		//valueBytse := valueFr.Bytes()
		//fmt.Println("value Byte", string(valueBytse[1:]))
	}
}
