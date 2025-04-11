// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

interface ISocketRegistry {
    /// @notice sets the socket for an operator only callable by the RegistryCoordinator
    function setOperatorSocket(bytes32 _operatorId, string memory _socket) external;

    /// @notice gets the stored socket for an operator
    function getOperatorSocket(bytes32 _operatorId) external view returns (string memory);
}
