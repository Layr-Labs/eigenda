# EigenDAThresholdRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/core/EigenDAThresholdRegistry.sol)

**Inherits:**
[EigenDAThresholdRegistryStorage](/src/core/EigenDAThresholdRegistryStorage.sol/abstract.EigenDAThresholdRegistryStorage.md), OwnableUpgradeable

**Author:**
Layr Labs, Inc.


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


```solidity
function addVersionedBlobParams(VersionedBlobParams memory _versionedBlobParams) external onlyOwner returns (uint16);
```

### _addVersionedBlobParams


```solidity
function _addVersionedBlobParams(VersionedBlobParams memory _versionedBlobParams) internal returns (uint16);
```

### getQuorumAdversaryThresholdPercentage

Gets the adversary threshold percentage for a quorum


```solidity
function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber)
    public
    view
    virtual
    returns (uint8 adversaryThresholdPercentage);
```

### getQuorumConfirmationThresholdPercentage

Gets the confirmation threshold percentage for a quorum


```solidity
function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber)
    public
    view
    virtual
    returns (uint8 confirmationThresholdPercentage);
```

### getIsQuorumRequired

Checks if a quorum is required


```solidity
function getIsQuorumRequired(uint8 quorumNumber) public view virtual returns (bool);
```

### getBlobParams

Returns the blob params for a given blob version


```solidity
function getBlobParams(uint16 version) external view returns (VersionedBlobParams memory);
```

