#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the script directory
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
SUBGRAPH_DIR="$(dirname "$SCRIPT_DIR")"
REPO_ROOT="$(dirname "$(dirname "$SUBGRAPH_DIR")")"
CONTRACTS_DIR="$REPO_ROOT/contracts"
ABI_DIR="$SUBGRAPH_DIR/abis"

echo -e "${YELLOW}EigenDA Operator State ABI Update Script${NC}"
echo "========================================"

# Check if yq is installed
if ! command -v yq &> /dev/null; then
    echo -e "${RED}Error: yq is not installed or not in PATH${NC}"
    echo "Please install yq: brew install yq (macOS) or see https://github.com/mikefarah/yq"
    exit 1
fi

# Navigate to contracts directory
cd "$CONTRACTS_DIR"

# Compile contracts
echo -e "${YELLOW}Compiling contracts...${NC}"
if [ -f "Makefile" ] && grep -q "compile:" "Makefile"; then
    make compile
elif [ -f "compile.sh" ]; then
    ./compile.sh
else
    echo -e "${RED}Error: No compilation script found${NC}"
    exit 1
fi

# Define the contracts to extract
# Using simple variables instead of associative array for shell compatibility
CONTRACTS="RegistryCoordinator BLSApkRegistry EjectionManager"

# Create abis directory if it doesn't exist
mkdir -p "$ABI_DIR"

# Track if any ABI changed
ABI_CHANGED=false

# Extract ABIs
for CONTRACT_NAME in $CONTRACTS; do
    ARTIFACT_PATH="$CONTRACTS_DIR/out/${CONTRACT_NAME}.sol/${CONTRACT_NAME}.json"
    
    echo -e "${YELLOW}Extracting ${CONTRACT_NAME} ABI...${NC}"
    
    if [ ! -f "$ARTIFACT_PATH" ]; then
        echo -e "${RED}Error: ${CONTRACT_NAME} artifact not found at $ARTIFACT_PATH${NC}"
        echo "Make sure the contract compiled successfully"
        exit 1
    fi
    
    # Extract ABI using yq and save to file
    ABI_FILE="$ABI_DIR/${CONTRACT_NAME}.json"
    yq -o=json '.abi' "$ARTIFACT_PATH" > "$ABI_FILE"
    
    # Verify the ABI was extracted correctly
    if [ ! -s "$ABI_FILE" ]; then
        echo -e "${RED}Error: Failed to extract ABI for ${CONTRACT_NAME} or ABI is empty${NC}"
        exit 1
    fi
    
    # Check if ABI changed
    cd "$SUBGRAPH_DIR"
    if ! git diff --quiet "$ABI_FILE" 2>/dev/null; then
        echo -e "${YELLOW}⚠ ${CONTRACT_NAME} ABI has changed!${NC}"
        ABI_CHANGED=true
    fi
done

# Report results
cd "$SUBGRAPH_DIR"
if [ "$ABI_CHANGED" = false ]; then
    echo -e "${GREEN}✓ All ABIs are up to date (no changes)${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠ One or more ABIs have changed!${NC}"
    echo "Changed files:"
    git diff --stat "$ABI_DIR"/*.json 2>/dev/null || echo "New ABI files created"
    
    # For CI, exit with non-zero code if ABI changed
    if [ "$CI" = "true" ]; then
        echo -e "${RED}CI Check Failed: ABIs have changed and need to be committed${NC}"
        exit 1
    fi
    
    exit 0
fi