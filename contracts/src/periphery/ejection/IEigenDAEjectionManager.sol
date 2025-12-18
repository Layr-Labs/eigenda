// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {BN254} from "lib/eigenlayer-middleware/src/libraries/BN254.sol";

interface IEigenDAEjectionManager {
    /// @notice Sets the delay for ejection processes.
    /// @param delay The number of seconds that must pass after initiation before an ejection can be completed.
    ///              This is also the time guaranteed to a challenger to cancel the ejection.
    function setDelay(uint64 delay) external;

    /// @notice Sets the cooldown for ejection processes.
    /// @param cooldown The number of seconds that must pass before a new ejection can be initiated after a previous one.
    function setCooldown(uint64 cooldown) external;

    /// @notice Starts the ejection process for an operator. Takes a deposit from the ejector.
    /// @param operator The address of the operator to eject.
    /// @param quorums The quorums associated with the ejection process.
    function startEjection(address operator, bytes memory quorums) external;

    /// @notice Cancels the ejection process initiated by a ejector.
    /// @dev Any ejector can cancel an ejection process, but the deposit is returned to the ejector who initiated it.
    function cancelEjectionByEjector(address operator) external;

    /// @notice Completes the ejection process for an operator. Transfers the deposit back to the ejector.
    /// @dev Any ejector can complete an ejection process, but the deposit is returned to the ejector who initiated it.
    function completeEjection(address operator, bytes memory quorums) external;

    /// @notice Cancels the ejection process for a given operator with their signature. Refunds the deposit to the recipient.
    /// @param operator The address of the operator whose ejection is being cancelled.
    /// @param apkG2 The G2 point of the operator's public key.
    /// @param sigma The BLS signature of the operator.
    /// @param recipient The address to which the gas refund will be sent.
    function cancelEjectionWithSig(
        address operator,
        BN254.G2Point memory apkG2,
        BN254.G1Point memory sigma,
        address recipient
    ) external;

    /// @notice Cancels the ejection process for the message sender. Refunds gas to the caller.
    function cancelEjection() external;

    /// @notice Returns the address of the ejector for a given operator. If the returned address is zero, then there is no ejection in progress.
    function getEjector(address operator) external view returns (address);

    /// @notice Returns whether an ejection process has been initiated for a given operator.
    function ejectionTime(address operator) external view returns (uint64);

    /// @notice Returns the timestamp of the last ejection proceeding initiated for a given operator.
    function lastEjectionInitiated(address operator) external view returns (uint64);

    /// @notice Returns the quorums associated with the ejection process for a given operator.
    function ejectionQuorums(address operator) external view returns (bytes memory);

    /// @notice Returns the delay for ejection processes.
    function ejectionDelay() external view returns (uint64);

    /// @notice Returns the cooldown for ejection initiations per operator.
    function ejectionCooldown() external view returns (uint64);
}
