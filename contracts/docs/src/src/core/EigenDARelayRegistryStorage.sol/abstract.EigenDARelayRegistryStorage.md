# EigenDARelayRegistryStorage
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/core/EigenDARelayRegistryStorage.sol)

This storage contract is separated from the logic to simplify the upgrade process.


## State Variables
### relayKeyToInfo
A mapping of relay keys to relay info


```solidity
mapping(uint32 => RelayInfo) public relayKeyToInfo;
```


### nextRelayKey
The next relay key to be used


```solidity
uint32 public nextRelayKey;
```


### __GAP
Storage gap for upgradeability


```solidity
uint256[48] private __GAP;
```


