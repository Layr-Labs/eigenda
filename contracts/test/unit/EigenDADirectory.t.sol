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
        assertEq(names.length, 2, "Should have 2 names after removal (address1 + address3)");

        // Verify address2 is not present
        for (uint256 i = 0; i < names.length; i++) {
            assertTrue(
                keccak256(bytes(names[i])) != keccak256(bytes("address2")), "address2 should not be in the list"
            );
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

    // ===========================
    // Config Registry: Basic Operations - BlockNumber
    // ===========================

    function test_addConfigBlockNumber_success() public {
        string memory configName = "testConfig";
        uint256 activationBlock = block.number + 100;
        bytes memory configValue = bytes("testValue");

        vm.prank(owner);
        directory.addConfigBlockNumber(configName, activationBlock, configValue);

        bytes32 nameDigest = keccak256(abi.encodePacked(configName));
        assertEq(directory.getNumCheckpointsBlockNumber(nameDigest), 1, "Should have 1 checkpoint");
        assertEq(
            keccak256(directory.getConfigBlockNumber(nameDigest, 0)),
            keccak256(configValue),
            "Config value should match"
        );
        assertEq(directory.getActivationBlockNumber(nameDigest, 0), activationBlock, "Activation block should match");
    }

    function test_addConfigBlockNumber_revertNonOwner() public {
        string memory configName = "testConfig";
        uint256 activationBlock = block.number + 100;
        bytes memory configValue = bytes("testValue");

        vm.prank(nonOwner);
        vm.expectRevert("Caller is not the owner");
        directory.addConfigBlockNumber(configName, activationBlock, configValue);
    }

    function test_addConfigBlockNumber_multipleCheckpoints() public {
        string memory configName = "testConfig";
        bytes32 nameDigest = keccak256(abi.encodePacked(configName));

        vm.startPrank(owner);
        directory.addConfigBlockNumber(configName, block.number + 100, bytes("value1"));
        directory.addConfigBlockNumber(configName, block.number + 200, bytes("value2"));
        directory.addConfigBlockNumber(configName, block.number + 300, bytes("value3"));
        vm.stopPrank();

        assertEq(directory.getNumCheckpointsBlockNumber(nameDigest), 3, "Should have 3 checkpoints");
        assertEq(
            keccak256(directory.getConfigBlockNumber(nameDigest, 0)), keccak256(bytes("value1")), "First value correct"
        );
        assertEq(
            keccak256(directory.getConfigBlockNumber(nameDigest, 1)), keccak256(bytes("value2")), "Second value correct"
        );
        assertEq(
            keccak256(directory.getConfigBlockNumber(nameDigest, 2)), keccak256(bytes("value3")), "Third value correct"
        );
    }

    // ===========================
    // Config Registry: Basic Operations - Timestamp
    // ===========================

    function test_addConfigTimeStamp_success() public {
        string memory configName = "testConfig";
        uint256 activationTime = block.timestamp + 100;
        bytes memory configValue = bytes("testValue");

        vm.prank(owner);
        directory.addConfigTimeStamp(configName, activationTime, configValue);

        bytes32 nameDigest = keccak256(abi.encodePacked(configName));
        assertEq(directory.getNumCheckpointsTimeStamp(nameDigest), 1, "Should have 1 checkpoint");
        assertEq(
            keccak256(directory.getConfigTimeStamp(nameDigest, 0)), keccak256(configValue), "Config value should match"
        );
        assertEq(directory.getActivationTimeStamp(nameDigest, 0), activationTime, "Activation time should match");
    }

    function test_addConfigTimeStamp_revertNonOwner() public {
        string memory configName = "testConfig";
        uint256 activationTime = block.timestamp + 100;
        bytes memory configValue = bytes("testValue");

        vm.prank(nonOwner);
        vm.expectRevert("Caller is not the owner");
        directory.addConfigTimeStamp(configName, activationTime, configValue);
    }

    function test_addConfigTimeStamp_multipleCheckpoints() public {
        string memory configName = "testConfig";
        bytes32 nameDigest = keccak256(abi.encodePacked(configName));

        vm.startPrank(owner);
        directory.addConfigTimeStamp(configName, block.timestamp + 100, bytes("value1"));
        directory.addConfigTimeStamp(configName, block.timestamp + 200, bytes("value2"));
        directory.addConfigTimeStamp(configName, block.timestamp + 300, bytes("value3"));
        vm.stopPrank();

        assertEq(directory.getNumCheckpointsTimeStamp(nameDigest), 3, "Should have 3 checkpoints");
        assertEq(
            keccak256(directory.getConfigTimeStamp(nameDigest, 0)), keccak256(bytes("value1")), "First value correct"
        );
        assertEq(
            keccak256(directory.getConfigTimeStamp(nameDigest, 1)), keccak256(bytes("value2")), "Second value correct"
        );
        assertEq(
            keccak256(directory.getConfigTimeStamp(nameDigest, 2)), keccak256(bytes("value3")), "Third value correct"
        );
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
        assertEq(keccak256(resultsBlock[0].value), keccak256(bytes("blockValue1")), "BlockNumber values should be correct");
        assertEq(keccak256(resultsTimestamp[0].value), keccak256(hex"aa"), "Timestamp values should be correct");
    }
}
