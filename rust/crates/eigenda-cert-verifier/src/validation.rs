use crate::{error::SignaturesVerificationError, types::NonSignerStakesAndSignature};

pub fn validate_inputs<'a>(
    signed_quorum_numbers: &[u8],
    reference_block_number: u32,
    current_block_number: u32,
    params: &NonSignerStakesAndSignature,
) -> Result<(), SignaturesVerificationError<'a>> {
    use SignaturesVerificationError::*;

    if signed_quorum_numbers.is_empty() {
        return Err(EmptyQuorumNumbers);
    }

    let signed_quorum_numbers_len = signed_quorum_numbers.len();

    let quorum_apks_len = params.quorum_apks.len();
    if signed_quorum_numbers_len != quorum_apks_len {
        return Err(QuorumNumbersAndQuorumApksLengthMismatch {
            signed_quorum_numbers_len,
            quorum_apks_len,
        });
    }

    let quorum_apk_indices_len = params.quorum_apk_indices.len();
    if signed_quorum_numbers_len != quorum_apk_indices_len {
        return Err(QuorumNumbersAndQuorumApkIndicesLengthMismatch {
            signed_quorum_numbers_len,
            quorum_apk_indices_len,
        });
    }

    let total_stake_indices_len = params.total_stake_indices.len();
    if signed_quorum_numbers_len != total_stake_indices_len {
        return Err(QuorumNumbersAndTotalStakeIndicesLengthMismatch {
            signed_quorum_numbers_len,
            total_stake_indices_len,
        });
    }

    let non_signer_stake_indices_len = params.non_signer_stake_indices.len();
    if signed_quorum_numbers_len != non_signer_stake_indices_len {
        return Err(QuorumNumbersAndNonSignerStakeIndicesLengthMismatch {
            signed_quorum_numbers_len,
            non_signer_stake_indices_len,
        });
    }

    if reference_block_number >= current_block_number {
        return Err(ReferenceBlockDoesNotPrecedeCurrentBlock {
            reference_block_number,
            current_block_number,
        });
    }

    let non_signer_pubkeys_len = params.non_signer_pubkeys.len();
    let non_signer_quorum_bitmap_indices_len = params.non_signer_quorum_bitmap_indices.len();
    if non_signer_pubkeys_len != non_signer_quorum_bitmap_indices_len {
        return Err(
            NonSignerPubkeysAndNonSignerQuorumBitmapIndicesLengthMismatch {
                non_signer_pubkeys_len,
                non_signer_quorum_bitmap_indices_len,
            },
        );
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use alloc::vec;

    use ark_bn254::G1Affine;

    use crate::{
        NonSignerStakesAndSignature, error::SignaturesVerificationError, types::ReferenceBlock,
        validation::validate_inputs,
    };

    #[test]
    fn validate_inputs_succeeds_given_valid_inputs() {
        let signed_quorum_numbers = vec![0; 2];
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 2],
            quorum_apk_indices: vec![0; 2],
            total_stake_indices: vec![0; 2],
            non_signer_stake_indices: vec![0; 2],
            ..Default::default()
        };
        let reference_block = ReferenceBlock::default();
        let current_block_number = reference_block.number + 1;

        let result = validate_inputs(
            &signed_quorum_numbers,
            reference_block.number,
            current_block_number,
            &params,
        );

        assert!(result.is_ok());
    }

    #[test]
    fn verify_signatures_fails_given_empty_signed_quorum_numbers() {
        let signed_quorum_numbers = vec![];
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature::default();
        let reference_block = ReferenceBlock::default();

        let result = validate_inputs(
            &signed_quorum_numbers,
            reference_block.number,
            current_block_number,
            &params,
        );

        assert_eq!(
            result.unwrap_err(),
            SignaturesVerificationError::EmptyQuorumNumbers
        );
    }

    #[test]
    fn verify_signatures_fails_due_to_signed_quorum_numbers_and_quorum_apks_length_mismatch() {
        let signed_quorum_numbers = vec![0u8; 1];
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 42],
            quorum_apk_indices: vec![0u8; 1],
            total_stake_indices: vec![0u8; 1],
            non_signer_stake_indices: vec![0u8; 1],
            non_signer_pubkeys: vec![G1Affine::default(); 1],
            non_signer_quorum_bitmap_indices: vec![0u8; 1],
            ..Default::default()
        };
        let reference_block = ReferenceBlock::default();

        let result = validate_inputs(
            &signed_quorum_numbers,
            reference_block.number,
            current_block_number,
            &params,
        );

        assert_eq!(
            result.unwrap_err(),
            SignaturesVerificationError::QuorumNumbersAndQuorumApksLengthMismatch {
                signed_quorum_numbers_len: signed_quorum_numbers.len(),
                quorum_apks_len: params.quorum_apks.len(),
            }
        );
    }

    #[test]
    fn verify_signatures_fails_due_to_signed_quorum_numbers_and_quorum_apk_indices_length_mismatch()
    {
        let signed_quorum_numbers = vec![0u8; 1];
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 1],
            quorum_apk_indices: vec![0u8; 42],
            total_stake_indices: vec![0u8; 1],
            non_signer_stake_indices: vec![0u8; 1],
            non_signer_pubkeys: vec![G1Affine::default(); 1],
            non_signer_quorum_bitmap_indices: vec![0u8; 1],
            ..Default::default()
        };
        let reference_block = ReferenceBlock::default();

        let result = validate_inputs(
            &signed_quorum_numbers,
            reference_block.number,
            current_block_number,
            &params,
        );

        assert_eq!(
            result.unwrap_err(),
            SignaturesVerificationError::QuorumNumbersAndQuorumApkIndicesLengthMismatch {
                signed_quorum_numbers_len: signed_quorum_numbers.len(),
                quorum_apk_indices_len: params.quorum_apk_indices.len(),
            }
        );
    }

    #[test]
    fn verify_signatures_fails_due_to_signed_quorum_numbers_and_total_stake_indices_length_mismatch()
     {
        let signed_quorum_numbers = vec![0u8; 1];
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 1],
            quorum_apk_indices: vec![0u8; 1],
            total_stake_indices: vec![0u8; 42],
            non_signer_stake_indices: vec![0u8; 1],
            non_signer_pubkeys: vec![G1Affine::default(); 1],
            non_signer_quorum_bitmap_indices: vec![0u8; 1],
            ..Default::default()
        };
        let reference_block = ReferenceBlock::default();

        let result = validate_inputs(
            &signed_quorum_numbers,
            reference_block.number,
            current_block_number,
            &params,
        );

        assert_eq!(
            result.unwrap_err(),
            SignaturesVerificationError::QuorumNumbersAndTotalStakeIndicesLengthMismatch {
                signed_quorum_numbers_len: signed_quorum_numbers.len(),
                total_stake_indices_len: params.total_stake_indices.len(),
            }
        );
    }

    #[test]
    fn verify_signatures_fails_due_to_signed_quorum_numbers_and_non_signer_stake_indices_length_mismatch()
     {
        let signed_quorum_numbers = vec![0u8; 1];
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 1],
            quorum_apk_indices: vec![0u8; 1],
            total_stake_indices: vec![0u8; 1],
            non_signer_stake_indices: vec![0u8; 42],
            non_signer_pubkeys: vec![G1Affine::default(); 1],
            non_signer_quorum_bitmap_indices: vec![0u8; 1],
            ..Default::default()
        };
        let reference_block = ReferenceBlock::default();

        let result = validate_inputs(
            &signed_quorum_numbers,
            reference_block.number,
            current_block_number,
            &params,
        );

        assert_eq!(
            result.unwrap_err(),
            SignaturesVerificationError::QuorumNumbersAndNonSignerStakeIndicesLengthMismatch {
                signed_quorum_numbers_len: signed_quorum_numbers.len(),
                non_signer_stake_indices_len: params.non_signer_stake_indices.len(),
            }
        );
    }

    #[test]
    fn verify_signatures_fails_due_to_non_signer_pubkeys_and_non_signer_quorum_bitmap_indices_length_mismatch()
     {
        let signed_quorum_numbers = vec![0u8; 1];
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 1],
            quorum_apk_indices: vec![0u8; 1],
            total_stake_indices: vec![0u8; 1],
            non_signer_stake_indices: vec![0u8; 1],
            non_signer_pubkeys: vec![G1Affine::default(); 42],
            non_signer_quorum_bitmap_indices: vec![0u8; 41],
            ..Default::default()
        };
        let reference_block = ReferenceBlock::default();

        let result = validate_inputs(
            &signed_quorum_numbers,
            reference_block.number,
            current_block_number,
            &params,
        );

        assert_eq!(
            result.unwrap_err(),
            SignaturesVerificationError::NonSignerPubkeysAndNonSignerQuorumBitmapIndicesLengthMismatch {
                non_signer_pubkeys_len: params.non_signer_pubkeys.len(),
                non_signer_quorum_bitmap_indices_len: params.non_signer_quorum_bitmap_indices.len(),
            }
        );
    }

    #[test]
    fn verify_signatures_fails_with_reference_block_not_preceding_current_block() {
        let signed_quorum_numbers = vec![0u8; 1];
        let current_block_number = 41u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 1],
            quorum_apk_indices: vec![0u8; 1],
            total_stake_indices: vec![0u8; 1],
            non_signer_stake_indices: vec![0u8; 1],
            non_signer_pubkeys: vec![G1Affine::default(); 1],
            non_signer_quorum_bitmap_indices: vec![0u8; 1],
            ..Default::default()
        };
        let reference_block = ReferenceBlock {
            number: 42,
            ..Default::default()
        };

        let result = validate_inputs(
            &signed_quorum_numbers,
            reference_block.number,
            current_block_number,
            &params,
        );

        assert_eq!(
            result.unwrap_err(),
            SignaturesVerificationError::ReferenceBlockDoesNotPrecedeCurrentBlock {
                reference_block_number: reference_block.number,
                current_block_number
            }
        );
    }
}
