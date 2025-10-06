//! Utilities for decoding storage proofs from Ethereum contracts
//!
//! This module provides helper functions for working with storage proofs
//! returned from Ethereum nodes, making it easier to extract the data needed to
//! validate EigenDA certificates.

use alloy_primitives::StorageKey;
use eigenda_verification::verification::cert::types::history::{HistoryError, Update};
use reth_trie_common::StorageProof;

use crate::extraction::CertExtractionError;

/// Find a storage proof by key and return error if missing
///
/// This is the primary function used by extractors to locate the storage proof
/// they need for decoding contract state.
///
/// # Arguments
/// * `proofs` - Array of storage proofs from the Ethereum node
/// * `key` - The storage key being sought
/// * `variable_name` - Name of the contract variable for error reporting
///
/// # Returns
/// Reference to the storage proof if found
///
/// # Errors
/// Returns [`CertExtractionError::MissingStorageProof`] if the key is not found
pub fn find_required_proof<'a, T>(
    proofs: &'a [StorageProof],
    key: &StorageKey,
    variable_name: T,
) -> Result<&'a StorageProof, CertExtractionError>
where
    T: std::fmt::Display,
{
    use crate::extraction::CertExtractionError::*;

    find_proof(proofs, key).ok_or_else(|| MissingStorageProof(variable_name.to_string()))
}

/// Find a storage proof by key
///
/// Low-level function that searches for a storage proof without error handling.
/// Only the fallible [`find_required_proof`] is publicly exposed.
///
/// # Arguments  
/// * `proofs` - Array of storage proofs from the Ethereum node
/// * `key` - The storage key being sought
///
/// # Returns
/// `Some` reference to the storage proof if found
/// `None` if not found
fn find_proof<'a>(proofs: &'a [StorageProof], key: &StorageKey) -> Option<&'a StorageProof> {
    proofs.iter().find(|proof| proof.key == *key)
}

/// Create an Update object from extracted block numbers and value
///
/// Helper function for constructing history update entries from contract storage.
/// Handles the validation of block number relationships.
///
/// # Arguments
/// * `update_block` - Block number when this value was updated
/// * `next_update_block` - Block number when this value will be/was superseded
/// * `value` - The actual value being tracked in history
///
/// # Returns
/// Update object for use in history tracking
///
/// # Errors
/// Returns [`HistoryError`] if block number relationships are invalid
pub fn create_update<T: Copy + std::fmt::Debug>(
    update_block: u32,
    next_update_block: u32,
    value: T,
) -> Result<Update<T>, HistoryError> {
    Update::new(update_block, next_update_block, value)
}

#[cfg(test)]
mod tests {
    use alloy_primitives::{B256, StorageKey, U256};
    use eigenda_verification::verification::cert::types::history::HistoryError;
    use reth_trie_common::StorageProof;

    use super::{create_update, find_required_proof};
    use crate::extraction::CertExtractionError;

    fn create_test_storage_proof(key: StorageKey, value: U256) -> StorageProof {
        StorageProof {
            key,
            value,
            ..Default::default()
        }
    }

    fn create_test_key(value: u8) -> StorageKey {
        StorageKey::from(B256::repeat_byte(value))
    }

    #[test]
    fn find_required_proof_success() {
        let key1 = create_test_key(1);
        let key2 = create_test_key(2);
        let key3 = create_test_key(3);

        let proof1 = create_test_storage_proof(key1, U256::from(100));
        let proof2 = create_test_storage_proof(key2, U256::from(200));
        let proof3 = create_test_storage_proof(key3, U256::from(300));

        let proofs = vec![proof1, proof2, proof3];

        let found_proof = find_required_proof(&proofs, &key2, "test_variable").unwrap();

        assert_eq!(found_proof.key, key2);
        assert_eq!(found_proof.value, U256::from(200));
    }

    #[test]
    fn find_required_proof_missing_key() {
        let key1 = create_test_key(1);
        let key2 = create_test_key(2);
        let missing_key = create_test_key(99);

        let proof1 = create_test_storage_proof(key1, U256::from(100));
        let proof2 = create_test_storage_proof(key2, U256::from(200));

        let proofs = vec![proof1, proof2];

        let err = find_required_proof(&proofs, &missing_key, "missing_variable").unwrap_err();
        assert!(
            matches!(err, CertExtractionError::MissingStorageProof(ref var_name) if var_name == "missing_variable")
        );
    }

    #[test]
    fn find_required_proof_empty_proofs() {
        let key = create_test_key(1);
        let proofs: Vec<StorageProof> = vec![];

        let err = find_required_proof(&proofs, &key, "empty_proofs").unwrap_err();
        assert!(
            matches!(err, CertExtractionError::MissingStorageProof(ref var_name) if var_name == "empty_proofs")
        );
    }

    #[test]
    fn create_update_success() {
        let update = create_update(100, 200, "test_value").unwrap();
        assert_eq!(update.update_block_number(), 100);
        assert_eq!(update.next_update_block_number(), 200);
        assert_eq!(*update.value(), "test_value");
    }

    #[test]
    fn create_update_same_block_numbers() {
        let err = create_update(100, 100, 42u32).unwrap_err();
        assert!(matches!(err, HistoryError::InvalidBlockOrder { .. }));
    }

    #[test]
    fn create_update_invalid_order() {
        let err = create_update(200, 100, 42u32).unwrap_err();
        assert!(matches!(err, HistoryError::InvalidBlockOrder { .. }));
    }
}
