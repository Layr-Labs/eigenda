package deploy

import (
	"bytes"
	"errors"
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
	useDocker    = false
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

// Writes envMap to a file.
func writeEnv(envMap map[string]string, filename string) {

	f, err := os.Create(filename)
	if err != nil {
		log.Panicf("Failed to write experiment to env. Error: %s", err)
	}

	for key, value := range envMap {
		if value == "" {
			continue
		}
		_, err = f.WriteString(fmt.Sprintf("%v=%v\n", key, value))
		if err != nil {
			log.Panicf("Failed to write experiment to env. Error: %s", err)
		}
	}

	// err := godotenv.Write(envMap, filename)

}

// Creates a directory if it doesn't exist.
func createDirectory(name string) {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(name, os.ModePerm)
		if err != nil {
			log.Panicf("Failed to create directory. Error: %s", err)
		}
	}
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

// Execute yarn command
func execYarnCmd(command string, args ...string) {
	log.Printf("Executing yarn with command: %s", command)

	args = append([]string{command}, args...)
	cmd := exec.Command("yarn", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Print(fmt.Sprint(err) + ": " + stderr.String())
		log.Panicf("Failed to execute yarn command (%s). Err: %s", command, err)
	} else {
		log.Print(out.String())
	}

	log.Print("yarn command ran successfully")
}

// Executes a forge script with a given rpc and private key
func execForgeScript(script, privateKey string, deployer *ContractDeployer, extraArgs []string) {
	log.Printf("Executing forge script with params %s, %s, %s", script, deployer.RPC, privateKey)

	// Execute forge script
	var cmd *exec.Cmd

	args := []string{"script", script,
		"--rpc-url", deployer.RPC,
		"--private-key", privateKey,
		"--broadcast"}

	if deployer.VerifyContracts {
		args = append(args, "--verify",
			"--verifier", "blockscout",
			"--verifier-url", deployer.VerifierURL)
	}

	if deployer.Slow {
		args = append(args, "--slow")
	}

	if len(extraArgs) > 0 {
		args = append(args, extraArgs...)
	}

	// The following code converts the forge call into a docker call
	if useDocker {
		pwd, _ := os.Getwd()
		argString := fmt.Sprintf("docker run -v %v:/app -w /app %v \"forge %v\"", pwd, foundryImage, strings.Join(args[:], " "))
		cmd = exec.Command("/bin/sh", "-c", argString)
	} else {
		cmd = exec.Command("forge", args...)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	log.Println("Executing forge script with command: ", cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Panicf("Failed to execute forge script: %s\n Err: %s\n--- std out ---\n%s\n--- std err ---\n%s\n",
			cmd, err, out.String(), stderr.String())
	} else {
		log.Print(out.String())
	}

	log.Print("Forge script ran successfully!")
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

	log.Printf("bash command succeeded with params: %s", cmd)
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

	//log.Print("Cast wallet command ran successfully")
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

	log.Print("Cast call command ran successfully")
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

	log.Print("Cast bn command ran successfully")
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
	testPath := filepath.Join(rootPath, fmt.Sprintf("inabox/testdata/%s", testName))
	err := os.MkdirAll(testPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create test directory: %s", err.Error())
	}

	// Copy the template to the new test directory
	templatePath := filepath.Join(rootPath, fmt.Sprintf("inabox/templates/%s", templateName))
	err = execCmd(
		"cp",
		[]string{templatePath, fmt.Sprintf("%s/config.yaml", testPath)}, []string{})
	if err != nil {
		return "", fmt.Errorf("failed to copy template to test directory: %s", err.Error())
	}

	return testName, nil

}

func GetLatestTestDirectory(rootPath string) (string, error) {
	files, err := os.ReadDir(filepath.Join(rootPath, "inabox", "testdata"))
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", errors.New("no default experiment available")
	}
	testname := files[len(files)-1].Name()
	return testname, nil
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
