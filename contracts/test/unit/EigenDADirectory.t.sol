// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {Test} from "lib/forge-std/src/Test.sol";
import {EigenDADirectory} from "src/core/EigenDADirectory.sol";
import {ConfigRegistryTypes} from "src/core/libraries/v3/config-registry/ConfigRegistryTypes.sol";
import {EigenDAAccessControl} from "src/core/EigenDAAccessControl.sol";
import {IEigenDAAddressDirectory} from "src/core/interfaces/IEigenDADirectory.sol";

contract EigenDADirectoryTest is Test {
    EigenDADirectory public directory;
    EigenDAAccessControl public accessControl;

    address owner = makeAddr("owner");
    address nonOwner = makeAddr("nonOwner");

    string constant CONFIG_NAME_BLOCKNUMBER = "testConfigBlockNumber";
    string constant CONFIG_NAME_TIMESTAMP = "testConfigTimestamp";

    function setUp() public {
        // Deploy access control with owner
        accessControl = new EigenDAAccessControl(owner);

        // Deploy and initialize DA Directory
        directory = new EigenDADirectory();
        directory.initialize(address(accessControl));
    }

    // ===========================
    // Address Directory: Basic Operations
    // ===========================

    function test_addAddress_success() public {
        address testAddress = address(0x1234);
        string memory testName = "testAddress";

        vm.prank(owner);
        vm.expectEmit(true, true, true, true);
        emit IEigenDAAddressDirectory.AddressAdded(testName, keccak256(abi.encodePacked(testName)), testAddress);
        directory.addAddress(testName, testAddress);

        assertEq(directory.getAddress(testName), testAddress, "Address should be set correctly");
    }

    function test_addAddress_revertZeroAddress() public {
        string memory testName = "testAddress";

        vm.prank(owner);
        vm.expectRevert(IEigenDAAddressDirectory.ZeroAddress.selector);
        directory.addAddress(testName, address(0));
    }

    function test_addAddress_revertAlreadyExists() public {
        address testAddress = address(0x1234);
        string memory testName = "testAddress";

        vm.startPrank(owner);
        directory.addAddress(testName, testAddress);

        vm.expectRevert(abi.encodeWithSelector(IEigenDAAddressDirectory.AddressAlreadyExists.selector, testName));
        directory.addAddress(testName, address(0x5678));
        vm.stopPrank();
    }

    function test_addAddress_revertNonOwner() public {
        address testAddress = address(0x1234);
        string memory testName = "testAddress";

        vm.prank(nonOwner);
        vm.expectRevert("Caller is not the owner");
        directory.addAddress(testName, testAddress);
    }

    function test_replaceAddress_success() public {
        address oldAddress = address(0x1234);
        address newAddress = address(0x5678);
        string memory testName = "testAddress";

        vm.startPrank(owner);
        directory.addAddress(testName, oldAddress);

        vm.expectEmit(true, true, true, true);
        emit IEigenDAAddressDirectory.AddressReplaced(
            testName, keccak256(abi.encodePacked(testName)), oldAddress, newAddress
        );
        directory.replaceAddress(testName, newAddress);
        vm.stopPrank();

        assertEq(directory.getAddress(testName), newAddress, "Address should be replaced");
    }

    function test_replaceAddress_revertDoesNotExist() public {
        string memory testName = "nonexistentAddress";
        address newAddress = address(0x5678);

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IEigenDAAddressDirectory.AddressDoesNotExist.selector, testName));
        directory.replaceAddress(testName, newAddress);
    }

    function test_replaceAddress_revertZeroAddress() public {
        address oldAddress = address(0x1234);
        string memory testName = "testAddress";

        vm.startPrank(owner);
        directory.addAddress(testName, oldAddress);

        vm.expectRevert(IEigenDAAddressDirectory.ZeroAddress.selector);
        directory.replaceAddress(testName, address(0));
        vm.stopPrank();
    }

    function test_replaceAddress_revertSameValue() public {
        address testAddress = address(0x1234);
        string memory testName = "testAddress";

        vm.startPrank(owner);
        directory.addAddress(testName, testAddress);

        vm.expectRevert(abi.encodeWithSelector(IEigenDAAddressDirectory.NewValueIsOldValue.selector, testAddress));
        directory.replaceAddress(testName, testAddress);
        vm.stopPrank();
    }

    function test_replaceAddress_revertNonOwner() public {
        address oldAddress = address(0x1234);
        address newAddress = address(0x5678);
        string memory testName = "testAddress";

        vm.prank(owner);
        directory.addAddress(testName, oldAddress);

        vm.prank(nonOwner);
        vm.expectRevert("Caller is not the owner");
        directory.replaceAddress(testName, newAddress);
    }

    function test_removeAddress_success() public {
        address testAddress = address(0x1234);
        string memory testName = "testAddress";

        vm.startPrank(owner);
        directory.addAddress(testName, testAddress);

        vm.expectEmit(true, true, true, true);
        emit IEigenDAAddressDirectory.AddressRemoved(testName, keccak256(abi.encodePacked(testName)));
        directory.removeAddress(testName);
        vm.stopPrank();

        assertEq(directory.getAddress(testName), address(0), "Address should be removed");
    }

    function test_removeAddress_revertDoesNotExist() public {
        string memory testName = "nonexistentAddress";

        vm.prank(owner);
        vm.expectRevert(abi.encodeWithSelector(IEigenDAAddressDirectory.AddressDoesNotExist.selector, testName));
        directory.removeAddress(testName);
    }

    function test_removeAddress_revertNonOwner() public {
        address testAddress = address(0x1234);
        string memory testName = "testAddress";

        vm.prank(owner);
        directory.addAddress(testName, testAddress);

        vm.prank(nonOwner);
        vm.expectRevert("Caller is not the owner");
        directory.removeAddress(testName);
    }

    function test_getAddress_byString() public {
        address testAddress = address(0x1234);
        string memory testName = "testAddress";

        vm.prank(owner);
        directory.addAddress(testName, testAddress);

        assertEq(directory.getAddress(testName), testAddress, "Should retrieve address by name");
    }

    function test_getAddress_byBytes32() public {
        address testAddress = address(0x1234);
        string memory testName = "testAddress";
        bytes32 nameDigest = keccak256(abi.encodePacked(testName));

        vm.prank(owner);
        directory.addAddress(testName, testAddress);

        assertEq(directory.getAddress(nameDigest), testAddress, "Should retrieve address by digest");
    }

    function test_getAddress_nonexistent() public view {
        string memory testName = "nonexistentAddress";
        assertEq(directory.getAddress(testName), address(0), "Should return zero address for nonexistent name");
    }

    function test_getName_success() public {
        address testAddress = address(0x1234);
        string memory testName = "testAddress";
        bytes32 nameDigest = keccak256(abi.encodePacked(testName));

        vm.prank(owner);
        directory.addAddress(testName, testAddress);

        assertEq(directory.getName(nameDigest), testName, "Should retrieve name by digest");
    }

    function test_getName_nonexistent() public view {
        bytes32 nonexistentDigest = keccak256(abi.encodePacked("nonexistent"));
        assertEq(directory.getName(nonexistentDigest), "", "Should return empty string for nonexistent digest");
    }

    function test_getAllNames_empty() public view {
        string[] memory names = directory.getAllNames();
        // Note: AccessControl is registered during initialization, so we expect 1 name
        assertEq(names.length, 1, "Should have one name (AccessControl) after initialization");
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
        assertEq(names.length, 3, "Should have 3 names after removal (address1, address3, access control)");

        // Verify address2 is not present
        for (uint256 i = 0; i < names.length; i++) {
            assertTrue(keccak256(bytes(names[i])) != keccak256(bytes("address2")), "address2 should not be in the list");
        }
    }

    // ===========================
    // Address Directory: Edge Cases
    // ===========================

    function test_addAndReplace_multipleTimes() public {
        string memory testName = "testAddress";

        vm.startPrank(owner);
        directory.addAddress(testName, address(0x1));
        assertEq(directory.getAddress(testName), address(0x1), "First address should be set");

        directory.replaceAddress(testName, address(0x2));
        assertEq(directory.getAddress(testName), address(0x2), "Second address should be set");

        directory.replaceAddress(testName, address(0x3));
        assertEq(directory.getAddress(testName), address(0x3), "Third address should be set");
        vm.stopPrank();
    }

    function test_removeAndReAdd() public {
        address testAddress = address(0x1234);
        string memory testName = "testAddress";

        vm.startPrank(owner);
        directory.addAddress(testName, testAddress);
        directory.removeAddress(testName);

        // Should be able to add again after removal
        directory.addAddress(testName, testAddress);
        assertEq(directory.getAddress(testName), testAddress, "Should be able to re-add after removal");
        vm.stopPrank();
    }
}
