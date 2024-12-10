package auth

import (
	"crypto/ecdsa"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	"github.com/aws/smithy-go/rand"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadingKeysFromFile(t *testing.T) {
	tu.InitializeRandom()

	publicKey, err := ReadPublicECDSAKeyFile("./test-public.pem")
	require.NoError(t, err)
	require.NotNil(t, publicKey)

	privateKey, err := ReadPrivateECDSAKeyFile("./test-private.pem")
	require.NoError(t, err)
	require.NotNil(t, privateKey)

	bytesToSign := tu.RandomBytes(32)

	signature, err := ecdsa.SignASN1(rand.Reader, privateKey, bytesToSign)
	require.NoError(t, err)

	isValid := ecdsa.VerifyASN1(publicKey, bytesToSign, signature)
	require.True(t, isValid)

	// Change some bytes in the signature, it should be invalid now
	signature2 := make([]byte, len(signature))
	copy(signature2, signature)
	signature2[0] = signature2[0] + 1
	isValid = ecdsa.VerifyASN1(publicKey, bytesToSign, signature2)
	require.False(t, isValid)

	// Change some bytes in the message, it should be invalid now
	bytesToSign2 := make([]byte, len(bytesToSign))
	copy(bytesToSign2, bytesToSign)
	bytesToSign2[0] = bytesToSign2[0] + 1
	isValid = ecdsa.VerifyASN1(publicKey, bytesToSign2, signature)
	require.False(t, isValid)

}
