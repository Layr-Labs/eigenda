## EigenDA Cert Verfier Deployer

This script can be used to deploy an EigenDACertVerifier contract with specified security thresholds and quorum numbers required

### Config

To set up the deployment, a config json should be placed in the 'config/' folder with the following structure:

```json
{
    "eigenDAServiceManager": "0x...",
    "eigenDAThresholdRegistry": "0x...",
    "eigenDARelayRegistry": "0x...",
    "registryCoordinator": "0x...",
    "operatorStateRetriever": "0x...",

    "defaultSecurityThresholds": {
        "0_confirmationThreshold": 55,
        "1_adversaryThreshold": 33
    },

    "quorumNumbersRequired": "0x0001"
}
```

Two sample configs are provided in the 'config/' folder for preprod and testnet environments.

### Deployment

To deploy the contract, run the following command passing in the path to the config file and appropriate keys

```bash
forge script script/deploy/certverifier/CertVerifierDeployer.s.sol:CertVerifierDeployer --sig "run(string)" <config.json> --rpc-url $RPC --private-key $PRIVATE_KEY -vvvv --etherscan-api-key $ETHERSCAN_API_KEY --verify --broadcast
```

The deployment will output the address of the deployed contract to a json file in the 'output/' folder named 'certverifier_deployment_data.json'

```json
{
    "eigenDACertVerifier": "0x..."
}
```
