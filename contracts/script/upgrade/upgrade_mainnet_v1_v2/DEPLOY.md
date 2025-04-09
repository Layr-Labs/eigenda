Steps to do the upgrade by each party, divided up into phases, where previous phases block the next.

## Phase 1

### Deployer 
* Deploy a new proxy admin owned by the DA Ops multisig. 
* Deploy Proxies + Implementations + Initialize for the following:
    * EigenDAThresholdRegistry
    * EigenDARelayRegistry
    * EigenDADisperserRegistry
    * PaymentVault
    * Pauser Registry
* Deploy Implementations for the following:
    * EjectionManager
    * RegistryCoordinator
    * EigenDAServiceManager

##  Phase 2

### Executor Multisig
The Core Ops multisig should call the following on the timelock.

* Call DA Proxy Admin to upgrade the following contracts:
    * EjectionManager
    * RegistryCoordinator
    * EigenDAServiceManager
* Transfer DA Proxy Admin ownership to DA Ops Multisig

### Core Ops Multisig
* Transfer ownership of the following contracts to the DA Ops Multisig (implementation level fns)
    * Registry Coordinator
    * EigenDAServiceManager
    * Ejection Manager
* Call setPauserRegistry(new pauser registry) on the RegistryCoordinator

### DA Ops Multisig
* Initialize parameters in new V2 registries
    * ThresholdRegistry
    * RelayRegistry
    * DisperserRegistry

## Phase 3
* Verify and test upgrades thoroughly during the timelock period.
* Consider upgrading other contracts, refer to notes below.
* Deploy periphery contracts like:
    * OperatorStateRetriever
    * CertVerifier

## Notes

* The SocketRegistry is not the same implementation as on master at the time of writing. Upgrading would remove the migrateOperatorSockets function.
* 