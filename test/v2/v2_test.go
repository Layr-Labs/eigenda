package v2

import (
	"context"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"sync"
	"testing"
	"time"
)

var (
	preprodConfig = &TestClientConfig{
		PrivateKeyFile:                "/users/cody/ws/keys/preprod-account.txt",
		DisperserHostname:             "disperser-preprod-holesky.eigenda.xyz",
		DisperserPort:                 443,
		EthRPCURLs:                    []string{"https://ethereum-holesky-rpc.publicnode.com"},
		BLSOperatorStateRetrieverAddr: "0x93545e3b9013CcaBc31E80898fef7569a4024C0C",
		EigenDAServiceManagerAddr:     "0x54A03db2784E3D0aCC08344D05385d0b62d4F432",
		SubgraphURL:                   "https://subgraph.satsuma-prod.com/51caed8fa9cb/eigenlabs/eigenda-operator-state-preprod-holesky/version/v0.7.0/api",
		KZGPath:                       "/Users/cody/ws/srs",
		SRSOrder:                      268435456,
		SRSNumberToLoad:               1, // 2097152 is default in production, no need to load so much for tests
	}

	preprodLock   sync.Mutex
	preprodClient *TestClient
)

// TODO: automatically download KZG points if they are not present

func getPreprodClient(t *testing.T) *TestClient {
	preprodLock.Lock()
	defer preprodLock.Unlock()

	if preprodClient == nil {
		preprodClient = NewTestClient(t, preprodConfig)
	}

	return preprodClient
}

// Tests the basic dispersal workflow:
// - disperse a blob
// - wait for it to be confirmed
// - read the blob from the relays
// - read the blob from the validators
func testBasicDispersal(t *testing.T, rand *random.TestRandom, payloadSize int) {
	client := getPreprodClient(t)

	data := rand.Bytes(payloadSize)

	quorums := make([]core.QuorumID, 2)
	quorums[0] = core.QuorumID(0)
	quorums[1] = core.QuorumID(1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client.DisperseAndVerify(ctx, data, quorums)
}

// Disperse a 0 byte payload.
// TODO this is expected to fail
//func TestEmptyBlobDispersal(t *testing.T) {
//	rand := random.NewTestRandom(t)
//	testBasicDispersal(t, rand, 0)
//}

// Disperse a 1 byte payload.
func TestMicroscopicBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	testBasicDispersal(t, rand, 1)
}

// Disperse a small payload (between 1KB and 2KB).
func TestSmallBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	dataLength := 1024 + rand.Intn(1024)
	testBasicDispersal(t, rand, dataLength)
}

// Disperse a medium payload (between 100KB and 200KB).
func TestMediumBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	dataLength := 1024 * (100 + rand.Intn(100))
	testBasicDispersal(t, rand, dataLength)
}

// Disperse a medium payload (between 1MB and 16MB).
func TestLargeBlobDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	dataLength := 1024 * 1024 * (1 + rand.Intn(16))
	testBasicDispersal(t, rand, dataLength)
}

// TODO size 0 blob
// TODO maximum size blob
