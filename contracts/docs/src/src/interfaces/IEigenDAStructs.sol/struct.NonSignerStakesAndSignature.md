# NonSignerStakesAndSignature
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/interfaces/IEigenDAStructs.sol)


```solidity
struct NonSignerStakesAndSignature {
    uint32[] nonSignerQuorumBitmapIndices;
    BN254.G1Point[] nonSignerPubkeys;
    BN254.G1Point[] quorumApks;
    BN254.G2Point apkG2;
    BN254.G1Point sigma;
    uint32[] quorumApkIndices;
    uint32[] totalStakeIndices;
    uint32[][] nonSignerStakeIndices;
}
```

