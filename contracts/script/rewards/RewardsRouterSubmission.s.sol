// SPDX-License-Identifier: MIT
pragma solidity =0.8.12;

import {ServiceManagerRewardsRouter} from "../../src/rewards/ServiceManagerRewardsRouter.sol";
import {ServiceManagerBase, IRewardsCoordinator, IServiceManager} from "../../lib/eigenlayer-middleware/src/ServiceManagerBase.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {IServiceManagerRewardsRouter} from "../../src/rewards/IServiceManagerRewardsRouter.sol";
import {EigenDAServiceManager} from "../../src/core/EigenDAServiceManager.sol";
import "../../lib/eigenlayer-middleware/src/StakeRegistry.sol";
import "forge-std/Test.sol";
import "forge-std/Script.sol";
import "forge-std/StdJson.sol";
import "src/interfaces/IEigenDAStructs.sol";

contract RewardsRouterSubmission is Script, Test {

    address rewardsRouter;
    bytes quorumsNumbers;
    address[] tokens;
    uint256[] amounts;
    uint32 startTimestamp;
    uint32 duration;

    // forge script script/rewards/RewardsRouterSubmission.s.sol:RewardsRouterSubmission --sig "runEOA(string)" <config.json> --rpc-url $RPC --private-key $KEY -vvvv --broadcast
    function runEOA(string memory json) external {
        _parseParams(json);
        bytes memory calldata_to_rewardsRouter = _getCalldataToRewardsRouter();
        vm.startBroadcast();
        (bool success, ) = rewardsRouter.call(calldata_to_rewardsRouter);
        require(success, "Failed to call rewardsRouter");
        vm.stopBroadcast();
    }

    // forge script script/rewards/RewardsRouterSubmission.s.sol:RewardsRouterSubmission --sig "runCalldata(string)" <config.json> --rpc-url $RPC -vvvv 
    function runCalldata(string memory json) external {
        _parseParams(json);
        _getCalldataToRewardsRouter();
    }

    function _parseParams(string memory json) internal {
        string memory path = string.concat("./script/rewards/", json);
        string memory data = vm.readFile(path);

        bytes memory raw = stdJson.parseRaw(data, ".rewardsRouter");
        rewardsRouter = abi.decode(raw, (address));

        raw = stdJson.parseRaw(data, ".quorumsNumbers");
        quorumsNumbers = abi.decode(raw, (bytes));

        raw = stdJson.parseRaw(data, ".tokens");
        tokens = abi.decode(raw, (address[]));

        raw = stdJson.parseRaw(data, ".amounts");
        amounts = abi.decode(raw, (uint256[]));

        raw = stdJson.parseRaw(data, ".startTimestamp");
        startTimestamp = abi.decode(raw, (uint32));

        raw = stdJson.parseRaw(data, ".duration");
        duration = abi.decode(raw, (uint32));
    }

    function _getCalldataToRewardsRouter() public returns (bytes memory _calldata_to_rewardsRouter) {
        require(quorumsNumbers.length == tokens.length, "quorumsNumbers and tokens must be the same length");
        require(quorumsNumbers.length == amounts.length, "quorumsNumbers and amounts must be the same length");

        address stakeRegistry = address(EigenDAServiceManager(address(ServiceManagerRewardsRouter(rewardsRouter).serviceManager())).stakeRegistry());
        IRewardsCoordinator.RewardsSubmission[] memory rewardsSubmissions = new IRewardsCoordinator.RewardsSubmission[](quorumsNumbers.length);

        for (uint256 i = 0; i < quorumsNumbers.length; i++) {
            uint8 quorumNumber = uint8(quorumsNumbers[i]);
            uint256 numStrats = StakeRegistry(stakeRegistry).strategyParamsLength(quorumNumber); 

            IRewardsCoordinator.StrategyAndMultiplier[] memory strategyAndMultipliers = new IRewardsCoordinator.StrategyAndMultiplier[](numStrats);
            for (uint256 j = 0; j < numStrats; j++) {
                (IStrategy strategy, uint96 multiplier) = StakeRegistry(stakeRegistry).strategyParams(quorumNumber, j);
                strategyAndMultipliers[j] = IRewardsCoordinator.StrategyAndMultiplier({
                    strategy: strategy,
                    multiplier: multiplier
                });
            }

            strategyAndMultipliers = _sort(strategyAndMultipliers);

            rewardsSubmissions[i] = IRewardsCoordinator.RewardsSubmission({
                strategiesAndMultipliers: strategyAndMultipliers,
                token: IERC20(tokens[i]),
                amount: amounts[i],
                startTimestamp: startTimestamp,
                duration: duration
            });
        }

        _calldata_to_rewardsRouter = abi.encodeWithSelector(
            ServiceManagerRewardsRouter.createAVSRewardsSubmission.selector,
            rewardsSubmissions
        );

        emit log_named_bytes("calldata_to_rewardsRouter", _calldata_to_rewardsRouter);
        return _calldata_to_rewardsRouter;
    }

    function _sort(IRewardsCoordinator.StrategyAndMultiplier[] memory strategyAndMultipliers) public pure returns (IRewardsCoordinator.StrategyAndMultiplier[] memory) {
        uint length = strategyAndMultipliers.length;
        for (uint i = 1; i < length; i++) {
            uint key = uint(uint160(address(strategyAndMultipliers[i].strategy)));
            uint96 multiplier = strategyAndMultipliers[i].multiplier;
            int j = int(i) - 1;
            while ((int(j) >= 0) && (uint(uint160(address(strategyAndMultipliers[uint(j)].strategy))) > key)) {
                strategyAndMultipliers[uint(j) + 1].strategy = strategyAndMultipliers[uint(j)].strategy;
                strategyAndMultipliers[uint(j) + 1].multiplier = strategyAndMultipliers[uint(j)].multiplier;
                j--;
            }
            strategyAndMultipliers[uint(j + 1)].strategy = IStrategy(address(uint160(key)));
            strategyAndMultipliers[uint(j + 1)].multiplier = multiplier;
        }
        return strategyAndMultipliers;
    }
}