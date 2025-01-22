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
		SRSNumberToLoad:               2097152,
	}

	preprodLock   sync.Mutex
	preprodClient *TestClient
)

func getPreprodClient(t *testing.T) *TestClient {
	preprodLock.Lock()
	defer preprodLock.Unlock()

	if preprodClient == nil {
		preprodClient = NewTestClient(t, preprodConfig)
	}

	return preprodClient
}

func TestSimpleDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)
	client := getPreprodClient(t)

	dataLength := 1024 + rand.Intn(1024)
	data := rand.Bytes(dataLength)

	quorums := make([]core.QuorumID, 2)
	quorums[0] = core.QuorumID(0)
	quorums[1] = core.QuorumID(1)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	client.DisperseAndVerify(ctx, data, quorums)
}
