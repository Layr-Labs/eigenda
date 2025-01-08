package clients

import (
	"context"
	aws2 "github.com/Layr-Labs/eigenda/common/aws"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var (
	dockertestPool     *dockertest.Pool
	dockertestResource *dockertest.Resource
)

const (
	localstackPort = "4570"
	localstackHost = "http://0.0.0.0:4570"
	region         = "us-east-1"
)

func setup(t *testing.T) {
	deployLocalStack := !(os.Getenv("DEPLOY_LOCALSTACK") == "false")

	_, b, _, _ := runtime.Caller(0)
	rootPath := filepath.Join(filepath.Dir(b), "../../..")
	changeDirectory(filepath.Join(rootPath, "inabox"))

	if deployLocalStack {
		var err error
		dockertestPool, dockertestResource, err = deploy.StartDockertestWithLocalstackContainer(localstackPort)
		require.NoError(t, err)
	}
}

func changeDirectory(path string) {
	err := os.Chdir(path)
	if err != nil {

		currentDirectory, err := os.Getwd()
		if err != nil {
			log.Printf("Failed to get current directory. Error: %s", err)
		}

		log.Panicf("Failed to change directories. CWD: %s, Error: %s", currentDirectory, err)
	}

	newDir, err := os.Getwd()
	if err != nil {
		log.Panicf("Failed to get working directory. Error: %s", err)
	}
	log.Printf("Current Working Directory: %s\n", newDir)
}

func teardown() {
	deployLocalStack := !(os.Getenv("DEPLOY_LOCALSTACK") == "false")

	if deployLocalStack {
		deploy.PurgeDockertestResources(dockertestPool, dockertestResource)
	}
}

func TestRequestSigning(t *testing.T) {
	rand := random.NewTestRandom(t)
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

			require.Nil(t, request.Signature)
			request.Signature = signature
			err = auth.VerifyStoreChunksRequest(publicAddress, request)
			require.NoError(t, err)

			// Changing a byte in the middle of the signature should make the verification fail
			badSignature := make([]byte, len(signature))
			copy(badSignature, signature)
			badSignature[10] = badSignature[10] + 1
			request.Signature = badSignature
			err = auth.VerifyStoreChunksRequest(publicAddress, request)
			require.Error(t, err)

			// Changing a byte in the middle of the request should make the verification fail
			request.DisperserID = request.DisperserID + 1
			request.Signature = signature
			err = auth.VerifyStoreChunksRequest(publicAddress, request)
			require.Error(t, err)
		}
	}
}
