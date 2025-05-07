// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {EigenDATypesV3} from "src/core/libraries/v3/EigenDATypesV3.sol";
import {SafeERC20, IERC20} from "lib/openzeppelin-contracts/contracts/token/ERC20/utils/SafeERC20.sol";

library DisperserRegistryStorage {
    string internal constant STORAGE_ID = "eigen.da.disperser.registry";
    bytes32 internal constant STORAGE_POSITION =
        keccak256(abi.encode(uint256(keccak256(abi.encodePacked(STORAGE_ID))) - 1)) & ~bytes32(uint256(0xff));

    struct Disperser {
        EigenDATypesV3.DisperserInfo info;
        EigenDATypesV3.LockedDisperserDeposit deposit;
        uint64 unlockTimestamp;
    }

    struct Layout {
        mapping(uint32 => Disperser) disperser;
        EigenDATypesV3.LockedDisperserDeposit depositParams;
        uint32 nextDisperserKey;
    }

    function layout() internal pure returns (Layout storage s) {
        bytes32 position = STORAGE_POSITION;
        assembly {
            s.slot := position
        }
    }
}

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
        DisperserRegistryStorage.Disperser storage disperser = s().disperser[disperserKey];

        require(disperserAddress != address(0), "Invalid disperser address");

        IERC20(s().depositParams.token).safeTransferFrom(msg.sender, address(this), s().depositParams.deposit);

        disperser.info =
            EigenDATypesV3.DisperserInfo({disperser: disperserAddress, registered: true, disperserURL: disperserURL});
        disperser.deposit = s().depositParams;
        disperser.unlockTimestamp = type(uint64).max;

        return disperserKey;
    }

    function deregisterDisperser(uint32 disperserKey) internal {
        DisperserRegistryStorage.Disperser storage disperser = s().disperser[disperserKey];
        EigenDATypesV3.LockedDisperserDeposit storage lockedDeposit = s().disperser[disperserKey].deposit;

        require(disperser.info.registered, "Disperser not registered");

        disperser.info.registered = false;
        disperser.unlockTimestamp = uint64(block.timestamp) + lockedDeposit.lockPeriod;
    }

    function withdraw(uint32 disperserKey) internal {
        DisperserRegistryStorage.Disperser storage disperser = s().disperser[disperserKey];
        EigenDATypesV3.LockedDisperserDeposit storage lockedDeposit = s().disperser[disperserKey].deposit;

        require(lockedDeposit.refund > 0, "No deposit to withdraw");
        require(disperser.unlockTimestamp <= block.timestamp, "Deposit is still locked");

        IERC20(lockedDeposit.token).safeTransfer(disperser.info.disperser, lockedDeposit.refund);
        lockedDeposit.refund = 0;
    }

    function setDepositParams(EigenDATypesV3.LockedDisperserDeposit memory depositParams) internal {
        require(depositParams.deposit >= depositParams.refund, "Deposit must be at least refund");
        require(depositParams.token != address(0), "Invalid token address");

        s().depositParams = depositParams;
    }

    function getDisperserInfo(uint32 disperserKey) internal view returns (EigenDATypesV3.DisperserInfo memory) {
        return s().disperser[disperserKey].info;
    }

    function getDisperser(uint32 disperserKey) internal view returns (address) {
        return s().disperser[disperserKey].info.disperser;
    }

    function getLockedDeposit(uint32 disperserKey)
        internal
        view
        returns (EigenDATypesV3.LockedDisperserDeposit memory, uint64 unlockTimestamp)
    {
        DisperserRegistryStorage.Disperser storage disperser = s().disperser[disperserKey];
        return (disperser.deposit, disperser.unlockTimestamp);
    }

    function getDepositParams() internal view returns (EigenDATypesV3.LockedDisperserDeposit memory) {
        return s().depositParams;
    }
}
