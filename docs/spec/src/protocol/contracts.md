# EigenDA Managed Contracts

This page describes EigenDA contracts that are managed by EigenDA related actors (see the exact [roles](#governance-roles)). For EigenDA-related contracts that are managed by rollups, see the [rollup managed contracts](../integration/spec/4-contracts.md) page.

> Warning: This page is incomplete and a work in progress as we are undergoing refactors of our contracts as well as some protocol upgrades. The details will change, but the information contained here should at least help to understand the important concepts.

## Middlewares Contracts

We make use of eigenlayer-middleware contracts, which are fully documented [here](https://github.com/Layr-Labs/eigenlayer-middleware/tree/dev/docs).

## EigenDA Specific Contracts

<!-- Section copied over from https://www.notion.so/eigen-labs/EigenDA-V2-Integration-Spec-12d13c11c3e0800e8968f31ef2c6a2b3?pvs=4#18513c11c3e08058a034ddc9523a3197 -->
<!-- TODO: arch to review and update -->

The smart contracts can be found in our [repo](https://github.com/Layr-Labs/eigenda/tree/master/contracts/src/core), and the deployment addresses on different chains can be found in the [Networks](https://docs.eigenda.xyz/networks/mainnet#contract-addresses) section of our docs.

![image.png](../../assets/integration/contracts-eigenda.png)

### EigenDAThreshold Registry

The [EigenDAThresholdRegistry](https://github.com/Layr-Labs/eigenda/blob/c4567f90e835678fae4749f184857dea10ff330c/contracts/src/core/EigenDAThresholdRegistryStorage.sol#L22) contains two sets of fundamental parameters:

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
A new BlobParam version is rarely introduced by the EigenDA Foundation Governance. When dispersing a blob, rollups explicitly specify the version they wish to use. Currently, only version `0` is defined, with the following parameters ((reference)[https://etherscan.io/address/0xdb4c89956eEa6F606135E7d366322F2bDE609F1]):

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

### EigenDARelayRegistry

Contains EigenDA network registered Relays’ Ethereum address and DNS hostname or IP address. `BlobCertificates` contain `relayKeys`, which can be transformed into that relay’s URL by calling [relayKeyToUrl](https://github.com/Layr-Labs/eigenda/blob/77d4442aa1b37bdc275173a6b27d917cc161474c/contracts/src/core/EigenDARelayRegistry.sol#L35).

### EigenDADisperserRegistry

Contains EigenDA network registered Dispersers’ Ethereum address. The EigenDA Network currently only supports a single Disperser, hosted by EigenLabs. The Disperser’s URL is currently static and unchanging, and can be found on our docs site in the [Networks](https://docs.eigenda.xyz/networks/mainnet) section.

## Governance Roles

<!-- TODO: import from https://www.notion.so/eigen-labs/EigenDA-V2-Governance-17513c11c3e0806999cfe5e8b9bf7e6a -->
<!-- Do we want to make public everything in that doc?? -->

TODO