// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {Test} from "lib/forge-std/src/Test.sol";
import {EigenDADirectory} from "src/core/EigenDADirectory.sol";
import {ConfigRegistryTypes} from "src/core/libraries/v3/config-registry/ConfigRegistryTypes.sol";
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

    string constant CONFIG_NAME_BLOCKNUMBER = "testConfigBlockNumber";
    string constant CONFIG_NAME_TIMESTAMP = "testConfigTimestamp";

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
        accessControl = new EigenDAAccessControl(owner);

        // Deploy and initialize DA Directory
        directory = new EigenDADirectory();

        vm.expectEmit(true, true, true, true);
        emit IEigenDAAddressDirectory.AddressAdded(
            AddressDirectoryConstants.ACCESS_CONTROL_NAME,
            keccak256(abi.encodePacked(AddressDirectoryConstants.ACCESS_CONTROL_NAME)),
            address(accessControl)
        );

        // Verify event and genesis state
        directory.initialize(address(accessControl));
    }

    function test_initialize_revertAlreadyInitialized() public {
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
        assertEq(directory.getAllNames().length, 2, "Two named entries should exist");

        vm.expectEmit(true, true, true, true);
        emit IEigenDAAddressDirectory.AddressReplaced(
            testNamedKey, keccak256(abi.encodePacked(testNamedKey)), oldAddress, newAddress
        );
        directory.replaceAddress(testNamedKey, newAddress);
        vm.stopPrank();
        assertEq(directory.getAllNames().length, 2, "Two named entries should still exist");
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
        assertEq(directory.getAllNames().length, 2);

        vm.expectEmit(true, true, true, true);
        emit IEigenDAAddressDirectory.AddressRemoved(testNamedKey, keccak256(abi.encodePacked(testNamedKey)));
        directory.removeAddress(testNamedKey);
        vm.stopPrank();

        assertEq(directory.getAllNames().length, 1);
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

    // ===========================
    // Config Registry: BlockNumber Config Tests
    // ===========================

    function test_getActiveAndFutureBlockNumberConfigs_emptyCheckpoints() public view {
        ConfigRegistryTypes.BlockNumberCheckpoint[] memory results =
            directory.getActiveAndFutureBlockNumberConfigs(CONFIG_NAME_BLOCKNUMBER, 100);

        assertEq(results.length, 0, "Should return empty array when no checkpoints exist");
    }

    function test_getActiveAndFutureBlockNumberConfigs_singleCheckpoint_beforeActivation() public {
        // Add a checkpoint at activation block 100
        vm.prank(owner);
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 100, bytes("value1"));

        // Query with activation block before the checkpoint
        ConfigRegistryTypes.BlockNumberCheckpoint[] memory results =
            directory.getActiveAndFutureBlockNumberConfigs(CONFIG_NAME_BLOCKNUMBER, 50);

        assertEq(results.length, 0, "Should return empty array when querying before first checkpoint");
    }

    function test_getActiveAndFutureBlockNumberConfigs_singleCheckpoint_atActivation() public {
        // Add a checkpoint at activation block 100
        vm.prank(owner);
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 100, bytes("value1"));

        // Query with activation block equal to the checkpoint
        ConfigRegistryTypes.BlockNumberCheckpoint[] memory results =
            directory.getActiveAndFutureBlockNumberConfigs(CONFIG_NAME_BLOCKNUMBER, 100);

        assertEq(results.length, 1, "Should return 1 checkpoint");
        assertEq(results[0].activationBlock, 100, "Should return checkpoint at block 100");
        assertEq(keccak256(results[0].value), keccak256(bytes("value1")), "Should return correct value");
    }

    function test_getActiveAndFutureBlockNumberConfigs_singleCheckpoint_afterActivation() public {
        // Add a checkpoint at activation block 100
        vm.prank(owner);
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 100, bytes("value1"));

        // Query with activation block after the checkpoint
        ConfigRegistryTypes.BlockNumberCheckpoint[] memory results =
            directory.getActiveAndFutureBlockNumberConfigs(CONFIG_NAME_BLOCKNUMBER, 150);

        assertEq(results.length, 1, "Should return 1 checkpoint (the active one)");
        assertEq(results[0].activationBlock, 100, "Should return checkpoint at block 100");
        assertEq(keccak256(results[0].value), keccak256(bytes("value1")), "Should return correct value");
    }

    function test_getActiveAndFutureBlockNumberConfigs_multipleCheckpoints_beforeAll() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 100, bytes("value1"));
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 200, bytes("value2"));
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 300, bytes("value3"));
        vm.stopPrank();

        // Query before all checkpoints
        ConfigRegistryTypes.BlockNumberCheckpoint[] memory results =
            directory.getActiveAndFutureBlockNumberConfigs(CONFIG_NAME_BLOCKNUMBER, 50);

        assertEq(results.length, 0, "Should return empty array when querying before all checkpoints");
    }

    function test_getActiveAndFutureBlockNumberConfigs_multipleCheckpoints_betweenCheckpoints() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 100, bytes("value1"));
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 200, bytes("value2"));
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 300, bytes("value3"));
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 400, bytes("value4"));
        vm.stopPrank();

        // Query at activation block 150 (between 100 and 200)
        ConfigRegistryTypes.BlockNumberCheckpoint[] memory results =
            directory.getActiveAndFutureBlockNumberConfigs(CONFIG_NAME_BLOCKNUMBER, 150);

        assertEq(results.length, 4, "Should return current + all future checkpoints");
        assertEq(results[0].activationBlock, 100, "First result should be currently active checkpoint");
        assertEq(keccak256(results[0].value), keccak256(bytes("value1")), "Should have correct value");
        assertEq(results[1].activationBlock, 200, "Second result should be next checkpoint");
        assertEq(keccak256(results[1].value), keccak256(bytes("value2")), "Should have correct value");
        assertEq(results[2].activationBlock, 300, "Third result should be next checkpoint");
        assertEq(keccak256(results[2].value), keccak256(bytes("value3")), "Should have correct value");
        assertEq(results[3].activationBlock, 400, "Fourth result should be next checkpoint");
        assertEq(keccak256(results[3].value), keccak256(bytes("value4")), "Should have correct value");
    }

    function test_getActiveAndFutureBlockNumberConfigs_multipleCheckpoints_atCheckpoint() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 100, bytes("value1"));
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 200, bytes("value2"));
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 300, bytes("value3"));
        vm.stopPrank();

        // Query at exact activation block 200
        ConfigRegistryTypes.BlockNumberCheckpoint[] memory results =
            directory.getActiveAndFutureBlockNumberConfigs(CONFIG_NAME_BLOCKNUMBER, 200);

        assertEq(results.length, 2, "Should return checkpoint at 200 and all future");
        assertEq(results[0].activationBlock, 200, "First result should be checkpoint at 200");
        assertEq(keccak256(results[0].value), keccak256(bytes("value2")), "Should have correct value");
        assertEq(results[1].activationBlock, 300, "Second result should be next checkpoint");
        assertEq(keccak256(results[1].value), keccak256(bytes("value3")), "Should have correct value");
    }

    function test_getActiveAndFutureBlockNumberConfigs_multipleCheckpoints_afterAll() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 100, bytes("value1"));
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 200, bytes("value2"));
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 300, bytes("value3"));
        vm.stopPrank();

        // Query after all checkpoints
        ConfigRegistryTypes.BlockNumberCheckpoint[] memory results =
            directory.getActiveAndFutureBlockNumberConfigs(CONFIG_NAME_BLOCKNUMBER, 500);

        assertEq(results.length, 1, "Should return only the last (currently active) checkpoint");
        assertEq(results[0].activationBlock, 300, "Should return last checkpoint");
        assertEq(keccak256(results[0].value), keccak256(bytes("value3")), "Should have correct value");
    }

    function test_getActiveAndFutureBlockNumberConfigs_manyCheckpoints() public {
        // Add 10 checkpoints
        vm.startPrank(owner);
        for (uint256 i = 1; i <= 10; i++) {
            directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, i * 100, abi.encode(i));
        }
        vm.stopPrank();

        // Query at 550 (between checkpoint 5 and 6)
        ConfigRegistryTypes.BlockNumberCheckpoint[] memory results =
            directory.getActiveAndFutureBlockNumberConfigs(CONFIG_NAME_BLOCKNUMBER, 550);

        assertEq(results.length, 6, "Should return checkpoint 5 through 10");
        assertEq(results[0].activationBlock, 500, "First should be currently active (checkpoint 5)");
        assertEq(keccak256(results[0].value), keccak256(abi.encode(5)), "Should have correct value");
        assertEq(results[5].activationBlock, 1000, "Last should be checkpoint 10");
        assertEq(keccak256(results[5].value), keccak256(abi.encode(10)), "Should have correct value");
    }

    // ===========================
    // Config Registry: Timestamp Config Tests
    // ===========================

    function test_getActiveAndFutureTimestampConfigs_emptyCheckpoints() public view {
        ConfigRegistryTypes.TimeStampCheckpoint[] memory results =
            directory.getActiveAndFutureTimestampConfigs(CONFIG_NAME_TIMESTAMP, 100);

        assertEq(results.length, 0, "Should return empty array when no checkpoints exist");
    }

    function test_getActiveAndFutureTimestampConfigs_singleCheckpoint_beforeActivation() public {
        // Add a checkpoint at activation timestamp 100
        vm.prank(owner);
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 100, bytes("value1"));

        // Query with activation timestamp before the checkpoint
        ConfigRegistryTypes.TimeStampCheckpoint[] memory results =
            directory.getActiveAndFutureTimestampConfigs(CONFIG_NAME_TIMESTAMP, 50);

        assertEq(results.length, 0, "Should return empty array when querying before first checkpoint");
    }

    function test_getActiveAndFutureTimestampConfigs_singleCheckpoint_atActivation() public {
        // Add a checkpoint at activation timestamp 100
        vm.prank(owner);
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 100, bytes("value1"));

        // Query with activation timestamp equal to the checkpoint
        ConfigRegistryTypes.TimeStampCheckpoint[] memory results =
            directory.getActiveAndFutureTimestampConfigs(CONFIG_NAME_TIMESTAMP, 100);

        assertEq(results.length, 1, "Should return 1 checkpoint");
        assertEq(results[0].activationTime, 100, "Should return checkpoint at timestamp 100");
        assertEq(keccak256(results[0].value), keccak256(bytes("value1")), "Should return correct value");
    }

    function test_getActiveAndFutureTimestampConfigs_singleCheckpoint_afterActivation() public {
        // Add a checkpoint at activation timestamp 100
        vm.prank(owner);
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 100, bytes("value1"));

        // Query with activation timestamp after the checkpoint
        ConfigRegistryTypes.TimeStampCheckpoint[] memory results =
            directory.getActiveAndFutureTimestampConfigs(CONFIG_NAME_TIMESTAMP, 150);

        assertEq(results.length, 1, "Should return 1 checkpoint (the active one)");
        assertEq(results[0].activationTime, 100, "Should return checkpoint at timestamp 100");
        assertEq(keccak256(results[0].value), keccak256(bytes("value1")), "Should return correct value");
    }

    function test_getActiveAndFutureTimestampConfigs_multipleCheckpoints_betweenCheckpoints() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 100, bytes("value1"));
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 200, bytes("value2"));
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 300, bytes("value3"));
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 400, bytes("value4"));
        vm.stopPrank();

        // Query at activation timestamp 150 (between 100 and 200)
        ConfigRegistryTypes.TimeStampCheckpoint[] memory results =
            directory.getActiveAndFutureTimestampConfigs(CONFIG_NAME_TIMESTAMP, 150);

        assertEq(results.length, 4, "Should return current + all future checkpoints");
        assertEq(results[0].activationTime, 100, "First result should be currently active checkpoint");
        assertEq(keccak256(results[0].value), keccak256(bytes("value1")), "Should have correct value");
        assertEq(results[1].activationTime, 200, "Second result should be next checkpoint");
        assertEq(keccak256(results[1].value), keccak256(bytes("value2")), "Should have correct value");
        assertEq(results[2].activationTime, 300, "Third result should be next checkpoint");
        assertEq(keccak256(results[2].value), keccak256(bytes("value3")), "Should have correct value");
        assertEq(results[3].activationTime, 400, "Fourth result should be next checkpoint");
        assertEq(keccak256(results[3].value), keccak256(bytes("value4")), "Should have correct value");
    }

    function test_getActiveAndFutureTimestampConfigs_multipleCheckpoints_atCheckpoint() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 100, bytes("value1"));
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 200, bytes("value2"));
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 300, bytes("value3"));
        vm.stopPrank();

        // Query at exact activation timestamp 200
        ConfigRegistryTypes.TimeStampCheckpoint[] memory results =
            directory.getActiveAndFutureTimestampConfigs(CONFIG_NAME_TIMESTAMP, 200);

        assertEq(results.length, 2, "Should return checkpoint at 200 and all future");
        assertEq(results[0].activationTime, 200, "First result should be checkpoint at 200");
        assertEq(keccak256(results[0].value), keccak256(bytes("value2")), "Should have correct value");
        assertEq(results[1].activationTime, 300, "Second result should be next checkpoint");
        assertEq(keccak256(results[1].value), keccak256(bytes("value3")), "Should have correct value");
    }

    function test_getActiveAndFutureTimestampConfigs_multipleCheckpoints_afterAll() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 100, bytes("value1"));
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 200, bytes("value2"));
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 300, bytes("value3"));
        vm.stopPrank();

        // Query after all checkpoints
        ConfigRegistryTypes.TimeStampCheckpoint[] memory results =
            directory.getActiveAndFutureTimestampConfigs(CONFIG_NAME_TIMESTAMP, 500);

        assertEq(results.length, 1, "Should return only the last (currently active) checkpoint");
        assertEq(results[0].activationTime, 300, "Should return last checkpoint");
        assertEq(keccak256(results[0].value), keccak256(bytes("value3")), "Should have correct value");
    }

    function test_getActiveAndFutureTimestampConfigs_variableLengthData() public {
        // Add checkpoints with different length data
        vm.startPrank(owner);
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 100, hex"010203");
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 200, hex"0102030405060708");
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 300, hex"01");
        vm.stopPrank();

        // Query at 150
        ConfigRegistryTypes.TimeStampCheckpoint[] memory results =
            directory.getActiveAndFutureTimestampConfigs(CONFIG_NAME_TIMESTAMP, 150);

        assertEq(results.length, 3, "Should return all from checkpoint 1 onwards");
        assertEq(keccak256(results[0].value), keccak256(hex"010203"), "Should handle 3-byte value");
        assertEq(keccak256(results[1].value), keccak256(hex"0102030405060708"), "Should handle 8-byte value");
        assertEq(keccak256(results[2].value), keccak256(hex"01"), "Should handle 1-byte value");
    }

    // ===========================
    // Config Registry: Edge Cases and Boundary Tests
    // ===========================

    function test_getActiveAndFutureBlockNumberConfigs_boundaryValues() public {
        // Add checkpoints at boundary values
        vm.startPrank(owner);
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, block.number, bytes("value1"));
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, type(uint256).max, bytes("value2"));
        vm.stopPrank();

        // Query at 0
        ConfigRegistryTypes.BlockNumberCheckpoint[] memory results =
            directory.getActiveAndFutureBlockNumberConfigs(CONFIG_NAME_BLOCKNUMBER, block.number);

        assertEq(results.length, 2, "Should return both checkpoints");
        assertEq(results[0].activationBlock, block.number, "Should include checkpoint at block.number");
        assertEq(results[1].activationBlock, type(uint256).max, "Should include checkpoint at max");
    }

    function test_separateConfigs_doNotInterfere() public {
        // Add checkpoints to both BlockNumber and Timestamp configs
        vm.startPrank(owner);
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 100, bytes("blockValue1"));
        directory.addConfigBlockNumber(CONFIG_NAME_BLOCKNUMBER, 200, bytes("blockValue2"));
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 100, hex"aa");
        directory.addConfigTimeStamp(CONFIG_NAME_TIMESTAMP, 200, hex"bb");
        vm.stopPrank();

        // Query both
        ConfigRegistryTypes.BlockNumberCheckpoint[] memory resultsBlock =
            directory.getActiveAndFutureBlockNumberConfigs(CONFIG_NAME_BLOCKNUMBER, 150);
        ConfigRegistryTypes.TimeStampCheckpoint[] memory resultsTimestamp =
            directory.getActiveAndFutureTimestampConfigs(CONFIG_NAME_TIMESTAMP, 150);

        // Verify they don't interfere with each other
        assertEq(resultsBlock.length, 2, "BlockNumber should have 2 checkpoints");
        assertEq(resultsTimestamp.length, 2, "Timestamp should have 2 checkpoints");
        assertEq(
            keccak256(resultsBlock[0].value), keccak256(bytes("blockValue1")), "BlockNumber values should be correct"
        );
        assertEq(keccak256(resultsTimestamp[0].value), keccak256(hex"aa"), "Timestamp values should be correct");
    }
}
