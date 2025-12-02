package main

import (
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/disperser/encoder"
	"github.com/Layr-Labs/eigenda/ejector"
	"github.com/Layr-Labs/eigenda/test/v2/load"
)

const configDocsDir = "../../../docs/config"

// This program generates markdown documentation for configuration structs.
func main() {
	err := config.DocumentConfig(load.DefaultTrafficGeneratorConfig, configDocsDir, true)
	enforce.NilError(err, "failed to generate docs for the traffic generator config")

	err = config.DocumentConfig(ejector.DefaultRootEjectorConfig, configDocsDir, true)
	enforce.NilError(err, "failed to generate docs for the ejector config")

	err = config.DocumentConfig(encoder.DefaultEncoderConfig, configDocsDir, true)
	enforce.NilError(err, "failed to generate docs for the encoder config")
}
