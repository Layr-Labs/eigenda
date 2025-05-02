// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier} from "src/periphery/cert/interfaces/IEigenDACertVerifier.sol";
import {IEigenDACertVerifierRouter} from "src/periphery/cert/interfaces/IEigenDACertVerifierRouter.sol";
import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";

contract EigenDACertVerifierRouter is IEigenDACertVerifierRouter, OwnableUpgradeable {
    mapping(uint32 => address) public certVerifiers;
    uint32[] public certVerifierABNs;

    event CertVerifierAdded(uint32 indexed rbn, address indexed certVerifier);

    error ABNNotInFuture(uint32 abn);
    error ABNNotGreaterThanLast(uint32 abn);
    error InvalidCertLength();
    error NoCertVerifierAvailable();
    error NoCertVerifierFound(uint32 rbn);

    /// IEigenDACertVerifierRouter ///

    function checkDACert(bytes calldata certBytes) external view returns (uint8) {
        return IEigenDACertVerifier(getCertVerifierAt(_getRBN(certBytes))).checkDACert(certBytes);
    }

    function getCertVerifierAt(uint32 abn) public view returns (address) {
        return certVerifiers[_findClosestRegisteredABN(abn)];
    }

    /// ADMIN ///

    function initialize(address _initialOwner) external initializer {
        _transferOwnership(_initialOwner);
    }

    function addCertVerifier(uint32 abn, address certVerifier) external onlyOwner {
        if (abn <= block.number) {
            revert ABNNotInFuture(abn);
        }
        if (certVerifierABNs.length > 0 && abn <= certVerifierABNs[certVerifierABNs.length - 1]) {
            revert ABNNotGreaterThanLast(abn);
        }
        certVerifiers[abn] = certVerifier;
        certVerifierABNs.push(abn);
        emit CertVerifierAdded(abn, certVerifier);
    }

    /// INTERNAL ///

    function _getRBN(bytes calldata certBytes) internal pure returns (uint32) {
        // 0:32 is the data offset
        // 32:64 is the batch header root
        // 64:96 is the RBN
        if (certBytes.length < 96) {
            revert InvalidCertLength();
        }
        return abi.decode(certBytes[64:96], (uint32));
    }

    /// @notice Given a reference block number, find the closest activation block number
    ///         registered in this contract that is less than or equal to the given reference block number.
    /// @param rbn The reference block number to find the closest ABN for
    /// @return The closest ABN registered in this contract that is less than or equal to the given ABN.
    function _findClosestRegisteredABN(uint32 rbn) internal view returns (uint32) {
        // It is assumed that the latest ABN are the most likely to be used.
        if (certVerifierABNs.length == 0) {
            revert NoCertVerifierAvailable();
        }

        uint256 abnMaxIndex = certVerifierABNs.length - 1; // cache to memory
        for (uint256 i; i < certVerifierABNs.length; i++) {
            uint32 certVerifierABNMem = certVerifierABNs[abnMaxIndex - i];
            if (certVerifierABNMem <= rbn) {
                return certVerifierABNMem;
            }
        }
        revert NoCertVerifierFound(rbn);
    }
}
