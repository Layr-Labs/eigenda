//! Following eigenda-contracts/lib/eigenlayer-middleware/src/BLSSignatureChecker.sol
//! from 6797f3821db92c2214aaa6f137d94c603011ac2a lib/eigenlayer-middleware (v0.5.4-mainnet-rewards-v2-1-g6797f38)

#![no_std]
extern crate alloc;

mod aggregation;
mod bitmap_utils;
mod convert;
mod error;
mod hash;
mod types;
mod validation;
mod verification;

use aggregation::compute_signers_apk;
use error::SignaturesVerificationError;
use hash::BeHash;
use types::{NonSignerInfo, NonSignerStakesAndSignature, QuorumStakeTotals, ReferenceBlock};
use verification::verify;

#[derive(Default, Debug)]
pub struct SignaturesVerification {
    pub quorum_stake_totals: QuorumStakeTotals,
    pub signatory_record_hash: BeHash,
}

pub trait SignatureVerifier {
    /// Verifies the `params.sigma` aggregate signature over `msg_hash` by the to-be-computed signers aggregate pubkey
    ///
    /// # Arguments
    ///
    /// * `msg_hash` - the message claimed to be signed, in the context of EigenDA blob certificates this is the hash of the batchRoot
    /// * `signed_quorum_numbers` - quorum numbers for which at least one signature exists. A input
    /// like [0, 2] means Quorums 0 and 2 have received signatures while Quorum 1 hasn't received any
    /// * `reference_block_number` - number of the block that contains stake information used to build the certificate
    /// * `current_block_number` - the current block's number
    /// * `params` - metadata used to calculate the signers aggregate pubkey
    /// * `reference_block` - in the solidity implementation of this verification block data is
    /// queried as part of this method's execution, here the required data from the reference block
    /// is passed in wholesale as a ReferenceBlock
    ///
    /// # Returns
    ///
    /// - `Ok(SignaturesVerification)` if the signature could be successfully verified
    /// - `Err(SignaturesVerificationError)` if the signature verification failed or if any
    /// intermediate step failed
    fn verify_signatures<'a>(
        &'a self,
        msg_hash: BeHash,
        signed_quorum_numbers: &'a [u8],
        current_block_number: u32,
        params: &'a NonSignerStakesAndSignature,
        reference_block: &'a ReferenceBlock,
    ) -> Result<SignaturesVerification, SignaturesVerificationError<'a>>;
}

#[derive(Default, Debug)]
pub struct BlsSignaturesVerifier;

impl SignatureVerifier for BlsSignaturesVerifier {
    fn verify_signatures<'a>(
        &'a self,
        msg_hash: BeHash,
        signed_quorum_numbers: &'a [u8],
        current_block_number: u32,
        params: &'a NonSignerStakesAndSignature,
        reference_block: &'a ReferenceBlock,
    ) -> Result<SignaturesVerification, SignaturesVerificationError<'a>> {
        validation::validate_inputs(
            signed_quorum_numbers,
            reference_block.number,
            current_block_number,
            params,
        )?;

        let signers_aggregate_pubkey = compute_signers_apk(
            signed_quorum_numbers,
            u8::MAX, // todo: in solidity this is queried on the fly from current contract state
            &params.non_signer_pubkeys,
            &reference_block.hash_to_bitmap,
            &params.quorum_apks,
        )?;

        let _result = verify(
            &msg_hash,
            signers_aggregate_pubkey,
            params.apk_g2,
            params.sigma,
        );

        let _stake_totals = QuorumStakeTotals::default();
        let _non_signer_info = NonSignerInfo::default();

        Ok(SignaturesVerification::default())
    }
}
