// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDAThresholdRegistry} from "src/core/interfaces/IEigenDAThresholdRegistry.sol";
import {IEigenDASignatureVerifier} from "src/core/interfaces/IEigenDASignatureVerifier.sol";
import {IEigenDACertVerifierBase} from "src/periphery/cert/interfaces/IEigenDACertVerifierBase.sol";
import {EigenDATypesV1 as DATypesV1} from "src/core/libraries/v1/EigenDATypesV1.sol";

interface IEigenDACertVerifier is IEigenDACertVerifierBase {
    /// @notice Returns the EigenDA threshold registry contract.
    /// @return The EigenDA threshold registry contract.
    function eigenDAThresholdRegistry() external view returns (IEigenDAThresholdRegistry);

    /// @notice Returns the EigenDA signature verifier contract.
    /// @return The EigenDA signature verifier contract.
    function eigenDASignatureVerifier() external view returns (IEigenDASignatureVerifier);

    /// @notice Returns the security thresholds used for certificate verification.
    /// @return The security thresholds used for certificate verification.
    function securityThresholds() external view returns (DATypesV1.SecurityThresholds memory);

    /// @notice Returns the quorum numbers required for certificate verification. All required quorums must meet the stake threshold for a certificate.
    /// @return The quorum numbers required for certificate verification in bytes format.
    function quorumNumbersRequired() external view returns (bytes memory);

    /// @notice Returns the EigenDA certificate version. Used off-chain to identify how to encode a certificate for this CertVerifier.
    /// @return The EigenDA certificate version.
    function certVersion() external view returns (uint64);
}
