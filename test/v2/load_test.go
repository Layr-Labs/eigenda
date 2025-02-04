package v2

import (
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/docker/go-units"
	"testing"
)

func TestLightLoad(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := getClient(t)

	config := DefaultLoadGeneratorConfig()
	config.AverageBlobSize = 100 * units.KiB
	config.BlobSizeStdDev = 50 * units.KiB
	config.BytesPerSecond = 100 * units.KiB

	generator := NewLoadGenerator(config, c, rand)
	generator.Start(true)
}
