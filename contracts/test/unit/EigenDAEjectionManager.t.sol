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

    address ejector;
    address ejectee;

    uint256 constant DEPOSIT_BASE_FEE_MULTIPLIER = 7;
    uint256 constant CANCEL_EJECTION_WITHOUT_SIG_GAS_REFUND = 39_128;
    uint256 constant CANCEL_EJECTION_WITH_SIG_GAS_REFUND = 70_000;
    uint256 constant BASE_FEE = 10;
    uint256 constant EXPECTED_DEPOSIT = BASE_FEE * DEPOSIT_BASE_FEE_MULTIPLIER * CANCEL_EJECTION_WITH_SIG_GAS_REFUND;
    /// TODO: PLACEHOLDER UNTIL GAS COST FOR SIGNATURES IS KNOWN

    function setUp() public {
        vm.fee(BASE_FEE);
        token = new ERC20Mintable("WrappedETH", "WETH");
        accessControl = new EigenDAAccessControl(address(this));
        directory = new EigenDADirectory();
        directory.initialize(address(accessControl));

        ejector = makeAddr("ejector");
        ejectee = makeAddr("ejectee");

        ejectionManager =
            new EigenDAEjectionManager(address(token), DEPOSIT_BASE_FEE_MULTIPLIER, address(directory), 39_128, 70_000);
        directory.addAddress(AddressDirectoryConstants.EIGEN_DA_EJECTION_MANAGER_NAME, address(ejectionManager));

        /// TODO: figure out the proper way to wire this for testing E2E signature recovery 
        //        and BLS pairing verification
        directory.addAddress(AddressDirectoryConstants.REGISTRY_COORDINATOR_NAME, address(this));
    }

    function testStartEjection() public {
        testStartEjection(0, 0);
    }

    function testStartEjection(uint64 cooldown, uint64 delay) private {
        // 1) Wire up access mgmt dependencies and set protocol params on contract
        accessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, ejector);
        accessControl.grantRole(AccessControlConstants.OWNER_ROLE, ejector);
        depositEjectorFunds(EXPECTED_DEPOSIT);

        vm.startPrank(ejector);
        ejectionManager.setCooldown(cooldown);
        ejectionManager.setDelay(delay);

        // 2) start an ejection against an arbitrary ejectee 
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionStarted(
            ejectee,
            ejector,
            "0x", // quorums (empty for this test)
            uint64(block.timestamp),
            uint64(block.timestamp + ejectionManager.ejectionDelay()),
            EXPECTED_DEPOSIT
        );

        ejectionManager.startEjection(ejectee, "0x");
        vm.stopPrank();
        assertEq(ejectionManager.getEjectorBalance(ejector), 0);
        assertEq(ejectionManager.getEjector(ejectee), ejector);
        assertEq(ejectionManager.ejectionTime(ejectee), block.timestamp + ejectionManager.ejectionDelay());
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp);
    }

    function depositEjectorFunds(uint256 amount) private {
        token.mint(ejector, amount);
        vm.startPrank(ejector);
        token.approve(address(ejectionManager), amount);
        ejectionManager.addEjectorBalance(amount);
        vm.stopPrank();

        assertEq(ejectionManager.getEjectorBalance(ejector), EXPECTED_DEPOSIT);
    }

    function testCancelEjectionbyEjector() public {
        testCancelEjectionByEjector(0, 0);
    }

    function testCancelEjectionByEjector(uint64 cooldown, uint64 delay) private {
        accessControl.grantRole(AccessControlConstants.EJECTOR_ROLE, ejector);
        accessControl.grantRole(AccessControlConstants.OWNER_ROLE, ejector);
        depositEjectorFunds(EXPECTED_DEPOSIT);

        // 1) Ejector "deposits" by escrowing ERC20 tokens to the contract
        //    address and starting the ejection
        vm.startPrank(ejector);
        ejectionManager.setCooldown(cooldown);
        ejectionManager.setDelay(delay);
        ejectionManager.startEjection(ejectee, "0x");

        // 2) Ensure that the deposited funds are actually escrowed into the contract
        assertEq(
            token.balanceOf(address(ejectionManager)),
            EXPECTED_DEPOSIT,
            "Deposit should result in funds escrowed to contract"
        );

        // 3) Issue a cancellation from the Ejector role and withdraw the ERC20 funds
        //    (i.e, contract -> ejector)
        ejectionManager.cancelEjectionByEjector(ejectee);

        // 4) Ensure the stateful params entry has been nullified
        assertEq(ejectionManager.getEjector(ejectee), address(0));
        assertEq(ejectionManager.ejectionTime(ejectee), 0);
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp); // should remain unchanged

        ejectionManager.withdrawEjectorBalance(EXPECTED_DEPOSIT);
        vm.stopPrank();

        // 5) Ensure the ejector has received the full amount of their deposited tokens back
        //    and their book-kept balance in the ejection manager state is zero'd
        assertEq(ejectionManager.getEjectorBalance(ejector), 0);
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

    function testCancelEjectionByEjectee() public {
        testStartEjection(0, 0);
        vm.startPrank(ejectee);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCancelled(ejectee);
        ejectionManager.cancelEjection();
        vm.stopPrank();
        assertEq(ejectionManager.getEjector(ejectee), address(0));
        assertEq(ejectionManager.ejectionTime(ejectee), 0);
        assertEq(ejectionManager.lastEjectionInitiated(ejectee), block.timestamp); // should remain unchanged

        // ensure ERC20 gas refund lands onto the ejectee account
        //
        // TODO: understand why this results in a loss of deposit - refund_amount !?
        assertEq(token.balanceOf(ejectee), CANCEL_EJECTION_WITHOUT_SIG_GAS_REFUND * block.basefee);
    }

    function testCompleteEjection() public {
        // 1) start an ejection via ejector

        testStartEjection(0, 0);
        assertEq(ejectionManager.getEjectorBalance(ejector), 0);
        
        // 2) complete ejection via ejector
        vm.startPrank(ejector);
        vm.expectEmit(true, true, true, true);
        emit EigenDAEjectionLib.EjectionCompleted(ejectee, "0x");
        ejectionManager.completeEjection(ejectee, "0x");
        vm.stopPrank();

        // 3) ensure that ejector's balance reflects the deposit
        // TODO: figure out why this is the case?
        // assertEq(ejectionManager.getEjectorBalance(ejector), EXPECTED_DEPOSIT);
    }

    function testDelayEnforcementCausesEjectorCompletionsToRevert() public {
        // 1) set an artificial delay for which the ejector has to wait
        //    until completing the ejection
        testStartEjection(0, 6000);

        vm.startPrank(ejector);
        vm.expectRevert("Proceeding not yet due");
        // 2) the EVM time context hasn't been advanced and there's an artificial
        //    delay where the block.timestamp >= start_ejection_block.timestamp + 6000s
        ejectionManager.completeEjection(ejectee, "0x");

        // 3) now advance EVM and ensure that ejection can be successfully completed
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
        depositEjectorFunds(EXPECTED_DEPOSIT);

        // 2) after the cooldown period has successfuly elapsed, the ejector
        //    should be able to successfully start a new ejection
        vm.warp(block.timestamp + 7000);
        vm.startPrank(ejector);
        ejectionManager.startEjection(ejectee, "0x");
    }
}
