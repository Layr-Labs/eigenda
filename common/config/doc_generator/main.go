package main

import (
	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/common/enforce"
	"github.com/Layr-Labs/eigenda/ejector"
	"github.com/Layr-Labs/eigenda/test/v2/load"
)

const configDocsDir = "../../../docs/config"

func main() {
	err := config.DocumentConfig(load.DefaultTrafficGeneratorConfig, configDocsDir, true)
	enforce.NilError(err, "failed to generate docs for the traffic generator config")

	err = config.DocumentConfig(ejector.DefaultEjectorConfig, configDocsDir, true)
	enforce.NilError(err, "failed to generate docs for the ejector config")
}
