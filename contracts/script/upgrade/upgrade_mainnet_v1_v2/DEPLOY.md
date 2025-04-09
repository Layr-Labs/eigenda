# V1 to V2 Upgrade Process

## Goals

* Upgrade from EigenDA V1 to EigenDA V2 to mainnet ASAP with the following constraints:
    * DA Proxy Admin is currently controlled by the executor multisig, which requires a 10 day timelocked tx by the Core Ops Multisig
    * Implementation level ownership of DA contracts is currently the Core Ops Multisig
    * DA Ops Multisig not yet fully formed
* Transfer ownership of all DA contracts to the DA Ops Multisig
* Transfer ownership of the DA Proxy Admin to the DA Ops Multisig

## Phases

Listed are the steps to do the upgrade by each party, divided up into phases, where previous phases block the next.

### Phase 1

#### Deployer 
* Deploy a new proxy admin owned by the DA Ops multisig. 
* Deploy Proxies + Implementations + Initialize for the following:
    * EigenDAThresholdRegistry
    * EigenDARelayRegistry
    * EigenDADisperserRegistry
    * PaymentVault
* Deploy Implementations for the following:
    * EjectionManager
    * RegistryCoordinator
    * EigenDAServiceManager

###  Phase 2

#### Executor Multisig
The Core Ops multisig should call the following on the timelock.

* Call DA Proxy Admin to upgrade the following contracts:
    * EjectionManager
    * RegistryCoordinator
    * EigenDAServiceManager

#### DA Ops Multisig
* Initialize parameters in new V2 registries
    * ThresholdRegistry
    * RelayRegistry
    * DisperserRegistry

### Phase 3

#### Executor Multisig
* Transfer DA Proxy Admin ownership to DA Ops Multisig

#### Core Ops Multisig
* Transfer ownership of the following contracts to the DA Ops Multisig (implementation level fns)
    * Registry Coordinator
    * EigenDAServiceManager
    * Ejection Manager

#### Verifiers
* Verify and test upgrades thoroughly during the timelock period
* Execute timelock transactions.

### Phase 4
* Consider upgrading other contracts, refer to notes below.
* Deploy periphery contracts like:
    * OperatorStateRetriever 
    * CertVerifier 
* Merge old and new DA Proxy Admin

## Deployer's Notes

* The SocketRegistry is not the same implementation as on master at the time of writing. Upgrading would remove the migrateOperatorSockets function.
* The ThresholdRegistry, DisperserRegistry, RelayRegistry, and PaymentVault are contracts new to V2
* A new PauserRegistry is needed to not share the same one controlled by the Core Ops multisig, but this task is deferred to later.
* EjectionManager and RegistryCoordinator are to be upgraded immediately because of significant logic diffs from master to what's on chain.
* IndexRegistry + StakeRegistry + BLS APK Registry only contain code level optimizations, so upgrade is deferred.
