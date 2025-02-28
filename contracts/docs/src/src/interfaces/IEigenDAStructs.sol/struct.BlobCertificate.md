# BlobCertificate
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/interfaces/IEigenDAStructs.sol)


```solidity
struct BlobCertificate {
    BlobHeaderV2 blobHeader;
    bytes signature;
    uint32[] relayKeys;
}
```

