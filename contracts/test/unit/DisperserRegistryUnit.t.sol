// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {DisperserRegistry} from "src/core/DisperserRegistry.sol";
import {DisperserRegistryLib} from "src/core/libraries/v3/disperser/DisperserRegistryLib.sol";
import {DisperserRegistryTypes} from "src/core/libraries/v3/disperser/DisperserRegistryTypes.sol";
import {DisperserRegistryStorage} from "src/core/libraries/v3/disperser/DisperserRegistryStorage.sol";
import {IDisperserRegistry} from "src/core/interfaces/IDisperserRegistry.sol";
import {Constants} from "src/core/libraries/Constants.sol";
import {AccessControlLib} from "src/core/libraries/AccessControlLib.sol";

import {ERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

import {Test} from "forge-std/Test.sol";

contract MockERC20 is ERC20 {
    constructor(string memory name, string memory symbol) ERC20(name, symbol) {}

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }

    function burn(address from, uint256 amount) external {
        _burn(from, amount);
    }
}

contract DisperserRegistryTestHarness is DisperserRegistry {
    function setNextDisperserKey(uint32 nextDisperserKey) external {
        DisperserRegistryStorage.Layout storage s = DisperserRegistryStorage.layout();
        s.nextDisperserKey = nextDisperserKey;
    }

    function consumeDisperserKey() external returns (uint32) {
        return DisperserRegistryLib.consumeDisperserKey();
    }

    function setDisperserOwner(uint32 disperserKey, address owner) external {
        DisperserRegistryStorage.Layout storage s = DisperserRegistryStorage.layout();
        s.disperser[disperserKey].owner = owner;
    }
}

contract DisperserRegistryUnit is Test {
    uint256 constant DEPOSIT_AMOUNT = 1 ether;
    uint256 constant REFUND_AMOUNT = 0.1 ether;
    uint256 constant UPDATE_FEE = 0.01 ether;
    uint64 constant LOCK_PERIOD = 1 hours;

    address constant OWNER = address(uint160(uint256(keccak256("owner"))));
    DisperserRegistryTestHarness public disperserRegistry;
    MockERC20 public mockToken;

    function setUp() public {
        mockToken = new MockERC20("MockToken", "MTK");
        disperserRegistry = new DisperserRegistryTestHarness();
        DisperserRegistryTypes.LockedDisperserDeposit memory depositParams = DisperserRegistryTypes
            .LockedDisperserDeposit({
            token: address(mockToken),
            deposit: DEPOSIT_AMOUNT,
            refund: REFUND_AMOUNT,
            lockPeriod: LOCK_PERIOD
        });
        disperserRegistry.initialize(OWNER, depositParams, UPDATE_FEE);
    }

    function test_InitializeDisperserRegistry() public view {
        DisperserRegistryTypes.LockedDisperserDeposit memory depositParams = disperserRegistry.getDepositParams();
        assertEq(depositParams.token, address(mockToken), "Token address mismatch");
        assertEq(depositParams.deposit, DEPOSIT_AMOUNT, "Deposit amount mismatch");
        assertEq(depositParams.refund, REFUND_AMOUNT, "Refund amount mismatch");
        assertEq(depositParams.lockPeriod, LOCK_PERIOD, "Lock period mismatch");
        assertEq(disperserRegistry.owner(), OWNER, "Owner address mismatch");
    }

    function test_ConsumeDisperserKey(uint32 startKey) public {
        vm.assume(startKey != type(uint32).max);
        disperserRegistry.setNextDisperserKey(startKey);
        assertEq(disperserRegistry.getNextDisperserKey(), startKey, "Next disperser key should match the set value");
        // Increment the key for the next call
        disperserRegistry.setNextDisperserKey(startKey);

        assertEq(
            disperserRegistry.consumeDisperserKey(), startKey, "Consumed disperser key should match the initial key"
        );
        assertEq(disperserRegistry.getNextDisperserKey(), startKey + 1, "Next disperser key should be incremented");
    }

    function test_RegisterDisperser(address caller, address disperserAddress, string memory disperserURL)
        public
        returns (uint32 disperserKey)
    {
        vm.assume(disperserAddress != address(0));
        vm.assume(caller != address(0) && caller != address(disperserRegistry));

        mockToken.mint(caller, DEPOSIT_AMOUNT);
        vm.prank(caller);
        mockToken.approve(address(disperserRegistry), DEPOSIT_AMOUNT);

        DisperserRegistryTypes.LockedDisperserDeposit memory depositParams = disperserRegistry.getDepositParams();

        vm.expectEmit(true, true, true, true);
        emit DisperserRegistryLib.DisperserDepositTaken(disperserRegistry.getNextDisperserKey(), depositParams);
        vm.expectEmit(true, true, true, true);
        emit DisperserRegistryLib.DisperserRegistered(
            disperserRegistry.getNextDisperserKey(), disperserAddress, caller, disperserURL
        );

        vm.prank(caller);
        disperserKey = disperserRegistry.registerDisperser(disperserAddress, disperserURL);

        assertEq(disperserRegistry.getDisperserAddress(disperserKey), disperserAddress, "Disperser address mismatch");
        assertEq(disperserRegistry.getDisperserURL(disperserKey), disperserURL, "Disperser URL mismatch");
        assertEq(disperserRegistry.getDisperserOwner(disperserKey), caller, "Disperser owner mismatch");
        DisperserRegistryTypes.LockedDisperserDeposit memory disperserDepositParams =
            disperserRegistry.getDisperserDepositParams(disperserKey);
        assertEq(disperserDepositParams.token, address(mockToken), "Deposit token address mismatch");
        assertEq(disperserDepositParams.deposit, DEPOSIT_AMOUNT, "Deposit amount mismatch");
        assertEq(disperserDepositParams.refund, REFUND_AMOUNT, "Refund amount mismatch");
        assertTrue(
            disperserRegistry.getDisperserDepositUnlockTime(disperserKey) == 0,
            "Deposit unlock time should be in the future"
        );
        assertEq(mockToken.balanceOf(caller), 0, "Caller balance should reflect deposit deduction");
        assertEq(
            mockToken.balanceOf(address(disperserRegistry)),
            DEPOSIT_AMOUNT,
            "Disperser registry balance should reflect deposit"
        );
        assertEq(
            disperserRegistry.getExcessBalance(address(mockToken)),
            DEPOSIT_AMOUNT - REFUND_AMOUNT,
            "Excess balance should match deposit amount"
        );
    }

    function test_RegisterDisperserFailsIfDisperserAddressIsZero(address caller, string memory disperserURL) public {
        vm.assume(caller != address(0));
        vm.prank(caller);
        vm.expectRevert(abi.encodeWithSelector(IDisperserRegistry.InvalidDisperserAddress.selector, address(0)));
        disperserRegistry.registerDisperser(address(0), disperserURL);
    }

    function test_TransferDisperserOwnership(
        address caller,
        address newOwner,
        address disperserAddress,
        string memory disperserURL
    ) public {
        vm.assume(newOwner != address(0));
        uint32 disperserKey = test_RegisterDisperser(caller, disperserAddress, disperserURL);

        vm.expectEmit(true, true, true, true);
        emit DisperserRegistryLib.DisperserOwnershipTransferred(disperserKey, newOwner);
        vm.prank(caller);
        disperserRegistry.transferDisperserOwnership(disperserKey, newOwner);
    }

    function test_TransferDisperserOwnershipFailsIfDisperserNotRegistered(address newOwner, uint32 disperserKey)
        public
    {
        vm.assume(newOwner != address(0));

        vm.expectRevert(abi.encodeWithSelector(IDisperserRegistry.DisperserNotRegistered.selector, disperserKey));
        vm.prank(address(0));
        disperserRegistry.transferDisperserOwnership(disperserKey, newOwner);
    }

    function test_TransferDisperserOwnershipFailsIfNewOwnerIsZero(
        address caller,
        address disperserAddress,
        string memory disperserURL
    ) public {
        uint32 disperserKey = test_RegisterDisperser(caller, disperserAddress, disperserURL);

        vm.expectRevert(abi.encodeWithSelector(IDisperserRegistry.InvalidNewOwner.selector, address(0)));
        vm.prank(caller);
        disperserRegistry.transferDisperserOwnership(disperserKey, address(0));
    }

    function test_UpdateDisperserInfo(address caller, address disperserAddress, string memory disperserURL) public {
        vm.assume(caller != address(0));
        uint32 disperserKey = test_RegisterDisperser(caller, address(1), "initial_url");

        mockToken.mint(caller, UPDATE_FEE);
        vm.prank(caller);
        mockToken.approve(address(disperserRegistry), UPDATE_FEE);

        vm.expectEmit(true, true, true, true);
        emit DisperserRegistryLib.DisperserUpdateFeeTaken(disperserKey, caller, UPDATE_FEE);
        vm.expectEmit(true, true, true, true);
        emit DisperserRegistryLib.DisperserUpdated(disperserKey, disperserAddress, disperserURL);
        vm.prank(caller);
        disperserRegistry.updateDisperserInfo(disperserKey, disperserAddress, disperserURL);

        assertEq(disperserRegistry.getDisperserURL(disperserKey), disperserURL, "Disperser URL should be updated");
        assertEq(
            disperserRegistry.getDisperserAddress(disperserKey), disperserAddress, "Disperser address should be updated"
        );
        assertEq(mockToken.balanceOf(caller), 0, "Caller balance should reflect update fee deduction");
        assertEq(
            mockToken.balanceOf(address(disperserRegistry)),
            DEPOSIT_AMOUNT + UPDATE_FEE,
            "Disperser registry balance should reflect update fee"
        );
    }

    function test_updateDisperserInfoFailsIfDisperserNotRegistered(address disperserAddress, string memory disperserURL)
        public
    {
        vm.prank(address(0));
        vm.expectRevert(abi.encodeWithSelector(IDisperserRegistry.DisperserNotRegistered.selector, 0));
        disperserRegistry.updateDisperserInfo(0, disperserAddress, disperserURL);
    }

    function test_UpdateDisperserInfoFailsIfDisperserAddressIsZero(
        address caller,
        address disperserAddress,
        string memory disperserURL
    ) public {
        uint32 disperserKey = test_RegisterDisperser(caller, disperserAddress, disperserURL);
        vm.prank(caller);
        vm.expectRevert(abi.encodeWithSelector(IDisperserRegistry.InvalidDisperserAddress.selector, address(0)));
        disperserRegistry.updateDisperserInfo(disperserKey, address(0), disperserURL);
    }

    function test_DeregisterDisperser(address caller, address disperserAddress, string memory disperserURL)
        public
        returns (uint32 disperserKey)
    {
        disperserKey = test_RegisterDisperser(caller, disperserAddress, disperserURL);

        vm.expectEmit(true, true, true, true);
        emit DisperserRegistryLib.DisperserDeregistered(disperserKey, uint64(block.timestamp + LOCK_PERIOD));
        vm.prank(caller);
        disperserRegistry.deregisterDisperser(disperserKey);

        assertEq(disperserRegistry.getDisperserAddress(disperserKey), address(0), "Disperser address should be cleared");
        assertEq(disperserRegistry.getDisperserURL(disperserKey), "", "Disperser URL should be cleared");
        assertEq(disperserRegistry.getDisperserOwner(disperserKey), caller, "Disperser owner should be unchanged");
        assertEq(
            disperserRegistry.getDisperserDepositUnlockTime(disperserKey),
            uint64(block.timestamp + LOCK_PERIOD),
            "Deposit unlock time should be set to future time"
        );
        assertEq(mockToken.balanceOf(caller), 0, "Caller should still have no balance after deregistration");
        assertEq(
            mockToken.balanceOf(address(disperserRegistry)),
            DEPOSIT_AMOUNT,
            "Disperser registry should retain deposit amount"
        );
    }

    function test_DeregisterDisperserFailsIfDisperserNotRegistered(uint32 disperserKey) public {
        vm.prank(address(0));
        vm.expectRevert(abi.encodeWithSelector(IDisperserRegistry.DisperserNotRegistered.selector, disperserKey));
        disperserRegistry.deregisterDisperser(disperserKey);
    }

    function test_WithdrawDisperserDeposit(address caller, address disperserAddress, string memory disperserURL)
        public
    {
        uint32 disperserKey = test_DeregisterDisperser(caller, disperserAddress, disperserURL);

        vm.warp(block.timestamp + LOCK_PERIOD); // Move time past the lock period
        vm.expectEmit(true, true, true, true);
        emit DisperserRegistryLib.DisperserRefundIssued(disperserKey, address(mockToken), REFUND_AMOUNT);
        vm.prank(caller);
        disperserRegistry.withdrawDisperserDeposit(disperserKey);
        assertEq(mockToken.balanceOf(caller), REFUND_AMOUNT, "Caller should receive the refund amount");
        assertEq(
            mockToken.balanceOf(address(disperserRegistry)),
            DEPOSIT_AMOUNT - REFUND_AMOUNT,
            "Disperser registry should retain excess balance after refund"
        );
    }

    function test_WithdrawDisperserDepositFailsIfDisperserNotDeregistered(
        address caller,
        address disperserAddress,
        string memory disperserURL
    ) public {
        uint32 disperserKey = test_RegisterDisperser(caller, disperserAddress, disperserURL);
        vm.expectRevert(abi.encodeWithSelector(IDisperserRegistry.DisperserNotDeregistered.selector, disperserKey));
        vm.prank(caller);
        disperserRegistry.withdrawDisperserDeposit(disperserKey);
    }

    function test_WithdrawDisperserDepositFailsIfRefundLocked(
        address caller,
        address disperserAddress,
        string memory disperserURL
    ) public {
        uint32 disperserKey = test_DeregisterDisperser(caller, disperserAddress, disperserURL);

        vm.expectRevert(
            abi.encodeWithSelector(
                IDisperserRegistry.RefundLocked.selector, disperserKey, uint64(block.timestamp + LOCK_PERIOD)
            )
        );
        vm.prank(caller);
        disperserRegistry.withdrawDisperserDeposit(disperserKey);
    }

    function test_WithdrawDisperserDepositFailsIfZeroRefund(
        address caller,
        address disperserAddress,
        string memory disperserURL
    ) public {
        vm.prank(OWNER);
        disperserRegistry.setDepositParams(
            DisperserRegistryTypes.LockedDisperserDeposit({
                token: address(mockToken),
                deposit: 0,
                refund: 0,
                lockPeriod: LOCK_PERIOD
            })
        );
        vm.prank(caller);
        uint32 disperserKey = disperserRegistry.registerDisperser(disperserAddress, disperserURL);
        vm.prank(caller);
        disperserRegistry.deregisterDisperser(disperserKey);
        vm.warp(block.timestamp + LOCK_PERIOD); // Move time past the lock period

        vm.expectRevert(abi.encodeWithSelector(IDisperserRegistry.ZeroRefund.selector, disperserKey));
        vm.prank(caller);
        disperserRegistry.withdrawDisperserDeposit(disperserKey);
    }

    function test_SetDepositParams(address newToken, uint256 newDeposit, uint256 newRefund, uint64 newLockPeriod)
        public
    {
        vm.assume(newToken != address(0));
        vm.assume(newDeposit > 0 && newRefund > 0 && newLockPeriod > 0 && newDeposit >= newRefund);

        DisperserRegistryTypes.LockedDisperserDeposit memory newDepositParams = DisperserRegistryTypes
            .LockedDisperserDeposit({token: newToken, deposit: newDeposit, refund: newRefund, lockPeriod: newLockPeriod});
        vm.expectEmit(true, true, true, true);
        emit DisperserRegistryLib.DepositParamsSet(newDepositParams);

        vm.prank(OWNER);
        disperserRegistry.setDepositParams(newDepositParams);
        DisperserRegistryTypes.LockedDisperserDeposit memory updatedDepositParams = disperserRegistry.getDepositParams();
        assertEq(updatedDepositParams.token, newToken, "Deposit token address should match the new token");
        assertEq(updatedDepositParams.deposit, newDeposit, "Deposit amount should match the new deposit");
        assertEq(updatedDepositParams.refund, newRefund, "Refund amount should match the new refund");
        assertEq(updatedDepositParams.lockPeriod, newLockPeriod, "Lock period should match the new lock period");
    }

    function test_SetUpdateFee(uint256 newUpdateFee) public {
        vm.assume(newUpdateFee > 0);

        vm.expectEmit(true, true, true, true);
        emit DisperserRegistryLib.UpdateFeeSet(newUpdateFee);
        vm.prank(OWNER);
        disperserRegistry.setUpdateFee(newUpdateFee);
        assertEq(disperserRegistry.getUpdateFee(), newUpdateFee, "Update fee should be updated");
    }

    function test_OwnerFunctions(address notOwner) public {
        vm.assume(notOwner != OWNER && notOwner != address(0));
        vm.startPrank(notOwner);

        vm.expectRevert(abi.encodeWithSelector(AccessControlLib.MissingRole.selector, Constants.OWNER_ROLE, notOwner));
        disperserRegistry.setDepositParams(
            DisperserRegistryTypes.LockedDisperserDeposit({
                token: address(mockToken),
                deposit: DEPOSIT_AMOUNT,
                refund: REFUND_AMOUNT,
                lockPeriod: LOCK_PERIOD
            })
        );

        vm.expectRevert(abi.encodeWithSelector(AccessControlLib.MissingRole.selector, Constants.OWNER_ROLE, notOwner));
        disperserRegistry.setUpdateFee(UPDATE_FEE);

        vm.expectRevert(abi.encodeWithSelector(AccessControlLib.MissingRole.selector, Constants.OWNER_ROLE, notOwner));
        disperserRegistry.transferOwnership(notOwner);
        vm.stopPrank();
    }

    function test_DisperserOwnerFunctions(uint32 disperserKey, address disperserOwner, address notDisperserOwner)
        public
    {
        vm.assume(disperserOwner != address(0) && notDisperserOwner != address(0));
        vm.assume(disperserOwner != notDisperserOwner);
        disperserRegistry.setDisperserOwner(disperserKey, disperserOwner);

        vm.startPrank(notDisperserOwner);
        vm.expectRevert(
            abi.encodeWithSelector(IDisperserRegistry.NotDisperserOwner.selector, disperserKey, disperserOwner)
        );
        disperserRegistry.transferDisperserOwnership(disperserKey, notDisperserOwner);

        vm.expectRevert(
            abi.encodeWithSelector(IDisperserRegistry.NotDisperserOwner.selector, disperserKey, disperserOwner)
        );
        disperserRegistry.updateDisperserInfo(disperserKey, address(1), "new_url");

        vm.expectRevert(
            abi.encodeWithSelector(IDisperserRegistry.NotDisperserOwner.selector, disperserKey, disperserOwner)
        );
        disperserRegistry.deregisterDisperser(disperserKey);

        vm.expectRevert(
            abi.encodeWithSelector(IDisperserRegistry.NotDisperserOwner.selector, disperserKey, disperserOwner)
        );
        disperserRegistry.withdrawDisperserDeposit(disperserKey);
    }
}
