# EigenDA Managed Contracts

This page describe EigenDA contracts that are manage by EigenDA related actors (see the exact [roles](#governance-roles)). For EigenDA-related contracts that are managed by rollups, see the [rollup managed contracts](../integration/contracts.md) page.

## Middlewares Contracts

## EigenDA Specific Contracts

<!-- Section copied over from https://www.notion.so/eigen-labs/EigenDA-V2-Integration-Spec-12d13c11c3e0800e8968f31ef2c6a2b3?pvs=4#18513c11c3e08058a034ddc9523a3197 -->
<!-- TODO: arch to review and update -->

The smart contracts can be found [here](https://github.com/Layr-Labs/eigenda/tree/master/contracts/src/core).



#### EigenDACertVerifier

Contains a single function verifyDACertV2 which is used to verify `certs`. This function’s logic is described in the [Cert Validation](https://www.notion.so/EigenDA-V2-Integration-Spec-12d13c11c3e0800e8968f31ef2c6a2b3?pvs=21) section.

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

The securityThresholds are currently immutable. These are the same as the [previously called](https://github.com/Layr-Labs/eigenda/blob/master/docs/spec/overview.md#security-model) liveness and safety thresholds:

- Confirmation Threshold (fka liveness threshold): minimum percentage of stake which an attacker must control in order to mount a liveness attack on the system.
- Adversary Threshold (fka safety threshold): total percentage of stake which an attacker must control in order to mount a first-order safety attack on the system.

Their values are

```solidity
defaultSecurityThresholdsV2 = {
	confirmationThreshold = ??,
	adversaryThreshold = ??,
}
```

A new BlobParam version is very infrequently introduced by the EigenDA Foundation Governance, and rollups can choose which version they wish to use when dispersing a blob. Currently there is only version 0 defined, with parameters:

```solidity
versionedBlobParams[0] = {
	maxNumOperators = ??,
	numChunks = 8192,
	codingRate = ??,
}
```

The five parameters are intricately related by this formula which is also verified onchain by the [verifyBlobSecurityParams](https://github.com/Layr-Labs/eigenda/blob/77d4442aa1b37bdc275173a6b27d917cc161474c/contracts/src/libraries/EigenDABlobVerificationUtils.sol#L386) function: 

$$
numChunks \cdot (1 - \frac{100}{\gamma * codingRate}) \geq maxNumOperators
$$

where $\gamma$ = confirmationThreshold - adversaryThreshold

### EigenDARelayRegistry

Contains EigenDA network registered Relays’ Ethereum address and DNS hostname or IP address. `BlobCertificate`s contain `relayKey`(s), which can be transformed into that relay’s URL by calling [relayKeyToUrl](https://github.com/Layr-Labs/eigenda/blob/77d4442aa1b37bdc275173a6b27d917cc161474c/contracts/src/core/EigenDARelayRegistry.sol#L35).

### DisperserRegistry

Contains EigenDA network registered Dispersers’ Ethereum address. The EigenDA Network currently only supports a single Disperser, hosted by EigenLabs. The Disperser’s URL is currently static and unchanging, and can be found on our docs site in the [Networks](https://docs.eigenda.xyz/networks/mainnet) section.

## Deployments

<!-- TODO: add deployed contract addresses table -->

## Governance Roles

<!-- TODO: import from https://www.notion.so/eigen-labs/EigenDA-V2-Governance-17513c11c3e0806999cfe5e8b9bf7e6a -->
<!-- Do we want to make public everything in that doc?? -->
