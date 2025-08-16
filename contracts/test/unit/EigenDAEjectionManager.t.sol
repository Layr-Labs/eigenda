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

    uint256 constant DEPOSIT_AMOUNT = 1e18;

    function setUp() public {
        token = new ERC20Mintable("TestToken", "TTK");
        accessControl = new EigenDAAccessControl(address(this));
        directory = new EigenDADirectory();
        directory.initialize(address(accessControl));
        ejectionManager = new EigenDAEjectionManager(address(token), DEPOSIT_AMOUNT, address(directory), 1);
        accessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, address(this));
        directory.addAddress(AddressDirectoryConstants.EIGEN_DA_EJECTION_MANAGER_NAME, address(ejectionManager));
    }

    /// 1. Takes a deposit from the caller.
    /// 2. Starts the ejection process for the operator.
    /// 2a. sets quorums
    /// 2b. sets proceeding time to timestamp + delay
    /// 2c. sets proceeding initiated time to current timestamp
    /// 3. Emits EjectionStarted event.
    function testStartEjection(address caller, address ejectee) public {
        accessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, caller);
        token.mint(caller, DEPOSIT_AMOUNT);

        vm.startPrank(caller);
        token.approve(address(ejectionManager), DEPOSIT_AMOUNT);
        ejectionManager.addEjectorBalance(DEPOSIT_AMOUNT);

        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionStarted(
            ejectee,
            caller,
            "0x", // quorums (empty for this test)
            uint64(block.timestamp),
            uint64(block.timestamp + ejectionManager.ejectionDelay())
        );

        ejectionManager.startEjection(ejectee, "0x");
        vm.stopPrank();
        assertEq(ejectionManager.getEjector(ejectee), caller);
        assertEq(ejectionManager.ejectionTime(ejectee), block.timestamp + ejectionManager.ejectionDelay());
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp);
    }

    function testCancelEjectionByEjector() public {
        // Add test logic for canceling ejection by ejector
    }

    function testCompleteEjection() public {
        // Add test logic for completing ejection
    }
}
