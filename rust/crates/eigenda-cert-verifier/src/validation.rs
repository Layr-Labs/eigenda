use crate::{error::SignaturesVerificationError, types::NonSignerStakesAndSignature};

pub fn validate_inputs<'a>(
    quorum_numbers: &[u8],
    reference_block_number: u32,
    current_block_number: u32,
    params: &NonSignerStakesAndSignature,
) -> Result<(), SignaturesVerificationError<'a>> {
    use SignaturesVerificationError::*;

    if quorum_numbers.is_empty() {
        return Err(EmptyQuorumNumbers);
    }

    let quorum_numbers_len = quorum_numbers.len();

    let quorum_apks_len = params.quorum_apks.len();
    if quorum_numbers_len != quorum_apks_len {
        return Err(QuorumNumbersAndQuorumApksLengthMismatch {
            quorum_numbers_len,
            quorum_apks_len,
        });
    }

    let quorum_apk_indices_len = params.quorum_apk_indices.len();
    if quorum_numbers_len != quorum_apk_indices_len {
        return Err(QuorumNumbersAndQuorumApkIndicesLengthMismatch {
            quorum_numbers_len,
            quorum_apk_indices_len,
        });
    }

    let total_stake_indices_len = params.total_stake_indices.len();
    if quorum_numbers_len != total_stake_indices_len {
        return Err(QuorumNumbersAndTotalStakeIndicesLengthMismatch {
            quorum_numbers_len,
            total_stake_indices_len,
        });
    }

    let non_signer_stake_indices_len = params.non_signer_stake_indices.len();
    if quorum_numbers_len != non_signer_stake_indices_len {
        return Err(QuorumNumbersAndNonSignerStakeIndicesLengthMismatch {
            quorum_numbers_len,
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
        BlsSignaturesVerifier, NonSignerStakesAndSignature, SignatureVerifier,
        error::SignaturesVerificationError, types::ReferenceBlock,
    };

    #[test]
    fn test_verify_signatures_fails_given_empty_quorum_numbers() {
        let signatures_verifier = BlsSignaturesVerifier::default();

        let msg_hash = [0u8; 32];
        let quorum_numbers = vec![];
        let reference_block_number = 0u32;
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature::default();
        let reference_block = ReferenceBlock::default();

        let signatures_verification = signatures_verifier.verify_signatures(
            msg_hash,
            &quorum_numbers,
            reference_block_number,
            current_block_number,
            &params,
            &reference_block,
        );

        assert_eq!(
            signatures_verification.unwrap_err(),
            SignaturesVerificationError::EmptyQuorumNumbers
        );
    }

    #[test]
    fn test_verify_signatures_fails_due_to_quorum_numbers_and_quorum_apks_length_mismatch() {
        let signatures_verifier = BlsSignaturesVerifier::default();

        let msg_hash = [0u8; 32];
        let quorum_numbers = vec![0u8; 1];
        let reference_block_number = 0u32;
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 42],
            quorum_apk_indices: vec![0u8; 1],
            total_stake_indices: vec![0u8; 1],
            non_signer_stake_indices: vec![0u8; 1],
            non_signer_pubkeys: vec![G1Affine::default(); 1],
            non_signer_quorum_bitmap_indices: vec![0u8; 1],
        };
        let reference_block = ReferenceBlock::default();

        let signatures_verification = signatures_verifier.verify_signatures(
            msg_hash,
            &quorum_numbers,
            reference_block_number,
            current_block_number,
            &params,
            &reference_block,
        );

        assert_eq!(
            signatures_verification.unwrap_err(),
            SignaturesVerificationError::QuorumNumbersAndQuorumApksLengthMismatch {
                quorum_numbers_len: quorum_numbers.len(),
                quorum_apks_len: params.quorum_apks.len(),
            }
        );
    }

    #[test]
    fn test_verify_signatures_fails_due_to_quorum_numbers_and_quorum_apk_indices_length_mismatch() {
        let signatures_verifier = BlsSignaturesVerifier::default();

        let msg_hash = [0u8; 32];
        let quorum_numbers = vec![0u8; 1];
        let reference_block_number = 0u32;
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 1],
            quorum_apk_indices: vec![0u8; 42],
            total_stake_indices: vec![0u8; 1],
            non_signer_stake_indices: vec![0u8; 1],
            non_signer_pubkeys: vec![G1Affine::default(); 1],
            non_signer_quorum_bitmap_indices: vec![0u8; 1],
        };
        let reference_block = ReferenceBlock::default();

        let signatures_verification = signatures_verifier.verify_signatures(
            msg_hash,
            &quorum_numbers,
            reference_block_number,
            current_block_number,
            &params,
            &reference_block,
        );

        assert_eq!(
            signatures_verification.unwrap_err(),
            SignaturesVerificationError::QuorumNumbersAndQuorumApkIndicesLengthMismatch {
                quorum_numbers_len: quorum_numbers.len(),
                quorum_apk_indices_len: params.quorum_apk_indices.len(),
            }
        );
    }

    #[test]
    fn test_verify_signatures_fails_due_to_quorum_numbers_and_total_stake_indices_length_mismatch()
    {
        let signatures_verifier = BlsSignaturesVerifier::default();

        let msg_hash = [0u8; 32];
        let quorum_numbers = vec![0u8; 1];
        let reference_block_number = 0u32;
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 1],
            quorum_apk_indices: vec![0u8; 1],
            total_stake_indices: vec![0u8; 42],
            non_signer_stake_indices: vec![0u8; 1],
            non_signer_pubkeys: vec![G1Affine::default(); 1],
            non_signer_quorum_bitmap_indices: vec![0u8; 1],
        };
        let reference_block = ReferenceBlock::default();

        let signatures_verification = signatures_verifier.verify_signatures(
            msg_hash,
            &quorum_numbers,
            reference_block_number,
            current_block_number,
            &params,
            &reference_block,
        );

        assert_eq!(
            signatures_verification.unwrap_err(),
            SignaturesVerificationError::QuorumNumbersAndTotalStakeIndicesLengthMismatch {
                quorum_numbers_len: quorum_numbers.len(),
                total_stake_indices_len: params.total_stake_indices.len(),
            }
        );
    }

    #[test]
    fn test_verify_signatures_fails_due_to_quorum_numbers_and_non_signer_stake_indices_length_mismatch()
     {
        let signatures_verifier = BlsSignaturesVerifier::default();

        let msg_hash = [0u8; 32];
        let quorum_numbers = vec![0u8; 1];
        let reference_block_number = 0u32;
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 1],
            quorum_apk_indices: vec![0u8; 1],
            total_stake_indices: vec![0u8; 1],
            non_signer_stake_indices: vec![0u8; 42],
            non_signer_pubkeys: vec![G1Affine::default(); 1],
            non_signer_quorum_bitmap_indices: vec![0u8; 1],
        };
        let reference_block = ReferenceBlock::default();

        let signatures_verification = signatures_verifier.verify_signatures(
            msg_hash,
            &quorum_numbers,
            reference_block_number,
            current_block_number,
            &params,
            &reference_block,
        );

        assert_eq!(
            signatures_verification.unwrap_err(),
            SignaturesVerificationError::QuorumNumbersAndNonSignerStakeIndicesLengthMismatch {
                quorum_numbers_len: quorum_numbers.len(),
                non_signer_stake_indices_len: params.non_signer_stake_indices.len(),
            }
        );
    }

    #[test]
    fn test_verify_signatures_fails_due_to_non_signer_pubkeys_and_non_signer_quorum_bitmap_indices_length_mismatch()
     {
        let signatures_verifier = BlsSignaturesVerifier::default();

        let msg_hash = [0u8; 32];
        let quorum_numbers = vec![0u8; 1];
        let reference_block_number = 0u32;
        let current_block_number = 1u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 1],
            quorum_apk_indices: vec![0u8; 1],
            total_stake_indices: vec![0u8; 1],
            non_signer_stake_indices: vec![0u8; 1],
            non_signer_pubkeys: vec![G1Affine::default(); 42],
            non_signer_quorum_bitmap_indices: vec![0u8; 41],
        };
        let reference_block = ReferenceBlock::default();

        let signatures_verification = signatures_verifier.verify_signatures(
            msg_hash,
            &quorum_numbers,
            reference_block_number,
            current_block_number,
            &params,
            &reference_block,
        );

        assert_eq!(
            signatures_verification.unwrap_err(),
            SignaturesVerificationError::NonSignerPubkeysAndNonSignerQuorumBitmapIndicesLengthMismatch {
                non_signer_pubkeys_len: params.non_signer_pubkeys.len(),
                non_signer_quorum_bitmap_indices_len: params.non_signer_quorum_bitmap_indices.len(),
            }
        );
    }

    #[test]
    fn test_verify_signatures_fails_with_reference_block_not_preceding_current_block() {
        let signatures_verifier = BlsSignaturesVerifier::default();

        let msg_hash = [0u8; 32];
        let quorum_numbers = vec![0u8; 1];
        let reference_block_number = 42u32;
        let current_block_number = 41u32;
        let params = NonSignerStakesAndSignature {
            quorum_apks: vec![G1Affine::default(); 1],
            quorum_apk_indices: vec![0u8; 1],
            total_stake_indices: vec![0u8; 1],
            non_signer_stake_indices: vec![0u8; 1],
            non_signer_pubkeys: vec![G1Affine::default(); 1],
            non_signer_quorum_bitmap_indices: vec![0u8; 1],
        };
        let reference_block = ReferenceBlock::default();

        let signatures_verification = signatures_verifier.verify_signatures(
            msg_hash,
            &quorum_numbers,
            reference_block_number,
            current_block_number,
            &params,
            &reference_block,
        );

        assert_eq!(
            signatures_verification.unwrap_err(),
            SignaturesVerificationError::ReferenceBlockDoesNotPrecedeCurrentBlock {
                reference_block_number,
                current_block_number
            }
        );
    }
}
