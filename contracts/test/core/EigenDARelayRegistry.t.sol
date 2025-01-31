// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {Test} from "forge-std/Test.sol";
import {EigenDARelayRegistry} from "../../../src/core/EigenDARelayRegistry.sol";
import {IEigenDAStructs} from "../../../src/interfaces/IEigenDAStructs.sol";

contract EigenDARelayRegistryTest is Test {
    EigenDARelayRegistry public registry;
    address public owner;
    address public user;
    address public relayAddress;
    string public relayURL;

    function setUp() public {
        owner = address(this);
        user = address(0x1);
        relayAddress = address(0x2);
        relayURL = "https://relay.example.com";

        registry = new EigenDARelayRegistry();
        registry.initialize(owner);
    }

    function test_Initialize() public {
        assertEq(registry.owner(), owner);
    }

    function test_AddRelayInfo() public {
        IEigenDAStructs.RelayInfo memory relayInfo = IEigenDAStructs.RelayInfo({
            relayAddress: relayAddress,
            relayURL: relayURL
        });

        // Test successful registration with correct deposit
        uint32 relayKey = registry.addRelayInfo{value: 10 ether}(relayInfo);
        assertEq(relayKey, 0);
        assertEq(registry.relayKeyToAddress(relayKey), relayAddress);
        assertEq(registry.relayKeyToUrl(relayKey), relayURL);
        assertEq(registry.relayRegistrants(relayKey), address(this));
        assertEq(registry.relayDeposits(relayKey), 10 ether);
        assertTrue(registry.isOwnerCreatedRelay(relayKey));
    }

    function test_AddRelayInfo_IncorrectDeposit() public {
        IEigenDAStructs.RelayInfo memory relayInfo = IEigenDAStructs.RelayInfo({
            relayAddress: relayAddress,
            relayURL: relayURL
        });

        // Test registration with incorrect deposit
        vm.expectRevert("Deposit of 10 ether required for registration");
        registry.addRelayInfo{value: 5 ether}(relayInfo);
    }

    function test_DeregisterRelay() public {
        // Register a relay
        IEigenDAStructs.RelayInfo memory relayInfo = IEigenDAStructs.RelayInfo({
            relayAddress: relayAddress,
            relayURL: relayURL
        });
        uint32 relayKey = registry.addRelayInfo{value: 10 ether}(relayInfo);

        // Record initial balance
        uint256 initialBalance = address(this).balance;

        // Deregister the relay
        registry.deregisterRelay(relayKey);

        // Verify state is cleared
        assertEq(registry.relayRegistrants(relayKey), address(0));
        assertEq(registry.relayDeposits(relayKey), 0);
        assertEq(registry.relayKeyToAddress(relayKey), address(0));
        assertEq(registry.relayKeyToUrl(relayKey), "");
        assertFalse(registry.isOwnerCreatedRelay(relayKey));

        // Verify refund (10% tax)
        assertEq(address(this).balance, initialBalance + 9 ether);
    }

    function test_DeregisterRelay_NotRegistrant() public {
        // Register a relay as owner
        IEigenDAStructs.RelayInfo memory relayInfo = IEigenDAStructs.RelayInfo({
            relayAddress: relayAddress,
            relayURL: relayURL
        });
        uint32 relayKey = registry.addRelayInfo{value: 10 ether}(relayInfo);

        // Try to deregister as non-registrant
        vm.prank(user);
        vm.expectRevert("Caller is not the registrant");
        registry.deregisterRelay(relayKey);
    }

    function test_DeregisterRelay_NoDeposit() public {
        vm.expectRevert("No deposit found");
        registry.deregisterRelay(0);
    }

    // Helper function to receive ETH
    receive() external payable {}
} 