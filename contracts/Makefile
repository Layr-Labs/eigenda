compile: deps
	forge build

# clean doesn't remove the bindings, because they are committed to the repo.
clean:
	forge clean

# make bindings compiles the contracts and creates go bindings
bindings: compile
	rm -rf ./bindings
	./generate-bindings.sh

fmt:
	forge fmt

fmt-check:
	forge fmt --check

deps:
	yarn