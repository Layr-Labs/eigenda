// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

interface IVersionedDerivationEigenDACertVerifier {
    /// @notice Returns the offchain derivation version used in certificate verification.
    function offchainDerivationVersion() external view returns (uint16);
}