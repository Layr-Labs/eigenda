package genenv

const (
	dockerBuildContext = "../../../"

	ethInitDockerfile   = "inabox/strategies/containers/eth-init-script/Dockerfile"
	graphInitDockerfile = "inabox/strategies/containers/graph-init-script/Dockerfile"
	awsInitDockerfile   = "inabox/strategies/containers/aws-init-script/Dockerfile"

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
