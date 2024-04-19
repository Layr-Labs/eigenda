build: 
	cd .. && go build -o ./bin/dataapi ./cmd/dataapi

test:
	go test -v .

generate-swagger:
	@echo "  >  Generating swagger..."
	swag init -g ../cmd/dataapi/main.go --parseDependency
	swag fmt
