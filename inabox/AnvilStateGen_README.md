# Anvil State Generation steps for `N` Operators

## Generate Anvil State for 4 Operators for Anvil Chain to run on Kubernetes:
1. Update InitialSupply in the contract to 100000 ether enough for 200 operators
[Click here to view the highlighted code on GitHub](https://github.com/Layr-Labs/eigenda/blob/7a16b44b8b06e770e15d372108df2fd220720697/contracts/script/SetUpEigenDA.s.sol#L58C38-L58C38)


```solidity
// Define the initial supply as 100000 ether
uint256 initialSupply = 100000 ether;
```

2. Update InABox testconfig-anvil.yaml with below for 20 Operators for Anvil Chain to run on Kubernetes:
```yaml
environment:
  name: "staging"
  type: "local"

deployers:
- name: "default"
  rpc: http://localhost:8545
  verifyContracts: false
  verifierUrl: http://localhost:4000/api
  deploySubgraphs: true
  slow: false

eigenda:
  deployer: "default"

privateKeys:
  file: /inabox/secrets
  ecdsaMap:
    default:
      privateKey: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
    batcher0:
      privateKey: 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d

services:
  counts:
    dispersers: 1
    operators: 20
  stakes:
    total: 100000e18
    distribution: [1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5]
  basePort: 32000
  variables:
    globals:
      HOSTNAME:
      TIMEOUT: 20s
      CHAIN_RPC: http://0.0.0.0:8545
      CHAIN_ID: 31337
      G1_PATH: /data/kzg/g1.point
      G2_PATH: /data/kzg/g2.point
      CACHE_PATH: /data/kzg/SRSTables
      SRS_ORDER: 300000
      CHALLENGE_ORDER: 300000
      STD_LOG_LEVEL: "trace"
      FILE_LOG_LEVEL: "trace"
      VERBOSE: true
      NUM_CONNECTIONS: 50
      AWS_ENDPOINT_URL:
      AWS_REGION: us-east-1
      AWS_ACCESS_KEY_ID:
      AWS_SECRET_ACCESS_KEY:
      ENCODER_ADDRESS: encoder.encoder.svc.cluster.local:34000
      USE_GRAPH: false
```

3. Run Anvil with below command in another terminal:
```
anvil --port 8545 --dump-state opr-state.json
```

Output:
```
forge script script/SetUpEigenDA.s.sol:SetupEigenDA --rpc-url http://127.0.0.1:8545 \
    --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
    --broadcast

forge script script/MockRollupDeployer.s.sol:MockRollupDeployer --rpc-url http://127.0.0.1:8545 \
    --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80 \
    --broadcast --sig run(address,bytes32) \
    0xc5a5C42992dECbae36851359345FE25997F5C42d fb89ee77edb64bdddc6f0e840cf2265e481be2810a8868e2853243ff89bdc24e

Generating variables
Test environment has successfully deployed!
```

Copy generated states to states directory in this repo [here](https://github.com/Layr-Labs/eigenda-devops/tree/master/charts/anvil-chain/states)
```
1. Copy the generated state: opr-state.json and build docker image. Instructions here https://github.com/Layr-Labs/eigenda-devops/blob/master/charts/anvil-chain/README.md
```

## Generate Anvil State for 200 Operators for Anvil Chain to run on Kubernetes:
1. Use Secrets from dir: `inabox/secrets/keys_for_200_operators.zip` 
2. Update testconfig-anvil.yaml to below

```yaml
environment:
  name: "staging"
  type: "local"

deployers:
- name: "default"
  rpc: http://localhost:8545
  verifyContracts: false
  verifierUrl: http://localhost:4000/api
  deploySubgraphs: true
  slow: false

eigenda:
  deployer: "default"

privateKeys:
  file: /inabox/secrets
  ecdsaMap:
    default:
      privateKey: 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
    batcher0:
      privateKey: 0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d

services:
  counts:
    dispersers: 1
    operators: 200
  stakes:
    total: 100000e18
    distribution: [1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5, 1.3, 2, 3, 5]
  basePort: 32000
  variables:
    globals:
      HOSTNAME:
      TIMEOUT: 20s
      CHAIN_RPC: http://0.0.0.0:8545
      CHAIN_ID: 31337
      G1_PATH: /data/kzg/g1.point
      G2_PATH: /data/kzg/g2.point
      CACHE_PATH: /data/kzg/SRSTables
      SRS_ORDER: 300000
      CHALLENGE_ORDER: 300000
      STD_LOG_LEVEL: "trace"
      FILE_LOG_LEVEL: "trace"
      VERBOSE: true
      NUM_CONNECTIONS: 50
      AWS_ENDPOINT_URL:
      AWS_REGION: us-east-1
      AWS_ACCESS_KEY_ID:
      AWS_SECRET_ACCESS_KEY:
      ENCODER_ADDRESS: encoder.encoder.svc.cluster.local:34000
      USE_GRAPH: false
```

