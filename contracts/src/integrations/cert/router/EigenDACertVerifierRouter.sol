// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {IEigenDACertVerifierBase} from "src/integrations/cert/interfaces/IEigenDACertVerifierBase.sol";
import {IVersionedEigenDACertVerifier} from "src/integrations/cert/interfaces/IVersionedEigenDACertVerifier.sol";
import {IEigenDACertVerifierRouter} from "src/integrations/cert/interfaces/IEigenDACertVerifierRouter.sol";
import {OwnableUpgradeable} from "lib/openzeppelin-contracts-upgradeable/contracts/access/OwnableUpgradeable.sol";

contract EigenDACertVerifierRouter is IEigenDACertVerifierRouter, OwnableUpgradeable {
    /// @notice A mapping from an activation block number (ABN) to a cert verifier address.
    mapping(uint32 => address) public certVerifiers;

    /// @notice The list of Activation Block Numbers (ABNs) for the cert verifiers.
    /// @dev The list is guaranteed to be in ascending order
    ///      and corresponds to the keys of the certVerifiers mapping.
    uint32[] public certVerifierABNs;

    event CertVerifierAdded(uint32 indexed activationBlockNumber, address indexed certVerifier);

    /// @notice Thrown when an activation block number (ABN) is not in the future when adding a cert verifier.
    error ABNNotInFuture(uint32 activationBlockNumber);
    /// @notice Thrown when an activation block number (ABN) is not greater than the last registered ABN when adding a cert verifier.
    error ABNNotGreaterThanLast(uint32 activationBlockNumber);
    /// @notice Thrown when the length of the cert verifier is not valid.
    error InvalidCertLength();
    /// @notice Thrown when a reference block number (RBN) is in the future when verifying a cert.
    error RBNInFuture(uint32 referenceBlockNumber);
    /// @notice Thrown when the length of input arrays that are expected to match do not match.
    error LengthMismatch();

    /// IEigenDACertVerifierRouter ///

    /// @inheritdoc IEigenDACertVerifierBase
    function checkDACert(bytes calldata abiEncodedCert) external view returns (uint8) {
        return IEigenDACertVerifierBase(getCertVerifierAt(_getRBN(abiEncodedCert))).checkDACert(abiEncodedCert);
    }

    function getCertVerifierAt(uint32 referenceBlockNumber) public view returns (address) {
        return certVerifiers[_findPrecedingRegisteredABN(referenceBlockNumber)];
    }

    /// ADMIN ///

    /// @notice Initializes the EigenDACertVerifierRouter.
    /// @param initialOwner The owner can add new cert verifiers. See addCertVerifier for security implications.
    /// @param initABNs A list of ABNs that will be initialized with cert verifiers
    /// @param initCertVerifiers A list of cert verifiers corresponding to initABNs.
    function initialize(address initialOwner, uint32[] memory initABNs, address[] memory initCertVerifiers)
        external
        initializer
    {
        _transferOwnership(initialOwner);
        if (initABNs.length != initCertVerifiers.length) {
            revert LengthMismatch();
        }
        // Add the first cert verifier. Because the first ABN might be zero, the initABN check cannot happen inside the loop with a naive implementation.
        uint256 lastABN;
        for (uint256 i; i < initABNs.length; i++) {
            if (initABNs[i] <= lastABN && i > 0) {
                revert ABNNotGreaterThanLast(initABNs[i]);
            }
            lastABN = initABNs[i];
            _addCertVerifier(initABNs[i], initCertVerifiers[i]);
        }
    }

    /// @notice Adds a cert verifier to the router.
    /// @param activationBlockNumber The block number at which the cert verifier will be activated. Must be in the future.
    /// @param certVerifier The address of the cert verifier to be added.
    /// @dev EigenDA recommends that a mechanism be implemented to ensure a cert verifier cannot be added too close to the current block number.
    ///      This is to prevent a malicious party from setting a cert verifier without enough time for other parties to react.
    ///      This could be a timelock, multisig transaction restriction on activationBlockNumber, delay, ownerless contract, etc..
    function addCertVerifier(uint32 activationBlockNumber, address certVerifier) external onlyOwner {
        // We disallow adding cert verifiers at the current block number to avoid a race condition of
        // adding a cert verifier at the current block and verifying in the same block
        if (activationBlockNumber <= block.number) {
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
        if (certBytes.length < 96) {
            revert InvalidCertLength();
        }
        return abi.decode(certBytes[64:96], (uint32));
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
        if (referenceBlockNumber > block.number) {
            revert RBNInFuture(referenceBlockNumber);
        }
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
