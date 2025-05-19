.PHONY: test coverage coverage-html clean

# Default target
all: test

# Run all tests
test:
	go test -v ./...

# Run tests with coverage
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

# Generate HTML coverage report
coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Clean up test artifacts
clean:
	rm -f coverage.out coverage.html