.PHONY: all build package clean

# Configuration
IMAGE_NAME=lambda-packager
CONTAINER_NAME=lambda-packager-container
LAMBDA_ZIP=function.zip

all: package

# Build Docker image
build:
	@echo "Building Docker image for Lambda packaging..."
	docker build -t $(IMAGE_NAME) .

# Create Lambda package using Docker
package: build
	@echo "Creating Lambda package using Docker..."
	docker run --name $(CONTAINER_NAME) $(IMAGE_NAME)
	docker cp $(CONTAINER_NAME):/tmp/function.zip ./$(LAMBDA_ZIP)
	docker rm $(CONTAINER_NAME)
	@echo "Lambda package created: $(LAMBDA_ZIP)"

# Clean up resources
clean:
	@echo "Cleaning up resources..."
	docker rmi $(IMAGE_NAME) || true
	rm -f $(LAMBDA_ZIP)

