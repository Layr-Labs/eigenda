//! Build script for eigenda-proxy crate
//! This script downloads the eigenda-proxy binary during build time
//! if the `managed-proxy` feature is enabled.
//! It places the binary in the OUT_DIR and sets an environment
//! variable `EIGENDA_PROXY_PATH` pointing to its location.
//! The ManagedProxy struct in the crate uses this path to launch
//! the embedded proxy.

fn main() {
    // Only download and setup the binary if the managed-proxy feature is enabled
    #[cfg(feature = "managed-proxy")]
    {
        use sha2::{Digest, Sha256};
        use std::env;
        use std::fs;
        use std::io::Write;
        use std::path::Path;

        let out_dir = env::var("OUT_DIR").expect("OUT_DIR not set");
        let binary_path = Path::new(&out_dir).join("eigenda-proxy");

        // Check if binary already exists to avoid re-downloading on every build
        if binary_path.exists() {
            println!("cargo:warning=eigenda-proxy binary already exists, skipping download");
            println!(
                "cargo:rustc-env=EIGENDA_PROXY_PATH={}",
                binary_path.display()
            );
            return;
        }

        let os = env::consts::OS;
        let arch = env::consts::ARCH;
        // Download URL for the eigenda-proxy binary
        // TODO(samlaf): once https://github.com/Layr-Labs/eigenda/pull/2379 is merged and the next release is cut,
        // update this URL to point to the latest eigenda release packaged proxy binary instead of this test one.
        let (download_url, sha256checksum) = match (os, arch) {
            ("macos", "aarch64") => (
                "https://github.com/samlaf/test-ci/releases/download/v0.1.2/eigenda-proxy-darwin-arm64",
                "3b72f724c51dec34379f85bd722ec9a021a3dcb07da937ca34674240ef4c3851",
            ),
            ("linux", "x86_64") => (
                "https://github.com/samlaf/test-ci/releases/download/v0.1.2/eigenda-proxy-linux-amd64",
                "b2d6e32d72fb4f88b8417bd7c85be9d64210d3b37c01ecfb7f6c48d741d3a6b4",
            ),
            _ => panic!(
                "Unsupported platform: {os}-{arch}. Only macOS ARM64 and Linux x86_64 are supported."
            ),
        };

        println!("cargo:warning=Downloading eigenda-proxy binary from {download_url}");

        // Download the binary
        let response = reqwest::blocking::get(download_url)
                .unwrap_or_else(|e| {
                    panic!("Failed to download eigenda-proxy binary from '{download_url}': {e}. Please check your network connectivity and ensure the URL is accessible.");
                });

        if !response.status().is_success() {
            panic!(
                "Failed to download eigenda-proxy: HTTP {}",
                response.status()
            );
        }

        let bytes = response.bytes().expect("Failed to read response bytes");

        // Verify SHA-256 checksum
        let mut hasher = Sha256::new();
        hasher.update(&bytes);
        let computed_hash = format!("{:x}", hasher.finalize());

        if computed_hash != sha256checksum {
            panic!(
                "SHA-256 checksum mismatch for eigenda-proxy binary!\n\
                    Expected: {sha256checksum}\n\
                    Computed: {computed_hash}\n\
                    The downloaded binary may be corrupted or compromised."
            );
        }

        println!("cargo:warning=SHA-256 checksum verified: {computed_hash}");

        // Write binary to OUT_DIR
        let mut file =
            fs::File::create(&binary_path).expect("Failed to create eigenda-proxy binary file");
        file.write_all(&bytes)
            .expect("Failed to write eigenda-proxy binary");

        // Make the binary executable on Unix systems
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let mut perms = file
                .metadata()
                .expect("Failed to get file metadata")
                .permissions();
            perms.set_mode(0o755);
            fs::set_permissions(&binary_path, perms).expect("Failed to set executable permissions");
        }

        println!(
            "cargo:warning=Downloaded eigenda-proxy binary to: {}",
            binary_path.display()
        );

        // Set environment variable pointing to the binary location
        println!(
            "cargo:rustc-env=EIGENDA_PROXY_PATH={}",
            binary_path.display()
        );
    }

    // Rerun build script if the download URL changes (though it's hardcoded)
    println!("cargo:rerun-if-changed=build.rs");
}
