// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;

import {IBLSApkRegistry} from "lib/eigenlayer-middleware/src/interfaces/IBLSApkRegistry.sol";
import {IStakeRegistry} from "lib/eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IIndexRegistry} from "lib/eigenlayer-middleware/src/interfaces/IIndexRegistry.sol";
import {IServiceManager} from "lib/eigenlayer-middleware/src/interfaces/IServiceManager.sol";
import {IRegistryCoordinator} from "lib/eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {ISocketRegistry} from "lib/eigenlayer-middleware/src/interfaces/ISocketRegistry.sol";
import {IEigenDAAddressDirectory} from "src/core/interfaces/IEigenDADirectory.sol";

abstract contract EigenDARegistryCoordinatorStorage is IRegistryCoordinator {
    /**
     *
     *                            CONSTANTS AND IMMUTABLES
     *
     */

    /// @notice The EIP-712 typehash for the `DelegationApproval` struct used by the contract
    bytes32 public constant OPERATOR_CHURN_APPROVAL_TYPEHASH = keccak256(
        "OperatorChurnApproval(address registeringOperator,bytes32 registeringOperatorId,OperatorKickParam[] operatorKickParams,bytes32 salt,uint256 expiry)OperatorKickParam(uint8 quorumNumber,address operator)"
    );
    /// @notice The EIP-712 typehash used for registering BLS public keys
    bytes32 public constant PUBKEY_REGISTRATION_TYPEHASH = keccak256("BN254PubkeyRegistration(address operator)");
    /// @notice The maximum value of a quorum bitmap
    uint256 internal constant MAX_QUORUM_BITMAP = type(uint192).max;
    /// @notice The basis point denominator
    uint16 internal constant BIPS_DENOMINATOR = 10000;
    /// @notice Index for flag that pauses operator registration
    uint8 internal constant PAUSED_REGISTER_OPERATOR = 0;
    /// @notice Index for flag that pauses operator deregistration
    uint8 internal constant PAUSED_DEREGISTER_OPERATOR = 1;
    /// @notice Index for flag pausing operator stake updates
    uint8 internal constant PAUSED_UPDATE_OPERATOR = 2;
    /// @notice The maximum number of quorums this contract supports
    uint8 internal constant MAX_QUORUM_COUNT = 192;

    IEigenDAAddressDirectory public immutable directory;

    /**
     *
     *                                    STATE
     *
     */

    /// @notice the current number of quorums supported by the registry coordinator
    uint8 public quorumCount;
    /// @notice maps quorum number => operator cap and kick params
    mapping(uint8 => OperatorSetParam) internal _quorumParams;
    /// @notice maps operator id => historical quorums they registered for
    mapping(bytes32 => QuorumBitmapUpdate[]) internal _operatorBitmapHistory;
    /// @notice maps operator address => operator id and status
    mapping(address => OperatorInfo) internal _operatorInfo;
    mapping(bytes32 => bool) private _deprecated_0;
    /// @notice mapping from quorum number to the latest block that all quorums were updated all at once
    mapping(uint8 => uint256) public quorumUpdateBlockNumber;

    address[] private _deprecated_2;
    address private _deprecated_1;
    /// @notice the address of the entity allowed to eject operators from the AVS
    address public ejector;

    /// @notice the last timestamp an operator was ejected
    mapping(address => uint256) public lastEjectionTimestamp;
    /// @notice the delay in seconds before an operator can reregister after being ejected
    uint256 public ejectionCooldown;

    constructor(address _directory) {
        directory = IEigenDAAddressDirectory(_directory);
    }

    // storage gap for upgradeability
    // slither-disable-next-line shadowing-state
    uint256[39] private __GAP;
}
