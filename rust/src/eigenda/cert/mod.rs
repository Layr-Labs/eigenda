//! EigenDA Certificate Types and Structures
//!
//! This module provides Rust types for working with EigenDA certificates, which are used
//! to verify the inclusion of data blobs in the EigenDA network. The module supports
//! both version 2 and version 3 certificates.
//!
//! ## Key Components
//!
//! - [`StandardCommitment`] - Main wrapper for versioned certificates with RLP serialization
//! - [`EigenDAVersionedCert`] - Enum representing different certificate versions
//! - [`EigenDACertV2`]/[`EigenDACertV3`] - Version-specific certificate structures
//! - [`BlobInclusionInfo`] - Information about blob inclusion in batches
//! - [`BatchHeaderV2`] - Batch header containing batch root and reference block
//! - [`G1Point`]/[`G2Point`] - Elliptic curve points for cryptographic operations
//! ```
//!
//! ## Notes
//!
//! - Due to sovereign-sdk compatibility constraints these types could not be imported from [rust-eigenda-v2-common](https://crates.io/crates/rust-eigenda-v2-common)

mod solidity;

use alloy_primitives::{B256, Bytes, FixedBytes, U256};
use alloy_rlp::{Decodable, RlpDecodable, RlpEncodable};
use alloy_rlp::{Encodable, Error};
use serde::{Deserialize, Serialize};
use thiserror::Error;

use crate::eigenda::verification::cert::convert;
use crate::eigenda::verification::cert::types::RelayKey;

/// Byte indicating a version 2 certificate.
const VERSION_2: u8 = 1;

/// Byte indicating a version 3 certificate.
const VERSION_3: u8 = 2;

/// Main wrapper for EigenDA certificates supporting multiple versions.
///
/// This structure provides a unified interface for working with different versions
/// of EigenDA certificates (V2 and V3). It handles RLP serialization/deserialization
/// and provides version-agnostic access to certificate data.
#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct StandardCommitment(EigenDAVersionedCert);

impl StandardCommitment {
    /// Parse a certificate from RLP-encoded bytes.
    ///
    /// The first byte indicates the certificate version (1 for V2, 2 for V3),
    /// followed by the RLP-encoded certificate data.
    ///
    /// # Arguments
    ///
    /// * `bytes` - The RLP-encoded certificate bytes including version prefix
    ///
    /// # Returns
    ///
    /// Returns a `StandardCommitment` on success, or a parse error if the data
    /// is invalid or uses an unsupported version.
    ///
    /// # Errors
    ///
    /// * [`StandardCommitmentParseError::InsufficientData`] - If bytes are empty
    /// * [`StandardCommitmentParseError::UnsupportedCertVersion`] - If version is not supported
    /// * [`StandardCommitmentParseError::InvalidRlpCert`] - If RLP decoding fails
    pub fn from_rlp_bytes(bytes: &[u8]) -> Result<Self, StandardCommitmentParseError> {
        let (cert_version, mut cert_bytes) = bytes
            .split_first()
            .ok_or(StandardCommitmentParseError::EmptyCommitment)?;

        let versioned_cert = match *cert_version {
            VERSION_2 => {
                let cert = EigenDACertV2::decode(&mut cert_bytes)
                    .map_err(StandardCommitmentParseError::InvalidRlpCert)?;
                EigenDAVersionedCert::V2(cert)
            }
            VERSION_3 => {
                let cert = EigenDACertV3::decode(&mut cert_bytes)
                    .map_err(StandardCommitmentParseError::InvalidRlpCert)?;
                EigenDAVersionedCert::V3(cert)
            }
            _ => {
                return Err(StandardCommitmentParseError::UnsupportedCertVersion(
                    *cert_version,
                ));
            }
        };

        Ok(Self(versioned_cert))
    }

    /// Serialize the certificate to RLP-encoded bytes.
    ///
    /// The output includes a version byte prefix followed by the RLP-encoded
    /// certificate data.
    ///
    /// # Returns
    ///
    /// Returns the complete certificate as RLP-encoded bytes with version prefix.
    pub fn to_rlp_bytes(&self) -> Bytes {
        let mut bytes = Vec::new();
        match &self.0 {
            EigenDAVersionedCert::V2(c) => {
                bytes.push(VERSION_2);
                c.encode(&mut bytes);
            }
            EigenDAVersionedCert::V3(c) => {
                bytes.push(VERSION_3);
                c.encode(&mut bytes);
            }
        }

        Bytes::from(bytes)
    }

    /// Get the reference block number used when constructing this certificate.
    ///
    /// The reference block number is used for verifying the certificate against
    /// the blockchain state at a specific block height.
    ///
    /// # Returns
    ///
    /// Returns the reference block number as a u64.
    pub fn reference_block(&self) -> u64 {
        match &self.0 {
            EigenDAVersionedCert::V2(c) => c.batch_header_v2.reference_block_number as u64,
            EigenDAVersionedCert::V3(c) => c.batch_header_v2.reference_block_number as u64,
        }
    }

    /// Get the blob header version from the certificate.
    ///
    /// # Returns
    ///
    /// Returns the blob header version number.
    pub fn version(&self) -> u16 {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => {
                cert.blob_inclusion_info
                    .blob_certificate
                    .blob_header
                    .version
            }
            EigenDAVersionedCert::V3(cert) => {
                cert.blob_inclusion_info
                    .blob_certificate
                    .blob_header
                    .version
            }
        }
    }

    /// Get hashes of public keys of non-signing validators.
    ///
    /// These are validators that did not participate in signing the certificate.
    ///
    /// # Returns
    ///
    /// Returns a vector of 32-byte hashes of non-signer public keys.
    pub fn non_signers_pk_hashes(&self) -> Vec<B256> {
        let pks = match &self.0 {
            EigenDAVersionedCert::V2(cert) => {
                cert.nonsigner_stake_and_signature.non_signer_pubkeys.iter()
            }
            EigenDAVersionedCert::V3(cert) => {
                cert.nonsigner_stake_and_signature.non_signer_pubkeys.iter()
            }
        };

        // not the same version of G1Point
        pks.map(convert::point_to_hash).collect()
    }

    /// Get indices in the quorum bitmap for non-signing validators.
    ///
    /// # Returns
    ///
    /// Returns a slice of indices corresponding to non-signers in the quorum bitmap.
    pub fn non_signer_quorum_bitmap_indices(&self) -> &[u32] {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => {
                &cert
                    .nonsigner_stake_and_signature
                    .non_signer_quorum_bitmap_indices
            }

            EigenDAVersionedCert::V3(cert) => {
                &cert
                    .nonsigner_stake_and_signature
                    .non_signer_quorum_bitmap_indices
            }
        }
    }

    /// Get the quorums that signed this certificate.
    ///
    /// # Returns
    ///
    /// Returns the quorum numbers as bytes.
    pub fn signed_quorum_numbers(&self) -> &Bytes {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => &cert.signed_quorum_numbers,
            EigenDAVersionedCert::V3(cert) => &cert.signed_quorum_numbers,
        }
    }

    /// Get indices of aggregate public keys for each quorum.
    ///
    /// # Returns
    ///
    /// Returns indices pointing to the aggregate public keys used for verification.
    pub fn quorum_apk_indices(&self) -> &[u32] {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => {
                &cert.nonsigner_stake_and_signature.quorum_apk_indices
            }
            EigenDAVersionedCert::V3(cert) => {
                &cert.nonsigner_stake_and_signature.quorum_apk_indices
            }
        }
    }

    /// Get indices of total stakes for non-signing operators.
    ///
    /// # Returns
    ///
    /// Returns indices for looking up total stake amounts of non-signers.
    pub fn non_signer_total_stake_indices(&self) -> &[u32] {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => {
                &cert.nonsigner_stake_and_signature.total_stake_indices
            }
            EigenDAVersionedCert::V3(cert) => {
                &cert.nonsigner_stake_and_signature.total_stake_indices
            }
        }
    }

    /// Get stake indices for non-signing operators per quorum.
    ///
    /// # Returns
    ///
    /// Returns a nested vector of stake indices, organized by quorum.
    pub fn non_signer_stake_indices(&self) -> &[Vec<u32>] {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => {
                &cert.nonsigner_stake_and_signature.non_signer_stake_indices
            }
            EigenDAVersionedCert::V3(cert) => {
                &cert.nonsigner_stake_and_signature.non_signer_stake_indices
            }
        }
    }

    /// Get a reference to the batch header.
    ///
    /// # Returns
    ///
    /// Returns a reference to the BatchHeaderV2 containing batch metadata.
    pub fn batch_header_v2(&self) -> &BatchHeaderV2 {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => &cert.batch_header_v2,
            EigenDAVersionedCert::V3(cert) => &cert.batch_header_v2,
        }
    }

    /// Get blob inclusion information.
    ///
    /// # Returns
    ///
    /// Returns blob inclusion metadata.
    pub fn blob_inclusion_info(&self) -> &BlobInclusionInfo {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => &cert.blob_inclusion_info,
            EigenDAVersionedCert::V3(cert) => &cert.blob_inclusion_info,
        }
    }

    /// Get non-signer stakes and signature information.
    ///
    /// # Returns
    ///
    /// Returns complete information about non-signers including their stakes and signatures.
    pub fn nonsigner_stake_and_signature(&self) -> &NonSignerStakesAndSignature {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => &cert.nonsigner_stake_and_signature,
            EigenDAVersionedCert::V3(cert) => &cert.nonsigner_stake_and_signature,
        }
    }
}

/// EigenDA versioned certificate enum.
///
/// This enum wraps different versions of EigenDA certificates, allowing
/// the system to handle multiple certificate formats transparently.
#[derive(Debug, PartialEq, Clone, Serialize, Deserialize)]
pub enum EigenDAVersionedCert {
    /// Version 2 certificate
    V2(EigenDACertV2),
    /// Version 3 certificate
    V3(EigenDACertV3),
}

/// EigenDA Certificate Version 2.
///
/// This structure represents a version 2 certificate containing all necessary
/// information for verifying blob inclusion in the EigenDA network.
#[derive(Debug, Clone, RlpEncodable, RlpDecodable, PartialEq, Serialize, Deserialize)]
pub struct EigenDACertV2 {
    /// Information about blob inclusion in the batch
    pub blob_inclusion_info: BlobInclusionInfo,
    /// Batch header containing batch metadata
    pub batch_header_v2: BatchHeaderV2,
    /// Non-signer information and signatures
    pub nonsigner_stake_and_signature: NonSignerStakesAndSignature,
    /// Numbers of quorums that signed this certificate
    pub signed_quorum_numbers: Bytes,
}

/// EigenDA Certificate Version 3.
///
/// This structure represents a version 3 certificate with the same core components
/// as V2 but potentially different field ordering or processing logic.
#[derive(Debug, Clone, RlpEncodable, RlpDecodable, PartialEq, Serialize, Deserialize)]
pub struct EigenDACertV3 {
    /// Batch header containing batch metadata
    pub batch_header_v2: BatchHeaderV2,
    /// Information about blob inclusion in the batch
    pub blob_inclusion_info: BlobInclusionInfo,
    /// Non-signer information and signatures
    pub nonsigner_stake_and_signature: NonSignerStakesAndSignature,
    /// Numbers of quorums that signed this certificate
    pub signed_quorum_numbers: Bytes,
}

/// Batch Header Version 2 as defined by the EigenDA protocol.
///
/// This version is separate from the certificate version. For example, Certificate V3
/// can use BatchHeaderV2 since V2 is a tag for the EigenDA protocol. The V2 suffix
/// matches the corresponding Solidity struct name.
///
/// Reference: [EigenDATypesV2.sol](https://github.com/Layr-Labs/eigenda/blob/510291b9be38cacbed8bc62125f6f9a14bd604e4/contracts/src/core/libraries/v2/EigenDATypesV2.sol#L47)
#[derive(Debug, Clone, RlpEncodable, RlpDecodable, PartialEq, Serialize, Deserialize)]
pub struct BatchHeaderV2 {
    /// 32-byte root hash of the batch merkle tree
    pub batch_root: [u8; 32],
    /// Ethereum block number used as reference point for operator set verification
    ///
    /// This block number serves as a "snapshot" of the EigenDA operator set state
    /// for signature verification. When operators sign batches, their stakes and
    /// registered quorums are validated against the historical state at this specific
    /// block number, ensuring that signature verification uses a consistent view of
    /// the operator set even if operators join/leave or update their stakes after
    /// creating their signatures.
    ///
    /// The reference block number must be:
    /// - Less than the current block number when verification occurs
    /// - Within the stale stakes window (if stale stakes are forbidden)
    /// - Used consistently across all operator state lookups during verification
    ///
    /// See: [BLSSignatureChecker.checkSignatures](https://github.com/Layr-Labs/eigenlayer-middleware/blob/dev/docs/BLSSignatureChecker.md#blssignaturecheckerchecksignatures)
    pub reference_block_number: u32,
}

impl BatchHeaderV2 {
    /// Convert this batch header to its Solidity representation.
    ///
    /// The V2 suffix matches the corresponding Solidity struct name in the EigenDA contracts.
    ///
    /// Reference: [EigenDATypesV2.sol](https://github.com/Layr-Labs/eigenda/blob/510291b9be38cacbed8bc62125f6f9a14bd604e4/contracts/src/core/libraries/v2/EigenDATypesV2.sol#L28)
    ///
    /// # Returns
    ///
    /// Returns a `solidity::BatchHeaderV2` struct for use in contract interactions.
    pub fn to_sol(&self) -> solidity::BatchHeaderV2 {
        solidity::BatchHeaderV2 {
            batchRoot: FixedBytes::<32>(self.batch_root),
            referenceBlockNumber: self.reference_block_number,
        }
    }
}

/// Information required to prove blob inclusion in a batch.
///
/// This structure contains all the data needed to verify that a specific blob
/// is included in the batch at the specified index.
#[derive(Debug, Clone, RlpEncodable, RlpDecodable, PartialEq, Serialize, Deserialize)]
pub struct BlobInclusionInfo {
    /// Certificate containing blob metadata and commitments
    pub blob_certificate: BlobCertificate,
    /// Index of the blob within the batch
    pub blob_index: u32,
    /// Merkle proof data for inclusion verification
    pub inclusion_proof: Bytes,
}

/// Certificate containing all necessary information about a blob.
///
/// This structure includes the blob header with commitments, signatures,
/// and relay keys used for blob retrieval.
#[derive(Debug, Clone, RlpEncodable, RlpDecodable, PartialEq, Serialize, Deserialize)]
pub struct BlobCertificate {
    /// Header containing blob metadata and commitments
    pub blob_header: BlobHeaderV2,
    /// Cryptographic signature over the blob data
    pub signature: Bytes,
    /// Keys for relaying/retrieving the blob data
    pub relay_keys: Vec<RelayKey>,
}

/// Blob Header Version 2 containing blob metadata and commitments.
///
/// This version is separate from the certificate version. For example, Certificate V3
/// can use BlobHeaderV2 since V2 is a tag for the EigenDA protocol.
#[derive(Debug, Clone, RlpEncodable, RlpDecodable, PartialEq, Serialize, Deserialize)]
pub struct BlobHeaderV2 {
    /// Version number of the blob header format
    pub version: u16,
    /// Numbers identifying which quorums store this blob
    pub quorum_numbers: Bytes,
    /// Cryptographic commitment to the blob data
    pub commitment: BlobCommitment,
    /// Hash of the payment header for this blob
    pub payment_header_hash: [u8; 32],
}

/// Cryptographic commitments for verifying blob data integrity.
///
/// This structure contains KZG polynomial commitments that allow verification
/// of blob data without requiring the full blob content.
#[derive(Debug, Clone, RlpEncodable, RlpDecodable, PartialEq, Serialize, Deserialize, Default)]
pub struct BlobCommitment {
    /// KZG commitment to the blob polynomial (G1 point)
    pub commitment: G1Point,
    /// Commitment to the length of the blob (G2 point)
    pub length_commitment: G2Point,
    /// Proof for the length commitment (G2 point)
    pub length_proof: G2Point,
    /// Actual length of the blob in bytes
    pub length: u32,
}

impl BlobCommitment {
    /// Convert this blob commitment to its Solidity representation.
    ///
    /// # Returns
    ///
    /// Returns a `solidity::BlobCommitment` struct for use in contract interactions.
    pub fn to_sol(&self) -> solidity::BlobCommitment {
        solidity::BlobCommitment {
            commitment: (&self.commitment).into(),
            lengthCommitment: (&self.length_commitment).into(),
            lengthProof: (&self.length_proof).into(),
            length: self.length,
        }
    }
}

/// A point on the BN254 elliptic curve G1 subgroup.
///
/// G1 points are used for cryptographic commitments and signatures in EigenDA.
/// The BN254 curve is also known as the alt-bn128 curve.
#[derive(
    Debug, Clone, Copy, RlpEncodable, RlpDecodable, PartialEq, Serialize, Deserialize, Default,
)]
pub struct G1Point {
    /// X coordinate of the point
    pub x: U256,
    /// Y coordinate of the point
    pub y: U256,
}

impl From<&G1Point> for solidity::G1Point {
    /// Convert a G1Point to its Solidity representation.
    fn from(value: &G1Point) -> Self {
        solidity::G1Point {
            X: value.x,
            Y: value.y,
        }
    }
}

/// A point on the BN254 elliptic curve G2 subgroup.
///
/// G2 points are used for pairing-based cryptographic operations. Each coordinate
/// is represented as a vector of two U256 values forming an element in the quadratic
/// extension field Fp2.
#[derive(Debug, Clone, RlpEncodable, RlpDecodable, PartialEq, Serialize, Deserialize)]
pub struct G2Point {
    /// X coordinate as an Fp2 element [x0, x1]
    pub x: Vec<U256>,
    /// Y coordinate as an Fp2 element [y0, y1]
    pub y: Vec<U256>,
}

impl From<&G2Point> for solidity::G2Point {
    /// Convert a G2Point to its Solidity representation.
    ///
    /// Maps the Fp2 coordinates to the fixed-size arrays expected by Solidity.
    fn from(value: &G2Point) -> solidity::G2Point {
        let mut x = [U256::default(); 2];
        x[0] = value.x[0];
        x[1] = value.x[1];

        let mut y = [U256::default(); 2];
        y[0] = value.y[0];
        y[1] = value.y[1];

        solidity::G2Point { X: x, Y: y }
    }
}

/// Information about validators who did not sign the certificate.
///
/// This structure contains all data needed to verify the aggregate signature
/// while accounting for validators that did not participate in signing.
#[derive(Debug, Clone, RlpEncodable, RlpDecodable, PartialEq, Serialize, Deserialize)]
pub struct NonSignerStakesAndSignature {
    /// Indices of non-signers in the quorum bitmap
    pub non_signer_quorum_bitmap_indices: Vec<u32>,
    /// Public keys of validators that did not sign
    pub non_signer_pubkeys: Vec<G1Point>,
    /// Aggregate public keys for each quorum
    pub quorum_apks: Vec<G1Point>,
    /// Aggregate public key in G2 for pairing verification
    pub apk_g2: G2Point,
    /// BLS signature aggregated from all signers
    pub sigma: G1Point,
    /// Indices for quorum aggregate public keys
    pub quorum_apk_indices: Vec<u32>,
    /// Indices for total stake lookups
    pub total_stake_indices: Vec<u32>,
    /// Nested indices for non-signer stakes per quorum
    pub non_signer_stake_indices: Vec<Vec<u32>>,
}

/// Errors that can occur when parsing a `StandardCommitment` from bytes.
#[derive(Debug, Error)]
pub enum StandardCommitmentParseError {
    /// Empty commitment data (tx calldata contains 0 bytes)
    #[error("Empty commitment data")]
    EmptyCommitment,
    /// Unsupported cert version
    #[error("Unsupported cert version {0}")]
    UnsupportedCertVersion(u8),
    /// The cert couldn't be parsed from the RLP format
    #[error("Invalid RLP Cert")]
    InvalidRlpCert(Error),
}

#[cfg(test)]
mod tests {
    use super::{StandardCommitment, StandardCommitmentParseError};

    #[test]
    fn v2_serialization_round_trip() {
        let commitment_hex = "02f90389e5a0c769488dd5264b3ef21dce7ee2d42fba43e1f83ff228f501223e38818cb14492833f44fcf901eff901caf9018180820001f90159f842a0012e810ffc0a83074b3d14db9e78bbae623f7770cac248df9e73fac6b9d59d17a02a916ffbbf9dde4b7ebe94191a29ff686422d7dcb3b47ecb03c6ada75a9c15c8f888f842a01811c8b4152fce9b8c4bae61a3d097e61dfc43dc7d45363d19e7c7f1374034ffa001edc62174217cdce60a4b52fa234ac0d96db4307dac9150e152ba82cbb4d2f1f842a00f423b0dbc1fe95d2e3f7dbac6c099e51dbf73400a4b3f26b9a29665b4ac58a8a01855a2bd56c0e8f4cc85ac149cf9a531673d0e89e22f0d6c4ae419ed7c5d2940f888f842a02667cbb99d60fa0d7f3544141d3d531dceeeb50b06e5a0cdc42338a359138ae4a00dff4c929d8f8a307c19bba6e8006fe6700f6554cef9eb3797944f89472ffb30f842a004c17a6225acd5b4e7d672a1eb298c5358f4f6f17d04fd1ee295d0c0d372fa84a024bc3ad4d5e54f54f71db382ce276f37ac3c260cc74306b832e8a3c93c7951d302a0e43e11e2405c2fd1d880af8612d969b654827e0ba23d9feb3722ccce6226fce7b8411ddf4553c79c0515516fd3c8b3ae6a756b05723f4d0ebe98a450c8bcc96cbb355ef07a44eeb56f831be73647e4da20e22fa859f984ee41d6efcd3692063b0b0601c2800101a0a69e552a6fc2ff75d32edaf5313642ddeebe60d2069435d12e266ce800e9e96bf9016bc0c0f888f842a00d45727a99053af8d38d4716ab83ace676096e7506b6b7aa6953e87bc04a023ca016c030c31dd1c94062948ecdce2e67c4e6626c16af0033dcdb7a96362c937d48f842a00a95fac74aba7e3fbd24bc62457ce6981803d8f5fef28871d3d5e2af05d50cd4a0117400693917cd50d9bc28d4ab4fadf93a23e771f303637f8d1f83cd0632c3fcf888f842a0301bfced3253e99e8d50f2fed62313a16d714013d022a4dc4294656276f10d1ba0152e047a83c326a9d81dac502ec429b662b58ee119ca4c8748a355b539c24131f842a01944b5b4a3e93d46b0fe4370128c6cdcd066ae6b036b019a20f8d22fe9a10d67a00ddf3421722967c0bd965b9fc9e004bf01183b6206fec8de65e40331d185372ef842a02db8fb278708abf8878ebf578872ab35ee914ad8196b78de16b34498222ac1c2a02ff9d9a5184684f4e14530bde3a61a2f9adaa74734dff104b61ba3d963a644dac68207388208b7c68209998209c5c2c0c0820001";
        let raw_commitment = hex::decode(commitment_hex).unwrap();

        let commitment = StandardCommitment::from_rlp_bytes(raw_commitment.as_slice()).unwrap();
        let raw_from_bytes = commitment.to_rlp_bytes();

        assert_eq!(&raw_commitment, &raw_from_bytes);
    }

    #[test]
    fn fail_insufficient_data() {
        let raw_commitment = [];
        let commitment = StandardCommitment::from_rlp_bytes(raw_commitment.as_slice());

        assert!(matches!(
            &commitment,
            Err(StandardCommitmentParseError::EmptyCommitment),
        ));
    }

    #[test]
    fn fail_wrong_version() {
        let raw_commitment = [3, 3];
        let commitment = StandardCommitment::from_rlp_bytes(raw_commitment.as_slice());

        assert!(matches!(
            &commitment,
            Err(StandardCommitmentParseError::UnsupportedCertVersion(_)),
        ));
    }

    #[test]
    fn fail_invalid_rl_cert() {
        let raw_commitment = [2, 3, 3, 3, 3, 3, 3];
        let commitment = StandardCommitment::from_rlp_bytes(raw_commitment.as_slice());

        assert!(matches!(
            &commitment,
            Err(StandardCommitmentParseError::InvalidRlpCert(_)),
        ));
    }
}
