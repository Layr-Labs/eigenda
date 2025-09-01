use std::fmt::Display;

use crate::eigenda::cert::{BatchHeaderV2, BlobCertificate, BlobHeaderV2};
use alloy_primitives::{B256, keccak256};
use alloy_sol_types::SolValue;
use smallvec::SmallVec;

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

pub fn keccak256_many<T: AsRef<[u8]>>(values: &[T]) -> B256 {
    let mut buffer: SmallVec<[u8; 256]> = SmallVec::with_capacity(values.len() * 32);

    for value in values {
        buffer.extend_from_slice(value.as_ref());
    }

    keccak256(&buffer)
}
