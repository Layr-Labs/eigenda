// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IEigenDACertVerifierBase {
    /// @notice Check a DA cert's validity
    /// @return status An enum value. Success is always mapped to 1, and other values are errors specific to each CertVerifier.
    function checkDACert(bytes calldata certBytes) external view returns (uint8 status);
}
