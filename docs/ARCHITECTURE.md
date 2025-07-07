# EigenDA Architecture

> This document describes the architecture of EigenDA v2, a secure, high-throughput, decentralized data availability service built on Ethereum and EigenLayer.

## Design Goals

EigenDA is designed to:

1. **Provide data availability guarantees** - Ensure that data remains available for a specified duration with cryptographic guarantees
2. **Scale horizontally** - System capacity grows linearly with the number of operators
3. **Minimize trust assumptions** - Leverage Ethereum's security through EigenLayer restaking
4. **Enable efficient verification** - Allow light clients to verify data availability without downloading full data
5. **Support high throughput** - Handle gigabytes per second of data throughput
6. **Maintain composability** - Integrate seamlessly with rollups and other Layer 2 solutions

## System Overview

EigenDA operates as a marketplace where:
- **Rollups** pay to store data temporarily with availability guarantees
- **Operators** provide storage and bandwidth in exchange for rewards
- **Ethereum** provides consensus and economic security through EigenLayer

The system uses erasure coding to distribute data across operators, requiring only a fraction of operators to be online for data recovery.

## Architecture

### Components

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│   Rollup    │────▶│ EigenDA Proxy│────▶│  Disperser  │
└─────────────┘     └──────────────┘     └──────┬──────┘
                                                 │
                    ┌────────────────────────────┼────────────────────────────┐
                    │                            │                            │
                    ▼                            ▼                            ▼
             ┌──────────┐                ┌──────────┐                ┌──────────┐
             │ DA Node  │                │ DA Node  │      ···       │ DA Node  │
             └──────────┘                └──────────┘                └──────────┘
                    │                            │                            │
                    └────────────────────────────┼────────────────────────────┘
                                                 │
                                                 ▼
                                          ┌─────────────┐
                                          │  Retriever  │
                                          └─────────────┘
```

#### 1. EigenDA Proxy (`/api/proxy`)
- **Purpose**: REST API gateway for rollup integration
- **Responsibilities**:
  - Accept blob dispersal requests from rollups
  - Handle encoding and commitment generation
  - Manage retries and fallback to alternative storage
  - Verify retrieved data integrity
- **Key Features**:
  - Multiple commitment schemas (Optimism Alt-DA, Simple, Generic)
  - Caching layer for frequently accessed blobs
  - Storage fallback (S3, EigenDA) for reliability

#### 2. Disperser (`/disperser`)
- **Purpose**: Coordinate blob distribution across the network
- **Responsibilities**:
  - Encode blobs using Reed-Solomon erasure coding
  - Calculate chunk assignments based on operator stakes
  - Distribute chunks to DA nodes
  - Aggregate signatures into attestations
  - Submit batch confirmations to Ethereum
- **Components**:
  - **API Server**: gRPC/REST endpoints for blob submission
  - **Batcher**: Groups blobs into batches for efficiency
  - **Encoder**: Performs erasure coding and KZG commitments
  - **Dispatcher**: Manages chunk distribution to operators
  - **Aggregator**: Collects and aggregates BLS signatures

#### 3. DA Nodes (`/node`)
- **Purpose**: Store and serve blob chunks
- **Responsibilities**:
  - Validate received chunks using KZG proofs
  - Store chunks with time-to-live (TTL)
  - Sign attestations confirming chunk receipt
  - Serve chunks during retrieval
- **Storage**: LittDB - custom key-value store optimized for write-heavy workloads
- **Networking**: gRPC for dispersal/retrieval, P2P for future enhancements

#### 4. Smart Contracts (`/contracts`)
- **Purpose**: On-chain protocol coordination and security
- **Key Contracts**:
  - **EigenDAServiceManager**: Main protocol contract
  - **RegistryCoordinator**: Operator registration and stake management
  - **PaymentVault**: Handles blob payments and operator rewards
  - **EigenDACertVerifier**: Verifies blob availability certificates
- **Integration**: Leverages EigenLayer for restaking and slashing

#### 5. Retriever (`/retriever`)
- **Purpose**: Reconstruct blobs from distributed chunks
- **Process**:
  1. Query DA nodes for chunks based on blob key
  2. Verify chunk integrity using KZG proofs
  3. Decode chunks using Reed-Solomon decoding
  4. Reconstruct and return original blob

#### 6. Relay (`/relay`)
- **Purpose**: Alternative retrieval path for improved performance
- **Features**:
  - Stores complete blobs (not chunks)
  - Provides fast retrieval without reconstruction
  - Acts as a cache layer for frequently accessed data

### Data Flow

#### Dispersal (Write Path)

1. **Submission**: Rollup submits blob to EigenDA Proxy
2. **Encoding**: Proxy forwards to Disperser, which:
   - Generates KZG commitment for the blob
   - Applies Reed-Solomon encoding to create chunks
   - Calculates chunk assignments based on operator stakes
3. **Distribution**: Disperser sends chunks to assigned DA nodes
4. **Attestation**: Each DA node:
   - Verifies chunk validity using KZG proof
   - Stores chunk in LittDB with TTL
   - Returns BLS signature as attestation
5. **Aggregation**: Disperser aggregates signatures when quorum is reached
6. **Confirmation**: Batch of attestations submitted to Ethereum
7. **Certificate**: Availability certificate returned to rollup

#### Retrieval (Read Path)

1. **Request**: Client requests blob using certificate
2. **Chunk Discovery**: Retriever identifies which operators hold chunks
3. **Download**: Retriever fetches minimum required chunks
4. **Verification**: Each chunk verified against KZG commitment
5. **Reconstruction**: Reed-Solomon decoding reconstructs original blob
6. **Validation**: Final blob validated against certificate commitment

### Encoding and Security

#### Reed-Solomon Erasure Coding

EigenDA uses systematic Reed-Solomon encoding:
- **Data chunks (k)**: Original data split into k pieces
- **Parity chunks (n-k)**: Additional redundancy chunks
- **Recovery threshold**: Any k out of n chunks can reconstruct data

Example with 3/5 encoding:
```
Original: [A, B, C]           (3 data chunks)
Encoded:  [A, B, C, P1, P2]   (2 parity chunks)
Recovery: Any 3 chunks sufficient
```

#### KZG Polynomial Commitments

- **Purpose**: Enable trustless chunk verification
- **Properties**:
  - Commitment is 48 bytes regardless of blob size
  - Individual chunks verifiable without full blob
  - Supports efficient batch verification
- **Implementation**: BN254 elliptic curve with trusted setup

#### BLS Signature Aggregation

- **Purpose**: Efficient on-chain verification of operator attestations
- **Process**:
  1. Each operator signs with BLS private key
  2. Signatures aggregated into single signature
  3. One on-chain verification for entire operator set

### State Management

#### On-Chain State (Ethereum)

- **Operator Registry**: Stakes, public keys, and metadata
- **Batch Headers**: Confirmation of dispersed batches
- **Payment Records**: Blob fees and operator rewards
- **Quorum Configuration**: Security parameters per quorum

#### Off-Chain State

##### Disperser State
- **Blob Metadata**: Status, confirmations, expiry
- **Batch State**: Current batch composition
- **Operator Assignments**: Chunk distribution mappings

##### Node State
- **Chunk Store**: Blob chunks with metadata
- **Expiration Index**: TTL-based cleanup tracking
- **Attestation Cache**: Recent signatures

### Quorum System

EigenDA supports multiple quorums with different security properties:

```
Quorum 0: High Security
- Adversary Threshold: 33%
- Confirmation Threshold: 67%
- Use Case: Critical rollup data

Quorum 1: Balanced
- Adversary Threshold: 25%
- Confirmation Threshold: 55%
- Use Case: Standard applications

Quorum 2: High Throughput
- Adversary Threshold: 20%
- Confirmation Threshold: 40%
- Use Case: High-volume, lower-value data
```

### Performance Optimizations

#### Encoding Pipeline
1. **Amortized KZG**: Batch proving reduces commitment cost
2. **ICICLE GPU**: Hardware acceleration for encoding
3. **Streaming**: Process blobs without full buffering

#### Storage Architecture
1. **LittDB**: Write-optimized key-value store
2. **Chunk Bundling**: Reduce storage overhead
3. **Lazy Deletion**: Mark-and-sweep for expired data

#### Network Optimization
1. **Parallel Distribution**: Concurrent chunk sending
2. **Connection Pooling**: Reuse gRPC connections
3. **Adaptive Timeouts**: Dynamic timeout adjustment

### Failure Modes and Recovery

#### Operator Failures
- **Below Threshold**: System continues normally
- **Above Threshold**: Temporary unavailability until operators recover
- **Permanent Loss**: Data unrecoverable if too many operators fail

#### Disperser Failures
- **Transient**: Retries handle temporary issues
- **Crash**: State recovered from persistent storage
- **Byzantine**: On-chain slashing for malicious behavior

#### Network Partitions
- **Partial**: Degraded throughput but continued operation
- **Complete**: System halts until connectivity restored

### Monitoring and Observability

#### Metrics
- **Throughput**: Blobs/sec, bytes/sec per component
- **Latency**: Percentiles for each operation stage
- **Success Rate**: Dispersal and retrieval success ratios
- **Storage**: Utilization and growth trends

#### Distributed Tracing
- Request flow tracking across components
- Bottleneck identification
- Error propagation analysis

### Future Enhancements

1. **Decentralized Disperser**: Remove single point of failure
2. **P2P Chunk Routing**: Direct operator communication
3. **Adaptive Encoding**: Dynamic redundancy based on demand
4. **Cross-Rollup Deduplication**: Shared data optimization
5. **ZK Light Clients**: Succinct availability proofs

## Security Considerations

### Trust Model

1. **Honest Majority**: Assumes honest operators control majority stake per quorum
2. **Ethereum Security**: Inherits security from L1 through EigenLayer
3. **Cryptographic Hardness**: Relies on discrete log and pairing assumptions

### Attack Vectors and Mitigations

1. **Data Withholding**
   - **Attack**: Operators claim to store but don't serve data
   - **Mitigation**: Slashing conditions and redundancy

2. **Sybil Attacks**
   - **Attack**: Single entity controls multiple operators
   - **Mitigation**: Stake requirements and EigenLayer delegation

3. **Denial of Service**
   - **Attack**: Flood system with invalid requests
   - **Mitigation**: Rate limiting and payment requirements

## Conclusion

EigenDA's architecture achieves its design goals through:
- **Modular design** enabling independent component scaling
- **Cryptographic techniques** providing trustless verification
- **Economic incentives** aligning operator and user interests
- **Ethereum integration** leveraging L1 security

The system demonstrates that high-throughput data availability is achievable without compromising decentralization or security.