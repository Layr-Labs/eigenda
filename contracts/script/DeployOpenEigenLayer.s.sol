// SPDX-License-Identifier: BUSL-1.1
pragma solidity =0.8.12;

import "@openzeppelin/contracts/token/ERC20/presets/ERC20PresetFixedSupply.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/beacon/UpgradeableBeacon.sol";

import "eigenlayer-core/contracts/interfaces/IETHPOSDeposit.sol";

import "eigenlayer-core/contracts/core/StrategyManager.sol";
import "eigenlayer-core/contracts/core/Slasher.sol";
import "eigenlayer-core/contracts/core/DelegationManager.sol";
import "eigenlayer-core/contracts/core/AVSDirectory.sol";
import "eigenlayer-core/contracts/core/RewardsCoordinator.sol";

import "eigenlayer-core/contracts/strategies/StrategyBaseTVLLimits.sol";

import "eigenlayer-core/contracts/pods/EigenPod.sol";
import "eigenlayer-core/contracts/pods/EigenPodManager.sol";

import "eigenlayer-core/contracts/permissions/PauserRegistry.sol";

import "eigenlayer-core/test/mocks/EmptyContract.sol";
import "eigenlayer-core/test/mocks/ETHDepositMock.sol";

import "forge-std/Script.sol";
import "forge-std/Test.sol";

// # To load the variables in the .env file
// source .env

// # To deploy and verify our contract
// forge script script/M1_Deploy.s.sol:Deployer_M1 --rpc-url $RPC_URL  --private-key $PRIVATE_KEY --broadcast -vvvv
contract DeployOpenEigenLayer is Script, Test {
    Vm cheats = Vm(VM_ADDRESS);

    uint32 CALCULATION_INTERVAL_SECONDS = 7 days;
    uint32 MAX_REWARDS_DURATION = 70 days;
    uint32 MAX_RETROACTIVE_LENGTH = 84 days;
    uint32 MAX_FUTURE_LENGTH = 28 days;
    uint32 GENESIS_REWARDS_TIMESTAMP = 1712188800;
    uint32 activationDelay = 7 days;
    uint16 globalCommissionBips = 1000;

    // struct used to encode token info in config file
    struct StrategyConfig {
        uint256 maxDeposits;
        uint256 maxPerDeposit;
        address tokenAddress;
        string tokenSymbol;
    }

    // EigenLayer Contracts
    ProxyAdmin public eigenLayerProxyAdmin;
    PauserRegistry public eigenLayerPauserReg;
    Slasher public slasher;
    Slasher public slasherImplementation;
    DelegationManager public delegation;
    DelegationManager public delegationImplementation;
    StrategyManager public strategyManager;
    StrategyManager public strategyManagerImplementation;
    EigenPodManager public eigenPodManager;
    EigenPodManager public eigenPodManagerImplementation;
    AVSDirectory public avsDirectory;
    AVSDirectory public avsDirectoryImplementation;
    RewardsCoordinator public rewardsCoordinator;
    RewardsCoordinator public rewardsCoordinatorImplementation;
    UpgradeableBeacon public eigenPodBeacon;
    EigenPod public eigenPodImplementation;
    StrategyBase public baseStrategyImplementation;

    EmptyContract public emptyContract;

    // the ETH2 deposit contract -- if not on mainnet, we deploy a mock as stand-in
    IETHPOSDeposit public ethPOSDeposit;

    // strategies deployed
    StrategyBaseTVLLimits[] public deployedStrategyArray;

    function _deployEigenLayer(
        address executorMultisig,
        address operationsMultisig,
        address pauserMultisig,
        StrategyConfig[] memory strategyConfigs
    ) internal {
        require(executorMultisig != address(0), "executorMultisig address not configured correctly!");
        require(operationsMultisig != address(0), "operationsMultisig address not configured correctly!");

        // deploy proxy admin for the ability to upgrade proxy contracts
        eigenLayerProxyAdmin = new ProxyAdmin();

        //deploy pauser registry
        {
            address[] memory pausers = new address[](3);
            pausers[0] = executorMultisig;
            pausers[1] = operationsMultisig;
            pausers[2] = pauserMultisig;
            eigenLayerPauserReg = new PauserRegistry(pausers, executorMultisig);
        }

        /**
         * First, deploy upgradeable proxy contracts that **will point** to the implementations. Since the implementation contracts are
         * not yet deployed, we give these proxies an empty contract as the initial implementation, to act as if they have no code.
         */
        emptyContract = new EmptyContract();
        delegation = DelegationManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenLayerProxyAdmin), ""))
        );
        avsDirectory = AVSDirectory(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenLayerProxyAdmin), ""))
        );
        strategyManager = StrategyManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenLayerProxyAdmin), ""))
        );
        slasher = Slasher(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenLayerProxyAdmin), ""))
        );
        eigenPodManager = EigenPodManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenLayerProxyAdmin), ""))
        );
        rewardsCoordinator = RewardsCoordinator(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(eigenLayerProxyAdmin), ""))
        );

        // ETH POS deposit is 0 address
        eigenPodImplementation = new EigenPod(
            ethPOSDeposit,
            eigenPodManager,
            1000 // temp genesis time
        );

        eigenPodBeacon = new UpgradeableBeacon(address(eigenPodImplementation));

        // Second, deploy the *implementation* contracts, using the *proxy contracts* as inputs
        delegationImplementation = new DelegationManager(strategyManager, slasher, eigenPodManager);
        avsDirectoryImplementation = new AVSDirectory(delegation);
        rewardsCoordinatorImplementation = new RewardsCoordinator(
            delegation,
            strategyManager,
            CALCULATION_INTERVAL_SECONDS,
            MAX_REWARDS_DURATION,
            MAX_RETROACTIVE_LENGTH,
            MAX_FUTURE_LENGTH,
            GENESIS_REWARDS_TIMESTAMP
        );
        strategyManagerImplementation = new StrategyManager(delegation, eigenPodManager, slasher);
        slasherImplementation = new Slasher(strategyManager, delegation);
        eigenPodManagerImplementation = new EigenPodManager(
            ethPOSDeposit,
            eigenPodBeacon,
            strategyManager,
            slasher,
            delegation
        );

        // Third, upgrade the proxy contracts to use the correct implementation contracts and initialize them.
        IStrategy[] memory _strategies;
        uint256[] memory _withdrawalDelayBlocks;
        eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(delegation))),
            address(delegationImplementation),
            abi.encodeWithSelector(DelegationManager.initialize.selector, executorMultisig, eigenLayerPauserReg, 0, 0, _strategies, _withdrawalDelayBlocks)
        );
        eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(avsDirectory))),
            address(avsDirectoryImplementation),
            abi.encodeWithSelector(AVSDirectory.initialize.selector, executorMultisig, eigenLayerPauserReg, 0)
        );
        eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(strategyManager))),
            address(strategyManagerImplementation),
            abi.encodeWithSelector(
                StrategyManager.initialize.selector,
                executorMultisig,
                operationsMultisig,
                eigenLayerPauserReg,
                0,
                0
            )
        );
        eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(rewardsCoordinator))),
            address(rewardsCoordinatorImplementation),
            abi.encodeWithSelector(
                RewardsCoordinator.initialize.selector,
                executorMultisig, 
                eigenLayerPauserReg,
                0, 
                executorMultisig,
                activationDelay,
                globalCommissionBips
            )
        );
        eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(slasher))),
            address(slasherImplementation),
            abi.encodeWithSelector(Slasher.initialize.selector, executorMultisig, eigenLayerPauserReg, 0)
        );
        eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenPodManager))),
            address(eigenPodManagerImplementation),
            abi.encodeWithSelector(
                EigenPodManager.initialize.selector,
                executorMultisig,
                eigenLayerPauserReg,
                0
            )
        );

        // deploy StrategyBaseTVLLimits contract implementation
        baseStrategyImplementation = new StrategyBaseTVLLimits(strategyManager);
        // create upgradeable proxies that each point to the implementation and initialize them
        for (uint256 i = 0; i < strategyConfigs.length; ++i) {
            deployedStrategyArray.push(
                StrategyBaseTVLLimits(
                    address(
                        new TransparentUpgradeableProxy(
                            address(baseStrategyImplementation),
                            address(eigenLayerProxyAdmin),
                            abi.encodeWithSelector(
                                StrategyBaseTVLLimits.initialize.selector,
                                strategyConfigs[i].maxPerDeposit,
                                strategyConfigs[i].maxDeposits,
                                IERC20(strategyConfigs[i].tokenAddress),
                                eigenLayerPauserReg
                            )
                        )
                    )
                )
            );
        }

        eigenLayerProxyAdmin.transferOwnership(executorMultisig);
        eigenPodBeacon.transferOwnership(executorMultisig);
    }
}