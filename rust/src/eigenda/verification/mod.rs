pub mod blob;
pub mod cert;

use sov_rollup_interface::da::BlockHeaderTrait;

use crate::{
    eigenda::{
        types::StandardCommitment,
        verification::{blob::BlobVerificationError, cert::error::CertVerificationError},
    },
    spec::{AncestorMetadata, EthereumBlockHeader},
};

/// Certificate recency validation
///
/// https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#1-rbn-recency-validation
pub fn verify_cert_recency(
    header: &EthereumBlockHeader,
    referenced_height: u64,
    cert_recency_window: u64,
) -> Result<(), CertVerificationError> {
    let inclusion_height = header.height();

    let recency_height = referenced_height + cert_recency_window;
    if inclusion_height > recency_height {
        return Err(CertVerificationError::RecencyWindowMissed(
            inclusion_height,
            recency_height,
        ));
    }

    Ok(())
}

/// Certificate validation
///
/// https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#2-cert-validation
pub fn verify_cert(
    header: &EthereumBlockHeader,
    ancestor: &AncestorMetadata,
    cert: &StandardCommitment,
) -> Result<(), CertVerificationError> {
    let current_block = header.height() as u32; // sol does `uint32(block.number)`
    let ancestor_data = ancestor
        .data
        .as_ref()
        .ok_or(CertVerificationError::AncestorDataMissing)?;

    let inputs = ancestor_data.extract(cert, current_block)?;

    cert::verify(inputs)?;
    Ok(())
}

/// Blob validation against the certificate
///
/// https://layr-labs.github.io/eigenda/integration/spec/6-secure-integration.html#3-blob-validation
pub fn verify_blob(
    _certificate: &StandardCommitment,
    _blob: &[u8],
) -> Result<(), BlobVerificationError> {
    // TODO: Verify the blob against the certificate. Doing that we have a
    // full validated chain of data

    blob::verify();

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::spec::EthereumBlockHeader;
    use alloy_consensus::Header;

    /// Helper function to create a mock EthereumBlockHeader with a given height
    fn create_mock_header(height: u64) -> EthereumBlockHeader {
        let mut header = Header::default();
        header.number = height;
        EthereumBlockHeader::from(header)
    }

    #[test]
    fn test_verify_cert_recency_success_within_window() {
        let referenced_height = 100;
        let cert_recency_window = 50;
        let inclusion_height = referenced_height + cert_recency_window; // exactly at the window boundary

        let header = create_mock_header(inclusion_height);

        let result = verify_cert_recency(&header, referenced_height, cert_recency_window);
        assert!(result.is_ok());
    }

    #[test]
    fn test_verify_cert_recency_success_before_window_expires() {
        let referenced_height = 100;
        let cert_recency_window = 50;
        let inclusion_height = referenced_height + cert_recency_window - 10; // well within the window

        let header = create_mock_header(inclusion_height);

        let result = verify_cert_recency(&header, referenced_height, cert_recency_window);
        assert!(result.is_ok());
    }

    #[test]
    fn test_verify_cert_recency_success_same_block() {
        let referenced_height = 100;
        let cert_recency_window = 50;
        let inclusion_height = referenced_height; // included in the same block as reference

        let header = create_mock_header(inclusion_height);

        let result = verify_cert_recency(&header, referenced_height, cert_recency_window);
        assert!(result.is_ok());
    }

    #[test]
    fn test_verify_cert_recency_failure_window_missed() {
        let referenced_height = 100;
        let cert_recency_window = 50;
        let inclusion_height = referenced_height + cert_recency_window + 1; // one block past the window

        let header = create_mock_header(inclusion_height);

        let result = verify_cert_recency(&header, referenced_height, cert_recency_window);
        assert!(result.is_err());

        if let Err(CertVerificationError::RecencyWindowMissed(actual_inclusion, expected_recency)) =
            result
        {
            assert_eq!(actual_inclusion, inclusion_height);
            assert_eq!(expected_recency, referenced_height + cert_recency_window);
        } else {
            panic!("Expected RecencyWindowMissed error");
        }
    }

    #[test]
    fn test_verify_cert_recency_failure_far_past_window() {
        let referenced_height = 100;
        let cert_recency_window = 50;
        let inclusion_height = referenced_height + cert_recency_window + 100; // far past the window

        let header = create_mock_header(inclusion_height);

        let result = verify_cert_recency(&header, referenced_height, cert_recency_window);
        assert!(result.is_err());

        if let Err(CertVerificationError::RecencyWindowMissed(actual_inclusion, expected_recency)) =
            result
        {
            assert_eq!(actual_inclusion, inclusion_height);
            assert_eq!(expected_recency, referenced_height + cert_recency_window);
        } else {
            panic!("Expected RecencyWindowMissed error");
        }
    }

    #[test]
    fn test_verify_cert_recency_zero_window() {
        let referenced_height = 100;
        let cert_recency_window = 0;
        let inclusion_height = referenced_height; // must be included in the same block

        let header = create_mock_header(inclusion_height);

        let result = verify_cert_recency(&header, referenced_height, cert_recency_window);
        assert!(result.is_ok());
    }

    #[test]
    fn test_verify_cert_recency_zero_window_failure() {
        let referenced_height = 100;
        let cert_recency_window = 0;
        let inclusion_height = referenced_height + 1; // one block past zero window

        let header = create_mock_header(inclusion_height);

        let result = verify_cert_recency(&header, referenced_height, cert_recency_window);
        assert!(result.is_err());

        if let Err(CertVerificationError::RecencyWindowMissed(actual_inclusion, expected_recency)) =
            result
        {
            assert_eq!(actual_inclusion, inclusion_height);
            assert_eq!(expected_recency, referenced_height + cert_recency_window);
        } else {
            panic!("Expected RecencyWindowMissed error");
        }
    }

    #[test]
    fn test_verify_cert_recency_large_window() {
        let referenced_height = 1000;
        let cert_recency_window = u64::MAX - referenced_height; // maximum possible window
        let inclusion_height = referenced_height + 1000; // well within the large window

        let header = create_mock_header(inclusion_height);

        let result = verify_cert_recency(&header, referenced_height, cert_recency_window);
        assert!(result.is_ok());
    }

    #[test]
    fn test_verify_cert_recency_edge_case_max_values() {
        let referenced_height = u64::MAX - 100;
        let cert_recency_window = 50;
        let inclusion_height = referenced_height + 25; // within window

        let header = create_mock_header(inclusion_height);

        let result = verify_cert_recency(&header, referenced_height, cert_recency_window);
        assert!(result.is_ok());
    }
}
