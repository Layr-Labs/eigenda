# Trustless Integration Upgrade

>Applies only to EigenDACertV4. “Trustless integration” = “secure integration”.

## Overview

This section describes a schema for deterministically upgrading an eigenda blob derivation pipeline. The eigenda blob derivation pipeline contains two components:
- onchain: cert verifier and cert verifier router
- offchain derivation: kzg verification, recency check, altda commitment parsing and other logics defined in [secure-integration](./6-secure-integration.md).

## Background

Consensus systems (L1/L2) typically upgrade logic via hardfork at block `X`:
- Before `X`, old logic executes; after `X`, new logic executes.
- Software must be backward compatible (able to execute both logics) and enforceable (disallow executing old logic after X without stalling consensus). 

### Onchain Integration Upgrade

Integrations upgrade onchain logic by adding a new [EigenDACertVerifier](./4-contracts.md#eigendacertverifier) to a [EigenDACertVerifierRouter](./4-contracts.md#eigendacertverifierrouter). Each verifier has an corresponding activationBlockNumber (ABN) within the `EigenDACertVerifierRouter`. The router uses the DACert's reference block number to determine which verifier to use by comparing against its ABN. More see section [contracts](./4-contracts.md)

This mechanism mirrors hardfork behavior: it is backward compatible and enforceable. Each `EigenDACertVerifier` also embeds a constant `offchain derivation version`, set at deployment, which governs off-chain logic.


### Offchain Integration Upgrade

EigenDA blob derivation includes substantial off-chain processing. The `offchain derivation version` (uint16) versions the entire off-chain logic. For example, the [recency window](./6-secure-integration.md#1-rbn-recency-validation) is `14400` when its `offchain derivation version = 0`; new versions may change the recency value, alter payload encoding, or introduce new new configs or validation rules.

To safely upgrade offchain logic, the L2 node’s eigenda-proxy must know when a new version becomes valid. With a new DACert type `EigenDACertV4`, this is enforced by requiring the certVerifier to check that the DACert’s offchain derivation version matches the constant value set by the contract. Once this check passes, off-chain code can safely use the `offchain derivation version` embedded in the DACert. Thus, onchain logic controls activation of offchain versioning, ensuring backward-compatible and enforceable upgrades.

### Note

Each L2 should deploy its own router. Using the EigenLabs-deployed router delegates upgrade scheduling to EigenLabs. See [contracts]((./4-contracts.md)) for router details and deployment guidance.
