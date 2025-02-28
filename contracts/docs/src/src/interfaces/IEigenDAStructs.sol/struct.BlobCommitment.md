# BlobCommitment
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/interfaces/IEigenDAStructs.sol)


```solidity
struct BlobCommitment {
    BN254.G1Point commitment;
    BN254.G2Point lengthCommitment;
    BN254.G2Point lengthProof;
    uint32 length;
}
```

