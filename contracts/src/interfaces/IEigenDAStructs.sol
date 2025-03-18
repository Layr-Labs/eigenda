// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {BN254} from "../../lib/eigenlayer-middleware/src/libraries/BN254.sol";



///////////////////////// V1 ///////////////////////////////

/**
 * @title QuorumBlobParam
 * @notice The parameters for a quorum in V1 verification
 */
struct QuorumBlobParam {
    uint8 quorumNumber;
    uint8 adversaryThresholdPercentage;
    uint8 confirmationThresholdPercentage; 
    uint32 chunkLength; 
}

/**
 * @title BlobHeader
 * @notice The header of a V1 blob
 */
struct BlobHeader {
    BN254.G1Point commitment; 
    uint32 dataLength; 
    QuorumBlobParam[] quorumBlobParams; 
}

/**
 * @title ReducedBatchHeader
 * @notice The reduced header of a V1 batch
 */
struct ReducedBatchHeader {
    bytes32 blobHeadersRoot;
    uint32 referenceBlockNumber;
}

/**
 * @title BatchHeader
 * @notice The header of a V1 batch
 */
struct BatchHeader {
    bytes32 blobHeadersRoot;
    bytes quorumNumbers; 
    bytes signedStakeForQuorums; 
    uint32 referenceBlockNumber;
}
    
/**
 * @title BatchMetadata
 * @notice The metadata of a V1 batch
 */
struct BatchMetadata {
    BatchHeader batchHeader; 
    bytes32 signatoryRecordHash; 
    uint32 confirmationBlockNumber; 
}

/**
 * @title BlobVerificationProof
 * @notice The proof of batch inclusion for a V1 blob
 */
struct BlobVerificationProof {
    uint32 batchId;
    uint32 blobIndex;
    BatchMetadata batchMetadata;
    bytes inclusionProof;
    bytes quorumIndices;
}

///////////////////////// V2 ///////////////////////////////

/**
 * @title VersionedBlobParams
 * @notice The parameters for a V2 blob version
 */
struct VersionedBlobParams {
    uint32 maxNumOperators;
    uint32 numChunks;
    uint8 codingRate;
}

/**
 * @title SecurityThresholds
 * @notice The security thresholds for a V2 blob version
 */
struct SecurityThresholds {
    uint8 confirmationThreshold;
    uint8 adversaryThreshold;
}

/**
 * @title RelayInfo
 * @notice The info for a V2 relay
 */
struct RelayInfo {
    address relayAddress;
    string relayURL;
}

/**
 * @title DisperserInfo
 * @notice The info for a V2 disperser
 */
struct DisperserInfo {
    address disperserAddress;
}

/**
 * @title BlobInclusionInfo
 * @notice The inclusion proof for a V2 blob 
 */
struct BlobInclusionInfo {
    BlobCertificate blobCertificate;
    uint32 blobIndex;
    bytes inclusionProof;
}

/**
 * @title BlobCertificate
 * @notice The certificate for a V2 blob
 */
struct BlobCertificate {
    BlobHeaderV2 blobHeader;
    bytes signature;
    uint32[] relayKeys;
}

/**
 * @title BlobHeaderV2
 * @notice The header of a V2 blob
 */
struct BlobHeaderV2 {
    uint16 version;
    bytes quorumNumbers;
    BlobCommitment commitment;
    bytes32 paymentHeaderHash;
}

/**
 * @title BlobCommitment
 * @notice The KZG commitment for a V2 blob
 */
struct BlobCommitment {
    BN254.G1Point commitment;
    BN254.G2Point lengthCommitment;
    BN254.G2Point lengthProof;
    uint32 length;
}

/**
 * @title SignedBatch
 * @notice A signed V2 batch
 */
struct SignedBatch {
    BatchHeaderV2 batchHeader;
    Attestation attestation;
}

/**
 * @title BatchHeaderV2
 * @notice The header of a V2 batch
 */
struct BatchHeaderV2 {
    bytes32 batchRoot;
    uint32 referenceBlockNumber;
}

/**
 * @title Attestation
 * @notice The attestation for a V2 batch
 */
struct Attestation {
    BN254.G1Point[] nonSignerPubkeys;
    BN254.G1Point[] quorumApks;
    BN254.G1Point sigma;
    BN254.G2Point apkG2;
    uint32[] quorumNumbers;
}

///////////////////////// SIGNATURE VERIFIER ///////////////////////////////

/**
 * @title NonSignerStakesAndSignature
 * @notice The non-signer stakes and signatures used for BLS signature verification of V1 and V2 batches
 */
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

/**
 * @title QuorumStakeTotals
 * @notice The total stakes signed for each quorum returned by the signature verifier
 */
struct QuorumStakeTotals {
    uint96[] signedStakeForQuorum;
    uint96[] totalStakeForQuorum;
}

/**
 * @title CheckSignaturesIndices
 * @notice The indices needed for checking signatures of V1 and V2 batches
 */
struct CheckSignaturesIndices {
    uint32[] nonSignerQuorumBitmapIndices;
    uint32[] quorumApkIndices;
    uint32[] totalStakeIndices;  
    uint32[][] nonSignerStakeIndices;
}