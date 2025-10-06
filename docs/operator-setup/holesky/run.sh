#!/bin/sh
# Path: run.sh

. ./.env

node_plugin_image="ghcr.io/layr-labs/eigenda/opr-nodeplugin:2.2.0"

# Check if V2 ports are defined
if [ -z "$NODE_V2_DISPERSAL_PORT" ]; then
  echo "ERROR: NODE_V2_DISPERSAL_PORT is not defined!"
fi
if [ -z "$NODE_V2_RETRIEVAL_PORT" ]; then
  echo "ERROR: NODE_V2_RETRIEVAL_PORT is not defined!"
fi
if [ -z "$NODE_V2_DISPERSAL_PORT" ] || [ -z "$NODE_V2_RETRIEVAL_PORT" ]; then
  echo "ERROR: Please update your .env file. See .env.example for reference."
    exit 1
fi
socket="$NODE_HOSTNAME":"${NODE_DISPERSAL_PORT}"\;"${NODE_RETRIEVAL_PORT}"\;"${NODE_V2_DISPERSAL_PORT}"\;"${NODE_V2_RETRIEVAL_PORT}"

# In all commands, We have to explicitly set the password again here because
# when docker run loads the `.env` file, it keeps the quotes around the password
# which causes the password to be incorrect.
# To test that try running `docker run --rm --env-file .env busybox /bin/sh -c 'echo $NODE_ECDSA_KEY_PASSWORD'`
# This will output password with single quote. Not sure why this happens.
optIn() {
  echo "You are about to opt-in to quorum: $1 with socket registration: $socket"
  echo "Confirm? [Y/n] "
  read -r answer
  if [ "$answer" = "n" ] || [ "$answer" = "N" ]; then
    echo "Operation cancelled"
    exit 1
  fi

  docker run --env-file .env \
  --rm \
  --volume "${NODE_ECDSA_KEY_FILE_HOST}":/app/operator_keys/ecdsa_key.json \
  --volume "${NODE_BLS_KEY_FILE_HOST}":/app/operator_keys/bls_key.json \
  --volume "${NODE_LOG_PATH_HOST}":/app/logs:rw \
  "$node_plugin_image" \
  --ecdsa-key-password "$NODE_ECDSA_KEY_PASSWORD" \
  --bls-key-password "$NODE_BLS_KEY_PASSWORD" \
  --operation opt-in \
  --socket "$socket" \
  --quorum-id-list "$1"
}

optOut() {
  echo "You are about to opt-out from quorum: $1 with socket registration: $socket"
  echo "Confirm? [Y/n] "
  read -r answer
  if [ "$answer" = "n" ] || [ "$answer" = "N" ]; then
    echo "Operation cancelled"
    exit 1
  fi

  docker run --env-file .env \
    --rm \
    --volume "${NODE_ECDSA_KEY_FILE_HOST}":/app/operator_keys/ecdsa_key.json \
    --volume "${NODE_BLS_KEY_FILE_HOST}":/app/operator_keys/bls_key.json \
    --volume "${NODE_LOG_PATH_HOST}":/app/logs:rw \
    "$node_plugin_image" \
    --ecdsa-key-password "$NODE_ECDSA_KEY_PASSWORD" \
    --bls-key-password "$NODE_BLS_KEY_PASSWORD" \
    --operation opt-out \
    --socket "$socket" \
    --quorum-id-list "$1"
}

listQuorums() {
  # we have to pass a dummy quorum-id-list as it is required by the plugin
  docker run --env-file .env \
    --rm \
    --volume "${NODE_ECDSA_KEY_FILE_HOST}":/app/operator_keys/ecdsa_key.json \
    --volume "${NODE_BLS_KEY_FILE_HOST}":/app/operator_keys/bls_key.json \
    --volume "${NODE_LOG_PATH_HOST}":/app/logs:rw \
    "$node_plugin_image" \
    --ecdsa-key-password "$NODE_ECDSA_KEY_PASSWORD" \
    --bls-key-password "$NODE_BLS_KEY_PASSWORD" \
    --socket "$socket" \
    --operation list-quorums \
    --quorum-id-list 0
}

updateSocket() {
  echo "You are about to update your socket registration to: $socket"
  echo "Confirm? [Y/n] "
  read -r answer
  if [ "$answer" = "n" ] || [ "$answer" = "N" ]; then
    echo "Operation cancelled"
    exit 1
  fi

  # we have to pass a dummy quorum-id-list as it is required by the plugin
  docker run --env-file .env \
    --rm \
    --volume "${NODE_ECDSA_KEY_FILE_HOST}":/app/operator_keys/ecdsa_key.json \
    --volume "${NODE_BLS_KEY_FILE_HOST}":/app/operator_keys/bls_key.json \
    --volume "${NODE_LOG_PATH_HOST}":/app/logs:rw \
    "$node_plugin_image" \
    --ecdsa-key-password "$NODE_ECDSA_KEY_PASSWORD" \
    --bls-key-password "$NODE_BLS_KEY_PASSWORD" \
    --socket "$socket" \
    --operation update-socket \
    --quorum-id-list 0
}

if [ "$1" = "opt-in" ]; then
  if [ -z "$2" ]; then
    echo "Please provide quorum number (0/1/0,1)"
    echo "Example Usage: ./run.sh opt-in 0"
    exit 1
  fi
  optIn "$2"
elif [ "$1" = "opt-out" ]; then
  if [ -z "$2" ]; then
    echo "Please provide quorum number (0/1/0,1)"
    echo "Example Usage: ./run.sh opt-out 0"
    exit 1
  fi
  optOut "$2"
elif [ "$1" = "list-quorums" ]; then
  listQuorums
elif [ "$1" = "update-socket" ]; then
  updateSocket
else
  echo "Invalid command"
fi
