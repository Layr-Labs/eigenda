# BlobVerificationProof
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDAStructs.sol)

The proof of batch inclusion for a V1 blob


```solidity
struct BlobVerificationProof {
    uint32 batchId;
    uint32 blobIndex;
    BatchMetadata batchMetadata;
    bytes inclusionProof;
    bytes quorumIndices;
}
```

