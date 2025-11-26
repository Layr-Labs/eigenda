// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import "../Env.sol";
import "./1-DeployImplementations.s.sol";
import {EOADeployer} from "zeus-templates/templates/EOADeployer.sol";
import {Encode} from "zeus-templates/utils/Encode.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {
    IPauserRegistry
} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/interfaces/IPauserRegistry.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {Pausable} from "lib/eigenlayer-middleware/lib/eigenlayer-contracts/src/contracts/permissions/Pausable.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

/// @title ExecuteUpgrade
/// @notice Execute upgrade of EigenDA implementations via timelock controller
contract ExecuteUpgrade is EOADeployer {
    using Env for *;
    using Encode for *;

    function _runAsEOA() internal override {
        // Get proxy admin
        ProxyAdmin proxyAdmin = ProxyAdmin(Env.proxyAdmin());

        // Upgrade ServiceManager with reinitialization
        {
            // Get current batch confirmers - we need to reconstruct this from events or state
            // For now, using empty array as we'll need to set them post-upgrade
            address[] memory batchConfirmers = new address[](0);

            proxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(Env.proxy.serviceManager()))),
                address(Env.impl.serviceManager()),
                abi.encodeWithSelector(
                    EigenDAServiceManager.initialize.selector,
                    Env.proxy.serviceManager().pauserRegistry(),
                    Pausable(address(Env.proxy.serviceManager())).paused(),
                    Env.impl.owner(), // newOwner
                    batchConfirmers,
                    Env.proxy.serviceManager().rewardsInitiator()
                )
            );
        }

        // Upgrade RegistryCoordinator with reinitialization
        {
            // Get quorum count to read operator set params
            uint8 quorumCount = Env.proxy.registryCoordinator().quorumCount();
            IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams =
                new IRegistryCoordinator.OperatorSetParam[](quorumCount);
            uint96[] memory minimumStakes = new uint96[](quorumCount);
            IStakeRegistry.StrategyParams[][] memory strategyParams = new IStakeRegistry.StrategyParams[][](quorumCount);

            StakeRegistry stakeRegistry = Env.proxy.stakeRegistry();

            // Read current configuration for each quorum
            for (uint8 i = 0; i < quorumCount; i++) {
                operatorSetParams[i] = Env.proxy.registryCoordinator().getOperatorSetParams(i);
                minimumStakes[i] = stakeRegistry.minimumStakeForQuorum(i);

                // Read strategy params for this quorum
                uint256 strategyParamsLength = stakeRegistry.strategyParamsLength(i);
                strategyParams[i] = new IStakeRegistry.StrategyParams[](strategyParamsLength);
                for (uint256 j = 0; j < strategyParamsLength; j++) {
                    strategyParams[i][j] = stakeRegistry.strategyParamsByIndex(i, j);
                }
            }

            proxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(Env.proxy.registryCoordinator()))),
                address(Env.impl.registryCoordinator()),
                abi.encodeWithSelector(
                    EigenDARegistryCoordinator.initialize.selector,
                    Env.impl.owner(), // newOwner
                    Env.proxy.registryCoordinator().ejector(),
                    Env.proxy.registryCoordinator().pauserRegistry(),
                    Pausable(address(Env.proxy.registryCoordinator())).paused(),
                    operatorSetParams,
                    minimumStakes,
                    strategyParams
                )
            );
        }

        // Upgrade ThresholdRegistry with reinitialization
        {
            // Read current threshold parameters
            bytes memory quorumAdversaryThresholdPercentages =
                Env.proxy.thresholdRegistry().quorumAdversaryThresholdPercentages();
            bytes memory quorumConfirmationThresholdPercentages =
                Env.proxy.thresholdRegistry().quorumConfirmationThresholdPercentages();
            bytes memory quorumNumbersRequired = Env.proxy.thresholdRegistry().quorumNumbersRequired();

            // Read versioned blob params
            uint16 nextBlobVersion = Env.proxy.thresholdRegistry().nextBlobVersion();
            DATypesV1.VersionedBlobParams[] memory versionedBlobParams =
                new DATypesV1.VersionedBlobParams[](nextBlobVersion);
            for (uint16 i = 0; i < nextBlobVersion; i++) {
                versionedBlobParams[i] = Env.proxy.thresholdRegistry().getBlobParams(i);
            }

            proxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(Env.proxy.thresholdRegistry()))),
                address(Env.impl.thresholdRegistry()),
                abi.encodeWithSelector(
                    EigenDAThresholdRegistry.initialize.selector,
                    Env.impl.owner(), // newOwner
                    quorumAdversaryThresholdPercentages,
                    quorumConfirmationThresholdPercentages,
                    quorumNumbersRequired,
                    versionedBlobParams
                )
            );
        }

        // Upgrade RelayRegistry with reinitialization
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.relayRegistry()))),
            address(Env.impl.relayRegistry()),
            abi.encodeWithSelector(EigenDARelayRegistry.initialize.selector, Env.impl.owner()) // newOwner
        );

        // Upgrade DisperserRegistry with reinitialization
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.disperserRegistry()))),
            address(Env.impl.disperserRegistry()),
            abi.encodeWithSelector(EigenDADisperserRegistry.initialize.selector, Env.impl.owner()) // newOwner
        );

        // Upgrade PaymentVault with reinitialization
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.paymentVault()))),
            address(Env.impl.paymentVault()),
            abi.encodeWithSelector(
                PaymentVault.initialize.selector,
                Env.impl.owner(), // newOwner
                Env.proxy.paymentVault().minNumSymbols(),
                Env.proxy.paymentVault().pricePerSymbol(),
                Env.proxy.paymentVault().priceUpdateCooldown(),
                Env.proxy.paymentVault().globalSymbolsPerPeriod(),
                Env.proxy.paymentVault().reservationPeriodInterval(),
                Env.proxy.paymentVault().globalRatePeriodInterval()
            )
        );
    }

    /// -----------------------------------------------------------------------
    /// 2) Post-upgrade assertions
    /// -----------------------------------------------------------------------

    function testScript() public virtual {
        // Hook for pre-test setup.
        _beforeTestScript();
        // Execute upgrade.
        runAsEOA();
        // Hook for post-upgrade assertions.
        _afterTestScript();
    }

    /// -----------------------------------------------------------------------
    /// Test hooks
    /// -----------------------------------------------------------------------

    function _beforeTestScript() internal view {}

    function _afterTestScript() internal view {}
}
