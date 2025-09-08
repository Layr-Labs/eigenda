use std::fmt::Display;

use crate::eigenda::cert::{BatchHeaderV2, BlobCertificate, BlobHeaderV2};
use alloy_primitives::{B256, Keccak256, keccak256};
use alloy_sol_types::SolValue;
use derive_more::{AsMut, AsRef, Deref, DerefMut, From, Into};

#[repr(transparent)]
#[derive(
    Debug, Clone, Copy, PartialEq, Eq, Hash, Deref, DerefMut, AsRef, AsMut, From, Into, Default,
)]
pub struct TruncHash(pub [u8; 24]);

impl Display for TruncHash {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "{}", hex::encode(self.0))
    }
}

pub trait HashExt {
    fn hash_ext(&self) -> B256;
}

impl HashExt for BlobCertificate {
    fn hash_ext(&self) -> B256 {
        let blob_header = self.blob_header.hash_ext();
        let encoded =
            (blob_header, self.signature.clone(), self.relay_keys.clone()).abi_encode_sequence();
        keccak256(&encoded)
    }
}

impl HashExt for BlobHeaderV2 {
    fn hash_ext(&self) -> B256 {
        let encoded = (
            self.version,
            self.quorum_numbers.clone(),
            self.commitment.to_sol(),
        )
            .abi_encode_sequence();

        let hashed = keccak256(&encoded);
        let encoded = (hashed, self.payment_header_hash).abi_encode();
        keccak256(&encoded)
    }
}

impl HashExt for BatchHeaderV2 {
    fn hash_ext(&self) -> B256 {
        let encoded = self.to_sol().abi_encode();
        keccak256(&encoded)
    }
}

pub fn streaming_keccak256<T: AsRef<[u8]>>(values: &[T]) -> B256 {
    let mut hasher = Keccak256::new();
    for v in values {
        hasher.update(v.as_ref());
    }
    hasher.finalize()
}
