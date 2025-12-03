//! EigenDA Proxy Client Library
//! This crate provides a client for interacting with an [EigenDA proxy](https://github.com/Layr-Labs/eigenda/tree/master/api/proxy).
//! Although we recommend running and managing the proxy as a separate service, this crate also provides
//! a managed proxy service that will spin up a proxy instance as a subprocess.

pub mod client;
pub use client::{EigenDaProxyConfig, ProxyClient, ProxyError};

#[cfg(feature = "managed-proxy")]
pub mod managed_proxy;
#[cfg(feature = "managed-proxy")]
pub use managed_proxy::ManagedProxy;
