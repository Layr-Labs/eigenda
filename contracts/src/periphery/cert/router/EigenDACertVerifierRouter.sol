// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifier} from "src/periphery/cert/interfaces/IEigenDACertVerifier.sol";
import {IEigenDACertVerifierRouter} from "src/periphery/cert/interfaces/IEigenDACertVerifierRouter.sol";
import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";

contract EigenDACertVerifierRouter is IEigenDACertVerifierRouter, OwnableUpgradeable {
    /// @notice A mapping from an activation block number (ABN) to a cert verifier address.
    mapping(uint32 => address) public certVerifiers;

    /// @notice The list of Activation Block Numbers (ABNs) for the cert verifiers.
    /// @dev The list is sorted in ascending order, and corresponds to the keys of the certVerifiers mapping.
    uint32[] public certVerifierABNs;

    event CertVerifierAdded(uint32 indexed abn, address indexed certVerifier);

    error ABNNotInFuture(uint32 abn);
    error ABNNotGreaterThanLast(uint32 abn);
    error InvalidCertLength();
    error NoCertVerifierFound(uint32 rbn);

    /// IEigenDACertVerifierRouter ///

    function checkDACert(bytes calldata certBytes) external view returns (uint8) {
        return IEigenDACertVerifier(getCertVerifierAt(_getRBN(certBytes[32:]))).checkDACert(certBytes);
    }

    function getCertVerifierAt(uint32 rbn) public view returns (address) {
        return certVerifiers[_findPrecedingRegisteredABN(rbn)];
    }

    /// ADMIN ///

    function initialize(address _initialOwner, address certVerifier) external initializer {
        _transferOwnership(_initialOwner);
        // Add a default cert verifier at block 0, which will be used for all blocks before the first ABN.
        _addCertVerifier(0, certVerifier);
        emit CertVerifierAdded(0, certVerifier);
    }

    function addCertVerifier(uint32 abn, address certVerifier) external onlyOwner {
        if (abn <= block.number) {
            revert ABNNotInFuture(abn);
        }
        if (abn <= certVerifierABNs[certVerifierABNs.length - 1]) {
            revert ABNNotGreaterThanLast(abn);
        }
        _addCertVerifier(abn, certVerifier);
    }

    /// INTERNAL ///

    function _addCertVerifier(uint32 abn, address certVerifier) internal {
        certVerifiers[abn] = certVerifier;
        certVerifierABNs.push(abn);
        emit CertVerifierAdded(abn, certVerifier);
    }

    function _getRBN(bytes calldata certBytes) internal pure returns (uint32) {
        // 0:32 is the batch header root
        // 32:64 is the RBN
        if (certBytes.length < 64) {
            revert InvalidCertLength();
        }
        return abi.decode(certBytes[32:64], (uint32));
    }

    /// @notice Given a reference block number, find the closest activation block number
    ///         registered in this contract that is less than or equal to the given reference block number.
    /// @param rbn The reference block number to find the closest ABN for
    /// @return abn The preceding ABN registered in this contract that is less than or equal to the given ABN.
    function _findPrecedingRegisteredABN(uint32 rbn) internal view returns (uint32 abn) {
        // It is assumed that the latest ABN are the most likely to be used.
        uint256 abnMaxIndex = certVerifierABNs.length - 1; // cache to memory
        for (uint256 i; i < certVerifierABNs.length; i++) {
            abn = certVerifierABNs[abnMaxIndex - i];
            if (abn <= rbn) {
                return abn;
            }
        }
    }
}
