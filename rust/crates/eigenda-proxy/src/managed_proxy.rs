//! Managed Proxy
//!
//! This module provides the ManagedProxy type for managing an eigenda-proxy binary.
//! It is only available when the `managed-proxy` feature is enabled.

use std::path::PathBuf;
use std::process::Stdio;
use tokio::process::{Child, Command};

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

    /// Start the embedded proxy and monitor it in the background.
    /// This spawns the process and returns the Child handle for further management.
    pub async fn start(&self, args: &[&str]) -> Result<Child, std::io::Error> {
        let binary_path = self.binary_path.clone();

        // Spawn the process
        let child = Command::new(&binary_path)
            .args(args)
            // Redirect stdout and stderr to null for now to not clutter output.
            // If needed, we could allow user to specify log file paths or pipe to parent stdout/stderr.
            .stdout(Stdio::null())
            .stderr(Stdio::null())
            .kill_on_drop(true)
            .spawn()?;

        Ok(child)
    }
}

#[cfg(test)]
mod tests {
    use std::os::unix::process::ExitStatusExt;

    use super::*;

    #[tokio::test]
    async fn test_proxy_version() {
        let mut proxy = ManagedProxy::new()
            .unwrap()
            .start(&["--version"])
            .await
            .unwrap();

        let status = proxy.wait().await.unwrap();
        assert!(status.success());
    }

    #[tokio::test]
    async fn test_start_and_kill_memstore_proxy() {
        let mut proxy = ManagedProxy::new()
            .unwrap()
            .start(&[
                "--memstore.enabled",
                "--apis.enabled=standard",
                "--eigenda.g1-path=../../../resources/srs/g1.point",
            ])
            .await
            .unwrap();

        // Give the proxy a moment to start up
        tokio::time::sleep(std::time::Duration::from_millis(3000)).await;

        let status = proxy.try_wait().unwrap();
        assert!(status.is_none(), "Proxy exited prematurely");

        proxy.start_kill().unwrap();

        let status = proxy.wait().await.unwrap();
        assert!(status.signal() == Some(9), "Proxy was not killed properly");
    }
}
