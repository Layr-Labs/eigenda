use alloy_sol_types::sol;

sol! {
    /// Version 2 batch header for EigenDA protocol
    ///
    /// Contains essential metadata about a batch of blobs in the EigenDA network.
    /// This header version is used across multiple certificate versions (V2, V3) as it
    /// represents the EigenDA protocol version rather than the certificate version.
    ///
    /// Reference: https://github.com/Layr-Labs/eigenda/blob/510291b9be38cacbed8bc62125f6f9a14bd604e4/contracts/src/core/libraries/v2/EigenDATypesV2.sol#L47
    #[derive(Debug)]
    struct BatchHeaderV2 {
        /// Merkle root hash that summarizes all data blobs in this batch
        ///
        /// This cryptographic commitment allows efficient verification of blob inclusion
        /// within the batch without needing to download all batch data.
        bytes32 batchRoot;
        
        /// Ethereum block number used as reference point for operator set verification
        ///
        /// This block number serves as a "snapshot" of the EigenDA operator set state
        /// for BLS signature verification. When operators sign batches, their stakes and
        /// registered quorums are validated against the historical state at this specific
        /// block number. This ensures signature verification uses a consistent view of
        /// the operator set even if operators join/leave or update their stakes after
        /// creating their signatures.
        ///
        /// The reference block number must be:
        /// - Less than the current block number when verification occurs
        /// - Within the stale stakes window (if stale stakes are forbidden)  
        /// - Used consistently across all operator state lookups during verification
        ///
        /// See: [BLSSignatureChecker.checkSignatures](https://github.com/Layr-Labs/eigenlayer-middleware/blob/dev/docs/BLSSignatureChecker.md#blssignaturecheckerchecksignatures)
        uint32 referenceBlockNumber;
    }

    /// Point on the BN254 G1 elliptic curve group
    ///
    /// G1 points are used in EigenDA for:
    /// - Public keys of operators in the network
    /// - Cryptographic commitments to blob data
    /// - Signature aggregation in the BLS signature scheme
    ///
    /// The BN254 curve is specifically chosen for its pairing-friendly properties
    /// which enable efficient zero-knowledge proofs and signature verification.
    #[derive(Debug)]
    struct G1Point {
        /// X coordinate of the point on the curve
        uint256 X;
        /// Y coordinate of the point on the curve  
        uint256 Y;
    }

    /// Point on the BN254 G2 elliptic curve group
    ///
    /// G2 points are used in EigenDA for:
    /// - Length commitments and proofs in polynomial commitments
    /// - Aggregated public keys in BLS signature schemes
    /// - Pairing operations for cryptographic verification
    ///
    /// Encoding of field elements is: X[1] * i + X[0]
    /// This is because of the (unknown to us) convention used in the bn254 pairing precompile contract
    /// "Elements a * i + b of F_p^2 are encoded as two elements of F_p, (a, b)."
    /// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-197.md#encoding
    #[derive(Debug)]
    struct G2Point {
        /// X coordinate as a field extension element [X0, X1] where X = X0 + X1*i
        uint256[2] X;
        /// Y coordinate as a field extension element [Y0, Y1] where Y = Y0 + Y1*i
        uint256[2] Y;
    }

    /// Cryptographic commitment to a data blob using polynomial commitments
    ///
    /// Contains the necessary cryptographic proofs to verify
    /// the integrity and properties of a data blob without downloading it.
    /// Uses KZG polynomial commitments over the BN254 curve
    #[derive(Debug)]
    struct BlobCommitment {
        /// KZG commitment to the blob data polynomial
        ///
        /// This G1 point represents a cryptographic binding to the entire blob
        /// content, allowing verification of the data's integrity.
        G1Point commitment;
        
        /// KZG commitment to the length of the blob
        ///
        /// Proves the claimed length of the blob data
        G2Point lengthCommitment;
        
        /// KZG proof for the length commitment
        ///
        /// Cryptographic proof that demonstrates the length commitment
        /// was computed correctly for the claimed blob length.
        G2Point lengthProof;
        
        /// Actual length of the blob data in bytes
        ///
        /// The proven length of the blob that corresponds to the
        /// `lengthCommitment` and `lengthProof`.
        uint32 length;
    }
}
