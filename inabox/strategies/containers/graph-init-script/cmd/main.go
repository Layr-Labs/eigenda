package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/Layr-Labs/eigenda/inabox/strategies/containers/config"
	graphinitscript "github.com/Layr-Labs/eigenda/inabox/strategies/containers/graph-init-script"
)

func main() {
	config := config.OpenCwdConfig()
	startBlock := GetLatestBlockNumber(config.Deployers[0].RPC)
	graphinitscript.DeploySubgraphs(config, startBlock)
}

// From the Foundry book: "Perform a call on an account without publishing a transaction."
func GetLatestBlockNumber(rpcUrl string) int {
	cmd := exec.Command("cast", "bn", "--rpc-url", rpcUrl)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Print(fmt.Sprint(err) + ": " + stderr.String())
		log.Panicf("Failed to execute cast wallet command. Err: %s", err)
	}

	log.Print("Cast bn command ran succesfully")
	blockNum, err := strconv.ParseInt(strings.Trim(out.String(), "\n"), 10, 0)
	if err != nil {
		log.Print(fmt.Sprint(err) + ": " + stderr.String())
		log.Panicf("Failed parse integer from blocknum string. Err: %s", err)
	}
	return int(blockNum)
}
