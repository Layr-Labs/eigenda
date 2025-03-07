# EigenDACertVerifier
[Git Source](https://github.com/Layr-Labs/eigenda/blob/f0d0dc5708f7e00684e5f5d89ab0227171768419/src/core/EigenDACertVerifier.sol)

**Inherits:**
[IEigenDACertVerifier](/src/interfaces/IEigenDACertVerifier.sol/interface.IEigenDACertVerifier.md)

A CertVerifier is an immutable contract that is used by a consumer to verify EigenDA blob certificates

For V2 verification this contract is deployed with immutable security thresholds and required quorum numbers,
to change these values or verification behavior a new CertVerifier must be deployed


## State Variables
### eigenDAThresholdRegistry
The EigenDAThresholdRegistry contract address


```solidity
IEigenDAThresholdRegistry public immutable eigenDAThresholdRegistry;
```


### eigenDABatchMetadataStorage
The EigenDABatchMetadataStorage contract address

*On L1 this contract is the EigenDA Service Manager contract*


```solidity
IEigenDABatchMetadataStorage public immutable eigenDABatchMetadataStorage;
```


### eigenDASignatureVerifier
The EigenDASignatureVerifier contract address

*On L1 this contract is the EigenDA Service Manager contract*


```solidity
IEigenDASignatureVerifier public immutable eigenDASignatureVerifier;
```


### eigenDARelayRegistry
The EigenDARelayRegistry contract address


```solidity
IEigenDARelayRegistry public immutable eigenDARelayRegistry;
```


### operatorStateRetriever
The EigenDA middleware OperatorStateRetriever contract address


```solidity
OperatorStateRetriever public immutable operatorStateRetriever;
```


### registryCoordinator
The EigenDA middleware RegistryCoordinator contract address


```solidity
IRegistryCoordinator public immutable registryCoordinator;
```


### securityThresholdsV2
The security thresholds checked against for V2 certificates


```solidity
SecurityThresholds public securityThresholdsV2;
```


### quorumNumbersRequiredV2
The quorum numbers required for V2 certificates


```solidity
bytes public quorumNumbersRequiredV2;
```


## Functions
### constructor


```solidity
constructor(
    IEigenDAThresholdRegistry _eigenDAThresholdRegistry,
    IEigenDABatchMetadataStorage _eigenDABatchMetadataStorage,
    IEigenDASignatureVerifier _eigenDASignatureVerifier,
    IEigenDARelayRegistry _eigenDARelayRegistry,
    OperatorStateRetriever _operatorStateRetriever,
    IRegistryCoordinator _registryCoordinator,
    SecurityThresholds memory _securityThresholdsV2,
    bytes memory _quorumNumbersRequiredV2
);
```

### verifyDACertV1

Verifies that a the blob cert is valid for the required quorums


```solidity
function verifyDACertV1(BlobHeader calldata blobHeader, BlobVerificationProof calldata blobVerificationProof)
    external
    view;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`blobHeader`|`BlobHeader`|The blob header to verify|
|`blobVerificationProof`|`BlobVerificationProof`|The blob cert verification proof to verify|


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

Verifies a blob cert using the immutable required quorums and security thresholds set in the constructor


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

Verifies a blob cert using the immutable required quorums and security thresholds set in the constructor


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

Returns the nonSignerStakesAndSignature for a given blob cert and signed batch for V2 verification


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

Verifies the security parameters for a blob cert for V2 verification


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

Verifies the security parameters for a blob cert for V2 verification


```solidity
function verifyDACertSecurityParams(uint16 version, SecurityThresholds memory securityThresholds) external view;
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`version`|`uint16`|The version of the blob to verify|
|`securityThresholds`|`SecurityThresholds`|The security thresholds to verify against|


### quorumAdversaryThresholdPercentages

Returns an array of bytes where each byte represents the adversary threshold percentage of the quorum at that index for V1 verification


```solidity
function quorumAdversaryThresholdPercentages() external view returns (bytes memory);
```

### quorumConfirmationThresholdPercentages

Returns an array of bytes where each byte represents the confirmation threshold percentage of the quorum at that index for V1 verification


```solidity
function quorumConfirmationThresholdPercentages() external view returns (bytes memory);
```

### quorumNumbersRequired

Returns an array of bytes where each byte represents the number of a required quorum for V1 verification


```solidity
function quorumNumbersRequired() public view returns (bytes memory);
```

### getQuorumAdversaryThresholdPercentage

Returns the adversary threshold percentage for a quorum for V1 verification


```solidity
function getQuorumAdversaryThresholdPercentage(uint8 quorumNumber) external view returns (uint8);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`quorumNumber`|`uint8`|The number of the quorum to get the adversary threshold percentage for|


### getQuorumConfirmationThresholdPercentage

Returns the confirmation threshold percentage for a quorum for V1 verification


```solidity
function getQuorumConfirmationThresholdPercentage(uint8 quorumNumber) external view returns (uint8);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`quorumNumber`|`uint8`|The number of the quorum to get the confirmation threshold percentage for|


### getIsQuorumRequired

Returns true if a quorum is required for V1 verification


```solidity
function getIsQuorumRequired(uint8 quorumNumber) external view returns (bool);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`quorumNumber`|`uint8`|The number of the quorum to check if it is required for V1 verification|


### getBlobParams

Returns the blob params for a given blob version


```solidity
function getBlobParams(uint16 version) public view returns (VersionedBlobParams memory);
```
**Parameters**

|Name|Type|Description|
|----|----|-----------|
|`version`|`uint16`|The version of the blob to get the params for|


