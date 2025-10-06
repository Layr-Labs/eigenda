//! Core EigenDA cryptographic verification primitives
//!
//! This module provides the fundamental cryptographic verification components for
//! EigenDA certificates and blob data. It implements the low-level verification
//! algorithms following the EigenDA protocol specification.
//!
//! ## Module Structure
//!
//! This crate contains the core verification primitives:
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
//! ## Architecture
//!
//! This module focuses on the cryptographic core of EigenDA verification and does
//! not handle:
//! - Ethereum state extraction and proof verification
//! - Rollup-specific integration logic
//! - Certificate recency validation (handled by higher-level adapters)
//!
//! The verification functions expect pre-validated inputs and focus purely on
//! cryptographic correctness.
//!
//! ## References
//!
//! - [EigenDA Protocol Specification](https://docs.eigenlayer.xyz/eigenda/overview/)
//! - [Certificate Verification Reference](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/integrations/cert/libraries/EigenDACertVerificationLib.sol)
//!
//! ## Verification Functions
//!
//! This module provides two key verification utilities:
//!
//! 1. **[`verify_cert_recency`]** - Certificate recency validation
//! 2. **[`verify_blob`]** - Blob commitment verification
//!
//! ### Certificate Recency Validation
//!
//! [`verify_cert_recency`] ensures certificates are used within an acceptable time window to prevent
//! stale certificate attacks. It validates that the certificate's reference block
//! is recent enough relative to the current inclusion block according to rollup parameters.
//!
//! ### Blob Verification
//!
//! [`verify_blob`] validates that blob data matches the cryptographic commitment in the certificate
//! using KZG polynomial commitments. This function serves as a convenient wrapper around
//! the core blob verification primitives.
//!
//! ## Integration with Other Crates
//!
//! This crate works together with:
//! - `eigenda_ethereum` for extracting Ethereum state data needed for verification
//! - `sov_eigenda_adapter` for rollup-specific integration logic
//!
//! ## References
//!
//! - [EigenDA Integration Specification](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html)
//! - [Sovereign SDK Documentation](https://docs.sovereign.xyz/)

use tracing::instrument;

use crate::cert::StandardCommitment;
use crate::verification::blob::error::BlobVerificationError;
use crate::verification::cert::error::CertVerificationError;

pub mod blob;
pub mod cert;

/// Validate certificate recency to prevent stale certificate attacks
///
/// Ensures that the certificate's reference block is recent enough relative to
/// the inclusion block. This prevents attackers from using old certificates
/// with outdated operator sets.
///
/// # Arguments
/// * `header` - Ethereum block header where the certificate is being included
/// * `referenced_height` - Block height referenced by the certificate
/// * `cert_recency_window` - Maximum allowed age of the certificate in blocks
///
/// # Returns
/// `Ok(())` if the certificate is within the recency window
///
/// # Errors
/// Returns [`CertVerificationError::RecencyWindowMissed`] if the certificate
/// is too old relative to the inclusion block.
///
/// # Reference
/// [EigenDA Specification - RBN Recency Validation](https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#1-rbn-recency-validation)
pub fn verify_cert_recency(
    inclusion_height: u64,
    referenced_height: u64,
    cert_recency_window: u64,
) -> Result<(), CertVerificationError> {
    let recency_height = referenced_height + cert_recency_window;
    if inclusion_height > recency_height {
        return Err(CertVerificationError::RecencyWindowMissed(
            inclusion_height,
            recency_height,
        ));
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
    use crate::verification::cert::error::CertVerificationError;
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
                CertVerificationError::RecencyWindowMissed(
                    inclusion_height,
                    referenced_height + cert_recency_window
                ),
                "{description}"
            );
        }
    }
}
