// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";

library EigenDATypesV2 {
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

    struct VersionedBlobParams {
        uint32 codingRate;
        uint32 reconstructionThreshold;
        uint32 numChunks;
        uint32 numUnits;
        uint32 samplesPerUnit;
    }
}
