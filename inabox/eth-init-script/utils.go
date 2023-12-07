package ethinitscript

import (
	"log"

	"github.com/Layr-Labs/eigenda/inabox/utils"
)

// Executes a forge script with a given rpc and private key
func execForgeScript(script, privateKey string, rpcUrl string, extraArgs []string) {
	args := []string{"script", script,
		"--rpc-url", rpcUrl,
		"--private-key", privateKey,
		"--broadcast",
	}

	// if deployer.VerifyContracts {
	// 	args = append(args, "--verify",
	// 		"--verifier", "blockscout",
	// 		"--verifier-url", deployer.VerifierURL)
	// }

	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	// The following code converts the forge call into a docker call
	err := utils.RunCommand("forge", args...)
	if err != nil {
		log.Fatal(err.Error() + "\n")
	}

	log.Print("Forge script ran succesfully!")
}
