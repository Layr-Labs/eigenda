package deployment

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Layr-Labs/eigenda/common/testinfra/containers"
	"gopkg.in/yaml.v3"
)

// SubgraphDeploymentConfig holds configuration for subgraph deployment
type SubgraphDeploymentConfig struct {
	RootPath      string                   `json:"root_path"`
	Subgraphs     []SubgraphConfig         `json:"subgraphs"`
	EigenDAConfig EigenDAContractAddresses `json:"eigenda_config"`
}

// SubgraphConfig represents configuration for a single subgraph
type SubgraphConfig struct {
	Name    string `json:"name"`
	Path    string `json:"path"` // relative path from subgraphs directory
	Enabled bool   `json:"enabled"`
}

// EigenDAContractAddresses contains the deployed contract addresses needed for subgraph configuration
type EigenDAContractAddresses struct {
	RegistryCoordinator string `json:"registry_coordinator"`
	BlsApkRegistry      string `json:"bls_apk_registry"`
	ServiceManager      string `json:"service_manager"`
}

// DefaultSubgraphDeploymentConfig returns a default configuration for subgraph deployment
func DefaultSubgraphDeploymentConfig(rootPath string) SubgraphDeploymentConfig {
	return SubgraphDeploymentConfig{
		RootPath: rootPath,
		Subgraphs: []SubgraphConfig{
			{
				Name:    "eigenda-operator-state",
				Path:    "eigenda-operator-state",
				Enabled: true,
			},
			{
				Name:    "eigenda-batch-metadata",
				Path:    "eigenda-batch-metadata",
				Enabled: true,
			},
		},
	}
}

// DeploySubgraphs deploys EigenDA subgraphs to a Graph Node using container
// This function replaces the functionality from inabox/deploy/subgraph.go
func DeploySubgraphs(ctx context.Context, graphNode *containers.GraphNodeContainer, config SubgraphDeploymentConfig, startBlock int) error {
	if graphNode == nil {
		return fmt.Errorf("graph node container is not available")
	}

	adminURL := graphNode.AdminURL()
	ipfsURL, err := graphNode.IPFSURL(ctx)
	if err != nil {
		// Fallback to localhost for backward compatibility
		ipfsURL = "http://localhost:5001"
		fmt.Printf("Warning: Could not get IPFS URL, using fallback: %s\n", ipfsURL)
	}

	return DeploySubgraphsWithURLs(config, adminURL, ipfsURL, startBlock)
}

// DeploySubgraphsWithURLs deploys EigenDA subgraphs using provided URLs
// This is the standalone function that can be called directly without containers
func DeploySubgraphsWithURLs(config SubgraphDeploymentConfig, adminURL, ipfsURL string, startBlock int) error {
	if ipfsURL == "" {
		return fmt.Errorf("IPFS URL must be provided")
	}

	fmt.Printf("Deploying subgraphs to Graph Node (Admin URL: %s, IPFS URL: %s)\n", adminURL, ipfsURL)

	// Deploy each enabled subgraph
	for _, subgraph := range config.Subgraphs {
		if !subgraph.Enabled {
			fmt.Printf("Skipping disabled subgraph: %s\n", subgraph.Name)
			continue
		}

		fmt.Printf("Deploying subgraph: %s\n", subgraph.Name)
		err := deploySubgraph(config.RootPath, subgraph, config.EigenDAConfig, adminURL, ipfsURL, startBlock)
		if err != nil {
			return fmt.Errorf("failed to deploy subgraph %s: %w", subgraph.Name, err)
		}
	}

	fmt.Println("Successfully deployed all subgraphs")
	return nil
}

// deploySubgraph deploys a single subgraph
func deploySubgraph(rootPath string, subgraph SubgraphConfig, eigenDAConfig EigenDAContractAddresses, adminURL, ipfsURL string, startBlock int) error {
	subgraphPath := filepath.Join(rootPath, "subgraphs", subgraph.Path)

	// Change to subgraph directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	fmt.Printf("deploySubgraph: Current directory: %s\n", originalDir)
	fmt.Printf("deploySubgraph: Root path: %s\n", rootPath)
	fmt.Printf("deploySubgraph: Attempting to change to subgraph path: %s\n", subgraphPath)
	
	defer func() { _ = os.Chdir(originalDir) }()

	if err := os.Chdir(subgraphPath); err != nil {
		return fmt.Errorf("failed to change to subgraph directory %s: %w", subgraphPath, err)
	}

	// Copy template files
	if err := copyTemplateFiles(); err != nil {
		return fmt.Errorf("failed to copy template files: %w", err)
	}

	// Update subgraph configuration with contract addresses
	if err := updateSubgraphConfig(subgraph, eigenDAConfig, startBlock); err != nil {
		return fmt.Errorf("failed to update subgraph config: %w", err)
	}

	// Install dependencies and generate code
	if err := runYarnCommands(); err != nil {
		return fmt.Errorf("failed to run yarn commands: %w", err)
	}

	// Deploy to Graph Node
	subgraphName := fmt.Sprintf("Layr-Labs/%s", subgraph.Name)
	if err := deployToGraphNode(subgraphName, adminURL, ipfsURL); err != nil {
		return fmt.Errorf("failed to deploy to graph node: %w", err)
	}

	return nil
}

// copyTemplateFiles copies the template files needed for subgraph deployment
func copyTemplateFiles() error {
	// Copy subgraph.yaml template
	if err := copyFile("templates/subgraph.yaml", "subgraph.yaml"); err != nil {
		return fmt.Errorf("failed to copy subgraph.yaml template: %w", err)
	}

	// Copy networks.json template
	if err := copyFile("templates/networks.json", "networks.json"); err != nil {
		return fmt.Errorf("failed to copy networks.json template: %w", err)
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// updateSubgraphConfig updates the subgraph configuration files with contract addresses
func updateSubgraphConfig(subgraph SubgraphConfig, eigenDAConfig EigenDAContractAddresses, startBlock int) error {
	// Update networks.json
	if err := updateNetworksConfig(eigenDAConfig, startBlock); err != nil {
		return fmt.Errorf("failed to update networks.json: %w", err)
	}

	// Update subgraph.yaml
	if err := updateSubgraphYAML(subgraph, eigenDAConfig, startBlock); err != nil {
		return fmt.Errorf("failed to update subgraph.yaml: %w", err)
	}

	return nil
}

// updateNetworksConfig updates networks.json with contract addresses
func updateNetworksConfig(eigenDAConfig EigenDAContractAddresses, startBlock int) error {
	data, err := os.ReadFile("networks.json")
	if err != nil {
		return fmt.Errorf("failed to read networks.json: %w", err)
	}

	var networks map[string]map[string]map[string]interface{}
	if err := json.Unmarshal(data, &networks); err != nil {
		return fmt.Errorf("failed to unmarshal networks.json: %w", err)
	}

	// Initialize devnet if it doesn't exist
	if networks["devnet"] == nil {
		networks["devnet"] = make(map[string]map[string]interface{})
	}

	// Update contract addresses for eigenda-operator-state subgraph
	contracts := []struct {
		name    string
		address string
	}{
		{"RegistryCoordinator", eigenDAConfig.RegistryCoordinator},
		{"RegistryCoordinator_Operator", eigenDAConfig.RegistryCoordinator},
		{"BLSApkRegistry", eigenDAConfig.BlsApkRegistry},
		{"BLSApkRegistry_Operator", eigenDAConfig.BlsApkRegistry},
		{"BLSApkRegistry_QuorumApkUpdates", eigenDAConfig.BlsApkRegistry},
		{"EigenDAServiceManager", eigenDAConfig.ServiceManager},
	}

	for _, contract := range contracts {
		if networks["devnet"][contract.name] == nil {
			networks["devnet"][contract.name] = make(map[string]interface{})
		}
		networks["devnet"][contract.name]["address"] = contract.address
		networks["devnet"][contract.name]["startBlock"] = startBlock
	}

	// Write back to file
	updatedData, err := json.MarshalIndent(networks, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal networks.json: %w", err)
	}

	if err := os.WriteFile("networks.json", updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write networks.json: %w", err)
	}

	fmt.Println("networks.json updated")
	return nil
}

// SubgraphYAML represents the structure of subgraph.yaml
type SubgraphYAML struct {
	DataSources []DataSource `yaml:"dataSources"`
	Schema      Schema       `yaml:"schema"`
	SpecVersion string       `yaml:"specVersion"`
}

type DataSource struct {
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

// updateSubgraphYAML updates subgraph.yaml with contract addresses
func updateSubgraphYAML(subgraph SubgraphConfig, eigenDAConfig EigenDAContractAddresses, startBlock int) error {
	data, err := os.ReadFile("subgraph.yaml")
	if err != nil {
		return fmt.Errorf("failed to read subgraph.yaml: %w", err)
	}

	var subgraphData SubgraphYAML
	if err := yaml.Unmarshal(data, &subgraphData); err != nil {
		return fmt.Errorf("failed to unmarshal subgraph.yaml: %w", err)
	}

	// Update data sources based on subgraph type
	if subgraph.Name == "eigenda-operator-state" {
		updateOperatorStateSubgraph(&subgraphData, eigenDAConfig, startBlock)
	} else if subgraph.Name == "eigenda-batch-metadata" {
		updateBatchMetadataSubgraph(&subgraphData, eigenDAConfig, startBlock)
	}

	// Write back to file
	updatedData, err := yaml.Marshal(&subgraphData)
	if err != nil {
		return fmt.Errorf("failed to marshal subgraph.yaml: %w", err)
	}

	if err := os.WriteFile("subgraph.yaml", updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write subgraph.yaml: %w", err)
	}

	fmt.Println("subgraph.yaml updated")
	return nil
}

// updateOperatorStateSubgraph updates data sources for eigenda-operator-state subgraph
func updateOperatorStateSubgraph(subgraphData *SubgraphYAML, eigenDAConfig EigenDAContractAddresses, startBlock int) {
	for i := range subgraphData.DataSources {
		// Remove 0x prefix from addresses
		switch i {
		case 0, 3: // RegistryCoordinator sources
			subgraphData.DataSources[i].Source.Address = strings.TrimPrefix(eigenDAConfig.RegistryCoordinator, "0x")
		case 1, 2, 4: // BLSApkRegistry sources
			subgraphData.DataSources[i].Source.Address = strings.TrimPrefix(eigenDAConfig.BlsApkRegistry, "0x")
		}
		subgraphData.DataSources[i].Source.StartBlock = startBlock
	}
}

// updateBatchMetadataSubgraph updates data sources for eigenda-batch-metadata subgraph
func updateBatchMetadataSubgraph(subgraphData *SubgraphYAML, eigenDAConfig EigenDAContractAddresses, startBlock int) {
	for i := range subgraphData.DataSources {
		// ServiceManager source
		subgraphData.DataSources[i].Source.Address = strings.TrimPrefix(eigenDAConfig.ServiceManager, "0x")
		subgraphData.DataSources[i].Source.StartBlock = startBlock
	}
}

// runYarnCommands runs the necessary yarn commands for subgraph deployment
func runYarnCommands() error {
	// Install dependencies
	if err := execCommand("yarn", "install"); err != nil {
		return fmt.Errorf("failed to run yarn install: %w", err)
	}

	// Generate code
	if err := execCommand("yarn", "codegen"); err != nil {
		return fmt.Errorf("failed to run yarn codegen: %w", err)
	}

	return nil
}

// deployToGraphNode deploys the subgraph to the Graph Node
func deployToGraphNode(subgraphName, adminURL, ipfsURL string) error {
	// Remove existing subgraph (ignore errors)
	_ = execCommand("npx", "graph", "remove", "--node", adminURL, subgraphName)

	// Create subgraph
	if err := execCommand("npx", "graph", "create", "--node", adminURL, subgraphName); err != nil {
		return fmt.Errorf("failed to create subgraph: %w", err)
	}

	// Deploy subgraph
	if err := execCommand("npx", "graph", "deploy", "--node", adminURL, "--ipfs", ipfsURL, "--version-label", "v0.0.1", subgraphName); err != nil {
		return fmt.Errorf("failed to deploy subgraph: %w", err)
	}

	fmt.Printf("Successfully deployed subgraph: %s\n", subgraphName)
	return nil
}

// execCommand executes a command and returns an error if it fails
func execCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}