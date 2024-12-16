package clients

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
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
	rootPath := filepath.Join(filepath.Dir(b), "../..")
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
		log.Panicf("Failed to change directories. Error: %s", err)
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
	createKeyOutput, err := keyManager.CreateKey(context.Background(), &kms.CreateKeyInput{
		KeySpec:  types.KeySpecEccNistP256,
		KeyUsage: types.KeyUsageTypeSignVerify,
	})
	require.NoError(t, err)

	keyID := *createKeyOutput.KeyMetadata.KeyId

	getPublicKeyOutput, err := keyManager.GetPublicKey(context.Background(), &kms.GetPublicKeyInput{
		KeyId: aws.String(keyID),
	})
	require.NoError(t, err)

	k, err := x509.ParsePKIXPublicKey(getPublicKeyOutput.PublicKey)
	require.NoError(t, err)

	publicKey := k.(*ecdsa.PublicKey)
	publicAddress := crypto.PubkeyToAddress(*publicKey)

	request := auth.RandomStoreChunksRequest(rand)
	request.Signature = nil

	signer := NewRequestSigner(region, localstackHost, keyID)

	signature, err := signer.SignStoreChunksRequest(context.Background(), request)
	require.NoError(t, err)

	require.Nil(t, request.Signature)
	request.Signature = signature

	err = auth.VerifyStoreChunksRequest(publicAddress, request)
	require.NoError(t, err)
}
