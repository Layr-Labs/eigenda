// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {BN254} from "eigenlayer-middleware/libraries/BN254.sol";

///////////////////////// V1 ///////////////////////////////

struct QuorumBlobParam {
    uint8 quorumNumber;
    uint8 adversaryThresholdPercentage;
    uint8 confirmationThresholdPercentage; 
    uint32 chunkLength; 
}

struct BlobHeader {
    BN254.G1Point commitment; 
    uint32 dataLength; 
    QuorumBlobParam[] quorumBlobParams; 
}

struct ReducedBatchHeader {
    bytes32 blobHeadersRoot;
    uint32 referenceBlockNumber;
}

struct BatchHeader {
    bytes32 blobHeadersRoot;
    bytes quorumNumbers; 
    bytes signedStakeForQuorums; 
    uint32 referenceBlockNumber;
}
    
struct BatchMetadata {
    BatchHeader batchHeader; 
    bytes32 signatoryRecordHash; 
    uint32 confirmationBlockNumber; 
}

struct BlobVerificationProof {
    uint32 batchId;
    uint32 blobIndex;
    BatchMetadata batchMetadata;
    bytes inclusionProof;
    bytes quorumIndices;
}

///////////////////////// V2 ///////////////////////////////

struct VersionedBlobParams {
    uint32 maxNumOperators;
    uint32 numChunks;
    uint8 codingRate;
}

struct SecurityThresholds {
    uint8 confirmationThreshold;
    uint8 adversaryThreshold;
}

struct RelayInfo {
    address relayAddress;
    string relayURL;
}

struct DisperserInfo {
    address disperserAddress;
}

struct BlobInclusionInfo {
    BlobCertificate blobCertificate;
    uint32 blobIndex;
    bytes inclusionProof;
}

struct BlobCertificate {
    BlobHeaderV2 blobHeader;
    bytes signature;
    uint32[] relayKeys;
}

struct BlobHeaderV2 {
    uint16 version;
    bytes quorumNumbers;
    BlobCommitment commitment;
    bytes32 paymentHeaderHash;
    uint32 salt;
}

struct BlobCommitment {
    BN254.G1Point commitment;
    BN254.G2Point lengthCommitment;
    BN254.G2Point lengthProof;
    uint32 length;
}

struct SignedBatch {
    BatchHeaderV2 batchHeader;
    Attestation attestation;
}

struct BatchHeaderV2 {
    bytes32 batchRoot;
    uint32 referenceBlockNumber;
}

struct Attestation {
    BN254.G1Point[] nonSignerPubkeys;
    BN254.G1Point[] quorumApks;
    BN254.G1Point sigma;
    BN254.G2Point apkG2;
    uint32[] quorumNumbers;
}

///////////////////////// SIGNATURE VERIFIER ///////////////////////////////

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

struct QuorumStakeTotals {
    uint96[] signedStakeForQuorum;
    uint96[] totalStakeForQuorum;
}

struct CheckSignaturesIndices {
    uint32[] nonSignerQuorumBitmapIndices;
    uint32[] quorumApkIndices;
    uint32[] totalStakeIndices;  
    uint32[][] nonSignerStakeIndices;
}