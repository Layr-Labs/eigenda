# EigenDARelayRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/core/EigenDARelayRegistry.sol)

**Inherits:**
OwnableUpgradeable, [EigenDARelayRegistryStorage](/src/core/EigenDARelayRegistryStorage.sol/abstract.EigenDARelayRegistryStorage.md), [IEigenDARelayRegistry](/src/interfaces/IEigenDARelayRegistry.sol/interface.IEigenDARelayRegistry.md)

A registry for EigenDA relay info

*This contract is append only and does not support updating or removing relay info*


## Functions
### constructor


```solidity
constructor();
```

### initialize


```solidity
function initialize(address _initialOwner) external initializer;
```

### addRelayInfo

Appends a relay info to the registry and returns the relay key


```solidity
function addRelayInfo(RelayInfo memory relayInfo) external onlyOwner returns (uint32);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`relayInfo`|`RelayInfo`|The relay info to add|


### relayKeyToAddress

Returns the relay address for a given relay key


```solidity
function relayKeyToAddress(uint32 key) external view returns (address);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`key`|`uint32`|The key of the relay to get the address for|


### relayKeyToUrl

Returns the relay URL for a given relay key


```solidity
function relayKeyToUrl(uint32 key) external view returns (string memory);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`key`|`uint32`|The key of the relay to get the URL for|


