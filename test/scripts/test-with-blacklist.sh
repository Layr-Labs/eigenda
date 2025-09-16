#!/usr/bin/env bash

# Runs all tests under the specified root, excluding any blacklisted packages/dirs.
# Usage: ./test-with-blacklist.sh <root> [blacklisted packages or dirs...]
set -euo pipefail

if [ "$#" -lt 1 ]; then
  echo "Usage: $0 <root> [blacklisted packages or dirs...]" >&2
  exit 1
fi

ROOT=$1
shift

# Resolve blacklist entries to concrete import paths.
EXCLUDED=""
for BL in "$@"; do
  # Normalize bare names relative to root; keep paths/imports/... as-is
  case "$BL" in
    .|./*|/*|*...*) cand="$BL" ;;     # already a path or has ... pattern
    *)              cand="$ROOT/$BL" ;;
  esac

  # If it's a directory and not already using ..., include subpackages
  if [ -d "$cand" ] && ! printf %s "$cand" | grep -q '\.\.\.'; then
    cand="$cand/..."
  fi

  if out=$(go list "$cand" 2>/dev/null); then
    EXCLUDED="$EXCLUDED
$out"
  else
    echo "Warning: '$BL' resolved to '$cand' but matched no packages" >&2
  fi
done

# All packages under root.
ALL=$(go list "$ROOT"/...)

# Filter out excluded packages (exact match).
PKGS=""
for p in $ALL; do
  if printf '%s\n' "$EXCLUDED" | grep -Fxq "$p"; then
    continue
  fi
  PKGS="$PKGS $p"
done

# Trim whitespace.
PKGS=$(echo "$PKGS" | xargs || true)

if [ -z "$PKGS" ]; then
  echo "No packages left to test after applying blacklist." >&2
  exit 0
fi

echo "Running tests for:"
printf '%s\n' $PKGS

# Run tests (coverage output like your prior script)
go test -short $PKGS -coverprofile=coverage.out
