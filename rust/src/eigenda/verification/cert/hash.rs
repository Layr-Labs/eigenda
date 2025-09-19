//! Cryptographic hashing utilities for EigenDA structures
//!
//! This module provides hashing functions and types for computing cryptographic
//! digests of EigenDA data structures, following the same hashing conventions
//! used in the on-chain smart contracts.

use std::fmt::Display;

use crate::eigenda::cert::{BatchHeaderV2, BlobCertificate, BlobHeaderV2};
use alloy_primitives::{B256, Keccak256, keccak256};
use alloy_sol_types::SolValue;
use derive_more::{AsMut, AsRef, Deref, DerefMut, From, Into};

/// A truncated 24-byte hash used for aggregate public key identification.
///
/// EigenDA uses truncated hashes of aggregate public keys to efficiently
/// identify and reference APKs in storage while maintaining collision
/// resistance for practical purposes.
#[repr(transparent)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, Hash, Deref, DerefMut, AsRef, AsMut, From, Into, Default,
)]
pub struct TruncHash(pub [u8; 24]);

impl Display for TruncHash {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", hex::encode(self.0))
    }
}

/// Extension trait for computing EigenDA-compatible hashes of data structures.
///
/// Provides standardized hashing methods that match the hashing logic
/// used in EigenDA smart contracts for consistent verification.
pub trait HashExt {
    /// Compute the EigenDA-compatible hash of this structure
    fn hash_ext(&self) -> B256;
}

impl HashExt for BlobCertificate {
    /// Hash a blob certificate using EigenDA's standard encoding
    ///
    /// Computes: keccak256(abi.encode(blob_header_hash, signature, relay_keys))
    fn hash_ext(&self) -> B256 {
        let blob_header = self.blob_header.hash_ext();
        let encoded =
            (blob_header, self.signature.clone(), self.relay_keys.clone()).abi_encode_sequence();
        keccak256(&encoded)
    }
}

impl HashExt for BlobHeaderV2 {
    /// Hash a blob header using EigenDA's standard encoding
    ///
    /// Two-step process:
    /// 1. Hash the core blob data: keccak256(abi.encode(version, quorum_numbers, commitment))
    /// 2. Hash with payment info: keccak256(abi.encode(core_hash, payment_header_hash))
    fn hash_ext(&self) -> B256 {
        let encoded = (
            self.version,
            self.quorum_numbers.clone(),
            self.commitment.to_sol(),
        )
            .abi_encode_sequence();

        let hashed = keccak256(&encoded);
        let encoded = (hashed, self.payment_header_hash).abi_encode();
        keccak256(&encoded)
    }
}

impl HashExt for BatchHeaderV2 {
    /// Hash a batch header using EigenDA's standard encoding
    ///
    /// Computes: keccak256(abi.encode(batch_root, reference_block_number))
    fn hash_ext(&self) -> B256 {
        let encoded = self.to_sol().abi_encode();
        keccak256(&encoded)
    }
}

/// Compute keccak256 hash over a sequence of byte arrays.
///
/// This is useful for hashing large amounts of data without concatenating
/// everything into memory first. Updates the hasher incrementally with
/// each provided byte array.
///
/// # Arguments
/// * `values` - Iterator of byte arrays to hash
///
/// # Returns
/// 32-byte keccak256 hash digest
pub fn streaming_keccak256<T: AsRef<[u8]>>(values: &[T]) -> B256 {
    let mut hasher = Keccak256::new();
    for v in values {
        hasher.update(v.as_ref());
    }
    hasher.finalize()
}

#[cfg(test)]
mod tests {
    use std::str::FromStr;

    use alloy_primitives::{B256, Bytes, keccak256};

    use crate::eigenda::{
        cert::{BatchHeaderV2, BlobCertificate, BlobCommitment, BlobHeaderV2},
        verification::cert::hash::{HashExt, TruncHash, streaming_keccak256},
    };

    #[test]
    fn blob_certificate_hash_ext() {
        let cert = BlobCertificate {
            blob_header: BlobHeaderV2 {
                version: 1,
                quorum_numbers: Bytes::from(vec![0u8, 1u8]),
                commitment: BlobCommitment::default(),
                payment_header_hash: [0u8; 32],
            },
            signature: Bytes::from(vec![1u8, 2u8, 3u8]),
            relay_keys: vec![],
        };

        let actual = cert.hash_ext();
        let expected =
            B256::from_str("0x7f8946919c6354b9dd8488a279fd919798adafc7a2023a308f766e157919c124")
                .unwrap();

        assert_eq!(actual, expected);
    }

    #[test]
    fn blob_header_v2_hash_ext() {
        let header = BlobHeaderV2 {
            version: 2,
            quorum_numbers: Bytes::from(vec![0u8]),
            commitment: BlobCommitment::default(),
            payment_header_hash: [1u8; 32],
        };

        let actual = header.hash_ext();
        let expected =
            B256::from_str("0x49508c922e2a74bfa0ae0e942aac3aacc28ababb4d4ffc823bb9fc5d3a858cca")
                .unwrap();

        assert_eq!(actual, expected);
    }

    #[test]
    fn batch_header_v2_hash_ext() {
        let header = BatchHeaderV2 {
            batch_root: [2u8; 32],
            reference_block_number: 12345,
        };

        let actual = header.hash_ext();
        let expected =
            B256::from_str("0xe231c6b7b4ff73c5300b4f46c8d880301e4f08356f9f7f307937a8b8ca397339")
                .unwrap();

        assert_eq!(actual, expected);
    }

    #[test]
    fn test_streaming_keccak256() {
        let values = vec![b"hello".as_slice(), b"world".as_slice()];
        let result = streaming_keccak256(&values);
        let expected = keccak256(b"helloworld");

        assert_eq!(result, expected);
    }

    #[test]
    fn trunc_hash_display() {
        let hash = TruncHash([
            1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
        ]);
        let actual = format!("{}", hash);
        let expected = "0102030405060708090a0102030405060708090a0b0c0d0e";
        assert_eq!(actual, expected);
    }
}
