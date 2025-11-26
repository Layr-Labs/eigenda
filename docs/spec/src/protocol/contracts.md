# EigenDA Protocol Contracts

This page describes EigenDA contracts that are managed by EigenDA related actors (see the exact [roles](#governance-roles)). For EigenDA-related contracts that are managed by rollups, see the [rollup managed contracts](../integration/spec/4-contracts.md) page.

> Warning: This page is incomplete and a work in progress as we are undergoing refactors of our contracts as well as some protocol upgrades. The details will change, but the information contained here should at least help to understand the important concepts.

## Overview
![image](../../assets/contracts-overview.png)

### Middleware Contracts

We make use of eigenlayer-middleware contracts, which are fully documented [here](https://github.com/Layr-Labs/eigenlayer-middleware/tree/dev/docs) and described [here](https://docs.eigencloud.xyz/eigenlayer/developers/concepts/eigenlayer-contracts/middleware-contracts). These contracts provide standard interfacing logic for operator state management and AVS representation. 

### Middleware Vendored Contracts

Some of the middleware contracts (e.g, `EjectionsManager`, `RegistryCoordinator`) have been directly vendored into the EigenDA project with minor modifications made.

### EigenDA Specific Contracts

The smart contracts can be found in our [repo](https://github.com/Layr-Labs/eigenda/tree/master/contracts/src/core), and the deployment addresses on different chains can be found in the [Networks](https://docs.eigenda.xyz/networks/mainnet#contract-addresses) section of our docs.


## Contracts Overview

| Contract Name                                                         | Project Category     | Deployed Behind ERC1967 Proxy? | isUsedByOffchainProtocol? |
|-----------------------------------------------------------------------|-----------------------|---------------------------------|----------------------------|
| [EigenDA Directory](#eigendadirectory)                                | [eigenda](#eigenda-specific-contracts)              | Yes                             | Yes                        |
| [Service Manager](#eigendaservicemanager)                             | [eigenda](#eigenda-specific-contracts)              | No                              | Yes                        |
| [Threshold Registry](#eigendathresholdregistry)                       | [eigenda](#eigenda-specific-contracts)              | Yes                             | Yes                        |
| [Relay Registry](#eigendarelayregistry)                               | [eigenda](#eigenda-specific-contracts)              | Yes                             | Yes                        |
| [Disperser Registry](#eigendadisperserregistry)                       | [eigenda](#eigenda-specific-contracts)              | Yes                             | Yes                        |
| [Payment Vault](#paymentvault)                                        | [eigenda](#eigenda-specific-contracts)              | Yes                             | Yes                        |
| [Pauser Registry](#pauserregistry)                                    | [middleware](#middleware-contracts)           | No                              | No                         |
| [BLS APK Registry](#blsapkblsapkregistry)                             | [middleware](#middleware-contracts)           | Yes                             | Yes                        |
| [Index Registry](#indexregistry)                                      | [middleware](#middleware-contracts)           | Yes                             | Yes                        |
| [Stake Registry](#stakeregistry)                                      | [middleware](#middleware-contracts)           | Yes                             | Yes                        |
| [Socket Registry](#socketregistry)                                    | [middleware](#middleware-contracts)           | Yes                             | Yes                        |
| [Operator State Retriever](#operatorstateretriever)                   | [middleware](#middleware-contracts)           | No                              | Yes                        |
| [Registry Coordinator](#eigendaregistrycoordinator)                   | [vendored middleware](#middleware-vendored-contracts)  | Yes                             | Yes                        |
| [Ejections Manager](#eigendaejectionsmanager)                         | [vendored middleware](#middleware-vendored-contracts)  | Yes                             | No                         |


<br />
<br />

------
### [`EigenDADirectory`](https://github.com/Layr-Labs/eigenda/blob/98a17e884de40a18ed9744e709ccc109adf273d3/contracts/src/core/EigenDADirectory.sol)
**Description**

This contract serves as the central discovery and reference point for all contracts composing the EigenDA system. It implements a lightweight namespace resolution protocol in which human-readable string keys are deterministically mapped to fixed storage slots containing `20-byte` contract address references.

**Access Mgmt**

- `Ownable` role that can do unilateral entry key modifications

**Offchain Usage**

This dynamic naming pattern requires off-chain management of canonical contract keys, allowing clients and services to retrieve on-chain system context from a single directory contract reference rather than requiring every contract address to be hard-coded or passed through environment configuration.

### [`EigenDAServiceManager`](https://github.com/Layr-Labs/eigenda/blob/98a17e884de40a18ed9744e709ccc109adf273d3/contracts/src/core/EigenDAServiceManager.sol)

**Description**
Used for onchain AVS registration with the EigenLayer protocol, EigenDA V1 batching, storing protocol params, rewards distribution, and referencing EigenDA protocol contracts:
- Inherits the [`ServiceManagerBase`](https://github.com/Layr-Labs/eigenlayer-middleware/blob/7314aef30b6a98c0156750f300b06bea629d0720/docs/ServiceManagerBase.md) for operator registration and rewards distribution.
- Manages batch settlement roles with callable function (i.e, `confirmBatch`) that allows for EigenDA V1 batches to be confirmed and settled into a storage commitment sequence.
- Stores protocol params (i.e, `BLOCK_STALE_MEASURE`, `BLOCK_STORE_DURATION`) for offchain ingestion by DA validator nodes.
- Stores non-callable references to other EigenDA protocol contracts in storage (i.e, [`DisperserRegistry`](#eigendadisperserregistry), [`ThresholdRegistry`](#eigendathresholdregistry), [`RelayRegistry`](#eigendarelayregistry), [`StakeRegistry`](#stakeregistry), [`PaymentVault`](#paymentvault)).

**Access Mgmt**

- `Pauser` role that can halt EigenDA V1 batch settlement
- `Ownable` role that can modify batch confirmer EOA allow-list, AVS metadata, `RewardsClaimee`, and `RewardsInitiator`
- `RegistryCoordinator` role that can register/de-register operators through routed calls to the `AVSDirectory` (i.e, `RegistryCoordinator` -> `EigenDAServiceManager` -> `AVSDirectory`)
- `RewardsInitiator` role that can create operator directed and general AVS rewards via routed calls to the `RewardsCoordinator` contract (i.e, `RewardsInitiator` -> `EigenDAServiceManager` -> `RewardsCoordinator`)


**Offchain Usage**

TODO

### [`EigenDAThresholdRegistry`](https://github.com/Layr-Labs/eigenda/blob/98a17e884de40a18ed9744e709ccc109adf273d3/contracts/src/core/EigenDAThresholdRegistry.sol)
**Description**
<!-- TODO: Cleanup this description and better coalesce wrt other contract doc entries -->
![image.png](../../assets/integration/contracts-eigenda.png)

The [EigenDAThresholdRegistry](https://github.com/Layr-Labs/eigenda/blob/c4567f90e835678fae4749f184857dea10ff330c/contracts/src/core/EigenDAThresholdRegistryStorage.sol#L22) contains two sets of protocol parameters:

```solidity

/// @notice mapping of blob version id to the params of the blob version
mapping(uint16 => VersionedBlobParams) public versionedBlobParams;
struct VersionedBlobParams {
    uint32 maxNumOperators;
    uint32 numChunks;
    uint8 codingRate;
}

/// @notice Immutable security thresholds for quorums
SecurityThresholds public defaultSecurityThresholdsV2;
struct SecurityThresholds {
    uint8 confirmationThreshold;
    uint8 adversaryThreshold;
}
```

The securityThresholds are currently immutable. Confirmation and adversary thresholds are sometimes also [referred to](https://docs.eigenda.xyz/overview#optimal-da-sharding) as liveness and safety thresholds:

- **Confirmation Threshold (aka liveness threshold)**: minimum percentage of stake which an attacker must control in order to mount a liveness attack on the system.
- **Adversary Threshold (aka safety threshold)**: total percentage of stake which an attacker must control in order to mount a first-order safety attack on the system.

Their default values are currently set as:

```solidity
defaultSecurityThresholdsV2 = {
    confirmationThreshold = 55,
    adversaryThreshold = 33,
}
```
A new BlobParam version is rarely introduced by the EigenDA Foundation Governance. When dispersing a blob, rollups explicitly specify the version they wish to use. Currently, only version `0` is defined, with the following parameters ([reference](https://etherscan.io/address/0xdb4c89956eEa6F606135E7d366322F2bDE609F1)):

```solidity
versionedBlobParams[0] = {
    maxNumOperators =  3537,
    numChunks = 8192,
    codingRate = 8,
}
```

The five parameters are intricately related by this formula which is also verified onchain by the [verifyBlobSecurityParams](https://github.com/Layr-Labs/eigenda/blob/77d4442aa1b37bdc275173a6b27d917cc161474c/contracts/src/libraries/EigenDABlobVerificationUtils.sol#L386) function:

$$
numChunks \cdot (1 - \frac{100}{\gamma * codingRate}) \geq maxNumOperators
$$

where $\gamma = confirmationThreshold - adversaryThreshold$

### [`EigenDARelayRegistry`](https://github.com/Layr-Labs/eigenda/blob/98a17e884de40a18ed9744e709ccc109adf273d3/contracts/src/core/EigenDARelayRegistry.sol)

**Description**

Contains EigenDA network registered Relays' Ethereum address and DNS hostname or IP address. `BlobCertificates` contain `relayKeys`, which can be transformed into that relay's URL by calling [relayKeyToUrl](https://github.com/Layr-Labs/eigenda/blob/77d4442aa1b37bdc275173a6b27d917cc161474c/contracts/src/core/EigenDARelayRegistry.sol#L35).

**Access Mgmt**

- `Ownable` role that can register new relay entries

**Offchain Usage**

TODO

### [`EigenDADisperserRegistry`](https://github.com/Layr-Labs/eigenda/blob/98a17e884de40a18ed9744e709ccc109adf273d3/contracts/src/core/EigenDADisperserRegistry.sol)

**Description**

Contains EigenDA network registered Dispersers' Ethereum address. The EigenDA Network currently only supports a single Disperser, hosted by EigenLabs. The Disperser's URL is currently static and unchanging, and can be found on our docs site in the [Networks](https://docs.eigenda.xyz/networks/mainnet) section.

**Access Mgmt**

- `Ownable` role that can register new dispersers

**Offchain Usage**

TODO

### [`PaymentVault`](https://github.com/Layr-Labs/eigenda/blob/98a17e884de40a18ed9744e709ccc109adf273d3/contracts/src/core/PaymentVault.sol)
**Description**

Payment contract used to escrow on-demand funds, hold user reservations, and define global payment parameters used by the network (i.e, `globalSymbolsPerPeriod`, `reservationPeriodInterval`, `globalRatePeriodInterval`).

**Access Mgmt**

- `Ownable` role that can set payment reservations

**Offchain Usage**

TODO

### [`PauserRegistry`](https://github.com/Layr-Labs/eigenlayer-contracts/blob/ac57bc1b28c83d9d7143c0da19167c148c3596a3/src/contracts/permissions/PauserRegistry.sol)

**Description**
Manages a stateful mapping of pausers that can be arbitrarily added or revoked. This contract is assumed to be deployed immutably. The pauser mapping is checked by caller:
- Mapping checked as prerequisite for pausing batch confirmation logic in [`EigenDAServiceManager`](#eigendaservicemanager)
- Mapping checked as prerequisite for pausing operator state update logic in [`RegistryCoordinator`](#eigendaregistrycoordinator)

**Access Mgmt**
- `Unpauser` (or admin) role that can set / remove existing pausers

**Offchain Usage**

TODO

### [`BLSApkRegistry`](https://github.com/Layr-Labs/eigenlayer-middleware/blob/2f7c93e38f56f292f247981a52bd3619a16b9918/src/BLSApkRegistry.sol)

**Description**
This contract stores each operator's BLS public key as well as per quorum aggregate public keys which are only updatable by the `RegistryCoordinator`.

**Access Mgmt**
- `RegistryCoordinator` role that can invoke aggregate key updates via the registration/de-registration of operators


**Offchain Usage**

TODO

### [`IndexRegistry`](https://github.com/Layr-Labs/eigenlayer-middleware/blob/2f7c93e38f56f292f247981a52bd3619a16b9918/src/IndexRegistry.sol)

**Description**
Maintains an ordered, historically versioned list of operators for each quorum, allowing the RegistryCoordinator to register or deregister operators while preserving full block-by-block history of operator counts and index assignments. It provides efficient read functions to reconstruct the operator set at any block

**Access Mgmt**
- `RegistryCoordinator` role that makes stateful updates when registering / deregistering quorum operators

**Offchain Usage**

TODO

### [`StakeRegistry`](https://github.com/Layr-Labs/eigenlayer-middleware/blob/2f7c93e38f56f292f247981a52bd3619a16b9918/src/StakeRegistry.sol)

**Description**
Stores stake updates bounded by block number and quorum strategy:
```solidity
    struct StakeUpdate {
        // the block number at which the stake amounts were updated and stored
        uint32 updateBlockNumber;
        // the block number at which the *next update* occurred.
        /// @notice This entry has the value **0** until another update takes place.
        uint32 nextUpdateBlockNumber;
        // stake weight for the quorum
        uint96 stake;
    }
```

**Access Mgmt**
- `Ownable` role that can deploy and modify staking strategies
- `RegistryCoordinator` role that makes stateful updates when registering / deregistering quorum operators


**Offchain Usage**

TODO

### [`SocketRegistry`](https://github.com/Layr-Labs/eigenlayer-middleware/blob/2f7c93e38f56f292f247981a52bd3619a16b9918/src/SocketRegistry.sol)

**Description**
Stores stateful mapping of `operator ID => socket` where socket is the operator's DNS hostname.

**Access Mgmt**
- `RegistryCoordinator` role that makes stateful updates when registering / deregistering quorum operators

**Offchain Usage**

TODO

### [`OperatorStateRetriever`](https://github.com/Layr-Labs/eigenlayer-middleware/blob/2f7c93e38f56f292f247981a52bd3619a16b9918/src/OperatorStateRetriever.sol)
**Description**

A stateless read-only contract that does exhaustive lookups against the registry coordinator for fetching operator metadata. This bundles stored procedure logic to avoid exhaustive RPC calls made to view functions by offchain EigenDA services.

**Access Mgmt**

N/A

**Offchain Usage**

TODO

### [`EigenDARegistryCoordinator`](https://github.com/Layr-Labs/eigenda/blob/98a17e884de40a18ed9744e709ccc109adf273d3/contracts/src/core/EigenDARegistryCoordinator.sol)

**Description**

This contract orchestrates operator lifecycle across EigenDA's stake, BLS key, index, and socket registries - handling:
- registration, deregistration
- churning
- stake-updates
- quorum creation/config
- historical quorum-bitmap tracking

**Access Mgmt**

- `Pauser` role that can halt operator state updates
- `Ownable` role that can add new quorums, operator set params, & ejector params / role changes
- `Ejector` role that can invoke an ejection function to forcibly deregister an operator


**Offchain Usage**

TODO

### [`EigenDAEjectionsManager`](https://github.com/Layr-Labs/eigenda/blob/98a17e884de40a18ed9744e709ccc109adf273d3/contracts/src/periphery/ejection/EigenDAEjectionManager.sol)

**Description**
Coordinates the lifecycle of ejecting non-responsive operators from EigenDA. It allows an `Ejector` role to queue and complete ejections. Each queued ejection has a corresponding bond attached by the `Ejector` where a targeted operator can cancel the ejection by providing a signature before it becomes "confirmable" after a number of `DelayBlocks`.

**Access Mgmt**
- `Ownable` role that can change public parameters (i.e, `DelayBlocks`, `CooldownBlocks`)
- `Ejector` role that can invoke an ejection function to forcibly deregister an operator

**Offchain Usage**

TODO


## Governance Roles

There are four key governance roles in the EigenDA contracts seen across network environments (i.e, `mainnet`, `hoodi-testnet`, `hoodi-preprod`, `sepolia-testnet`):
- [ERC1967](https://eips.ethereum.org/EIPS/eip-1967) `ProxyAdmin` that can upgrade implementation contracts
- `Owner` that can perform sensitive stateful operations across protocol contracts
- `Pauser` that can halt stateful updates on the `ServiceManager` and `RegistryCoordinator` contracts. This role is managed by the immutable [`PauserRegistry`](#pauserregistry) contract
- `Ejector` that can initialize and complete ejection requests via the [`EjectionsManager`](#eigendaejectionsmanager) contract