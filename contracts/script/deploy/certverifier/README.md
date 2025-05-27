## EigenDA V2 Cert Verfier Deployer

This script can be used to deploy an immutable EigenDACertVerifier contract for EigenDA V2 with custom security thresholds and quorum numbers. The deployment should only be performed on Ethereum L1 testnet or mainnet environment and is not currently supported on L2s.

### Config

To set up the deployment, a config json should be placed in the `config/` folder with the following structure:

```json
{
    "eigenDAServiceManager": "0x...",
    "eigenDAThresholdRegistry": "0x...",

    "defaultSecurityThresholds": {
        "0_confirmationThreshold": 55,
        "1_adversaryThreshold": 33
    },

    "quorumNumbersRequired": "0x0001"
}
```

Three sample configs are provided in the `config/` folder for preprod and testnet environments on holesky as well as testnet environment on sepolia.

### Deployment

To deploy the contract, run the following command passing in the path to the config file, the output path, and appropriate keys

```bash
forge script script/deploy/certverifier/CertVerifierDeployerV2.s.sol:CertVerifierDeployerV2 --sig "run(string, string)" <config.json> <output.json> --rpc-url $RPC --private-key $PRIVATE_KEY -vvvv --etherscan-api-key $ETHERSCAN_API_KEY --verify --broadcast
```

The deployment will output the address of the deployed contract to a json file in the `output/` folder named `certverifier_deployment_data.json`

```json
{
    "eigenDACertVerifier": "0x..."
}
```
