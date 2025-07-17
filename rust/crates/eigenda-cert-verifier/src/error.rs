use thiserror::Error;

#[derive(Error, Debug, PartialEq)]
pub enum SignaturesVerificationError<'a> {
    #[error("Empty quorum numbers")]
    EmptyQuorumNumbers,

    #[error(
        "Quorum numbers length ({signed_quorum_numbers_len}) does not match quorum apks length ({quorum_apks_len})"
    )]
    QuorumNumbersAndQuorumApksLengthMismatch {
        signed_quorum_numbers_len: usize,
        quorum_apks_len: usize,
    },

    #[error(
        "Quorum numbers length ({signed_quorum_numbers_len}) does not match quorum apk indices length ({quorum_apk_indices_len})"
    )]
    QuorumNumbersAndQuorumApkIndicesLengthMismatch {
        signed_quorum_numbers_len: usize,
        quorum_apk_indices_len: usize,
    },

    #[error(
        "Quorum numbers length ({signed_quorum_numbers_len}) does not match total stake indices length ({total_stake_indices_len})"
    )]
    QuorumNumbersAndTotalStakeIndicesLengthMismatch {
        signed_quorum_numbers_len: usize,
        total_stake_indices_len: usize,
    },

    #[error(
        "Quorum numbers length ({signed_quorum_numbers_len}) does not match non signer stake indices length ({non_signer_stake_indices_len})"
    )]
    QuorumNumbersAndNonSignerStakeIndicesLengthMismatch {
        signed_quorum_numbers_len: usize,
        non_signer_stake_indices_len: usize,
    },

    #[error(
        "Non-signer pubkeys length ({non_signer_pubkeys_len}) does not match non-signer quorum bitmap indices length ({non_signer_quorum_bitmap_indices_len})"
    )]
    NonSignerPubkeysAndNonSignerQuorumBitmapIndicesLengthMismatch {
        non_signer_pubkeys_len: usize,
        non_signer_quorum_bitmap_indices_len: usize,
    },

    #[error(
        "Reference block ({reference_block_number}) must precede current block ({current_block_number})"
    )]
    ReferenceBlockDoesNotPrecedeCurrentBlock {
        reference_block_number: u32,
        current_block_number: u32,
    },

    #[error(
        "Bit indices length ({bit_indices_len}) exceeds max byte slice length ({max_bit_indices_len})"
    )]
    BitIndicesGreaterThanMaxLength {
        bit_indices_len: usize,
        max_bit_indices_len: usize,
    },

    #[error("Bit indices are not unique: {bit_indices:?}")]
    BitIndicesNotUnique { bit_indices: &'a [u8] },

    #[error("Bit indices are not ordered: {bit_indices:?}")]
    BitIndicesNotSorted { bit_indices: &'a [u8] },

    #[error(
        "One or more bit index out of {bit_indices:?} is greater than or equal to the provided upper bound: {upper_bound_bit_index:?}"
    )]
    BitIndexGreaterThanOrEqualToUpperBound {
        bit_indices: &'a [u8],
        upper_bound_bit_index: u8,
    },

    #[error("Expected to find a matching bitmap from the hash of a signer's pubkey")]
    SignerBitmapNotFoundOrPubkeyAtInfinity,
}
