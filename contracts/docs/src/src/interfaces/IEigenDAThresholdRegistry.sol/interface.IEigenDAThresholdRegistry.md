# IEigenDAThresholdRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDAThresholdRegistry.sol)


## Functions
### quorumAdversaryThresholdPercentages

Returns an array of bytes where each byte represents the adversary threshold percentage of the quorum at that index


```solidity
function quorumAdversaryThresholdPercentages() external view returns (bytes memory);
```

### quorumConfirmationThresholdPercentages

Returns an array of bytes where each byte represents the confirmation threshold percentage of the quorum at that index


```solidity
function quorumConfirmationThresholdPercentages() external view returns (bytes memory);
```

### quorumNumbersRequired

Returns an array of bytes where each byte represents the number of a required quorum


```solidity
function quorumNumbersRequired() external view returns (bytes memory);
```

### getQuorumAdversaryThresholdPercentage

Returns the adversary threshold percentage for a quorum for V1 verification


```solidity
function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) external view returns (uint8);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`quorumNumber`|`uint8`|The number of the quorum to get the adversary threshold percentage for|


### getQuorumConfirmationThresholdPercentage

Returns the confirmation threshold percentage for a quorum for V1 verification


```solidity
function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) external view returns (uint8);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`quorumNumber`|`uint8`|The number of the quorum to get the confirmation threshold percentage for|


### getIsQuorumRequired

Returns true if a quorum is required for V1 verification


```solidity
function getIsQuorumRequired(uint8 quorumNumber) external view returns (bool);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`quorumNumber`|`uint8`|The number of the quorum to check if it is required for V1 verification|


### getBlobParams

Returns the blob params for a given blob version


```solidity
function getBlobParams(uint16 version) external view returns (VersionedBlobParams memory);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`version`|`uint16`|The version of the blob to get the params for|


## Events
### VersionedBlobParamsAdded
Emitted when a new blob version is added to the registry


```solidity
event VersionedBlobParamsAdded(uint16 indexed version, VersionedBlobParams versionedBlobParams);
```

