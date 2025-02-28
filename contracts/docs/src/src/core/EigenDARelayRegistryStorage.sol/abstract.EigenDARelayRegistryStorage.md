# EigenDARelayRegistryStorage
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/core/EigenDARelayRegistryStorage.sol)

**Author:**
Layr Labs, Inc.

This storage contract is separate from the logic to simplify the upgrade process.


## State Variables
### relayKeyToInfo

```solidity
mapping(uint32 => RelayInfo) public relayKeyToInfo;
```


### nextRelayKey

```solidity
uint32 public nextRelayKey;
```


### __GAP

```solidity
uint256[48] private __GAP;
```


