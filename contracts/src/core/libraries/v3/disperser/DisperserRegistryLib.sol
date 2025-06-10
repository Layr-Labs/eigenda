// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";
import {DisperserRegistryTypes} from "src/core/libraries/v3/disperser/DisperserRegistryTypes.sol";
import {DisperserRegistryStorage} from "src/core/libraries/v3/disperser/DisperserRegistryStorage.sol";
import {IDisperserRegistry} from "src/core/interfaces/IDisperserRegistry.sol";

library DisperserRegistryLib {
    using SafeERC20 for IERC20;

    /// @notice Emitted when a new disperser is registered.
    event DisperserRegistered(
        uint32 indexed disperserKey, address indexed disperserAddress, address indexed owner, string disperserURL
    );

    /// @notice Emitted when a deposit is taken for a disperser.
    event DisperserDepositTaken(
        uint32 indexed disperserKey, DisperserRegistryTypes.LockedDisperserDeposit depositParams
    );

    /// @notice Emitted when ownership of a disperser is transferred.
    event DisperserOwnershipTransferred(uint32 indexed disperserKey, address indexed newOwner);

    /// @notice Emitted when a disperser's update fee is taken.
    event DisperserUpdateFeeTaken(uint32 indexed disperserKey, address indexed owner, uint256 updateFee);

    /// @notice Emitted when a disperser's address or URL is updated.
    event DisperserUpdated(uint32 indexed disperserKey, address indexed disperserAddress, string disperserURL);

    /// @notice Emitted when a disperser is deregistered.
    event DisperserDeregistered(uint32 indexed disperserKey, uint64 unlockTimestamp);

    /// @notice Emitted when a refund is issued for a deregistered disperser.
    event DisperserRefundIssued(uint32 indexed disperserKey, address indexed token, uint256 refundAmount);

    /// @notice Emitted when the deposit parameters for dispersers are set.
    event DepositParamsSet(DisperserRegistryTypes.LockedDisperserDeposit depositParams);

    /// @notice Emitted when the update fee for dispersers is set.
    event UpdateFeeSet(uint256 updateFee);

    /// @notice Emitted when a disperser key is consumed.
    event DisperserKeyConsumed(uint32 indexed disperserKey);

    function s() internal pure returns (DisperserRegistryStorage.Layout storage) {
        return DisperserRegistryStorage.layout();
    }

    /// @notice Consumes a disperser key and increments the nextDisperserKey counter.
    function consumeDisperserKey() internal returns (uint32) {
        uint32 disperserKey = s().nextDisperserKey;
        s().nextDisperserKey++;
        emit DisperserKeyConsumed(disperserKey);
        return disperserKey;
    }

    /// @notice Registers a new disperser with the given address and URL. Takes a deposit from the caller and registers the current deposit parameters.
    /// @param disperserAddress The address of the disperser.
    /// @param disperserURL The URL of the disperser.
    /// @return disperserKey The key assigned to the registered disperser.
    function registerDisperser(address disperserAddress, string memory disperserURL)
        internal
        returns (uint32 disperserKey)
    {
        disperserKey = consumeDisperserKey();
        DisperserRegistryTypes.DisperserInfo storage disperser = s().disperser[disperserKey];
        DisperserRegistryTypes.LockedDisperserDeposit memory depositParams = s().depositParams; // memory copy for gas efficiency

        if (disperserAddress == address(0)) {
            revert IDisperserRegistry.InvalidDisperserAddress(disperserAddress);
        }

        if (depositParams.deposit > 0) {
            IERC20(depositParams.token).safeTransferFrom(msg.sender, address(this), depositParams.deposit);
            s().excess[depositParams.token] += depositParams.deposit - depositParams.refund; // we've already checked deposit >= refund
            disperser.deposit = depositParams;
            emit DisperserDepositTaken(disperserKey, depositParams);
        }

        disperser.disperser = disperserAddress;
        disperser.disperserURL = disperserURL;
        disperser.owner = msg.sender;

        emit DisperserRegistered(disperserKey, disperserAddress, msg.sender, disperserURL);

        return disperserKey;
    }

    /// @notice Transfers ownership of a disperser to a new owner.
    /// @param disperserKey The key of the disperser to transfer ownership of.
    /// @param newOwner The address of the new owner.
    function transferDisperserOwnership(uint32 disperserKey, address newOwner) internal {
        DisperserRegistryTypes.DisperserInfo storage disperser = s().disperser[disperserKey];

        if (disperser.disperser == address(0)) {
            revert IDisperserRegistry.DisperserNotRegistered(disperserKey);
        }
        if (newOwner == address(0)) {
            revert IDisperserRegistry.InvalidNewOwner(newOwner);
        }

        disperser.owner = newOwner;
        emit DisperserOwnershipTransferred(disperserKey, newOwner);
    }

    /// @notice Updates the disperser's address and URL. Takes an update fee if applicable.
    /// @param disperserKey The key of the disperser to update.
    /// @param disperserAddress The new address of the disperser.
    /// @param disperserURL The new URL of the disperser.
    function updateDisperserInfo(uint32 disperserKey, address disperserAddress, string memory disperserURL) internal {
        DisperserRegistryTypes.DisperserInfo storage disperser = s().disperser[disperserKey];

        if (disperser.disperser == address(0)) {
            revert IDisperserRegistry.DisperserNotRegistered(disperserKey);
        }
        if (disperserAddress == address(0)) {
            revert IDisperserRegistry.InvalidDisperserAddress(disperserAddress);
        }

        address token = s().depositParams.token;
        uint256 updateFee = s().updateFee;
        if (updateFee > 0) {
            IERC20(token).safeTransferFrom(msg.sender, address(this), updateFee);
            s().excess[token] += updateFee;
            emit DisperserUpdateFeeTaken(disperserKey, msg.sender, updateFee);
        }

        disperser.disperser = disperserAddress;
        disperser.disperserURL = disperserURL;
        emit DisperserUpdated(disperserKey, disperserAddress, disperserURL);
    }

    /// @notice Deregisters a disperser, marking it as inactive and setting an unlock timestamp for the deposit.
    /// @param disperserKey The key of the disperser to deregister.
    /// @dev The deposit can be withdrawn after the unlock timestamp has passed.
    ///      The disperser is marked as inactive by setting its address and URL to zero.
    ///      The unlock timestamp is set to the current block timestamp plus the lock period defined in the deposit parameters.
    ///      This function does not refund the deposit immediately; it must be withdrawn separately after the unlock period.
    function deregisterDisperser(uint32 disperserKey) internal {
        DisperserRegistryTypes.DisperserInfo storage disperser = s().disperser[disperserKey];
        DisperserRegistryTypes.LockedDisperserDeposit storage lockedDeposit = s().disperser[disperserKey].deposit;

        if (disperser.disperser == address(0)) {
            revert IDisperserRegistry.DisperserNotRegistered(disperserKey);
        }

        disperser.disperser = address(0);
        disperser.disperserURL = "";
        disperser.unlockTimestamp = uint64(block.timestamp) + lockedDeposit.lockPeriod;
        emit DisperserDeregistered(disperserKey, disperser.unlockTimestamp);
    }

    /// @notice Withdraws the deposit from a deregistered disperser.
    /// @param disperserKey The key of the disperser to withdraw from.
    /// @dev The deposit can only be withdrawn if the unlock timestamp has passed.
    ///      The refund amount is transferred to the disperser's owner.
    ///      The refund amount is set to zero after the transfer to prevent re-entrancy.
    function withdrawDisperserDeposit(uint32 disperserKey) internal {
        DisperserRegistryTypes.DisperserInfo storage disperser = s().disperser[disperserKey];
        DisperserRegistryTypes.LockedDisperserDeposit storage lockedDeposit = s().disperser[disperserKey].deposit;

        if (disperser.unlockTimestamp == 0) {
            revert IDisperserRegistry.DisperserNotDeregistered(disperserKey);
        }
        if (lockedDeposit.refund == 0) {
            revert IDisperserRegistry.ZeroRefund(disperserKey);
        }
        if (disperser.unlockTimestamp > block.timestamp) {
            revert IDisperserRegistry.RefundLocked(disperserKey, disperser.unlockTimestamp);
        }

        IERC20(lockedDeposit.token).safeTransfer(disperser.owner, lockedDeposit.refund);
        emit DisperserRefundIssued(disperserKey, lockedDeposit.token, lockedDeposit.refund);
        lockedDeposit.refund = 0;
        disperser.unlockTimestamp = 0;
    }

    /// @notice Sets the deposit parameters for disperser deposits.
    /// @param depositParams The deposit parameters to set.
    function setDepositParams(DisperserRegistryTypes.LockedDisperserDeposit memory depositParams) internal {
        if (depositParams.deposit < depositParams.refund) {
            revert IDisperserRegistry.DepositMustBeAtLeastRefund(depositParams.deposit, depositParams.refund);
        }
        if (depositParams.token == address(0)) {
            revert IDisperserRegistry.InvalidTokenAddress(depositParams.token);
        }

        s().depositParams = depositParams;
        emit DepositParamsSet(depositParams);
    }

    /// @notice Sets the update fee for dispersers.
    /// @param updateFee The update fee to set.
    function setUpdateFee(uint256 updateFee) internal {
        s().updateFee = updateFee;
        emit UpdateFeeSet(updateFee);
    }

    /// @notice Returns the current deposit parameters for dispersers.
    function getDepositParams() internal view returns (DisperserRegistryTypes.LockedDisperserDeposit memory) {
        return s().depositParams;
    }

    /// @notice Returns the deposit parameters for a specific disperser.
    /// @dev No check is performed to ensure the disperser is registered.
    function getDisperserDepositParams(uint32 disperserKey)
        internal
        view
        returns (DisperserRegistryTypes.LockedDisperserDeposit memory)
    {
        return s().disperser[disperserKey].deposit;
    }

    /// @notice Returns the address of a disperser by its key.
    /// @dev No check is performed to ensure the disperser is registered.
    function getDisperserAddress(uint32 disperserKey) internal view returns (address) {
        return s().disperser[disperserKey].disperser;
    }

    /// @notice Returns the owner of a disperser by its key.
    /// @dev No check is performed to ensure the disperser is registered.
    function getDisperserOwner(uint32 disperserKey) internal view returns (address) {
        return s().disperser[disperserKey].owner;
    }

    /// @notice Returns the unlock timestamp for a disperser's deposit.
    /// @dev No check is performed to ensure the disperser is registered.
    function getDisperserUnlockTime(uint32 disperserKey) internal view returns (uint64) {
        return s().disperser[disperserKey].unlockTimestamp;
    }

    /// @notice Returns the URL of a disperser by its key.
    /// @dev No check is performed to ensure the disperser is registered.
    function getDisperserURL(uint32 disperserKey) internal view returns (string memory) {
        return s().disperser[disperserKey].disperserURL;
    }

    /// @notice Returns the excess balance of a token held by the registry.
    function getExcessBalance(address token) internal view returns (uint256) {
        return s().excess[token];
    }

    /// @notice Returns the next disperser key that will be used for registration.
    function getNextDisperserKey() internal view returns (uint32) {
        return s().nextDisperserKey;
    }

    /// @notice Returns the update fee for dispersers.
    function getUpdateFee() internal view returns (uint256) {
        return s().updateFee;
    }
}
