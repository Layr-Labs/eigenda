// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {Test} from "lib/forge-std/src/Test.sol";
import {EigenDADirectory} from "src/core/EigenDADirectory.sol";
import {AddressDirectoryConstants} from "src/core/libraries/v3/address-directory/AddressDirectoryConstants.sol";
import {EigenDAAccessControl} from "src/core/EigenDAAccessControl.sol";
import {IEigenDAAddressDirectory} from "src/core/interfaces/IEigenDADirectory.sol";

contract EigenDADirectoryTest is Test {
    EigenDADirectory public directory;
    EigenDAAccessControl public accessControl;

    address owner = makeAddr("owner");
    address nonOwner = makeAddr("nonOwner");

    address testAddress = makeAddr("testAddr");
    string testNamedKey = "testNamedKey";

    function setUp() public {
        // Deploy AccessControl with owner
        accessControl = new EigenDAAccessControl(owner);

        // Deploy and initialize DA Directory
        directory = new EigenDADirectory();
        directory.initialize(address(accessControl));
    }

    // ===========================
    // Address Directory: Basic Operations
    // ===========================

    function test_initialize() public {
        string[] memory names = directory.getAllNames();
        assertNotEq(
            directory.getAddress(AddressDirectoryConstants.ACCESS_CONTROL_NAME),
            address(0x0),
            "AccessControl contract should have entry"
        );
        assertEq(names.length, 1, "Should have one name (AccessControl) after initialization");

        vm.expectRevert("AlreadyInitialized()");
        directory.initialize(address(0));
    }

    function test_addAddress_success() public {
        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit IEigenDAAddressDirectory.AddressAdded(testNamedKey, keccak256(abi.encodePacked(testNamedKey)), testAddress);
        directory.addAddress(testNamedKey, testAddress);

        assertEq(directory.getAddress(testNamedKey), testAddress, "Address should be set correctly");
    }

    function test_addAddress_revertZeroAddress() public {
        vm.prank(owner);
        vm.expectRevert(IEigenDAAddressDirectory.ZeroAddress.selector);
        directory.addAddress(testNamedKey, address(0));
    }

    function test_addAddress_revertAlreadyExists() public {
        vm.startPrank(owner);
        directory.addAddress(testNamedKey, testAddress);

        vm.expectRevert(abi.encodeWithSelector(IEigenDAAddressDirectory.AddressAlreadyExists.selector, testNamedKey));
        directory.addAddress(testNamedKey, address(0x5678));
        vm.stopPrank();
    }

    function test_addAddress_revertNonOwner() public {
        vm.prank(nonOwner);
        vm.expectRevert("Caller is not the owner");
        directory.addAddress(testNamedKey, testAddress);
    }

    function test_replaceAddress_success() public {
        address oldAddress = address(0x1234);
        address newAddress = address(0x5678);

        vm.startPrank(owner);
        directory.addAddress(testNamedKey, oldAddress);

        vm.expectEmit(true, true, true, true);
        emit IEigenDAAddressDirectory.AddressReplaced(
            testNamedKey, keccak256(abi.encodePacked(testNamedKey)), oldAddress, newAddress
        );
        directory.replaceAddress(testNamedKey, newAddress);
        vm.stopPrank();

        assertEq(directory.getAddress(testNamedKey), newAddress, "Address should be replaced");
    }

    function test_replaceAddress_revertDoesNotExist() public {
        address newAddress = address(0x5678);

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IEigenDAAddressDirectory.AddressDoesNotExist.selector, testNamedKey));
        directory.replaceAddress(testNamedKey, newAddress);
    }

    function test_replaceAddress_revertZeroAddress() public {
        address oldAddress = address(0x1234);

        vm.startPrank(owner);
        directory.addAddress(testNamedKey, oldAddress);

        vm.expectRevert(IEigenDAAddressDirectory.ZeroAddress.selector);
        directory.replaceAddress(testNamedKey, address(0));
        vm.stopPrank();
    }

    function test_replaceAddress_revertSameValue() public {
        vm.startPrank(owner);
        directory.addAddress(testNamedKey, testAddress);

        vm.expectRevert(abi.encodeWithSelector(IEigenDAAddressDirectory.NewValueIsOldValue.selector, testAddress));
        directory.replaceAddress(testNamedKey, testAddress);
        vm.stopPrank();
    }

    function test_replaceAddress_revertNonOwner() public {
        address oldAddress = address(0x1234);
        address newAddress = address(0x5678);

        vm.prank(owner);
        directory.addAddress(testNamedKey, oldAddress);

        vm.prank(nonOwner);
        vm.expectRevert("Caller is not the owner");
        directory.replaceAddress(testNamedKey, newAddress);
    }

    function test_removeAddress_success() public {
        vm.startPrank(owner);
        directory.addAddress(testNamedKey, testAddress);

        vm.expectEmit(true, true, true, true);
        emit IEigenDAAddressDirectory.AddressRemoved(testNamedKey, keccak256(abi.encodePacked(testNamedKey)));
        directory.removeAddress(testNamedKey);
        vm.stopPrank();

        assertEq(directory.getAddress(testNamedKey), address(0), "Address should be removed");
    }

    function test_removeAddress_revertDoesNotExist() public {
        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IEigenDAAddressDirectory.AddressDoesNotExist.selector, testNamedKey));
        directory.removeAddress(testNamedKey);
    }

    function test_removeAddress_revertNonOwner() public {
        vm.prank(owner);
        directory.addAddress(testNamedKey, testAddress);

        vm.prank(nonOwner);
        vm.expectRevert("Caller is not the owner");
        directory.removeAddress(testNamedKey);
    }

    function test_getAddress_byString() public {
        vm.prank(owner);
        directory.addAddress(testNamedKey, testAddress);

        assertEq(directory.getAddress(testNamedKey), testAddress, "Should retrieve address by name");
    }

    function test_getAddress_byBytes32() public {
        address localTestAddress = address(0x1234);
        string memory localTestKeyName = "testAddress";
        bytes32 nameDigest = keccak256(abi.encodePacked(localTestKeyName));

        vm.prank(owner);
        directory.addAddress(localTestKeyName, localTestAddress);

        assertEq(directory.getAddress(nameDigest), localTestAddress, "Should retrieve address by digest");
    }

    function test_getAddress_nonexistent() public view {
        string memory unknownTestNameKey = "nonexistentAddress";
        assertEq(
            directory.getAddress(unknownTestNameKey), address(0), "Should return zero address for nonexistent name"
        );
    }

    function test_getName_success() public {
        bytes32 nameDigest = keccak256(abi.encodePacked(testNamedKey));

        vm.prank(owner);
        directory.addAddress(testNamedKey, testAddress);

        assertEq(directory.getName(nameDigest), testNamedKey, "Should retrieve name by digest");
    }

    function test_getName_nonexistent() public view {
        bytes32 nonexistentDigest = keccak256(abi.encodePacked("nonexistent"));
        assertEq(directory.getName(nonexistentDigest), "", "Should return empty string for nonexistent digest");
    }

    function test_getAllNames_multipleAddresses() public {
        vm.startPrank(owner);
        directory.addAddress("address1", address(0x1));
        directory.addAddress("address2", address(0x2));
        directory.addAddress("address3", address(0x3));
        vm.stopPrank();

        string[] memory names = directory.getAllNames();
        assertEq(names.length, 4, "Should have 4 names (3 added + AccessControl)");

        // Verify the added names are present (order not guaranteed)
        bool foundAddress1 = false;
        bool foundAddress2 = false;
        bool foundAddress3 = false;

        for (uint256 i = 0; i < names.length; i++) {
            if (keccak256(bytes(names[i])) == keccak256(bytes("address1"))) foundAddress1 = true;
            if (keccak256(bytes(names[i])) == keccak256(bytes("address2"))) foundAddress2 = true;
            if (keccak256(bytes(names[i])) == keccak256(bytes("address3"))) foundAddress3 = true;
        }

        assertTrue(foundAddress1, "address1 should be in the list");
        assertTrue(foundAddress2, "address2 should be in the list");
        assertTrue(foundAddress3, "address3 should be in the list");
    }

    function test_getAllNames_afterRemoval() public {
        vm.startPrank(owner);
        directory.addAddress("address1", address(0x1));
        directory.addAddress("address2", address(0x2));
        directory.addAddress("address3", address(0x3));
        directory.removeAddress("address2");
        vm.stopPrank();

        string[] memory names = directory.getAllNames();
        assertEq(names.length, 3, "Should have 3 names after removal (address1, address3, AccessControl)");

        // Verify address2 is not present
        for (uint256 i = 0; i < names.length; i++) {
            assertTrue(keccak256(bytes(names[i])) != keccak256(bytes("address2")), "address2 should not be in the list");
        }
    }

    // ===========================
    // Address Directory: Edge Cases
    // ===========================

    function test_addAndReplace_multipleTimes() public {
        vm.startPrank(owner);
        directory.addAddress(testNamedKey, address(0x1));
        assertEq(directory.getAddress(testNamedKey), address(0x1), "First address should be set");

        directory.replaceAddress(testNamedKey, address(0x2));
        assertEq(directory.getAddress(testNamedKey), address(0x2), "Second address should be set");

        directory.replaceAddress(testNamedKey, address(0x3));
        assertEq(directory.getAddress(testNamedKey), address(0x3), "Third address should be set");
        vm.stopPrank();
    }

    function test_removeAndReAdd() public {
        vm.startPrank(owner);
        directory.addAddress(testNamedKey, testAddress);
        directory.removeAddress(testNamedKey);

        // Should be able to add again after removal
        directory.addAddress(testNamedKey, testAddress);
        assertEq(directory.getAddress(testNamedKey), testAddress, "Should be able to re-add after removal");
        vm.stopPrank();
    }
}
