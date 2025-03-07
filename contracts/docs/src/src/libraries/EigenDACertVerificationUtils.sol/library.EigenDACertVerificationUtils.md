# EigenDACertVerificationUtils
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/libraries/EigenDACertVerificationUtils.sol)

Library of functions to be used by smart contracts wanting to verify submissions of blob certificates on EigenDA.


## State Variables
### THRESHOLD_DENOMINATOR

```solidity
uint256 public constant THRESHOLD_DENOMINATOR = 100;
```


## Functions
### _verifyDACertV1ForQuorums

Verifies a V1 blob certificate for a set of quorums


```solidity
function _verifyDACertV1ForQuorums(
    IEigenDAThresholdRegistry eigenDAThresholdRegistry,
    IEigenDABatchMetadataStorage batchMetadataStorage,
    BlobHeader calldata blobHeader,
    BlobVerificationProof calldata blobVerificationProof,
    bytes memory requiredQuorumNumbers
) internal view;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`eigenDAThresholdRegistry`|`IEigenDAThresholdRegistry`|is the ThresholdRegistry contract address|
|`batchMetadataStorage`|`IEigenDABatchMetadataStorage`|is the BatchMetadataStorage contract address|
|`blobHeader`|`BlobHeader`|is the blob header to verify|
|`blobVerificationProof`|`BlobVerificationProof`|is the blob verification proof to verify|
|`requiredQuorumNumbers`|`bytes`|is the required quorum numbers to verify against|


### _verifyDACertsV1ForQuorums

Verifies a set of V1 blob certificates for a set of quorums


```solidity
function _verifyDACertsV1ForQuorums(
    IEigenDAThresholdRegistry eigenDAThresholdRegistry,
    IEigenDABatchMetadataStorage batchMetadataStorage,
    BlobHeader[] calldata blobHeaders,
    BlobVerificationProof[] calldata blobVerificationProofs,
    bytes memory requiredQuorumNumbers
) internal view;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`eigenDAThresholdRegistry`|`IEigenDAThresholdRegistry`|is the ThresholdRegistry contract address|
|`batchMetadataStorage`|`IEigenDABatchMetadataStorage`|is the BatchMetadataStorage contract address|
|`blobHeaders`|`BlobHeader[]`|is the set of blob headers to verify|
|`blobVerificationProofs`|`BlobVerificationProof[]`|is the set of blob verification proofs to verify for each blob header|
|`requiredQuorumNumbers`|`bytes`|is the required quorum numbers to verify against|


### _verifyDACertV2ForQuorums

Verifies a V2 blob certificate for a set of quorums


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
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`eigenDAThresholdRegistry`|`IEigenDAThresholdRegistry`|is the ThresholdRegistry contract address|
|`signatureVerifier`|`IEigenDASignatureVerifier`|is the SignatureVerifier contract address|
|`eigenDARelayRegistry`|`IEigenDARelayRegistry`|is the RelayRegistry contract address|
|`batchHeader`|`BatchHeaderV2`|is the batch header to verify|
|`blobInclusionInfo`|`BlobInclusionInfo`|is the blob inclusion proof to verify against the batch|
|`nonSignerStakesAndSignature`|`NonSignerStakesAndSignature`|is the non-signer stakes and signatures to check the signature against|
|`securityThresholds`|`SecurityThresholds`|are the confirmation and adversary thresholds to verify|
|`requiredQuorumNumbers`|`bytes`|is the required quorum numbers to verify against|
|`signedQuorumNumbers`|`bytes`|are the quorum numbers that signed on the batch|


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

Verifies a V2 blob certificate for a set of quorums from a signed batch


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
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`eigenDAThresholdRegistry`|`IEigenDAThresholdRegistry`|is the ThresholdRegistry contract address|
|`signatureVerifier`|`IEigenDASignatureVerifier`|is the SignatureVerifier contract address|
|`eigenDARelayRegistry`|`IEigenDARelayRegistry`|is the RelayRegistry contract address|
|`operatorStateRetriever`|`OperatorStateRetriever`|is the OperatorStateRetriever contract address|
|`registryCoordinator`|`IRegistryCoordinator`|is the RegistryCoordinator contract address|
|`signedBatch`|`SignedBatch`|is the signed batch to verify|
|`blobInclusionInfo`|`BlobInclusionInfo`|is the blob inclusion proof to verify against the batch|
|`securityThresholds`|`SecurityThresholds`|are the confirmation and adversary thresholds to verify|
|`requiredQuorumNumbers`|`bytes`|is the required quorum numbers to verify against|


### _getNonSignerStakesAndSignature

*Internal function to get the non-signer stakes and signature from the Attestation of a signed batch*


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

*Internal function to verify the security parameters of a V2 blob certificate*


```solidity
function _verifyDACertSecurityParams(
    VersionedBlobParams memory blobParams,
    SecurityThresholds memory securityThresholds
) internal pure;
```

### _verifyRelayKeysSet

*Internal function to verify that the provided relay keys are set on the RelayRegistry*


```solidity
function _verifyRelayKeysSet(IEigenDARelayRegistry eigenDARelayRegistry, uint32[] memory relayKeys) internal view;
```

