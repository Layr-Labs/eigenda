// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierBase} from "src/interfaces/IEigenDACertVerifier.sol";

interface IEigenDACertVerifierRouter is IEigenDACertVerifierBase {
    function addCertVerifier(uint32 referenceBlockNumber, address certVerifier) external;

    function getCertVerifierAt(uint32 rbn) external view returns (IEigenDACertVerifierBase);

    function certVerifiers(uint32 referenceBlockNumber) external view returns (IEigenDACertVerifierBase);

    function certVerifierRBNs(uint256 index) external view returns (uint32);
}
