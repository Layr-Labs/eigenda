# IEigenDARelayRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDARelayRegistry.sol)

A registry for EigenDA relay info


## Functions
### addRelayInfo

Appends a relay info to the registry and returns the relay key


```solidity
function addRelayInfo(RelayInfo memory relayInfo) external returns (uint32);
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


## Events
### RelayAdded
Emitted when a relay is added to the registry


```solidity
event RelayAdded(address indexed relay, uint32 indexed key, string relayURL);
```

