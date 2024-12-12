package auth

import (
	"context"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	wmock "github.com/Layr-Labs/eigenda/core/mock"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

// TODO:
//  - test good request
//  - test bad request
//  - test request from disperser that doesn't exist
//  - test auth caching
//  - test key caching
//  - test cache sizes

//// Verify that public key can be converted to an eth address and back
//func TestPubKeyRoundTrip(t *testing.T) {
//	rand := random.NewTestRandom(t)
//
//	publicKey, _ := rand.ECDSA()
//
//	ethAddress := crypto.PubkeyToAddress(*publicKey)
//
//	publicKey2, err := crypto.UnmarshalPubkey(ethAddress.Bytes())
//	require.NoError(t, err)
//
//	require.Equal(t, publicKey, publicKey2)
//}

func TestValidRequest(t *testing.T) {
	rand := random.NewTestRandom(t)

	start := rand.Time()

	publicKey, privateKey := rand.ECDSA()
	disperserAddress := crypto.PubkeyToAddress(*publicKey)

	//disperserAddressBytes := crypto.FromECDSAPub(publicKey)
	//disperserAddress := gethcommon.BytesToAddress(disperserAddressBytes)

	chainReader := wmock.MockWriter{}
	chainReader.Mock.On("GetDisperserAddress", uint32(0)).Return(disperserAddress, nil)

	authenticator, err := NewRequestAuthenticator(
		context.Background(),
		&chainReader,
		10,
		time.Minute,
		time.Minute,
		start)
	require.NoError(t, err)

	request := randomStoreChunksRequest(rand)
	request.DisperserID = 0
	signature, err := SignStoreChunksRequest(privateKey, request)
	require.NoError(t, err)
	request.Signature = signature

	err = authenticator.AuthenticateStoreChunksRequest(context.Background(), "localhost", request, start)
	require.NoError(t, err)
}
