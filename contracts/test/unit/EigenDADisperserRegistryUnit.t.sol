// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import "../MockEigenDADeployer.sol";

contract EigenDADisperserRegistryUnit is MockEigenDADeployer {

    event DisperserAdded(uint32 indexed disperserKey, address disperserAddress);
    event DisperserRemoved(uint32 indexed disperserKey, address registrant);

    // Test user
    address public user = address(0x123);
    uint32 public disperserKey = 1;
    address public disperserAddress = address(0x456);

    function setUp() virtual public {
        _deployDA();
        // Fund the user and owner accounts for tests
        vm.deal(user, 11 ether);
        vm.deal(registryCoordinatorOwner, 11 ether);
    }

    function test_initalize() public {
        assertEq(eigenDADisperserRegistry.owner(), registryCoordinatorOwner);
        vm.expectRevert("Initializable: contract is already initialized");
        eigenDADisperserRegistry.initialize(address(this));
    }

    function test_registerDisperserInfo() public {
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
        });

        vm.expectEmit(address(eigenDADisperserRegistry));
        emit DisperserAdded(disperserKey, disperserAddress);
        
        vm.prank(user);
        eigenDADisperserRegistry.setDisperserInfo{value: 10 ether}(disperserKey, disperserInfo);

        // Verify registration data
        assertEq(eigenDADisperserRegistry.disperserKeyToAddress(disperserKey), disperserAddress);
        assertEq(eigenDADisperserRegistry.disperserRegistrants(disperserKey), user);
        assertEq(eigenDADisperserRegistry.disperserDeposits(disperserKey), 10 ether);
        assertEq(eigenDADisperserRegistry.isOwnerCreatedDisperser(disperserKey), false);
    }

    function test_updateDisperserInfo() public {
        // First register
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
        });
        vm.prank(user);
        eigenDADisperserRegistry.setDisperserInfo{value: 10 ether}(disperserKey, disperserInfo);
        
        // Then update with a new disperser address
        address newDisperserAddress = address(0x789);
        DisperserInfo memory newDisperserInfo = DisperserInfo({
            disperserAddress: newDisperserAddress
        });
        
        vm.prank(user);
        eigenDADisperserRegistry.setDisperserInfo{value: 10 ether}(disperserKey, newDisperserInfo);
        
        // Verify updated data
        assertEq(eigenDADisperserRegistry.disperserKeyToAddress(disperserKey), newDisperserAddress);
        assertEq(eigenDADisperserRegistry.disperserDeposits(disperserKey), 20 ether);
    }

    function test_deregisterDisperser() public {
        // First register
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
        });
        vm.prank(user);
        eigenDADisperserRegistry.setDisperserInfo{value: 10 ether}(disperserKey, disperserInfo);
        
        // Check user balance before deregistration
        uint256 balanceBefore = user.balance;
        
        // Deregister
        vm.expectEmit(address(eigenDADisperserRegistry));
        emit DisperserRemoved(disperserKey, user);
        
        vm.prank(user);
        eigenDADisperserRegistry.deregisterDisperser(disperserKey);
        
        // Check the refund (should be 90% of 0.1 ether)
        uint256 expectedRefund = 0.09 ether; // 0.1 ether - 10% tax
        assertEq(user.balance, balanceBefore + expectedRefund);
        
        // Verify everything is cleared
        assertEq(eigenDADisperserRegistry.disperserKeyToAddress(disperserKey), address(0));
        assertEq(eigenDADisperserRegistry.disperserRegistrants(disperserKey), address(0));
        assertEq(eigenDADisperserRegistry.disperserDeposits(disperserKey), 0);
        assertEq(eigenDADisperserRegistry.isOwnerCreatedDisperser(disperserKey), false);
    }

    function test_setDisperserInfo_revert_insufficientDeposit() public {
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
        });

        vm.expectRevert("Deposit of 10 ether required for registration");
        vm.prank(user);
        eigenDADisperserRegistry.setDisperserInfo{value: 0.05 ether}(disperserKey, disperserInfo);
    }

    function test_setDisperserInfo_revert_unauthorizedUpdate() public {
        // First register with user
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
        });
        vm.prank(user);
        eigenDADisperserRegistry.setDisperserInfo{value: 10 ether}(disperserKey, disperserInfo);
        
        // Then try to update with a different address
        address anotherUser = address(0x789);
        vm.deal(anotherUser, 1 ether);
        
        vm.expectRevert("Caller is not the registrant");
        vm.prank(anotherUser);
        eigenDADisperserRegistry.setDisperserInfo{value: 10 ether}(disperserKey, disperserInfo);
    }

    function test_deregisterDisperser_revert_unauthorizedDeregistration() public {
        // First register with user
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
        });
        vm.prank(user);
        eigenDADisperserRegistry.setDisperserInfo{value: 10 ether}(disperserKey, disperserInfo);
        
        // Then try to deregister with a different address
        address anotherUser = address(0x789);
        
        vm.expectRevert("Caller is not the registrant");
        vm.prank(anotherUser);
        eigenDADisperserRegistry.deregisterDisperser(disperserKey);
    }
    
    function test_ownerCreatedDisperser() public {
        uint32 ownerDisperserKey = 2;
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
        });

        vm.expectEmit(address(eigenDADisperserRegistry));
        emit DisperserAdded(ownerDisperserKey, disperserAddress);
        
        // Register as owner with deposit
        vm.prank(registryCoordinatorOwner);
        eigenDADisperserRegistry.setDisperserInfo{value: 10 ether}(ownerDisperserKey, disperserInfo);

        // Verify registration data
        assertEq(eigenDADisperserRegistry.disperserKeyToAddress(ownerDisperserKey), disperserAddress);
        assertEq(eigenDADisperserRegistry.disperserRegistrants(ownerDisperserKey), registryCoordinatorOwner);
        assertEq(eigenDADisperserRegistry.disperserDeposits(ownerDisperserKey), 10 ether);
        assertEq(eigenDADisperserRegistry.isOwnerCreatedDisperser(ownerDisperserKey), true);
    }
    
    function test_ownerUpdateDisperser() public {
        // First register as owner
        uint32 ownerDisperserKey = 2;
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
        });
        
        vm.prank(registryCoordinatorOwner);
        eigenDADisperserRegistry.setDisperserInfo{value: 10 ether}(ownerDisperserKey, disperserInfo);
        
        // Then update with a new disperser address with deposit
        address newDisperserAddress = address(0x789);
        DisperserInfo memory newDisperserInfo = DisperserInfo({
            disperserAddress: newDisperserAddress
        });
        
        vm.prank(registryCoordinatorOwner);
        eigenDADisperserRegistry.setDisperserInfo{value: 10 ether}(ownerDisperserKey, newDisperserInfo);
        
        // Verify updated data
        assertEq(eigenDADisperserRegistry.disperserKeyToAddress(ownerDisperserKey), newDisperserAddress);
        assertEq(eigenDADisperserRegistry.disperserDeposits(ownerDisperserKey), 20 ether);
        assertEq(eigenDADisperserRegistry.isOwnerCreatedDisperser(ownerDisperserKey), true);
    }
    
    function test_ownerDeregisterDisperser() public {
        // First register as owner with deposit
        uint32 ownerDisperserKey = 2;
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
        });
        
        vm.prank(registryCoordinatorOwner);
        eigenDADisperserRegistry.setDisperserInfo{value: 10 ether}(ownerDisperserKey, disperserInfo);
        
        // Check owner balance before deregistration
        uint256 balanceBefore = registryCoordinatorOwner.balance;
        
        // Deregister
        vm.expectEmit(address(eigenDADisperserRegistry));
        emit DisperserRemoved(ownerDisperserKey, registryCoordinatorOwner);
        
        vm.prank(registryCoordinatorOwner);
        eigenDADisperserRegistry.deregisterDisperser(ownerDisperserKey);
        
        // Check the refund (should be 90% of 10 ether)
        uint256 expectedRefund = 9 ether; // 10 ether - 10% tax
        assertEq(registryCoordinatorOwner.balance, balanceBefore + expectedRefund);
        
        // Verify everything is cleared
        assertEq(eigenDADisperserRegistry.disperserKeyToAddress(ownerDisperserKey), address(0));
        assertEq(eigenDADisperserRegistry.disperserRegistrants(ownerDisperserKey), address(0));
        assertEq(eigenDADisperserRegistry.disperserDeposits(ownerDisperserKey), 0);
        assertEq(eigenDADisperserRegistry.isOwnerCreatedDisperser(ownerDisperserKey), false);
    }
    
    function test_owner_revert_insufficientDeposit() public {
        // Even the owner must pay the full deposit
        uint32 ownerDisperserKey = 2;
        DisperserInfo memory disperserInfo = DisperserInfo({
            disperserAddress: disperserAddress
        });

        vm.expectRevert("Deposit of 10 ether required for registration");
        vm.prank(registryCoordinatorOwner);
        eigenDADisperserRegistry.setDisperserInfo{value: 0.05 ether}(ownerDisperserKey, disperserInfo);
    }
}
