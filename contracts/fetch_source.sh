#!/bin/bash

# Check if an address was provided
if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Usage: $0 <contract_address> <etherscan_api_key>"
    exit 1
fi

# Store the address and API key
address=$(echo "$1" | tr '[:upper:]' '[:lower:]')
api_key="$2"

# Create the directory if it doesn't exist
mkdir -p "./sources/$address"

# Export the Etherscan API key and fetch the source code
export ETHERSCAN_API_KEY="$api_key"
cast source "$address" -d "./sources/$address"

# Get contract name from Etherscan
echo "Getting contract name from Etherscan..."
response=$(curl -s "https://api.etherscan.io/api?module=contract&action=getsourcecode&address=$address&apikey=$api_key")
contract_name=$(echo "$response" | jq -r '.result[0].ContractName')
echo "Contract name: $contract_name"

# Note: foundry.toml must be copied manually as needed

# Build the contracts
echo "Building contracts..."
(cd "./sources/$address" && forge build)

echo "Source code fetched and compiled in ./sources/$address/"