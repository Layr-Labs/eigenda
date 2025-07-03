#!/bin/sh

# Wait 2 seconds to ensure minio is finished bootstrapping
# TODO: Update this to do event based polling on minio server directly vs semi-arbitrary timeout
sleep 2s

# Configure MinIO client (mc)
echo "Configuring MinIO client..."
mc alias set local http://minio:9000 minioadmin minioadmin

# Ensure the bucket exists
echo "Creating bucket: eigenda-proxy-test..."
mc mb local/eigenda-proxy-test || echo "Bucket already exists."

echo "Bucket setup complete."
