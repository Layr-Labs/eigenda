pub mod client;
pub use client::{EigenDaProxyConfig, ProxyClient, ProxyError};

#[cfg(feature = "managed-proxy")]
pub mod managed_proxy;
#[cfg(feature = "managed-proxy")]
pub use managed_proxy::ManagedProxy;
