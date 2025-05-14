// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierBase, IEigenDACertVerifier} from "src/periphery/cert/interfaces/IEigenDACertVerifier.sol";
import {IEigenDACertVerifierRouter} from "src/periphery/cert/interfaces/IEigenDACertVerifierRouter.sol";
import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";

contract EigenDACertVerifierRouter is IEigenDACertVerifierRouter, OwnableUpgradeable {
    /// @notice A mapping from an activation block number (ABN) to a cert verifier address.
    mapping(uint32 => address) public certVerifiers;

    /// @notice The list of Activation Block Numbers (ABNs) for the cert verifiers.
    /// @dev The list is guaranteed to be in ascending order
    ///      and corresponds to the keys of the certVerifiers mapping.
    uint32[] public certVerifierABNs;

    event CertVerifierAdded(uint32 indexed activationBlockNumber, address indexed certVerifier);

    error ABNNotInFuture(uint32 activationBlockNumber);
    error ABNNotGreaterThanLast(uint32 activationBlockNumber);
    error InvalidCertLength();

    /// IEigenDACertVerifierRouter ///

    /// @inheritdoc IEigenDACertVerifierBase
    function checkDACert(bytes calldata abiEncodedCert) external view returns (uint8) {
        return IEigenDACertVerifier(getCertVerifierAt(_getRBN(abiEncodedCert))).checkDACert(abiEncodedCert);
    }

    function getCertVerifierAt(uint32 referenceBlockNumber) public view returns (address) {
        return certVerifiers[_findPrecedingRegisteredABN(referenceBlockNumber)];
    }

    /// ADMIN ///

    function initialize(address _initialOwner, address certVerifier) external initializer {
        _transferOwnership(_initialOwner);
        // Add a default cert verifier at block 0, which will be used for all blocks before the first ABN.
        _addCertVerifier(0, certVerifier);
    }

    function addCertVerifier(uint32 activationBlockNumber, address certVerifier) external onlyOwner {
        if (activationBlockNumber < block.number) {
            revert ABNNotInFuture(activationBlockNumber);
        }
        if (activationBlockNumber <= certVerifierABNs[certVerifierABNs.length - 1]) {
            revert ABNNotGreaterThanLast(activationBlockNumber);
        }
        _addCertVerifier(activationBlockNumber, certVerifier);
    }

    /// INTERNAL ///

    function _addCertVerifier(uint32 activationBlockNumber, address certVerifier) internal {
        certVerifiers[activationBlockNumber] = certVerifier;
        certVerifierABNs.push(activationBlockNumber);
        emit CertVerifierAdded(activationBlockNumber, certVerifier);
    }

    function _getRBN(bytes calldata certBytes) internal pure returns (uint32) {
        // 0:32 is the pointer to the start of the byte array.
        // 32:64 is the batch header root
        // 64:96 is the RBN
        if (certBytes.length < 64) {
            revert InvalidCertLength();
        }
        return abi.decode(certBytes[32:64], (uint32));
    }

    /// @notice Given a reference block number, find the closest activation block number
    ///         registered in this contract that is less than or equal to the given reference block number.
    /// @param referenceBlockNumber The reference block number to find the closest ABN for
    /// @return activationBlockNumber The preceding ABN registered in this contract that is less than or equal to the given ABN.
    function _findPrecedingRegisteredABN(uint32 referenceBlockNumber)
        internal
        view
        returns (uint32 activationBlockNumber)
    {
        // It is assumed that the latest ABN are the most likely to be used.
        uint256 abnMaxIndex = certVerifierABNs.length - 1; // cache to memory
        for (uint256 i; i < certVerifierABNs.length; i++) {
            activationBlockNumber = certVerifierABNs[abnMaxIndex - i];
            if (activationBlockNumber <= referenceBlockNumber) {
                return activationBlockNumber;
            }
        }
    }
}
