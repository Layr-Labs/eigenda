package load

import (
	"encoding/json"
	"fmt"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func parseConfig(t *testing.T, configFile string) *LoadGeneratorConfig {
	configFile = client.ResolveTildeInPath(t, configFile)
	configFileBytes, err := os.ReadFile(configFile)
	require.NoError(t, err)

	config := &LoadGeneratorConfig{}
	err = json.Unmarshal(configFileBytes, config)
	require.NoError(t, err)

	return config
}

func TestLoad(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := client.GetClient(t)

	config := parseConfig(t, "../config/load/1mb_s-10mb-0x.json")

	generator := NewLoadGenerator(config, c, rand)

	signals := make(chan os.Signal)
	go func() {
		<-signals
		fmt.Printf("Stop Requested\n")
		generator.Stop()
	}()

	generator.Start(true)
}
