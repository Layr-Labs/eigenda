package graphinitscript

import (
	"encoding/json"
	"log"
	"path/filepath"
	"strings"

	"github.com/Layr-Labs/eigenda/inabox/strategies/containers/config"
	"gopkg.in/yaml.v3"
)

type EigenDAOperatorStateSubgraphUpdater struct {
	c *config.Config
}

func (u EigenDAOperatorStateSubgraphUpdater) UpdateSubgraph(s *Subgraph, startBlock int) {
	s.DataSources[0].Source.Address = strings.TrimPrefix(u.c.EigenDA.RegistryCoordinatorWithIndices, "0x")
	s.DataSources[0].Source.StartBlock = startBlock
	s.DataSources[1].Source.Address = strings.TrimPrefix(u.c.EigenDA.PubkeyRegistry, "0x")
	s.DataSources[1].Source.StartBlock = startBlock
	s.DataSources[2].Source.Address = strings.TrimPrefix(u.c.EigenDA.PubkeyCompendium, "0x")
	s.DataSources[2].Source.StartBlock = startBlock
	s.DataSources[3].Source.Address = strings.TrimPrefix(u.c.EigenDA.PubkeyCompendium, "0x")
	s.DataSources[3].Source.StartBlock = startBlock
	s.DataSources[4].Source.Address = strings.TrimPrefix(u.c.EigenDA.RegistryCoordinatorWithIndices, "0x")
	s.DataSources[4].Source.StartBlock = startBlock
	s.DataSources[5].Source.Address = strings.TrimPrefix(u.c.EigenDA.PubkeyRegistry, "0x")
	s.DataSources[5].Source.StartBlock = startBlock
}

func (u EigenDAOperatorStateSubgraphUpdater) UpdateNetworks(n Networks, startBlock int) {
	n["devnet"]["BLSRegistryCoordinatorWithIndices"]["address"] = u.c.EigenDA.RegistryCoordinatorWithIndices
	n["devnet"]["BLSRegistryCoordinatorWithIndices"]["startBlock"] = startBlock
	n["devnet"]["BLSRegistryCoordinatorWithIndices_Operator"]["address"] = u.c.EigenDA.RegistryCoordinatorWithIndices
	n["devnet"]["BLSRegistryCoordinatorWithIndices_Operator"]["startBlock"] = startBlock

	n["devnet"]["BLSPubkeyRegistry"]["address"] = u.c.EigenDA.PubkeyRegistry
	n["devnet"]["BLSPubkeyRegistry"]["startBlock"] = startBlock
	n["devnet"]["BLSPubkeyRegistry_QuorumApkUpdates"]["address"] = u.c.EigenDA.PubkeyRegistry
	n["devnet"]["BLSPubkeyRegistry_QuorumApkUpdates"]["startBlock"] = startBlock

	n["devnet"]["BLSPubkeyCompendium"]["address"] = u.c.EigenDA.PubkeyCompendium
	n["devnet"]["BLSPubkeyCompendium"]["startBlock"] = startBlock
	n["devnet"]["BLSPubkeyCompendium_Operator"]["address"] = u.c.EigenDA.PubkeyCompendium
	n["devnet"]["BLSPubkeyCompendium_Operator"]["startBlock"] = startBlock
}

type EigenDAUIMonitoringUpdater struct {
	c *config.Config
}

func (u EigenDAUIMonitoringUpdater) UpdateSubgraph(s *Subgraph, startBlock int) {
	s.DataSources[0].Source.Address = strings.TrimPrefix(u.c.EigenDA.ServiceManager, "0x")
	s.DataSources[0].Source.StartBlock = startBlock
}

func (u EigenDAUIMonitoringUpdater) UpdateNetworks(n Networks, startBlock int) {
	n["devnet"]["EigenDAServiceManager"]["address"] = u.c.EigenDA.ServiceManager
	n["devnet"]["EigenDAServiceManager"]["startBlock"] = startBlock
}

func DeploySubgraphs(cfg *config.Config, startBlock int) {
	log.Println("Deploying subgraphs")
	updateSubgraph(EigenDAOperatorStateSubgraphUpdater{c: cfg}, "/subgraphs/eigenda-operator-state", startBlock)
	updateSubgraph(EigenDAUIMonitoringUpdater{c: cfg}, "/subgraphs/eigenda-batch-metadata", startBlock)
}

func updateSubgraph(updater subgraphUpdater, path string, startBlock int) {

	// Update networks.json
	networkFile := filepath.Join(path, "networks.json")
	networkData := config.ReadFile(networkFile)

	var networkTemplate Networks
	if err := json.Unmarshal(networkData, &networkTemplate); err != nil {
		log.Panicf("Failed to unmarshal networks.json. Error: %s", err)
	}
	updater.UpdateNetworks(networkTemplate, startBlock)

	networkJson, err := json.MarshalIndent(networkTemplate, "", "  ")
	if err != nil {
		log.Panicf("Error: %s", err.Error())
	}
	config.WriteFile(networkFile, networkJson)
	log.Print("networks.json written")

	// Update subgraph.yaml
	subgraphFile := filepath.Join(path, "subgraph.yaml")
	subgraphTemplateData := config.ReadFile(subgraphFile)

	var sub Subgraph
	if err := yaml.Unmarshal(subgraphTemplateData, &sub); err != nil {
		log.Panicf("Error %s:", err.Error())
	}
	updater.UpdateSubgraph(&sub, startBlock)
	subgraphYaml, err := yaml.Marshal(&sub)
	if err != nil {
		log.Panic(err)
	}
	config.WriteFile(subgraphFile, subgraphYaml)
	log.Print("subgraph.yaml written")
}
