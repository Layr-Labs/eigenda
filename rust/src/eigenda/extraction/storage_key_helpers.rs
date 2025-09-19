//! Ethereum storage key generation utilities
//!
//! This module provides functions for generating storage keys used to access
//! Ethereum contract storage slots. It implements the standard Ethereum storage
//! layout rules for different data types.
//!
//! ## Storage Layout Rules
//!
//! Ethereum uses a specific storage layout for different data structures:
//! - Simple variables: stored directly at their slot number
//! - Mappings: `keccak256(abi.encode(key, slot))`
//! - Dynamic arrays: `keccak256(slot)` for base, then sequential slots
//! - Nested mappings: Multiple levels of keccak256 hashing
//!
//! ## References
//! - [Solidity Storage Layout](https://docs.soliditylang.org/en/latest/internals/layout_in_storage.html)

use alloy_primitives::{StorageKey, U256, keccak256};
use alloy_sol_types::SolValue;

/// Generate a simple storage key from a slot number
///
/// Used for basic state variables that occupy a single storage slot.
/// The slot number directly corresponds to the storage location.
///
/// # Arguments
/// * `slot` - The storage slot number
///
/// # Returns
/// Storage key for the slot
pub fn simple_slot_key(slot: u64) -> StorageKey {
    U256::from(slot).into()
}

/// Generate storage key for a mapping value
///
/// Implements the Ethereum mapping storage rule:
/// `storage_key = keccak256(abi.encode(key, slot))`
///
/// # Arguments
/// * `key` - The mapping key to look up
/// * `slot` - The storage slot of the mapping variable
///
/// # Returns
/// Storage key for the mapping value
pub fn mapping_key(key: U256, slot: u64) -> StorageKey {
    let slot = U256::from(slot);
    keccak256((key, slot).abi_encode())
}

/// Generate storage key for dynamic array element
///
/// Implements the Ethereum dynamic array storage rule:
/// `storage_key = keccak256(keccak256(abi.encode(key, slot))) + index`
///
/// The first keccak256 gives the array length location, the second gives
/// the data start location, then we add the index.
///
/// # Arguments
/// * `key` - The mapping key that contains the array
/// * `slot` - The storage slot of the mapping variable
/// * `index` - The array index to access
///
/// # Returns
/// Storage key for the array element
pub fn dynamic_array_key(key: U256, slot: u64, index: u32) -> StorageKey {
    let slot = U256::from(slot);
    let length_base = keccak256((key, slot).abi_encode());
    let data_base: U256 = keccak256(length_base).into();
    (data_base + U256::from(index)).into()
}

/// Generate storage key for nested mapping with dynamic array
///
/// Implements the storage rule for nested mappings containing arrays:
/// `storage_key = keccak256(keccak256(abi.encode(second_key, keccak256(abi.encode(first_key, slot))))) + index`
///
/// This handles structures like `mapping(address => mapping(uint256 => SomeStruct[]))`
///
/// # Arguments
/// * `first_key` - The first-level mapping key
/// * `slot` - The storage slot of the outer mapping variable  
/// * `second_key` - The second-level mapping key
/// * `index` - The array index to access
///
/// # Returns
/// Storage key for the nested array element
pub fn nested_dynamic_array_key(
    first_key: U256,
    slot: u64,
    second_key: U256,
    index: u32,
) -> StorageKey {
    let slot = U256::from(slot);
    let b1 = keccak256((first_key, slot).abi_encode());
    let b2 = keccak256((second_key, b1).abi_encode());
    let data_base: U256 = keccak256(b2).into();
    (data_base + U256::from(index)).into()
}

#[cfg(test)]
mod tests {
    use super::*;
    use alloy_primitives::hex;

    #[test]
    fn simple_slot_key_test() {
        let result = simple_slot_key(150);
        let value = hex!("0000000000000000000000000000000000000000000000000000000000000096");
        let expected = StorageKey::from(value);
        assert_eq!(result, expected);
    }

    #[test]
    fn mapping_key_test() {
        let result = mapping_key(U256::from(42), 5);
        let value = hex!("d3e7a847b0e4be9f2ff1f88564b0a771bb9789c2c82f98679296a6042483791d");
        let expected = StorageKey::from(value);
        assert_eq!(result, expected);
    }

    #[test]
    fn dynamic_array_key_test() {
        let result = dynamic_array_key(U256::from(0x123), 10, 5);
        let value = hex!("7fe76a52931b48d767fa7e54a1d7007662ab2827fd4b83ca6b158f06dbdbed88");
        let expected = StorageKey::from(value);
        assert_eq!(result, expected);
    }

    #[test]
    fn nested_dynamic_array_key_test() {
        let result = nested_dynamic_array_key(U256::from(0x456), 15, U256::from(0x789), 3);
        let value = hex!("7b559e449c242de80687a166a5b9feebff23ad66e81b26e687aa932f8ef0afca");
        let expected = StorageKey::from(value);
        assert_eq!(result, expected);
    }
}
