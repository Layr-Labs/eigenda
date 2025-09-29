package directory

// All contracts that the EigenDA offchain code interacts with should be defined here.
// It is ok to remove contracts from this list if the offchain code doesn't interact with them anymore.

// When you add to this list, make sure you keep things in alphabetical order.

const (
	CertVerifierRouter     ContractName = "CERT_VERIFIER_ROUTER"
	EigenDAEjectionManager ContractName = "EIGEN_DA_EJECTION_MANAGER"
	OperatorStateRetriever ContractName = "OPERATOR_STATE_RETRIEVER"
	PaymentVault           ContractName = "PAYMENT_VAULT"
	RegistryCoordinator    ContractName = "REGISTRY_COORDINATOR"
	RelayRegistry          ContractName = "RELAY_REGISTRY"
	ServiceManager         ContractName = "SERVICE_MANAGER"
)

// a list of all contracts currently known to the EigenDA offchain code.
var knownContracts = []ContractName{
	CertVerifierRouter,
	EigenDAEjectionManager,
	OperatorStateRetriever,
	PaymentVault,
	RegistryCoordinator,
	RelayRegistry,
	ServiceManager,
}
