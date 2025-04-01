## EigenDA Cert Verfier Deployer

This script can be used to deploy an EigenDACertVerifier contract with custom security thresholds and quorum numbers. The deployment should only be performed on Ethereum L1 testnet or mainnet environment and is not supported on L2s.

### Config

To set up the deployment, a config json should be placed in the `config/` folder with the following structure:

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

Two sample configs are provided in the `config/` folder for preprod and testnet environments.

### Deployment

To deploy the contract, run the following command passing in the path to the config file, the output path, and appropriate keys

```bash
forge script script/deploy/certverifier/CertVerifierDeployer.s.sol:CertVerifierDeployer --sig "run(string, string)" <config.json> <output.json> --rpc-url $RPC --private-key $PRIVATE_KEY -vvvv --etherscan-api-key $ETHERSCAN_API_KEY --verify --broadcast
```

The deployment will output the address of the deployed contract to a json file in the `output/` folder named `certverifier_deployment_data.json`

```json
{
    "eigenDACertVerifier": "0x..."
}
```

#### EigenDA V1

To deploy with just EigenDA V1 verification enabled, you only need to pass valid addresses for the following dependency contracts into your json config:
- `eigenDAThresholdRegistry`
- `eigenDAServiceManager`

V2 dependency addresses can be set to `0x..0`.

**Currently** V1 only verification is a short-term feature expression of the contract which will be deprecated coinciding the eventual deprecation fo the EigenDA V1 network. There are no explicit checks to prevent someone from unintendedly calling V2 verification functions which will revert nor is it explicitly obvious which protocols the cert verifier supports (i.e, V1, V2, V1 && V2). Please use at your own caution! 