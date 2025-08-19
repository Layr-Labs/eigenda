use alloy_consensus::{proofs::calculate_transaction_root, transaction::SignerRecoverable};
use alloy_primitives::B256;
use bytes::Bytes;
use reth_trie_common::proof::ProofVerificationError;
use serde::{Deserialize, Serialize};
use sov_rollup_interface::da::{
    BlobReaderTrait, BlockHeaderTrait, DaSpec, DaVerifier, RelevantBlobs, RelevantProofs,
};
use thiserror::Error;

use crate::{
    eigenda::verification::{
        blob::BlobVerificationError, verify_blob, verify_cert, verify_cert_recency,
    },
    ethereum::{extract_certificate, get_ancestor},
    spec::{
        AncestorMetadata, BlobWithSender, EigenDaSpec, EthereumAddress, EthereumBlockHeader,
        EthereumHash, NamespaceId, TransactionWithBlob,
    },
};

/// Errors that may occur when verifying with the [`EigenDaVerifier`].
#[derive(Debug, Error)]
pub enum VerifierError {
    #[error(transparent)]
    CompletenessError(#[from] CompletenessProofError),

    #[error(transparent)]
    InclusionError(#[from] InclusionProofError),
}

/// A verifier verifies that rollup transactions are indeed included in the
/// block, that no rollup tx from the block is missing and that blobs were
/// extracted correctly.
#[derive(Debug, Clone)]
pub struct EigenDaVerifier {
    batch_namespace: NamespaceId,
    proof_namespace: NamespaceId,
    cert_recency_window: u64,
}

impl EigenDaVerifier {
    pub fn verify_transactions(
        &self,
        block_header: &EthereumBlockHeader,
        namespace: NamespaceId,
        blobs_with_senders: &[BlobWithSender],
        completeness_proof: EigenDaCompletenessProof,
        inclusion_proof: EigenDaInclusionProof,
    ) -> Result<(), VerifierError> {
        // Verify completeness proof, proving that all transactions in a
        // specified set are in the block
        let (proven_transactions, proven_ancestors) = completeness_proof.verify(block_header)?;

        // Verify that the provided `relevant_blobs` are the only ones that
        // contain rollup data, and that they are correctly build from the block
        inclusion_proof.verify(
            block_header,
            namespace,
            &proven_transactions,
            &proven_ancestors,
            blobs_with_senders,
            self.cert_recency_window,
        )?;

        Ok(())
    }
}

impl DaVerifier for EigenDaVerifier {
    type Spec = EigenDaSpec;

    type Error = VerifierError;

    fn new(params: <Self::Spec as DaSpec>::ChainParams) -> Self {
        EigenDaVerifier {
            proof_namespace: params.rollup_proof_namespace,
            batch_namespace: params.rollup_batch_namespace,
            cert_recency_window: params.cert_recency_window,
        }
    }

    fn verify_relevant_tx_list(
        &self,
        block_header: &<Self::Spec as DaSpec>::BlockHeader,
        relevant_blobs: &RelevantBlobs<<Self::Spec as DaSpec>::BlobTransaction>,
        relevant_proofs: RelevantProofs<
            <Self::Spec as DaSpec>::InclusionMultiProof,
            <Self::Spec as DaSpec>::CompletenessProof,
        >,
    ) -> Result<(), Self::Error> {
        // Verify the `proof` namespace
        self.verify_transactions(
            block_header,
            self.proof_namespace,
            &relevant_blobs.proof_blobs,
            relevant_proofs.proof.completeness_proof,
            relevant_proofs.proof.inclusion_proof,
        )?;
        // Verify the `batch` namespace
        self.verify_transactions(
            block_header,
            self.batch_namespace,
            &relevant_blobs.batch_blobs,
            relevant_proofs.batch.completeness_proof,
            relevant_proofs.batch.inclusion_proof,
        )?;

        Ok(())
    }
}

/// Errors that may occur when verifying the [`EigenDaCompletenessProof`].
#[derive(Debug, Error)]
pub enum CompletenessProofError {
    #[error("Incorrect ancestry chain")]
    IncorrectAncestry,

    #[error("Recomputed merkle root ({0}) doesn't match the one in header ({1})")]
    MerkleRootMismatch(B256, B256),

    #[error("Error occurred while verifying the Ethereum account proof: {0}")]
    ProofVerificationError(#[from] ProofVerificationError),
}

/// A proof of completeness of the transactions in the block.
///
/// This proof holds all the transactions from the block, in order. It also
/// holds ancestors which are guaranteed to be contiguous and represent a direct
/// ancestor chain of the block.
#[derive(Debug, Serialize, Deserialize)]
pub struct EigenDaCompletenessProof {
    ancestors: Vec<AncestorMetadata>,
    transactions: Vec<TransactionWithBlob>,
}

impl EigenDaCompletenessProof {
    /// Create a new completeness proof.
    pub fn new(ancestors: Vec<AncestorMetadata>, transactions: Vec<TransactionWithBlob>) -> Self {
        Self {
            ancestors,
            transactions,
        }
    }

    /// Verify that the proof holds the complete list of transactions for the
    /// block. Also verify that the ancestors in the proof represent the correct
    /// ancestry chain and the block is a direct descendant.
    ///
    /// Upon success, the proof returns a vector of transactions that were
    /// proven to represent a whole transaction set of a specific block. It also
    /// returns a verified ancestor chain.
    ///
    /// # Errors
    ///
    /// This function will return an error if:
    ///   - ancestry chain is not continuous
    ///   - the header is not a direct descendant of the ancestry chain
    ///   - state data of the specific ancestor is incorrect
    ///   - recomputed transaction_root is different from the root in the header
    pub fn verify(
        self,
        header: &EthereumBlockHeader,
    ) -> Result<(Vec<TransactionWithBlob>, Vec<AncestorMetadata>), CompletenessProofError> {
        // Iterate through the ancestors and check their parent/child
        // relationships. If the direct ancestry is not valid, the error is
        // returned. In case when we have no ancestors the `last_ancestor` is None.
        let last_ancestor = self.ancestors.iter().try_fold(
            None,
            |parent: Option<&AncestorMetadata>, maybe_child| {
                // If there is no parent to check against, we set the ancestry as valid
                let valid_ancestry = parent
                    .map(|parent| parent.header.is_parent(&maybe_child.header))
                    .unwrap_or(true);

                if valid_ancestry {
                    Ok(Some(maybe_child))
                } else {
                    Err(CompletenessProofError::IncorrectAncestry)
                }
            },
        )?;

        // Check if last ancestor is our actual parent
        if let Some(ancestor) = last_ancestor
            && !ancestor.header.is_parent(header)
        {
            return Err(CompletenessProofError::IncorrectAncestry);
        }

        // Validate the data state of the ancestor blocks
        for ancestor in &self.ancestors {
            if let Some(data) = &ancestor.data {
                data.verify(ancestor.header.as_ref().state_root)?;
            }
        }

        // Verify that the proof holds the complete list of transactions for the block
        let transactions = self.transactions.iter().map(|t| &t.tx).collect::<Vec<_>>();
        let calculated_transactions_root = calculate_transaction_root(&transactions);
        if calculated_transactions_root != header.as_ref().transactions_root {
            return Err(CompletenessProofError::MerkleRootMismatch(
                calculated_transactions_root,
                header.as_ref().transactions_root,
            ));
        }

        Ok((self.transactions, self.ancestors))
    }
}

/// Errors that may occur when verifying the [`EigenDaInclusionProof`].
#[derive(Debug, Error)]
pub enum InclusionProofError {
    #[error("Transaction ({0}) in proof wasn't part of the completeness proof")]
    NotProvenTransaction(B256),

    #[error("Ancestor missing, height({0})")]
    AncestorMissing(u64),

    #[error("Proof incomplete, some relevant transactions are missing")]
    ProofIncomplete,

    #[error("Blob missing for transaction ({0})")]
    MissingBlob(B256),

    #[error("Additional blob ({0}) provided, irrelevant for the rollup")]
    IrrelevantBlob(B256),

    #[error("Transaction has incorrect sender ({1}), expected ({0})")]
    IncorrectSender(EthereumAddress, EthereumAddress),

    #[error("Transaction has incorrect hash ({1}), expected ({0})")]
    IncorrectBlobHash(EthereumHash, EthereumHash),

    #[error("Malformed data for transaction ({0})")]
    IncorrectBlobData(EthereumHash),

    #[error("Error occurred while verifying EigenDA blob: {0}")]
    BlobVerificationError(#[from] BlobVerificationError),

    #[error("Error occurred while tried to recover a transaction sender")]
    RecoverSenderError,
}

/// A proof of inclusion of the rollup transactions from the block.
///
/// This proof holds all transactions from the block which may be relevant to
/// the rollup, in order. The proof will check if the rollup's transactions are
/// the only one in block within rollup's namespace and if they are correctly
/// extracted from the block.
#[derive(Debug, Serialize, Deserialize)]
pub struct EigenDaInclusionProof {
    transactions: Vec<TransactionWithBlob>,
}

impl EigenDaInclusionProof {
    /// Create a new inclusion proof from a complete and ordered list of
    /// transactions that can be relevant for the rollup.
    pub fn new(maybe_relevant_txs: Vec<TransactionWithBlob>) -> Self {
        Self {
            transactions: maybe_relevant_txs,
        }
    }

    /// Verify that the proof holds transactions extracted from all transactions
    /// within given namespace, no more, no less.
    fn verify_transactions(
        &self,
        namespace: NamespaceId,
        proven_transactions: &[TransactionWithBlob],
    ) -> Result<(), InclusionProofError> {
        // Transaction hashes related to the transactions in the proof
        let mut transaction_hashes = self
            .transactions
            .iter()
            .map(|TransactionWithBlob { tx, .. }| tx.hash());

        // Transactions in a namespace extracted from the transactions that were
        // already proven to be a complete set.
        let mut namespace_transaction_hashes = proven_transactions
            .iter()
            .filter(|TransactionWithBlob { tx, .. }| namespace.contains(tx))
            .map(|TransactionWithBlob { tx, .. }| tx.hash());

        loop {
            match (
                transaction_hashes.next(),
                namespace_transaction_hashes.next(),
            ) {
                // Compare hash from proof with proven one
                (Some(hash), Some(proven_hash)) => {
                    if hash != proven_hash {
                        return Err(InclusionProofError::NotProvenTransaction(*hash));
                    }
                }
                // Extra transactions in proof
                (Some(hash), None) => {
                    return Err(InclusionProofError::NotProvenTransaction(*hash));
                }
                // Proof is missing a transaction
                (None, Some(_)) => return Err(InclusionProofError::ProofIncomplete),
                // Done
                (None, None) => return Ok(()),
            }
        }
    }

    /// Verify that all transactions with the valid certificates have a valid
    /// data blob. The transactions with an invalid certificates are ignored. If
    /// there is a single transaction with the valid certificate but invalid or
    /// missing data blob, the proof fails.
    fn verify_certs_and_blobs(
        &self,
        header: &EthereumBlockHeader,
        proven_ancestors: &[AncestorMetadata],
        cert_recency_window: u64,
    ) -> impl Iterator<Item = Result<(EthereumHash, EthereumAddress, Bytes), InclusionProofError>>
    {
        // Returning of iterator might be a bit convoluted. But it's nice
        // because we can skip having to allocate a temporary vector for
        // verified senders with blobs. The idea here is. Ignore all
        // transactions with the invalid certificates. If the certificate is
        // valid but the data blob is not, return a proof error. If both are
        // valid, construct a validated blob with sender.
        self.transactions
            .iter()
            .flat_map(move |TransactionWithBlob { tx, blob }| {
                // Skipping malformed cert
                let cert = extract_certificate(tx)?;

                // Skipping cert with failed recency check
                let referenced_height = cert.reference_block();
                if verify_cert_recency(header, referenced_height, cert_recency_window).is_err() {
                    return None;
                }

                let Some(ancestor) =
                    get_ancestor(proven_ancestors, header.height(), referenced_height)
                else {
                    return Some(Err(InclusionProofError::AncestorMissing(referenced_height)));
                };

                // Skipping invalid cert
                if verify_cert(header, ancestor, &cert).is_err() {
                    return None;
                }

                // The certificate is proven to be valid. The corresponding blob
                // should exist and it should be valid.
                let Some(blob) = blob.as_ref() else {
                    return Some(Err(InclusionProofError::MissingBlob(*tx.hash())));
                };
                if let Err(err) = verify_blob(&cert, blob) {
                    return Some(Err(err.into()));
                };

                // The relationship is valid
                let Ok(sender) = tx.recover_signer() else {
                    return Some(Err(InclusionProofError::RecoverSenderError));
                };
                let hash = EthereumHash::from(*tx.hash());
                let sender = EthereumAddress::from(sender);
                Some(Ok((hash, sender, blob.clone())))
            })
    }

    /// Verify that the given `blobs_with_senders` list form a complete set of
    /// rollup's batches in a single block, and that all supplied transactions
    /// are correctly extracted from that block.
    ///
    /// The proof accepts all transactions previously verified for completeness
    /// in the block. Based on that, the proof will verify if the list of its
    /// transactions forms complete namespace i.e. if there are no transactions
    /// in the block with the same namespace not included in proof.
    ///
    /// Then, since the proof holds all the transactions part of the namespace,
    /// we can check those transactions against the provided transactions in the
    /// `blobs_with_senders` list. Each transaction is checked for data and
    /// sender equality, in order they appear.
    ///
    /// # Errors
    ///
    /// This function will return an error if:
    ///   - there is an extra `BlobWithSender` provided which wasn't found in the block
    ///   - the list of provided `BlobWithSender`s is not complete, there are more in the block
    ///   - provided sender is different from the one in the proven transaction
    ///   - provided hash is different from the hash of the proven transaction
    ///   - provided blob data is different from the blob retrieved and proven to be part of the transaction
    pub fn verify(
        &self,
        header: &EthereumBlockHeader,
        namespace: NamespaceId,
        proven_transactions: &[TransactionWithBlob],
        proven_ancestors: &[AncestorMetadata],
        blobs_with_senders: &[BlobWithSender],
        cert_recency_window: u64,
    ) -> Result<(), InclusionProofError> {
        // Verify transactions contained by the proof
        self.verify_transactions(namespace, proven_transactions)?;
        // Verify certificates and blobs
        let mut valid_proven_blobs =
            self.verify_certs_and_blobs(header, proven_ancestors, cert_recency_window);

        let mut blobs_with_senders = blobs_with_senders.iter();

        // Compare proven transactions to the provided transactions represented
        // as blobs with senders
        loop {
            let (hash, sender, blob, provided) =
                match (valid_proven_blobs.next(), blobs_with_senders.next()) {
                    // We have a proven blob and some provided BlobWithSender
                    (Some(Ok((hash, sender, blob))), Some(provided)) => {
                        (hash, sender, blob, provided)
                    }
                    // `blob_with_sender` not provided
                    (Some(Ok((hash, ..))), None) => {
                        return Err(InclusionProofError::MissingBlob(*hash));
                    }
                    // The certificate/blob verification resulted in an error
                    (Some(Err(err)), _) => {
                        return Err(err);
                    }
                    // Missing transaction for the provided `blob_with_sender``
                    (None, Some(provided)) => {
                        return Err(InclusionProofError::IrrelevantBlob(*provided.hash()));
                    }
                    // We are finished
                    (None, None) => break,
                };

            // Check if sender is the same
            if sender != provided.sender {
                return Err(InclusionProofError::IncorrectSender(
                    sender,
                    provided.sender,
                ));
            }

            // Check if hash is the same
            if hash != provided.hash() {
                return Err(InclusionProofError::IncorrectBlobHash(
                    hash,
                    provided.hash(),
                ));
            }

            // Check if data read by the rollup was correct
            let consumed_data = provided.verified_data();
            if consumed_data.is_empty() {
                // Nothing to check
                continue;
            }
            if consumed_data.len() > blob.len() {
                // Provided transaction has more data than original one
                return Err(InclusionProofError::IncorrectBlobData(hash));
            }
            if consumed_data[..] != blob[..consumed_data.len()] {
                // Data mismatch
                return Err(InclusionProofError::IncorrectBlobData(hash));
            }
        }

        Ok(())
    }
}
