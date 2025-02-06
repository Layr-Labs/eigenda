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

func parseConfig(configFile string) (*LoadGeneratorConfig, error) {
	configFile, err := client.ResolveTildeInPath(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config file: %v", err)
	}
	configFileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	config := &LoadGeneratorConfig{}
	err = json.Unmarshal(configFileBytes, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %v", err)
	}

	return config, nil
}

func TestLoad(t *testing.T) {
	rand := random.NewTestRandom(t)
	c := client.GetTestClient(t)

	config, err := parseConfig("../config/load/1mb_s-10mb-0x.json")
	require.NoError(t, err)

	generator := NewLoadGenerator(config, c, rand)

	signals := make(chan os.Signal)
	go func() {
		<-signals
		fmt.Printf("Stop Requested\n")
		generator.Stop()
	}()

	generator.Start(true)
}
