# EigenDARelayRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/core/EigenDARelayRegistry.sol)

**Inherits:**
OwnableUpgradeable, [EigenDARelayRegistryStorage](/src/core/EigenDARelayRegistryStorage.sol/abstract.EigenDARelayRegistryStorage.md), [IEigenDARelayRegistry](/src/interfaces/IEigenDARelayRegistry.sol/interface.IEigenDARelayRegistry.md)

**Author:**
Layr Labs, Inc.


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


```solidity
function addRelayInfo(RelayInfo memory relayInfo) external onlyOwner returns (uint32);
```

### relayKeyToAddress


```solidity
function relayKeyToAddress(uint32 key) external view returns (address);
```

### relayKeyToUrl


```solidity
function relayKeyToUrl(uint32 key) external view returns (string memory);
```

