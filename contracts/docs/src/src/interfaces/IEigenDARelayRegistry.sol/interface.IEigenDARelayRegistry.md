# IEigenDARelayRegistry
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/interfaces/IEigenDARelayRegistry.sol)


## Functions
### addRelayInfo


```solidity
function addRelayInfo(RelayInfo memory relayInfo) external returns (uint32);
```

### relayKeyToAddress


```solidity
function relayKeyToAddress(uint32 key) external view returns (address);
```

### relayKeyToUrl


```solidity
function relayKeyToUrl(uint32 key) external view returns (string memory);
```

## Events
### RelayAdded

```solidity
event RelayAdded(address indexed relay, uint32 indexed key, string relayURL);
```

