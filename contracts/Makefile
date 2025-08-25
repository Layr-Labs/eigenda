deps:
	yarn

compile: deps
	forge build

# make bindings compiles the contracts and creates go bindings
bindings: compile
	rm -rf ./bindings
	./generate-bindings.sh

fmt:
	forge fmt

fmt-check:
	forge fmt --check