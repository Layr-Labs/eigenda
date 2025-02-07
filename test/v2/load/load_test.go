package load

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/test/v2/client"
)

const targetConfigFile = "../config/load/100kb_s-1mb-3x.json"

func TestLoad(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := client.GetTestClient(t, []core.QuorumID{0, 1})

	config, err := ReadConfigFile(targetConfigFile)
	require.NoError(t, err)

	generator := NewLoadGenerator(config, c, rand)

	signals := make(chan os.Signal)
	go func() {
		<-signals
		generator.Stop()
	}()

	generator.Start(true)
}
