# EigenDACertVerifierRouter Deployment

This directory contains the deployment script for the EigenDACertVerifierRouter contract.

## Overview

The EigenDACertVerifierRouter is a routing contract that directs certificate verification requests to the appropriate cert verifier contract based on the reference block number (RBN) in the certificate. This contract is deployed as implementation behind an OpenZeppelin [ERC1967](https://eips.ethereum.org/EIPS/eip-1967) proxy.


## Deployment

To deploy the EigenDACertVerifierRouter, use the following command:

```shell
forge script script/deploy/router/CertVerifierRouterDeployer.s.sol:CertVerifierRouterDeployer \
  --sig "run(string, string)" <config.json> <output.json> \
  --rpc-url $RPC \
  --private-key $PRIVATE_KEY \
  -vvvv \
  --etherscan-api-key $ETHERSCAN_API_KEY \
  --verify \
  --broadcast
```

### Configuration

Create a configuration file in the `config/` directory with the following format:

```json
{
  "initialOwner": "0x0000000000000000000000000000000000000001",
  "initABNConfigs" : [
    {
      "blockNumber": 0,
      "certVerifier": "0x0000000000000000000000000000000000000002"
    }
  ],
  "proxyAdmin": "0x0000000000000000000000000000000000000003"
}
```

- The `initialOwner` parameter specifies the address that will be set as the owner of the deployed router contract.
- The `initABNConfigs` specifies the activation block numbers that each initial cert verifier will be placed at with respect to block history, and the address of each.
- The `proxyAdmin` parameter specifies the address of the proxy admin for the transparent proxy.

### Post-Deployment

After deployment, the router is initialized with the provided initial cert verifier at block height 0. The owner will need to call `addCertVerifier(uint32 abn, address certVerifier)` to register additional cert verifiers with their activation block numbers (ABNs).

The deployment script will write the deployment addresses to an output JSON file in the format:

```json
{
  "eigenDACertVerifierRouter": "0x...",
  "eigenDACertVerifierRouterImplementation": "0x..."
}
```