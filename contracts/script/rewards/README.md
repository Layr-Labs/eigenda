## EigenDA Reward Submission Script

This script can be used to submit rewards to operators of EigenDA quorums through the ServiceManagerRewardsRouter contract

### Config

To set up the deployment, a config json should be placed in the `rewards/` folder with the following structure:

```json
{
    "rewardsRouter": "0x...",
    "quorumsNumbers": "0x0001",
    "tokens": [
        "0x...",
        "0x..."
    ],
    "amounts": [
        1000000000000000000,
        1000000000000000000
    ],
    "startTimestamp": 1,
    "duration": 2
}
```
The `rewardsRouter` should be the address of the ServiceManagerRewardsRouter contract.

The `quorumsNumbers` should be a hex string of the quorum numbers to submit rewards to.

The `tokens` should be a list of addresses of the tokens to submit as rewards where the index of the token in the list corresponds to the index of the quorum number in the `quorumsNumbers` string.

The `amounts` should be a list of the amounts to submit as rewards where the index of the amount in the list corresponds to the index of the quorum number in the `quorumsNumbers` string and the token in the `tokens` list.

The `startTimestamp` should be the timestamp of when the rewards submission begins.

The `duration` should be the duration of the rewards submission in seconds.

### Execution

To submit rewards, you can run the following commands to either execute the script from an EOA or log the calldata and submit it from a multisig respectively.

```bash
forge script script/rewards/RewardsRouterSubmission.s.sol:RewardsRouterSubmission --sig "runEOA(string)" <config.json> --rpc-url $RPC --private-key $PRIVATE_KEY -vvvv --broadcast
```

```bash
forge script script/rewards/RewardsRouterSubmission.s.sol:RewardsRouterSubmission --sig "runCalldata(string)" <config.json> --rpc-url $RPC -vvvv 
```

