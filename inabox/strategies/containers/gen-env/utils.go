package genenv

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
	testPath := filepath.Join(rootPath, fmt.Sprintf("inabox/strategies/containers/testdata/%s", testName))
	err := os.MkdirAll(testPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create test directory: %s", err.Error())
	}

	// Copy the template to the new test directory
	templatePath := filepath.Join(rootPath, fmt.Sprintf("inabox/strategies/containers/templates/%s", templateName))
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
