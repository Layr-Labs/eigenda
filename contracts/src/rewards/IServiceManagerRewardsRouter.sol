// SPDX-License-Identifier: MIT
pragma solidity ^0.8.9;

import {ServiceManagerBase, IRewardsCoordinator, IServiceManager} from "../../lib/eigenlayer-middleware/src/ServiceManagerBase.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/**
 * @title Interface for ServiceManagerRewardsRouter which is used to push rewards submissions through the ServiceManager
**/
interface IServiceManagerRewardsRouter {

    /**
     * @notice Returns whether the address is allowed to push rewards
     * @param _address The address to check
     * @return bool true if the address is allowed to push rewards, false otherwise
     */
    function isAllowedToPushRewards(address _address) external view returns (bool);

     /**
     * @notice Sets an addresses permission to push rewards
     * @param _address The address to set as allowed
     * @param _allowed Whether the address is allowed to push rewards
     */
    function setAllowedToPushRewards(address _address, bool _allowed) external;

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
    ) external;

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
    ) external;
}