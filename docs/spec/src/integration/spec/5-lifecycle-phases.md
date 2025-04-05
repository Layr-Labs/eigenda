
# Lifecycle Phases

## Encoding

This phase occurs inside the eigenda-proxy, because the proxy acts as the “bridge” between the Rollup Domain and Data Availability Domain (see [lifecycle](./2-rollup-payload-lifecycle.md) diagram).

A `payload` consists of an arbitrary byte array. The DisperseBlob endpoint accepts an `encodedPayload`, which needs to be a bn254 field element array.

## BlobHeader Construction

The BlobHeader contains 4 main sections that we need to construct.

**Version**

The blobHeader version refers to one of the versionedBlobParams struct defined in the [EigenDAThresholdRegistry](./4-contracts.md#eigendathreshold-registry) contract.

**QuorumNumbers**

QuorumNumbers represents a list a quorums that are required to sign over and make the blob available. Quorum 0 represents the ETH quorum, quorum 1 represents the EIGEN quorum, and both of these are required. Custom quorums can also be added to this list.

**BlobCommitment**

The BlobCommitment is binding commitment for an EigenDA Blob. Because of the length field, a BlobCommitment can only represent a single unique `Blob`. It is also used by the disperser to convince EigenDA validators that the chunks that they have received are indeed part of the blob (or it's reed-solomon extension). It can either be computed locally from the blob, or one can ask the disperser to generate it via the `GetBlobCommitment` endpoint.

```protobuf
message BlobCommitment {
  // A G1 commitment to the blob data.
  bytes commitment = 1;
  // A G2 commitment to the blob data.
  bytes length_commitment = 2;
    // Used with length_commitment to assert the correctness of the `length` field below.
  bytes length_proof = 3;
  // Length in bn254 field elements (32 bytes) of the blob. Must be a power of 2.
  uint32 length = 4;
}
```

Unlike Ethereum blobs which are all 128KiB, EigenDA blobs can be any power of 2 length between 32KiB and 16MiB (currently), and so the `commitment` alone is not sufficient to prevent certain attacks:

- Why is a commitment to the length of the blob necessary?
    
    There are different variants of the attack. The basic invariant the system needs to satisfy is that with the chunks from sufficient set of validators, you can get back the full blob. So the total size of the chunks held by these validators needs to exceed the blob size. If I don't know the blob size (or at least an upper bound), there's no way for the system to validate this invariant.
    Here’s a simple example. Assume a network of 8 DA nodes, and coding ratio 1/2. For a `blob` containing 128 field elements (FEs), each node gets 128*2/8=32 FEs, meaning that any 4 nodes can join forces and reconstruct the data. Now assume a world without length proof; a malicious disperser receives the same blob, uses the same commitment, but claims that the blob only had length 4 FEs. He sends each node 4*2/8=1 FE. The chunks submitted to the nodes match the commitment, so the nodes accept and sign over the blob’s batch. But now there are only 8 FEs in the system, which is not enough to reconstruct the original blob (need at least 128 for that).
    

> Note that the length here is the length of the blob (power of 2), which is different from the payload_length encoded as part of the `PayloadHeader` in the `blob` itself (see the [encoding section](#encoding)).
> 

**PaymentHeader**

The paymentHeader specifies how the blob dispersal to the network will be paid for. There are 2 modes of payment, the permissionless pay-per-blob model and the permissioned reserved-bandwidth approach. See the [Payments](https://docs.eigenda.xyz/releases/payments) release doc for full details; we will only describe how to set these 4 fields here.

```protobuf
message PaymentHeader {
  // The account ID of the disperser client. This should be a hex-encoded string of the ECDSA public key
  // corresponding to the key used by the client to sign the BlobHeader.
  string account_id = 1;
  // UNIX timestamp in nanoseconds at the time of the dispersal request.
  // Used to determine the reservation period, for the reserved-bandwidth payment model.
  int64 timestamp = 2;
  // Total amount of tokens paid by the requesting account, including the current request.
  // Used for the pay-per-blob payment model.
  bytes cumulative_payment = 3;
}
```

Users who want to pay-per-blob need to set the cumulative_payment. Users who have already paid for reserved-bandwidth should instead set the timestamp. If both are set, reserved-bandwidth will be used first, and cumulative_payment only used if the entire bandwidth for the current reservation period has been used up.

An rpc call to the Disperser’s `GetPaymentState` method can be made to query the current state of an `account_id`. A client can query for this information on startup, cache it, and then update it manually when making pay-per-blob payments. In this way, it can keep track of the cumulative_payment and set it correctly for subsequent dispersals.

## Blob Dispersal

The `DisperseBlob` method takes a `blob` and `blob_header` as input. Dispersal entails taking a blob, reed-solomon encoding it into chunks, dispersing those to the EigenDA nodes, retrieving their signatures, creating a `cert` from them, and returning that cert to the client. The disperser batches blobs for a few seconds before dispersing them to nodes, so an entire dispersal process can exceed 10 seconds. For this reason, the API has been designed asynchronously with 2 relevant methods:

```protobuf
// Async call which queues up the blob for processing and immediately returns.
rpc DisperseBlob(DisperseBlobRequest) returns (DisperseBlobReply) {}
// Polled for the blob status updates, until a terminal status is received
rpc GetBlobStatus(BlobStatusRequest) returns (BlobStatusReply) {}

message DisperseBlobRequest {
  bytes blob = 1;
  common.v2.BlobHeader blob_header = 2;
  bytes signature = 3;
}
message BlobStatusReply {
  BlobStatus status = 1;
  SignedBatch signed_batch = 2;
  BlobVerificationInfo blob_verification_info = 3;
}

// Intermediate states: QUEUED, ENCODED
// Terminal states: CERTIFIED, UNKNOWN, FAILED, INSUFFICIENT_SIGNATURES
enum BlobStatus {
  UNKNOWN = 0; // functionally equivalent to FAILED but for unknown unknown bugs
  QUEUED = 1; // Initial state after a DisperseBlob call returns
  ENCODED = 2; // Reed-Solomon encoded into chunks ready to be dispersed to DA Nodes
  CERTIFIED = 3; // blob has been dispersed and attested by NA nodes
  FAILED = 4; // permanent failure (for reasons other than insufficient signatures)
  INSUFFICIENT_SIGNATURES = 5;
}
```

After a successful DisperseBlob rpc call, `BlobStatus.QUEUED` is returned. To retrieve a `cert`, the `GetBlobStatus` rpc shall be polled until a terminal status is reached. If `BlobStatus.CERTIFIED` is received, the `signed_batch` and `blob_verification_info` fields of the `BlobStatusReply` will be returned and can be used to create the `cert` . Any other terminal status indicates failure, and a new blob dispersal will need to be made.

**Failover to EthDA**

The proxy can be configured to retry `FAILED` dispersal n times, after which it returns to the rollup a `503` HTTP status code which rollup batchers can use to failover to EthDA. See [here](https://github.com/ethereum-optimism/specs/issues/434) for more info.

## BlobStatusReply → Cert

This is not necessarily part of the spec but is currently needed given that the disperser doesn’t actually return a cert, so we need a bit of data processing to transform its returned value into a Cert. The transformation is visualized in the [Ultra High Res Diagram](../spec.md#ultra-high-resolution-diagram). The main difference is just calling the [`getNonSignerStakesAndSignature`](https://github.com/Layr-Labs/eigenda/blob/d9cf91e22b6812f85151f4d83aecc96bae967316/contracts/src/core/EigenDABlobVerifier.sol#L222) helper function within the new `EigenDACertVerifier` contract to create the `NonSignerStakesAndSignature` struct. The following pseudocode below exemplifies this necessary preprocessing step:

```python

class CertV2:
    batch_header: any  # You can specify proper types here
    blob_verification_proof: any
    nonsigner_stake_sigs: any

def get_cert_v2(blob_header_hash, blob_verifier_binding) -> CertV2:
    """
    V2 cert construction pseudocode
    @param blob_header_hash: key used for referencing blob status from disperser
    @param blob_verifier_binding: ABI contract binding used for generating nonsigner metadata
    @return v2_cert: EigenDA V2 certificate used by rollup 
    """
  # Call the disperser for the info needed to construct the cert
    blob_status_reply = disperser_client.get_blob_status(blob_header_hash)
    
    # Validate the blob_header received, since it uniquely identifies
    # an EigenDA dispersal.
    blob_header_hash_from_reply = blob_status_reply
                                                                .blob_verification_info
                                                                .blob_certificate
                                                                .blob_header
                                                                .Hash()
    if blob_header_hash != blob_header_hash_from_reply {
        throw/raise/panic
    }

    # Extract first 2 cert fields from blob status reply
    batch_header = blob_status_reply.signed_batch.batch_header
    blob_verification_proof = blob_status_reply.blob_verification_info
    
    # Construct NonSignerStakesAndSignature
    nonsigner_stake_sigs = blob_verifier_binding.getNonSignerStakesAndSignature(
                                                     blob_status_reply.signed_batch)
                                                 
  return Cert(batch_header, blob_verification_proof, nonsigner_stake_sigs)
```

## Posting to Ethereum

The proxy converts the `cert` to an [`altda-commitment`](./3-datastructs.md#altdacommitment) ready to be submitted to the batcher’s inbox without any further modifications by the rollup stack.

## Retrieval

There are two main blob retrieval paths:

1. decentralized retrieval: retrieve chunks from Validators are recreate the `blob` from them.
2. centralized retrieval: the same [Relay API](https://docs.eigenda.xyz/releases/v2#relay-interfaces) that Validators use to download chunks, can also be used to retrieve full blobs.

EigenDA V2 has a new [Relay API](https://docs.eigenda.xyz/releases/v2#relay-interfaces) for retrieving blobs from the disperser. The `GetBlob` method takes a `blob_key` as input, which is a synonym for `blob_header_hash`. Note that `BlobCertificate` (different from `DACert`!) contains an array of `relay_keys`, which are the relays that can serve that specific blob. A relay’s URL can be retrieved from the [relayKeyToUrl](https://github.com/Layr-Labs/eigenda/blob/9a4bdc099b98f6e5116b11778f0cf1466f13779c/contracts/src/core/EigenDARelayRegistry.sol#L35) function on the EigenDARelayRegistry.sol contract.

## Decoding

Decoding performs the exact reverse operations that [Encoding](#encoding) did.