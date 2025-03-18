# IEigenDADisperserRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDADisperserRegistry.sol)

A registry for EigenDA disperser info


## Functions
### setDisperserInfo

Sets the disperser info for a given disperser key


```solidity
function setDisperserInfo(uint32 _disperserKey, DisperserInfo memory _disperserInfo) external;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`_disperserKey`|`uint32`|The key of the disperser to set the info for|
|`_disperserInfo`|`DisperserInfo`|The info to set for the disperser|


### disperserKeyToAddress

Returns the disperser address for a given disperser key


```solidity
function disperserKeyToAddress(uint32 key) external view returns (address);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`key`|`uint32`|The key of the disperser to get the address for|


## Events
### DisperserAdded
Emitted when a disperser is added to the registry


```solidity
event DisperserAdded(uint32 indexed key, address indexed disperser);
```

