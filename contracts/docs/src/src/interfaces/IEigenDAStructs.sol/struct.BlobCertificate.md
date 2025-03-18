# BlobCertificate
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/interfaces/IEigenDAStructs.sol)

The certificate for a V2 blob


```solidity
struct BlobCertificate {
    BlobHeaderV2 blobHeader;
    bytes signature;
    uint32[] relayKeys;
}
```

