# EigenDADisperserRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/core/EigenDADisperserRegistry.sol)

**Inherits:**
OwnableUpgradeable, [EigenDADisperserRegistryStorage](/src/core/EigenDADisperserRegistryStorage.sol/abstract.EigenDADisperserRegistryStorage.md), [IEigenDADisperserRegistry](/src/interfaces/IEigenDADisperserRegistry.sol/interface.IEigenDADisperserRegistry.md)

A registry for EigenDA disperser info


## Functions
### constructor


```solidity
constructor();
```

### initialize


```solidity
function initialize(address _initialOwner) external initializer;
```

### setDisperserInfo

Sets the disperser info for a given disperser key


```solidity
function setDisperserInfo(uint32 _disperserKey, DisperserInfo memory _disperserInfo) external onlyOwner;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_disperserKey`|`uint32`|The key of the disperser to set the info for|
|`_disperserInfo`|`DisperserInfo`|The info to set for the disperser|


### disperserKeyToAddress

Returns the disperser address for a given disperser key


```solidity
function disperserKeyToAddress(uint32 _key) external view returns (address);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_key`|`uint32`|The key of the disperser to get the address for|


