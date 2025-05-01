// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IEigenDACertVerifier {
    /// @notice Check a DA cert's validity, and revert if invalid.
    function verifyDACert(bytes calldata certBytes) external view;

    /// @notice Check a DA cert's validity
    /// @return status An enum value. Success is always mapped to 1, and other values are errors specific to each CertVerifier.
    function checkDACert(bytes calldata certBytes) external view returns (uint8 status);

    /// @notice Returns the certificate version. Used off-chain to identify how to encode a certificate for this CertVerifier.
    function CERT_VERSION() external view returns (uint8);
}
