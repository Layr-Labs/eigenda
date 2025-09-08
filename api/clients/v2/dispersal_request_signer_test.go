package clients

import (
	"context"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/api/hashing"
	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/Layr-Labs/eigenda/testbed"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

var (
	localstackContainer *testbed.LocalStackContainer
)

const (
	localstackPort = "4579"
	localstackHost = "http://0.0.0.0:4579"
	region         = "us-east-1"
)

func setup(t *testing.T) {
	deployLocalStack := (os.Getenv("DEPLOY_LOCALSTACK") != "false")

	if deployLocalStack {
		var err error
		cfg := testbed.DefaultLocalStackConfig()
		cfg.Services = []string{"s3, kms"}
		cfg.Port = localstackPort
		cfg.Host = "0.0.0.0"

		localstackContainer, err = testbed.NewLocalStackContainer(context.Background(), cfg)
		require.NoError(t, err)
	}
}

func teardown() {
	deployLocalStack := (os.Getenv("DEPLOY_LOCALSTACK") != "false")

	if deployLocalStack {
		_ = localstackContainer.Terminate(context.Background())
	}
}

func TestRequestSigning(t *testing.T) {
	rand := random.NewTestRandom()
	setup(t)
	defer teardown()

	keyManager := kms.New(kms.Options{
		Region:       region,
		BaseEndpoint: aws.String(localstackHost),
	})

	for i := 0; i < 10; i++ {
		createKeyOutput, err := keyManager.CreateKey(context.Background(), &kms.CreateKeyInput{
			KeySpec:  types.KeySpecEccSecgP256k1,
			KeyUsage: types.KeyUsageTypeSignVerify,
		})
		require.NoError(t, err)

		keyID := *createKeyOutput.KeyMetadata.KeyId

		key, err := aws2.LoadPublicKeyKMS(context.Background(), keyManager, keyID)
		require.NoError(t, err)

		publicAddress := crypto.PubkeyToAddress(*key)

		for j := 0; j < 10; j++ {
			request := auth.RandomStoreChunksRequest(rand)
			request.Signature = nil

			signer, err := NewDispersalRequestSigner(context.Background(), region, localstackHost, keyID)
			require.NoError(t, err)

			// Test a valid signature.
			signature, err := signer.SignStoreChunksRequest(context.Background(), request)
			require.NoError(t, err)

			require.Nil(t, request.GetSignature())
			request.Signature = signature
			hash, err := auth.VerifyStoreChunksRequest(publicAddress, request)
			require.NoError(t, err)
			expectedHash, err := hashing.HashStoreChunksRequest(request)
			require.NoError(t, err)
			require.Equal(t, expectedHash, hash)

			// Changing a byte in the middle of the signature should make the verification fail
			badSignature := make([]byte, len(signature))
			copy(badSignature, signature)
			badSignature[10] = badSignature[10] + 1
			request.Signature = badSignature
			hash, err = auth.VerifyStoreChunksRequest(publicAddress, request)
			require.Error(t, err)
			require.Nil(t, hash)

			// Changing a byte in the middle of the request should make the verification fail
			request.DisperserID = request.GetDisperserID() + 1
			request.Signature = signature
			hash, err = auth.VerifyStoreChunksRequest(publicAddress, request)
			require.Error(t, err)
			require.Nil(t, hash)
		}
	}
}
