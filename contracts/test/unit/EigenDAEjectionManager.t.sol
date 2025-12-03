// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

import {Test} from "forge-std/Test.sol";

import {EigenDAEjectionManager} from "src/periphery/ejection/EigenDAEjectionManager.sol";
import {EigenDAEjectionLib} from "src/periphery/ejection/libraries/EigenDAEjectionLib.sol";

import {EigenDAAccessControl} from "src/core/EigenDAAccessControl.sol";
import {EigenDADirectory} from "src/core/EigenDADirectory.sol";

import {AccessControlConstants} from "src/core/libraries/v3/access-control/AccessControlConstants.sol";
import {AddressDirectoryConstants} from "src/core/libraries/v3/address-directory/AddressDirectoryConstants.sol";

import {ERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

contract ERC20Mintable is ERC20 {
    constructor(string memory name, string memory symbol) ERC20(name, symbol) {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}

contract EigenDAEjectionManagerTest is Test {
    EigenDADirectory directory;
    EigenDAAccessControl accessControl;
    EigenDAEjectionManager ejectionManager;
    ERC20Mintable token;

    uint256 constant DEPOSIT_BASE_FEE_MULTIPLIER = 7;
    uint256 constant CANCEL_EJECTION_WITHOUT_SIG_GAS_REFUND = 39_128;
    uint256 constant CANCEL_EJECTION_WITH_SIG_GAS_REFUND = 70_000;
    uint256 constant BASE_FEE = 10;
    uint256 constant EXPECTED_DEPOSIT = BASE_FEE * DEPOSIT_BASE_FEE_MULTIPLIER * CANCEL_EJECTION_WITH_SIG_GAS_REFUND;
    /// TODO: PLACEHOLDER UNTIL GAS COST FOR SIGNATURES IS KNOWN

    function setUp() public {
        vm.fee(BASE_FEE);
        token = new ERC20Mintable("TestToken", "TTK");
        accessControl = new EigenDAAccessControl(address(this));
        directory = new EigenDADirectory();
        directory.initialize(address(accessControl));
        ejectionManager =
            new EigenDAEjectionManager(address(token), DEPOSIT_BASE_FEE_MULTIPLIER, address(directory), 39_128, 70_000);
        accessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, address(this));
        directory.addAddress(AddressDirectoryConstants.EIGEN_DA_EJECTION_MANAGER_NAME, address(ejectionManager));
        directory.addAddress(AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME, address(this));
    }

    function testStartEjection(address caller, address ejectee) public {
        testStartEjection(caller, ejectee, 0, 0);
    }

    /// 1. Takes a deposit from the caller.
    /// 2. Starts the ejection process for the operator.
    /// 2a. sets quorums
    /// 2b. sets proceeding time to timestamp + delay
    /// 2c. sets proceeding initiated time to current timestamp
    /// 3. Emits EjectionStarted event.
    function testStartEjection(address caller, address ejectee, uint64 cooldown, uint64 delay) private {
        accessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, caller);
        accessControl.grantRole(AccessControlConstants.OWNER_ROLE, caller);

        vm.assume(caller != address(0) && ejectee != address(0) && caller != ejectee);
        token.mint(caller, EXPECTED_DEPOSIT);

        vm.startPrank(caller);
        ejectionManager.setCooldown(cooldown);
        ejectionManager.setDelay(delay);

        token.approve(address(ejectionManager), EXPECTED_DEPOSIT);
        ejectionManager.addEjectorBalance(EXPECTED_DEPOSIT);

        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionStarted(
            ejectee,
            caller,
            "0x", // quorums (empty for this test)
            uint64(block.timestamp),
            uint64(block.timestamp + ejectionManager.ejectionDelay()),
            EXPECTED_DEPOSIT
        );

        ejectionManager.startEjection(ejectee, "0x");
        vm.stopPrank();
        assertEq(ejectionManager.getEjector(ejectee), caller);
        assertEq(ejectionManager.ejectionTime(ejectee), block.timestamp + ejectionManager.ejectionDelay());
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp);
    }

    function testCancelEjectionByEjector(address ejector, address operator) public {
        accessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, ejector);
        token.mint(ejector, EXPECTED_DEPOSIT);

        // 1) Ejector "deposits" by escrowing ERC20 tokens to the contract
        //    address and starting the ejection
        vm.startPrank(ejector);
        token.approve(address(ejectionManager), EXPECTED_DEPOSIT);
        ejectionManager.addEjectorBalance(EXPECTED_DEPOSIT);
        ejectionManager.startEjection(operator, "0x");

        // 2) Ensure that the deposited funds are actually escrowed into the contract
        assertEq(
            token.balanceOf(address(ejectionManager)),
            EXPECTED_DEPOSIT,
            "Deposit should result in funds escrowed to contract"
        );

        // 3) Issue a cancellation from the Ejector role and withdraw the ERC20 funds
        //    (i.e, contract -> ejector)
        ejectionManager.cancelEjectionByEjector(operator);

        // 4) Ensure the stateful params entry has been nullified
        assertEq(ejectionManager.getEjector(operator), address(0));
        assertEq(ejectionManager.ejectionTime(operator), 0);
        assertEq(ejectionManager.lastEjectionInitiated(operator), block.timestamp); // should remain unchanged

        ejectionManager.withdrawEjectorBalance(EXPECTED_DEPOSIT);
        vm.stopPrank();

        // 5) Ensure the ejector has received the full amount of their deposited tokens back
        assertEq(
            token.balanceOf(address(ejectionManager)),
            0,
            "Ejections manager should not have any escrowed collateral tokens after ejector withdraw"
        );
        assertEq(
            token.balanceOf(address(ejector)),
            EXPECTED_DEPOSIT,
            "withdrawn tokens should be fully reissued to the ejector"
        );
    }

    function testCancelEjectionByEjectee(address caller, address ejectee) public {
        testStartEjection(caller, ejectee, 0, 0);
        vm.startPrank(ejectee);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCancelled(ejectee);
        ejectionManager.cancelEjection();
        vm.stopPrank();
        assertEq(ejectionManager.getEjector(ejectee), address(0));
        assertEq(ejectionManager.ejectionTime(ejectee), 0);
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp); // should remain unchanged

        // ensure ERC20 gas refund lands onto the ejectee account
        assertGt(token.balanceOf(ejectee), 0);
    }

    function testCompleteEjection(address caller, address ejectee) public {
        testStartEjection(caller, ejectee, 0, 0);
        vm.startPrank(caller);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCompleted(ejectee, "0x");
        ejectionManager.completeEjection(ejectee, "0x");
        vm.stopPrank();
    }

    function testDelayEnforcementCausesAttemptedCompletionsToRevert(address caller, address ejectee) public {
        // set an artificial delay for which the ejector has to wait
        // until completing the ejection
        testStartEjection(caller, ejectee, 0, 6000);

        vm.startPrank(caller);
        vm.expectRevert("Proceeding not yet due");
        // the EVM time context hasn't been advanced and there's an artificial
        // delay where the block.timestamp >= start_ejection_block.timestamp + 1000
        ejectionManager.completeEjection(ejectee, "0x");
        vm.stopPrank();
    }

    function testCoolDownEnforcementCausesAttemptedCompletionsToRevert(address caller, address ejectee) public {
        // 1) Ensure the test starts with a clean slate for the ejectee
        //    by warping time forward past any potential cooldown from previous fuzz iterations
        vm.warp(block.timestamp + 7000);

        // 2) set an artificial delay for which the ejector has to wait
        //    until completing the ejection
        testStartEjection(caller, ejectee, 6000, 0);

        // 3) have the ejectee successfully cancel the ejection
        vm.startPrank(ejectee);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCancelled(ejectee);
        ejectionManager.cancelEjection();
        vm.stopPrank();

        // 4) now the ejector attempts to eject the ejectee which reverts due
        //    to cooldown enforcements since block time hasn't advanced
        vm.startPrank(caller);
        vm.expectRevert();
        ejectionManager.startEjection(ejectee, "0x");
        vm.stopPrank();
    }
}
