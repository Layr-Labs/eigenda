# EigenDADisperserRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/core/EigenDADisperserRegistry.sol)

**Inherits:**
OwnableUpgradeable, [EigenDADisperserRegistryStorage](/src/core/EigenDADisperserRegistryStorage.sol/abstract.EigenDADisperserRegistryStorage.md), [IEigenDADisperserRegistry](/src/interfaces/IEigenDADisperserRegistry.sol/interface.IEigenDADisperserRegistry.md)

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

### setDisperserInfo


```solidity
function setDisperserInfo(uint32 _disperserKey, DisperserInfo memory _disperserInfo) external onlyOwner;
```

### disperserKeyToAddress


```solidity
function disperserKeyToAddress(uint32 _key) external view returns (address);
```

