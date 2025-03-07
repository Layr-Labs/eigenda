# EigenDADisperserRegistryStorage
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/core/EigenDADisperserRegistryStorage.sol)

This storage contract is separated from the logic to simplify the upgrade process.


## State Variables
### disperserKeyToInfo
A mapping of disperser keys to disperser info


```solidity
mapping(uint32 => DisperserInfo) public disperserKeyToInfo;
```


### __GAP
Storage gap for upgradeability


```solidity
uint256[49] private __GAP;
```


