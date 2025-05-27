// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierBase} from "src/periphery/cert/interfaces/IEigenDACertVerifierBase.sol";
import {IVersionedEigenDACertVerifier} from "src/periphery/cert/interfaces/IVersionedEigenDACertVerifier.sol";

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";

import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

/// @notice The IEigenDACertVerifier interface provides the getters necessary to transform a BlobStatusReply received after a dispersal
///         into a Blob Certificate that can be verified by the EigenDACertVerifier that implements this interface version.
// IEigenDACertVerifier provides the getters necessary to transform a BlobStatusReply received after a dispersal into a Cert that can be verified by the EigenDACertVerifier that implements this interface version.
interface IEigenDACertVerifier {
    /// @notice Returns the EigenDAThresholdRegistry contract.
    function eigenDAThresholdRegistry() external view returns (IEigenDAThresholdRegistry);

    /// @notice Returns the EigenDASignatureVerifier contract.
    function eigenDASignatureVerifier() external view returns (IEigenDASignatureVerifier);

    /// @notice Returns the security thresholds required for EigenDA certificate verification.
    function securityThresholds() external view returns (DATypesV1.SecurityThresholds memory);

    /// @notice Returns the quorum numbers required in bytes format for certificate verification.
    function quorumNumbersRequired() external view returns (bytes memory);
}
