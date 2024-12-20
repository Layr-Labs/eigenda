build:
	cd .. && go build -o ./bin/dataapi ./cmd/dataapi

test:
	go test -v .

generate-swagger-v1:
	@echo "  >  Generating v1 swagger..."
	swag init -g swagger.go --parseDependency --parseInternal --output docs/v1 --instanceName V1 --packageName v1 --dir ./v1 --parseDepth 1
	swag fmt --dir ./v1

generate-swagger-v2:
	@echo "  >  Generating v2 swagger..."
	swag init -g swagger.go --parseDependency --parseInternal --output docs/v2 --instanceName V2 --packageName v2 --dir ./v2 --parseDepth 1
	swag fmt --dir ./v2

generate-swagger: generate-swagger-v1 generate-swagger-v2
