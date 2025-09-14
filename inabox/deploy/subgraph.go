package deploy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

type subgraphUpdater interface {
	updateSubgraph(s *Subgraph, startBlock int)
	updateNetworks(n Networks, startBlock int)
}

type eigenDAOperatorStateSubgraphUpdater struct {
	c *Config
}

func (u eigenDAOperatorStateSubgraphUpdater) updateSubgraph(s *Subgraph, startBlock int) {
	s.DataSources[0].Source.Address = strings.TrimPrefix(u.c.EigenDA.RegistryCoordinator, "0x")
	s.DataSources[0].Source.StartBlock = startBlock
	s.DataSources[1].Source.Address = strings.TrimPrefix(u.c.EigenDA.BlsApkRegistry, "0x")
	s.DataSources[1].Source.StartBlock = startBlock
	s.DataSources[2].Source.Address = strings.TrimPrefix(u.c.EigenDA.BlsApkRegistry, "0x")
	s.DataSources[2].Source.StartBlock = startBlock
	s.DataSources[3].Source.Address = strings.TrimPrefix(u.c.EigenDA.RegistryCoordinator, "0x")
	s.DataSources[3].Source.StartBlock = startBlock
	s.DataSources[4].Source.Address = strings.TrimPrefix(u.c.EigenDA.BlsApkRegistry, "0x")
	s.DataSources[4].Source.StartBlock = startBlock
}

func (u eigenDAOperatorStateSubgraphUpdater) updateNetworks(n Networks, startBlock int) {
	n["devnet"]["RegistryCoordinator"]["address"] = u.c.EigenDA.RegistryCoordinator
	n["devnet"]["RegistryCoordinator"]["startBlock"] = startBlock
	n["devnet"]["RegistryCoordinator_Operator"]["address"] = u.c.EigenDA.RegistryCoordinator
	n["devnet"]["RegistryCoordinator_Operator"]["startBlock"] = startBlock

	n["devnet"]["BLSApkRegistry"]["address"] = u.c.EigenDA.BlsApkRegistry
	n["devnet"]["BLSApkRegistry"]["startBlock"] = startBlock
	n["devnet"]["BLSApkRegistry_Operator"]["address"] = u.c.EigenDA.BlsApkRegistry
	n["devnet"]["BLSApkRegistry_Operator"]["startBlock"] = startBlock
	n["devnet"]["BLSApkRegistry_QuorumApkUpdates"]["address"] = u.c.EigenDA.BlsApkRegistry
	n["devnet"]["BLSApkRegistry_QuorumApkUpdates"]["startBlock"] = startBlock
}

type eigenDAUIMonitoringUpdater struct {
	c *Config
}

func (u eigenDAUIMonitoringUpdater) updateSubgraph(s *Subgraph, startBlock int) {
	s.DataSources[0].Source.Address = strings.TrimPrefix(u.c.EigenDA.ServiceManager, "0x")
	s.DataSources[0].Source.StartBlock = startBlock
}

func (u eigenDAUIMonitoringUpdater) updateNetworks(n Networks, startBlock int) {
	n["devnet"]["EigenDAServiceManager"]["address"] = u.c.EigenDA.ServiceManager
	n["devnet"]["EigenDAServiceManager"]["startBlock"] = startBlock
}

func (env *Config) deploySubgraphs(startBlock int) error {
	if !env.Environment.IsLocal() {
		return nil
	}

	fmt.Println("Deploying Subgraph")
	if err := env.deploySubgraph(
		eigenDAOperatorStateSubgraphUpdater{c: env},
		"eigenda-operator-state",
		startBlock,
	); err != nil {
		return fmt.Errorf("failed to deploy eigenda-operator-state subgraph: %w", err)
	}

	if err := env.deploySubgraph(eigenDAUIMonitoringUpdater{c: env}, "eigenda-batch-metadata", startBlock); err != nil {
		return fmt.Errorf("failed to deploy eigenda-batch-metadata subgraph: %w", err)
	}

	return nil
}

func (env *Config) deploySubgraph(updater subgraphUpdater, path string, startBlock int) error {

	subgraphPath := filepath.Join(env.rootPath, "subgraphs", path)
	if err := changeDirectory(subgraphPath); err != nil {
		return fmt.Errorf("error changing directories: %w", err)
	}

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil {
		logger.Info("Successfully changed to absolute path", "path", cwd)
	}

	logger.Debug("Executing bash command", "command", `cp "./templates/subgraph.yaml" "./"`)
	if err := execBashCmd(`cp "./templates/subgraph.yaml" "./"`); err != nil {
		return fmt.Errorf("failed to copy subgraph.yaml: %w", err)
	}

	logger.Debug("Executing bash command", "command", `cp "./templates/networks.json" "./"`)
	if err := execBashCmd(`cp "./templates/networks.json" "./"`); err != nil {
		return fmt.Errorf("failed to copy networks.json: %w", err)
	}

	env.updateSubgraph(updater, subgraphPath, startBlock)

	logger.Debug("Executing yarn install")
	if err := execYarnCmd("install"); err != nil {
		return fmt.Errorf("failed to execute yarn install: %w", err)
	}

	logger.Debug("Executing yarn codegen")
	if err := execYarnCmd("codegen"); err != nil {
		return fmt.Errorf("failed to execute yarn codegen: %w", err)
	}

	logger.Debug("Executing yarn remove-local")
	if err := execYarnCmd("remove-local"); err != nil {
		return fmt.Errorf("failed to execute yarn remove-local: %w", err)
	}

	logger.Debug("Executing yarn create-local")
	if err := execYarnCmd("create-local"); err != nil {
		return fmt.Errorf("failed to execute yarn create-local: %w", err)
	}

	logger.Debug("Executing yarn deploy-local")
	if err := execYarnCmd("deploy-local", "--version-label", "v0.0.1"); err != nil {
		return fmt.Errorf("failed to execute yarn deploy-local: %w", err)
	}

	return nil
}

func (env *Config) updateSubgraph(updater subgraphUpdater, path string, startBlock int) {
	const (
		networkFile  = "networks.json"
		subgraphFile = "subgraph.yaml"
	)

	currDir, _ := os.Getwd()
	if err := changeDirectory(currDir); err != nil {
		logger.Fatal("Error changing directories", "error", err)
	}
	defer func() {
		if err := changeDirectory(currDir); err != nil {
			logger.Fatal("Error changing directories", "error", err)
		}
	}()

	// Log the current working directory (absolute path)
	if cwd, err := os.Getwd(); err == nil {
		logger.Info("Successfully changed to absolute path", "path", cwd)
	}

	networkData, err := readFile(networkFile)
	if err != nil {
		logger.Fatal("Error reading networks.json", "error", err)
	}

	var networkTemplate Networks
	if err := json.Unmarshal([]byte(networkData), &networkTemplate); err != nil {
		logger.Fatal("Failed to unmarshal networks.json", "error", err)
	}
	updater.updateNetworks(networkTemplate, startBlock)
	networkJson, err := json.MarshalIndent(networkTemplate, "", "  ")
	if err != nil {
		logger.Fatal("Error marshaling networks.json", "error", err)
	}

	if err := writeFile(networkFile, networkJson); err != nil {
		logger.Fatal("Error writing networks.json", "error", err)
	}
	logger.Info("networks.json written")

	subgraphTemplateData, err := readFile(subgraphFile)
	if err != nil {
		logger.Fatal("Error reading subgraph.yaml", "error", err)
	}

	var sub Subgraph
	if err := yaml.Unmarshal(subgraphTemplateData, &sub); err != nil {
		logger.Fatal("Error unmarshaling subgraph.yaml", "error", err)
	}
	updater.updateSubgraph(&sub, startBlock)
	subgraphYaml, err := yaml.Marshal(&sub)
	if err != nil {
		logger.Fatal("Error marshaling subgraph", "error", err)
	}
	if err := writeFile(subgraphFile, subgraphYaml); err != nil {
		logger.Fatal("Error writing subgraph.yaml", "error", err)
	}
	logger.Info("subgraph.yaml written")
}
