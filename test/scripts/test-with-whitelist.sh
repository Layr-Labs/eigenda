#!/usr/bin/env bash

# Runs tests only for explicitly whitelisted packages (or directories).
# Usage: ./test-with-whitelist.sh <root> <whitelisted packages or dirs...>
set -euo pipefail

if [ "$#" -lt 2 ]; then
  echo "Usage: $0 <root> <whitelisted packages or dirs...>" >&2
  exit 1
fi

ROOT=$1
shift

PKGS=""

for WL in "$@"; do
  # Normalize bare names relative to root; keep paths/imports/... as-is
  case "$WL" in
    .|./*|/*|*...*) cand="$WL" ;;   # already a path or has ... pattern
    *)              cand="$ROOT/$WL" ;;
  esac

  # If it's a directory and not already using ..., include subpackages
  if [ -d "$cand" ] && ! printf %s "$cand" | grep -q '\.\.\.'; then
    cand="$cand/..."
  fi

  if out=$(go list "$cand" 2>/dev/null); then
    PKGS="$PKGS $out"
  else
    echo "Warning: '$WL' resolved to '$cand' but matched no packages" >&2
  fi
done

# Trim leading/trailing spaces
PKGS=$(echo "$PKGS" | xargs)

if [ -z "$PKGS" ]; then
  echo "No packages matched the whitelist." >&2
  exit 0
fi

echo "Running tests for whitelist:"
printf '%s\n' $PKGS

CI=true go test -short $PKGS -coverprofile=coverage.out
