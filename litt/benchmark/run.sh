#!/usr/bin/env bash

# This script is used to run the LittDB benchmark.

# Find the directory of this script
SCRIPT_DIR=$(dirname "$(readlink -f "$0")")

go run "$SCRIPT_DIR/cmd/main.go" "$@"
