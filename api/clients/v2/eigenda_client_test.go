package v2

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/api/clients"
	"github.com/Layr-Labs/eigenda/api/clients/codecs"
	codecsmock "github.com/Layr-Labs/eigenda/api/clients/codecs/mock"
	clientsmock "github.com/Layr-Labs/eigenda/api/clients/mock"
	testrandom "github.com/Layr-Labs/eigenda/common/testutils/random"
	core "github.com/Layr-Labs/eigenda/core/v2"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type ClientTester struct {
	Random          *testrandom.TestRandom
	Client          *EigenDAClient
	MockRelayClient *clientsmock.MockRelayClient
	MockCodec       *codecsmock.BlobCodec
}

func (c *ClientTester) assertExpectations(t *testing.T) {
	c.MockRelayClient.AssertExpectations(t)
	c.MockCodec.AssertExpectations(t)
}

// buildClientTester sets up a client with mocks necessary for testing
func buildClientTester(t *testing.T) ClientTester {
	logger := logging.NewNoopLogger()
	clientConfig := &EigenDAClientConfig{}

	mockRelayClient := clientsmock.MockRelayClient{}
	mockCodec := codecsmock.BlobCodec{}

	random := testrandom.NewTestRandom()

	client, err := NewEigenDAClient(
		logger,
		random.Rand,
		clientConfig,
		&mockRelayClient,
		&mockCodec)

	assert.NotNil(t, client)
	assert.Nil(t, err)

	return ClientTester{
		Random:          random,
		Client:          client,
		MockRelayClient: &mockRelayClient,
		MockCodec:       &mockCodec,
	}
}

// TestGetBlobSuccess tests that a blob is received without error in the happy case
func TestGetBlobSuccess(t *testing.T) {
	tester := buildClientTester(t)

	blobKey := core.BlobKey(tester.Random.RandomBytes(32))
	blobBytes := tester.Random.RandomBytes(100)

	relayKeys := make([]core.RelayKey, 1)
	relayKeys[0] = tester.Random.Uint32()
	blobCert := core.BlobCertificate{
		RelayKeys: relayKeys,
	}

	tester.MockRelayClient.On("GetBlob", mock.Anything, relayKeys[0], blobKey).Return(blobBytes, nil).Once()
	tester.MockCodec.On("DecodeBlob", blobBytes).Return(tester.Random.RandomBytes(50), nil).Once()

	blob, err := tester.Client.GetBlob(context.Background(), blobKey, blobCert)

	assert.NotNil(t, blob)
	assert.Nil(t, err)

	tester.assertExpectations(t)
}

// TestRandomRelayRetries verifies correct behavior when some relays from the certificate do not respond with the blob,
// requiring the client to retry with other relays.
func TestRandomRelayRetries(t *testing.T) {
	tester := buildClientTester(t)

	blobKey := core.BlobKey(tester.Random.RandomBytes(32))
	blobBytes := tester.Random.RandomBytes(100)

	relayCount := 100
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	blobCert := core.BlobCertificate{
		RelayKeys: relayKeys,
	}

	// for this test, only a single relay is online
	// we will be asserting that it takes a different amount of retries to dial this relay, since the array of relay keys to try is randomized
	onlineRelayKey := relayKeys[tester.Random.Intn(len(relayKeys))]

	offlineKeyMatcher := func(relayKey core.RelayKey) bool { return relayKey != onlineRelayKey }
	onlineKeyMatcher := func(relayKey core.RelayKey) bool { return relayKey == onlineRelayKey }
	var failedCallCount int
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.MatchedBy(offlineKeyMatcher), blobKey).Return(nil, fmt.Errorf("offline relay")).Run(func(args mock.Arguments) {
		failedCallCount++
	})
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.MatchedBy(onlineKeyMatcher), blobKey).Return(blobBytes, nil)
	tester.MockCodec.On("DecodeBlob", mock.Anything).Return(tester.Random.RandomBytes(50), nil)

	// keep track of how many tries various blob retrievals require
	// this allows us to assert that there is variability, i.e. that relay call order is actually random
	requiredTries := map[int]bool{}

	for i := 0; i < relayCount; i++ {
		failedCallCount = 0
		blob, err := tester.Client.GetBlob(context.Background(), blobKey, blobCert)
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
	tester := buildClientTester(t)

	blobKey := core.BlobKey(tester.Random.RandomBytes(32))

	relayCount := 10
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	blobCert := core.BlobCertificate{
		RelayKeys: relayKeys,
	}

	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(nil, fmt.Errorf("offline relay"))

	blob, err := tester.Client.GetBlob(context.Background(), blobKey, blobCert)
	assert.Nil(t, blob)
	assert.NotNil(t, err)

	tester.assertExpectations(t)
}

// TestNoRelaysInCert tests that having no relay keys in the cert is handled gracefully
func TestNoRelaysInCert(t *testing.T) {
	tester := buildClientTester(t)

	blobKey := core.BlobKey(tester.Random.RandomBytes(32))

	// cert has no listed relay keys
	blobCert := core.BlobCertificate{
		RelayKeys: []core.RelayKey{},
	}

	blob, err := tester.Client.GetBlob(context.Background(), blobKey, blobCert)
	assert.Nil(t, blob)
	assert.NotNil(t, err)

	tester.assertExpectations(t)
}

// TestGetBlobReturns0Len verifies that a 0 length blob returned from a relay is handled gracefully, and that the client retries after such a failure
func TestGetBlobReturns0Len(t *testing.T) {
	tester := buildClientTester(t)

	blobKey := core.BlobKey(tester.Random.RandomBytes(32))

	relayCount := 10
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	blobCert := core.BlobCertificate{
		RelayKeys: relayKeys,
	}

	// the first GetBlob will return a 0 len blob
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return([]byte{}, nil).Once()
	// the second call will return random bytes
	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(tester.Random.RandomBytes(100), nil).Once()

	tester.MockCodec.On("DecodeBlob", mock.Anything).Return(tester.Random.RandomBytes(50), nil)

	// the call to the first relay will fail with a 0 len blob returned. the call to the second relay will succeed
	blob, err := tester.Client.GetBlob(context.Background(), blobKey, blobCert)
	assert.NotNil(t, blob)
	assert.Nil(t, err)

	tester.assertExpectations(t)
}

// TestFailedDecoding verifies that a failed blob decode is handled gracefully, and that the client retries after such a failure
func TestFailedDecoding(t *testing.T) {
	tester := buildClientTester(t)

	blobKey := core.BlobKey(tester.Random.RandomBytes(32))

	relayCount := 10
	relayKeys := make([]core.RelayKey, relayCount)
	for i := 0; i < relayCount; i++ {
		relayKeys[i] = tester.Random.Uint32()
	}
	blobCert := core.BlobCertificate{
		RelayKeys: relayKeys,
	}

	tester.MockRelayClient.On("GetBlob", mock.Anything, mock.Anything, blobKey).Return(tester.Random.RandomBytes(100), nil)

	tester.MockCodec.On("DecodeBlob", mock.Anything).Return(nil, fmt.Errorf("decode failed")).Once()
	tester.MockCodec.On("DecodeBlob", mock.Anything).Return(tester.Random.RandomBytes(50), nil).Once()

	// decoding will fail the first time, but succeed the second time
	blob, err := tester.Client.GetBlob(context.Background(), blobKey, blobCert)
	assert.NotNil(t, blob)
	assert.Nil(t, err)

	tester.assertExpectations(t)
}

// TestErrorFreeClose tests the happy case, where none of the internal closes yield an error
func TestErrorFreeClose(t *testing.T) {
	tester := buildClientTester(t)

	tester.MockRelayClient.On("Close").Return(nil).Once()

	err := tester.Client.Close()
	assert.Nil(t, err)

	tester.assertExpectations(t)
}

// TestErrorClose tests what happens when subcomponents throw errors when being closed
func TestErrorClose(t *testing.T) {
	tester := buildClientTester(t)

	tester.MockRelayClient.On("Close").Return(fmt.Errorf("close failed")).Once()

	err := tester.Client.Close()
	assert.NotNil(t, err)

	tester.assertExpectations(t)
}

// TestGetCodec checks that the codec used in construction is returned by GetCodec
func TestGetCodec(t *testing.T) {
	tester := buildClientTester(t)

	assert.Equal(t, tester.MockCodec, tester.Client.GetCodec())

	tester.assertExpectations(t)
}

// TestBuilder tests that the method that builds the client from config doesn't throw any obvious errors
func TestBuilder(t *testing.T) {
	clientConfig := &EigenDAClientConfig{
		BlobEncodingVersion:   codecs.DefaultBlobEncoding,
		PointVerificationMode: IFFT,
	}

	sockets := make(map[core.RelayKey]string)
	sockets[core.RelayKey(44)] = "socketVal"

	relayClientConfig := &clients.RelayClientConfig{
		Sockets:           sockets,
		UseSecureGrpcFlag: true,
	}

	client, err := BuildEigenDAClient(
		logging.NewNoopLogger(),
		clientConfig,
		relayClientConfig)

	assert.NotNil(t, client)
	assert.Nil(t, err)

	assert.NotNil(t, client.GetCodec())
}
