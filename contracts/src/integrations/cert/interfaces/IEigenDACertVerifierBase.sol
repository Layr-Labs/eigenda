// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IEigenDACertVerifierBase {
    /// @notice Check a DA cert's validity
    /// @param abiEncodedCert The ABI encoded certificate. Any cert verifier should decode this ABI encoding based on the certificate version.
    /// @return status An enum value. Success is always mapped to 1, and other values are errors specific to each CertVerifier.
    /// @dev This function should never revert, because it needs to be proven using Risc0's Steel library, which doesn't support reverts.
    //       See https://github.com/risc0/risc0-ethereum/issues/438
    function checkDACert(bytes calldata abiEncodedCert) external view returns (uint8 status);
}
