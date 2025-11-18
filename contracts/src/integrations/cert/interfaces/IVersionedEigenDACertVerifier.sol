// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IVersionedEigenDACertVerifier {
    /// @notice Returns the EigenDA certificate version. Used off-chain to identify how to encode a certificate for this CertVerifier.
    /// @return The EigenDA certificate version.
    function certVersion() external view returns (uint8);
}
