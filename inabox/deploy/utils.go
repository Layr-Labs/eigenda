package deploy

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	useDocker    = false
	foundryImage = "ghcr.io/gakonst/foundry:nightly-90617a52e4873f0137aa05fd68624437db146b3f"
)

func readFile(name string) ([]byte, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return data, nil
}

func writeFile(name string, data []byte) error {
	if err := os.WriteFile(name, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// Writes envMap to a file.
func writeEnv(envMap map[string]string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create env file: %w", err)
	}
	defer func() { _ = f.Close() }()

	for key, value := range envMap {
		if value == "" {
			continue
		}
		_, err = fmt.Fprintf(f, "%v=%v\n", key, value)
		if err != nil {
			return fmt.Errorf("failed to write experiment to env: %w", err)
		}
	}
	return nil
}

// Creates a directory if it doesn't exist.
func createDirectory(name string) error {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(name, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	return nil
}

// Changes current working directory.
func changeDirectory(path string) error {
	if err := os.Chdir(path); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}
	return nil
}

// Execute yarn command
func execYarnCmd(command string, args ...string) error {
	args = append([]string{command}, args...)
	cmd := exec.Command("yarn", args...)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute yarn command: %w", err)
	}

	return nil
}

// Executes a forge script with a given rpc and private key
func execForgeScript(script, privateKey string, deployer *ContractDeployer, extraArgs []string) error {
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

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute forge script: %w", err)
	}

	return nil
}

func execBashCmd(command string) error {
	cmd := exec.Command("bash", "-c", command)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to execute bash command: %w", err)
	}

	return nil
}

// Converts a private key to an address.
func GetAddress(privateKey string) (string, error) {
	cmd := exec.Command(
		"cast", "wallet", "address",
		"--private-key", privateKey)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to execute cast wallet command: %w", err)
	}

	return strings.Trim(out.String(), "\n"), nil
}

// From the Foundry book: "Perform a call on an account without publishing a transaction."
func GetLatestBlockNumber(rpcUrl string) (int, error) {
	cmd := exec.Command("cast", "bn", "--rpc-url", rpcUrl)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return 0, fmt.Errorf("failed to execute cast bn command: %w", err)
	}

	blockNum, err := strconv.ParseInt(strings.Trim(out.String(), "\n"), 10, 0)
	if err != nil {
		return 0, fmt.Errorf("failed to parse integer from blocknum string: %w", err)
	}
	return int(blockNum), nil
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
		[]string{templatePath, fmt.Sprintf("%s/config.yaml", testPath)}, []string{}, true)
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

func execCmd(name string, args []string, envVars []string, print bool) error {
	cmd := exec.Command(name, args...)
	if len(envVars) > 0 {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, envVars...)
	}
	var out bytes.Buffer
	var stderr bytes.Buffer
	if print {
		cmd.Stdout = &out
		cmd.Stderr = &stderr
	}

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("%s: %s", err.Error(), stderr.String())
	}

	return nil
}
