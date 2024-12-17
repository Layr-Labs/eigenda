package clients

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509/pkix"
	"encoding/asn1"
	"errors"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/inabox/deploy"
	"github.com/Layr-Labs/eigenda/node/auth"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/cryptobyte"
	"log"
	"math/big"
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

type publicKeyInfo struct {
	Raw       asn1.RawContent
	Algorithm pkix.AlgorithmIdentifier
	PublicKey asn1.BitString
}

func ParsePublicKey(keyBytes []byte) (*ecdsa.PublicKey, error) {
	pki := publicKeyInfo{}
	rest, err := asn1.Unmarshal(keyBytes, &pki)

	if err != nil {
		return nil, err
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("trailing data after public key (%d bytes)", len(rest))
	}

	rightAlignedKey := cryptobyte.String(pki.PublicKey.RightAlign())

	x, y := elliptic.Unmarshal(crypto.S256(), rightAlignedKey)
	if x == nil {
		return nil, errors.New("x509: failed to unmarshal elliptic curve point")
	}

	return &ecdsa.PublicKey{
		Curve: crypto.S256(),
		X:     x,
		Y:     y,
	}, nil
}

type signatureInfo struct {
	R *big.Int
	S *big.Int
}

// AddRecoveryID computes the recovery ID for a given signature and public key and adds it to the signature.
func AddRecoveryID(hash []byte, pubKey *ecdsa.PublicKey, partialSignature []byte) error {
	for v := 0; v < 4; v++ {
		partialSignature[64] = byte(v)
		recoveredPubKey, err := secp256k1.RecoverPubkey(hash, partialSignature)
		if err != nil {
			return fmt.Errorf("failed to recover public key: %w", err)
		}

		x, y := elliptic.Unmarshal(secp256k1.S256(), recoveredPubKey)
		if x.Cmp(pubKey.X) == 0 && y.Cmp(pubKey.Y) == 0 {
			return nil
		}
	}

	return fmt.Errorf("no valid recovery ID found")
}

// pad32 pads a byte slice to 32 bytes, inserting zeros at the beginning if necessary.
func pad32(bytes []byte) []byte {
	if len(bytes) == 32 {
		return bytes
	}

	padded := make([]byte, 32)
	copy(padded[32-len(bytes):], bytes)
	return padded
}

// ParseSignature parses a signature from AWS into eth format
func ParseSignature(
	publicKey *ecdsa.PublicKey,
	hash []byte,
	signatureBytes []byte) ([]byte, error) {

	si := signatureInfo{}
	rest, err := asn1.Unmarshal(signatureBytes, &si)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal signature: %w", err)
	}
	if len(rest) > 0 {
		return nil, fmt.Errorf("trailing data after signature (%d bytes)", len(rest))
	}

	rBytes := pad32(si.R.Bytes())
	sBytes := pad32(si.S.Bytes())

	result := make([]byte, 65)
	copy(result[0:32], rBytes)
	copy(result[32:64], sBytes)

	err = AddRecoveryID(hash, publicKey, result)
	if err != nil {
		return nil, fmt.Errorf("failed to compute recovery ID: %w", err)
	}

	return result, nil
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

		getPublicKeyOutput, err := keyManager.GetPublicKey(context.Background(), &kms.GetPublicKeyInput{
			KeyId: aws.String(keyID),
		})
		require.NoError(t, err)

		key, err := ParsePublicKey(getPublicKeyOutput.PublicKey)
		require.NoError(t, err)

		publicAddress := crypto.PubkeyToAddress(*key)

		for j := 0; j < 10; j++ {
			request := auth.RandomStoreChunksRequest(rand)
			request.Signature = nil

			signer := NewRequestSigner(region, localstackHost, keyID)

			// Test a valid signature.
			signature, err := signer.SignStoreChunksRequest(context.Background(), request)
			require.NoError(t, err)

			// TODO this doesn't belong here
			signature, err = ParseSignature(
				key,
				auth.HashStoreChunksRequest(request),
				signature)

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
