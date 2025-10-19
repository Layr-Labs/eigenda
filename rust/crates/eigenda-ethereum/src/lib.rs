//! EigenDA data extraction from Ethereum contract storage
//!
//! Provides utilities for extracting and decoding data from EigenDA
//! protocol smart contracts deployed on Ethereum. It enables verification of
//! blob certificates by fetching the necessary on-chain state data.
//!
//! ## Architecture
//!
//! The extraction system follows a trait-based approach:
//! - [`StorageKeyProvider`]: Generates storage keys for contract data
//! - [`DataDecoder`]: Decodes storage proofs into typed data structures
//!
//! ## Key Components
//!
//! - **Extractors**: Specialized types for extracting specific data from contracts
//! - **Contract Interfaces**: High-level interfaces for each EigenDA contract
//! - **Storage Helpers**: Utilities for generating Ethereum storage keys
//! - **Decode Helpers**: Utilities for parsing storage proofs
//!
//! ## Contract Data Extracted
//!
//! - Quorum configurations and counts
//! - Operator stake histories and bitmap histories  
//! - Aggregated public key (APK) histories
//! - Blob versioning parameters
//! - Security thresholds and required quorum numbers
//! - Stale stake prevention settings (feature-gated)

/// Smart contract interfaces and data structures for EigenDA contracts.
#[cfg(feature = "native")]
pub mod contracts;

/// Ethereum provider utilities and helper functions.
#[cfg(feature = "native")]
pub mod provider;
