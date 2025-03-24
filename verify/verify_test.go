package verify

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"runtime"
	"testing"

	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	"github.com/Layr-Labs/eigenda/api/grpc/common"
	"github.com/Layr-Labs/eigenda/encoding/kzg"
	"github.com/Layr-Labs/eigenda/encoding/rs"
	"github.com/stretchr/testify/require"
)

func TestCommitmentVerification(t *testing.T) {
	t.Parallel()

	var data = []byte("inter-subjective and not objective!")

	x, err := hex.DecodeString("1021d699eac68ce312196d480266e8b82fd5fe5c4311e53313837b64db6df178")
	require.NoError(t, err)

	y, err := hex.DecodeString("02efa5a7813233ae13f32bae9b8f48252fa45c1b06a5d70bed471a9bea8d98ae")
	require.NoError(t, err)

	c := &common.G1Commitment{
		X: x,
		Y: y,
	}

	kzgConfig := kzg.KzgConfig{
		G1Path:          "../resources/g1.point",
		G2Path:          "../resources/g2.point",
		G2TrailingPath:  "../resources/g2.trailing.point",
		CacheDir:        "../resources/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	cfg := &Config{
		VerifyCerts: false,
	}

	v, err := NewVerifier(cfg, kzgConfig, nil)
	require.NoError(t, err)

	// Happy path verification
	codec := codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec())
	blob, err := codec.EncodeBlob(data)
	require.NoError(t, err)
	err = v.VerifyCommitment(c, blob)
	require.NoError(t, err)

	// failure with wrong data
	fakeData, err := codec.EncodeBlob([]byte("I am an imposter!!"))
	require.NoError(t, err)
	err = v.VerifyCommitment(c, fakeData)
	require.Error(t, err)
}

func TestCommitmentWithTooLargeBlob(t *testing.T) {

	var dataRand [2000 * 32]byte
	_, err := rand.Read(dataRand[:])
	require.NoError(t, err)
	data := dataRand[:]

	kzgConfig := kzg.KzgConfig{
		G1Path:          "../resources/g1.point",
		G2Path:          "../resources/g2.point",
		G2TrailingPath:  "../resources/g2.trailing.point",
		CacheDir:        "../resources/SRSTables",
		SRSOrder:        3000,
		SRSNumberToLoad: 3000,
		NumWorker:       uint64(runtime.GOMAXPROCS(0)),
	}

	cfg := &Config{
		VerifyCerts: false,
	}

	v, err := NewVerifier(cfg, kzgConfig, nil)
	require.NoError(t, err)

	// Some wrong commitment just to pass in function
	x, err := hex.DecodeString("1021d699eac68ce312196d480266e8b82fd5fe5c4311e53313837b64db6df178")
	require.NoError(t, err)

	y, err := hex.DecodeString("02efa5a7813233ae13f32bae9b8f48252fa45c1b06a5d70bed471a9bea8d98ae")
	require.NoError(t, err)

	c := &common.G1Commitment{
		X: x,
		Y: y,
	}

	// Happy path verification
	codec := codecs.NewIFFTCodec(codecs.NewDefaultBlobCodec())
	blob, err := codec.EncodeBlob(data)
	require.NoError(t, err)

	inputFr, err := rs.ToFrArray(blob)
	require.NoError(t, err)

	err = v.VerifyCommitment(c, blob)
	msg := fmt.Sprintf(
		"cannot verify commitment because the number of stored srs in the memory is insufficient, have %v need %v",
		kzgConfig.SRSNumberToLoad,
		len(inputFr),
	)
	require.EqualError(t, err, msg)

}
