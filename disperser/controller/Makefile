build:
	cd .. && go build -o ./bin/controller ./cmd/controller

test:
	go test -v ./...

run: build
	@echo "  >  Running controller..."
	cd .. && ./bin/controller
