// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {IRegistryCoordinator, RegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {VersionedBlobParams} from "src/interfaces/IEigenDAStructs.sol";
import {IPauserRegistry} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";

/// @dev This is the struct needed by the multisig to initialize the contracts.
struct CalldataInitParams {
    CalldataRegistryCoordinatorParams registryCoordinatorParams;
    CalldataThresholdRegistryParams thresholdRegistryParams;
    CalldataServiceManagerParams serviceManagerParams;
}

struct CalldataRegistryCoordinatorParams {
    RegistryCoordinator.OperatorSetParam[] operatorSetParams;
    uint96[] minimumStakes;
    IStakeRegistry.StrategyParams[][] strategyParams;
}

struct CalldataThresholdRegistryParams {
    bytes quorumAdversaryThresholdPercentages;
    bytes quorumConfirmationThresholdPercentages;
    bytes quorumNumbersRequired;
    VersionedBlobParams[] versionedBlobParams;
}

struct CalldataServiceManagerParams {
    address[] batchConfirmers;
}

struct ImmutableInitParams {
    ProxyAdmin proxyAdmin;
    address initialOwner;
    IPauserRegistry pauserRegistry;
    uint256 initialPausedStatus;
    DeployedAddresses proxies;
    DeployedAddresses implementations;
    ImmutableRegistryCoordinatorParams registryCoordinatorParams;
    ImmutablePaymentVaultParams paymentVaultParams;
    ImmutableServiceManagerParams serviceManagerParams;
}

struct DeployedAddresses {
    address indexRegistry;
    address stakeRegistry;
    address socketRegistry;
    address blsApkRegistry;
    address registryCoordinator;
    address thresholdRegistry;
    address relayRegistry;
    address paymentVault;
    address disperserRegistry;
    address serviceManager;
}

struct ImmutableRegistryCoordinatorParams {
    address churnApprover;
    address ejector;
}

struct ImmutablePaymentVaultParams {
    uint64 minNumSymbols;
    uint64 pricePerSymbol;
    uint64 priceUpdateCooldown;
    uint64 globalSymbolsPerPeriod;
    uint64 reservationPeriodInterval;
    uint64 globalRatePeriodInterval;
}

struct ImmutableServiceManagerParams {
    address rewardsInitiator;
}
