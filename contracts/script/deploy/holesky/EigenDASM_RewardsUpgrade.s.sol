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
 */
contract ServiceManagerBaseUpgrade is ExistingDeploymentParser {
    // Hardcode these values to your needs
    address public serviceManager = 0xD4A7E1Bd8015057293f0D0A557088c286942e84b;
    address public serviceManagerImplementation = 0xdBFd6C8582b58590C4AFF40BdF15488A086bC672;

    ProxyAdmin public avsProxyAdmin = ProxyAdmin(0xB043055dd967A382577c2f5261fA6428f2905c15);
    address deployerAddress = 0xDA29BB71669f46F2a779b4b62f03644A84eE3479;
    address registryCoordinator = 0x53012C69A189cfA2D9d29eb6F19B32e0A2EA3490;
    address stakeRegistry = 0xBDACD5998989Eec814ac7A0f0f6596088AA2a270;

    function run(string memory deployArg) external {
        // 1. Setup and parse existing EigenLayer Holesky preprod contracts
        _parseInitialDeploymentParams(
            "script/deploy/holesky/config/eigenlayer.config.json"
        );
        _parseDeployedContracts(
            "script/deploy/holesky/config/eigenlayer_addresses.config.json"
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
        IPauserRegistry pauserRegistry = IPauserRegistry(0x85Ef7299F8311B25642679edBF02B62FA2212F06);
        address emptyContract = 0x9690d52B1Ce155DB2ec5eCbF5a262ccCc7B3A6D2;
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
        // actual strategy addresses
        address stETHStrategy = 0x7D704507b76571a51d9caE8AdDAbBFd0ba0e63d3;
        address rETHStrategy = 0x3A8fBdf9e77DFc25d09741f51d3E181b25d0c4E0;

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

        // IERC20 token = new ERC20PresetFixedSupply(
        //     "HARRYPOTTEROBAMASONIC10INU",
        //     "BITCOIN",
        //     mockTokenInitialSupply,
        //     msg.sender
        // );

        IERC20 token = IERC20(0x3B78576F7D6837500bA3De27A60c7f594934027E);


        // must be in multiples of weeks i.e startTimestamp % 604800 == 0
        uint32 startTimestamp = 1714608000 + 8 weeks;
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

        // Set rewardsInitiator
        // EigenDAServiceManager(serviceManager).setRewardsInitiator(msg.sender);

        token.approve(serviceManager, amount);
        EigenDAServiceManager(serviceManager).createAVSRewardsSubmission(rewardsSubmissions);
    }

    /// @dev check implementation address set properly
    function _verifyUpgrade() internal virtual {
        // Preprod RewardsCoordinator
        require(
            address(rewardsCoordinator) == 0xAcc1fb458a1317E886dB376Fc8141540537E68fE,
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