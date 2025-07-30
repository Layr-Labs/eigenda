#!/bin/bash
set -euo pipefail

# Incrementally fetch git history until we find the merge base with the target branch
# Usage: ./fetch-base-incrementally.sh <base-ref>

if [ $# -eq 0 ]; then
    echo "Error: No base branch specified"
    echo "Usage: $0 <base-ref>"
    exit 1
fi

BASE_REF="$1"
echo "Fetching git history to find merge base with $BASE_REF branch..."

# Check if we have the merge base and print it
check_merge_base() {
    if MERGE_BASE=$(git merge-base HEAD "origin/$BASE_REF" 2>/dev/null); then
        echo "✓ Found merge base: $(git rev-parse --short "$MERGE_BASE")"
        return 0
    fi
    return 1
}

FETCH_INCREMENT=1 # How many commits to deepen on each iteration
MAX_ITERATIONS=100   # The number of times to deepen before we just fetch everything

# Start with a minimal fetch of just the base ref
git fetch --depth=1 origin "$BASE_REF" >/dev/null 2>&1

# Now deepen to get history from both branches
echo "→ Fetching up to $FETCH_INCREMENT commits of shared history..."
git fetch --deepen=$FETCH_INCREMENT origin HEAD "$BASE_REF"

# Check if we already have the merge base
if check_merge_base; then
    exit 0
fi

# Incrementally deepen by FETCH_INCREMENT commits at a time
for i in $(seq 1 $MAX_ITERATIONS); do
    echo "→ Need more history, fetching up to $FETCH_INCREMENT additional commits..."
    git fetch --deepen=$FETCH_INCREMENT origin HEAD "$BASE_REF"
    
    if check_merge_base; then
        exit 0
    fi
done

# Final fallback to full fetch if needed
echo "⚠ Branch history is deep, fetching all commits..."
git fetch --unshallow origin "$BASE_REF" || git fetch origin "$BASE_REF"

# Verify we found it
if ! check_merge_base; then
    echo "✗ Failed to find merge base with $BASE_REF even after full fetch"
    exit 1
fi
