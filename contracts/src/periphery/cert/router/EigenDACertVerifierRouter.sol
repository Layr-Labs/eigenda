// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier} from "src/periphery/cert/interfaces/IEigenDACertVerifier.sol";
import {IEigenDACertVerifierRouter} from "src/periphery/cert/interfaces/IEigenDACertVerifierRouter.sol";
import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";

contract EigenDACertVerifierRouter is IEigenDACertVerifierRouter, OwnableUpgradeable {
    mapping(uint32 => address) public certVerifiers;
    uint32[] public certVerifierRBNs;

    event CertVerifierAdded(uint32 indexed rbn, address indexed certVerifier);

    error RBNNotInFuture(uint32 rbn);
    error RBNNotGreaterThanLast(uint32 rbn);
    error InvalidCertLength();
    error NoCertVerifierAvailable();
    error NoCertVerifierFound(uint32 rbn);

    /// IEigenDACertVerifierRouter ///

    function checkDACert(bytes calldata certBytes) external view returns (uint8) {
        return IEigenDACertVerifier(getCertVerifierAt(_getRBN(certBytes))).checkDACert(certBytes);
    }

    function getCertVerifierAt(uint32 rbn) public view returns (address) {
        return certVerifiers[_findClosestRegisteredRBN(rbn)];
    }

    /// ADMIN ///

    function initialize(address _initialOwner) external initializer {
        _transferOwnership(_initialOwner);
    }

    function addCertVerifier(uint32 rbn, address certVerifier) external onlyOwner {
        if (rbn <= block.number) {
            revert RBNNotInFuture(rbn);
        }
        if (certVerifierRBNs.length > 0 && rbn <= certVerifierRBNs[certVerifierRBNs.length - 1]) {
            revert RBNNotGreaterThanLast(rbn);
        }
        certVerifiers[rbn] = certVerifier;
        certVerifierRBNs.push(rbn);
        emit CertVerifierAdded(rbn, certVerifier);
    }

    /// INTERNAL ///

    function _getRBN(bytes calldata certBytes) internal pure returns (uint32) {
        if (certBytes.length < 36) {
            revert InvalidCertLength();
        }
        return abi.decode(certBytes[32:36], (uint32));
    }

    /// @notice Given an RBN, find the closest RBN registered in this contract that is less than or equal to the given RBN.
    /// @param referenceBlockNumber The reference block number to find the closest RBN for
    /// @return closestRBN The closest RBN registered in this contract that is less than or equal to the given RBN.
    function _findClosestRegisteredRBN(uint32 referenceBlockNumber) internal view returns (uint32) {
        // It is assumed that the latest RBNs are the most likely to be used.
        if (certVerifierRBNs.length == 0) {
            revert NoCertVerifierAvailable();
        }

        uint256 rbnMaxIndex = certVerifierRBNs.length - 1; // cache to memory
        for (uint256 i; i < certVerifierRBNs.length; i++) {
            uint32 certVerifierRBNMem = certVerifierRBNs[rbnMaxIndex - i];
            if (certVerifierRBNMem <= referenceBlockNumber) {
                return certVerifierRBNMem;
            }
        }
        revert NoCertVerifierFound(referenceBlockNumber);
    }
}
