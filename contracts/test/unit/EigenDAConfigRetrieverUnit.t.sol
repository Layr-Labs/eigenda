// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {Test} from "lib/forge-std/src/Test.sol";
import {EigenDAConfigRetriever} from "src/periphery/cfg-getter/EigenDAConfigRetriever.sol";
import {EigenDADirectory} from "src/core/EigenDADirectory.sol";
import {ConfigRegistryTypes} from "src/core/libraries/v3/config-registry/ConfigRegistryTypes.sol";
import {EigenDAAccessControl} from "src/core/EigenDAAccessControl.sol";

contract EigenDAConfigRetrieverUnit is Test {
    EigenDAConfigRetriever public configRetriever;
    EigenDADirectory public directory;
    EigenDAAccessControl public accessControl;

    address owner = address(uint160(uint256(keccak256(abi.encodePacked("owner")))));
    address nonOwner = address(uint160(uint256(keccak256(abi.encodePacked("nonOwner")))));

    string constant CONFIG_NAME_BYTES32 = "testConfig";
    string constant CONFIG_NAME_BYTES = "testConfigBytes";

    function setUp() public {
        // Deploy access control with owner
        accessControl = new EigenDAAccessControl(owner);

        // Deploy directory
        directory = new EigenDADirectory();
        directory.initialize(address(accessControl));

        // Deploy config retriever
        configRetriever = new EigenDAConfigRetriever(address(directory));
    }

    // ===========================
    // Bytes32 Config Tests
    // ===========================

    function test_getActiveAndFutureBytes32Configs_emptyCheckpoints() public view {
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results =
            configRetriever.getActiveAndFutureBytes32Configs(CONFIG_NAME_BYTES32, 100);

        assertEq(results.length, 0, "Should return empty array when no checkpoints exist");
    }

    function test_getActiveAndFutureBytes32Configs_singleCheckpoint_beforeActivation() public {
        // Add a checkpoint at activation key 100
        vm.prank(owner);
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 100, bytes32(uint256(1)));

        // Query with activation key before the checkpoint
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results =
            configRetriever.getActiveAndFutureBytes32Configs(CONFIG_NAME_BYTES32, 50);

        assertEq(results.length, 0, "Should return empty array when querying before first checkpoint");
    }

    function test_getActiveAndFutureBytes32Configs_singleCheckpoint_atActivation() public {
        // Add a checkpoint at activation key 100
        vm.prank(owner);
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 100, bytes32(uint256(1)));

        // Query with activation key equal to the checkpoint
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results =
            configRetriever.getActiveAndFutureBytes32Configs(CONFIG_NAME_BYTES32, 100);

        assertEq(results.length, 1, "Should return 1 checkpoint");
        assertEq(results[0].activationKey, 100, "Should return checkpoint at key 100");
        assertEq(results[0].value, bytes32(uint256(1)), "Should return correct value");
    }

    function test_getActiveAndFutureBytes32Configs_singleCheckpoint_afterActivation() public {
        // Add a checkpoint at activation key 100
        vm.prank(owner);
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 100, bytes32(uint256(1)));

        // Query with activation key after the checkpoint
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results =
            configRetriever.getActiveAndFutureBytes32Configs(CONFIG_NAME_BYTES32, 150);

        assertEq(results.length, 1, "Should return 1 checkpoint (the active one)");
        assertEq(results[0].activationKey, 100, "Should return checkpoint at key 100");
        assertEq(results[0].value, bytes32(uint256(1)), "Should return correct value");
    }

    function test_getActiveAndFutureBytes32Configs_multipleCheckpoints_beforeAll() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 100, bytes32(uint256(1)));
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 200, bytes32(uint256(2)));
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 300, bytes32(uint256(3)));
        vm.stopPrank();

        // Query before all checkpoints
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results =
            configRetriever.getActiveAndFutureBytes32Configs(CONFIG_NAME_BYTES32, 50);

        assertEq(results.length, 0, "Should return empty array when querying before all checkpoints");
    }

    function test_getActiveAndFutureBytes32Configs_multipleCheckpoints_betweenCheckpoints() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 100, bytes32(uint256(1)));
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 200, bytes32(uint256(2)));
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 300, bytes32(uint256(3)));
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 400, bytes32(uint256(4)));
        vm.stopPrank();

        // Query at activation key 150 (between 100 and 200)
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results =
            configRetriever.getActiveAndFutureBytes32Configs(CONFIG_NAME_BYTES32, 150);

        assertEq(results.length, 4, "Should return current + all future checkpoints");
        assertEq(results[0].activationKey, 100, "First result should be currently active checkpoint");
        assertEq(results[0].value, bytes32(uint256(1)), "Should have correct value");
        assertEq(results[1].activationKey, 200, "Second result should be next checkpoint");
        assertEq(results[1].value, bytes32(uint256(2)), "Should have correct value");
        assertEq(results[2].activationKey, 300, "Third result should be next checkpoint");
        assertEq(results[2].value, bytes32(uint256(3)), "Should have correct value");
        assertEq(results[3].activationKey, 400, "Fourth result should be next checkpoint");
        assertEq(results[3].value, bytes32(uint256(4)), "Should have correct value");
    }

    function test_getActiveAndFutureBytes32Configs_multipleCheckpoints_atCheckpoint() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 100, bytes32(uint256(1)));
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 200, bytes32(uint256(2)));
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 300, bytes32(uint256(3)));
        vm.stopPrank();

        // Query at exact activation key 200
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results =
            configRetriever.getActiveAndFutureBytes32Configs(CONFIG_NAME_BYTES32, 200);

        assertEq(results.length, 2, "Should return checkpoint at 200 and all future");
        assertEq(results[0].activationKey, 200, "First result should be checkpoint at 200");
        assertEq(results[0].value, bytes32(uint256(2)), "Should have correct value");
        assertEq(results[1].activationKey, 300, "Second result should be next checkpoint");
        assertEq(results[1].value, bytes32(uint256(3)), "Should have correct value");
    }

    function test_getActiveAndFutureBytes32Configs_multipleCheckpoints_afterAll() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 100, bytes32(uint256(1)));
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 200, bytes32(uint256(2)));
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 300, bytes32(uint256(3)));
        vm.stopPrank();

        // Query after all checkpoints
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results =
            configRetriever.getActiveAndFutureBytes32Configs(CONFIG_NAME_BYTES32, 500);

        assertEq(results.length, 1, "Should return only the last (currently active) checkpoint");
        assertEq(results[0].activationKey, 300, "Should return last checkpoint");
        assertEq(results[0].value, bytes32(uint256(3)), "Should have correct value");
    }

    function test_getActiveAndFutureBytes32Configs_manyCheckpoints() public {
        // Add 10 checkpoints
        vm.startPrank(owner);
        for (uint256 i = 1; i <= 10; i++) {
            directory.addConfigBytes32(CONFIG_NAME_BYTES32, i * 100, bytes32(i));
        }
        vm.stopPrank();

        // Query at 550 (between checkpoint 5 and 6)
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results =
            configRetriever.getActiveAndFutureBytes32Configs(CONFIG_NAME_BYTES32, 550);

        assertEq(results.length, 6, "Should return checkpoint 5 through 10");
        assertEq(results[0].activationKey, 500, "First should be currently active (checkpoint 5)");
        assertEq(results[0].value, bytes32(uint256(5)), "Should have correct value");
        assertEq(results[5].activationKey, 1000, "Last should be checkpoint 10");
        assertEq(results[5].value, bytes32(uint256(10)), "Should have correct value");
    }

    // ===========================
    // Bytes Config Tests
    // ===========================

    function test_getActiveAndFutureBytesConfigs_emptyCheckpoints() public view {
        ConfigRegistryTypes.BytesCheckpoint[] memory results =
            configRetriever.getActiveAndFutureBytesConfigs(CONFIG_NAME_BYTES, 100);

        assertEq(results.length, 0, "Should return empty array when no checkpoints exist");
    }

    function test_getActiveAndFutureBytesConfigs_singleCheckpoint_beforeActivation() public {
        // Add a checkpoint at activation key 100
        vm.prank(owner);
        directory.addConfigBytes(CONFIG_NAME_BYTES, 100, hex"01");

        // Query with activation key before the checkpoint
        ConfigRegistryTypes.BytesCheckpoint[] memory results =
            configRetriever.getActiveAndFutureBytesConfigs(CONFIG_NAME_BYTES, 50);

        assertEq(results.length, 0, "Should return empty array when querying before first checkpoint");
    }

    function test_getActiveAndFutureBytesConfigs_singleCheckpoint_atActivation() public {
        // Add a checkpoint at activation key 100
        vm.prank(owner);
        directory.addConfigBytes(CONFIG_NAME_BYTES, 100, hex"01");

        // Query with activation key equal to the checkpoint
        ConfigRegistryTypes.BytesCheckpoint[] memory results =
            configRetriever.getActiveAndFutureBytesConfigs(CONFIG_NAME_BYTES, 100);

        assertEq(results.length, 1, "Should return 1 checkpoint");
        assertEq(results[0].activationKey, 100, "Should return checkpoint at key 100");
        assertEq(keccak256(results[0].value), keccak256(hex"01"), "Should return correct value");
    }

    function test_getActiveAndFutureBytesConfigs_singleCheckpoint_afterActivation() public {
        // Add a checkpoint at activation key 100
        vm.prank(owner);
        directory.addConfigBytes(CONFIG_NAME_BYTES, 100, hex"01");

        // Query with activation key after the checkpoint
        ConfigRegistryTypes.BytesCheckpoint[] memory results =
            configRetriever.getActiveAndFutureBytesConfigs(CONFIG_NAME_BYTES, 150);

        assertEq(results.length, 1, "Should return 1 checkpoint (the active one)");
        assertEq(results[0].activationKey, 100, "Should return checkpoint at key 100");
        assertEq(keccak256(results[0].value), keccak256(hex"01"), "Should return correct value");
    }

    function test_getActiveAndFutureBytesConfigs_multipleCheckpoints_betweenCheckpoints() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigBytes(CONFIG_NAME_BYTES, 100, hex"01");
        directory.addConfigBytes(CONFIG_NAME_BYTES, 200, hex"02");
        directory.addConfigBytes(CONFIG_NAME_BYTES, 300, hex"03");
        directory.addConfigBytes(CONFIG_NAME_BYTES, 400, hex"04");
        vm.stopPrank();

        // Query at activation key 150 (between 100 and 200)
        ConfigRegistryTypes.BytesCheckpoint[] memory results =
            configRetriever.getActiveAndFutureBytesConfigs(CONFIG_NAME_BYTES, 150);

        assertEq(results.length, 4, "Should return current + all future checkpoints");
        assertEq(results[0].activationKey, 100, "First result should be currently active checkpoint");
        assertEq(keccak256(results[0].value), keccak256(hex"01"), "Should have correct value");
        assertEq(results[1].activationKey, 200, "Second result should be next checkpoint");
        assertEq(keccak256(results[1].value), keccak256(hex"02"), "Should have correct value");
        assertEq(results[2].activationKey, 300, "Third result should be next checkpoint");
        assertEq(keccak256(results[2].value), keccak256(hex"03"), "Should have correct value");
        assertEq(results[3].activationKey, 400, "Fourth result should be next checkpoint");
        assertEq(keccak256(results[3].value), keccak256(hex"04"), "Should have correct value");
    }

    function test_getActiveAndFutureBytesConfigs_multipleCheckpoints_atCheckpoint() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigBytes(CONFIG_NAME_BYTES, 100, hex"01");
        directory.addConfigBytes(CONFIG_NAME_BYTES, 200, hex"02");
        directory.addConfigBytes(CONFIG_NAME_BYTES, 300, hex"03");
        vm.stopPrank();

        // Query at exact activation key 200
        ConfigRegistryTypes.BytesCheckpoint[] memory results =
            configRetriever.getActiveAndFutureBytesConfigs(CONFIG_NAME_BYTES, 200);

        assertEq(results.length, 2, "Should return checkpoint at 200 and all future");
        assertEq(results[0].activationKey, 200, "First result should be checkpoint at 200");
        assertEq(keccak256(results[0].value), keccak256(hex"02"), "Should have correct value");
        assertEq(results[1].activationKey, 300, "Second result should be next checkpoint");
        assertEq(keccak256(results[1].value), keccak256(hex"03"), "Should have correct value");
    }

    function test_getActiveAndFutureBytesConfigs_multipleCheckpoints_afterAll() public {
        // Add multiple checkpoints
        vm.startPrank(owner);
        directory.addConfigBytes(CONFIG_NAME_BYTES, 100, hex"01");
        directory.addConfigBytes(CONFIG_NAME_BYTES, 200, hex"02");
        directory.addConfigBytes(CONFIG_NAME_BYTES, 300, hex"03");
        vm.stopPrank();

        // Query after all checkpoints
        ConfigRegistryTypes.BytesCheckpoint[] memory results =
            configRetriever.getActiveAndFutureBytesConfigs(CONFIG_NAME_BYTES, 500);

        assertEq(results.length, 1, "Should return only the last (currently active) checkpoint");
        assertEq(results[0].activationKey, 300, "Should return last checkpoint");
        assertEq(keccak256(results[0].value), keccak256(hex"03"), "Should have correct value");
    }

    function test_getActiveAndFutureBytesConfigs_variableLengthData() public {
        // Add checkpoints with different length data
        vm.startPrank(owner);
        directory.addConfigBytes(CONFIG_NAME_BYTES, 100, hex"010203");
        directory.addConfigBytes(CONFIG_NAME_BYTES, 200, hex"0102030405060708");
        directory.addConfigBytes(CONFIG_NAME_BYTES, 300, hex"01");
        vm.stopPrank();

        // Query at 150
        ConfigRegistryTypes.BytesCheckpoint[] memory results =
            configRetriever.getActiveAndFutureBytesConfigs(CONFIG_NAME_BYTES, 150);

        assertEq(results.length, 3, "Should return all from checkpoint 1 onwards");
        assertEq(keccak256(results[0].value), keccak256(hex"010203"), "Should handle 3-byte value");
        assertEq(keccak256(results[1].value), keccak256(hex"0102030405060708"), "Should handle 8-byte value");
        assertEq(keccak256(results[2].value), keccak256(hex"01"), "Should handle 1-byte value");
    }

    // ===========================
    // Edge Cases and Boundary Tests
    // ===========================

    function test_getActiveAndFutureBytes32Configs_boundaryValues() public {
        // Add checkpoints at boundary values
        vm.startPrank(owner);
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 0, bytes32(uint256(1)));
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, type(uint256).max, bytes32(uint256(2)));
        vm.stopPrank();

        // Query at 0
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results =
            configRetriever.getActiveAndFutureBytes32Configs(CONFIG_NAME_BYTES32, 0);

        assertEq(results.length, 2, "Should return both checkpoints");
        assertEq(results[0].activationKey, 0, "Should include checkpoint at 0");
        assertEq(results[1].activationKey, type(uint256).max, "Should include checkpoint at max");
    }

    function test_constructor_setsConfigRegistry() public view {
        assertEq(
            address(configRetriever.configRegistry()),
            address(directory),
            "Constructor should set configRegistry correctly"
        );
    }

    function test_separateConfigs_doNotInterfere() public {
        // Add checkpoints to both bytes32 and bytes configs
        vm.startPrank(owner);
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 100, bytes32(uint256(1)));
        directory.addConfigBytes32(CONFIG_NAME_BYTES32, 200, bytes32(uint256(2)));
        directory.addConfigBytes(CONFIG_NAME_BYTES, 100, hex"aa");
        directory.addConfigBytes(CONFIG_NAME_BYTES, 200, hex"bb");
        vm.stopPrank();

        // Query both
        ConfigRegistryTypes.Bytes32Checkpoint[] memory results32 =
            configRetriever.getActiveAndFutureBytes32Configs(CONFIG_NAME_BYTES32, 150);
        ConfigRegistryTypes.BytesCheckpoint[] memory resultsBytes =
            configRetriever.getActiveAndFutureBytesConfigs(CONFIG_NAME_BYTES, 150);

        // Verify they don't interfere with each other
        assertEq(results32.length, 2, "Bytes32 should have 2 checkpoints");
        assertEq(resultsBytes.length, 2, "Bytes should have 2 checkpoints");
        assertEq(results32[0].value, bytes32(uint256(1)), "Bytes32 values should be correct");
        assertEq(keccak256(resultsBytes[0].value), keccak256(hex"aa"), "Bytes values should be correct");
    }
}

