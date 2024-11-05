#!/usr/bin/env bash

# This script finds and deletes all compiled protobufs.

# The location where this script can be found.
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

API_DIR="${SCRIPT_DIR}/.."
GRPC_DIR="${API_DIR}/grpc"

if [ -d "${GRPC_DIR}" ]; then
  # Delete all compiled protobufs
  find "${GRPC_DIR}" -name '*.pb.go' -type f | xargs rm -rf
  # Delete all empty directories
  find "${GRPC_DIR}" -type d -empty -delete
fi

DISPERSER_DIR="$SCRIPT_DIR/../../disperser"
DISPERSER_GRPC_DIR="$DISPERSER_DIR/api/grpc"
if [ -d "${DISPERSER_GRPC_DIR}" ]; then
  # Delete all compiled protobufs
  find "${DISPERSER_GRPC_DIR}" -name '*.pb.go' -type f | xargs rm -rf
  # Delete all empty directories
  find "${DISPERSER_GRPC_DIR}" -type d -empty -delete
fi