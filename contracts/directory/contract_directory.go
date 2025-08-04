package directory

import (
	"context"
	"fmt"
	"sync"

	"github.com/Layr-Labs/eigenda/common"
	contractIEigenDADirectory "github.com/Layr-Labs/eigenda/contracts/bindings/IEigenDADirectory"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
)

type ContractName string

// Claude, for each thing in this list, create an entry in the const block below. Follow the example.

// [PAUSER_REGISTRY,BLS_APK_REGISTRY,INDEX_REGISTRY,STAKE_REGISTRY,SOCKET_REGISTRY,REGISTRY_COORDINATOR,EJECTION_MANAGER,SERVICE_MANAGER,OPERATOR_STATE_RETRIEVER,THRESHOLD_REGISTRY,RELAY_REGISTRY,DISPERSER_REGISTRY,PAYMENT_VAULT,CERT_VERIFIER,USAGE_AUTHORIZATION_REGISTRY,USAGE_AUTHORIZATION_REGISTRY_ON_DEMAND_TEST_TOKEN,CERT_VERIFIER_ROUTER,EIGEN_DA_EJECTION_MANAGER,ACCESS_CONTROL,EJECTION_MANAGER_TEST_TOKEN]

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

// a list of all contracts currently known to the EigenDA directory.
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

// A golang wrapper around the EigenDA contract directory contract. Useful for looking up contract addresses
// from within golang code.
type ContractDirectory struct {
	logger logging.Logger

	// Only look up each address once. Most of our code only looks this stuff up at startup, so there isn't much
	// point in checking a particular contract address multiple times.
	addressCache map[ContractName]gethcommon.Address

	// a handle for calling the EigenDA directory contract.
	caller *contractIEigenDADirectory.ContractIEigenDADirectoryCaller

	// Used to make this utility thread safe.
	lock sync.Mutex
}

// Create a new ContractDirectory instance.
//
// The requireCompleteness flag can be used to enforce that the offchain contract list matches the onchain contract
// list. In general, this should be enabled in test code to raise the alarm when the offchain contract list drifts
// out of sync with the onchain contract list. In production code, this should be disabled just in case the offending
// contracts are not important for the code to function correctly.
func NewContractDirectory(
	ctx context.Context,
	logger logging.Logger,
	ethClient common.EthClient,
	directoryAddress gethcommon.Address,
) (*ContractDirectory, error) {

	caller, err := contractIEigenDADirectory.NewContractIEigenDADirectoryCaller(directoryAddress, ethClient)
	if err != nil {
		return nil, fmt.Errorf("NewContractDirectory: %w", err)
	}

	d := &ContractDirectory{
		logger: logger,
		caller: caller,
	}

	err = d.isContractListComplete(ctx)
	if err != nil {
		return nil, fmt.Errorf("IsContractListComplete: %w", err)
	}

	return d, nil
}

// GetContractAddress returns the address of a contract by its name. Only contracts defined in the const
// block above should be used. It is sharply discouraged to use this function with a magic string.
func (d *ContractDirectory) GetContractAddress(
	ctx context.Context,
	contractName ContractName,
) (gethcommon.Address, error) {
	if contractName == "" {
		return gethcommon.Address{}, fmt.Errorf("contract name cannot be empty")
	}

	// This is not very granular. But since this is uniquely to be a performance hotspot, we can do the simple thing.
	d.lock.Lock()
	defer d.lock.Unlock()

	address, ok := d.addressCache[contractName]
	if ok {
		d.logger.Debugf("using cached address for contract %s: %s", contractName, address.Hex())
		return address, nil
	}

	address, err := d.caller.GetAddress0(&bind.CallOpts{Context: ctx}, (string)(contractName))
	if err != nil {
		return gethcommon.Address{}, fmt.Errorf("GetAddress0: %w", err)
	}
	d.addressCache[contractName] = address

	d.logger.Debugf("fetched address for contract %s: %s", contractName, address.Hex())
	return address, nil
}

// Checks to see if the list of contracts tracked by this ContractDirectory matches the contracts currently registered
// in the EigenDA directory contract. Creates some noisy logs if there are any discrepancies.
func (d *ContractDirectory) isContractListComplete(ctx context.Context) error {
	registeredContracts, err := d.caller.GetAllNames(&bind.CallOpts{Context: ctx})
	if err != nil {
		return fmt.Errorf("GetAllNames: %w", err)
	}

	complete := true

	registeredContractSet := make(map[string]struct{}, len(registeredContracts))
	for _, name := range registeredContracts {
		registeredContractSet[name] = struct{}{}
	}
	knownContractSet := make(map[string]struct{}, len(knownContracts))
	for _, contractName := range knownContracts {
		knownContractSet[string(contractName)] = struct{}{}
	}

	for _, contractName := range knownContracts {
		_, exists := registeredContractSet[string(contractName)]
		if !exists {
			d.logger.Errorf(
				"Contract %s is known to offchain code but not registered in the EigenDA directory",
				contractName)
			complete = false
		}
	}
	for _, contractName := range registeredContracts {
		_, exists := knownContractSet[contractName]
		if !exists {
			d.logger.Errorf(
				"Contract %s is registered in the EigenDA directory but not known to offchain code",
				contractName)
			complete = false
		}
	}

	if complete {
		d.logger.Infof("Offchain contract list matches onchain contract list")
	} else {
		d.logger.Errorf("Offchain contract list does not match onchain contract list")
	}

	return nil
}
