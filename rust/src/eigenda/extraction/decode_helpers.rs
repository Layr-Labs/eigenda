//! Utilities for decoding storage proofs from Ethereum contracts
//!
//! This module provides helper functions for working with storage proofs
//! returned from Ethereum nodes, making it easier to extract the data needed to
//! validate EigenDA certificates.

use alloy_primitives::StorageKey;
use reth_trie_common::StorageProof;

use crate::eigenda::{
    extraction::CertExtractionError,
    verification::cert::types::history::{HistoryError, Update},
};

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
    use CertExtractionError::*;

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
