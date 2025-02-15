package random

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

// Tests that random seeding produces random results, and that consistent seeding produces consistent results
func TestSetup(t *testing.T) {
	testRandom1 := NewTestRandom()
	x := testRandom1.Int()

	testRandom2 := NewTestRandom()
	y := testRandom2.Int()

	require.NotEqual(t, x, y)

	seed := rand.Int63()
	testRandom3 := NewTestRandom(seed)
	a := testRandom3.Int()

	testRandom4 := NewTestRandom(seed)
	b := testRandom4.Int()

	require.Equal(t, a, b)
}

func TestReset(t *testing.T) {
	random := NewTestRandom()

	a := random.Uint64()
	b := random.Uint64()
	c := random.Uint64()
	d := random.Uint64()

	random.Reset()

	require.Equal(t, a, random.Uint64())
	require.Equal(t, b, random.Uint64())
	require.Equal(t, c, random.Uint64())
	require.Equal(t, d, random.Uint64())
}

func TestECDSAKeyGeneration(t *testing.T) {
	random := NewTestRandom()

	// We should not get the same key pair twice in a row
	public1, private1, err := random.ECDSA()
	require.NoError(t, err)
	public2, private2, err := random.ECDSA()
	require.NoError(t, err)

	assert.NotEqual(t, &public1, &public2)
	assert.NotEqual(t, &private1, &private2)

	// Getting keys should result in deterministic generator state.
	generatorState := random.Uint64()
	random.Reset()
	_, _, err = random.ECDSA()
	require.NoError(t, err)
	_, _, err = random.ECDSA()
	require.NoError(t, err)
	require.Equal(t, generatorState, random.Uint64())

	// Keypair should be valid.
	data := random.Bytes(32)

	signature, err := crypto.Sign(data, private1)
	require.NoError(t, err)

	signingPublicKey, err := crypto.SigToPub(data, signature)
	require.NoError(t, err)
	require.Equal(t, &public1, &signingPublicKey)
}

func TestBLSKeyGeneration(t *testing.T) {
	random := NewTestRandom()

	// We should not get the same key pair twice in a row
	keypair1, err := random.BLS()
	require.NoError(t, err)
	keypair2, err := random.BLS()
	require.NoError(t, err)

	require.False(t, keypair1.PrivKey.Equal(keypair2.PrivKey))
	require.False(t, keypair1.PubKey.Equal(keypair2.PubKey.G1Affine))

	// Getting keys should result in deterministic generator state.
	generatorState := random.Uint64()
	random.Reset()
	_, err = random.BLS()
	require.NoError(t, err)
	_, err = random.BLS()
	require.NoError(t, err)
	require.Equal(t, generatorState, random.Uint64())

	// Keys should be deterministic.
	random.Reset()
	keypair3, err := random.BLS()
	require.NoError(t, err)
	require.True(t, keypair1.PrivKey.Equal(keypair3.PrivKey))
	require.True(t, keypair1.PubKey.Equal(keypair3.PubKey.G1Affine))

	// Keypair should be valid.
	data := random.Bytes(32)
	signature := keypair1.SignMessage(([32]byte)(data))

	isValid := signature.Verify(keypair1.GetPubKeyG2(), ([32]byte)(data))
	require.True(t, isValid)
}
