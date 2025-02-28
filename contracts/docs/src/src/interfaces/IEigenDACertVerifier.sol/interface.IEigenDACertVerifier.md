# IEigenDACertVerifier
[Git Source](https://github.com/Layr-Labs/eigenda/blob/538f0525d9ff112a8ba32701edaf2860a0ad7306/src/interfaces/IEigenDACertVerifier.sol)

**Inherits:**
[IEigenDAThresholdRegistry](/src/interfaces/IEigenDAThresholdRegistry.sol/interface.IEigenDAThresholdRegistry.md)


## Functions
### verifyDACertV1

Verifies a the blob cert is valid for the required quorums


```solidity
function verifyDACertV1(BlobHeader calldata blobHeader, BlobVerificationProof calldata blobVerificationProof)
    external
    view;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`blobHeader`|`BlobHeader`|The blob header to verify|
|`blobVerificationProof`|`BlobVerificationProof`|The blob cert verification proof to verify against|


### verifyDACertsV1

Verifies a batch of blob certs for the required quorums


```solidity
function verifyDACertsV1(BlobHeader[] calldata blobHeaders, BlobVerificationProof[] calldata blobVerificationProofs)
    external
    view;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`blobHeaders`|`BlobHeader[]`|The blob headers to verify|
|`blobVerificationProofs`|`BlobVerificationProof[]`|The blob cert verification proofs to verify against|


### verifyDACertV2

Verifies a blob cert for the specified quorums with the default security thresholds


```solidity
function verifyDACertV2(
    BatchHeaderV2 calldata batchHeader,
    BlobInclusionInfo calldata blobInclusionInfo,
    NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
    bytes memory signedQuorumNumbers
) external view;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchHeader`|`BatchHeaderV2`|The batch header of the blob|
|`blobInclusionInfo`|`BlobInclusionInfo`|The inclusion proof for the blob cert|
|`nonSignerStakesAndSignature`|`NonSignerStakesAndSignature`|The nonSignerStakesAndSignature to verify the blob cert against|
|`signedQuorumNumbers`|`bytes`|The signed quorum numbers corresponding to the nonSignerStakesAndSignature|


### verifyDACertV2FromSignedBatch

Verifies a blob cert for the specified quorums with the default security thresholds


```solidity
function verifyDACertV2FromSignedBatch(SignedBatch calldata signedBatch, BlobInclusionInfo calldata blobInclusionInfo)
    external
    view;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`signedBatch`|`SignedBatch`|The signed batch to verify the blob cert against|
|`blobInclusionInfo`|`BlobInclusionInfo`|The inclusion proof for the blob cert|


### verifyDACertV2ForZKProof

Thin try/catch wrapper around verifyDACertV2 that returns false instead of panicing

*The Steel library (https://github.com/risc0/risc0-ethereum/tree/main/crates/steel)
currently has a limitation that it can only create zk proofs for functions that return a value*


```solidity
function verifyDACertV2ForZKProof(
    BatchHeaderV2 calldata batchHeader,
    BlobInclusionInfo calldata blobInclusionInfo,
    NonSignerStakesAndSignature calldata nonSignerStakesAndSignature,
    bytes memory signedQuorumNumbers
) external view returns (bool);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`batchHeader`|`BatchHeaderV2`|The batch header of the blob|
|`blobInclusionInfo`|`BlobInclusionInfo`|The inclusion proof for the blob cert|
|`nonSignerStakesAndSignature`|`NonSignerStakesAndSignature`|The nonSignerStakesAndSignature to verify the blob cert against|
|`signedQuorumNumbers`|`bytes`|The signed quorum numbers corresponding to the nonSignerStakesAndSignature|


### getNonSignerStakesAndSignature

Returns the nonSignerStakesAndSignature for a given blob cert and signed batch


```solidity
function getNonSignerStakesAndSignature(SignedBatch calldata signedBatch)
    external
    view
    returns (NonSignerStakesAndSignature memory);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`signedBatch`|`SignedBatch`|The signed batch to get the nonSignerStakesAndSignature for|

**Returns**

|Name|Type|Description|
|----|----|-----------|
|`<none>`|`NonSignerStakesAndSignature`|nonSignerStakesAndSignature The nonSignerStakesAndSignature for the given signed batch attestation|


### verifyDACertSecurityParams

Verifies the security parameters for a blob cert


```solidity
function verifyDACertSecurityParams(VersionedBlobParams memory blobParams, SecurityThresholds memory securityThresholds)
    external
    view;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`blobParams`|`VersionedBlobParams`|The blob params to verify|
|`securityThresholds`|`SecurityThresholds`|The security thresholds to verify against|


### verifyDACertSecurityParams

Verifies the security parameters for a blob cert


```solidity
function verifyDACertSecurityParams(uint16 version, SecurityThresholds memory securityThresholds) external view;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`version`|`uint16`|The version of the blob to verify|
|`securityThresholds`|`SecurityThresholds`|The security thresholds to verify against|


