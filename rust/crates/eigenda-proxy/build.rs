use std::env;
use std::path::Path;
use std::process::Command;

fn main() {
    // Get the workspace root directory (eigenda/)
    let manifest_dir = env::var("CARGO_MANIFEST_DIR").unwrap();
    let workspace_root = Path::new(&manifest_dir)
        .parent()
        .and_then(|p| p.parent())
        .and_then(|p| p.parent())
        .expect("Failed to find workspace root");

    let api_proxy_dir = workspace_root.join("api/proxy");
    let proxy_binary_path = api_proxy_dir.join("bin/eigenda-proxy");

    println!("cargo:rerun-if-changed={}", api_proxy_dir.display());
    println!("cargo:rerun-if-changed={}", proxy_binary_path.display());

    // Check if the binary already exists
    if !proxy_binary_path.exists() {
        println!(
            "cargo:warning=Proxy binary not found at {}, building it...",
            proxy_binary_path.display()
        );

        // Build the proxy using make
        let status = Command::new("make")
            .current_dir(&api_proxy_dir)
            .status()
            .expect("Failed to run 'make' in api/proxy directory");

        if !status.success() {
            panic!("Failed to build eigenda-proxy using make");
        }

        // Verify the binary was created
        if !proxy_binary_path.exists() {
            panic!(
                "Proxy binary was not created at expected path: {}",
                proxy_binary_path.display()
            );
        }
    }

    // Set the environment variable that will be used by include_bytes!
    let proxy_binary_absolute = proxy_binary_path
        .canonicalize()
        .expect("Failed to canonicalize proxy binary path");

    println!(
        "cargo:rustc-env=EMBEDDED_PROXY_BINARY_PATH={}",
        proxy_binary_absolute.display()
    );
    println!(
        "cargo:warning=Embedding proxy binary from: {}",
        proxy_binary_absolute.display()
    );
}
