// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";
import {DisperserRegistryTypes} from "src/core/libraries/v3/disperser/DisperserRegistryTypes.sol";
import {DisperserRegistryStorage} from "src/core/libraries/v3/disperser/DisperserRegistryStorage.sol";

library DisperserRegistryLib {
    using SafeERC20 for IERC20;

    function s() internal pure returns (DisperserRegistryStorage.Layout storage) {
        return DisperserRegistryStorage.layout();
    }

    function consumeDisperserKey() internal returns (uint32) {
        uint32 disperserKey = s().nextDisperserKey;
        s().nextDisperserKey++;
        return disperserKey;
    }

    function registerDisperser(address disperserAddress, string memory disperserURL)
        internal
        returns (uint32 disperserKey)
    {
        disperserKey = consumeDisperserKey();
        DisperserRegistryTypes.DisperserInfo storage disperser = s().disperser[disperserKey];
        DisperserRegistryTypes.LockedDisperserDeposit memory depositParams = s().depositParams; // memory copy for gas efficiency

        require(disperserAddress != address(0), "Invalid disperser address");

        if (depositParams.deposit > 0) {
            IERC20(depositParams.token).safeTransferFrom(msg.sender, address(this), depositParams.deposit);
            s().excess[depositParams.token] += depositParams.deposit - depositParams.refund; // we've already checked deposit >= refund
        }

        disperser.disperser = disperserAddress;
        disperser.disperserURL = disperserURL;
        disperser.deposit = depositParams;
        disperser.unlockTimestamp = type(uint64).max;
        disperser.owner = msg.sender;

        return disperserKey;
    }

    function transferDisperserOwnership(uint32 disperserKey, address newOwner) internal {
        DisperserRegistryTypes.DisperserInfo storage disperser = s().disperser[disperserKey];

        require(disperser.disperser != address(0), "Disperser not registered");
        require(newOwner != address(0), "Invalid new owner");

        disperser.owner = newOwner;
    }

    function updateDisperserInfo(uint32 disperserKey, address disperserAddress, string memory disperserURL) internal {
        DisperserRegistryTypes.DisperserInfo storage disperser = s().disperser[disperserKey];

        require(disperser.disperser != address(0), "Disperser not registered");
        require(disperserAddress != address(0), "Invalid disperser address");

        address token = s().depositParams.token;
        uint256 updateFee = s().updateFee;
        if (updateFee > 0) {
            IERC20(token).safeTransferFrom(msg.sender, address(this), updateFee);
            s().excess[token] += updateFee;
        }

        disperser.disperser = disperserAddress;
        disperser.disperserURL = disperserURL;
    }

    function deregisterDisperser(uint32 disperserKey) internal {
        DisperserRegistryTypes.DisperserInfo storage disperser = s().disperser[disperserKey];
        DisperserRegistryTypes.LockedDisperserDeposit storage lockedDeposit = s().disperser[disperserKey].deposit;

        require(disperser.disperser != address(0), "Disperser not registered");

        disperser.disperser = address(0);
        disperser.disperserURL = "";
        disperser.unlockTimestamp = uint64(block.timestamp) + lockedDeposit.lockPeriod;
    }

    function withdraw(uint32 disperserKey) internal {
        DisperserRegistryTypes.DisperserInfo storage disperser = s().disperser[disperserKey];
        DisperserRegistryTypes.LockedDisperserDeposit storage lockedDeposit = s().disperser[disperserKey].deposit;

        require(lockedDeposit.refund > 0, "No deposit to withdraw");
        require(disperser.unlockTimestamp <= block.timestamp, "Deposit is still locked");

        IERC20(lockedDeposit.token).safeTransfer(disperser.owner, lockedDeposit.refund);
        lockedDeposit.refund = 0;
    }

    function setDepositParams(DisperserRegistryTypes.LockedDisperserDeposit memory depositParams) internal {
        require(depositParams.deposit >= depositParams.refund, "Deposit must be at least refund");
        require(depositParams.token != address(0), "Invalid token address");

        s().depositParams = depositParams;
    }

    function getDepositParams() internal view returns (DisperserRegistryTypes.LockedDisperserDeposit memory) {
        return s().depositParams;
    }

    function getDisperserAddress(uint32 disperserKey) internal view returns (address) {
        return s().disperser[disperserKey].disperser;
    }

    function getDisperserOwner(uint32 disperserKey) internal view returns (address) {
        return s().disperser[disperserKey].owner;
    }

    function getDisperserUnlockTimestamp(uint32 disperserKey) internal view returns (uint64) {
        return s().disperser[disperserKey].unlockTimestamp;
    }

    function getDisperserURL(uint32 disperserKey) internal view returns (string memory) {
        return s().disperser[disperserKey].disperserURL;
    }
}
