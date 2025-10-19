//! EigenDA client library for blob and certificate verification
//!
//! This crate provides comprehensive functionality for working with EigenDA blobs, including
//! certificate parsing, state extraction from Ethereum contracts, cryptographic verification,
//! and handling of cryptographic proofs.
//!
//! ## Main Components
//!
//! - [`cert`] - Certificate data structures and parsing
//! - [`error`] - Unified error types for verification operations
//! - [`extraction`] - Contract state extraction and proof processing
//! - [`verification`] - Cryptographic verification algorithms (certificates and blobs)

/// Certificate data structures and parsing for EigenDA certificates.
pub mod cert;
/// Error types for EigenDA verification.
pub mod error;
/// Certificate state extraction from Ethereum contract storage proofs.
pub mod extraction;
/// Cryptographic verification algorithms for certificates and blobs.
pub mod verification;
