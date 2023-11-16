package ethinitscript

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/Layr-Labs/eigenda/inabox/strategies/containers/config"
	"github.com/joho/godotenv"
)

const (
	foundryImage = "ghcr.io/gakonst/foundry:nightly-90617a52e4873f0137aa05fd68624437db146b3f"
)

func readFile(name string) []byte {
	data, err := os.ReadFile(name)
	if err != nil {
		log.Panicf("Failed to read file. Error: %s", err)
	}

	return data
}

func writeFile(name string, data []byte) {
	err := os.WriteFile(name, data, 0644)
	if err != nil {
		log.Panicf("Failed to write file. Err: %s", err)
	}
}

// Reads and loads env map from a given file
func ReadEnv(filename string) map[string]string {
	err := godotenv.Load(filename)
	if err != nil {
		log.Panicf("Failed to load env file. Error: %s", err)
	}

	env, err := godotenv.Read(filename)
	if err != nil {
		log.Panicf("Failed to read env file. Error: %s", err)
	}

	return env
}

// Execute yarn command
func execYarnCmd(command string) {
	log.Printf("Executing yarn with command: %s", command)

	cmd := exec.Command("yarn", command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Panicf("Failed to execute yarn command (%s). Err: %s", command, err)
	}

	log.Print("yarn command ran succesfully")
}

// Executes a forge script with a given rpc and private key
func execForgeScript(script, privateKey string, deployer *config.ContractDeployer, extraArgs []string) {
	args := []string{"script", script,
		"--rpc-url", deployer.RPC,
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
	err := RunCommand("forge", args...)
	if err != nil {
		log.Fatal(err.Error() + "\n")
	}

	log.Print("Forge script ran succesfully!")
}

func RunCommand(name string, args ...string) error {
	log.Printf("Running command: %s\n", strings.Join(append([]string{name}, args...), " "))
	cmd := exec.Command(name, args...)

	// Set the output to the corresponding os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command and wait for it to finish
	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

// Changes current working directory.
func changeDirectory(path string) {
	err := os.Chdir(path)
	if err != nil {
		log.Panicf("Failed to change directories. Error: %s", err)
	}

	newDir, err := os.Getwd()
	if err != nil {
		log.Panicf("Failed to get working directory. Error: %s", err)
	}
	log.Printf("Current Working Directory: %s\n", newDir)
}
