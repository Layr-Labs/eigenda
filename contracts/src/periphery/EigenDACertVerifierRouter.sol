// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier, IEigenDACertVerifierBase} from "src/interfaces/IEigenDACertVerifier.sol";

contract EigenDACertVerifierRouter is IEigenDACertVerifierBase {
    mapping(uint32 => address) public certVerifiers;

    function verifyDACert(bytes calldata certBytes) external view {
        uint32 rbn = getRBN(certBytes);
        address certVerifier = certVerifiers[rbn];
        require(certVerifier != address(0), "Cert verifier not found");
        IEigenDACertVerifier(certVerifier).verifyDACert(certBytes);
    }

    function checkDACert(bytes calldata certBytes) external view returns (uint8) {
        uint32 rbn = getRBN(certBytes);
        address certVerifier = certVerifiers[rbn];
        require(certVerifier != address(0), "Cert verifier not found");
        return IEigenDACertVerifier(certVerifier).checkDACert(certBytes);
    }

    function getRBN(bytes calldata certBytes) internal pure returns (uint32) {
        return abi.decode(certBytes[32:36], (uint32));
    }
}
