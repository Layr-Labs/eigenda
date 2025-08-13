# Release Runbook

This runbook is meant to be terse and to the point. 
For more detailed information, best practices, and philosophical stances, see our [release management process](release-process.md) guide.

> NOTE: WE ONLY DEPLOY V2 ARTIFACTS. V1 DISPERSERS ARE FROZEN ON 0.9.2. THIS SIMPLIFIES THE DEPLOYMENT PROCESS.

1. Create release/X.Y.Z branch using [eigenda-releaser action](https://github.com/Layr-Labs/eigenda/actions/workflows/eigenda-releaser.yaml)
    1. Input is X.Y.Z
    2. Creates branch on latest commit on master
    3. Needs approval from someone else
2. Create and push vX.Y.Z-rc.1 tag
3. Prepare release notes
   1. Create a new "Release X.Y.Z" page on our [notion release management](https://www.notion.so/eigen-labs/Monorepo-Release-Mgmt-21f13c11c3e0802b9d7fcf4173a49d12) page
   2. Use Claude to generate release notes, following [Release 2.1.0](https://www.notion.so/Release-2-1-0-22413c11c3e0809ea585eafa44a3fc23?pvs=21) as an example
   3. Go through the release notes in DETAIL and make sure you understand all the implications of the changes.
   4. Amend the release notes to be human digestable
4. Announce code freeze at that commit on slack
5. Deploy to preprod (disperser + all validators)
   1. TODO: how?
6. Few hours preprod load test (see [load test](#load-test) section below)
    1. if any bugs are found:
    2. make any necessary fixes on master
    3. cherry-pick to release/X.Y.Z branch
    4. create vX.Y.Z-rc.N+1 tag
    5. go to step 5
7. Create vX.Y.Z tag
8. Promote holesky disperser + our own validator to v.X.Y.Z
9.  Promote sepolia disperser + validators
10. Promote mainnet-beta disperser (???)
11. Release v.X.Y.Z:
    1.  Publish release on eigenda repo using release notes created in step 3 (but mark as pre-release on github!)
    2.  Update https://github.com/Layr-Labs/eigenda-operator-setup, example: https://github.com/Layr-Labs/eigenda-operator-setup/pull/170/files + make a release on that repo
    3.  Make operator announcement on slack
12. Soak test for a couple days (run load generator with low traffic)
    1. if any bugs are found that need to be patched immediately, go to step 2 and start vX.Y.Z+1 release process
13. Promote mainnet disperser
14. Update github release:
    1. Update title: [Testnet] EigenDA v2.1.0 (Pre-release) ——> [Mainnet] EigenDA v2.1.0
    2. Uncheck the "pre-release" option
15. Make operator + rollup announcement on slack

## Load Test

```yaml
# Example 30MiB/s Load Gen Config
replicaCount: 8
"MBPerSecond": 3.75
```

### Starting Load Test:

1. git clone [eigenda-devops](https://github.com/Layr-Labs/eigenda-devops) repo
2. Modify [replicaCount](https://github.com/Layr-Labs/eigenda-devops/blob/75ca819594789e87d2db0f462f96f386bd5b0291/charts/traffic-generator-v2/values/eigenda-preprod/us-east-1/holesky/values.yaml#L15)
3. Modify [MBPerSecond](https://github.com/Layr-Labs/eigenda-devops/blob/75ca819594789e87d2db0f462f96f386bd5b0291/charts/traffic-generator-v2/values/eigenda-preprod/us-east-1/holesky/values.yaml#L80)
4. `cd charts/traffic-generator-v2`
5. `helm upgrade --install --atomic --namespace=traffic-generator-v2 --values values/eigenda-preprod/us-east-1/holesky/values.yaml traffic-generator-v2 .`