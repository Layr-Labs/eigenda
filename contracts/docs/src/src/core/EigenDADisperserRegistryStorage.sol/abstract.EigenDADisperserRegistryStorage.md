# EigenDADisperserRegistryStorage
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/core/EigenDADisperserRegistryStorage.sol)

**Author:**
Layr Labs, Inc.

This storage contract is separate from the logic to simplify the upgrade process.


## State Variables
### disperserKeyToInfo

```solidity
mapping(uint32 => DisperserInfo) public disperserKeyToInfo;
```


### __GAP

```solidity
uint256[49] private __GAP;
```


