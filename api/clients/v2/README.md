# EigenDA Core Clients v2

![Core Client Diagram](assets/core_clients_v2.svg)

## Overview

The EigenDA Core Clients v2 library provides a comprehensive Go SDK for interacting with the EigenDA network. This library enables developers to store and retrieve data blobs on EigenDA's decentralized storage network with built-in verification, authentication, and payment management.

## Architecture

The library consists of several key components:

### Core Components

- **DisperserClient**: Manages communication with disperser servers for blob storage
- **PayloadDisperser**: High-level interface for dispersing data payloads to EigenDA
- **PayloadRetriever**: Interface for retrieving data from EigenDA (supports both relay and validator retrieval)
- **CertVerifier**: Verifies EigenDA certificates and blob commitments
- **Accountant**: Manages payment reservations and on-demand payment authentication

### Retrieval Methods

1. **Relay Retrieval** (Recommended): Retrieves data from EigenDA relay network for optimal performance
2. **Validator Retrieval**: Direct retrieval from EigenDA validator nodes as an alternative method

### Verification System

- **Certificate Verification**: Validates EigenDA certificates using on-chain data
- **Commitment Verification**: Verifies KZG commitments for data integrity
- **Proof Verification**: Validates cryptographic proofs for data chunks

## Quick Start

### Basic Usage Example

```go
package main

import (
    "context"
    "fmt"
    
    "github.com/Layr-Labs/eigenda/api/clients/v2"
    "github.com/Layr-Labs/eigenda/api/clients/v2/payloaddispersal"
    "github.com/Layr-Labs/eigenda/api/clients/v2/payloadretrieval"
)

func main() {
    ctx := context.Background()
    privateKey := "your_private_key_here"
    
    // Create payload disperser
    payloadDisperser, err := createPayloadDisperser(privateKey)
    if err != nil {
        panic(fmt.Sprintf("create payload disperser: %v", err))
    }
    defer payloadDisperser.Close()
    
    // Disperse data to EigenDA
    data := []byte("Hello EigenDA!")
    eigenDACert, err := payloadDisperser.SendPayload(ctx, data)
    if err != nil {
        panic(fmt.Sprintf("send payload: %v", err))
    }
    
    // Create payload retriever
    payloadRetriever, err := createRelayPayloadRetriever()
    if err != nil {
        panic(fmt.Sprintf("create relay payload retriever: %v", err))
    }
    defer payloadRetriever.Close()
    
    // Retrieve data from EigenDA
    retrievedData, err := payloadRetriever.GetPayload(ctx, eigenDACert)
    if err != nil {
        panic(fmt.Sprintf("get payload: %v", err))
    }
    
    fmt.Printf("Original: %s\n", string(data))
    fmt.Printf("Retrieved: %s\n", string(retrievedData))
}
```

## Configuration

### PayloadClientConfig

Configure how payloads are processed and stored:

```go
config := &clients.PayloadClientConfig{
    PayloadPolynomialForm: codecs.PolynomialFormEval,  // Data processing form
    BlobVersion:           0,                           // Blob version from threshold registry
}
```

### Network Configuration

For different networks, configure the appropriate endpoints and contract addresses:

```go
// Holesky Testnet (example)
const (
    ethRPCURL                    = "https://ethereum-holesky-rpc.publicnode.com"
    disperserHostname           = "disperser-testnet-holesky.eigenda.xyz" 
    certVerifierRouterAddress   = "0x7F40A8e1B62aa1c8Afed23f6E8bAe0D340A4BC4e"
    registryCoordinatorAddress  = "0x53012C69A189cfA2D9d29eb6F19B32e0A2EA3490"
)
```

## Payment Methods

The library supports two payment methods:

1. **Reservations**: Pre-paid storage capacity for consistent usage
2. **On-demand Payments**: Pay-per-use model by sending funds to the payment vault

Authentication is handled automatically by the `Accountant` component based on your configured payment method.

## Key Features

- **Automatic Retry Logic**: Built-in retry mechanisms for network failures
- **Concurrent Operations**: Parallel processing for improved performance  
- **Metrics Integration**: Optional metrics collection with `SequenceProbe`
- **Flexible Configuration**: Customizable timeouts, retry policies, and network settings
- **Security**: Built-in cryptographic verification and authentication
- **Multiple Retrieval Options**: Choose between relay or validator retrieval based on your needs

## Examples

See the `examples/` directory for comprehensive usage examples:

- `client_construction.go`: How to construct and configure clients
- `example_relay_retrieval_test.go`: Relay-based data retrieval
- `example_validator_retrieval_test.go`: Validator-based data retrieval

## Dependencies

- Go 1.19 or later
- Ethereum client for on-chain interactions
- KZG cryptographic libraries for proof generation/verification

## Support

For questions and support, please refer to the main EigenDA documentation or open an issue in the repository.
