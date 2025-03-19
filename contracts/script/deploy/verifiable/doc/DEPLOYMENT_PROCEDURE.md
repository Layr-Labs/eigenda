# Deployment Procedure

This document details the deployment procedure using the `DeployVerifiable.sol` script. Verifying the correctness of the contracts is out of the scope of this document. The process of deploying EigenDA contracts with this script is divided into the following distinct phases:

## Deployment

* Use [placeholder.config.json](../../config/placeholder.config.json) as a template to create a configuration file with all parameters that the contracts need to be initialized with. 
* Any deployer can run the deployment script with `forge script DeployVerifiable` with the DEPLOY_CONFIG_PATH env variable pointing to the config file, and other deployer related parameters that can be found in [foundry's documentation](https://book.getfoundry.sh/reference/forge/forge-script). The contracts should be verified on etherscan
    * A full example command using a private key on Holesky would be `forge script DeployVerifiable --private-key $PRIVATE_KEY --rpc-url $HOLESKY_RPC_URL --etherscan-api-key $ETHERSCAN_API_KEY --verify --verifier etherscan --chain holesky --broadcast`
* After the script is run, a `DeploymentInitializer` contract will be deployed. The contract's address, github commit, and configuration file (if not in the commit) is to be sent over to the intended contract owner for verification.

## Verification

The intended owner of the contract should verify the following:

* The DeploymentInitializer contract received from the deployer is verified on etherscan. This contract contains getters for all contract addresses involved in the deployment, as well as all initialization variables that are statically sized. Any initialization parameters that are dynamically sized are to be submitted as calldata by the owner.
* For all contracts, verify:
    * the contract is verified on etherscan
    * the constructor arguments match what is expected
    * initial state is correct, elaborated on below for different contract types
* For proxies, verify:
    * the owner is the proxy admin listed in DeploymentInitializer
    * there were no successful transactions on the proxy except construction
    * the implementation contract is the designated mock or empty contract.
* Verify all non-contract initialization parameters against the expected configuration.
* Generate calldata with the dynamically sized initialization parameters to call the initializeDeployment(params) DeploymentInitializer, and simulate the behavior of this transaction on the tool of your choice. (e.g. Tenderly, foundry tests)
    * All proxies should be upgraded to their respective implementations
    * All proxies should be initialized properly. What needs to be checked is specific to each contract, so refer to the specific contracts.
    * The proxy admin should be changed to the intended owner

After these checks, the intended owner can proceed to initialize the deployment using the calldata detailed above.