#!/usr/bin/env bash

# This script is used to run the LittDB benchmark.

# Find the directory of this script
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

(
    cd "$SCRIPT_DIR/.." || exit 1
    make build || exit 1

    BIN_PATH=$(cd "$(dirname './bin')" && pwd)/$(basename './bin')
    BINARY_PATH="${BIN_PATH}/benchmark"
)

$BINARY_PATH "$@"
