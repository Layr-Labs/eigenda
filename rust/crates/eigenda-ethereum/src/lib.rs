//! Ethereum integration utilities for EigenDA
//!
//! Provides utilities for interacting with EigenDA smart contracts deployed on Ethereum.
//! This crate focuses on contract bindings and provider functionality for fetching
//! blockchain data.
//!
//! ## Key Components
//!
//! - **[`contracts`]** - Smart contract interfaces and data structures for EigenDA contracts
//! - **[`provider`]** - Ethereum provider utilities and helper functions for fetching state
//!
//! ## Architecture Notes
//!
//! This crate handles the Ethereum interaction layer. For certificate state extraction
//! and verification, see the `eigenda-verification` crate which contains:
//! - Contract storage proof extraction
//! - State data decoding
//! - Cryptographic verification

/// Smart contract interfaces and data structures for EigenDA contracts.
pub mod contracts;

/// Ethereum provider utilities and helper functions.
pub mod provider;
