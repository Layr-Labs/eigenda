// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

import {Test} from "forge-std/Test.sol";
import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import {ProxyAdmin} from "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import {EigenDAEjectionManager} from "src/periphery/ejection/EigenDAEjectionManager.sol";
import {EigenDAEjectionLib} from "src/periphery/ejection/libraries/EigenDAEjectionLib.sol";

import {EigenDAAccessControl} from "src/core/EigenDAAccessControl.sol";
import {EigenDADirectory} from "src/core/EigenDADirectory.sol";

import {AccessControlConstants} from "src/core/libraries/v3/access-control/AccessControlConstants.sol";
import {AddressDirectoryConstants} from "src/core/libraries/v3/address-directory/AddressDirectoryConstants.sol";

contract EigenDAEjectionManagerTest is Test {
    EigenDADirectory directory;
    EigenDAAccessControl accessControl;
    EigenDAEjectionManager ejectionManager;
    EigenDAEjectionManager ejectionManagerImplementation;
    ProxyAdmin proxyAdmin;

    address ejector;
    address ejectee;
    address untestedDep;

    uint256 constant DEPOSIT_BASE_FEE_MULTIPLIER = 7;
    uint256 constant CANCEL_EJECTION_WITHOUT_SIG_GAS_REFUND = 39_128;
    uint256 constant CANCEL_EJECTION_WITH_SIG_GAS_REFUND = 70_000;
    uint256 constant BASE_FEE = 10;
    uint256 constant EXPECTED_DEPOSIT = BASE_FEE * DEPOSIT_BASE_FEE_MULTIPLIER * CANCEL_EJECTION_WITH_SIG_GAS_REFUND;
    /// TODO: PLACEHOLDER UNTIL GAS COST FOR SIGNATURES IS KNOWN
    /// TODO: Add tests that ensure multiple ejections can be ran at once by a single ejector (1 ejector : N ejectees)
    ///       Also (N ejector : N ejectees)

    function setUp() public {
        vm.fee(BASE_FEE);
        accessControl = new EigenDAAccessControl(address(this));
        directory = new EigenDADirectory();
        directory.initialize(address(accessControl));

        ejector = makeAddr("ejector");
        ejectee = makeAddr("ejectee");

        untestedDep = makeAddr("untestedCalleeAddr");

        // Deploy proxy admin
        proxyAdmin = new ProxyAdmin();

        // Deploy implementation
        ejectionManagerImplementation = new EigenDAEjectionManager();

        // Encode initialize call
        bytes memory initData = abi.encodeWithSelector(
            EigenDAEjectionManager.initialize.selector,
            address(accessControl),
            address(untestedDep),
            address(untestedDep),
            address(this),
            DEPOSIT_BASE_FEE_MULTIPLIER,
            CANCEL_EJECTION_WITHOUT_SIG_GAS_REFUND,
            CANCEL_EJECTION_WITH_SIG_GAS_REFUND
        );

        // Deploy proxy
        TransparentUpgradeableProxy proxy =
            new TransparentUpgradeableProxy(address(ejectionManagerImplementation), address(proxyAdmin), initData);

        // Cast proxy to EigenDAEjectionManager
        ejectionManager = EigenDAEjectionManager(address(proxy));

        directory.addAddress(AddressDirectoryConstants.EIGEN_DA_EJECTION_MANAGER_NAME, address(ejectionManager));
    }

    function testStartEjection() public {
        testStartEjection(0, 0);
    }

    function testStartEjection(uint64 cooldown, uint64 delay) private {
        // 0) Wire up access mgmt dependencies and set protocol params on contract
        accessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, ejector);
        accessControl.grantRole(AccessControlConstants.OWNER_ROLE, ejector);

        vm.startPrank(ejector);
        ejectionManager.setCooldown(cooldown);
        ejectionManager.setDelay(delay);

        // 1) start an ejection against an arbitrary ejectee
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionStarted(
            ejectee,
            ejector,
            "0x", // quorums (empty for this test)
            uint64(block.timestamp),
            uint64(block.timestamp + ejectionManager.ejectionDelay())
        );

        ejectionManager.startEjection(ejectee, "0x");
        vm.stopPrank();

        // 2) verify that ejectee record was properly created
        assertEq(ejectionManager.getEjector(ejectee), ejector);
        assertEq(ejectionManager.ejectionTime(ejectee), block.timestamp + ejectionManager.ejectionDelay());
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp);
    }

    function testCancelEjectionByEjector() public {
        testCancelEjectionByEjector(0, 0);
    }

    function testCancelEjectionByEjector(uint64 cooldown, uint64 delay) private {
        // 0) grant roles
        accessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, ejector);
        accessControl.grantRole(AccessControlConstants.OWNER_ROLE, ejector);

        // 1) Ejector starts ejection for ejectee after setting contract params
        vm.startPrank(ejector);
        ejectionManager.setCooldown(cooldown);
        ejectionManager.setDelay(delay);
        ejectionManager.startEjection(ejectee, "0x");

        // 2) Issue a cancellation from the Ejector role
        ejectionManager.cancelEjectionByEjector(ejectee);

        // 3) Ensure the ejectee record has been nullified
        assertEq(ejectionManager.getEjector(ejectee), address(0));
        assertEq(ejectionManager.ejectionTime(ejectee), 0);
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp); // should remain unchanged

        vm.stopPrank();
    }

    function testCancelEjectionByEjectee() public {
        // 0) Start the ejection
        testStartEjection(0, 0);

        // 1) Cancel the ejection on behalf of the ejectee
        vm.startPrank(ejectee);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCancelled(ejectee);
        ejectionManager.cancelEjection();
        vm.stopPrank();

        // 2) Ensure the ejectee record is nullified
        assertEq(ejectionManager.getEjector(ejectee), address(0));
        assertEq(ejectionManager.ejectionTime(ejectee), 0);
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp); // should remain unchanged
    }

    function testCompleteEjection() public {
        // 0) start an ejection via ejector

        testStartEjection(0, 0);

        // 1) complete ejection via ejector
        vm.startPrank(ejector);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCompleted(ejectee, "0x");
        ejectionManager.completeEjection(ejectee, "0x");
        vm.stopPrank();

        // 2) ensure that ejectee's record is nullified and the
        //    ejector's book-kept balance reincorporates the initial deposit amount
        assertEq(ejectionManager.getEjector(ejectee), address(0));
        assertEq(ejectionManager.ejectionTime(ejectee), 0);
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp); // should remain unchanged
    }

    function testDelayEnforcementCausesEjectorCompletionsToRevert() public {
        // 0) set an artificial delay for which the ejector has to wait
        //    until completing the ejection
        testStartEjection(0, 6000);

        vm.startPrank(ejector);
        vm.expectRevert("Proceeding not yet due");
        // 1) the EVM time context hasn't been advanced and there's an artificial
        //    delay where the block.timestamp >= start_ejection_block.timestamp + 6000s
        ejectionManager.completeEjection(ejectee, "0x");

        // 2) now advance EVM and ensure that ejection can be successfully completed
        //    by ejector
        vm.warp(block.timestamp + 7000);
        ejectionManager.completeEjection(ejectee, "0x");

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
        vm.startPrank(ejector);
        ejectionManager.startEjection(ejectee, "0x");

        // 3) after the cooldown period has successfully elapsed, the ejector
        //    should be able to successfully start a new ejection
        vm.warp(block.timestamp + 7000);
        ejectionManager.startEjection(ejectee, "0x");
        vm.stopPrank();
    }
}
