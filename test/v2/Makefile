clean:
	rm -rf bin 2>/dev/null || true

build: clean
	go build -o bin/load load/main/load_main.go

test:
	cd correctness && go test

generate-load: build
	./bin/load $(ARGS)
