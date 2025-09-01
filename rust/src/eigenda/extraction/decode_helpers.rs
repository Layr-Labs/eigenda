use alloy_primitives::StorageKey;
use reth_trie_common::StorageProof;

use crate::eigenda::{
    extraction::CertExtractionError,
    verification::cert::types::history::{HistoryError, Update},
};

/// Find a storage proof by key and return error if missing
pub fn find_required_proof<'a, T>(
    proofs: &'a [StorageProof],
    key: &StorageKey,
    variable_name: T,
) -> Result<&'a StorageProof, CertExtractionError>
where
    T: std::fmt::Display,
{
    use CertExtractionError::*;

    find_proof(proofs, key).ok_or_else(|| MissingStorageProof(variable_name.to_string()))
}

/// Find a storage proof by key
pub fn find_proof<'a>(proofs: &'a [StorageProof], key: &StorageKey) -> Option<&'a StorageProof> {
    proofs.iter().find(|proof| proof.key == *key)
}

/// Create an Update object from extracted block numbers and value
pub fn create_update<T: Copy + std::fmt::Debug>(
    update_block: u32,
    next_update_block: u32,
    value: T,
) -> Result<Update<T>, HistoryError> {
    Update::new(update_block, next_update_block, value)
}
