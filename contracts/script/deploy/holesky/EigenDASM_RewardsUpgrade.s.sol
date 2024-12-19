// SPDX-License-Identifier: BUSL-1.1
/*
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
 * @title ServiceManagerBaseUpgrade for Preprod contracts.
 * Assumes EOA deploying has permissions to call the proxyAdmin to upgrade.
 *
 *
 * Local Fork: Deploy/Upgrade RewardsCoordinator
 * anvil --fork-url $RPC_HOLESKY
 * forge script script/deploy/holesky/EigenDASM_RewardsUpgrade.s.sol:ServiceManagerBaseUpgrade --rpc-url http://127.0.0.1:8545 --private-key $PRIVATE_KEY --broadcast -vvvv --sig "run(string memory deployArg)" upgrade
 * forge script script/deploy/holesky/EigenDASM_RewardsUpgrade.s.sol:ServiceManagerBaseUpgrade --rpc-url http://127.0.0.1:8545 --private-key $PRIVATE_KEY --broadcast -vvvv --sig "run(string memory deployArg)" deploy
 * forge script script/deploy/holesky/EigenDASM_RewardsUpgrade.s.sol:ServiceManagerBaseUpgrade --rpc-url http://127.0.0.1:8545 --private-key $PRIVATE_KEY --broadcast -vvvv --sig "run(string memory deployArg)" createAVSRewardsSubmission
 *
 * Upgrade Holesky testnet: Deploy/Upgrade RewardsCoordinator
 * forge script script/deploy/holesky/EigenDASM_RewardsUpgrade.s.sol:ServiceManagerBaseUpgrade --rpc-url $RPC_HOLESKY --private-key $PRIVATE_KEY --broadcast --verify -vvvv --sig "run(string memory deployArg)" upgrade
 * forge script script/deploy/holesky/EigenDASM_RewardsUpgrade.s.sol:ServiceManagerBaseUpgrade --rpc-url $RPC_HOLESKY --private-key $PRIVATE_KEY --broadcast --verify -vvvv --sig "run(string memory deployArg)" deploy
 * forge script script/deploy/holesky/EigenDASM_RewardsUpgrade.s.sol:ServiceManagerBaseUpgrade --rpc-url $RPC_HOLESKY --private-key $PRIVATE_KEY --broadcast --verify -vvvv --sig "run(string memory deployArg)" createAVSRewardsSubmission
 *//*
contract ServiceManagerBaseUpgrade is ExistingDeploymentParser {
    // Hardcode these values to your needs
    address public serviceManager = 0x54A03db2784E3D0aCC08344D05385d0b62d4F432;
    address public serviceManagerImplementation = 0xFe779fB43280A92cd85466312E2AE8A4F1A48007;
    ProxyAdmin public avsProxyAdmin = ProxyAdmin(0x9Fd7E279f5bD692Dc04792151E14Ad814FC60eC1);
    address deployerAddress = 0xDA29BB71669f46F2a779b4b62f03644A84eE3479;
    address registryCoordinator = 0x2c61EA360D6500b58E7f481541A36B443Bc858c6;
    address stakeRegistry = 0x53668EBf2e28180e38B122c641BC51Ca81088871;

    function run(string memory deployArg) external {
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

        if (keccak256(abi.encode(deployArg)) == keccak256(abi.encode("upgrade"))) {
            _upgradeServiceManager();
        } else if (keccak256(abi.encode(deployArg)) == keccak256(abi.encode("deploy"))) {
            _deployServiceManager();
        } else if (keccak256(abi.encode(deployArg)) == keccak256(abi.encode("createAVSRewardsSubmission"))) {
            _createAVSRewardsSubmission();
        }

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

    function _deployServiceManager() internal virtual {
        IPauserRegistry pauserRegistry = IPauserRegistry(0x9Ab2FEAf0465f0eD51Fc2b663eF228B418c9Dad1);
        address emptyContract = 0xc08b788d587F927b49665b90ab35D5224965f3d9;
        uint256 initialPausedStatus = 0;
        address initialOwner = deployerAddress;
        address[] memory batchConfirmers;

        // 1. Deploy new ServiceManager implementation contract
        serviceManagerImplementation = address(
            new EigenDAServiceManager(
                avsDirectory,
                rewardsCoordinator,
                IRegistryCoordinator(registryCoordinator),
                IStakeRegistry(stakeRegistry)
            )
        );

        // 2. Deploy new TUPS and initialize
        serviceManager = address(
            new TransparentUpgradeableProxy(emptyContract, address(avsProxyAdmin), "")
        );

        avsProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(serviceManager))),
            address(serviceManagerImplementation),
            abi.encodeWithSelector(
                EigenDAServiceManager.initialize.selector,
                eigenLayerPauserReg,
                initialPausedStatus,
                initialOwner,
                batchConfirmers
            )
        );
    }

    /// @notice Example createAVSRewardsSubmission call with the ServiceManager
    function _createAVSRewardsSubmission() internal {
        uint256 mockTokenInitialSupply = 1e30;
        address stETHStrategy = 0x5C8b55722f421556a2AAfb7A3EA63d4c3e514312;
        address rETHStrategy = 0x87f6C7d24b109919eB38295e3F8298425e6331D9;

        IRewardsCoordinator.StrategyAndMultiplier[] memory strategyAndMultipliers = new IRewardsCoordinator.StrategyAndMultiplier[](2);
        // Strategy addresses must be in ascending order
        strategyAndMultipliers[0] = IRewardsCoordinator.StrategyAndMultiplier({
            strategy: IStrategy(stETHStrategy),
            multiplier: 1e18
        });
        strategyAndMultipliers[1] = IRewardsCoordinator.StrategyAndMultiplier({
            strategy: IStrategy(rETHStrategy),
            multiplier: 1e18
        });

        IERC20 token = new ERC20PresetFixedSupply(
            "dog wif hat",
            "MOCK1",
            mockTokenInitialSupply,
            msg.sender
        );
        // must be in multiples of weeks i.e startTimestamp % 604800 == 0
        uint32 startTimestamp = 1714608000;
        // must be in multiples of weeks i.e duration % 604800 == 0
        uint32 duration = 1 weeks;
        // amount <= 1e38 - 1
        uint256 amount = 100e18;

        // 2. Create RewardsSubmission input param
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
            address(rewardsCoordinator) == 0xb22Ef643e1E067c994019A4C19e403253C05c2B0,
            "ServiceManagerBaseUpgrade: RewardsCoordinator address is incorrect"
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
*/