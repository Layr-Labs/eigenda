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

echo -e "${YELLOW}EigenDA Batch Metadata ABI Update Script${NC}"
echo "========================================="

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

# Extract EigenDAServiceManager ABI
echo -e "${YELLOW}Extracting EigenDAServiceManager ABI...${NC}"
SERVICE_MANAGER_ARTIFACT="$CONTRACTS_DIR/out/EigenDAServiceManager.sol/EigenDAServiceManager.json"

if [ ! -f "$SERVICE_MANAGER_ARTIFACT" ]; then
    echo -e "${RED}Error: EigenDAServiceManager artifact not found at $SERVICE_MANAGER_ARTIFACT${NC}"
    echo "Make sure the contract compiled successfully"
    exit 1
fi

# Create abis directory if it doesn't exist
mkdir -p "$ABI_DIR"

# Extract ABI using yq and save to file
echo -e "${YELLOW}Saving ABI to $ABI_DIR/EigenDAServiceManager.json${NC}"
yq -o=json '.abi' "$SERVICE_MANAGER_ARTIFACT" > "$ABI_DIR/EigenDAServiceManager.json"

# Verify the ABI was extracted correctly
if [ ! -s "$ABI_DIR/EigenDAServiceManager.json" ]; then
    echo -e "${RED}Error: Failed to extract ABI or ABI is empty${NC}"
    exit 1
fi

# Check if ABI changed
cd "$SUBGRAPH_DIR"
if git diff --quiet "$ABI_DIR/EigenDAServiceManager.json" 2>/dev/null; then
    echo -e "${GREEN}✓ ABI is up to date (no changes)${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠ ABI has changed!${NC}"
    echo "Changes detected in EigenDAServiceManager ABI:"
    git diff --stat "$ABI_DIR/EigenDAServiceManager.json" 2>/dev/null || echo "New ABI file created"
    
    # For CI, exit with non-zero code if ABI changed
    if [ "$CI" = "true" ]; then
        echo -e "${RED}CI Check Failed: ABI has changed and needs to be committed${NC}"
        exit 1
    fi
    
    exit 0
fi