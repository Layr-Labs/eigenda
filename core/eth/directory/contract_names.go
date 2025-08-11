package directory

// All contracts that the EigenDA offchain code interacts with should be defined here.
// It is ok to remove contracts from this list if the offchain code doesn't interact with them anymore.

const (
	ServiceManager         ContractName = "SERVICE_MANAGER"
	EigenDAEjectionManager ContractName = "EIGEN_DA_EJECTION_MANAGER"
)

// a list of all contracts currently known to the EigenDA offchain code.
var knownContracts = []ContractName{
	ServiceManager,
	EigenDAEjectionManager,
}
