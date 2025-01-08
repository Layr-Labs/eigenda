build:
	cd .. && go build -o ./bin/dataapi ./cmd/dataapi

test:
	go test -v ./...

generate-swagger-v1:
	@echo "  >  Generating v1 swagger..."
	swag init -g ../cmd/dataapi/main.go --parseDependency --output docs/v1 --instanceName V1 --packageName v1 --parseDepth 0 --exclude ./v2 --dir .
	swag fmt --dir . --exclude ./v2/server_v2.go

generate-swagger-v2:
	@echo "  >  Generating v2 swagger..."
	swag init -g swagger.go --parseDependency --output docs/v2 --instanceName V2 --packageName v2 --dir ./v2 --parseDepth 0
	swag fmt --dir ./v2

generate-swagger: generate-swagger-v1 generate-swagger-v2
