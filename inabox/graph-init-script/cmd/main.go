package main

import (
	"github.com/Layr-Labs/eigenda/inabox/config"
	graphinitscript "github.com/Layr-Labs/eigenda/inabox/graph-init-script"
	"github.com/Layr-Labs/eigenda/inabox/utils"
)

func main() {
	config := config.OpenCwdConfig()
	startBlock := utils.GetLatestBlockNumber(config.Deployers[0].RPC)
	graphinitscript.DeploySubgraphs(config, startBlock)
}
