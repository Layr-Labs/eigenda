//! EigenDA client library for blob and certificate verification
//!
//! This crate provides functionality for working with EigenDA blobs, including
//! verification of blob certificates and handling of cryptographic proofs.

pub mod cert;
/// Error types for EigenDA verification.
pub mod error;
/// Certificate extraction logic and state data processing.
pub mod extraction;
pub mod verification;
