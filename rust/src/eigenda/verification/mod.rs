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
    cert: &StandardCommitment,
    cert_recency_window: u64,
) -> Result<(), CertVerificationError> {
    let referenced_height = cert.reference_block();
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
