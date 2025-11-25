// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import "../Env.sol";
import {DeployImplementations} from "./1-DeployImplementations.s.sol";
import {MultisigBuilder} from "zeus-templates/templates/MultisigBuilder.sol";
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
contract ExecuteUpgrade is MultisigBuilder, DeployImplementations {
    using Env for *;
    using Encode for *;

    function _runAsMultisig() internal override prank(Env.executorMultisig()) {
        // Get proxy admin
        ProxyAdmin proxyAdmin = ProxyAdmin(Env.proxyAdmin());

        // Get the new owner address (this should be set to the target owner for the migration)
        address newOwner = Env.executorMultisig(); // or use a specific owner from env

        // Upgrade ServiceManager with reinitialization
        {
            // Get current batch confirmers - we need to reconstruct this from events or state
            // For now, using empty array as we'll need to set them post-upgrade
            address[] memory batchConfirmers = new address[](0);

            proxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(Env.proxy.serviceManager()))),
                address(Env.impl.serviceManager()),
                abi.encodeWithSignature(
                    "initialize(address,uint256,address,address[],address)",
                    Env.proxy.serviceManager().pauserRegistry(),
                    Pausable(address(Env.proxy.serviceManager())).paused(),
                    newOwner,
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
                abi.encodeWithSignature(
                    "initialize(address,address,address,uint256,(uint32,uint16,uint16)[],uint96[],(address,uint96)[][])",
                    newOwner,
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
                abi.encodeWithSignature(
                    "initialize(address,bytes,bytes,bytes,(uint32,uint32,uint8)[])",
                    newOwner,
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
            abi.encodeWithSignature("initialize(address)", newOwner)
        );

        // Upgrade DisperserRegistry with reinitialization
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.disperserRegistry()))),
            address(Env.impl.disperserRegistry()),
            abi.encodeWithSignature("initialize(address)", newOwner)
        );

        // Upgrade PaymentVault with reinitialization
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(Env.proxy.paymentVault()))),
            address(Env.impl.paymentVault()),
            abi.encodeWithSignature(
                "initialize(address,uint64,uint64,uint64,uint64,uint64,uint64)",
                newOwner,
                Env.proxy.paymentVault().minNumSymbols(),
                Env.proxy.paymentVault().pricePerSymbol(),
                Env.proxy.paymentVault().priceUpdateCooldown(),
                Env.proxy.paymentVault().globalSymbolsPerPeriod(),
                Env.proxy.paymentVault().reservationPeriodInterval(),
                Env.proxy.paymentVault().globalRatePeriodInterval()
            )
        );
    }

    function testScript() public virtual override {
        // 1 - Deploy new implementations
        runAsEOA();

        // 2 - Execute upgrades via multisig
        execute();

        // 3 - Validate all upgrades were successful
        _validateUpgrades();
    }

    /// @dev Validate that all upgrades were successful
    function _validateUpgrades() internal view {
        address newOwner = Env.executorMultisig();

        // Validate ServiceManager upgrade
        address serviceManagerImpl = Env._getProxyImpl(address(Env.proxy.serviceManager()));
        assertTrue(
            serviceManagerImpl == address(Env.impl.serviceManager()),
            "ServiceManager proxy should point to new implementation"
        );
        assertTrue(Env.proxy.serviceManager().owner() == newOwner, "ServiceManager owner should be updated");

        // Validate RegistryCoordinator upgrade
        address registryCoordinatorImpl = Env._getProxyImpl(address(Env.proxy.registryCoordinator()));
        assertTrue(
            registryCoordinatorImpl == address(Env.impl.registryCoordinator()),
            "RegistryCoordinator proxy should point to new implementation"
        );
        assertTrue(Env.proxy.registryCoordinator().owner() == newOwner, "RegistryCoordinator owner should be updated");

        // Validate ThresholdRegistry upgrade
        address thresholdRegistryImpl = Env._getProxyImpl(address(Env.proxy.thresholdRegistry()));
        assertTrue(
            thresholdRegistryImpl == address(Env.impl.thresholdRegistry()),
            "ThresholdRegistry proxy should point to new implementation"
        );
        assertTrue(Env.proxy.thresholdRegistry().owner() == newOwner, "ThresholdRegistry owner should be updated");

        // Validate RelayRegistry upgrade
        address relayRegistryImpl = Env._getProxyImpl(address(Env.proxy.relayRegistry()));
        assertTrue(
            relayRegistryImpl == address(Env.impl.relayRegistry()),
            "RelayRegistry proxy should point to new implementation"
        );
        assertTrue(Env.proxy.relayRegistry().owner() == newOwner, "RelayRegistry owner should be updated");

        // Validate DisperserRegistry upgrade
        address disperserRegistryImpl = Env._getProxyImpl(address(Env.proxy.disperserRegistry()));
        assertTrue(
            disperserRegistryImpl == address(Env.impl.disperserRegistry()),
            "DisperserRegistry proxy should point to new implementation"
        );
        assertTrue(Env.proxy.disperserRegistry().owner() == newOwner, "DisperserRegistry owner should be updated");

        // Validate PaymentVault upgrade
        address paymentVaultImpl = Env._getProxyImpl(address(Env.proxy.paymentVault()));
        assertTrue(
            paymentVaultImpl == address(Env.impl.paymentVault()),
            "PaymentVault proxy should point to new implementation"
        );
        assertTrue(Env.proxy.paymentVault().owner() == newOwner, "PaymentVault owner should be updated");

        // Validate ProxyAdmin ownership
        address proxyAdmin = Env._getProxyAdmin(address(Env.proxy.serviceManager()));
        assertTrue(proxyAdmin == Env.proxyAdmin(), "ServiceManager proxy admin should be correct");

        proxyAdmin = Env._getProxyAdmin(address(Env.proxy.registryCoordinator()));
        assertTrue(proxyAdmin == Env.proxyAdmin(), "RegistryCoordinator proxy admin should be correct");

        proxyAdmin = Env._getProxyAdmin(address(Env.proxy.thresholdRegistry()));
        assertTrue(proxyAdmin == Env.proxyAdmin(), "ThresholdRegistry proxy admin should be correct");

        proxyAdmin = Env._getProxyAdmin(address(Env.proxy.relayRegistry()));
        assertTrue(proxyAdmin == Env.proxyAdmin(), "RelayRegistry proxy admin should be correct");

        proxyAdmin = Env._getProxyAdmin(address(Env.proxy.disperserRegistry()));
        assertTrue(proxyAdmin == Env.proxyAdmin(), "DisperserRegistry proxy admin should be correct");

        proxyAdmin = Env._getProxyAdmin(address(Env.proxy.paymentVault()));
        assertTrue(proxyAdmin == Env.proxyAdmin(), "PaymentVault proxy admin should be correct");
    }
}
