// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {IEigenDACertVerifierBase} from "src/periphery/cert/interfaces/IEigenDACertVerifierBase.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

interface IVersionedEigenDACertVerifier {
    /// @notice Returns the EigenDA certificate version. Used off-chain to identify how to encode a certificate for this CertVerifier.
    /// @return The EigenDA certificate version.
    function certVersion() external view returns (uint64);
}
