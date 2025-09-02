# sov-eigenda-adapter

`sov-eigenda-adapter` is an adapter making [EigenDA](https://docs.eigencloud.xyz/products/eigenda/core-concepts/overview) compatible with the Sovereign SDK.

## EigenDA Integration

The `sov-eigenda-adapter` integrates with EigenDA's certificate verification system and requires connection to Ethereum mainnet for contract state verification.

## How it Works

All of `sov-eigenda-adapter` boils down to two trait implementations:
 - [`DaService`](https://github.com/Sovereign-Labs/sovereign-sdk/blob/nightly/crates/rollup-interface/src/node/da.rs#L112)
 - [`DaVerifier`](https://github.com/Sovereign-Labs/sovereign-sdk/blob/nightly/crates/rollup-interface/src/state_machine/da.rs#L56)

## The DaService Trait

The job of the `DaService` is to allow the Sovereign SDK's node software to communicate with a DA layer. It has two related responsibilities. The first is to interact with DA layer nodes - retrieving data for the rollup as it becomes available.
The second is to process that data into the form expected by the `DaVerifier`. 

`sov-eigenda-adapter`'s DA service communicates with both Ethereum nodes and EigenDA proxy services. 
For each Ethereum block, the service:

1. Fetches all transactions and identifies those targeting rollup namespaces
2. Extracts EigenDA certificates from relevant transactions
3. Retrieves blob data from EigenDA proxy using certificate information
4. Gathers Ethereum state proofs for all contracts needed for certificate verification
5. Retrieves state data needed for certificate verification.
6. Packages everything into completeness and inclusion proofs for the verifier

The service handles retries, caching, and rate limiting for both Ethereum and EigenDA interactions.

## The DaVerifier Trait

The `DaVerifier` trait is responsible for verifying a set of `BlobTransactions` fetched from a Data Availability (DA) layer block,
ensuring these transactions are both **complete** and **included**.

### Proving inclusion and completeness of rollups data

EigenDA stores rollup data as blobs with associated certificates that prove the data's availability. Each certificate contains:
- Batch metadata including reference block numbers
- Blob inclusion proofs demonstrating the blob is part of the certified batch
- BLS aggregate signatures from EigenDA operators attesting to data availability

Rollup transactions are submitted to Ethereum with EigenDA certificates embedded in the transaction data. The certificates reference specific Ethereum blocks for verification context.

#### Checking _completeness_ of the data

Completeness verification ensures all transactions belonging to the rollup namespace are captured from the Ethereum block. This involves:

1. **Transaction Root Verification**: Computing the merkle root of all provided transactions and comparing against the Ethereum block's transaction root
2. **Namespace Filtering**: Identifying all transactions targeting the rollup's namespace addresses (batch and proof namespaces)
3. **Certificate state Validation**: Verifying the states used to verify the certificates are valid

#### Checking _inclusion_ of the data

Inclusion verification validates that each blob matches its certificate and is properly extracted:

1. **Certificate Validation**: Each EigenDA certificate is verified against the referenced Ethereum block state using account state
2. **Certificate Recency**: Certificates must be submitted within the punctuality window (reference block number + cert recency window)
3. **Blob Verification**: The actual blob data is verified against the certificate's commitment and inclusion proofs
4. **Sender Verification**: Transaction sender addresses are verified through signature recovery
5. **State Proof Verification**: All contract state used for verification is proven against Ethereum block state roots
