# EigenDAThresholdRegistryStorage
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/core/EigenDAThresholdRegistryStorage.sol)

**Inherits:**
[IEigenDAThresholdRegistry](/src/interfaces/IEigenDAThresholdRegistry.sol/interface.IEigenDAThresholdRegistry.md)

**Author:**
Layr Labs, Inc.

This storage contract is separate from the logic to simplify the upgrade process.


## State Variables
### quorumAdversaryThresholdPercentages
The adversary threshold percentage for the quorum at position `quorumNumber`


```solidity
bytes public quorumAdversaryThresholdPercentages;
```


### quorumConfirmationThresholdPercentages
The confirmation threshold percentage for the quorum at position `quorumNumber`


```solidity
bytes public quorumConfirmationThresholdPercentages;
```


### quorumNumbersRequired
The set of quorum numbers that are required


```solidity
bytes public quorumNumbersRequired;
```


### nextBlobVersion
The next blob version id to be added


```solidity
uint16 public nextBlobVersion;
```


### versionedBlobParams
mapping of blob version id to the params of the blob version


```solidity
mapping(uint16 => VersionedBlobParams) public versionedBlobParams;
```


### __GAP

```solidity
uint256[45] private __GAP;
```


