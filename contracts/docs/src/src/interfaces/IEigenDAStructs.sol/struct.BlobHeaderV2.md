# BlobHeaderV2
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/interfaces/IEigenDAStructs.sol)


```solidity
struct BlobHeaderV2 {
    uint16 version;
    bytes quorumNumbers;
    BlobCommitment commitment;
    bytes32 paymentHeaderHash;
}
```

