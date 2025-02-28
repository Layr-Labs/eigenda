# EigenDACertVerificationUtils
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/libraries/EigenDACertVerificationUtils.sol)

**Author:**
Layr Labs, Inc.


## State Variables
### THRESHOLD_DENOMINATOR

```solidity
uint256 public constant THRESHOLD_DENOMINATOR = 100;
```


## Functions
### _verifyDACertV1ForQuorums


```solidity
function _verifyDACertV1ForQuorums(
    IEigenDAThresholdRegistry eigenDAThresholdRegistry,
    IEigenDABatchMetadataStorage batchMetadataStorage,
    BlobHeader calldata blobHeader,
    BlobVerificationProof calldata blobVerificationProof,
    bytes memory requiredQuorumNumbers
) internal view;
```

### _verifyDACertsV1ForQuorums


```solidity
function _verifyDACertsV1ForQuorums(
    IEigenDAThresholdRegistry eigenDAThresholdRegistry,
    IEigenDABatchMetadataStorage batchMetadataStorage,
    BlobHeader[] calldata blobHeaders,
    BlobVerificationProof[] calldata blobVerificationProofs,
    bytes memory requiredQuorumNumbers
) internal view;
```

### _verifyDACertV2ForQuorums


```solidity
function _verifyDACertV2ForQuorums(
    IEigenDAThresholdRegistry eigenDAThresholdRegistry,
    IEigenDASignatureVerifier signatureVerifier,
    IEigenDARelayRegistry eigenDARelayRegistry,
    BatchHeaderV2 memory batchHeader,
    BlobInclusionInfo memory blobInclusionInfo,
    NonSignerStakesAndSignature memory nonSignerStakesAndSignature,
    SecurityThresholds memory securityThresholds,
    bytes memory requiredQuorumNumbers,
    bytes memory signedQuorumNumbers
) internal view;
```

### verifyDACertV2ForQuorumsExternal

*External function needed for try-catch wrapper*


```solidity
function verifyDACertV2ForQuorumsExternal(
    IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
    IEigenDASignatureVerifier _signatureVerifier,
    IEigenDARelayRegistry _eigenDARelayRegistry,
    BatchHeaderV2 memory _batchHeader,
    BlobInclusionInfo memory _blobInclusionInfo,
    NonSignerStakesAndSignature memory _nonSignerStakesAndSignature,
    SecurityThresholds memory _securityThresholds,
    bytes memory _requiredQuorumNumbers,
    bytes memory _signedQuorumNumbers
) external view;
```

### _verifyDACertV2ForQuorumsFromSignedBatch


```solidity
function _verifyDACertV2ForQuorumsFromSignedBatch(
    IEigenDAThresholdRegistry eigenDAThresholdRegistry,
    IEigenDASignatureVerifier signatureVerifier,
    IEigenDARelayRegistry eigenDARelayRegistry,
    OperatorStateRetriever operatorStateRetriever,
    IRegistryCoordinator registryCoordinator,
    SignedBatch memory signedBatch,
    BlobInclusionInfo memory blobInclusionInfo,
    SecurityThresholds memory securityThresholds,
    bytes memory requiredQuorumNumbers
) internal view;
```

### _getNonSignerStakesAndSignature


```solidity
function _getNonSignerStakesAndSignature(
    OperatorStateRetriever operatorStateRetriever,
    IRegistryCoordinator registryCoordinator,
    SignedBatch memory signedBatch
)
    internal
    view
    returns (NonSignerStakesAndSignature memory nonSignerStakesAndSignature, bytes memory signedQuorumNumbers);
```

### _verifyDACertSecurityParams


```solidity
function _verifyDACertSecurityParams(
    VersionedBlobParams memory blobParams,
    SecurityThresholds memory securityThresholds
) internal pure;
```

### _verifyRelayKeysSet


```solidity
function _verifyRelayKeysSet(IEigenDARelayRegistry eigenDARelayRegistry, uint32[] memory relayKeys) internal view;
```

