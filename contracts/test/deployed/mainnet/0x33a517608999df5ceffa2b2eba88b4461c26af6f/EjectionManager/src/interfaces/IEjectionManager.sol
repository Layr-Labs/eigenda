// SPDX-License-Identifier: MIT
pragma solidity ^0.8.12;

/**
 * @title Interface for a contract that ejects operators from an AVSs RegistryCoordinator
 * @author Layr Labs, Inc.
 */
interface IEjectionManager {

    /// @notice A quorum's ratelimit parameters
    struct QuorumEjectionParams {
        uint32 rateLimitWindow; // Time delta to track ejection over
        uint16 ejectableStakePercent; // Max stake to be ejectable per time delta
    }

    /// @notice A stake ejection event
    struct StakeEjection {
        uint256 timestamp; // Timestamp of the ejection
        uint256 stakeEjected; // Amount of stake ejected at the timestamp
    }

    ///@notice Emitted when the ejector address is set
    event EjectorUpdated(address ejector, bool status);
    ///@notice Emitted when the ratelimit parameters for a quorum are set
    event QuorumEjectionParamsSet(uint8 quorumNumber, uint32 rateLimitWindow, uint16 ejectableStakePercent);
    ///@notice Emitted when an operator is ejected
    event OperatorEjected(bytes32 operatorId, uint8 quorumNumber);
    ///@notice Emitted when operators are ejected for a quroum 
    event QuorumEjection(uint32 ejectedOperators, bool ratelimitHit);

   /**
     * @notice Ejects operators from the AVSs registryCoordinator under a ratelimit
     * @param _operatorIds The ids of the operators to eject for each quorum
     */
    function ejectOperators(bytes32[][] memory _operatorIds) external;

    /**
     * @notice Sets the ratelimit parameters for a quorum
     * @param _quorumNumber The quorum number to set the ratelimit parameters for
     * @param _quorumEjectionParams The quorum ratelimit parameters to set for the given quorum
     */
    function setQuorumEjectionParams(uint8 _quorumNumber, QuorumEjectionParams memory _quorumEjectionParams) external;

    /**
     * @notice Sets the address permissioned to eject operators under a ratelimit
     * @param _ejector The address to permission
     */
    function setEjector(address _ejector, bool _status) external;

    /**
     * @notice Returns the amount of stake that can be ejected for a quorum at the current block.timestamp
     * @param _quorumNumber The quorum number to view ejectable stake for
     */
    function amountEjectableForQuorum(uint8 _quorumNumber) external view returns (uint256);
}
