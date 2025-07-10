//! Following eigenda-contracts/lib/eigenlayer-middleware/src/BLSSignatureChecker.sol
//! from 6797f3821db92c2214aaa6f137d94c603011ac2a lib/eigenlayer-middleware (v0.5.4-mainnet-rewards-v2-1-g6797f38)

// #![no_std]
extern crate alloc;

mod bitmap_utils;
mod error;
mod hashing;
mod types;
mod validation;

use ark_bn254::G1Affine;
use error::SignaturesVerificationError;
use types::{NonSignerInfo, NonSignerStakesAndSignature, QuorumStakeTotals};

#[derive(Default, Debug)]
pub struct SignaturesVerification {
    pub quorum_stake_totals: QuorumStakeTotals,
    pub signatory_record_hash: [u8; 32],
}

pub trait SignatureVerifier {
    fn verify_signatures(
        &self,
        msg_hash: [u8; 32],
        quorum_numbers: &[u8],
        reference_block_number: u32,
        current_block_number: u32,
        params: &NonSignerStakesAndSignature,
    ) -> Result<SignaturesVerification, SignaturesVerificationError<'_>>;
}

#[derive(Default, Debug)]
pub struct BlsSignaturesVerifier;

impl SignatureVerifier for BlsSignaturesVerifier {
    fn verify_signatures(
        &self,
        _msg_hash: [u8; 32],
        quorum_numbers: &[u8],
        reference_block_number: u32,
        current_block_number: u32,
        params: &NonSignerStakesAndSignature,
    ) -> Result<SignaturesVerification, SignaturesVerificationError<'_>> {
        validation::validate_inputs(
            quorum_numbers,
            reference_block_number,
            current_block_number,
            params,
        )?;

        // pseudo-code to think this through:
        // let overall_aggregate_non_signer_pubkey =  = quorums.iter().map(|quorum|
        //      let aggregate_non_signer_pubkey_of_quorum = quorum.pubkeys.iter()
        //          .filter(|pubkey| !quorum.pubkey2signature.contains(pubkey)).sum();
        //      aggregate_non_signer_pubkey_of_quorum
        // .sum();
        // is the subtraction done per quorum?
        let _apk = G1Affine::default();

        let _stake_totals = QuorumStakeTotals::default();
        let _non_signer_info = NonSignerInfo::default();

        // todo: this is queried from on-chain state at the current block for all I can tell
        // why is this not queried at the reference block?
        // it'd be ok if quorum_len could never shrink
        let quorum_len = 10; // registryCoordinator.quorumCount()
        let _signing_quorum_bitmap =
            bitmap_utils::bit_indices_to_bitmap(quorum_numbers, quorum_len);

        for non_signer_pubkey in &params.non_signer_pubkeys {
            let _hash = hashing::hash_g1_point(non_signer_pubkey);
        }

        Ok(SignaturesVerification::default())
    }
}

#[cfg(test)]
mod tests {
    // use super::*;

    #[test]
    fn test_sth() {}
}
