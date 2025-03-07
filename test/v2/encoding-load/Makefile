BINARY_NAME=encoding-load-generator
# Define config file parameter
CONFIG_FILE?=

.PHONY: build
build:
	go build -o $(BINARY_NAME) ./main

.PHONY: clean
clean:
	rm -f $(BINARY_NAME)

.PHONY: run
run: build
	@if [ -z "$(CONFIG_FILE)" ]; then \
		echo "Error: CONFIG_FILE must be provided"; \
		echo "Usage: make run CONFIG_FILE=<config_file_path>"; \
		echo "       or use CONFIG_FILE=- to read from GENERATOR_ENCODING_LOAD environment variable"; \
		exit 1; \
	fi
	./$(BINARY_NAME) $(CONFIG_FILE) 