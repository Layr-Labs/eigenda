#!/usr/bin/env bash

# This script finds and deletes all compiled protobufs.

# The location where this script can be found.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

API_DIR="${SCRIPT_DIR}/.."
GRPC_DIR="${API_DIR}/grpc"
find "${GRPC_DIR}" -name '*.pb.go' -type f | xargs rm -rf

DISPERSER_DIR="$SCRIPT_DIR/../../disperser"
DISPERSER_GRPC_DIR="$DISPERSER_DIR/api/grpc"
find "${DISPERSER_GRPC_DIR}" -name '*.pb.go' -type f | xargs rm -rf
