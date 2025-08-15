package directory

// All contracts that the EigenDA offchain code interacts with should be defined here.
// It is ok to remove contracts from this list if the offchain code doesn't interact with them anymore.

const (
	ServiceManager            ContractName = "SERVICE_MANAGER"
	BLSOperatorStateRetriever ContractName = "OPERATOR_STATE_RETRIEVER"
	EjectionManager           ContractName = "EJECTION_MANAGER"
	RegistryCoordinator       ContractName = "REGISTRY_COORDINATOR"
	RelayRegistry             ContractName = "RELAY_REGISTRY"
)

// a list of all contracts currently known to the EigenDA offchain code.
var knownContracts = []ContractName{
	ServiceManager,
	BLSOperatorStateRetriever,
	EjectionManager,
	RegistryCoordinator,
	RelayRegistry,
}
