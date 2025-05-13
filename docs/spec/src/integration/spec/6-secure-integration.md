# Secure Integration

This page is meant to be read by eigenda and rollup developers who are writing a secure integration and need to understand the details. For users who just want a high-level understanding of what a secure integration is, please visit our [secure integration overview](https://docs.eigenda.xyz/integrations-guides/rollup-guides/integrations-overview) page instead.

## Validity Conditions

EigenDA is a service that assures the availability and integrity of payloads posted to it for 14 days.
When deriving a rollup chain by running its derivation pipeline, only EigenDA certs that satisfy 3 validity conditions are considered valid and used:
1. RBN Recency Validation - ensure that the cert's reference block number (RBN) is not too old with respect to the L1 block at which the cert was included in the rollup's batcher-inbox. This ensures that the blob on EigenDA has sufficient availability time left (out of the 14 day period) in order to be downloadable if needed during a rollup fault proof window.
2. Cert Validation - ensures sufficient operator stake has signed to make the blob available, for all specified quorums. The stake is obtained onchain at a given reference block number (RBN) specified inside the cert.
3. Blob Validation - ensures that the blob used is consistent with the KZG commitment inside the Cert.

### 1. RBN Recency Validation

This check is related to time guarantees. Timing is especially important for optimistic rollups where 
the blobs need to be available during an interactive fault proof window in order to 
protect the rollup's safety.

![](../../assets/integration/recency-window-timeline.png)

Looking at the timing diagram above, in order for the blobs to be available during the entire challenge period, we need the fault proof game to start < 7days after the cert is posted to the rollup's batcher inbox. In order to uphold this guarantee, what we need to do is simply to have rollups' derivation pipelines reject certs whose DA availability period started a long time ago. However, from the cert itself, there is no way to know when the cert was signed and made available. The only information available on the cert itself is cert.RBN, the reference block number chosen by the disperser at which to anchor operator stakes. But that happens to be before validators sign, so it is enough to bound how far that can be from the cert's inclusion block.

Rollups must thus enforce that
```
cert.L1InclusionBlock - cert.RBN < RecencyWindowSize
```

This has a second security implication. A malicious EigenDA disperser could have chosen a reference block number (RBN) that is very old, where the stake of operators was very different from the current one, due to operators withdrawing stake for example.

> To give a concrete example with a rollup stack, optimism has a [sequencerWindow](https://docs.optimism.io/stack/rollup/derivation-pipeline#sequencer-window) which forces batches to land onchain in a timely fashion (12h). This filtering however, happens in the [BatchQueue](https://specs.optimism.io/protocol/derivation.html#batch-queue) stage of the derivation pipeline (DP), and doesn't prevent the DP being stalled in the [L1Retrieval](https://specs.optimism.io/protocol/derivation.html#l1-retrieval) stage by an old cert having been submitted whose blob is no longer available on EigenDA. To prevent this, we need the recencyWindow filtering to happen during the L1Retrieval stage of the DP.
>
> Despite its semantics being slightly different, sequencerWindow and recencyWindow are related concepts, and in order to not force another config change on op altda forks, we suggest using the same value as the `SequencerWindowSize` for the `RecencyWindowSize`, namely 12h.

![image.png](../../assets/integration/cert-rbn-recency-window.png)


### 2. Cert Validation

Cert validation is done inside the EigenDACertVerifier contract, which EigenDA deploys as-is, but is also available for rollups to modify and deploy on their own. Specifically, [verifyDACertV2](https://github.com/Layr-Labs/eigenda/blob/ee092f345dfbc37fce3c02f99a756ff446c5864a/contracts/src/periphery/cert/v2/EigenDACertVerifierV2.sol#L72) is the entry point for validation. This could either be called during a normal eth transaction (either for pessimistic “bridging” like EigenDA V1 used to do, or when uploading a Blob Field Element to a one-step-proof’s [preimage contract](https://specs.optimism.io/fault-proof/index.html#pre-image-oracle)), or be zk proven using a library like [Steel](https://github.com/risc0/risc0-ethereum/blob/main/crates/steel/docs/what-is-steel.md).

The [cert verification](https://github.com/Layr-Labs/eigenda/blob/ee092f345dfbc37fce3c02f99a756ff446c5864a/contracts/src/periphery/cert/v2/EigenDACertVerificationV2Lib.sol#L122) logic consists of:

1. [merkleInclusion](https://github.com/Layr-Labs/eigenda/blob/ee092f345dfbc37fce3c02f99a756ff446c5864a/contracts/src/periphery/cert/v2/EigenDACertVerificationV2Lib.sol#L132): 
2. verify `sigma` (operators’ bls signature) over `batchRoot` using the `NonSignerStakesAndSignature` struct
3. verify blob security params (blob_params + security thresholds)
4. verify each quorum part of the blob_header has met its threshold

### 3. Blob Validation

There are different required validation steps, depending on whether the client is retrieving or dispersing a blob.

Retrieval (whether data is coming from relays, or directly from DA nodes):

1. Verify that received blob length is ≤ the `length` in the cert’s `BlobCommitment`
2. Verify that the blob length claimed in the `BlobCommitment` is greater than `0`
3. Verify that the blob length claimed in the `BlobCommitment` is a power of two
4. Verify that the payload length claimed in the encoded payload header is ≤ the maximum permissible payload length, as calculated from the `length` in the `BlobCommitment`
    1. The maximum permissible payload length is computed by looking at the claimed blob length, and determining how many bytes would remain if you were to remove the encoding which is performed when converting a `payload` into an `encodedPayload`. This presents an upper bound for payload length: e.g. “If the `payload` were any bigger than `X`, then the process of converting it to an `encodedPayload` would have yielded a `blob` of larger size than claimed”
5. If the bytes received for the blob are longer than necessary to convey the payload, as determined by the claimed payload length, then verify that all extra bytes are `0x0`.
    1. Due to how padding of a blob works, it’s possible that there may be trailing `0x0` bytes, but there shouldn’t be any trailing bytes that aren’t equal to `0x0`.
6. Verify the KZG commitment. This can either be done:
    1. directly: recomputing the commitment using SRS points and checking that the two commitments match (this is the current implemented way)
    2. indirectly: verifying a point opening using Fiat-Shamir (see this [issue](https://github.com/Layr-Labs/eigenda/issues/1037))

Dispersal:

1. If the `BlobCertificate` was generated using the disperser’s `GetBlobCommitment` RPC endpoint, verify its contents:
    1. verify KZG commitment
    2. verify that `length` matches the expected value, based on the blob that was actually sent
    3. verify the `lengthProof` using the `length` and `lengthCommitment`
2. After dispersal, verify that the `BlobKey` actually dispersed by the disperser matches the locally computed `BlobKey`

Note: The verification steps in point 1. for dispersal are not currently implemented. This route only makes sense for clients that want to avoid having large amounts of SRS data, but KZG commitment verification via Fiat-Shamir is required to do the verification without this data. Until the alternate verification method is implemented, usage of `GetBlobCommitment` places a correctness trust assumption on the disperser generating the commitment.

### Rollup Stack Secure Integrations

|                     | Nitro V1       | OP V1 (insecure) | Nitro V2       | OP V2                                                                                |
| ------------------- | -------------- | ---------------- | -------------- | ------------------------------------------------------------------------------------ |
| Cert Verification   | SequencerInbox | x                | one-step proof | one-step proof: done in preimage oracle contract when uploading a blob field element |
| Blob Verification   | one-step proof | x                | one-step proof | one-step proof                                                                       |
| Timing Verification | SequencerInbox | x                | SequencerInbox | one-step proof (?)                                                                   |