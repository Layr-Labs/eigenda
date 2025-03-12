# EigenDAThresholdRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/core/EigenDAThresholdRegistry.sol)

**Inherits:**
[EigenDAThresholdRegistryStorage](/src/core/EigenDAThresholdRegistryStorage.sol/abstract.EigenDAThresholdRegistryStorage.md), OwnableUpgradeable

This contract is used for storing:
- The threshold percentages used for V1 certificate verification
- The parameters for the blob versions used for V2 certificate verification


## Functions
### constructor


```solidity
constructor();
```

### initialize


```solidity
function initialize(
    address _initialOwner,
    bytes memory _quorumAdversaryThresholdPercentages,
    bytes memory _quorumConfirmationThresholdPercentages,
    bytes memory _quorumNumbersRequired,
    VersionedBlobParams[] memory _versionedBlobParams
) external initializer;
```

### addVersionedBlobParams

Appends a new blob version to the registry

*This function is append only and cannot be used to update existing blob versions*


```solidity
function addVersionedBlobParams(VersionedBlobParams memory _versionedBlobParams) external onlyOwner returns (uint16);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_versionedBlobParams`|`VersionedBlobParams`|The blob version parameters to add|


### _addVersionedBlobParams

Internal function to append a new blob version to the registry


```solidity
function _addVersionedBlobParams(VersionedBlobParams memory _versionedBlobParams) internal returns (uint16);
```

### getQuorumAdversaryThresholdPercentage

Returns the adversary threshold percentage for a quorum for V1 verification


```solidity
function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber)
    public
    view
    virtual
    returns (uint8 adversaryThresholdPercentage);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`quorumNumber`|`uint8`|The number of the quorum to get the adversary threshold percentage for|


### getQuorumConfirmationThresholdPercentage

Returns the confirmation threshold percentage for a quorum for V1 verification


```solidity
function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber)
    public
    view
    virtual
    returns (uint8 confirmationThresholdPercentage);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`quorumNumber`|`uint8`|The number of the quorum to get the confirmation threshold percentage for|


### getIsQuorumRequired

Returns true if a quorum is required for V1 verification


```solidity
function getIsQuorumRequired(uint8 quorumNumber) public view virtual returns (bool);
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


