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
echo "Fetching history to find merge base with $BASE_REF"

# How many commits to deepen on each iteration
FETCH_INCREMENT=50

# The number of times to deepen before we just fetch everything
MAX_ITERATIONS=4

# Start with a shallow fetch
git fetch --depth=$FETCH_INCREMENT origin "$BASE_REF"

# Check if we already have the merge base
if git merge-base HEAD "origin/$BASE_REF" >/dev/null 2>&1; then
    echo "✓ Found merge base in initial fetch"
    exit 0
fi

# Incrementally deepen by FETCH_INCREMENT commits at a time
for i in $(seq 1 $MAX_ITERATIONS); do
    echo "→ Merge base not found after $((i * FETCH_INCREMENT)) commits, deepening..."
    git fetch --deepen=$FETCH_INCREMENT origin "$BASE_REF"
    
    if git merge-base HEAD "origin/$BASE_REF" >/dev/null 2>&1; then
        echo "✓ Found merge base after fetching ~$(((i + 1) * FETCH_INCREMENT)) commits"
        exit 0
    fi
done

# Final fallback to full fetch if needed
echo "⚠ Merge base still not found after $((MAX_ITERATIONS * FETCH_INCREMENT + FETCH_INCREMENT)) commits, fetching full history..."
git fetch --unshallow origin "$BASE_REF" || git fetch origin "$BASE_REF"

# Verify we found it
if git merge-base HEAD "origin/$BASE_REF" >/dev/null 2>&1; then
    echo "✓ Successfully found merge base with full history"
else
    echo "✗ Failed to find merge base with $BASE_REF even after full fetch"
    exit 1
fi
