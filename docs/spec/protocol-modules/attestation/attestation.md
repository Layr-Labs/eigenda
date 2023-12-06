# Attestation Module

The attestation module of the EigenDA protocol is implemented by the EigenDA smart contracts which provide the following interfaces:
- A `confirmBatch` method for confirming data store confirmation.
- `register` and `deregister` methods for registering and deregistering DA nodes.

**Confirmation behavior**
The `confirmBatch` interface upholds the following system properties:
1. Sufficient stake checking: A blob is only accepted on-chain when signatures from operators having sufficient stake on each quorum are presented.
2. Reorg behavior: On chain confirmations behave properly under chain reorgs or forks.
3. Confirmer permissioning: Only permissioned dispersers can confirm blobs.

**Operator registration guards**
The `register` and `deregister` interfaces uphold the following system properties:
1. DA nodes cannot register if their delegated stake is insufficient.
2. DA nodes cannot deregister if they are still responsible for storing data.

This document discusses how these properties are achieved by the attestation protocol.

## Confirmation Behavior

### Sufficient stake checking

The [BLSRegistry.sol](../contracts-registry.md) maintains the `pubkeyToStakeHistory` and `pubKeyToIndexHistory` storage variables, which allow for the current stake and index of each operator to be retrieved for an arbitrary block number. These variables are updated whenever DA nodes register or deregister.

TODO: Describe quorum storage variables.

Whenever the `confirmBatch` method of the [ServiceMananger.sol](../contracts-service-manager.md) is called, the following checks are used to ensure that sufficient stake is held by the set of signatories.
- Signature verification. The signature of the calculated aggregate public key is verified.
- Quorum threshold check. The total stake of the signing DA nodes is verified to be above the quorum thresholds for each quorum.

### Reorg behavior

One aspect of the chain behavior of which the attestation protocol must be aware is that of chain reorganization. The following requirements relate to chain reorganizations:
1. Signed attestations should remain valid under reorgs so that a disperser never needs to resend the data and gather new signatures.
2. If an attestation is reorged out, a disperser should always be able to simply resubmit it after a specific waiting period.
3. Payloads constructed by a disperser and sent to DA nodes should never be rejected due to reorgs.

These requirements result in the following design choices:
- Chunk allotments should be based on registration state from a finalized block.
- If an attestation is reorged out and if the transaction containing the header of a batch is not present within `BLOCK_STALE_MEASURE` blocks since `referenceBlockNumber` and the block that is `BLOCK_STALE_MEASURE` blocks since `referenceBlockNumber`  is finalized, then the disperser should again start a new dispersal with that blob of data. Otherwise, the disperser must not re-submit another transaction containing the header of a batch associated with the same blob of data.
- Payment payloads sent to DA nodes should only take into account finalized attestations.

The first and second decision satisfies requirements 1 and 2. The three decisions together satisfy requirement 3.

Whenever the `confirmBatch` method of the [ServiceMananger.sol](../contracts-service-manager.md) is called, the following checks are used to ensure that only finalized registration state is utilized:
- Stake staleness check. The `referenceBlockNumber` is verified to be within `BLOCK_STALE_MEASURE` blocks before the confirmation block.This is to make sure that batches using outdated stakes are not confirmed. It is assured that stakes from within `BLOCK_STALE_MEASURE` blocks before confirmation are valid by delaying removal of stake by `BLOCK_STALE_MEASURE + MAX_DURATION_BLOCKS`.

### Confirmer Permissioning

TODO: Specify how confirmer is permissioned


## Operator registration guards

TODO: Describe these guards
