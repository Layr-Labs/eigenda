// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IEigenDASemVer {
    /// @notice Returns the semantic version of the contract implementation. Refer to https://semver.org/
    function semver() external view returns (uint8 major, uint8 minor, uint8 patch);
}
