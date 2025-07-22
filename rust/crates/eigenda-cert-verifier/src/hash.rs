use alloc::vec::Vec;
use core::iter::once;
use tiny_keccak::{Hasher, Keccak};

use crate::types::{BlockNumber, NonSigner};

pub type BeHash = [u8; 32];
pub type TruncatedBeHash = [u8; 24];

pub fn keccak256(inputs: &[&[u8]]) -> BeHash {
    let mut hasher = Keccak::v256();
    for input in inputs {
        hasher.update(input);
    }
    let mut output = [0u8; 32];
    hasher.finalize(&mut output);
    output
}

pub fn signature_record(reference_block: BlockNumber, non_signers: &[NonSigner]) -> BeHash {
    let first = reference_block.to_be_bytes();

    let inputs = once(first.as_slice())
        .chain(
            non_signers
                .iter()
                .map(|non_signer| non_signer.pk_hash.as_slice()),
        )
        .collect::<Vec<_>>();

    keccak256(&inputs)
}

#[cfg(test)]
mod tests {
    use alloc::vec::Vec;

    use crate::{hash, types::NonSigner};

    #[test]
    fn signature_record_hash_success() {
        let reference_block = 42;
        let non_signers = [[42u8; 32], [43u8; 32]]
            .into_iter()
            .map(|pk_hash| NonSigner {
                pk_hash,
                ..Default::default()
            })
            .collect::<Vec<_>>();

        let actual = hash::signature_record(reference_block, &non_signers);
        let expected = [
            98, 139, 21, 105, 137, 7, 68, 235, 47, 165, 71, 215, 6, 47, 69, 231, 217, 9, 25, 96,
            61, 240, 244, 80, 244, 59, 71, 232, 252, 217, 178, 41,
        ];
        assert_eq!(actual, expected);
    }
}
