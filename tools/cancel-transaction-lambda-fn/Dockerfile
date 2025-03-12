FROM python:3.12-slim

WORKDIR /app

# Install zip utility
RUN apt-get update && apt-get install -y zip && rm -rf /var/lib/apt/lists/*

# Copy source files
COPY lambda_function.py .
COPY lambda_helper.py .
COPY requirements.txt .

# Install dependencies with Lambda compatibility
RUN pip install --platform manylinux2014_x86_64 \
    --implementation cp \
    --python-version 3.12 \
    --only-binary=:all: \
    --target=. \
    -r requirements.txt

# Create the Lambda package
RUN zip -r /tmp/function.zip .

# Default command does nothing, we'll use docker cp to get the zip
CMD ["echo", "Lambda package created in /tmp/function.zip"]
