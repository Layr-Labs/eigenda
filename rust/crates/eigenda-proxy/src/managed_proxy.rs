use std::fs;
use std::io::Write;
use std::path::PathBuf;
use std::process::Stdio;
use tokio::process::Command;

/// Embedded binary data
const EMBEDDED_BINARY: &[u8] = include_bytes!("../../../../api/proxy/bin/eigenda-proxy");

pub struct ManagedProxy {
    binary_path: PathBuf,
}

pub struct ProxyHandle {
    binary_path: PathBuf,
    child: tokio::process::Child,
}

impl ManagedProxy {
    /// Create a new ManagedProxy instance by extracting the embedded binary
    pub fn new() -> Result<Self, std::io::Error> {
        // Create a temporary directory for the binary
        let temp_dir = std::env::temp_dir();

        // Add .exe extension to make it executable on Windows
        #[cfg(windows)]
        let binary_path = temp_dir.join("eigenda-proxy-embedded.exe");
        #[cfg(not(windows))]
        let binary_path = temp_dir.join("eigenda-proxy-embedded");

        // Write the embedded binary to the temporary location
        let mut file = fs::File::create(&binary_path)?;
        file.write_all(EMBEDDED_BINARY)?;

        // Make the binary executable (Unix-like systems only)
        // On Windows, the .exe extension handles executability
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let mut perms = file.metadata()?.permissions();
            perms.set_mode(0o755);
            file.set_permissions(perms)?;
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

        Ok(ProxyHandle { binary_path, child })
    }
}

impl ProxyHandle {
    /// Stop the running proxy
    pub async fn stop(&mut self) -> Result<(), std::io::Error> {
        self.child.kill().await
    }

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
            eprintln!("Warning: Failed to kill proxy process: {}", e);
        }

        // Clean up the temporary binary file
        if let Err(e) = fs::remove_file(&self.binary_path) {
            eprintln!(
                "Warning: Failed to remove temporary proxy binary {}: {}",
                self.binary_path.display(),
                e
            );
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
            .start(&["--help"])
            .await
            .unwrap();

        proxy.wait().await.unwrap();
    }
}
