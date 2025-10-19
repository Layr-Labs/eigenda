# EigenDA-Sovereign SDK DA Adapter Backend

[![Rust](https://img.shields.io/badge/rust-1.88%2B-orange.svg)](https://www.rust-lang.org)
[![License](https://img.shields.io/badge/license-MIT%20OR%20Apache--2.0-blue.svg)](#license)

Implements the necessary [EigenDA](https://docs.eigencloud.xyz/products/eigenda/core-concepts/overview) backend infrastructure to support the implementation of a Data Availability Adapter for Sovereign SDK based on EigenLayer's EigenDA, enabling rollups to use it as their data availability layer with full cryptographic verification.

## 🏗️ Architecture

The project is built using a modular architecture with specialized crates:

### Core Crates

| Crate | Purpose | Key Features |
|-------|---------|--------------|
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
git clone https://github.com/Layr-Labs/eigenda.git
cd eigenda/rust

# Build all crates
cargo build --release

# Run tests
cargo test
```

## ⚙️ Configuration

The crates provide modular components for EigenDA integration that can be composed based on your rollup's needs. Key configuration points include:

- **Ethereum RPC endpoint** for contract interaction
- **EigenDA Proxy URL** for blob operations
- **Rollup namespace** for transaction filtering

## 🔧 How It Works

These crates provide the foundational components needed to integrate EigenDA with Sovereign SDK rollups:

### Core Capabilities

**Ethereum Integration** (`eigenda-ethereum`)
- Contract interaction and state queries
- Ethereum block monitoring
- State proof generation

**Proxy Communication** (`eigenda-proxy`)
- Blob submission and retrieval
- Certificate management
- Retry logic and error handling

**Cryptographic Verification** (`eigenda-verification`)
- ✅ EigenDA certificate validation
- ✅ BLS aggregate signature verification
- ✅ KZG commitment proof validation
- ✅ Ethereum state proof verification
- ✅ Operator stake extraction and validation

**SRS Data** (`eigenda-srs-data`)
- BN254 curve parameters for KZG operations
- Structured reference string management

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


## 🛠️ Development

### Project Structure

```
eigenda/rust/
├── crates/
│   ├── eigenda-ethereum/        # Ethereum contract utilities
│   ├── eigenda-proxy/           # EigenDA proxy client
│   ├── eigenda-verification/    # Cryptographic verification
│   └── eigenda-srs-data/        # Structured reference string data
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
