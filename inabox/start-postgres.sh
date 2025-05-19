#!/bin/bash

# This script starts a PostgreSQL container for inabox deployment
# and initializes it with the necessary schemas for EigenDA

POSTGRES_PORT=${1:-5433}
POSTGRES_USER=${2:-postgres}
POSTGRES_PASSWORD=${3:-postgres}
POSTGRES_DB=${4:-eigenda}

echo "Starting PostgreSQL container on port $POSTGRES_PORT..."

# Check if container already exists
if docker ps -a | grep -q postgres-eigenda; then
  echo "PostgreSQL container already exists, stopping and removing it..."
  docker stop postgres-eigenda >/dev/null 2>&1
  docker rm postgres-eigenda >/dev/null 2>&1
fi

# Start PostgreSQL container
docker run -d \
  --name postgres-eigenda \
  -e POSTGRES_USER=$POSTGRES_USER \
  -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
  -e POSTGRES_DB=$POSTGRES_DB \
  -p $POSTGRES_PORT:5432 \
  postgres:14

# Wait for PostgreSQL to be ready
echo "Waiting for PostgreSQL to start..."
for i in {1..30}; do
  if docker exec postgres-eigenda pg_isready -U $POSTGRES_USER > /dev/null 2>&1; then
    echo "PostgreSQL is ready"
    break
  fi
  echo -n "."
  sleep 1
done

# Initialize database schemas
echo "Initializing PostgreSQL schemas..."

docker exec -i postgres-eigenda psql -U $POSTGRES_USER -d $POSTGRES_DB << EOF

-- Create tables for blob metadata
CREATE TABLE IF NOT EXISTS blob_metadata (
    blob_key BYTEA PRIMARY KEY,
    blob_header JSONB NOT NULL,
    signature BYTEA NOT NULL,
    requested_at BIGINT NOT NULL,
    requested_at_bucket BYTEA NOT NULL,
    requested_at_blob_key BYTEA NOT NULL,
    blob_status INTEGER NOT NULL,
    updated_at BIGINT NOT NULL,
    account_id VARCHAR(42) NOT NULL,
    expiry BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_blob_metadata_blob_status_updated_at ON blob_metadata (blob_status, updated_at);
CREATE INDEX IF NOT EXISTS idx_blob_metadata_account_id_requested_at ON blob_metadata (account_id, requested_at);
CREATE INDEX IF NOT EXISTS idx_blob_metadata_requested_at_bucket_key ON blob_metadata (requested_at_bucket, requested_at_blob_key);

-- Create tables for blob certificates
CREATE TABLE IF NOT EXISTS blob_certificates (
    blob_key BYTEA PRIMARY KEY,
    blob_certificate JSONB NOT NULL,
    fragment_info JSONB NOT NULL
);

-- Create tables for batch headers
CREATE TABLE IF NOT EXISTS batch_headers (
    batch_header_hash BYTEA PRIMARY KEY,
    batch_header JSONB NOT NULL
);

-- Create tables for batches
CREATE TABLE IF NOT EXISTS batches (
    batch_header_hash BYTEA PRIMARY KEY,
    batch_info JSONB NOT NULL
);

-- Create tables for dispersal requests
CREATE TABLE IF NOT EXISTS dispersal_requests (
    batch_header_hash BYTEA NOT NULL,
    operator_id BYTEA NOT NULL,
    dispersal_request JSONB NOT NULL,
    dispersed_at BIGINT NOT NULL,
    PRIMARY KEY (batch_header_hash, operator_id)
);

CREATE INDEX IF NOT EXISTS idx_dispersal_requests_operator_dispersed_at ON dispersal_requests (operator_id, dispersed_at);

-- Create tables for dispersal responses
CREATE TABLE IF NOT EXISTS dispersal_responses (
    batch_header_hash BYTEA NOT NULL,
    operator_id BYTEA NOT NULL,
    dispersal_response JSONB NOT NULL,
    responded_at BIGINT NOT NULL,
    PRIMARY KEY (batch_header_hash, operator_id)
);

CREATE INDEX IF NOT EXISTS idx_dispersal_responses_operator_responded_at ON dispersal_responses (operator_id, responded_at);

-- Create tables for attestations
CREATE TABLE IF NOT EXISTS attestations (
    batch_header_hash BYTEA PRIMARY KEY,
    attestation JSONB NOT NULL,
    attested_at BIGINT NOT NULL,
    attested_at_bucket VARCHAR(64) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_attestations_attested_at_bucket_attested_at ON attestations (attested_at_bucket, attested_at);

-- Create tables for blob inclusion info
CREATE TABLE IF NOT EXISTS blob_inclusion_info (
    blob_key BYTEA NOT NULL,
    batch_header_hash BYTEA NOT NULL,
    inclusion_info JSONB NOT NULL,
    PRIMARY KEY (blob_key, batch_header_hash)
);

EOF

# Check if initialization was successful
if [ $? -eq 0 ]; then
  echo "PostgreSQL initialization completed successfully!"
else
  echo "Error: PostgreSQL initialization failed!"
  exit 1
fi

echo "PostgreSQL is now ready for use with the following connection parameters:"
echo "  Host: localhost"
echo "  Port: $POSTGRES_PORT"
echo "  User: $POSTGRES_USER"
echo "  Password: $POSTGRES_PASSWORD"
echo "  Database: $POSTGRES_DB"