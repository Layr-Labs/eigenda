//! EigenDA protocol/interaction implementation for Sovereign SDK
//!
//! ## Overview
//!
//! EigenDA is a data availability solution built on EigenLayer that provides:
//! - **High throughput**: Optimized for large-scale data availability
//! - **Cost efficiency**: Reduced costs through erasure coding and off-chain storage  
//! - **Cryptographic security**: BLS signature aggregation and KZG commitments
//! - **Ethereum integration**: Uses Ethereum for coordination and verification
//!
//! ## Architecture
//!
//! The implementation is organized into several key modules:
//!
//! ### [`cert`] - Certificate Management
//! - Parsing and encoding EigenDA certificates (V2/V3 formats)
//! - Solidity type definitions for contract interaction
//! - Certificate validation and metadata extraction
//!
//! ### [`extraction`] - On-chain Data Extraction  
//! - Ethereum contract storage access
//! - Operator stake and quorum information retrieval
//! - Historical data and proof extraction
//!
//! ### [`proxy`] - Network Communication
//! - Direct communication with EigenDA disperser nodes
//! - Blob submission and retrieval operations
//! - Network protocol implementation
//!
//! ### [`verification`] - Cryptographic Verification
//! - Certificate signature verification using BLS aggregation
//! - Blob data integrity checks with KZG proofs
//! - Stake-weighted quorum validation
//! - Security threshold enforcement
//!
//! ## Features
//!
//! - `native` - Mode of execution where non-zkvm operations can be performed like disk or network access
//! - `stale-stakes-forbidden` - Adds stale stake prevention verification
//!
//! ## References
//!
//! - [EigenDA Documentation](https://docs.eigenlayer.xyz/eigenda/overview/)
//! - [EigenDA Contracts](https://github.com/Layr-Labs/eigenda/tree/master/contracts)
//! - [EigenLayer Protocol](https://docs.eigenlayer.xyz/)

pub mod extraction;

#[cfg(feature = "native")]
pub mod proxy;

pub mod cert;
pub mod verification;
