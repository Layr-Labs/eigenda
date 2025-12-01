//! EigenDA Proxy Client Library
//!
//! This crate provides a client for interacting with an [EigenDA proxy](https://github.com/Layr-Labs/eigenda/tree/master/api/proxy)
//! It supports storing and retrieving blob data through the EigenDA network
//! using standard commitments and certificates.

pub mod client;
pub use client::{EigenDaProxyConfig, ProxyClient, ProxyError};
