# NonSignerStakesAndSignature
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDAStructs.sol)

The non-signer stakes and signatures used for BLS signature verification of V1 and V2 batches


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

