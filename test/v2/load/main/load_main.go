package main

import (
	"fmt"
	"os"

	"github.com/Layr-Labs/eigenda/test/v2/client"
	"github.com/Layr-Labs/eigenda/test/v2/load"
	"github.com/stretchr/testify/require"
)

func main() {
	if len(os.Args) != 3 {
		panic(fmt.Sprintf("Expected 3 args, got %d. Usage: %s <env_file> <load_file>.\n"+
			"If '-' is passed in lieu of a config file, the config file path is read from the environment variable "+
			"$GENERATOR_ENV or $GENERATOR_LOAD, respectively.\n",
			len(os.Args), os.Args[0]))
	}

	envFile := os.Args[1]
	if envFile == "-" {
		envFile = os.Getenv("GENERATOR_ENV")
		if envFile == "" {
			panic("$GENERATOR_ENV not set")
		}
	}

	loadFile := os.Args[2]
	if loadFile == "-" {
		loadFile = os.Getenv("GENERATOR_LOAD")
		if loadFile == "" {
			panic("$GENERATOR_LOAD not set")
		}
	}

	c, err := client.GetClient(envFile)
	if err != nil {
		panic(err)
	}

	config, err := load.ReadConfigFile(loadFile)
	require.NoError(nil, err)

	generator := load.NewLoadGenerator(config, c)

	signals := make(chan os.Signal)
	go func() {
		<-signals
		generator.Stop()
	}()

	generator.Start(true)
}
