# BatchHeader
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDAStructs.sol)

The header of a V1 batch


```solidity
struct BatchHeader {
    bytes32 blobHeadersRoot;
    bytes quorumNumbers;
    bytes signedStakeForQuorums;
    uint32 referenceBlockNumber;
}
```

