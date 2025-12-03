//! Managed Proxy
//!
//! This module provides the ManagedProxy type for managing an eigenda-proxy binary.
//! It is only available when the `managed-proxy` feature is enabled.

use std::path::PathBuf;
use std::process::Stdio;
use tokio::process::Command;

/// Path to the downloaded eigenda-proxy binary (set by build.rs when managed-proxy feature is enabled)
const EIGENDA_PROXY_PATH: &str = env!("EIGENDA_PROXY_PATH");

/// ManagedProxy struct that handles launching the proxy binary as a subprocess.
/// It is currently kept very minimal and doesn't do any monitoring, health checks, piping proxy output, etc.
pub struct ManagedProxy {
    binary_path: PathBuf,
}

impl ManagedProxy {
    /// Create a new ManagedProxy instance using the downloaded binary
    pub fn new() -> Result<Self, std::io::Error> {
        let binary_path = PathBuf::from(EIGENDA_PROXY_PATH);

        // Verify the binary exists
        if !binary_path.exists() {
            return Err(std::io::Error::new(
                std::io::ErrorKind::NotFound,
                format!(
                    "eigenda-proxy binary not found at {}. This should have been downloaded during build.",
                    binary_path.display()
                ),
            ));
        }

        Ok(Self { binary_path })
    }

    /// Start the embedded proxy and monitor it in the background
    /// This spawns the process and monitors it using tokio::select!
    pub async fn start(&self, args: &[&str]) -> Result<ProxyHandle, std::io::Error> {
        let binary_path = self.binary_path.clone();

        // Spawn the process
        let child = Command::new(&binary_path)
            .args(args)
            .stdout(Stdio::null())
            .spawn()?;

        Ok(ProxyHandle { child })
    }
}

/// Handle to the running managed proxy process, which is currently simply
/// a wrapper around the [tokio::process::Child].
pub struct ProxyHandle {
    child: tokio::process::Child,
}

impl ProxyHandle {
    /// Stop the running proxy
    pub async fn stop(&mut self) -> Result<(), std::io::Error> {
        self.child.kill().await
    }

    /// Wait for the proxy process to exit. The proxy is a long-standing process which is not meant to exit normally,
    /// but we recommend calling this inside a tokio::select! along with other tasks to monitor for unexpected exits
    /// and crash the client.
    pub async fn wait(&mut self) -> Result<std::process::ExitStatus, std::io::Error> {
        self.child.wait().await
    }
}

impl Drop for ProxyHandle {
    // Note: We can't await in Drop, so monitor task cleanup happens in background
    // For proper cleanup, users should call stop() before dropping
    fn drop(&mut self) {
        // Attempt to kill the child process
        if let Err(e) = self.child.start_kill() {
            eprintln!("Warning: Failed to kill proxy process: {e}");
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[tokio::test]
    async fn test_basic_start_stop() {
        let mut proxy = ManagedProxy::new()
            .unwrap()
            .start(&["--version"])
            .await
            .unwrap();

        let status = proxy.wait().await.unwrap();
        assert!(status.success());
    }
}
