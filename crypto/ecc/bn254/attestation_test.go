package bn_test

import (
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestNewG1Point(t *testing.T) {
	x := big.NewInt(1)
	y := big.NewInt(2)
	p := core.NewG1Point(x, y)

	assert.Equal(t, x, p.X.BigInt(new(big.Int)))
	assert.Equal(t, y, p.Y.BigInt(new(big.Int)))
}

func TestG1Point_Add(t *testing.T) {
	x1, y1 := big.NewInt(1), big.NewInt(2)
	x2, y2 := big.NewInt(3), big.NewInt(4)
	p1 := core.NewG1Point(x1, y1)
	p2 := core.NewG1Point(x2, y2)

	p1.Add(p2)

	assert.NotNil(t, p1)
}

func TestG1Point_Sub(t *testing.T) {
	x1, y1 := big.NewInt(1), big.NewInt(2)
	x2, y2 := big.NewInt(3), big.NewInt(4)
	p1 := core.NewG1Point(x1, y1)
	p2 := core.NewG1Point(x2, y2)

	p1.Sub(p2)

	assert.NotNil(t, p1)
}

func TestG1Point_SerializeDeserialize(t *testing.T) {
	x := big.NewInt(1)
	y := big.NewInt(2)
	p := core.NewG1Point(x, y)
	serialized := p.Serialize()

	deserialized, err := p.Deserialize(serialized)
	assert.NoError(t, err)
	assert.Equal(t, p, deserialized)
}

func TestKeyPair_SignAndVerify(t *testing.T) {
	sk, err := core.GenRandomBlsKeys()
	assert.NoError(t, err)

	message := [32]byte{1, 2, 3, 4}
	sig := sk.SignMessage(message)

	pubkey := sk.GetPubKeyG2()
	assert.True(t, sig.Verify(pubkey, message))
}

func TestKeyPair_GetOperatorID(t *testing.T) {
	sk, err := core.GenRandomBlsKeys()
	assert.NoError(t, err)

	operatorID := sk.GetPubKeyG1().GetOperatorID()
	assert.NotEmpty(t, operatorID)
}

func TestMakePubkeyRegistrationData(t *testing.T) {
	sk, err := core.GenRandomBlsKeys()
	assert.NoError(t, err)

	operatorAddress := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	registrationData := sk.MakePubkeyRegistrationData(operatorAddress)
	assert.NotNil(t, registrationData)
}
