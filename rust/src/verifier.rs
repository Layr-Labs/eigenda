use alloy_consensus::proofs::calculate_transaction_root;
use alloy_consensus::transaction::SignerRecoverable;
use alloy_consensus::{EthereumTxEnvelope, TxEip4844};
use alloy_primitives::B256;
use bytes::Bytes;
use reth_trie_common::proof::ProofVerificationError;
use serde::{Deserialize, Serialize};
use serde_with::serde_as;
use sov_rollup_interface::da::{
    BlobReaderTrait, BlockHeaderTrait, DaSpec, DaVerifier, RelevantBlobs, RelevantProofs,
};
use thiserror::Error;
use tracing::instrument;

use crate::eigenda::verification::blob::codec::decode_encoded_payload;
use crate::eigenda::verification::blob::error::BlobVerificationError;
use crate::eigenda::verification::{verify_blob, verify_cert, verify_cert_recency};
use crate::ethereum::extract_certificate;
use crate::spec::{
    BlobWithSender, CertificateStateData, EigenDaSpec, EthereumAddress, EthereumBlockHeader,
    EthereumHash, NamespaceId, TransactionWithBlob,
};

/// Errors that may occur when verifying with the [`EigenDaVerifier`].
#[derive(Debug, Error)]
#[allow(clippy::large_enum_variant)]
pub enum VerifierError {
    #[error(transparent)]
    /// Error verifying completeness proof.
    CompletenessError(#[from] CompletenessProofError),

    #[error(transparent)]
    /// Error verifying inclusion proof.
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
    #![allow(clippy::result_large_err)]
    #[instrument(skip_all, fields(block_height = block_header.height()))]
    /// Verify that provided blobs are correctly included in the block.
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
        let proven_transactions = completeness_proof.verify(block_header)?;

        // Verify that the provided `blobs_with_senders` are the only ones that
        // contain rollup data, and that they are correctly built from the block
        inclusion_proof.verify(
            block_header,
            namespace,
            &proven_transactions,
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

    #[instrument(skip_all, fields(block_height = block_header.height()))]
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
    #[error("Recomputed merkle root ({0}) doesn't match the one in header ({1})")]
    /// Merkle root mismatch between computed and header values.
    MerkleRootMismatch(B256, B256),
}

/// A proof of completeness of the transactions in the block.
///
/// This proof holds all the transactions from the block, in order.
#[serde_as]
#[derive(Debug, Serialize, Deserialize)]
pub struct EigenDaCompletenessProof {
    /// serde_as is required due to risc0 serialization.
    #[serde_as(as = "Vec<alloy_consensus::serde_bincode_compat::EthereumTxEnvelope<'_>>")]
    transactions: Vec<EthereumTxEnvelope<TxEip4844>>,
}

impl EigenDaCompletenessProof {
    /// Create a new completeness proof.
    pub fn new(transactions: Vec<EthereumTxEnvelope<TxEip4844>>) -> Self {
        Self { transactions }
    }

    /// Verify that the proof holds the complete list of transactions for the
    /// block.
    ///
    /// Upon success, the proof returns a vector of transactions that were
    /// proven to represent a whole transaction set of a specific block.
    ///
    /// # Errors
    ///
    /// This function will return an error if:
    ///   - Recomputed transaction_root is different from the root in the header
    #[allow(clippy::result_large_err)]
    pub fn verify(
        self,
        header: &EthereumBlockHeader,
    ) -> Result<Vec<EthereumTxEnvelope<TxEip4844>>, CompletenessProofError> {
        // Verify that the proof holds the complete list of transactions for the block
        let calculated_transactions_root = calculate_transaction_root(&self.transactions);
        if calculated_transactions_root != header.as_ref().transactions_root {
            return Err(CompletenessProofError::MerkleRootMismatch(
                calculated_transactions_root,
                header.as_ref().transactions_root,
            ));
        }

        Ok(self.transactions)
    }
}

/// Errors that may occur when verifying the [`EigenDaInclusionProof`].
#[derive(Debug, Error)]
pub enum InclusionProofError {
    #[error("Transaction ({0}) in proof wasn't part of the completeness proof")]
    /// Transaction in proof wasn't part of completeness proof.
    NotProvenTransaction(B256),

    /// The ancestry chain is incorrect
    #[error("Incorrect ancestry chain")]
    IncorrectAncestry,

    #[error("Certificate state is missing for transaction ({0})")]
    /// Certificate state missing for transaction.
    CertStateMissing(B256),

    #[error("Proof incomplete, some relevant transactions are missing")]
    /// Proof incomplete, some relevant transactions missing.
    ProofIncomplete,

    #[error("Blob missing for transaction ({0})")]
    /// Blob missing for transaction.
    MissingBlob(B256),

    #[error("Additional blob ({0}) provided, irrelevant for the rollup")]
    /// Additional blob provided that's irrelevant for rollup.
    IrrelevantBlob(B256),

    #[error("Transaction has incorrect sender ({1}), expected ({0})")]
    /// Transaction has incorrect sender address.
    IncorrectSender(EthereumAddress, EthereumAddress),

    #[error("Transaction has incorrect hash ({1}), expected ({0})")]
    /// Transaction has incorrect blob hash.
    IncorrectBlobHash(EthereumHash, EthereumHash),

    #[error("Malformed data for transaction ({0})")]
    /// Malformed blob data for transaction.
    IncorrectBlobData(EthereumHash),

    #[error("Error occurred while verifying EigenDA blob: {0}")]
    /// Error verifying EigenDA blob.
    BlobVerificationError(#[from] BlobVerificationError),

    #[error("Error occurred while tried to recover a transaction sender")]
    /// Error recovering transaction sender.
    RecoverSenderError,

    #[error("Error occurred while verifying the Ethereum account proof: {0}")]
    /// Error verifying Ethereum account proof.
    ProofVerificationError(#[from] ProofVerificationError),
}

/// A proof of inclusion of the rollup transactions from the block.
///
/// This proof holds all transactions from the block which may be relevant to
/// the rollup, in order. The proof will check if the rollup's transactions are
/// the only one in block within rollup's namespace and if they are correctly
/// extracted from the block.
#[derive(Debug, Serialize, Deserialize)]
pub struct EigenDaInclusionProof {
    #[cfg(feature = "use-rbn-state")]
    ancestors: Vec<EthereumBlockHeader>,
    transactions: Vec<TransactionWithBlob>,
}

impl EigenDaInclusionProof {
    /// Create a new inclusion proof from a complete and ordered list of
    /// transactions that can be relevant for the rollup.
    pub fn new(
        #[cfg(feature = "use-rbn-state")] ancestors: Vec<EthereumBlockHeader>,
        maybe_relevant_txs: Vec<TransactionWithBlob>,
    ) -> Self {
        Self {
            #[cfg(feature = "use-rbn-state")]
            ancestors,
            transactions: maybe_relevant_txs,
        }
    }

    /// Get the [`EthereumBlockHeader`] for the specific referenced block. The
    /// `ancestors` are expected to be a contiguous chain of ancestors preceding the
    /// `current_height`.
    #[cfg(feature = "use-rbn-state")]
    #[instrument(skip_all, fields(block_height = current_height))]
    pub fn get_ancestor(
        &self,
        current_height: u64,
        referenced_height: u64,
    ) -> Option<&EthereumBlockHeader> {
        // Check that the referenced height is always smaller from the current_height
        if current_height <= referenced_height {
            return None;
        }

        // Safety: We know that the referenced_height is always smaller from current_height.
        let diff = current_height - referenced_height;
        let ancestors_len = self.ancestors.len() as u64;

        // Check that the referenced height is in the vector
        if ancestors_len < diff {
            return None;
        }

        // Safety: We know that the `diff` <= `ancestors_len`
        let index = (ancestors_len - diff) as usize;
        Some(&self.ancestors[index])
    }

    /// Verify that the ancestors in the proof represent the correct ancestry
    /// chain and the block header is a direct descendant.
    ///
    /// # Errors
    ///
    /// This function will return an error if:
    ///   - Ancestry chain is not continuous
    ///   - The header is not a direct descendant of the ancestry chain
    #[cfg(feature = "use-rbn-state")]
    fn verify_ancestors(&self, header: &EthereumBlockHeader) -> Result<(), InclusionProofError> {
        // Iterate through the ancestors and check their parent/child
        // relationships. If the direct ancestry is not valid, the error is
        // returned. In case when we have no ancestors the `last_ancestor` is None.
        let last_ancestor = self.ancestors.iter().try_fold(
            None,
            |parent: Option<&EthereumBlockHeader>, maybe_child| {
                // If there is no parent to check against, we set the ancestry as valid
                let valid_ancestry = parent
                    .map(|parent| parent.is_parent(maybe_child))
                    .unwrap_or(true);

                if valid_ancestry {
                    Ok(Some(maybe_child))
                } else {
                    Err(InclusionProofError::IncorrectAncestry)
                }
            },
        )?;

        // Check if last ancestor is our actual parent
        if let Some(ancestor) = last_ancestor
            && !ancestor.is_parent(header)
        {
            return Err(InclusionProofError::IncorrectAncestry);
        }

        Ok(())
    }

    fn verify_cert_state(
        &self,
        header: &EthereumBlockHeader,
        #[cfg(feature = "use-rbn-state")] referenced_height: u64,
        state: &CertificateStateData,
    ) -> Result<(), InclusionProofError> {
        // Verify the certificate state against the referenced block header state root
        #[cfg(feature = "use-rbn-state")]
        {
            let current_height = header.height();
            let rbn_header = self
                .get_ancestor(current_height, referenced_height)
                .ok_or(InclusionProofError::IncorrectAncestry)?;

            state.verify(rbn_header.as_ref().state_root)?;
        }

        // Verify the certificate state against the current block header state root
        #[cfg(not(feature = "use-rbn-state"))]
        {
            state.verify(header.as_ref().state_root)?;
        }

        Ok(())
    }

    /// Verify that the proof holds transactions extracted from all transactions
    /// within given namespace, no more, no less.
    /// Also verify that the certificate states are correct./
    ///
    /// # Errors
    ///
    /// This function will return an error if:
    ///   - Certificate state data is incorrect
    ///   - Transaction in proof wasn't part of the completeness proof
    ///   - Proof is incomplete, some relevant transactions are missing
    #[instrument(skip_all)]
    fn verify_transactions(
        &self,
        namespace: NamespaceId,
        proven_transactions: &[EthereumTxEnvelope<TxEip4844>],
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
            .filter(|tx| namespace.contains(*tx))
            .map(|tx| tx.hash());

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
    #[instrument(skip_all, fields(block_height = header.height()))]
    fn verify_certs_and_blobs(
        &self,
        header: &EthereumBlockHeader,
        cert_recency_window: u64,
    ) -> impl Iterator<Item = Result<(EthereumHash, EthereumAddress, Bytes), InclusionProofError>>
    {
        // Returning of iterator might be a bit convoluted. But it's nice
        // because we can skip having to allocate a temporary vector for
        // verified senders with blobs. The idea here is:
        // - Ignore all transactions with the invalid certificates.
        // - If the certificate is valid but the data blob is not, return a proof error.
        // - If both are valid, construct a validated blob with sender.
        self.transactions.iter().flat_map(
            move |TransactionWithBlob {
                      tx,
                      encoded_payload,
                      cert_state,
                  }| {
                // Skipping malformed cert
                let cert = extract_certificate(tx)?;

                // Check the recency. Skipping certs with failed recency check
                let referenced_height = cert.reference_block();
                verify_cert_recency(header, referenced_height, cert_recency_window).ok()?;

                // State should be set, so we can verify the cert
                let Some(state) = cert_state.as_ref() else {
                    return Some(Err(InclusionProofError::CertStateMissing(*tx.hash())));
                };

                // Verify the cert state. The state should always be valid
                if let Err(err) = self.verify_cert_state(
                    header,
                    #[cfg(feature = "use-rbn-state")]
                    referenced_height,
                    state,
                ) {
                    return Some(Err(err));
                };

                // Verify the cert. We are skipping it if it's invalid.
                verify_cert(header, state, &cert).ok()?;

                // The encoded payload should be set
                let Some(encoded_payload) = encoded_payload.as_ref() else {
                    return Some(Err(InclusionProofError::MissingBlob(*tx.hash())));
                };

                // Verify the encoded payload against the certificate
                if let Err(err) = verify_blob(&cert, encoded_payload) {
                    return Some(Err(err.into()));
                };

                // Decode an encoded payload. The blob is dropped if it can't be decoded.
                let blob = decode_encoded_payload(encoded_payload).ok()?;

                // Recover the signer of the transaction
                let Ok(sender) = tx.recover_signer() else {
                    return Some(Err(InclusionProofError::RecoverSenderError));
                };

                let hash = EthereumHash::from(*tx.hash());
                let sender = EthereumAddress::from(sender);
                Some(Ok((hash, sender, Bytes::from(blob))))
            },
        )
    }

    /// Verify that the given `blobs_with_senders` list form a complete set of
    /// rollup's batches in a single block, and that all supplied transactions
    /// are correctly extracted from that block.
    ///
    /// The proof accepts all transactions previously verified for completeness
    /// in the block. Based on that, the proof will verify if the list of its
    /// transactions forms a complete namespace i.e. if there are no transactions
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
    ///   - Ancestry chain contained by the proof is incorrect.
    ///   - There is an extra `BlobWithSender` provided which wasn't found in the block.
    ///   - The list of provided `BlobWithSender`s is not complete, there are more in the block.
    ///   - Provided sender is different from the one in the proven transaction.
    ///   - Provided hash is different from the hash of the proven transaction.
    ///   - Provided blob data is different from the blob retrieved and proven to be part of the transaction.
    #[instrument(skip_all, fields(block_height = header.height()))]
    pub fn verify(
        &self,
        header: &EthereumBlockHeader,
        namespace: NamespaceId,
        proven_transactions: &[EthereumTxEnvelope<TxEip4844>],
        blobs_with_senders: &[BlobWithSender],
        cert_recency_window: u64,
    ) -> Result<(), InclusionProofError> {
        // Verify ancestry chain contained by the proof
        #[cfg(feature = "use-rbn-state")]
        self.verify_ancestors(header)?;
        // Verify transactions contained by the proof
        self.verify_transactions(namespace, proven_transactions)?;
        // Verify certificates and blobs
        let mut valid_proven_blobs = self.verify_certs_and_blobs(header, cert_recency_window);

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
            if sender != provided.sender() {
                return Err(InclusionProofError::IncorrectSender(
                    sender,
                    provided.sender(),
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

#[cfg(test)]
mod tests {
    use std::str::FromStr;

    use alloy_consensus::{EthereumTxEnvelope, Header, SignableTransaction, TxEip1559, TxEnvelope};
    use alloy_primitives::{TxKind, address};
    use alloy_signer::Signature;
    use bytes::Bytes;

    use super::*;
    use crate::spec::EthereumBlockHeader;

    pub fn create_test_header(
        number: u64,
        parent_hash: B256,
        transactions_root: B256,
    ) -> EthereumBlockHeader {
        let header = Header {
            number,
            parent_hash,
            transactions_root,
            timestamp: 1234567890,
            ..Default::default()
        };
        EthereumBlockHeader::from(header)
    }

    pub fn create_test_transaction_with_blob() -> TransactionWithBlob {
        let tx_data = TxEip1559 {
            to: TxKind::Call(address!("0x1234567890123456789012345678901234567890")),
            ..Default::default()
        };

        let signature = Signature::from_str(
            "0xaa231fbe0ed2b5418e6ba7c19bee2522852955ec50996c02a2fe3e71d30ddaf1645baf4823fea7cb4fcc7150842493847cfb6a6d63ab93e8ee928ee3f61f503500"
        ).expect("could not parse 0x-prefixed signature");
        let tx_signed = tx_data.into_signed(signature);
        let tx = EthereumTxEnvelope::Eip1559(tx_signed);
        TransactionWithBlob {
            tx,
            encoded_payload: Some(Bytes::from(b"test blob data".to_vec())),
            cert_state: None,
        }
    }

    #[test]
    fn test_completeness_proof_new() {
        let transaction_with_blob = create_test_transaction_with_blob();
        let transactions = vec![transaction_with_blob.tx.clone()];

        let proof = EigenDaCompletenessProof::new(transactions.clone());

        assert_eq!(proof.transactions, transactions);
    }

    #[test]
    fn test_completeness_proof_verify_empty() {
        let header = create_test_header(
            1,
            B256::default(),
            calculate_transaction_root::<TxEnvelope>(&[]),
        );
        let proof = EigenDaCompletenessProof::new(vec![]);

        let result = proof.verify(&header);
        assert!(result.is_ok());
        let transactions = result.unwrap();
        assert!(transactions.is_empty());
    }

    #[test]
    fn test_completeness_proof_verify_merkle_root_mismatch() {
        let header = create_test_header(1, B256::default(), B256::from_slice(&[1; 32]));
        let transaction = create_test_transaction_with_blob();
        let proof = EigenDaCompletenessProof::new(vec![transaction.tx]);

        let result = proof.verify(&header);
        assert!(matches!(
            result,
            Err(CompletenessProofError::MerkleRootMismatch(_, _))
        ));
    }

    #[test]
    fn test_completeness_proof_verify_with_valid_transactions() {
        let tx = create_test_transaction_with_blob();
        let tx_refs = vec![tx.tx.clone()];
        let tx_root = calculate_transaction_root(&tx_refs);

        let header = create_test_header(1, B256::default(), tx_root);
        let proof = EigenDaCompletenessProof::new(tx_refs.clone());

        let result = proof.verify(&header);
        assert!(result.is_ok());
        let verified_transactions = result.unwrap();
        assert_eq!(verified_transactions, tx_refs);
    }

    #[test]
    fn test_eigenda_verifier_new() {
        use crate::spec::RollupParams;

        let batch_namespace =
            NamespaceId::from_str("0x1234567890123456789012345678901234567890").unwrap();
        let proof_namespace =
            NamespaceId::from_str("0x9876543210987654321098765432109876543210").unwrap();
        let cert_recency_window = 100;

        let params = RollupParams {
            rollup_batch_namespace: batch_namespace,
            rollup_proof_namespace: proof_namespace,
            cert_recency_window,
        };

        let verifier = EigenDaVerifier::new(params);

        assert_eq!(verifier.batch_namespace, batch_namespace);
        assert_eq!(verifier.proof_namespace, proof_namespace);
        assert_eq!(verifier.cert_recency_window, cert_recency_window);
    }

    #[test]
    fn test_completeness_proof_error_display() {
        let hash1 = B256::from_slice(&[1; 32]);
        let hash2 = B256::from_slice(&[2; 32]);
        let merkle_error = CompletenessProofError::MerkleRootMismatch(hash1, hash2);
        let merkle_string = format!("{merkle_error}");
        assert!(merkle_string.contains("Recomputed merkle root"));
    }

    #[test]
    fn test_inclusion_proof_error_display() {
        let tx_hash = B256::from_slice(&[1; 32]);
        let error = InclusionProofError::NotProvenTransaction(tx_hash);
        let error_string = format!("{error}");
        assert!(error_string.contains("wasn't part of the completeness proof"));

        let incomplete_error = InclusionProofError::ProofIncomplete;
        let incomplete_string = format!("{incomplete_error}");
        assert!(incomplete_string.contains("some relevant transactions are missing"));
    }
}

#[cfg(feature = "use-rbn-state")]
#[cfg(test)]
mod use_rbn_state_tests {
    use std::str::FromStr;

    use alloy_primitives::{Address, B256};
    use bytes::Bytes;

    use crate::{
        spec::{BlobWithSender, NamespaceId},
        verifier::{EigenDaInclusionProof, InclusionProofError, tests},
    };

    #[test]
    fn test_inclusion_proof_new() {
        let transactions = vec![tests::create_test_transaction_with_blob()];
        let proof = EigenDaInclusionProof::new(vec![], transactions.clone());

        assert_eq!(proof.transactions, transactions);
    }

    #[test]
    fn test_inclusion_proof_verify_transactions_empty() {
        let proof = EigenDaInclusionProof::new(vec![], vec![]);
        let namespace =
            NamespaceId::from_str("0x1234567890123456789012345678901234567890").unwrap();

        let result = proof.verify_transactions(namespace, &[]);
        assert!(result.is_ok());
    }

    #[test]
    fn test_inclusion_proof_verify_transactions_not_proven() {
        let transaction = tests::create_test_transaction_with_blob();
        let proof = EigenDaInclusionProof::new(vec![], vec![transaction]);
        let namespace =
            NamespaceId::from_str("0x1234567890123456789012345678901234567890").unwrap();

        let result = proof.verify_transactions(namespace, &[]);
        assert!(matches!(
            result,
            Err(InclusionProofError::NotProvenTransaction(_))
        ));
    }

    #[test]
    fn test_inclusion_proof_verify_transactions_proof_incomplete() {
        let proven_tx = tests::create_test_transaction_with_blob();
        let proof = EigenDaInclusionProof::new(vec![], vec![]);
        let namespace =
            NamespaceId::from_str("0x1234567890123456789012345678901234567890").unwrap();

        let result = proof.verify_transactions(namespace, &[proven_tx.tx]);
        assert!(matches!(result, Err(InclusionProofError::ProofIncomplete)));
    }

    #[test]
    fn test_inclusion_proof_verify_empty_proof_and_transactions() {
        let header = tests::create_test_header(1, B256::default(), B256::default());
        let namespace =
            NamespaceId::from_str("0x1234567890123456789012345678901234567890").unwrap();
        let proof = EigenDaInclusionProof::new(vec![], vec![]);

        let result = proof.verify(&header, namespace, &[], &[], 100);
        assert!(result.is_ok());
    }

    #[test]
    fn test_inclusion_proof_verify_missing_blob_error() {
        let header = tests::create_test_header(1, B256::default(), B256::default());
        let namespace =
            NamespaceId::from_str("0x1234567890123456789012345678901234567890").unwrap();

        let tx_hash = B256::from_slice(&[1; 32]);
        let blob_sender = Address::from_str("0x1234567890123456789012345678901234567890").unwrap();
        let blob_with_sender = BlobWithSender::new(blob_sender, tx_hash, Bytes::new());

        let proof = EigenDaInclusionProof::new(vec![], vec![]);

        let result = proof.verify(&header, namespace, &[], &[blob_with_sender], 100);
        assert!(matches!(
            result,
            Err(InclusionProofError::IrrelevantBlob(_))
        ));
    }

    #[test]
    fn test_inclusion_proof_verify_incorrect_sender() {
        let header = tests::create_test_header(1, B256::default(), B256::default());
        let namespace =
            NamespaceId::from_str("0x1234567890123456789012345678901234567890").unwrap();

        let tx_hash = B256::from_slice(&[1; 32]);
        let wrong_sender = Address::from_str("0x9876543210987654321098765432109876543210").unwrap();
        let blob_with_sender = BlobWithSender::new(wrong_sender, tx_hash, Bytes::new());

        let proof = EigenDaInclusionProof::new(vec![], vec![]);

        let result = proof.verify(&header, namespace, &[], &[blob_with_sender], 100);
        assert!(matches!(
            result,
            Err(InclusionProofError::IrrelevantBlob(_))
        ));
    }

    #[test]
    fn test_inclusion_proof_verify_incorrect_hash() {
        let header = tests::create_test_header(1, B256::default(), B256::default());
        let namespace =
            NamespaceId::from_str("0x1234567890123456789012345678901234567890").unwrap();

        let wrong_hash = B256::from_slice(&[99; 32]);
        let sender = Address::from_str("0x1234567890123456789012345678901234567890").unwrap();
        let blob_with_sender = BlobWithSender::new(sender, wrong_hash, Bytes::new());

        let proof = EigenDaInclusionProof::new(vec![], vec![]);

        let result = proof.verify(&header, namespace, &[], &[blob_with_sender], 100);
        assert!(matches!(
            result,
            Err(InclusionProofError::IrrelevantBlob(_))
        ));
    }

    #[test]
    fn test_inclusion_proof_verify_data_length_mismatch() {
        let header = tests::create_test_header(1, B256::default(), B256::default());
        let namespace =
            NamespaceId::from_str("0x1234567890123456789012345678901234567890").unwrap();

        let tx_hash = B256::from_slice(&[1; 32]);
        let sender = Address::from_str("0x1234567890123456789012345678901234567890").unwrap();
        let long_data = Bytes::from(vec![1u8; 1000]);
        let blob_with_sender = BlobWithSender::new(sender, tx_hash, long_data);

        let proof = EigenDaInclusionProof::new(vec![], vec![]);

        let result = proof.verify(&header, namespace, &[], &[blob_with_sender], 100);
        assert!(matches!(
            result,
            Err(InclusionProofError::IrrelevantBlob(_))
        ));
    }
}
