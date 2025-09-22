//! Sovereign SDK adapter for EigenDA data availability layer.
//!
//! This crate provides integration between Sovereign SDK rollups and EigenDA,
//! enabling efficient blob storage and verification on the EigenDA network.

mod eigenda;
mod ethereum;
#[cfg(feature = "native")]
/// Service layer providing EigenDA client functionality and Ethereum integration.
pub mod service;
/// Core types and trait implementations for EigenDA data availability specification.
pub mod spec;
/// Cryptographic verification of transaction inclusion and completeness proofs.
pub mod verifier;
