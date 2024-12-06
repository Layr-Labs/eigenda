#!/usr/bin/env bash

# This script generates a new ECDSA keypair for the disperser service.

openssl ecparam -name prime256v1 -genkey -noout -out eigenda-disperser-private.pem
openssl ec -in eigenda-disperser-private.pem -pubout -out eigenda-disperser-public.pem
