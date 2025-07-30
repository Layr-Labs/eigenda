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
echo "Debug: Current HEAD is $(git rev-parse --short HEAD)"
echo "Debug: Current branch has $(git rev-list --count HEAD 2>/dev/null || echo 'unknown') commits"
echo "Debug: Git log shows: $(git log --oneline -5 2>/dev/null || echo 'no history available')"

# How many commits to deepen on each iteration
FETCH_INCREMENT=50

# The number of times to deepen before we just fetch everything
MAX_ITERATIONS=4

# Start with a minimal fetch of just the base ref
echo "Debug: Fetching 1 commit from origin/$BASE_REF to establish the ref"
git fetch --depth=1 origin "$BASE_REF"
echo "Debug: After fetch, origin/$BASE_REF is at $(git rev-parse --short "origin/$BASE_REF" 2>/dev/null || echo 'not available')"

# Now deepen to get history from both branches
echo "Debug: Deepening by $FETCH_INCREMENT commits to fetch history from both branches"
git fetch --deepen=$FETCH_INCREMENT
echo "Debug: After deepening, HEAD has $(git rev-list --count HEAD 2>/dev/null || echo 'unknown') commits"
echo "Debug: origin/$BASE_REF has $(git rev-list --count "origin/$BASE_REF" 2>/dev/null || echo 'unknown') commits"

# Check if we already have the merge base
if MERGE_BASE=$(git merge-base HEAD "origin/$BASE_REF" 2>/dev/null); then
    echo "✓ Found merge base in initial fetch: $MERGE_BASE"
    echo "Debug: Commits from merge base to HEAD: $(git rev-list --count "$MERGE_BASE"..HEAD 2>/dev/null || echo 'unknown')"
    echo "Debug: Commits from merge base to origin/$BASE_REF: $(git rev-list --count "$MERGE_BASE".."origin/$BASE_REF" 2>/dev/null || echo 'unknown')"
    exit 0
else
    echo "Debug: merge-base command failed, checking available history..."
    echo "Debug: HEAD has $(git rev-list --count HEAD 2>/dev/null || echo 'unknown') commits"
    echo "Debug: origin/$BASE_REF has $(git rev-list --count "origin/$BASE_REF" 2>/dev/null || echo 'unknown') commits"
fi

# Incrementally deepen by FETCH_INCREMENT commits at a time
for i in $(seq 1 $MAX_ITERATIONS); do
    TOTAL_DEEPENED=$((i * FETCH_INCREMENT + FETCH_INCREMENT))
    echo "→ Merge base not found after initial $FETCH_INCREMENT + $((i * FETCH_INCREMENT)) = $TOTAL_DEEPENED commits, deepening..."
    echo "Debug: Before deepen - HEAD has $(git rev-list --count HEAD 2>/dev/null || echo 'unknown') commits"
    git fetch --deepen=$FETCH_INCREMENT
    echo "Debug: After deepen - HEAD has $(git rev-list --count HEAD 2>/dev/null || echo 'unknown') commits"
    echo "Debug: After deepen - origin/$BASE_REF has $(git rev-list --count "origin/$BASE_REF" 2>/dev/null || echo 'unknown') commits"
    
    if MERGE_BASE=$(git merge-base HEAD "origin/$BASE_REF" 2>/dev/null); then
        echo "✓ Found merge base after fetching ~$TOTAL_DEEPENED commits: $MERGE_BASE"
        echo "Debug: Commits from merge base to HEAD: $(git rev-list --count "$MERGE_BASE"..HEAD 2>/dev/null || echo 'unknown')"
        echo "Debug: Commits from merge base to origin/$BASE_REF: $(git rev-list --count "$MERGE_BASE".."origin/$BASE_REF" 2>/dev/null || echo 'unknown')"
        exit 0
    fi
done

# Final fallback to full fetch if needed
TOTAL_ATTEMPTED=$((MAX_ITERATIONS * FETCH_INCREMENT + FETCH_INCREMENT))
echo "⚠ Merge base still not found after $TOTAL_ATTEMPTED commits, fetching full history..."
git fetch --unshallow origin "$BASE_REF" || git fetch origin "$BASE_REF"

# Verify we found it
if git merge-base HEAD "origin/$BASE_REF" >/dev/null 2>&1; then
    echo "✓ Successfully found merge base with full history"
else
    echo "✗ Failed to find merge base with $BASE_REF even after full fetch"
    exit 1
fi
