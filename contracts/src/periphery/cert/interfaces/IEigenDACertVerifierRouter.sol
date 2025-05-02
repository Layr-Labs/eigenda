// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierBase} from "src/periphery/cert/interfaces/IEigenDACertVerifierBase.sol";

interface IEigenDACertVerifierRouter is IEigenDACertVerifierBase {
    /// @notice Returns the address for a cert verifier at a given ABN.
    function getCertVerifierAt(uint32 abn) external view returns (address);
}
