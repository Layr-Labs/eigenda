package localinit

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

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
func execForgeScript(script, privateKey string, deployer *ContractDeployer, extraArgs []string) {
	log.Printf("Executing forge script with params %s, %s, %s", script, deployer.RPC, privateKey)

	// Execute forge script
	// TODO: Figure out how to make this script run faster. I think it may have something to do with this https://github.com/foundry-rs/foundry/blob/master/crates/forge/bin/cmd/script/broadcast.rs#L155
	//       But the "--slow" parameter is not being passed here so not sure.
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

func execBashCmd(command string) {
	log.Printf("Executing bash with command: %s", command)

	cmd := exec.Command("bash", "-c", command)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Print(fmt.Sprint(err) + ": " + stderr.String())
		log.Panicf("Failed to execute bash command. Err: %s", err)
	} else {
		log.Print(out.String())
	}

	log.Printf("bash command succeded with params: %s", cmd)
}

// Converts a private key to an address.
func GetAddress(privateKey string) string {
	cmd := exec.Command(
		"cast", "wallet", "address",
		"--private-key", privateKey)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Print(fmt.Sprint(err) + ": " + stderr.String())
		log.Panicf("Failed to execute cast wallet command. Err: %s", err)
	}

	//log.Print("Cast wallet command ran succesfully")
	return strings.Trim(out.String(), "\n")
}

// From the Foundry book: "Perform a call on an account without publishing a transaction."
func CallContract(destination string, signature string, rpcUrl string) string {
	cmd := exec.Command(
		"cast", "call", destination, signature,
		"--rpc-url", rpcUrl)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Print(fmt.Sprint(err) + ": " + stderr.String())
		log.Panicf("Failed to execute cast wallet command. Err: %s", err)
	}

	log.Print("Cast call command ran succesfully")
	return strings.Trim(out.String(), "\n")
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

// Create a new test directory and copy the template to it.
func CreateNewTestDirectory(templateName, rootPath string) (string, error) {

	// Get the current date time with format '+%dD-%mM-%YY-%HH-%MM-%SS'
	testName := time.Now().Format("2006Y-01M-02D-15H-04M-05S")

	// Create the new test directory
	testPath := filepath.Join(rootPath, fmt.Sprintf("inabox-docker/testdata/%s", testName))
	err := os.MkdirAll(testPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create test directory: %s", err.Error())
	}

	// Copy the template to the new test directory
	templatePath := filepath.Join(rootPath, fmt.Sprintf("inabox-docker/templates/%s", templateName))
	err = execCmd(
		"cp",
		[]string{templatePath, fmt.Sprintf("%s/config.yaml", testPath)}, []string{})
	if err != nil {
		return "", fmt.Errorf("failed to copy template to test directory: %s", err.Error())
	}

	return testName, nil

}

func execCmd(name string, args []string, envVars []string) error {
	cmd := exec.Command(name, args...)
	if len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	// TODO: When these are uncommented, the deployer sometimes fails to start anvil
	// cmd.Stdout = &out
	// cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), stderr.String())
	}
	fmt.Print(out.String())
	return nil
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
