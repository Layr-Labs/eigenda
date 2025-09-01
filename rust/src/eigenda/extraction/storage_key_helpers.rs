use alloy_primitives::{StorageKey, U256, keccak256};
use alloy_sol_types::SolValue;

/// Generate a simple storage key from a slot number: keccak256(abi.encode(slot))
pub fn simple_slot_key(slot: u64) -> StorageKey {
    U256::from(slot).into()
}

/// Generate storage key for a mapping: keccak256(abi.encode(key, slot))
pub fn mapping_key(key: U256, slot: u64) -> StorageKey {
    let slot = U256::from(slot);
    keccak256((key, slot).abi_encode())
}

/// Generate storage key for dynamic array element: keccak256(keccak256(abi.encode(key, slot))) + index
pub fn dynamic_array_key(key: U256, slot: u64, index: u32) -> StorageKey {
    let slot = U256::from(slot);
    let length_base = keccak256((key, slot).abi_encode());
    let data_base: U256 = keccak256(length_base).into();
    (data_base + U256::from(index)).into()
}

/// Generate nested mapping key for dynamic array:
/// keccak256(abi.encode(second_key,    keccak256(abi.encode(first_key, slot))    )) + index
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
