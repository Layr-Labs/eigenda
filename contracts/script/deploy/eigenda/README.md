# EigenDA Deployment Script

This is the deployment script that is being used to deploy any fresh deployments of EigenDA on a new network. It is meant to replace the older deployment scripts after any dependencies on them are removed.

## Running the Script

A mainnet beta configuration is included in this folder. You can run the script with any configuration by setting the environment variable DEPLOY_CONFIG_PATH.

To run the script, you can run the following command with the DEPLOY_CONFIG_PATH environment variable set:

`forge script DeployEigenDA --rpc-url XXX --broadcast`

Please refer to [foundry's documentation](https://getfoundry.sh/forge/reference/forge-script) to set up your wallet, API keys, verification as necessary based on your use case.