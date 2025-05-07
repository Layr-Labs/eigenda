# Secure Integration

This page is meant to be read by eigenda and rollup developers who are writing a secure integration and need to understand the details. For users who just want a high-level understanding of what a secure integration is, please visit our [secure integration overview](https://docs.eigenda.xyz/integrations-guides/rollup-guides/integrations-overview) page instead.

## Validity Conditions

EigenDA is a service that assures the availability and integrity of payloads posted to it for 14 days.
When deriving a rollup chain by running its derivation pipeline, only EigenDA certs that satisfy 3 validity conditions are considered valid and used:
1. RBN Recency Validation - ensure that the cert's reference block number (RBN) is not too old with respect to the L1 block at which the cert was included in the rollup's batcher-inbox. This ensures that the blob on EigenDA has sufficient availability time left (out of the 14 day period) in order to be downloadable if needed during a rollup fraud proof window.
2. Cert Validation - ensures sufficient operator stake has signed to make the blob available, for all specified quorums. The stake is obtained onchain at a given reference block number (RBN) specified inside the cert.
3. Blob Validation - ensures that the blob used is consistent with the KZG commitment inside the Cert.

### 1. RBN Recency Validation

First, we assume that every rollup stack has a mechanism to force blobs to be included in the rollup’s batcher-inbox in a timely matter, such that the cert's blob is still available on EigenDA for long enough to cover the fraud proof duration (typically 7 days), or long enough to be downloadable by zk rollup participants.

> For example, optimism has a [sequencerWindow](https://docs.optimism.io/stack/rollup/derivation-pipeline#sequencer-window) which forces batches to land onchain in a timely fashion (12h). This indirectly also forces EigenDA certs to land onchain in a timely fashion. Therefore, it is safe to assume that the certs land onchain < SequencerWindowSize after having been signed.
>
> OP stack rollups thus need to make sure that `SequencerWindow` << 14 days.

This guarantee alone however is not enough, because the EigenDA disperser could have chosen a reference block number (RBN) that is very old, such that the stake of operators does not represent the true current value. Because EigenDA V2 doesn't pessimistically bridge certs onchain like EigenDA V1, enforcing RBN to not be too old is thus the responsibility of the rollup.

> Rollups must enforce that `cert.RBN` > `cert.L1InclusionBlock` - `RecencyWindowSize`

We suggest using the same value as the `SequencerWindowSize` for the `RecencyWindowSize`, namely 12h.

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