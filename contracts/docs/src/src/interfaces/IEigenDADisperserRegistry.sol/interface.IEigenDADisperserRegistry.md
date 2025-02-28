# IEigenDADisperserRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/interfaces/IEigenDADisperserRegistry.sol)


## Functions
### setDisperserInfo


```solidity
function setDisperserInfo(uint32 _disperserKey, DisperserInfo memory _disperserInfo) external;
```

### disperserKeyToAddress


```solidity
function disperserKeyToAddress(uint32 key) external view returns (address);
```

## Events
### DisperserAdded

```solidity
event DisperserAdded(uint32 indexed key, address indexed disperser);
```

