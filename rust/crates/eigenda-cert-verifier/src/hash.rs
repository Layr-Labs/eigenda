use alloy_primitives::B256;
use alloy_sol_types::SolValue;
use core::iter::once;
use eigenda_cert::{BatchHeaderV2, BlobCertificate, BlobHeaderV2};
use tiny_keccak::{Hasher, Keccak};

pub type TruncatedB256 = [u8; 24];

pub trait HashExt {
    fn hash_ext(&self) -> B256;
}

impl HashExt for BlobCertificate {
    fn hash_ext(&self) -> B256 {
        let blob_header = self.blob_header.hash_ext();
        let encoded = (blob_header, self.signature.clone(), self.relay_keys.clone()).abi_encode();
        keccak_v256(once(encoded))
    }
}

impl HashExt for BlobHeaderV2 {
    fn hash_ext(&self) -> B256 {
        let encoded = (
            self.version,
            self.quorum_numbers.clone(),
            self.commitment.to_sol(),
        )
            .abi_encode();

        let hashed = keccak_v256(once(encoded));
        let encoded = (hashed, self.payment_header_hash).abi_encode();
        keccak_v256(once(encoded))
    }
}

impl HashExt for BatchHeaderV2 {
    fn hash_ext(&self) -> B256 {
        let batch_header = self.to_sol().abi_encode();
        keccak_v256(once(batch_header))
    }
}

pub fn keccak_v256(inputs: impl Iterator<Item = impl AsRef<[u8]>>) -> B256 {
    let mut hasher = Keccak::v256();
    for input in inputs {
        hasher.update(input.as_ref());
    }
    let mut output = [0u8; 32];
    hasher.finalize(&mut output);
    output.into()
}

pub fn signature_record(reference_block: u32, pk_hashes: &[B256]) -> B256 {
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
