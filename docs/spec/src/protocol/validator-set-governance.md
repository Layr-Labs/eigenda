# Decentralized Validator Set Governance

## Overview

EigenDA's validator set governance manages validator entry and exit in a decentralized way. This document describes the ejection and churning protocols that govern how validators leave and join the EigenDA validator set.

The protocol includes:
- Ejection: dispersers may eject under-performing validators, with validators able to cancel ejections.
- Churner: an on-chain function that removes the validator with the smallest amount of stake to allow a validator to join when the validator set is full.

## 1. Ejection Protocol

The ejection protocol maintains EigenDA's liveness and quality of service by allowing dispersers to eject honest but under-performing validators.

### 1.1 Protocol Actors

| Actor | Role | Implementation |
|-------|------|----------------|
| **Ejector** (Disperser) | Monitors validator performance and initiates ejections | [`ejector/`](https://github.com/Layr-Labs/eigenda/tree/master/ejector) |
| **Ejectee** (Validator) | Monitors ejection attempts and defends against unjust ejections | [`node/ejection/ejection_sentinel.go`](https://github.com/Layr-Labs/eigenda/blob/master/node/ejection/ejection_sentinel.go) |
| **Ejection Manager** | Smart contract coordinating ejection lifecycle | [`EigenDAEjectionManager.sol`](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/ejection/EigenDAEjectionManager.sol) |

### 1.2 Ejection Initiation

The ejection lifecycle is managed by the `BeginEjection()` method in [`ejector/ejection_manager.go:127-193`](https://github.com/Layr-Labs/eigenda/blob/master/ejector/ejection_manager.go#L127-L193), which performs all pre-flight checks before initiating an on-chain ejection.

#### 1.2.1 Ejector Authorization

Only authorized dispersers can initiate ejections. Authorized disperser addresses are stored in an allow-list within the `EigenDAEjectionManager` contract. Initially, this list contains only the EigenDA disperser operated by EigenLabs. The list can be expanded as additional dispersers become available.

**Implementation**: The contract enforces this via the `onlyEjector` modifier ([`EigenDAEjectionManager.sol:66-69`](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/ejection/EigenDAEjectionManager.sol#L66-L69)), which checks the `EJECTOR_ROLE` using AccessControl.

#### 1.2.2 Automatic Ejection Decision-Making

The disperser monitors validator performance over a configurable time window (`performance_evaluation_window`, default: 10 minutes) and computes each validator's `signing_rate`.

A validator becomes eligible for ejection only when **all** of the following conditions are met:

1. **Zero signing rate**: The validator's `signing_rate` is zero over the evaluation window
2. **Cool-down period elapsed**: `DISPERSER_COOL_DOWN` has passed since the last ejection attempt against this validator
3. **Selective non-participation**: Other validators show non-zero signing rates during the same period

These rules prevent ejections during network-wide outages and limit wasted transaction fees when dealing with potentially malicious validators who repeatedly cancel ejections while being under-performing.

**Implementation**: The evaluation logic is in [`ejector/ejector.go:102-184`](https://github.com/Layr-Labs/eigenda/blob/master/ejector/ejector.go#L102-L184). The ejection criterion is implemented as:

```go
// ejector/ejector.go:146
isEjectable := signingRate.GetSignedBatches() == 0 && signingRate.GetUnsignedBatches() > 0
```

This ensures a validator is only ejectable if they signed zero batches but there were batches to sign (selective non-participation). The evaluation window is configured via `EjectionCriteriaTimeWindow` in [`ejector/ejector_config.go:41-45`](https://github.com/Layr-Labs/eigenda/blob/master/ejector/ejector_config.go#L41-L45).

#### 1.2.3 Non-Ejection List

The disperser maintains a non-ejection list to handle validators that repeatedly cancel ejections without actually performing their duties. When a validator's failed ejection attempts reach `MAX_FAILURE_TIMES`, they are added to this list and **automatic ejection stops**. Human intervention is then required to deal with these validators. This list can also be manually configured.

**Implementation**: The non-ejection list (called `ejectionBlacklist`) is maintained in [`ejector/ejection_manager.go:54-60`](https://github.com/Layr-Labs/eigenda/blob/master/ejector/ejection_manager.go#L54-L60). Failed attempts are tracked in the `failedEjectionAttempts` map (lines 68-72), and validators are added to the blacklist in `handleAbortedEjection` (lines 384-412). The threshold is configured via `MaxConsecutiveFailedEjectionAttempts` in [`ejector/ejector_config.go:53-54`](https://github.com/Layr-Labs/eigenda/blob/master/ejector/ejector_config.go#L53-L54), with a default value of 5.

#### 1.2.4 Manual Ejection

In addition to automatic ejection based on performance monitoring, dispersers can manually initiate ejections against specific validators.

### 1.3 Ejection Logic in the Smart Contract

The `EigenDAEjectionManager` contract enforces the following constraints before accepting an ejection request:

1. **Rate Limiting**: At least `EJECTION_COOL_DOWN` (30 minutes) must have passed since the previous ejection attempt against the same validator
2. **Concurrency Control**: At most one active ejection is allowed per validator at any given time

Upon accepting a valid ejection request, the contract:
1. Records the ejection in contract storage
2. Starts a cancellation window of duration `RESPONSE_TIME` (30 minutes)
3. Emits an ejection event that validators monitor

**Implementation**: The constraint checks are enforced in [`EigenDAEjectionLib.sol`](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/ejection/libraries/EigenDAEjectionLib.sol):

```solidity
// EigenDAEjectionLib.sol:36-42
require(ejectee.record.proceedingTime == 0, "Ejection already in progress");
require(ejectee.lastProceedingInitiated + s().cooldown <= block.timestamp, "Ejection cooldown not met");

ejectee.record.ejector = ejector;
ejectee.record.quorums = quorums;
ejectee.record.proceedingTime = uint64(block.timestamp) + s().delay;
ejectee.lastProceedingInitiated = uint64(block.timestamp);
```

The first `require` enforces concurrency control (one ejection per validator), the second enforces the cooldown period, and the `delay` parameter sets the cancellation window duration.

### 1.4 Validator Defense (Cancellation)

#### 1.4.1 Ejection Monitoring

Each validator node runs an ejection sentinel ([`node/ejection/ejection_sentinel.go`](https://github.com/Layr-Labs/eigenda/blob/master/node/ejection/ejection_sentinel.go)) that continuously monitors the `EigenDAEjectionManager` contract for ejection events targeting that validator.

#### 1.4.2 Cancellation Modes

Validators operate in one of two modes, configurable via a trusted dispersers list:

| Mode | Condition | Behavior |
|------|-----------|----------|
| **Mode 1** | Ejector is in trusted dispersers list | No cancellation sent (validator trusts ejector's judgment) |
| **Mode 2** | Ejector is not in trusted dispersers list | Cancel if validator is online and running compliant software version |

**Default Configuration**: The trusted dispersers list is empty by default, meaning validators operate in Mode 2 for all ejectors.

**Note**: Validators must configure a wallet to submit cancellation transactions. Until most validators have set up their cancellation infrastructure, only the EigenDA disperser will be authorized as a valid ejector.


#### 1.4.3 Cancellation Process

To cancel an ejection, the validator:

1. **Generates cancellation message** containing:
   - Chain ID (identifying which L1 blockchain)
   - Validator's address
   - Block height at which the ejection was initiated

2. **Signs the message** using the validator's BLS private key

3. **Submits transaction** to `EigenDAEjectionManager` containing the signed cancellation message

If the cancellation is received within the `RESPONSE_TIME` window and the signature is valid, the ejection is canceled and the validator remains in the validator set.

### 1.5 Ejection Finalization

If no valid cancellation is received before the `RESPONSE_TIME` window expires, any disperser can finalize the ejection by submitting a finalizing transaction to the contract. Upon finalization, the validator is deregistered from the EigenDA validator set via a call to [`EigenDARegistryCoordinator`](../contracts.md#eigendaregistrycoordinator).


### 1.7 Rejoining After Ejection

Validators that have been ejected are subject to a cool-down period of **1 day** before they can rejoin the validator set.

### 1.8 Protocol Parameters

| Parameter | Value | Description | Implementation |
|-----------|-------|-------------|----------------|
| `RESPONSE_TIME` | 30 minutes | Cancellation window duration | `delay` in [`EigenDAEjectionStorage.sol:40-42`](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/ejection/libraries/EigenDAEjectionStorage.sol#L40-L42) |
| `EJECTION_COOL_DOWN` | 30 minutes | Minimum time between ejection attempts for same validator | `cooldown` in [`EigenDAEjectionStorage.sol:40-42`](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/ejection/libraries/EigenDAEjectionStorage.sol#L40-L42) |
| `DISPERSER_COOL_DOWN` | 24 hours (default) | Cool-down before retrying ejection after failed attempt | `EjectionRetryDelay` in [`ejector/ejector_config.go:50-51`](https://github.com/Layr-Labs/eigenda/blob/master/ejector/ejector_config.go#L50-L51) |
| `MAX_FAILURE_TIMES` | 5 (default) | Failed ejection attempts before adding to non-ejection list | `MaxConsecutiveFailedEjectionAttempts` in [`ejector/ejector_config.go:53-54`](https://github.com/Layr-Labs/eigenda/blob/master/ejector/ejector_config.go#L53-L54) |
| `performance_evaluation_window` | 10 minutes (default) | Time window for computing signing rate | `EjectionCriteriaTimeWindow` in [`ejector/ejector_config.go:41-45`](https://github.com/Layr-Labs/eigenda/blob/master/ejector/ejector_config.go#L41-L45) |
| Rejoin cool-down | 1 day | Wait time before ejected validator can rejoin | (Contract-level parameter) |

### 1.9 Implementation References

| Component | Path |
|-----------|------|
| Ejector service | [`ejector/`](https://github.com/Layr-Labs/eigenda/tree/master/ejector) |
| Ejection sentinel | [`node/ejection/ejection_sentinel.go`](https://github.com/Layr-Labs/eigenda/blob/master/node/ejection/ejection_sentinel.go) |
| Ejection manager contract | [`contracts/src/periphery/ejection/EigenDAEjectionManager.sol`](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/ejection/EigenDAEjectionManager.sol) |
| Ejection library | [`contracts/src/periphery/ejection/libraries/EigenDAEjectionLib.sol`](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/ejection/libraries/EigenDAEjectionLib.sol) |
| Ejection types | [`contracts/src/periphery/ejection/libraries/EigenDAEjectionTypes.sol`](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/ejection/libraries/EigenDAEjectionTypes.sol) |
| Ejection storage | [`contracts/src/periphery/ejection/libraries/EigenDAEjectionStorage.sol`](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/periphery/ejection/libraries/EigenDAEjectionStorage.sol) |

---

## 2. Churning Protocol

The churning protocol governs how new validators join the EigenDA validator set when the maximum validator capacity has been reached. The churning logic is computed entirely on-chain.

### 2.1 Overview

When the validator set is at maximum capacity, a new validator can only join by "churning out" an existing validator with the smallest stake. The smart contract automatically identifies and ejects the smallest-stake validator to make room for the higher-stake incoming validator.

### 2.2 On-Chain Churn Selection

The [`EigenDARegistryCoordinator`](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/EigenDARegistryCoordinator.sol) contract implements the churn selection logic:

1. A new validator attempts to register and the validator set is at capacity
2. The contract iterates through all current validators in the set and identifies the validator with the smallest stake
3. Automatically deregisters the smallest-stake validator
4. Registers the new validator 

**Implementation**: The main registration logic is in `registerOperator()` ([`EigenDARegistryCoordinator.sol:108-142`](https://github.com/Layr-Labs/eigenda/blob/master/contracts/src/core/EigenDARegistryCoordinator.sol#L108-L142)), which checks if the operator count exceeds `maxOperatorCount` and calls `_churnOperator()`. The `_churnOperator()` function performs an exhaustive search:

```solidity
// EigenDARegistryCoordinator.sol:157-178
function _churnOperator(uint8 quorumNumber) internal {
    bytes32[] memory operatorList = indexRegistry().getOperatorListAtBlockNumber(quorumNumber, uint32(block.number));
    require(operatorList.length > 0, "RegCoord._churnOperator: no operators to churn");

    // Find the operator with the lowest stake
    bytes32 operatorToChurn;
    uint96 lowestStake = type(uint96).max;
    for (uint256 i; i < operatorList.length; i++) {
        uint96 operatorStake = stakeRegistry().getCurrentStake(operatorList[i], quorumNumber);
        if (operatorStake < lowestStake) {
            lowestStake = operatorStake;
            operatorToChurn = operatorList[i];
        }
    }

    // Deregister the operator with the lowest stake
    bytes memory quorumNumbers = new bytes(1);
    quorumNumbers[0] = bytes1(uint8(quorumNumber));
    _deregisterOperator({operator: blsApkRegistry().pubkeyHashToOperator(operatorToChurn), quorumNumbers: quorumNumbers});
}
```

This iterates through all operators to find the one with minimum stake and automatically deregisters them.
