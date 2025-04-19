#!/bin/bash

# Check if API key is provided
if [ -z "$1" ]; then
    echo "Usage: $0 <etherscan_api_key>"
    exit 1
fi

# Store the API key
api_key="$1"

# Define all contract addresses with their names
declare -A CONTRACTS=(
    ["IndexRegistry"]="0x1ae0b73118906f39d5ed30ae4a484ce2f479a14c"
    ["StakeRegistry"]="0x1c468cf7089d263c2f53e2579b329b16abc4dd96"
    ["EjectionManager"]="0x33a517608999df5ceffa2b2eba88b4461c26af6f"
    ["EigenDAServiceManager"]="0x58fde694db83e589abb21a6fe66cb20ce5554a07"
    ["SocketRegistry"]="0x5b60105ced5207d6ad217bf2d426e133454ecfb4"
    ["BLSApkRegistry"]="0x5d0b9ce2e277daf508528e9f6bf6314e79e4ed2b"
    ["RegistryCoordinator"]="0xdcabf0be991d4609096cce316df08d091356e03f"
)

# Fetch all contracts
echo "Starting to fetch all contracts..."

for name in "${!CONTRACTS[@]}"; do
    address="${CONTRACTS[$name]}"
    echo "Fetching $name ($address)..."
    ./fetch_source.sh "$address" "$api_key"
    echo "Completed $name"
    echo "------------------------"
done

echo "All contracts fetched successfully!"