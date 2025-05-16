# EigenDACertVerifierRouter Immutable Deployment

This directory contains the deployment script for the EigenDACertVerifierRouter contract.

## Overview

The EigenDACertVerifierRouter is a routing contract that directs certification verification requests to the appropriate cert verifier contract based on the reference block number (RBN) in the certificate.

## Deployment

To deploy an immutable EigenDACertVerifierRouter, use the following command:

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
  "initialCertVerifier": "0x0000000000000000000000000000000000000000"
}
```

- The `initialOwner` parameter specifies the address that will be set as the owner of the deployed router contract.
- The `initialCertVerifier` parameter specifies the initial address of the cert verifier initialized at block height 0.

### Post-Deployment

After deployment, the router contract is deployed but doesn't have any cert verifiers registered. The owner will need to call `addCertVerifier(uint32 abn, address certVerifier)` to register cert verifiers with their activation block numbers (ABNs).

The deployment script will write the router's address to an output JSON file in the format:

```json
{
  "eigenDACertVerifierRouter": "0x..."
}
```