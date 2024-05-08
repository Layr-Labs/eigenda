// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {ExistingDeploymentParser, PaymentCoordinator} from "eigenlayer-scripts/utils/ExistingDeploymentParser.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";

import {EigenDAServiceManager} from "../../../src/core/EigenDAServiceManager.sol";

/**
 * @title ServiceManagerBaseUpgrade for Preprod contracts. 
 * Assumes EOA deploying has permissions to call the proxyAdmin to upgrade.
 *
 *
 * Local Fork: Deploy/Upgrade PaymentCoordinator
 * anvil --fork-url $RPC_HOLESKY
 * forge script script/deploy/holesky/EigenDASM_PaymentsUpgrade.s.sol:ServiceManagerBaseUpgrade --rpc-url http://127.0.0.1:8545 --private-key $PRIVATE_KEY --broadcast -vvvv --verify
 *
 * Upgrade Holesky testnet: Deploy/Upgrade PaymentCoordinator
 * forge script script/deploy/holesky/EigenDASM_PaymentsUpgrade.s.sol:ServiceManagerBaseUpgrade --rpc-url $RPC_HOLESKY --private-key $PRIVATE_KEY --broadcast -vvvv --verify
 */
contract ServiceManagerBaseUpgrade is ExistingDeploymentParser {
    // Hardcode these values to your needs
    address public serviceManager = 0x54A03db2784E3D0aCC08344D05385d0b62d4F432;
    address public serviceManagerImplementation = 0xFe779fB43280A92cd85466312E2AE8A4F1A48007;
    ProxyAdmin public avsProxyAdmin = ProxyAdmin(0x9Fd7E279f5bD692Dc04792151E14Ad814FC60eC1);
    address deployerAddress = 0xDA29BB71669f46F2a779b4b62f03644A84eE3479;
    address registryCoordinator = 0x2c61EA360D6500b58E7f481541A36B443Bc858c6;
    address stakeRegistry = 0x53668EBf2e28180e38B122c641BC51Ca81088871;

    function run() external {
        // 1. Setup and parse existing EigenLayer Holesky preprod contracts
        _parseInitialDeploymentParams(
            "script/deploy/holesky/config/eigenlayer_preprod.config.json"
        );
        _parseDeployedContracts(
            "script/deploy/holesky/config/eigenlayer_preprod_addresses.config.json"
        );

        // 2. broadcast deployment
        vm.startBroadcast();

        emit log_named_address("Deployer Address", msg.sender);
        _upgradeServiceManager();

        vm.stopBroadcast();


        // 3. Sanity Checks
        _verifyUpgrade();

        // Verify Eigenlayer contracts parsed from config
        _verifyContractPointers();
        _verifyImplementations();
        _verifyContractsInitialized({isInitialDeployment: false});
        _verifyInitializationParams();
    }

    /// @dev Should override this to change to your specific upgrade needs
    function _upgradeServiceManager() internal virtual {
        // 1. Deploy new ServiceManager implementation contract
        serviceManagerImplementation = address(
            new EigenDAServiceManager(
                avsDirectory,
                paymentCoordinator,
                IRegistryCoordinator(registryCoordinator),
                IStakeRegistry(stakeRegistry)
            )
        );

        // 2. Upgrade the ServiceManager proxy to the new implementation
        avsProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(serviceManager))),
            address(serviceManagerImplementation)
        );
    }

    /// @dev 
    function _verifyUpgrade() internal virtual {
        // Preprod PaymentCoordinator
        require(
            address(paymentCoordinator) == 0xb22Ef643e1E067c994019A4C19e403253C05c2B0,
            "ServiceManagerBaseUpgrade: PaymentCoordinator address is incorrect"
        );
        require(
            avsProxyAdmin.getProxyImplementation(
                TransparentUpgradeableProxy(payable(serviceManager))
            ) == serviceManagerImplementation,
            "ServiceManagerBaseUpgrade: ServiceMananger implementation initially set incorrectly"
        );
        require(
            msg.sender == deployerAddress,
            "ServiceManagerBaseUpgrade: deployer address is incorrect"
        );
    }
}
