package deploy

import (
	"fmt"
	"log"
)

const (
	dockerBuildContext = "../"

	localinitDockerfile = "inabox-docker/local-init/Dockerfile"

	churnerImage      = "ghcr.io/layr-labs/eigenda/churner:local"
	churnerDockerfile = "churner/Dockerfile"

	disImage      = "ghcr.io/layr-labs/eigenda/disperser:local"
	disDockerfile = "disperser/disperser.Dockerfile"

	encoderImage      = "ghcr.io/layr-labs/eigenda/encoder:local"
	encoderDockerfile = "disperser/encoder.Dockerfile"

	batcherImage      = "ghcr.io/layr-labs/eigenda/batcher:local"
	batcherDockerfile = "disperser/batcher.Dockerfile"

	nodeImage      = "ghcr.io/layr-labs/eigenda/node:local"
	nodeDockerfile = "node/cmd/Dockerfile"

	retrieverImage      = "ghcr.io/layr-labs/eigenda/retriever:local"
	retrieverDockerfile = "retriever/Dockerfile"
)

func (env *Config) GenerateServiceConfig() {
	// Create a new experiment and deploy the contracts
	err := env.loadPrivateKeys()
	if err != nil {
		log.Panicf("could not load private keys: %v", err)
	}

	fmt.Println("Generating service config variables")
	env.GenerateAllVariables()
}
