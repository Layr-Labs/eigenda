// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {ServiceManagerBase, IRewardsCoordinator, IServiceManager} from "../../lib/eigenlayer-middleware/src/ServiceManagerBase.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IServiceManagerRewardsRouter} from "./IServiceManagerRewardsRouter.sol";

/**
 * @title Entrypoint for pushing rewards submissions through the ServiceManager
 * @dev This contract should be set as the `rewardsInitiator` on the ServiceManager contract
**/
contract ServiceManagerRewardsRouter is OwnableUpgradeable, IServiceManagerRewardsRouter {
    using SafeERC20 for IERC20;

    /// @notice The ServiceManager contract with rewards interface
    IServiceManager public immutable serviceManager;

    /// @notice Mapping of addresses that are allowed to push rewards
    mapping(address => bool) public isAllowedToPushRewards;

    modifier onlyAllowedToPushRewards() {
        require(isAllowedToPushRewards[msg.sender], "ServiceManagerRewardsRouter: Address not permissioned to push rewards");
        _;
    }
 
    constructor(
        IServiceManager _serviceManager
    ) {
        serviceManager = _serviceManager;
    }

    function initialize(address _initialOwner) public initializer {
        _transferOwnership(_initialOwner);
    }

     /**
     * @notice Sets an addresses permission to push rewards
     * @param _address The address to set as allowed
     * @param _allowed Whether the address is allowed to push rewards
     */
    function setAllowedToPushRewards(address _address, bool _allowed) external onlyOwner {
        isAllowedToPushRewards[_address] = _allowed;
    }

    /**
     * @notice Creates a new rewards submission to the EigenLayer RewardsCoordinator contract, to be split amongst the
     * set of stakers delegated to operators who are registered to this `avs`
     * @param rewardsSubmissions The rewards submissions being created
     * @dev Only callable by a permissioned rewardsInitiator address
     * @dev The duration of the `rewardsSubmission` cannot exceed `MAX_REWARDS_DURATION`
     * @dev The tokens are sent to the `RewardsCoordinator` contract
     * @dev Strategies must be in ascending order of addresses to check for duplicates
     * @dev This function will revert if the `rewardsSubmission` is malformed,
     * e.g. if the `strategies` and `weights` arrays are of non-equal lengths
     */
    function createAVSRewardsSubmission(
        IRewardsCoordinator.RewardsSubmission[] calldata rewardsSubmissions
    ) public virtual onlyAllowedToPushRewards {

        for (uint256 i = 0; i < rewardsSubmissions.length; ++i) {
            // transfer token to ServiceManagerRewardsRouter and approve ServiceManager to transfer again in createAVSRewardsSubmission() call
            rewardsSubmissions[i].token.safeTransferFrom(
                msg.sender,
                address(this),
                rewardsSubmissions[i].amount
            );
            rewardsSubmissions[i].token.safeIncreaseAllowance(
                address(serviceManager),
                rewardsSubmissions[i].amount
            );
        }

        serviceManager.createAVSRewardsSubmission(rewardsSubmissions);
    }

    /**
     * @notice Creates a new operator-directed rewards submission, to be split amongst the operators and
     * set of stakers delegated to operators who are registered to this `avs`.
     * @param operatorDirectedRewardsSubmissions The operator-directed rewards submissions being created.
     * @dev Only callable by a permissioned rewardsInitiator address
     * @dev The duration of the `rewardsSubmission` cannot exceed `MAX_REWARDS_DURATION`
     * @dev The tokens are sent to the `RewardsCoordinator` contract
     * @dev This contract needs a token approval of sum of all `operatorRewards` in the `operatorDirectedRewardsSubmissions`, before calling this function.
     * @dev Strategies must be in ascending order of addresses to check for duplicates
     * @dev Operators must be in ascending order of addresses to check for duplicates.
     * @dev This function will revert if the `operatorDirectedRewardsSubmissions` is malformed.
     * @dev This function may fail to execute with a large number of submissions due to gas limits. Use a
     * smaller array of submissions if necessary.
     */
    function createOperatorDirectedAVSRewardsSubmission(
        IRewardsCoordinator.OperatorDirectedRewardsSubmission[]
            calldata operatorDirectedRewardsSubmissions
    ) public virtual onlyAllowedToPushRewards {

        for (
            uint256 i = 0;
            i < operatorDirectedRewardsSubmissions.length;
            ++i
        ) {
            // Calculate total amount of token to transfer
            uint256 totalAmount = 0;
            for (
                uint256 j = 0;
                j <
                operatorDirectedRewardsSubmissions[i].operatorRewards.length;
                ++j
            ) {
                totalAmount += operatorDirectedRewardsSubmissions[i]
                    .operatorRewards[j]
                    .amount;
            }

            // Transfer token to ServiceManagerRewardsRouter and approve ServiceManager to transfer again
            // in createOperatorDirectedAVSRewardsSubmission() call
            operatorDirectedRewardsSubmissions[i].token.safeTransferFrom(
                msg.sender,
                address(this),
                totalAmount
            );
            operatorDirectedRewardsSubmissions[i].token.safeIncreaseAllowance(
                address(serviceManager),
                totalAmount
            );
        }

        serviceManager.createOperatorDirectedAVSRewardsSubmission(operatorDirectedRewardsSubmissions);
    }
}