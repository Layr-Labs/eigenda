use alloy_sol_types::sol;

// commented-out structs left intentionally to serve as documentation of which types were not require
sol! {
    ///////////////////////// BN254 ///////////////////////////////

    #[derive(Default, Debug)]
    struct G1Point {
        uint256 X;
        uint256 Y;
    }

    // Encoding of field elements is: X[1] * i + X[0]
    #[derive(Default, Debug)]
    struct G2Point {
        uint256[2] X;
        uint256[2] Y;
    }

    ///////////////////////// V2 ///////////////////////////////

    #[derive(Default, Debug)]
    struct VersionedBlobParams {
        uint32 maxNumOperators;
        uint32 numChunks;
        uint8 codingRate;
    }

    #[derive(Default, Debug)]
    struct SecurityThresholds {
        uint8 confirmationThreshold;
        uint8 adversaryThreshold;
    }

    #[derive(Default, Debug)]
    struct RelayInfo {
        address relayAddress;
        string relayURL;
    }

    // struct DisperserInfo {
    //     address disperserAddress;
    // }

    #[derive(Default, Debug)]
    struct BlobInclusionInfo {
        BlobCertificate blobCertificate;
        uint32 blobIndex;
        bytes inclusionProof;
    }

    #[derive(Default, Debug)]
    struct BlobCertificate {
        BlobHeaderV2 blobHeader;
        bytes signature;
        uint32[] relayKeys;
    }

    #[derive(Default, Debug)]
    struct BlobHeaderV2 {
        uint16 version;
        bytes quorumNumbers;
        BlobCommitment commitment;
        bytes32 paymentHeaderHash;
    }

    #[derive(Default, Debug)]
    struct BlobCommitment {
        G1Point commitment;
        G2Point lengthCommitment;
        G2Point lengthProof;
        uint32 length;
    }

    // struct SignedBatch {
    //     BatchHeaderV2 batchHeader;
    //     Attestation attestation;
    // }

    #[derive(Default, Debug)]
    struct BatchHeaderV2 {
        bytes32 batchRoot;
        uint32 referenceBlockNumber;
    }

    // struct Attestation {
    //     BN254.G1Point[] nonSignerPubkeys;
    //     BN254.G1Point[] quorumApks;
    //     BN254.G1Point sigma;
    //     BN254.G2Point apkG2;
    //     uint32[] quorumNumbers;
    // }

    ///////////////////////// SIGNATURE VERIFIER ///////////////////////////////

    #[derive(Default, Debug)]
    struct NonSignerStakesAndSignature {
        uint32[] nonSignerQuorumBitmapIndices;
        G1Point[] nonSignerPubkeys;
        G1Point[] quorumApks;
        G2Point apkG2;
        G1Point sigma;
        uint32[] quorumApkIndices;
        uint32[] totalStakeIndices;
        uint32[][] nonSignerStakeIndices;
    }

    // struct QuorumStakeTotals {
    //     uint96[] signedStakeForQuorum;
    //     uint96[] totalStakeForQuorum;
    // }

    // struct CheckSignaturesIndices {
    //     uint32[] nonSignerQuorumBitmapIndices;
    //     uint32[] quorumApkIndices;
    //     uint32[] totalStakeIndices;
    //     uint32[][] nonSignerStakeIndices;
    // }
}
