
# Lifecycle Phases

## Encoding

This phase occurs inside the eigenda-proxy, because the proxy acts as the “bridge” between the Rollup Domain and Data Availability Domain (see [lifecycle](./2-rollup-payload-lifecycle.md) diagram).

A `payload` consists of an arbitrary byte array. The DisperseBlob endpoint accepts an `encodedPayload`, which needs to be a bn254 field element array.

## BlobHeader Construction

The BlobHeader contains 4 main sections that we need to construct.

**Version**

The blobHeader version refers to one of the versionedBlobParams struct defined in the [EigenDAThresholdRegistry](./4-contracts.md#eigendathreshold-registry) contract.

**QuorumNumbers**

QuorumNumbers represents a list a quorums that are required to sign over and make the blob available. Quorum 0 represents the ETH quorum, quorum 1 represents the EIGEN quorum, and both of these are always required. Custom quorums can also be added to this list.

**BlobCommitment**

The BlobCommitment is a binding commitment for an EigenDA Blob. Because of the length field, a BlobCommitment can only represent a single unique `Blob`. It is also used by the disperser to convince EigenDA validators that the chunks that they have received are indeed part of the blob (or its reed-solomon extension). It can either be computed locally from the blob, or one can ask the disperser to generate it via the `GetBlobCommitment` endpoint.

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

Users who want to pay-per-blob need to set the cumulative_payment. `timestamp` is used by users who have paid for reserved-bandwidth. If both are set, reserved-bandwidth will be used first, and cumulative_payment only used if the entire bandwidth for the current reservation period has been used up.

An rpc call to the Disperser’s `GetPaymentState` method can be made to query the current state of an `account_id`. A client can query for this information on startup, cache it, and then update it manually when making pay-per-blob payments. In this way, it can keep track of the cumulative_payment and set it correctly for subsequent dispersals.

## Secure Dispersal
### Diagram
![image.png](../../assets/integration/secure-blob-dispersal.png)


### System Flow

1. Using `latest_block_number` (lbn) number fetched from ETH RPC node, *Proxy* calls the router to get the `verifier` address *most likely* (if using `EigenDACertVerifierRouter`) to be committed to by the RBN returned by the EigenDA disperser

2. Using the `verifier`, Proxy fetches the `required_quorums` an embeds them into the `BlobHeader` as part of the disperser request

3. *Proxy* submits the payload blob request to the EigenDA disperser and waits for a `BlobStatusReply` (BSR)

4. While querying the disperser, *Proxy* periodically checks against the confirmation threshold as it’s updated in real-time by the disperser ([reference](#blob-dispersal-with-eigenda-disperser)) using `reference_block_number` (rbn) returned in the `BlobStatusReply`

5. *Proxy* calls the `verifier`'s `certVersion()` method to get the `cert_version`

6.  *Proxy* casts the `DACert` into a structured ABI binding type using the `cert_version` to dictate which certificate representation to use

7. *Proxy* then passes ABI encoded cert bytes via a call to the `verifier`'s `checkDACert` function which returns a `verification_status_code`

8. Using the `verification_status_code`, proxy determines whether to return the certificate (`CertV2Lib.StatusCode.SUCCESS`) or retry a subsequent dispersal attempt


### Adding New Verifiers — Synchronization Risk

There is a synchronization risk that can temporarily cause dispersals to fail when adding a new `verifier'` to the `EigenDACertVerifierRouter` at a future activation block number (`abn'`). If `latest_block < abn'` **and** `rbn >= abn'`, dispersals may fail if the `required_quorums` set differs between `verifier` and `verifier'`. In this case, the quorums included in the client's `BlobHeader` (based on the old verifier) would not match those expected by `checkDACert` (using the new verifier). This mismatch results in **at most** a few failed dispersals, which will resolve once `latest_block >= abn'` and `reference_block_number >= abn'`, ensuring verifier consistency.


### Blob Dispersal With EigenDA Disperser

The `DisperseBlob` method takes a `blob` and `blob_header` as input. Dispersal entails taking a blob, reed-solomon encoding it into chunks, dispersing those to the EigenDA nodes, retrieving their signatures, creating a `DACert` from them, and returning that cert to the client. The disperser batches blobs for a few seconds before dispersing them to nodes, so an entire dispersal process can exceed 10 seconds. For this reason, the API has been designed asynchronously with 2 relevant methods:

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

// Intermediate states: QUEUED, ENCODED, GATHERING_SIGNATURES
// Terminal states: UNKNOWN, COMPLETE, FAILED
enum BlobStatus {
  UNKNOWN = 0; // functionally equivalent to FAILED but for unknown unknown bugs
  QUEUED = 1; // Initial state after a DisperseBlob call returns
  ENCODED = 2; // Reed-Solomon encoded into chunks ready to be dispersed to DA Nodes
  GATHERING_SIGNATURES = 3; // blob chunks are actively being transmitted to validators
  COMPLETE = 4; // blob has been dispersed and attested by DA nodes
  FAILED = 5;
}
```

After a successful DisperseBlob RPC call, the disperser returns `BlobStatus.QUEUED`. To retrieve a cert, the GetBlobStatus RPC should be polled until a terminal status is reached.

If `BlobStatus.GATHERING_SIGNATURES` is returned, the `signed_batch` and `blob_verification_info` fields will be present in the `BlobStatusReply`. These can be used to construct a `DACert`, which may be verified immediately against the configured threshold parameters stored in the `EigenDACertVerifier` contract. If the verification passes, the certificate can be accepted early. If verification fails, polling should continue.

Once `BlobStatus.COMPLETE` is returned, it indicates that the disperser has stopped collecting additional signatures, typically due to reaching a timeout or encountering an issue. While the `signed_batch` and `blob_verification_info` fields will be populated and can be used to construct a `DACert`, the `DACert` could still be invalid if an insufficient amount of signatures were collected in-regards to the threshold parameters.

Any other terminal status indicates failure, and a new blob dispersal will need to be made.

**Failover to Native Rollup DA**

The proxy can be configured to retry `UNKNOWN`, `FAILED`, & `COMPLETE` dispersal `n` times, after which it returns to the rollup a `503` HTTP status code which rollup batchers can use to failover to EthDA or native rollup DA offerings (e.g, arbitrum anytrust). See [here](https://github.com/ethereum-optimism/specs/issues/434) for more info on the OP implementation and [here](https://hackmd.io/@epociask/SJUyIZlZkx) for Arbitrum. 

### BlobStatusReply → Cert

This is not necessarily part of the spec but is currently needed given that the disperser doesn't actually return a cert, so we need a bit of data processing to transform its returned value into a Cert. The transformation is visualized in the [Ultra High Res Diagram](../spec.md#ultra-high-resolution-diagram). 

In the updated implementation, a `CertBuilder` constructs the DA Cert through direct communication with the `OperatorStateRetriever` contract, which provides the necessary information about operator stake states. This approach ensures accurate on-chain data for certificate verification. The following pseudocode demonstrates this process:

```python
class DACert:
    batch_header: any
    blob_verification_proof: any
    nonsigner_stake_sigs: any
    cert_version: uint8
    signedQuorumNumbers: bytes

def get_da_cert(blob_header_hash, operator_state_retriever, cert_version_uint8) -> DACert:
    """
    DA Cert construction pseudocode with OperatorStateRetriever
    @param blob_header_hash: key used for referencing blob status from disperser
    @param operator_state_retriever: ABI contract binding for retrieving operator state data
    @param cert_version_uint8: uint8 version of the certificate format to use
    @return DACert: EigenDA certificate used by rollup 
    """
    # Call the disperser for the info needed to construct the cert
    blob_status_reply = disperser_client.get_blob_status(blob_header_hash)
    
    # Validate the blob_header received, since it uniquely identifies
    # an EigenDA dispersal.
    blob_header_hash_from_reply = blob_status_reply.blob_verification_info.blob_certificate.blob_header.Hash()
    if blob_header_hash \!= blob_header_hash_from_reply:
        throw/raise/panic
    
    # Extract first 2 cert fields from blob status reply
    batch_header = blob_status_reply.signed_batch.batch_header
    blob_verification_proof = blob_status_reply.blob_verification_info
    
    # Get the reference block number from the batch header
    reference_block_number = batch_header.reference_block_number
    
    # Get quorum IDs from the blob header
    quorum_numbers = blob_verification_info.blob_certificate.blob_header.quorum_numbers
    
    # Retrieve operator state data directly from the OperatorStateRetriever contract
    operator_states = operator_state_retriever.getOperatorState(
        reference_block_number,
        quorum_numbers,
        blob_status_reply.signed_batch.signatures
    )
    
    # Construct NonSignerStakesAndSignature using the operator state data
    nonsigner_stake_sigs = construct_nonsigner_stakes_and_signature(
        operator_states,
        blob_status_reply.signed_batch.signatures
    )

    signed_quorum_numbers = blob_status_reply.signed_batch.quorum_numbers
    
    return DACert(batch_header, blob_verification_proof, nonsigner_stake_sigs, cert_version_uint8, signed_quorum_numbers)
```
## Posting to Ethereum

The proxy converts the `DACert` to an [`altda-commitment`](./3-datastructs.md#altdacommitment) ready to be submitted to the batcher’s inbox without any further modifications by the rollup stack.

## Secure Retrieval

### System Diagram
![image.png](../../assets/integration/secure-blob-retrieval.png)


### System Flow

1. A *Rollup Node* queries *Proxy’s* `/get` endpoint to fetch batch contents associated with an encoded DA commitment.

2. *Proxy* decodes the `cert_version` for the DA commitment and uses an internal mapping of `cert_version` ⇒ `cert_abi_struct` to deserialize into the structured binding cert type.

3. *Proxy* submits ABI encoded cert bytes to `EigenDACertVerifier` read call via the `checkDAcert` method, which returns a `verification_status_code`.

4. *Proxy* interprets the `verification_status_code` to understand how to acknowledge the certificate's validity. If the verification fails, *Proxy* returns an HTTP **418 I'm a teapot** status code, indicating to a secure rollup that it should disregard the certificate and treat it as an empty batch in its derivation pipeline.

5. Assuming a valid certificate, *Proxy* attempts to query EigenDA [retrieval paths](#retrieval-paths) for the underlying blob contents.

6. Once fetched, *Proxy* verifies the blob's KZG commitments to ensure tamper resistance (i.e., confirming that what's returned from EigenDA matches what was committed to during dispersal).

7. *Proxy* decodes the underlying blob into a `payload` type, which is returned to the *Rollup Node*.

### Retrieval Paths
There are two main blob retrieval paths:

1. **decentralized retrieval:** retrieve erasure coded chunks from Validators and recreate the `blob` from them.
2. **centralized retrieval:** the same [Relay API](https://docs.eigenda.xyz/releases/v2#relay-interfaces) that Validators use to download chunks, can also be used to retrieve full blobs.

EigenDA V2 has a new [Relay API](https://docs.eigenda.xyz/releases/v2#relay-interfaces) for retrieving blobs from the disperser. The `GetBlob` method takes a `blob_key` as input, which is a synonym for `blob_header_hash`. Note that `BlobCertificate` (different from `DACert`!) contains an array of `relay_keys`, which are the relays that can serve that specific blob. A relay’s URL can be retrieved from the [relayKeyToUrl](https://github.com/Layr-Labs/eigenda/blob/9a4bdc099b98f6e5116b11778f0cf1466f13779c/contracts/src/core/EigenDARelayRegistry.sol#L35) function on the EigenDARelayRegistry.sol contract.

## Decoding

Decoding performs the exact reverse operations that [Encoding](#encoding) did.