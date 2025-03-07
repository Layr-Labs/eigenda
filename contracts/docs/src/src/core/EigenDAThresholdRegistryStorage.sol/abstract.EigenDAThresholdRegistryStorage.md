# EigenDAThresholdRegistryStorage
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/core/EigenDAThresholdRegistryStorage.sol)

**Inherits:**
[IEigenDAThresholdRegistry](/src/interfaces/IEigenDAThresholdRegistry.sol/interface.IEigenDAThresholdRegistry.md)

This storage contract is separated from the logic to simplify the upgrade process.


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
Storage gap for upgradeability


```solidity
uint256[45] private __GAP;
```


