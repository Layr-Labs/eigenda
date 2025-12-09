use alloy_primitives::B256;
use reth_trie_common::proof::ProofVerificationError;
use thiserror::Error;

use crate::cert::StandardCommitmentParseError;
use crate::extraction::CertExtractionError;
use crate::verification::blob::error::BlobVerificationError;
use crate::verification::cert::error::CertVerificationError;

/// Errors that can occur during EigenDA verification.
#[derive(Error, Debug, PartialEq)]
pub enum EigenDaVerificationError {
    /// Transaction is not EIP1559
    #[error("Transaction is not EIP1559")]
    TxNotEip1559(B256),

    /// Standard commitment parse error
    #[error(transparent)]
    StandardCommitmentParseError(#[from] StandardCommitmentParseError),

    /// Certificate is too old relative to the current block (recency validation failure)
    #[error("The recency window was missed, inclusion_height ({0}), recency height ({1})")]
    RecencyWindowMissed(u64, u64),

    /// Certificate verification error
    #[error(transparent)]
    CertVerificationError(#[from] CertVerificationError),

    /// Proof verification error
    #[error(transparent)]
    ProofVerificationError(#[from] ProofVerificationError),

    /// Certificate extraction error
    #[error(transparent)]
    CertExtractionError(#[from] CertExtractionError),

    /// Certificate state missing for transaction.
    #[error("Certificate state is missing for transaction ({0})")]
    MissingCertState(B256),

    /// Blob missing for transaction.
    #[error("Blob missing for transaction ({0})")]
    MissingBlob(B256),

    /// Blob verification error
    #[error(transparent)]
    BlobVerificationError(#[from] BlobVerificationError),
}
