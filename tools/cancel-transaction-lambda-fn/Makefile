.PHONY: start clean test package-lambda create-key deploy-lambda invoke-lambda all

# Configuration variables
AWS_ENDPOINT=http://localhost:4566
LAMBDA_NAME=eth-tx-canceller
LAMBDA_ROLE=eth-tx-role
LAMBDA_ZIP=function.zip

all: start package-lambda create-key deploy-lambda invoke-lambda

# Start the Docker Compose environment
start:
	@echo "Starting Docker Compose environment..."
	docker-compose up -d
	@echo "Waiting for services to initialize..."
	sleep 5

# Clean up resources
clean:
	@echo "Cleaning up resources..."
	docker-compose down
	rm -f $(LAMBDA_ZIP) lambda_output.json .kms_key_id

# Create Lambda deployment package
package-lambda:
	@echo "Creating Lambda deployment package..."
	mkdir -p lambda_package
	cp lambda_function.py lambda_package/
	cp lambda_helper.py lambda_package/
	pip install --implementation cp --only-binary=:all: -r requirements.txt -t lambda_package/
	cd lambda_package && zip -r ../$(LAMBDA_ZIP) *
	rm -rf lambda_package
	@echo "Lambda package created: $(LAMBDA_ZIP)"

# Create KMS key for signing
create-key:
	@echo "Creating KMS key..."
	aws --endpoint-url=$(AWS_ENDPOINT) kms create-key \
		--description "Ethereum signer key" \
		--key-usage SIGN_VERIFY \
		--customer-master-key-spec ECC_SECG_P256K1 \
		--query 'KeyMetadata.KeyId' --output text > .kms_key_id
	@echo "Created KMS key: $$(cat .kms_key_id)"

# Create IAM role and deploy Lambda
deploy-lambda:
	@echo "Creating IAM role..."
	aws --endpoint-url=$(AWS_ENDPOINT) iam create-role \
		--role-name $(LAMBDA_ROLE) \
		--assume-role-policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"Service":"lambda.amazonaws.com"},"Action":"sts:AssumeRole"}]}'
	
	@echo "Adding KMS permissions to role..."
	aws --endpoint-url=$(AWS_ENDPOINT) iam put-role-policy \
		--role-name $(LAMBDA_ROLE) \
		--policy-name lambda-kms-policy \
		--policy-document '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Action":["kms:GetPublicKey","kms:Sign"],"Resource":"*"}]}'
	
	@echo "Deploying Lambda function..."
	aws --endpoint-url=$(AWS_ENDPOINT) lambda create-function \
		--function-name $(LAMBDA_NAME) \
		--runtime python3.10 \
		--handler lambda_function.lambda_handler \
		--timeout 30 \
		--memory-size 512 \
		--role arn:aws:iam::000000000000:role/$(LAMBDA_ROLE) \
		--zip-file fileb://$(LAMBDA_ZIP) \
		--environment "Variables={ETH_RPC=http://anvil:8545,PRIORITY_FEE_GWEI=1,NONCE=0,KMS_KEY_ID=$$(cat .kms_key_id),LOG_LEVEL=DEBUG}"
	
	@echo "Lambda function deployed"
	aws --endpoint-url=$(AWS_ENDPOINT) lambda get-function --function-name $(LAMBDA_NAME)
	
	@echo "Waiting for Lambda function to initialize..."
	sleep 5

# Invoke the Lambda function
invoke-lambda:
	@echo "Invoking Lambda function..."
	aws --endpoint-url=$(AWS_ENDPOINT) lambda invoke \
		--function-name $(LAMBDA_NAME) \
		--payload '{}' \
		--cli-binary-format raw-in-base64-out \
		--log-type Tail \
		lambda_output.json | grep -v "Latency\|END\|START\|REPORT" | jq -r '.LogResult' | base64 -d || true
	
	@echo "Lambda execution result:"
	@if [ -f lambda_output.json ]; then cat lambda_output.json; else echo "No output file created"; fi

# Simple full test
test: clean start package-lambda create-key deploy-lambda invoke-lambda

