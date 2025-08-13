//! Following eigenda-contracts/lib/eigenlayer-middleware/src/BLSSignatureChecker.sol
//! from 6797f3821db92c2214aaa6f137d94c603011ac2a lib/eigenlayer-middleware (v0.5.4-mainnet-rewards-v2-1-g6797f38)

#![no_std]
extern crate alloc;

pub mod bitmap;
mod check;
pub mod convert;
pub mod error;
pub mod hash;
mod signature;
pub mod types;
pub mod verification;
