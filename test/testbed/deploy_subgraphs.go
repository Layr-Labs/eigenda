package testbed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Layr-Labs/eigensdk-go/logging"
	"gopkg.in/yaml.v3"
)

// Subgraph yaml
type Subgraph struct {
	DataSources []DataSources `yaml:"dataSources"`
	Schema      Schema        `yaml:"schema"`
	SpecVersion string        `yaml:"specVersion"`
}

type DataSources struct {
	Kind    string  `yaml:"kind"`
	Mapping Mapping `yaml:"mapping"`
	Name    string  `yaml:"name"`
	Network string  `yaml:"network"`
	Source  Source  `yaml:"source"`
}

type Schema struct {
	File string `yaml:"file"`
}

type Source struct {
	Abi        string `yaml:"abi"`
	Address    string `yaml:"address"`
	StartBlock int    `yaml:"startBlock"`
}

type Mapping struct {
	Abis          []Abis         `yaml:"abis"`
	ApiVersion    string         `yaml:"apiVersion"`
	Entities      []string       `yaml:"entities"`
	EventHandlers []EventHandler `yaml:"eventHandlers"`
	BlockHandlers []BlockHandler `yaml:"blockHandlers"`
	File          string         `yaml:"file"`
	Kind          string         `yaml:"kind"`
	Language      string         `yaml:"language"`
}

type Abis struct {
	File string `yaml:"file"`
	Name string `yaml:"name"`
}

type EventHandler struct {
	Event   string `yaml:"event"`
	Handler string `yaml:"handler"`
}

type BlockHandler struct {
	Handler string `yaml:"handler"`
}

type Networks map[string]map[string]map[string]any

type SubgraphUpdater interface {
	UpdateSubgraph(s *Subgraph, startBlock int)
	UpdateNetworks(n Networks, startBlock int)
}

type EigenDAOperatorStateSubgraphUpdater struct {
	RegistryCoordinator string
	BlsApkRegistry      string
}

func (u EigenDAOperatorStateSubgraphUpdater) UpdateSubgraph(s *Subgraph, startBlock int) {
	s.DataSources[0].Source.Address = strings.TrimPrefix(u.RegistryCoordinator, "0x")
	s.DataSources[0].Source.StartBlock = startBlock
	s.DataSources[1].Source.Address = strings.TrimPrefix(u.BlsApkRegistry, "0x")
	s.DataSources[1].Source.StartBlock = startBlock
	s.DataSources[2].Source.Address = strings.TrimPrefix(u.BlsApkRegistry, "0x")
	s.DataSources[2].Source.StartBlock = startBlock
	s.DataSources[3].Source.Address = strings.TrimPrefix(u.RegistryCoordinator, "0x")
	s.DataSources[3].Source.StartBlock = startBlock
	s.DataSources[4].Source.Address = strings.TrimPrefix(u.BlsApkRegistry, "0x")
	s.DataSources[4].Source.StartBlock = startBlock
}

func (u EigenDAOperatorStateSubgraphUpdater) UpdateNetworks(n Networks, startBlock int) {
	n["devnet"]["RegistryCoordinator"]["address"] = u.RegistryCoordinator
	n["devnet"]["RegistryCoordinator"]["startBlock"] = startBlock
	n["devnet"]["RegistryCoordinator_Operator"]["address"] = u.RegistryCoordinator
	n["devnet"]["RegistryCoordinator_Operator"]["startBlock"] = startBlock

	n["devnet"]["BLSApkRegistry"]["address"] = u.BlsApkRegistry
	n["devnet"]["BLSApkRegistry"]["startBlock"] = startBlock
	n["devnet"]["BLSApkRegistry_Operator"]["address"] = u.BlsApkRegistry
	n["devnet"]["BLSApkRegistry_Operator"]["startBlock"] = startBlock
	n["devnet"]["BLSApkRegistry_QuorumApkUpdates"]["address"] = u.BlsApkRegistry
	n["devnet"]["BLSApkRegistry_QuorumApkUpdates"]["startBlock"] = startBlock
}

type EigenDAUIMonitoringUpdater struct {
	ServiceManager string
}

func (u EigenDAUIMonitoringUpdater) UpdateSubgraph(s *Subgraph, startBlock int) {
	s.DataSources[0].Source.Address = strings.TrimPrefix(u.ServiceManager, "0x")
	s.DataSources[0].Source.StartBlock = startBlock
}

func (u EigenDAUIMonitoringUpdater) UpdateNetworks(n Networks, startBlock int) {
	n["devnet"]["EigenDAServiceManager"]["address"] = u.ServiceManager
	n["devnet"]["EigenDAServiceManager"]["startBlock"] = startBlock
}

// SubgraphDeploymentConfig contains configuration for deploying subgraphs
type SubgraphDeploymentConfig struct {
	RootPath            string
	RegistryCoordinator string
	BlsApkRegistry      string
	ServiceManager      string
	Logger              logging.Logger
}

// DeploySubgraphs deploys the subgraphs for EigenDA
func DeploySubgraphs(config SubgraphDeploymentConfig, startBlock int) error {
	if config.Logger == nil {
		return fmt.Errorf("logger is required")
	}

	config.Logger.Info("Deploying Subgraphs", "startBlock", startBlock)

	// Deploy eigenda-operator-state subgraph
	if err := deploySubgraph(
		config,
		EigenDAOperatorStateSubgraphUpdater{
			RegistryCoordinator: config.RegistryCoordinator,
			BlsApkRegistry:      config.BlsApkRegistry,
		},
		"eigenda-operator-state",
		startBlock,
	); err != nil {
		return fmt.Errorf("failed to deploy eigenda-operator-state subgraph: %w", err)
	}

	// Deploy eigenda-batch-metadata subgraph
	if err := deploySubgraph(
		config,
		EigenDAUIMonitoringUpdater{ServiceManager: config.ServiceManager},
		"eigenda-batch-metadata",
		startBlock,
	); err != nil {
		return fmt.Errorf("failed to deploy eigenda-batch-metadata subgraph: %w", err)
	}

	return nil
}

func deploySubgraph(config SubgraphDeploymentConfig, updater SubgraphUpdater, path string, startBlock int) error {
	if config.Logger == nil {
		return fmt.Errorf("logger is required")
	}

	config.Logger.Info("Deploying Subgraph", "path", path, "startBlock", startBlock)

	subgraphPath := filepath.Join(config.RootPath, "subgraphs", path)
	if err := os.Chdir(subgraphPath); err != nil {
		return fmt.Errorf("error changing directories: %w", err)
	}

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil && config.Logger != nil {
		config.Logger.Info("Successfully changed to absolute path", "path", cwd)
	}

	config.Logger.Debug("Executing bash command", "command", `cp "./templates/subgraph.yaml" "./"`)
	if err := execBashCmd(`cp "./templates/subgraph.yaml" "./"`); err != nil {
		return fmt.Errorf("failed to copy subgraph.yaml: %w", err)
	}

	config.Logger.Debug("Executing bash command", "command", `cp "./templates/networks.json" "./"`)
	if err := execBashCmd(`cp "./templates/networks.json" "./"`); err != nil {
		return fmt.Errorf("failed to copy networks.json: %w", err)
	}

	if err := updateSubgraph(config, updater, startBlock); err != nil {
		return fmt.Errorf("failed to update subgraph: %w", err)
	}

	config.Logger.Debug("Executing yarn install")
	if err := execYarnCmd("install"); err != nil {
		return fmt.Errorf("failed to execute yarn install: %w", err)
	}

	config.Logger.Debug("Executing yarn codegen")
	if err := execYarnCmd("codegen"); err != nil {
		return fmt.Errorf("failed to execute yarn codegen: %w", err)
	}

	config.Logger.Debug("Executing yarn remove-local")
	if err := execYarnCmd("remove-local"); err != nil {
		return fmt.Errorf("failed to execute yarn remove-local: %w", err)
	}

	config.Logger.Debug("Executing yarn create-local")
	if err := execYarnCmd("create-local"); err != nil {
		return fmt.Errorf("failed to execute yarn create-local: %w", err)
	}

	config.Logger.Debug("Executing yarn deploy-local")
	if err := execYarnCmd("deploy-local", "--version-label", "v0.0.1"); err != nil {
		return fmt.Errorf("failed to execute yarn deploy-local: %w", err)
	}

	return nil
}

func updateSubgraph(config SubgraphDeploymentConfig, updater SubgraphUpdater, startBlock int) error {
	const (
		networkFile  = "networks.json"
		subgraphFile = "subgraph.yaml"
	)

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil && config.Logger != nil {
		config.Logger.Info("Current working directory", "path", cwd)
	}

	networkData, err := os.ReadFile(networkFile)
	if err != nil {
		return fmt.Errorf("error reading networks.json: %w", err)
	}

	var networkTemplate Networks
	if err := json.Unmarshal(networkData, &networkTemplate); err != nil {
		return fmt.Errorf("failed to unmarshal networks.json: %w", err)
	}
	updater.UpdateNetworks(networkTemplate, startBlock)
	networkJson, err := json.MarshalIndent(networkTemplate, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling networks.json: %w", err)
	}

	if err := os.WriteFile(networkFile, networkJson, 0644); err != nil {
		return fmt.Errorf("error writing networks.json: %w", err)
	}
	if config.Logger != nil {
		config.Logger.Info("networks.json written")
	}

	subgraphTemplateData, err := os.ReadFile(subgraphFile)
	if err != nil {
		return fmt.Errorf("error reading subgraph.yaml: %w", err)
	}

	var sub Subgraph
	if err := yaml.Unmarshal(subgraphTemplateData, &sub); err != nil {
		return fmt.Errorf("error unmarshaling subgraph.yaml: %w", err)
	}
	updater.UpdateSubgraph(&sub, startBlock)
	subgraphYaml, err := yaml.Marshal(&sub)
	if err != nil {
		return fmt.Errorf("error marshaling subgraph: %w", err)
	}
	if err := os.WriteFile(subgraphFile, subgraphYaml, 0644); err != nil {
		return fmt.Errorf("error writing subgraph.yaml: %w", err)
	}

	config.Logger.Info("subgraph.yaml written")
	return nil
}

// Helper functions for executing commands

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
