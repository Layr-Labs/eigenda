#!/usr/bin/env bash

# Script to install git hooks for the EigenDA repository
# This script works correctly in both regular git repos and git worktrees
#
# Usage:
#   ./scripts/install-hooks.sh         # Install hooks (overwrites existing)
#   mise run install-hooks             # Recommended: Install via mise

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get the repository root
REPO_ROOT=$(git rev-parse --show-toplevel)

# Get the git common directory (handles both regular repos and worktrees)
# This ensures hooks are installed in the shared location for worktrees
GIT_COMMON_DIR=$(git rev-parse --git-common-dir)

# The hooks directory
HOOKS_DIR="$GIT_COMMON_DIR/hooks"

# Source hooks directory
SOURCE_HOOKS_DIR="$REPO_ROOT/scripts/hooks"

echo -e "${YELLOW}Installing git hooks...${NC}"
echo "Repository root: $REPO_ROOT"
echo "Git hooks directory: $HOOKS_DIR"

# Ensure hooks directory exists
if [ ! -d "$HOOKS_DIR" ]; then
    echo -e "${RED}Error: Hooks directory does not exist: $HOOKS_DIR${NC}"
    exit 1
fi

# Install pre-commit hook
PRE_COMMIT_SOURCE="$SOURCE_HOOKS_DIR/pre-commit"
PRE_COMMIT_TARGET="$HOOKS_DIR/pre-commit"

if [ ! -f "$PRE_COMMIT_SOURCE" ]; then
    echo -e "${RED}Error: Source pre-commit hook not found: $PRE_COMMIT_SOURCE${NC}"
    exit 1
fi

# Check if hook already exists and remove it
if [ -f "$PRE_COMMIT_TARGET" ] || [ -L "$PRE_COMMIT_TARGET" ]; then
    echo -e "${YELLOW}Pre-commit hook already exists, overwriting...${NC}"
    rm -f "$PRE_COMMIT_TARGET"
fi

# Copy the hook (we use cp instead of symlink for better portability)
cp "$PRE_COMMIT_SOURCE" "$PRE_COMMIT_TARGET"
chmod +x "$PRE_COMMIT_TARGET"

echo -e "${GREEN}✓ Pre-commit hook installed successfully${NC}"
echo ""
echo "The following checks will run before each commit:"
echo "  - Linting (golangci-lint)"
echo "  - Go mod tidy check"
echo "  - Format checking (Go and contracts)"
echo ""
echo -e "${YELLOW}Note:${NC} You can bypass these checks using: git commit --no-verify"
echo -e "${YELLOW}Note:${NC} Make sure you have run 'mise install' to set up all required tools"

exit 0
