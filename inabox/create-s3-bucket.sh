#!/bin/bash
set -e

S3_BUCKET="test-eigenda-blobstore"
S3_REGION="us-east-1"

if AWS_ACCESS_KEY_ID=localstack AWS_SECRET_ACCESS_KEY=localstack \
    aws s3api head-bucket --endpoint-url=$AWS_URL --bucket "$S3_BUCKET" 2>/dev/null; then
    echo "Bucket $S3_BUCKET already exists"
else
    echo "Creating bucket $S3_BUCKET"
   AWS_ACCESS_KEY_ID=localstack AWS_SECRET_ACCESS_KEY=localstack aws s3api create-bucket \
            --endpoint-url=$AWS_URL \
            --bucket "$S3_BUCKET" \
            --region "$S3_REGION" 
fi
