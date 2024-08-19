// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {ERC20PresetFixedSupply} from "@openzeppelin/contracts/token/ERC20/presets/ERC20PresetFixedSupply.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {
    ExistingDeploymentParser,
    RewardsCoordinator,
    IRewardsCoordinator,
    IPauserRegistry,
    IStrategy,
    IERC20
} from "eigenlayer-scripts/utils/ExistingDeploymentParser.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
    
import {EigenDAServiceManager} from "../../../src/core/EigenDAServiceManager.sol";

/**
 * @title ServiceManagerBaseUpgrade for Mainnet contracts addresses
 * NOTE: This script will not actually run on mainnet for EigenDA as it assumes the EOA deploying
 * has permissions to call the proxyAdmin to upgrade. This is not the case as contract upgrades are behind a
 * 10-day timelock deploy with a multisig. This is meant to provide a template for upgrading the ServiceManagerBase for AVSs
 * and to create a valid AVS rewards submission.
 *
 *
 * Local Fork: Deploy/Upgrade RewardsCoordinator
 * anvil --fork-url $RPC_MAINNET
 * forge script script/deploy/mainnet/EigenDASM_RewardsUpgrade_example.s.sol:ServiceManagerBaseUpgrade --rpc-url http://127.0.0.1:8545 -vvvv
 */
contract ServiceManagerBaseUpgrade is ExistingDeploymentParser {
    // Hardcode these values to your needs
    address public serviceManager = 0x870679E138bCdf293b7Ff14dD44b70FC97e12fc0;
    address public serviceManagerImplementation = 0xCDFFF07d5b8AcdAd13607615118a2e65030f5be1;

    ProxyAdmin public avsProxyAdmin = ProxyAdmin(0x8247EF5705d3345516286B72bFE6D690197C2E99);
    address registryCoordinator = 0x0BAAc79acD45A023E19345c352d8a7a83C4e5656;
    address stakeRegistry = 0x006124Ae7976137266feeBFb3F4D2BE4C073139D;

    // owner address to prank who has permissions to upgrade. Actual mainnet upgrades go through a 10-day timelock
    // but for purposes of a example script of upgrading and creating a AVS rewards submission with mainnet addresses,
    // we will prank this prank address in our script tests
    address deployerAddress = 0x369e6F597e22EaB55fFb173C6d9cD234BD699111;

    function run() external {
        // 1. Setup and parse existing EigenLayer Mainnet contracts
        _parseInitialDeploymentParams(
            "script/deploy/mainnet/config/mainnet.config.json"
        );
        _parseDeployedContracts(
            "script/deploy/mainnet/config/mainnet_addresses.json"
        );

        // 2. prank to set rewardsInitiator
        // Set rewardsInitiator
        vm.prank(EigenDAServiceManager(serviceManager).owner());
        EigenDAServiceManager(serviceManager).setRewardsInitiator(deployerAddress);

        // 3. test deployment upgrade and createAVSRewardsSubmission
        vm.startPrank(deployerAddress);

        emit log_named_address("Deployer Address", msg.sender);
        _upgradeServiceManager();
        _createAVSRewardsSubmission();

        vm.stopPrank(); 

        // 4. Sanity Checks
        _verifyUpgrade();

        // Verify Eigenlayer contracts parsed from config
        _verifyContractPointers();
        _verifyImplementations();
        // comment out initalize checks because mainnet initialize interfaces are different from current implementation
        // ex. EigenPodManager from pending PEPE upgrade
        // _verifyContractsInitialized();
        _verifyInitializationParams();
    }

    /// @dev Should override this to change to your specific upgrade needs
    function _upgradeServiceManager() internal virtual {
        // 1. Deploy new ServiceManager implementation contract
        serviceManagerImplementation = address(
            new EigenDAServiceManager(
                avsDirectory,
                rewardsCoordinator,
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

    /// @notice Example createAVSRewardsSubmission call with the ServiceManager
    function _createAVSRewardsSubmission() internal {
        
        uint256 mockTokenInitialSupply = 1e30;
        // actual strategy addresses
        address stETHStrategy = 0x93c4b944D05dfe6df7645A86cd2206016c51564D;
        address rETHStrategy = 0x1BeE69b7dFFfA4E2d53C2a2Df135C388AD25dCD2;

        IRewardsCoordinator.StrategyAndMultiplier[] memory strategyAndMultipliers = new IRewardsCoordinator.StrategyAndMultiplier[](2);
        // Strategy addresses must be in ascending order
        strategyAndMultipliers[0] = IRewardsCoordinator.StrategyAndMultiplier({
            strategy: IStrategy(rETHStrategy),
            multiplier: 1e18
        });
        strategyAndMultipliers[1] = IRewardsCoordinator.StrategyAndMultiplier({
            strategy: IStrategy(stETHStrategy),
            multiplier: 1e18
        });

        IERC20 token = new ERC20PresetFixedSupply(
            "HARRYPOTTEROBAMASONIC10INU",
            "BITCOIN",
            mockTokenInitialSupply,
            deployerAddress
        );

        // must be in multiples of weeks i.e startTimestamp % 604800 == 0
        uint32 startTimestamp = 1710979200 + 2 weeks;
        // must be in multiples of weeks i.e duration % 604800 == 0
        uint32 duration = 10 weeks;
        // amount <= 1e38 - 1
        uint256 amount = 5000000e18;

        // Create RewardsSubmission input param
        IRewardsCoordinator.RewardsSubmission[]
            memory rewardsSubmissions = new IRewardsCoordinator.RewardsSubmission[](1);
        rewardsSubmissions[0] = IRewardsCoordinator.RewardsSubmission({
            strategiesAndMultipliers: strategyAndMultipliers,
            token: token,
            amount: amount,
            startTimestamp: startTimestamp,
            duration: duration
        });


        token.approve(serviceManager, amount);
        EigenDAServiceManager(serviceManager).createAVSRewardsSubmission(rewardsSubmissions);
    }

    /// @dev check implementation address set properly
    function _verifyUpgrade() internal virtual {
        // Preprod RewardsCoordinator
        require(
            address(rewardsCoordinator) == 0x7750d328b314EfFa365A0402CcfD489B80B0adda,
            "ServiceManagerBaseUpgrade: RewardsCoordinator address is incorrect"
        );
        require(
            avsProxyAdmin.getProxyImplementation(
                TransparentUpgradeableProxy(payable(serviceManager))
            ) == serviceManagerImplementation,
            "ServiceManagerBaseUpgrade: ServiceMananger implementation initially set incorrectly"
        );
    }
}
