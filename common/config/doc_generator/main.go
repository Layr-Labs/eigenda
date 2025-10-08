package main

import (
	"fmt"

	"github.com/Layr-Labs/eigenda/common/config"
	"github.com/Layr-Labs/eigenda/test/v2/load"
)

const configDocsDir = "docs/config"

func main() {
	fmt.Printf("Generating config docs in %q...\n", configDocsDir) // TODO
	config.DocumentConfig(load.DefaultTrafficGeneratorConfig, configDocsDir, true)
}
