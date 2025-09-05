// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {Test} from "forge-std/Test.sol";

import {EigenDAEjectionManager} from "src/periphery/ejection/EigenDAEjectionManager.sol";
import {EigenDAEjectionLib} from "src/periphery/ejection/libraries/EigenDAEjectionLib.sol";
import {EigenDAEjectionTypes} from "src/periphery/ejection/libraries/EigenDAEjectionTypes.sol";

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
    uint256 constant CANCEL_EJECTION_WITHOUT_SIG_GAS_REFUND = 39128;
    uint256 constant CANCEL_EJECTION_WITH_SIG_GAS_REFUND = 70000;
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
            new EigenDAEjectionManager(address(token), DEPOSIT_BASE_FEE_MULTIPLIER, address(directory), 39128, 70000);
        accessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, address(this));
        directory.addAddress(AddressDirectoryConstants.EIGEN_DA_EJECTION_MANAGER_NAME, address(ejectionManager));
        directory.addAddress(AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME, address(this));
    }

    /// 1. Takes a deposit from the caller.
    /// 2. Starts the ejection process for the operator.
    /// 2a. sets quorums
    /// 2b. sets proceeding time to timestamp + delay
    /// 2c. sets proceeding initiated time to current timestamp
    /// 3. Emits EjectionStarted event.
    function testStartEjection(address caller, address ejectee) public {
        accessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, caller);
        vm.assume(caller != address(0) && ejectee != address(0) && caller != ejectee);
        token.mint(caller, EXPECTED_DEPOSIT);

        vm.startPrank(caller);
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

    function testCancelEjectionByEjector(address caller, address ejectee) public {
        testStartEjection(caller, ejectee);
        vm.startPrank(caller);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCancelled(ejectee);
        vm.startSnapshotGas("CANCEL EJECTION");
        ejectionManager.cancelEjectionByEjector(ejectee);
        vm.stopSnapshotGas("CANCEL EJECTION");
        vm.stopPrank();
        assertEq(ejectionManager.getEjector(ejectee), address(0));
        assertEq(ejectionManager.ejectionTime(ejectee), 0);
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp); // should remain unchanged
    }

    function testCancelEjectionByEjectee(address caller, address ejectee) public {
        testStartEjection(caller, ejectee);
        vm.startPrank(ejectee);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCancelled(ejectee);
        ejectionManager.cancelEjection();
        vm.stopPrank();
        assertEq(ejectionManager.getEjector(ejectee), address(0));
        assertEq(ejectionManager.ejectionTime(ejectee), 0);
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp); // should remain unchanged
    }

    function testCompleteEjection(address caller, address ejectee) public {
        testStartEjection(caller, ejectee);
        vm.startPrank(caller);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCompleted(ejectee, "0x");
        ejectionManager.completeEjection(ejectee, "0x");
        vm.stopPrank();
    }
}
