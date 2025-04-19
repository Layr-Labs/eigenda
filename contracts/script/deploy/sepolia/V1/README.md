# EigenDA V1 Deployment on Sepolia

## Overview

This directory contains scripts for deploying EigenDA V1 contracts to the Sepolia testnet. The purpose of these scripts is to test the upgrade path from V1 to V2 contracts on a testnet before executing the same upgrade on Ethereum mainnet.

Rather than deploying arbitrary versions of the contracts, these scripts deploy the exact same contract code that is currently running on mainnet. This ensures that the upgrade testing accurately represents what will happen on mainnet.

## How It Works

The deployment process:

1. **Fetch contract source code** from Etherscan for existing mainnet implementations
2. **Compile each contract** individually in its respective compilation environment, creating deployment artifacts
3. **Deploy the artifacts** to Sepolia with Sepolia-specific parameters
4. **Verify** that the deployed contracts match the mainnet implementations

The fetching scripts not only download the source code but also compile each contract within its own isolated environment, making the compiled artifacts available to the forge deployment scripts.

This approach ensures we're testing the true upgrade path from the exact same contract code that's deployed on mainnet, minimizing the risk of unexpected issues during the actual mainnet upgrade.

## End-to-End Instructions

### Prerequisites

- An Etherscan API key (required for fetching contract source code)

### Setup Environment Variables

To run the complete deployment process:

```bash
# 1. Fetch mainnet contract sources and compile them in isolated environments
cd /workspaces/eigenda/contracts
./fetch_all_mainnet_v1_da_sources.sh YOUR_ETHERSCAN_API_KEY

# 2. Compile the deployment scripts
forge clean
forge build

# 3. Deploy V1 contracts to Sepolia
forge script script/deploy/sepolia/V1/DeployV1Contracts.s.sol:DeployV1Contracts --rpc-url $SEPOLIA_RPC_URL --private-key $PRIVATE_KEY --broadcast

# 4. Verify the deployed contracts
forge script script/deploy/sepolia/V1/VerifyV1Contracts.s.sol:VerifyV1Contracts --rpc-url $SEPOLIA_RPC_URL
```

## Script Descriptions

### 1. `fetch_all_mainnet_v1_da_sources.sh`

Fetches all the EigenDA V1 contract source code from Etherscan for the current mainnet deployments and compiles each contract in its own isolated environment.

**Usage:**
```bash
./fetch_all_mainnet_v1_da_sources.sh YOUR_ETHERSCAN_API_KEY
```

This script calls `fetch_source.sh` for each contract address, passing the Etherscan API key. It organizes the fetched source code in the `contracts/sources` directory, and compiles each contract independently. The compilation creates artifact files that are later used by the forge deployment scripts.

### 2. `fetch_source.sh`

Helper script that fetches a specific contract's source code from Etherscan and compiles it in an isolated environment.

**Usage:**
```bash
./fetch_source.sh <contract_address> <etherscan_api_key>
```

**Example:**
```bash
./fetch_source.sh 0x1ae0b73118906f39d5ed30ae4a484ce2f479a14c ABC123DEF456GHI789
```

This script:
1. Downloads the contract source code and its dependencies
2. Sets up a dedicated compilation environment matching the original
3. Compiles the contract to produce the same bytecode as mainnet
4. Makes the compiled artifact available for the deployment scripts

### 3. `ConfigV1Lib.sol`

Library that provides functions to read configuration values from a TOML file. Used by the deployment scripts to get Sepolia-specific configuration parameters.

### 4. `DeployV1Contracts.s.sol`

Main deployment script that handles the deployment of the V1 contracts to Sepolia.

**Usage:**
```bash
forge script script/deploy/sepolia/V1/DeployV1Contracts.s.sol:DeployV1Contracts --rpc-url $SEPOLIA_RPC_URL --private-key $PRIVATE_KEY --broadcast
```

This script:
1. Reads the configuration from a TOML file
2. Creates proxy admin and infrastructure contracts
3. Deploys the registry contracts (StakeRegistry, IndexRegistry, etc.)
4. Sets up all contracts with the right initialization parameters
5. Links all the contracts together

### 5. `VerifyV1Contracts.s.sol`

Verification script that checks the deployed contracts to ensure they're set up correctly.

**Usage:**
```bash
forge script script/deploy/sepolia/V1/VerifyV1Contracts.s.sol:VerifyV1Contracts --rpc-url $SEPOLIA_RPC_URL
```

This script checks that:
1. All contracts are initialized with correct parameters
2. Contract relationships are set up properly (e.g., registries linked to the coordinator)

## Configuration

The deployment is configured using a TOML file located at `script/deploy/sepolia/V1/config/sepolia.config.toml`. This file contains:

- Contract source locations
- InitParams for each contract
- EigenDA-specific configuration

You can modify this file to adjust the deployment parameters for your Sepolia testnet deployment.

## Notes

- **Important**: Contracts deployed this way cannot be verified using Forge's built-in verification mechanism (e.g., `--verify` flag). This is because the deployment uses pre-compiled artifacts, and the source code is not directly available to the deployment script. If verification on Etherscan is needed, it would have to be done manually using the source code in the `contracts/sources` directory.
- The fetch scripts require an Etherscan API key to download the contract source code. You can obtain one by creating an account on [Etherscan](https://etherscan.io/myapikey).