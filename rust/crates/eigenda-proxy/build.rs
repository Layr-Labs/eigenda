use std::env;
use std::fs;
use std::io::Write;
use std::path::Path;

fn main() {
    let out_dir = env::var("OUT_DIR").expect("OUT_DIR not set");
    let binary_path = Path::new(&out_dir).join("eigenda-proxy");

    // Download URL for the eigenda-proxy binary
    // TODO(samlaf): once https://github.com/Layr-Labs/eigenda/pull/2379 is merged and the next release is cut,
    // update this URL to point to the latest eigenda release packaged proxy binary instead of this test one.
    let download_url = "https://github.com/samlaf/samlaf.github.io/releases/download/test/eigenda-proxy";

    // Check if binary already exists to avoid re-downloading on every build
    if !binary_path.exists() {
        println!("cargo:warning=Downloading eigenda-proxy binary from {}", download_url);

        // Download the binary
        let response = reqwest::blocking::get(download_url)
            .expect("Failed to download eigenda-proxy binary");

        if !response.status().is_success() {
            panic!(
                "Failed to download eigenda-proxy: HTTP {}",
                response.status()
            );
        }

        let bytes = response.bytes().expect("Failed to read response bytes");

        // Write binary to OUT_DIR
        let mut file = fs::File::create(&binary_path)
            .expect("Failed to create eigenda-proxy binary file");
        file.write_all(&bytes)
            .expect("Failed to write eigenda-proxy binary");

        // Make the binary executable on Unix systems
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let mut perms = file.metadata()
                .expect("Failed to get file metadata")
                .permissions();
            perms.set_mode(0o755);
            fs::set_permissions(&binary_path, perms)
                .expect("Failed to set executable permissions");
        }

        println!("cargo:warning=Downloaded eigenda-proxy binary to: {}", binary_path.display());
    }

    // Set environment variable pointing to the binary location
    println!("cargo:rustc-env=EIGENDA_PROXY_PATH={}", binary_path.display());

    // Rerun build script if the download URL changes (though it's hardcoded)
    println!("cargo:rerun-if-changed=build.rs");
}
