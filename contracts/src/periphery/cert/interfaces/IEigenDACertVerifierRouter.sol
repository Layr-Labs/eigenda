// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierBase} from "src/periphery/cert/interfaces/IEigenDACertVerifierBase.sol";

interface IEigenDACertVerifierRouter is IEigenDACertVerifierBase {
    /// @notice Returns the address for the active cert verifier at a given reference block number.
    ///         The reference block number must not be in the future.
    function getCertVerifierAt(uint32 referenceBlockNumber) external view returns (address);
}
