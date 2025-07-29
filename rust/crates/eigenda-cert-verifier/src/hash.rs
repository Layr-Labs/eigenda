use alloy_primitives::FixedBytes;
use alloy_sol_types::SolValue;
use core::iter::once;
use tiny_keccak::{Hasher, Keccak};

use crate::types::solidity::{BatchHeaderV2, BlobCertificate, BlobHeaderV2};

pub type TruncatedKeccak256Hash = [u8; 24];
pub type Keccak256Hash = FixedBytes<32>;

impl BlobCertificate {
    pub fn hash_keccak_v256(&self) -> Keccak256Hash {
        let blob_header = self.blobHeader.hash_keccak_v256();
        let encoded = (blob_header, self.signature.clone(), self.relayKeys.clone()).abi_encode();
        keccak_v256(once(encoded))
    }
}

impl BlobHeaderV2 {
    pub fn hash_keccak_v256(&self) -> Keccak256Hash {
        let encoded = (
            self.version,
            self.quorumNumbers.clone(),
            self.commitment.clone(),
        )
            .abi_encode();

        let hashed = keccak_v256(once(encoded));
        let encoded = (hashed, self.paymentHeaderHash).abi_encode();
        keccak_v256(once(encoded))
    }
}

impl BatchHeaderV2 {
    pub fn hash_keccak_v256(&self) -> Keccak256Hash {
        let batch_header = self.abi_encode();
        keccak_v256(once(batch_header))
    }
}

pub fn keccak_v256(inputs: impl Iterator<Item = impl AsRef<[u8]>>) -> Keccak256Hash {
    let mut hasher = Keccak::v256();
    for input in inputs {
        hasher.update(input.as_ref());
    }
    let mut output = [0u8; 32];
    hasher.finalize(&mut output);
    output.into()
}

pub fn signature_record(reference_block: u32, pk_hashes: &[Keccak256Hash]) -> Keccak256Hash {
    let first = reference_block.to_be_bytes();
    let inputs = once(first.as_ref()).chain(pk_hashes.iter().map(|pk_hash| pk_hash.as_ref()));
    keccak_v256(inputs)
}

#[cfg(test)]
mod tests {
    use alloc::vec::Vec;

    use crate::hash;

    #[test]
    fn signature_record_hash_success() {
        let reference_block = 42;
        let pk_hashes = [[42u8; 32], [43u8; 32]]
            .into_iter()
            .map(|pk_hash| pk_hash.into())
            .collect::<Vec<_>>();

        let actual = hash::signature_record(reference_block, &pk_hashes);
        let expected = [
            98, 139, 21, 105, 137, 7, 68, 235, 47, 165, 71, 215, 6, 47, 69, 231, 217, 9, 25, 96,
            61, 240, 244, 80, 244, 59, 71, 232, 252, 217, 178, 41,
        ];
        assert_eq!(actual, expected);
    }
}
