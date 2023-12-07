package genenv

const (
	dockerBuildContext = "../"

	ethInitDockerfile   = "inabox/eth-init-script/cmd/Dockerfile"
	graphInitDockerfile = "inabox/graph-init-script/cmd/Dockerfile"
	awsInitDockerfile   = "inabox/aws-init-script/cmd/Dockerfile"

	churnerImage      = "ghcr.io/layr-labs/eigenda/churner:local"
	churnerDockerfile = "churner/cmd/Dockerfile"

	disImage      = "ghcr.io/layr-labs/eigenda/disperser:local"
	disDockerfile = "disperser/cmd/apiserver/Dockerfile"

	encoderImage      = "ghcr.io/layr-labs/eigenda/encoder:local"
	encoderDockerfile = "disperser/cmd/encoder/Dockerfile"

	batcherImage      = "ghcr.io/layr-labs/eigenda/batcher:local"
	batcherDockerfile = "disperser/cmd/batcher/Dockerfile"

	nodeImage      = "ghcr.io/layr-labs/eigenda/node:local"
	nodeDockerfile = "node/cmd/Dockerfile"

	retrieverImage      = "ghcr.io/layr-labs/eigenda/retriever:local"
	retrieverDockerfile = "retriever/cmd/Dockerfile"
)
