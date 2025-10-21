//! Core EigenDA cryptographic verification primitives
//!
//! This module provides the fundamental cryptographic verification components for
//! EigenDA certificates and blob data. It implements the low-level verification
//! algorithms following the EigenDA protocol specification.
//!
//! ## Module Structure
//!
//! This module contains the core verification primitives:
//!
//! - **[`cert`]** - Certificate cryptographic verification
//!   - BLS signature aggregation and verification
//!   - Stake-weighted quorum validation
//!   - Security threshold enforcement
//!   - Operator state consistency checks
//!
//! - **[`blob`]** - Blob data integrity verification
//!   - KZG polynomial commitment verification
//!   - Blob encoding validation
//!
//! ## High-Level API
//!
//! This module provides convenient high-level functions for common verification workflows:
//!
//! - **[`extract_certificate`]** - Extracts an EigenDA certificate from an EIP-4844 transaction
//! - **[`verify_and_extract_blob`]** - All-in-one verification: recency, certificate, and blob extraction
//! - **[`verify_cert_recency`]** - Certificate recency validation to prevent stale certificate attacks
//! - **[`verify_blob`]** - Blob commitment verification using KZG proofs
//!
//! ## Low-Level API
//!
//! For fine-grained control, use the submodules directly:
//! - [`cert::verify`] - Certificate-only verification with extracted state data
//! - [`blob::verify_blob`] - Blob-only verification
//!
//! ## Integration with Other Modules
//!
//! This module works together with:
//! - [`crate::extraction`] - For extracting and verifying Ethereum contract state
//! - [`crate::error`] - For unified error handling across verification operations
//!
//! ## References
//!
//! - [EigenDA Protocol Specification](https://docs.eigenlayer.xyz/eigenda/overview/)
//! - [Certificate Verification Reference](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/integrations/cert/libraries/EigenDACertVerificationLib.sol)
//! - [EigenDA Integration Specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html)
//! - [Sovereign SDK Documentation](https://docs.sovereign.xyz/)

use alloy_consensus::{EthereumTxEnvelope, Header, Transaction, TxEip4844};
use alloy_primitives::B256;
use bytes::Bytes;
use tracing::instrument;

use crate::cert::StandardCommitment;
use crate::error::EigenDaVerificationError;
use crate::extraction::CertStateData;
use crate::verification::blob::codec::decode_encoded_payload;
use crate::verification::blob::error::BlobVerificationError;

/// Blob integrity verification using KZG polynomial commitments.
pub mod blob;
/// Certificate cryptographic verification using BLS signatures.
pub mod cert;

/// Extracts an EigenDA certificate from an EIP-4844 transaction.
///
/// Parses the transaction input data to extract a [`StandardCommitment`] certificate.
///
/// # Arguments
/// * `tx` - EIP-4844 transaction envelope containing the certificate
///
/// # Returns
/// The parsed [`StandardCommitment`] certificate
///
/// # Errors
/// - [`EigenDaVerificationError::TxNotEip1559`] if transaction is not EIP-1559 format
/// - [`EigenDaVerificationError::StandardCommitmentParseError`] if certificate parsing fails
pub fn extract_certificate(
    tx: &EthereumTxEnvelope<TxEip4844>,
) -> Result<StandardCommitment, EigenDaVerificationError> {
    use EigenDaVerificationError::*;

    let signed_tx = tx.as_eip1559().ok_or_else(|| TxNotEip1559(*tx.hash()))?;
    let rlp_bytes = signed_tx.input();
    let cert = StandardCommitment::from_rlp_bytes(rlp_bytes)?;
    Ok(cert)
}

/// Verifies an EigenDA certificate and extracts the blob data.
///
/// This is a high-level convenience function that performs the complete verification workflow:
/// 1. Validates certificate recency
/// 2. Verifies contract state proofs against the state root
/// 3. Extracts verification inputs from proven state
/// 4. Verifies certificate cryptographically (BLS signatures, quorum stakes, thresholds)
/// 5. Verifies blob data matches the certificate commitment (KZG proof)
/// 6. Decodes and returns the blob payload
///
/// # Arguments
/// * `tx` - Transaction hash (for error reporting)
/// * `cert` - The certificate to verify
/// * `cert_state` - Optional contract state data with proofs
/// * `cert_state_header` - Block header containing the state root for verification
/// * `inclusion_height` - Block height where certificate is included
/// * `referenced_height` - Block height referenced by the certificate
/// * `cert_recency_window` - Maximum allowed certificate age in blocks
/// * `encoded_payload` - Optional encoded blob payload to verify
///
/// # Returns
/// Decoded blob data as [`Bytes`]
///
/// # Errors
/// - [`EigenDaVerificationError::MissingCertState`] if cert_state is None
/// - [`EigenDaVerificationError::ProofVerificationError`] if state proofs are invalid
/// - [`EigenDaVerificationError::MissingBlob`] if encoded_payload is None
/// - [`EigenDaVerificationError::BlobVerificationError`] if blob verification fails
#[allow(clippy::too_many_arguments)]
pub fn verify_and_extract_payload(
    tx: B256,
    cert: &StandardCommitment,
    cert_state: &Option<CertStateData>,
    cert_state_header: &Header,
    inclusion_height: u64,
    referenced_height: u64,
    cert_recency_window: u64,
    encoded_payload: &Option<Bytes>,
) -> Option<Result<Bytes, EigenDaVerificationError>> {
    use EigenDaVerificationError::*;

    // if certificate recency verification fails: ignore
    verify_cert_recency(inclusion_height, referenced_height, cert_recency_window).ok()?;

    let cert_state = match cert_state.as_ref() {
        Some(cert_state) => cert_state,
        None => return Some(Err(MissingCertState(tx))),
    };

    if let Err(err) = cert_state.verify(cert_state_header.state_root) {
        return Some(Err(ProofVerificationError(err)));
    }

    // if certificate extraction fails: ignore
    let current_block = inclusion_height as u32;
    let inputs = cert_state.extract(cert, current_block).ok()?;

    // if certificate verification fails: ignore
    cert::verify(inputs).ok()?;

    let encoded_payload = match encoded_payload.as_ref() {
        Some(encoded_payload) => encoded_payload,
        None => return Some(Err(MissingBlob(tx))),
    };

    if let Err(err) = verify_blob(cert, encoded_payload) {
        return Some(Err(BlobVerificationError(err)));
    }

    // if encoded_payload decode fails: ignore
    let payload = decode_encoded_payload(encoded_payload).ok()?;

    Some(Ok(Bytes::from(payload)))
}

/// Validate certificate recency to prevent stale certificate attacks
///
/// Ensures that the certificate's reference block is recent enough relative to
/// the inclusion block. This prevents attackers from using old certificates
/// with outdated operator sets.
///
/// # Arguments
/// * `inclusion_height` - Block height where the certificate is being included
/// * `referenced_height` - Block height referenced by the certificate
/// * `cert_recency_window` - Maximum allowed age of the certificate in blocks
///
/// # Returns
/// `Ok(())` if the certificate is within the recency window
///
/// # Errors
/// Returns [`EigenDaVerificationError::RecencyWindowMissed`] if the certificate
/// is too old relative to the inclusion block.
///
/// # Reference
/// [EigenDA Specification - RBN Recency Validation](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#1-rbn-recency-validation)
pub fn verify_cert_recency(
    inclusion_height: u64,
    referenced_height: u64,
    cert_recency_window: u64,
) -> Result<(), EigenDaVerificationError> {
    use EigenDaVerificationError::*;

    let recency_height = referenced_height + cert_recency_window;
    if inclusion_height > recency_height {
        return Err(RecencyWindowMissed(inclusion_height, recency_height));
    }

    Ok(())
}

/// Validate encoded payload against certificate commitment
///
/// Verifies that the provided encoded payload matches the cryptographic
/// commitment contained in the certificate using KZG polynomial commitments.
///
/// # Arguments
/// * `cert` - Certificate containing the blob commitment
/// * `encoded_payload` - Encoded payload to validate
///
/// # Returns
/// `Ok(())` if the encoded payload matches the certificate commitment
///
/// # Errors
/// Returns [`BlobVerificationError`] if:
/// - Encoded payload doesn't match the commitment
/// - KZG proof verification fails
/// - Commitment is malformed
///
/// # Reference
/// [EigenDA Specification - Blob Validation](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation)
#[instrument(skip_all)]
pub fn verify_blob(
    cert: &StandardCommitment,
    encoded_payload: &[u8],
) -> Result<(), BlobVerificationError> {
    let blob_commitment = &cert
        .blob_inclusion_info()
        .blob_certificate
        .blob_header
        .commitment;

    blob::verify(blob_commitment, encoded_payload)
}

#[cfg(test)]
mod tests {
    // use alloy_consensus::Header;
    use crate::error::EigenDaVerificationError;
    use crate::verification::verify_cert_recency;

    #[test]
    fn verify_cert_recency_success_cases() {
        // Test cases: (description, referenced_height, cert_recency_window, inclusion_height_offset)
        let test_cases = [
            ("exactly at window boundary", 100, 50, 50),
            ("well within window", 100, 50, 40),
            ("same block as reference", 100, 50, 0),
            ("zero window success", 100, 0, 0),
            ("large window", 1000, u64::MAX - 1000, 1000),
            ("edge case max values", u64::MAX - 100, 50, 25),
        ];

        for (description, referenced_height, cert_recency_window, inclusion_offset) in test_cases {
            let inclusion_height = referenced_height + inclusion_offset;
            let result =
                verify_cert_recency(inclusion_height, referenced_height, cert_recency_window);
            assert_eq!(result, Ok(()), "{description}");
        }
    }

    #[test]
    fn verify_cert_recency_failure_cases() {
        // Test cases: (description, referenced_height, cert_recency_window, inclusion_height_offset)
        let test_cases = [
            ("one block past window", 100, 50, 51),
            ("far past window", 100, 50, 150),
            ("zero window failure", 100, 0, 1),
        ];

        for (description, referenced_height, cert_recency_window, inclusion_offset) in test_cases {
            let inclusion_height = referenced_height + inclusion_offset;
            let err = verify_cert_recency(inclusion_height, referenced_height, cert_recency_window)
                .unwrap_err();

            assert_eq!(
                err,
                EigenDaVerificationError::RecencyWindowMissed(
                    inclusion_height,
                    referenced_height + cert_recency_window
                ),
                "{description}"
            );
        }
    }
}
