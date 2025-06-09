// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {IRegistryCoordinator, RegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {IPauserRegistry} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";
import "forge-std/StdToml.sol";

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
    DATypesV1.VersionedBlobParams[] versionedBlobParams;
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
    ImmutableDisperserRegistryParams disperserRegistryParams;
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

struct ImmutableDisperserRegistryParams {
    uint256 deposit;
    uint256 refund;
    address token;
    uint64 lockPeriod;
}

library InitParamsLib {
    function initialOwner(string memory configData) internal pure returns (address) {
        return stdToml.readAddress(configData, ".initialOwner");
    }

    function pausers(string memory configData) internal pure returns (address[] memory) {
        return stdToml.readAddressArray(configData, ".initParams.core.pauserRegistry.pausers");
    }

    function unpauser(string memory configData) internal pure returns (address) {
        return stdToml.readAddress(configData, ".initParams.core.pauserRegistry.unpauser");
    }

    function rewardsCoordinator(string memory configData) internal pure returns (address) {
        return stdToml.readAddress(configData, ".initParams.shared.rewardsCoordinator");
    }

    function avsDirectory(string memory configData) internal pure returns (address) {
        return stdToml.readAddress(configData, ".initParams.shared.avsDirectory");
    }

    function delegationManager(string memory configData) internal pure returns (address) {
        return stdToml.readAddress(configData, ".initParams.shared.delegationManager");
    }

    function initialPausedStatus(string memory configData) internal pure returns (uint256) {
        return stdToml.readUint(configData, ".initParams.shared.initialPausedStatus");
    }

    function registryCoordinatorParams(string memory configData)
        internal
        pure
        returns (ImmutableRegistryCoordinatorParams memory)
    {
        return ImmutableRegistryCoordinatorParams({
            churnApprover: stdToml.readAddress(configData, ".initParams.middleware.registryCoordinator.churnApprover"),
            ejector: stdToml.readAddress(configData, ".initParams.middleware.registryCoordinator.ejector")
        });
    }

    function paymentVaultParams(string memory configData) internal pure returns (ImmutablePaymentVaultParams memory) {
        return ImmutablePaymentVaultParams({
            minNumSymbols: uint64(stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.minNumSymbols")),
            pricePerSymbol: uint64(stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.pricePerSymbol")),
            priceUpdateCooldown: uint64(
                stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.priceUpdateCooldown")
            ),
            globalSymbolsPerPeriod: uint64(
                stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.globalSymbolsPerPeriod")
            ),
            reservationPeriodInterval: uint64(
                stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.reservationPeriodInterval")
            ),
            globalRatePeriodInterval: uint64(
                stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.globalRatePeriodInterval")
            )
        });
    }

    function serviceManagerParams(string memory configData)
        internal
        pure
        returns (ImmutableServiceManagerParams memory)
    {
        return ImmutableServiceManagerParams({
            rewardsInitiator: stdToml.readAddress(configData, ".initParams.eigenDA.serviceManager.rewardsInitiator")
        });
    }

    function disperserRegistryParams(string memory configData)
        internal
        pure
        returns (ImmutableDisperserRegistryParams memory)
    {
        return ImmutableDisperserRegistryParams({
            deposit: stdToml.readUint(configData, ".initParams.eigenDA.disperserRegistry.deposit"),
            refund: stdToml.readUint(configData, ".initParams.eigenDA.disperserRegistry.refund"),
            token: stdToml.readAddress(configData, ".initParams.eigenDA.disperserRegistry.token"),
            lockPeriod: uint64(stdToml.readUint(configData, ".initParams.eigenDA.disperserRegistry.lockPeriod"))
        });
    }

    function operatorSetParams(string memory configData)
        internal
        pure
        returns (IRegistryCoordinator.OperatorSetParam[] memory)
    {
        bytes memory operatorConfigsRaw =
            stdToml.parseRaw(configData, ".initParams.middleware.registryCoordinator.operatorSetParams");
        return abi.decode(operatorConfigsRaw, (IRegistryCoordinator.OperatorSetParam[]));
    }

    function minimumStakes(string memory configData) internal pure returns (uint96[] memory) {
        bytes memory stakesConfigsRaw =
            stdToml.parseRaw(configData, ".initParams.middleware.registryCoordinator.minimumStakes");
        return abi.decode(stakesConfigsRaw, (uint96[]));
    }

    function strategyParams(string memory configData)
        internal
        pure
        returns (IStakeRegistry.StrategyParams[][] memory)
    {
        bytes memory strategyConfigsRaw =
            stdToml.parseRaw(configData, ".initParams.middleware.registryCoordinator.strategyParams");
        return abi.decode(strategyConfigsRaw, (IStakeRegistry.StrategyParams[][]));
    }

    function quorumAdversaryThresholdPercentages(string memory configData) internal pure returns (bytes memory) {
        return
            stdToml.readBytes(configData, ".initParams.eigenDA.thresholdRegistry.quorumAdversaryThresholdPercentages");
    }

    function quorumConfirmationThresholdPercentages(string memory configData) internal pure returns (bytes memory) {
        return stdToml.readBytes(
            configData, ".initParams.eigenDA.thresholdRegistry.quorumConfirmationThresholdPercentages"
        );
    }

    function quorumNumbersRequired(string memory configData) internal pure returns (bytes memory) {
        return stdToml.readBytes(configData, ".initParams.eigenDA.thresholdRegistry.quorumNumbersRequired");
    }

    function versionedBlobParams(string memory configData)
        internal
        pure
        returns (DATypesV1.VersionedBlobParams[] memory)
    {
        bytes memory versionedBlobParamsRaw =
            stdToml.parseRaw(configData, ".initParams.eigenDA.thresholdRegistry.versionedBlobParams");
        return abi.decode(versionedBlobParamsRaw, (DATypesV1.VersionedBlobParams[]));
    }

    function batchConfirmers(string memory configData) internal pure returns (address[] memory) {
        return stdToml.readAddressArray(configData, ".initParams.eigenDA.serviceManager.batchConfirmers");
    }

    function calldataInitParams(string memory configData) internal pure returns (CalldataInitParams memory) {
        return CalldataInitParams({
            registryCoordinatorParams: CalldataRegistryCoordinatorParams({
                operatorSetParams: operatorSetParams(configData),
                minimumStakes: minimumStakes(configData),
                strategyParams: strategyParams(configData)
            }),
            thresholdRegistryParams: CalldataThresholdRegistryParams({
                quorumAdversaryThresholdPercentages: quorumAdversaryThresholdPercentages(configData),
                quorumConfirmationThresholdPercentages: quorumConfirmationThresholdPercentages(configData),
                quorumNumbersRequired: quorumNumbersRequired(configData),
                versionedBlobParams: versionedBlobParams(configData)
            }),
            serviceManagerParams: CalldataServiceManagerParams({batchConfirmers: batchConfirmers(configData)})
        });
    }
}
