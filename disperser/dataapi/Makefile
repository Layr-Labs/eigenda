build:
	cd .. && go build -o ./bin/dataapi ./cmd/dataapi

test:
	go test -v ./...

generate-swagger:
	@echo "  >  Generating swagger..."
	swag init -g swagger.go --parseDependency --output docs/v2 --instanceName V2 --packageName v2 --dir ./v2 --parseDepth 0
	swag fmt --dir ./v2

run: build
	@echo "  >  Running dataapi..."
	cd .. && ./bin/dataapi