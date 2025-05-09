// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierBase} from "src/periphery/cert/interfaces/IEigenDACertVerifierBase.sol";

interface IEigenDACertVerifier is IEigenDACertVerifierBase {
    /// @notice Returns the certificate version. Used off-chain to identify how to encode a certificate for this CertVerifier.
    function certVersion() external view returns (uint64);
}
