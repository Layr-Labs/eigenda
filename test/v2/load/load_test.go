package load

import (
	"os"
	"testing"

	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/docker/go-units"
)

func TestLoad(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := client.GetTestClient(t, []core.QuorumID{0, 1})

	config := DefaultLoadGeneratorConfig()
	config.AverageBlobSize = 100 * units.KiB
	config.BlobSizeStdDev = 50 * units.KiB
	config.BytesPerSecond = 100 * units.KiB

	generator := NewLoadGenerator(config, c, rand)

	signals := make(chan os.Signal)
	go func() {
		<-signals
		generator.Stop()
	}()

	generator.Start(true)
}
