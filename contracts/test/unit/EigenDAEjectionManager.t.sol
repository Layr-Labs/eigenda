// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

import {EigenDAEjectionManager} from "src/periphery/ejection/EigenDAEjectionManager.sol";
import {EigenDAEjectionLib} from "src/periphery/ejection/libraries/EigenDAEjectionLib.sol";

import {AccessControlConstants} from "src/core/libraries/v3/access-control/AccessControlConstants.sol";

import {MockEigenDADeployer} from "test/MockEigenDADeployer.sol";

contract EigenDAEjectionManagerTest is MockEigenDADeployer {
    address testEjector;
    address ejectee;

    /// TODO: Add tests that ensure multiple ejections can be ran at once by a single ejector (1 ejector : N ejectees)
    ///       Also (N ejector : N ejectees)

    function setUp() public {
        // Deploy all mock contracts including EigenDAEjectionManager
        _deployDA();

        testEjector = makeAddr("testEjector");
        ejectee = makeAddr("ejectee");

        // Grant roles as the registryCoordinatorOwner who has DEFAULT_ADMIN_ROLE
        vm.startPrank(registryCoordinatorOwner);
        eigenDAAccessControl.grantRole(eigenDAAccessControl.DEFAULT_ADMIN_ROLE(), address(this));
        eigenDAAccessControl.grantRole(AccessControlConstants.OWNER_ROLE, address(this));
        vm.stopPrank();
    }

    function testStartEjection() public {
        testStartEjection(0, 0);
    }

    function testStartEjection(uint64 cooldown, uint64 delay) private {
        // 0) Wire up access mgmt dependencies and set protocol params on contract
        eigenDAAccessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, testEjector);
        eigenDAAccessControl.grantRole(AccessControlConstants.OWNER_ROLE, testEjector);

        vm.startPrank(testEjector);
        eigenDAEjectionManager.setCooldown(cooldown);
        eigenDAEjectionManager.setDelay(delay);

        // 1) start an ejection against an arbitrary ejectee
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionStarted(
            ejectee,
            testEjector,
            "0x", // quorums (empty for this test)
            uint64(block.timestamp),
            uint64(block.timestamp + eigenDAEjectionManager.ejectionDelay())
        );

        eigenDAEjectionManager.startEjection(ejectee, "0x");
        vm.stopPrank();

        // 2) verify that ejectee record was properly created
        assertEq(eigenDAEjectionManager.getEjector(ejectee), testEjector);
        assertEq(eigenDAEjectionManager.ejectionTime(ejectee), block.timestamp + eigenDAEjectionManager.ejectionDelay());
        assertEq(eigenDAEjectionManager.lastEjectionInitiated(ejectee), block.timestamp);
    }

    function testCancelEjectionByEjector() public {
        testCancelEjectionByEjector(0, 0);
    }

    function testCancelEjectionByEjector(uint64 cooldown, uint64 delay) private {
        // 0) grant roles
        eigenDAAccessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, testEjector);
        eigenDAAccessControl.grantRole(AccessControlConstants.OWNER_ROLE, testEjector);

        // 1) Ejector starts ejection for ejectee after setting contract params
        vm.startPrank(testEjector);
        eigenDAEjectionManager.setCooldown(cooldown);
        eigenDAEjectionManager.setDelay(delay);
        eigenDAEjectionManager.startEjection(ejectee, "0x");

        // 2) Issue a cancellation from the Ejector role
        eigenDAEjectionManager.cancelEjectionByEjector(ejectee);

        // 3) Ensure the ejectee record has been nullified
        assertEq(eigenDAEjectionManager.getEjector(ejectee), address(0));
        assertEq(eigenDAEjectionManager.ejectionTime(ejectee), 0);
        assertEq(eigenDAEjectionManager.lastEjectionInitiated(ejectee), block.timestamp); // should remain unchanged

        vm.stopPrank();
    }

    function testCancelEjectionByEjectee() public {
        // 0) Start the ejection
        testStartEjection(0, 0);

        // 1) Cancel the ejection on behalf of the ejectee
        vm.startPrank(ejectee);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCancelled(ejectee);
        eigenDAEjectionManager.cancelEjection();
        vm.stopPrank();

        // 2) Ensure the ejectee record is nullified
        assertEq(eigenDAEjectionManager.getEjector(ejectee), address(0));
        assertEq(eigenDAEjectionManager.ejectionTime(ejectee), 0);
        assertEq(eigenDAEjectionManager.lastEjectionInitiated(ejectee), block.timestamp); // should remain unchanged
    }

    function testCompleteEjection() public {
        // 0) start an ejection via ejector

        testStartEjection(0, 0);

        // 1) complete ejection via ejector
        vm.startPrank(testEjector);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCompleted(ejectee, "0x");
        eigenDAEjectionManager.completeEjection(ejectee, "0x");
        vm.stopPrank();

        // 2) ensure that ejectee's record is nullified and the
        //    ejector's book-kept balance reincorporates the initial deposit amount
        assertEq(eigenDAEjectionManager.getEjector(ejectee), address(0));
        assertEq(eigenDAEjectionManager.ejectionTime(ejectee), 0);
        assertEq(eigenDAEjectionManager.lastEjectionInitiated(ejectee), block.timestamp); // should remain unchanged
    }

    function testDelayEnforcementCausesEjectorCompletionsToRevert() public {
        // 0) set an artificial delay for which the ejector has to wait
        //    until completing the ejection
        testStartEjection(0, 6000);

        vm.startPrank(testEjector);
        vm.expectRevert("Proceeding not yet due");
        // 1) the EVM time context hasn't been advanced and there's an artificial
        //    delay where the block.timestamp >= start_ejection_block.timestamp + 6000s
        eigenDAEjectionManager.completeEjection(ejectee, "0x");

        // 2) now advance EVM and ensure that ejection can be successfully completed
        //    by ejector
        vm.warp(block.timestamp + 7000);
        eigenDAEjectionManager.completeEjection(ejectee, "0x");

        vm.stopPrank();
    }

    function testCoolDownEnforcementCausesAttemptedCompletionsToRevert() public {
        // 0) warp the time context

        vm.warp(block.timestamp + 7000);
        // 1) set an artificial cooldown period for which the ejector has to wait
        //    until completing the ejection
        testCancelEjectionByEjector(6000, 0);

        // 2) ensure that a too-early attempted ejector completion reverts
        vm.expectRevert("Ejection cooldown not met");
        vm.startPrank(testEjector);
        eigenDAEjectionManager.startEjection(ejectee, "0x");

        // 3) after the cooldown period has successfully elapsed, the ejector
        //    should be able to successfully start a new ejection
        vm.warp(block.timestamp + 7000);
        eigenDAEjectionManager.startEjection(ejectee, "0x");
        vm.stopPrank();
    }
}
