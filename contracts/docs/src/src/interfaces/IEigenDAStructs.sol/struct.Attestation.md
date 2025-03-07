# Attestation
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDAStructs.sol)

The attestation for a V2 batch


```solidity
struct Attestation {
    BN254.G1Point[] nonSignerPubkeys;
    BN254.G1Point[] quorumApks;
    BN254.G1Point sigma;
    BN254.G2Point apkG2;
    uint32[] quorumNumbers;
}
```

