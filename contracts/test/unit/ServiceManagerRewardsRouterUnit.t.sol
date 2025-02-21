// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import "../MockEigenDADeployer.sol";
import {ServiceManagerRewardsRouter} from "../../src/rewards/ServiceManagerRewardsRouter.sol";
import {IServiceManagerRewardsRouter} from "../../src/rewards/IServiceManagerRewardsRouter.sol";

contract ServiceManagerRewardsRouterUnit is MockEigenDADeployer {

    IServiceManagerRewardsRouter serviceManagerRewardsRouter;
    IServiceManagerRewardsRouter serviceManagerRewardsRouterImplementation;
    
    function setUp() virtual public {
        _deployDA();

        serviceManagerRewardsRouterImplementation = new ServiceManagerRewardsRouter(eigenDAServiceManager);
        serviceManagerRewardsRouter = IServiceManagerRewardsRouter(
            address(
                new TransparentUpgradeableProxy(address(serviceManagerRewardsRouterImplementation), address(proxyAdmin), "")
            )
        );
        cheats.prank(proxyAdminOwner);
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(serviceManagerRewardsRouter))),
            address(serviceManagerRewardsRouterImplementation),
            abi.encodeWithSelector(
                ServiceManagerRewardsRouter.initialize.selector,
                registryCoordinatorOwner
            )
        );

        vm.prank(registryCoordinatorOwner);
        eigenDAServiceManager.setRewardsInitiator(address(serviceManagerRewardsRouter));
    }

    function test_setAllowedToPushRewards() public {
        vm.prank(registryCoordinatorOwner);
        serviceManagerRewardsRouter.setAllowedToPushRewards(rewardsInitiator, true);
        assertEq(serviceManagerRewardsRouter.isAllowedToPushRewards(rewardsInitiator), true);

        vm.prank(registryCoordinatorOwner);
        serviceManagerRewardsRouter.setAllowedToPushRewards(rewardsInitiator, false);
        assertEq(serviceManagerRewardsRouter.isAllowedToPushRewards(rewardsInitiator), false);
    }

    function test_setAllowedToPushRewards_reverts_NotOwner() public {
        vm.expectRevert("Ownable: caller is not the owner");
        serviceManagerRewardsRouter.setAllowedToPushRewards(rewardsInitiator, true);
    }

    function test_createAVSRewardsSubmission() public {
        vm.prank(registryCoordinatorOwner);
        serviceManagerRewardsRouter.setAllowedToPushRewards(rewardsInitiator, true);

        IRewardsCoordinator.RewardsSubmission[] memory rewardsSubmissions = new IRewardsCoordinator.RewardsSubmission[](0);
            
        vm.prank(rewardsInitiator);
        serviceManagerRewardsRouter.createAVSRewardsSubmission(rewardsSubmissions);
    }

    function test_createAVSRewardsSubmission_reverts_NotAllowed() public {
        vm.prank(registryCoordinatorOwner);
        serviceManagerRewardsRouter.setAllowedToPushRewards(rewardsInitiator, false);

        IRewardsCoordinator.RewardsSubmission[] memory rewardsSubmissions = new IRewardsCoordinator.RewardsSubmission[](0);

        vm.expectRevert("ServiceManagerRewardsRouter: Address not permissioned to push rewards");
        serviceManagerRewardsRouter.createAVSRewardsSubmission(rewardsSubmissions);
    }

    function test_createOperatorDirectedAVSRewardsSubmission() public {
        vm.prank(registryCoordinatorOwner);
        serviceManagerRewardsRouter.setAllowedToPushRewards(rewardsInitiator, true);

        IRewardsCoordinator.OperatorDirectedRewardsSubmission[] memory operatorDirectedRewardsSubmissions = new IRewardsCoordinator.OperatorDirectedRewardsSubmission[](0);

        vm.prank(rewardsInitiator);
        serviceManagerRewardsRouter.createOperatorDirectedAVSRewardsSubmission(operatorDirectedRewardsSubmissions);
    }

    function test_createOperatorDirectedAVSRewardsSubmission_reverts_NotAllowed() public {
        vm.prank(registryCoordinatorOwner);
        serviceManagerRewardsRouter.setAllowedToPushRewards(rewardsInitiator, false);

        IRewardsCoordinator.OperatorDirectedRewardsSubmission[] memory operatorDirectedRewardsSubmissions = new IRewardsCoordinator.OperatorDirectedRewardsSubmission[](0);

        vm.expectRevert("ServiceManagerRewardsRouter: Address not permissioned to push rewards");
        serviceManagerRewardsRouter.createOperatorDirectedAVSRewardsSubmission(operatorDirectedRewardsSubmissions);
    }
}