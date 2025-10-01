package main

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/ejector"
)

func main() {
	// Loads config from environment variables.
	// If there are CLI arguments, they are treated as config file paths and loaded.
	cfg, err := config.ParseConfigFromCLI(ejector.DefaultEjectorConfig, ejector.EjectorConfigEnvPrefix)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg) // TODO
}
