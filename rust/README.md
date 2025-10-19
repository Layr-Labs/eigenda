# Sovereign EigenDA Adapter

[![Rust](https://img.shields.io/badge/rust-1.88%2B-orange.svg)](https://www.rust-lang.org)
[![License](https://img.shields.io/badge/license-MIT%20OR%20Apache--2.0-blue.svg)](#license)

A production-ready adapter that integrates [EigenDA](https://docs.eigencloud.xyz/products/eigenda/core-concepts/overview) with the Sovereign SDK, enabling rollups to use EigenDA as their data availability layer with full cryptographic verification.

## 🏗️ Architecture

The adapter is built using a modular architecture with specialized crates:

### Core Crates

| Crate | Purpose | Key Features |
|-------|---------|--------------|
| **`sov-eigenda-adapter`** | Main adapter implementing Sovereign SDK DA traits | `DaService` and `DaVerifier` implementations |
| **`eigenda-ethereum`** | Ethereum contract interaction | Provider utilities, contract bindings |
| **`eigenda-proxy`** | EigenDA proxy service communication | Blob retrieval, certificate generation, retry logic |
| **`eigenda-verification`** | Cryptographic verification, validation, and state extraction | Certificate parsing, storage proofs, operator stake extraction, BLS signatures, commitment proofs |
| **`eigenda-srs-data`** | Structured reference string data | BN254 curve parameters for KZG commitments |

## 🚀 Quick Start

### Prerequisites

- ✅ **Ethereum Node**: Access to Ethereum mainnet RPC
- ✅ **EigenDA Proxy**: Connection to EigenDA proxy service

```bash
# Clone the repository
git clone https://github.com/eiger-co/sov-eigenda-adapter.git
cd sov-eigenda-adapter

# Build all crates
cargo build --release

# Run tests
cargo test
```

## ⚙️ Configuration

The adapter requires configuration for both Ethereum and EigenDA connections:

```rust
use sov_eigenda_adapter::EigenDaConfig;

let config = EigenDaConfig {
    ethereum_rpc_url: "https://mainnet.infura.io/v3/your-key".to_string(),
    eigenda_proxy_url: "http://localhost:3100".to_string(),
    rollup_namespace: "your-rollup-namespace".to_string(),
    // Additional configuration options...
};
```

## 🔧 How It Works

The adapter implements two core Sovereign SDK traits:

### [`DaService`](https://github.com/Sovereign-Labs/sovereign-sdk/blob/nightly/crates/rollup-interface/src/node/da.rs#L112)

Handles communication with the DA layer:

1. **Ethereum Monitoring** - Watches Ethereum blocks for rollup transactions
2. **Certificate Extraction** - Identifies and extracts EigenDA certificates from transactions
3. **Blob Retrieval** - Fetches blob data from EigenDA proxy using certificates
4. **State Proof Generation** - Gathers Ethereum state proofs for verification
5. **Data Packaging** - Prepares completeness and inclusion proofs for the verifier

### [`DaVerifier`](https://github.com/Sovereign-Labs/sovereign-sdk/blob/nightly/crates/rollup-interface/src/state_machine/da.rs#L56)

Cryptographically verifies DA data integrity:

#### Completeness Verification
- ✅ Transaction root verification against Ethereum block
- ✅ Namespace filtering for rollup-specific transactions
- ✅ Certificate state validation

#### Inclusion Verification  
- ✅ EigenDA certificate validation against Ethereum state
- ✅ Certificate recency within punctuality window
- ✅ Blob commitment verification using KZG proofs
- ✅ BLS aggregate signature verification
- ✅ State proof verification against block state roots

## 🧪 Testing

Run the full test suite:

```bash
# Unit tests
cargo test

# Integration tests
cargo test --test integration

# Benchmarks
cargo bench
```

### Test Categories

- **Unit Tests** - Individual component testing
- **Integration Tests** - End-to-end verification workflows
- **Property Tests** - Fuzz testing for edge cases
- **Performance Tests** - Benchmarking verification operations

## 📊 Examples

Explore the [`examples/`](examples/) directory for complete implementations:

- **[Demo Rollup](examples/demo-rollup/)** - Full rollup implementation using EigenDA

## 🛠️ Development

### Project Structure

```
sov-eigenda-adapter/
├── crates/
│   ├── sov-eigenda-adapter/     # Main adapter implementation
│   ├── eigenda-ethereum/        # Ethereum contract utilities
│   ├── eigenda-proxy/           # EigenDA proxy client
│   ├── eigenda-verification/    # Cryptographic verification
│   └── eigenda-srs-data/        # Structured reference string data
├── examples/                    # Example implementations
│   ├── demo-rollup/            # Complete rollup example
```

### Building from Source

```bash
# Development build
cargo build

# Release build with optimizations
cargo build --release

# Build specific crate
cargo build -p eigenda-verification
```

### Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality  
4. Ensure all tests pass
5. Submit a pull request

## 🔒 Security

This adapter implements production-grade security measures:

- **State Proof Verification** - All contract state is cryptographically proven
- **Certificate Validation** - Full BLS signature verification
- **Punctuality Checks** - Prevents stale certificate acceptance
- **Commitment Verification** - KZG proof validation for blob integrity

## 📝 License

This project is licensed under

- [MIT License](LICENSE)
