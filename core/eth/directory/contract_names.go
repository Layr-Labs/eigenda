package directory

// TODO: pare down this list just to the contracts we actually use offchain

// All contracts that the EigenDA offchain code interacts with should be defined here.
// It is ok to remove contracts from this list if the offchain code doesn't interact with them anymore.

const (
	PauserRegistry                              ContractName = "PAUSER_REGISTRY"
	BlsApkRegistry                              ContractName = "BLS_APK_REGISTRY"
	IndexRegistry                               ContractName = "INDEX_REGISTRY"
	StakeRegistry                               ContractName = "STAKE_REGISTRY"
	SocketRegistry                              ContractName = "SOCKET_REGISTRY"
	RegistryCoordinator                         ContractName = "REGISTRY_COORDINATOR"
	EjectionManager                             ContractName = "EJECTION_MANAGER"
	ServiceManager                              ContractName = "SERVICE_MANAGER"
	OperatorStateRetriever                      ContractName = "OPERATOR_STATE_RETRIEVER"
	ThresholdRegistry                           ContractName = "THRESHOLD_REGISTRY"
	RelayRegistry                               ContractName = "RELAY_REGISTRY"
	DisperserRegistry                           ContractName = "DISPERSER_REGISTRY"
	PaymentVault                                ContractName = "PAYMENT_VAULT"
	CertVerifier                                ContractName = "CERT_VERIFIER"
	UsageAuthorizationRegistry                  ContractName = "USAGE_AUTHORIZATION_REGISTRY"
	UsageAuthorizationRegistryOnDemandTestToken ContractName = "USAGE_AUTHORIZATION_REGISTRY_ON_DEMAND_TEST_TOKEN"
	CertVerifierRouter                          ContractName = "CERT_VERIFIER_ROUTER"
	EigenDAEjectionManager                      ContractName = "EIGEN_DA_EJECTION_MANAGER"
	AccessControl                               ContractName = "ACCESS_CONTROL"
	EjectionManagerTestToken                    ContractName = "EJECTION_MANAGER_TEST_TOKEN"
)

// a list of all contracts currently known to the EigenDA offchain code.
var knownContracts = []ContractName{
	PauserRegistry,
	BlsApkRegistry,
	IndexRegistry,
	StakeRegistry,
	SocketRegistry,
	RegistryCoordinator,
	EjectionManager,
	ServiceManager,
	OperatorStateRetriever,
	ThresholdRegistry,
	RelayRegistry,
	DisperserRegistry,
	PaymentVault,
	CertVerifier,
	UsageAuthorizationRegistry,
	UsageAuthorizationRegistryOnDemandTestToken,
	CertVerifierRouter,
	EigenDAEjectionManager,
	AccessControl,
	EjectionManagerTestToken,
}
