//! Core EigenDA cryptographic verification primitives
//!
//! This module provides the fundamental cryptographic verification components for
//! EigenDA certificates and blob data. It implements the low-level verification
//! algorithms following the EigenDA protocol specification.
//!
//! ## Module Structure
//!
//! This crate contains the core verification primitives:
//!
//! - **[`cert`]** - Certificate cryptographic verification
//!   - BLS signature aggregation and verification
//!   - Stake-weighted quorum validation
//!   - Security threshold enforcement
//!   - Operator state consistency checks
//!
//! - **[`blob`]** - Blob data integrity verification
//!   - KZG polynomial commitment verification
//!   - Blob encoding validation
//!
//! ## Architecture
//!
//! This module focuses on the cryptographic core of EigenDA verification and does
//! not handle:
//! - Ethereum state extraction and proof verification
//! - Rollup-specific integration logic
//! - Certificate recency validation (handled by higher-level adapters)
//!
//! The verification functions expect pre-validated inputs and focus purely on
//! cryptographic correctness.
//!
//! ## References
//!
//! - [EigenDA Protocol Specification](https://docs.eigenlayer.xyz/eigenda/overview/)
//! - [Certificate Verification Reference](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/integrations/cert/libraries/EigenDACertVerificationLib.sol)

pub mod blob;
pub mod cert;
