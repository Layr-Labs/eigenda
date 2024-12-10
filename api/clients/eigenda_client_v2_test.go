package clients_test

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	codecsmock "github.com/Layr-Labs/eigenda/api/clients/codecs/mock"
	clientsmock "github.com/Layr-Labs/eigenda/api/clients/mock"
	tu "github.com/Layr-Labs/eigenda/common/testutils"
	v2 "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"testing"
	"time"
)

type ClientV2Tester struct {
	ClientV2        *clients.EigenDAClientV2
	MockRelayClient *clientsmock.MockRelayClient
	MockCodec       *codecsmock.BlobCodec
}

func (c *ClientV2Tester) assertExpectations(t *testing.T) {
	c.MockRelayClient.AssertExpectations(t)
	c.MockCodec.AssertExpectations(t)
}

// buildClientV2Tester sets up a V2 client, with mocks necessary for testing
func buildClientV2Tester(t *testing.T) ClientV2Tester {
	tu.InitializeRandom()
	logger := logging.NewNoopLogger()
	clientConfig := &clients.EigenDAClientConfig{}

	mockRelayClient := clientsmock.MockRelayClient{}
	mockCodec := codecsmock.BlobCodec{}

	// TODO (litt3): use TestRandom once the PR merges https://github.com/Layr-Labs/eigenda/pull/976
	random := rand.New(rand.NewSource(rand.Int63()))

	client, err := clients.NewEigenDAClientV2(
		logger,
		random,
		clientConfig,
		&mockRelayClient,
		&mockCodec)

	assert.NotNil(t, client)
	assert.Nil(t, err)

	return ClientV2Tester{
		ClientV2:        client,
		MockRelayClient: &mockRelayClient,
		MockCodec:       &mockCodec,
	}
}

// TestGetBlobSuccess tests that a blob is received without error in the happy case
func TestGetBlobSuccess(t *testing.T) {
	tester := buildClientV2Tester(t)

	blobKey := v2.BlobKey(tu.RandomBytes(32))
	blobBytes := tu.RandomBytes(100)

	relayKeys := make([]v2.RelayKey, 1)
	relayKeys[0] = rand.Uint32()
	blobCert := v2.BlobCertificate{
		RelayKeys: relayKeys,
	}

	tester.MockRelayClient.On("GetBlob", mock.Anything, relayKeys[0], blobKey).Return(blobBytes, nil).Once()
	tester.MockCodec.On("DecodeBlob", blobBytes).Return(tu.RandomBytes(50), nil).Once()

	blob, err := tester.ClientV2.GetBlob(context.Background(), blobKey, blobCert)

	assert.NotNil(t, blob)
	assert.Nil(t, err)

	tester.assertExpectations(t)
}

// TestRandomRelayRetries verifies correct behavior when some relays from the certificate do not respond with the blob,
// requiring the client to retry with other relays.
func TestRandomRelayRetries(t *testing.T) {
	tester := buildClientV2Tester(t)

	blobKey := v2.BlobKey(tu.RandomBytes(32))
	blobBytes := tu.RandomBytes(100)

	relayCount := 100
	relayKeys := make([]v2.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = rand.Uint32()
	}
	blobCert := v2.BlobCertificate{
		RelayKeys: relayKeys,
	}

	// for this test, only a single relay is online
	// we will be asserting that it takes a different amount of retries to dial this relay, since the array of relay keys to try is randomized
	onlineRelayKey := relayKeys[rand.Intn(len(relayKeys))]

	offlineKeyMatcher := func(relayKey v2.RelayKey) bool { return relayKey != onlineRelayKey }
	onlineKeyMatcher := func(relayKey v2.RelayKey) bool { return relayKey == onlineRelayKey }
	var failedCallCount int
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.MatchedBy(offlineKeyMatcher), blobKey).Return(nil, fmt.Errorf("offline relay")).Run(func(args mock.Arguments) {
		failedCallCount++
	})
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.MatchedBy(onlineKeyMatcher), blobKey).Return(blobBytes, nil)
	tester.MockCodec.On("DecodeBlob", mock.Anything).Return(tu.RandomBytes(50), nil)

	// keep track of how many tries various blob retrievals require
	// this allows us to assert that there is variability, i.e. that relay call order is actually random
	requiredTries := map[int]bool{}

	for i := 0; i < relayCount; i++ {
		failedCallCount = 0
		blob, err := tester.ClientV2.GetBlob(context.Background(), blobKey, blobCert)
		assert.NotNil(t, blob)
		assert.Nil(t, err)

		requiredTries[failedCallCount] = true
	}

	// with 100 random tries, with possible values between 1 and 100, we can very confidently assert that there are at least 10 unique values
	assert.Greater(t, len(requiredTries), 10)

	tester.assertExpectations(t)
}

// TestNoRelayResponse tests functionality when none of the relays listed in the blob certificate respond
func TestNoRelayResponse(t *testing.T) {
	tester := buildClientV2Tester(t)

	blobKey := v2.BlobKey(tu.RandomBytes(32))

	relayCount := 10
	relayKeys := make([]v2.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = rand.Uint32()
	}
	blobCert := v2.BlobCertificate{
		RelayKeys: relayKeys,
	}

	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(nil, fmt.Errorf("offline relay"))

	blob, err := tester.ClientV2.GetBlob(context.Background(), blobKey, blobCert)
	assert.Nil(t, blob)
	assert.NotNil(t, err)

	tester.assertExpectations(t)
}

// TestNoRelaysInCert tests that having no relay keys in the cert is handled gracefully
func TestNoRelaysInCert(t *testing.T) {
	tester := buildClientV2Tester(t)

	blobKey := v2.BlobKey(tu.RandomBytes(32))

	// cert has no listed relay keys
	blobCert := v2.BlobCertificate{
		RelayKeys: []v2.RelayKey{},
	}

	blob, err := tester.ClientV2.GetBlob(context.Background(), blobKey, blobCert)
	assert.Nil(t, blob)
	assert.NotNil(t, err)

	tester.assertExpectations(t)
}

// TestGetBlobReturns0Len verifies that a 0 length blob returned from a relay is handled gracefully, and that the client retries after such a failure
func TestGetBlobReturns0Len(t *testing.T) {
	tester := buildClientV2Tester(t)

	blobKey := v2.BlobKey(tu.RandomBytes(32))

	relayCount := 10
	relayKeys := make([]v2.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = rand.Uint32()
	}
	blobCert := v2.BlobCertificate{
		RelayKeys: relayKeys,
	}

	// the first GetBlob will return a 0 len blob
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return([]byte{}, nil).Once()
	// the second call will return random bytes
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(tu.RandomBytes(100), nil).Once()

	tester.MockCodec.On("DecodeBlob", mock.Anything).Return(tu.RandomBytes(50), nil)

	// the call to the first relay will fail with a 0 len blob returned. the call to the second relay will succeed
	blob, err := tester.ClientV2.GetBlob(context.Background(), blobKey, blobCert)
	assert.NotNil(t, blob)
	assert.Nil(t, err)

	tester.assertExpectations(t)
}

// TestFailedDecoding verifies that a failed blob decode is handled gracefully, and that the client retries after such a failure
func TestFailedDecoding(t *testing.T) {
	tester := buildClientV2Tester(t)

	blobKey := v2.BlobKey(tu.RandomBytes(32))

	relayCount := 10
	relayKeys := make([]v2.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = rand.Uint32()
	}
	blobCert := v2.BlobCertificate{
		RelayKeys: relayKeys,
	}

	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(tu.RandomBytes(100), nil)

	tester.MockCodec.On("DecodeBlob", mock.Anything).Return(nil, fmt.Errorf("decode failed")).Once()
	tester.MockCodec.On("DecodeBlob", mock.Anything).Return(tu.RandomBytes(50), nil).Once()

	// decoding will fail the first time, but succeed the second time
	blob, err := tester.ClientV2.GetBlob(context.Background(), blobKey, blobCert)
	assert.NotNil(t, blob)
	assert.Nil(t, err)

	tester.assertExpectations(t)
}

// TestErrorFreeClose tests the happy case, where none of the internal closes yield an error
func TestErrorFreeClose(t *testing.T) {
	tester := buildClientV2Tester(t)

	tester.MockRelayClient.On("Close").Return(nil).Once()

	err := tester.ClientV2.Close()
	assert.Nil(t, err)

	tester.assertExpectations(t)
}

// TestErrorClose tests what happens when subcomponents throw errors when being closed
func TestErrorClose(t *testing.T) {
	tester := buildClientV2Tester(t)

	tester.MockRelayClient.On("Close").Return(fmt.Errorf("close failed")).Once()

	err := tester.ClientV2.Close()
	assert.NotNil(t, err)

	tester.assertExpectations(t)
}

// TestGetCodec checks that the codec used in construction is returned by GetCodec
func TestGetCodec(t *testing.T) {
	tester := buildClientV2Tester(t)

	assert.Equal(t, tester.MockCodec, tester.ClientV2.GetCodec())

	tester.assertExpectations(t)
}

// TestBuilder tests that the method that builds the client from config doesn't throw any obvious errors
func TestBuilder(t *testing.T) {
	clientConfig := &clients.EigenDAClientConfig{
		StatusQueryTimeout:           10 * time.Minute,
		StatusQueryRetryInterval:     50 * time.Millisecond,
		ResponseTimeout:              10 * time.Second,
		ConfirmationTimeout:          5 * time.Second,
		CustomQuorumIDs:              []uint{},
		SignerPrivateKeyHex:          "75f9e29cac7f5774d106adb355ef294987ce39b7863b75bb3f2ea42ca160926d",
		DisableTLS:                   false,
		PutBlobEncodingVersion:       codecs.DefaultBlobEncoding,
		DisablePointVerificationMode: false,
		WaitForFinalization:          true,
		RPC:                          "http://localhost:8080",
		EthRpcUrl:                    "http://localhost:8545",
		SvcManagerAddr:               "0x1234567890123456789012345678901234567890",
	}

	relayClientConfig := &clients.RelayClientConfig{
		Sockets:           make(map[v2.RelayKey]string),
		UseSecureGrpcFlag: true,
	}

	clientV2, err := clients.BuildEigenDAClientV2(
		logging.NewNoopLogger(),
		clientConfig,
		relayClientConfig)

	assert.NotNil(t, clientV2)
	assert.Nil(t, err)

	assert.NotNil(t, clientV2.GetCodec())
}
