// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

library DisperserRegistryTypes {
    /// @notice Represents the storage layout for a particular disperser.
    /// @param disperser The address of the disperser.
    /// @param owner The address of the owner of the disperser.
    /// @param unlockTimestamp The time when the disperser's deposit can be retrieved.
    /// @param disperserURL The disperser's URL.
    /// @param deposit The registered parameters of a disperser's deposit at the time of registration.
    struct DisperserInfo {
        address disperser;
        address owner;
        uint64 unlockTimestamp;
        string disperserURL;
        LockedDisperserDeposit deposit;
    }

    /// @notice Represents the parameters of a disperser deposit set on registration.
    /// @param deposit The amount of the deposit.
    /// @param refund The amount to be refunded after the lock period.
    /// @param token The address of the token used for the deposit.
    /// @param lockPeriod The duration for which the deposit is locked after deregistration.
    struct LockedDisperserDeposit {
        uint256 deposit;
        uint256 refund;
        address token;
        uint64 lockPeriod;
    }
}
