## Subraph deployment on Alchemy

A subgraph deployment on Alchemy requires the following steps:
1. install required deps:
```bash
cd subgraph/<subgraph_name>
yarn [install]
```

2. generate subgraph code
```bash
yarn codegen
```

Next, export the deploy key used for deploying subgraphs - go to the [Alchemy Dashboard](https://subgraphs.alchemy.com/dashboard) and click on any subgraph; it will be shown on `Deploy New Version` textbox:
```bash
export DEPLOY_KEY=...
```

Keep in mind that unlike our `graph-node` Satsuma [does enforce subgraph version control](https://docs.alchemy.com/reference/subgraph-versioning), which means that any update in a graph code or configuration requires a new version deployment. Redeploying the same version returns an error:
```bash
✖ Failed to deploy to Graph node https://subgraphs.alchemy.com/api/subgraphs/deploy: Subgraph version already exists
UNCAUGHT EXCEPTION: Error: EEXIT: 1
```

After a new version is published, it should be promoted to Live for DB access. Non-promoted subgraphs are not available for SQL querying since the deployment ID is only added to the DB for Live versions.

### Deployment ID

Redeploying the same subgraph without source code changes* under a new version does not change the deployment ID. However, changing `network.json` values does change the deployment ID.

\* In this scenario a new version deployment basically takes an instant since it reuses existing indexed data.

### Deployment

## mainnet
```bash
graph deploy eigenda-ejection-manager-holesky \
  --version-label v0.0.1 \
  --network mainnet \
  --network-file networks/mainnet.json
  --node https://subgraphs.alchemy.com/api/subgraphs/deploy \
  --deploy-key $DEPLOY_KEY \
  --ipfs https://ipfs.satsuma.xyz
```

## testnet

```bash
graph deploy eigenda-ejection-manager-holesky \
  --version-label v0.0.1 \
  --network holesky \
  --network-file networks/testnet.json
  --node https://subgraphs.alchemy.com/api/subgraphs/deploy \
  --deploy-key $DEPLOY_KEY \
  --ipfs https://ipfs.satsuma.xyz
```

## preprod

```bash
graph deploy eigenda-ejection-manager-preprod-holesky \
  --version-label v0.0.1 \
  --network holesky \
  --network-file networks/preprod.json
  --node https://subgraphs.alchemy.com/api/subgraphs/deploy \
  --deploy-key $DEPLOY_KEY \
  --ipfs https://ipfs.satsuma.xyz
```
