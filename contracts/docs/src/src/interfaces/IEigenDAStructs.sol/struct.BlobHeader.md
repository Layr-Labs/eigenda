# BlobHeader
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/interfaces/IEigenDAStructs.sol)


```solidity
struct BlobHeader {
    BN254.G1Point commitment;
    uint32 dataLength;
    QuorumBlobParam[] quorumBlobParams;
}
```

