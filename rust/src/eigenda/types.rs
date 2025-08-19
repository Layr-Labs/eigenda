use alloy_primitives::{B256, Bytes};
use alloy_rlp::{Decodable, Encodable, Error};
use eigenda_cert::{
    BatchHeaderV2, BlobInclusionInfo, EigenDACertV2, EigenDACertV3, EigenDAVersionedCert,
    NonSignerStakesAndSignature,
};
use serde::{Deserialize, Serialize};
use thiserror::Error;

use crate::eigenda::verification::cert::convert;

#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct StandardCommitment(EigenDAVersionedCert);

/// Byte indicating a version 2 certificate.
const VERSION_2: u8 = 1;

/// Byte indicating a version 3 certificate.
const VERSION_3: u8 = 2;

impl StandardCommitment {
    pub fn from_rlp_bytes(bytes: &[u8]) -> Result<Self, StandardCommitmentParseError> {
        let (cert_version, mut cert_bytes) = bytes
            .split_first()
            .ok_or(StandardCommitmentParseError::InsufficientData)?;

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

    /// Get reference block used when constructing this certificate.
    pub fn reference_block(&self) -> u64 {
        match &self.0 {
            EigenDAVersionedCert::V2(c) => c.batch_header_v2.reference_block_number as u64,
            EigenDAVersionedCert::V3(c) => c.batch_header_v2.reference_block_number as u64,
        }
    }

    pub fn to_rlp_bytes(&self) -> Vec<u8> {
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

        bytes
    }

    pub fn relay_keys(&self) -> &[u32] {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => &cert.blob_inclusion_info.blob_certificate.relay_keys,
            EigenDAVersionedCert::V3(cert) => &cert.blob_inclusion_info.blob_certificate.relay_keys,
        }
    }

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

    pub fn signed_quorum_numbers(&self) -> &Bytes {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => &cert.signed_quorum_numbers,
            EigenDAVersionedCert::V3(cert) => &cert.signed_quorum_numbers,
        }
    }

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

    pub fn batch_header_v2(&self) -> &BatchHeaderV2 {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => &cert.batch_header_v2,
            EigenDAVersionedCert::V3(cert) => &cert.batch_header_v2,
        }
    }

    pub fn blob_inclusion_info(&self) -> &BlobInclusionInfo {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => &cert.blob_inclusion_info,
            EigenDAVersionedCert::V3(cert) => &cert.blob_inclusion_info,
        }
    }

    pub fn nonsigner_stake_and_signature(&self) -> &NonSignerStakesAndSignature {
        match &self.0 {
            EigenDAVersionedCert::V2(cert) => &cert.nonsigner_stake_and_signature,
            EigenDAVersionedCert::V3(cert) => &cert.nonsigner_stake_and_signature,
        }
    }
}

#[derive(Debug, Error)]
pub enum StandardCommitmentParseError {
    /// Invalid cert metadata
    #[error("Insufficient commitment data")]
    InsufficientData,
    /// Unsupported cer version
    #[error("Unsupported cert version {0}")]
    UnsupportedCertVersion(u8),
    /// The cert couldn't be parsed from the RPL format
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

        assert_eq!(&raw_commitment, raw_from_bytes.as_slice());
    }

    #[test]
    fn fail_insufficient_data() {
        let raw_commitment = [];
        let commitment = StandardCommitment::from_rlp_bytes(raw_commitment.as_slice());

        assert!(matches!(
            &commitment,
            Err(StandardCommitmentParseError::InsufficientData),
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
