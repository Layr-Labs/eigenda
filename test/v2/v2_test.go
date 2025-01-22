package v2

import (
	"context"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"testing"
	"time"
)

func TestSimpleDispersal(t *testing.T) {
	rand := random.NewTestRandom(t)

	client := NewTestClient(t)

	dataLength := 1024 + rand.Intn(1024)
	data := rand.Bytes(dataLength)

	quorums := make([]core.QuorumID, 2)
	quorums[0] = core.QuorumID(0)
	quorums[1] = core.QuorumID(1)

	client.DispersePayload(context.Background(), 10*time.Minute, data, quorums)

}
