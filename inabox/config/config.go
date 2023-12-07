package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Layr-Labs/eigenda/inabox/utils"
)

func NewTestConfig(rootPath, testName string) (testEnv *Config) {
	return OpenConfig(filepath.Join(rootPath, "inabox/testdata", testName, "config.yaml"))
}

func OpenConfig(file string) (testEnv *Config) {
	data := utils.MustReadFile(file)
	var config Config
	utils.MustUnmarshalYaml[Config](data, &config)
	return &config
}

func OpenCwdConfig() *Config {
	cwd, err := os.Getwd()
	if err != nil {
		log.Panicf("Couldn't get cwd: %v", err)
	}
	return OpenConfig(filepath.Join(cwd, "config.yaml"))
}

func (env *Config) GetDeployer(name string) (*ContractDeployer, bool) {
	for _, deployer := range env.Deployers {
		if deployer.Name == name {
			return deployer, true
		}
	}
	return nil, false
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
	err = utils.RunCommand("cp", templatePath, filepath.Join(testPath, "config.yaml"))
	if err != nil {
		return "", fmt.Errorf("failed to copy template to test directory: %s", err.Error())
	}

	return testName, nil
}
