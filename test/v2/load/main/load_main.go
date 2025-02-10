package main

import (
	"fmt"
	"github.com/Layr-Labs/eigenda/common/testutils/random"
	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/Layr-Labs/eigenda/test/v2/load"
	"github.com/stretchr/testify/require"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		panic(fmt.Sprintf("Expected 3 args, got %d. Usage: %s <env_file> <load_file>\n",
			len(os.Args), os.Args[0]))
	}

	envFile := os.Args[1]
	loadFile := os.Args[2]

	client.SetTargetConfigFile(envFile)
	c, err := client.GetClient()
	if err != nil {
		panic(err)
	}

	rand := random.NewTestRandom(nil)

	config, err := load.ReadConfigFile(loadFile)
	require.NoError(nil, err)

	generator := load.NewLoadGenerator(config, c, rand)

	signals := make(chan os.Signal)
	go func() {
		<-signals
		generator.Stop()
	}()

	generator.Start(true)
}
