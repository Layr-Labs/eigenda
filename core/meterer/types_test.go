package meterer_test

import (
	"math/big"
	"testing"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/core/meterer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEIP712Signer(t *testing.T) {
	chainID := big.NewInt(17000)
	verifyingContract := common.HexToAddress("0x1234000000000000000000000000000000000000")
	signer := meterer.NewEIP712Signer(chainID, verifyingContract)

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	commitment := core.NewG1Point(big.NewInt(123), big.NewInt(456))

	header := &meterer.BlobHeader{
		BinIndex:          0,
		CumulativePayment: 1000,
		Commitment:        *commitment,
		DataLength:        1024,
		QuorumNumbers:     []uint8{1},
	}

	t.Run("SignBlobHeader", func(t *testing.T) {
		signature, err := signer.SignBlobHeader(header, privateKey)
		require.NoError(t, err)
		assert.NotEmpty(t, signature)
	})

	t.Run("RecoverSender", func(t *testing.T) {
		signature, err := signer.SignBlobHeader(header, privateKey)
		require.NoError(t, err)

		header.Signature = signature
		recoveredAddress, err := signer.RecoverSender(header)
		require.NoError(t, err)

		expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
		assert.Equal(t, expectedAddress, recoveredAddress)
	})
}

func TestConstructBlobHeader(t *testing.T) {
	chainID := big.NewInt(17000)
	verifyingContract := common.HexToAddress("0x1234000000000000000000000000000000000000")
	signer := meterer.NewEIP712Signer(chainID, verifyingContract)

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	commitment := core.NewG1Point(big.NewInt(123), big.NewInt(456))

	header, err := meterer.ConstructBlobHeader(
		signer,
		0,           // binIndex
		1000,        // cumulativePayment
		*commitment, // core.G1Point
		1024,        // dataLength
		[]uint8{1},
		privateKey,
	)

	require.NoError(t, err)
	assert.NotNil(t, header)
	assert.NotEmpty(t, header.Signature)

	recoveredAddress, err := signer.RecoverSender(header)
	require.NoError(t, err)

	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	assert.Equal(t, expectedAddress, recoveredAddress)
}

func TestEIP712SignerWithDifferentKeys(t *testing.T) {
	chainID := big.NewInt(17000)
	verifyingContract := common.HexToAddress("0x1234000000000000000000000000000000000000")
	signer := meterer.NewEIP712Signer(chainID, verifyingContract)

	privateKey1, err := crypto.GenerateKey()
	require.NoError(t, err)

	privateKey2, err := crypto.GenerateKey()
	require.NoError(t, err)

	commitment := core.NewG1Point(big.NewInt(123), big.NewInt(456))

	header, err := meterer.ConstructBlobHeader(
		signer,
		0,
		1000,
		*commitment,
		1024,
		[]uint8{1},
		privateKey1,
	)

	require.NoError(t, err)
	assert.NotNil(t, header)
	assert.NotEmpty(t, header.Signature)

	recoveredAddress, err := signer.RecoverSender(header)
	require.NoError(t, err)

	expectedAddress1 := crypto.PubkeyToAddress(privateKey1.PublicKey)
	expectedAddress2 := crypto.PubkeyToAddress(privateKey2.PublicKey)

	assert.Equal(t, expectedAddress1, recoveredAddress)
	assert.NotEqual(t, expectedAddress2, recoveredAddress)
}

func TestEIP712SignerWithModifiedHeader(t *testing.T) {
	chainID := big.NewInt(17000)
	verifyingContract := common.HexToAddress("0x1234000000000000000000000000000000000000")
	signer := meterer.NewEIP712Signer(chainID, verifyingContract)

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	commitment := core.NewG1Point(big.NewInt(123), big.NewInt(456))

	header, err := meterer.ConstructBlobHeader(
		signer,
		0,
		1000,
		*commitment,
		1024,
		[]uint8{1},
		privateKey,
	)

	require.NoError(t, err)
	assert.NotNil(t, header)
	assert.NotEmpty(t, header.Signature)
	recoveredAddress, err := signer.RecoverSender(header)
	require.NoError(t, err)

	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	assert.Equal(t, expectedAddress, recoveredAddress, "Recovered address should match the address derived from the private key")

	header.AccountID = "modifiedAccount"

	addr, err := signer.RecoverSender(header)
	require.NotEqual(t, expectedAddress, addr)
}

func TestEIP712SignerWithDifferentChainID(t *testing.T) {
	chainID1 := big.NewInt(17000)
	chainID2 := big.NewInt(17001)
	verifyingContract := common.HexToAddress("0x1234000000000000000000000000000000000000")
	signer1 := meterer.NewEIP712Signer(chainID1, verifyingContract)
	signer2 := meterer.NewEIP712Signer(chainID2, verifyingContract)

	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	commitment := core.NewG1Point(big.NewInt(123), big.NewInt(456))

	header, err := meterer.ConstructBlobHeader(
		signer1,
		0,
		1000,
		*commitment,
		1024,
		[]uint8{1},
		privateKey,
	)

	require.NoError(t, err)
	assert.NotNil(t, header)
	assert.NotEmpty(t, header.Signature)

	// Try to recover the sender using a signer with a different chain ID
	sender, err := signer2.RecoverSender(header)
	expectedAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	require.NotEqual(t, expectedAddress, sender)
}
