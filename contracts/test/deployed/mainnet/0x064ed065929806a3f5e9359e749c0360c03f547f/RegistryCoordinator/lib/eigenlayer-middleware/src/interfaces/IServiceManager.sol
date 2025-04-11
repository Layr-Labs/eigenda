// SPDX-License-Identifier: BUSL-1.1
pragma solidity >=0.5.0;

import {IRewardsCoordinator} from "eigenlayer-contracts/src/contracts/interfaces/IRewardsCoordinator.sol";
import {IServiceManagerUI} from "./IServiceManagerUI.sol";

/**
 * @title Minimal interface for a ServiceManager-type contract that forms the single point for an AVS to push updates to EigenLayer
 * @author Layr Labs, Inc.
 */
interface IServiceManager is IServiceManagerUI {
    /**
     * @notice Creates a new rewards submission to the EigenLayer RewardsCoordinator contract, to be split amongst the
     * set of stakers delegated to operators who are registered to this `avs`
     * @param rewardsSubmissions The rewards submissions being created
     * @dev Only callable by the permissioned rewardsInitiator address
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
     * @notice Creates a new operator-directed rewards submission on behalf of an AVS, to be split amongst the operators and
     * set of stakers delegated to operators who are registered to the `avs`.
     * @param operatorDirectedRewardsSubmissions The operator-directed rewards submissions being created
     * @dev Only callable by the permissioned rewardsInitiator address
     * @dev The duration of the `rewardsSubmission` cannot exceed `MAX_REWARDS_DURATION`
     * @dev The tokens are sent to the `RewardsCoordinator` contract
     * @dev This contract needs a token approval of sum of all `operatorRewards` in the `operatorDirectedRewardsSubmissions`, before calling this function.
     * @dev Strategies must be in ascending order of addresses to check for duplicates
     * @dev Operators must be in ascending order of addresses to check for duplicates.
     * @dev This function will revert if the `operatorDirectedRewardsSubmissions` is malformed.
     */
    function createOperatorDirectedAVSRewardsSubmission(
        IRewardsCoordinator.OperatorDirectedRewardsSubmission[]
            calldata operatorDirectedRewardsSubmissions
    ) external;

    /**
     * @notice Forwards a call to Eigenlayer's RewardsCoordinator contract to set the address of the entity that can call `processClaim` on behalf of this contract.
     * @param claimer The address of the entity that can call `processClaim` on behalf of the earner
     * @dev Only callable by the owner.
     */
    function setClaimerFor(address claimer) external;

    // EVENTS
    event RewardsInitiatorUpdated(
        address prevRewardsInitiator,
        address newRewardsInitiator
    );
}
