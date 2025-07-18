// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import {IRegistryCoordinator, RegistryCoordinator} from "lib/eigenlayer-middleware/src/RegistryCoordinator.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";
import {IPauserRegistry} from
    "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";
import {IEjectionManager} from "lib/eigenlayer-middleware/src/interfaces/IEjectionManager.sol";
import "forge-std/StdToml.sol";
import {EigenDATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

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

    function churnApprover(string memory configData) internal pure returns (address) {
        return stdToml.readAddress(configData, ".initParams.middleware.registryCoordinator.churnApprover");
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

    function minimumStakes(string memory configData) internal pure returns (uint96[] memory res) {
        uint256[] memory minimumStakesRaw =
            stdToml.readUintArray(configData, ".initParams.middleware.registryCoordinator.minimumStakes");
        res = new uint96[](minimumStakesRaw.length);
        for (uint256 i = 0; i < minimumStakesRaw.length; i++) {
            res[i] = uint96(minimumStakesRaw[i]);
        }
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

    function rewardsInitiator(string memory configData) internal pure returns (address) {
        return stdToml.readAddress(configData, ".initParams.eigenDA.serviceManager.rewardsInitiator");
    }

    function ejectors(string memory configData) internal pure returns (address[] memory) {
        return stdToml.readAddressArray(configData, ".initParams.middleware.ejectionManager.ejectors");
    }

    function quorumEjectionParams(string memory configData)
        internal
        pure
        returns (IEjectionManager.QuorumEjectionParams[] memory)
    {
        bytes memory quorumEjectionParamsRaw =
            stdToml.parseRaw(configData, ".initParams.middleware.ejectionManager.quorumEjectionParams");
        return abi.decode(quorumEjectionParamsRaw, (IEjectionManager.QuorumEjectionParams[]));
    }

    function minNumSymbols(string memory configData) internal pure returns (uint64) {
        return uint64(stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.minNumSymbols"));
    }

    function pricePerSymbol(string memory configData) internal pure returns (uint64) {
        return uint64(stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.pricePerSymbol"));
    }

    function priceUpdateCooldown(string memory configData) internal pure returns (uint64) {
        return uint64(stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.priceUpdateCooldown"));
    }

    function globalSymbolsPerPeriod(string memory configData) internal pure returns (uint64) {
        return uint64(stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.globalSymbolsPerPeriod"));
    }

    function reservationPeriodInterval(string memory configData) internal pure returns (uint64) {
        return uint64(stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.reservationPeriodInterval"));
    }

    function globalRatePeriodInterval(string memory configData) internal pure returns (uint64) {
        return uint64(stdToml.readUint(configData, ".initParams.eigenDA.paymentVault.globalRatePeriodInterval"));
    }

    function certVerifierSecurityThresholds(string memory configData)
        internal
        pure
        returns (EigenDATypesV1.SecurityThresholds memory thresholds)
    {
        thresholds.confirmationThreshold =
            uint8(stdToml.readUint(configData, ".initParams.eigenDA.certVerifier.confirmationThreshold"));
        thresholds.adversaryThreshold =
            uint8(stdToml.readUint(configData, ".initParams.eigenDA.certVerifier.adversaryThreshold"));
    }

    function certVerifierQuorumNumbersRequired(string memory configData) internal pure returns (bytes memory) {
        uint256[] memory certQuorumNumbersRequired =
            stdToml.readUintArray(configData, ".initParams.eigenDA.certVerifier.quorumNumbersRequired");

        // encode each quorum number as a single byte
        bytes memory quorumNumbersRequiredBytes = new bytes(certQuorumNumbersRequired.length);
        for (uint256 i = 0; i < certQuorumNumbersRequired.length; i++) {
            quorumNumbersRequiredBytes[i] = bytes1(uint8(certQuorumNumbersRequired[i]));
        }
        return quorumNumbersRequiredBytes;
    }
}
