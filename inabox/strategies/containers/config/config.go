package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func OpenConfig(file string) (testEnv *Config) {
	data := ReadFile(file)
	err := yaml.Unmarshal(data, &testEnv)
	if err != nil {
		log.Panicf("Error %s:", err.Error())
	}

	return
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
