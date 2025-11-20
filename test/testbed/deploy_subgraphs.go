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

type Networks map[string]any

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
	// Update the devnet template with actual contract addresses
	n["network"] = "devnet"
	n["RegistryCoordinator_address"] = u.RegistryCoordinator
	n["RegistryCoordinator_startBlock"] = startBlock
	n["BLSApkRegistry_address"] = u.BlsApkRegistry
	n["BLSApkRegistry_startBlock"] = startBlock
	// EjectionManager is set to zero address for now
	n["EjectionManager_address"] = "0x0000000000000000000000000000000000000000"
	n["EjectionManager_startBlock"] = startBlock
}

type EigenDAUIMonitoringUpdater struct {
	ServiceManager string
}

func (u EigenDAUIMonitoringUpdater) UpdateSubgraph(s *Subgraph, startBlock int) {
	s.DataSources[0].Source.Address = strings.TrimPrefix(u.ServiceManager, "0x")
	s.DataSources[0].Source.StartBlock = startBlock
}

func (u EigenDAUIMonitoringUpdater) UpdateNetworks(n Networks, startBlock int) {
	// Update the devnet template with actual contract addresses
	n["network"] = "devnet"
	n["EigenDAServiceManager_address"] = u.ServiceManager
	n["EigenDAServiceManager_startBlock"] = startBlock
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
	subgraphsRootPath := filepath.Join(config.RootPath, "subgraphs")

	// Install dependencies in the parent subgraphs directory first
	config.Logger.Debug("Installing parent subgraphs dependencies")
	if err := execYarnCmd("install", subgraphsRootPath, config.Logger); err != nil {
		return fmt.Errorf("failed to install parent subgraphs dependencies: %w", err)
	}

	// Update the devnet template and generate subgraph.yaml using mustache
	if err := updateSubgraph(config, updater, startBlock, subgraphPath); err != nil {
		return fmt.Errorf("failed to update subgraph: %w", err)
	}

	config.Logger.Debug("Executing yarn install")
	if err := execYarnCmd("install", subgraphPath, config.Logger); err != nil {
		return fmt.Errorf("failed to execute yarn install: %w", err)
	}

	config.Logger.Debug("Executing yarn prepare:devnet")
	if err := execYarnCmd("prepare:devnet", subgraphPath, config.Logger); err != nil {
		return fmt.Errorf("failed to execute yarn prepare:devnet %w", err)
	}

	config.Logger.Debug("Executing yarn codegen")
	if err := execYarnCmd("codegen", subgraphPath, config.Logger); err != nil {
		return fmt.Errorf("failed to execute yarn codegen: %w", err)
	}

	config.Logger.Debug("Executing yarn remove-local")
	if err := execYarnCmd("remove-local", subgraphPath, config.Logger); err != nil {
		return fmt.Errorf("failed to execute yarn remove-local: %w", err)
	}

	config.Logger.Debug("Executing yarn create-local")
	if err := execYarnCmd("create-local", subgraphPath, config.Logger); err != nil {
		return fmt.Errorf("failed to execute yarn create-local: %w", err)
	}

	config.Logger.Debug("Executing yarn deploy-local")
	if err := execYarnCmd("deploy-local", subgraphPath, config.Logger, "--version-label", "v0.0.1"); err != nil {
		return fmt.Errorf("failed to execute yarn deploy-local: %w", err)
	}

	return nil
}

func updateSubgraph(
	config SubgraphDeploymentConfig,
	updater SubgraphUpdater,
	startBlock int,
	subgraphPath string,
) error {
	// Path to the devnet template file
	devnetTemplatePath := filepath.Join(subgraphPath, "templates", "devnet.json")

	// Read the devnet template
	templateData, err := os.ReadFile(devnetTemplatePath)
	if err != nil {
		return fmt.Errorf("error reading templates/devnet.json: %w", err)
	}

	// Parse the template
	var devnetTemplate Networks
	if err := json.Unmarshal(templateData, &devnetTemplate); err != nil {
		return fmt.Errorf("failed to unmarshal templates/devnet.json: %w", err)
	}

	// Update the template with actual contract addresses and start blocks
	updater.UpdateNetworks(devnetTemplate, startBlock)

	// Write the updated template back
	updatedJson, err := json.MarshalIndent(devnetTemplate, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling templates/devnet.json: %w", err)
	}

	if err := os.WriteFile(devnetTemplatePath, updatedJson, 0644); err != nil {
		return fmt.Errorf("error writing templates/devnet.json: %w", err)
	}
	config.Logger.Info("templates/devnet.json written")

	return nil
}

// Helper functions for executing commands

func execYarnCmd(command string, workingDir string, logger logging.Logger, args ...string) error {
	cmdArgs := append([]string{command}, args...)
	cmd := exec.Command("yarn", cmdArgs...)
	cmd.Dir = workingDir

	logger.Debug("Executing yarn command", "command", cmd.String(), "workingDir", workingDir)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logger.Error("Yarn command failed", "stdout", out.String(), "stderr", stderr.String())
		return fmt.Errorf("failed to execute yarn command: %w", err)
	}

	return nil
}

func execBashCmd(command string, workingDir string, logger logging.Logger) error {
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = workingDir

	logger.Debug("Executing bash command", "command", command, "workingDir", workingDir)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logger.Error("Bash command failed", "stdout", out.String(), "stderr", stderr.String())
		return fmt.Errorf("failed to execute bash command: %w", err)
	}

	return nil
}
