package main

import (
	"context"
	"fmt"
	"github.com/Layr-Labs/eigenda/common"
	"github.com/Layr-Labs/eigenda/relay"
)

// main is the entrypoint for the light node.
func main() {

	config, err := relay.LoadConfigWithViper()
	if err != nil {
		panic(fmt.Sprintf("fatal error loading config: %s", err))
	}

	logger, err := common.NewLogger(config.Log)
	if err != nil {
		panic(fmt.Sprintf("fatal error creating logger: %s", err))
	}
	logger.Info(fmt.Sprintf("Relay configuration: %#v", config))

	// TODO
	_, err = relay.NewServer(
		context.Background(),
		logger,
		config,
		nil,
		nil,
		nil)
	if err != nil {
		panic(fmt.Sprintf("fatal error creating relay server: %s", err))
	}
}
