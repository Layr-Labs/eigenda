# Serves the mdbook docs located in this directory.
# Will open a browser window to view the docs.
serve: install-deps
	mdbook serve . --open

build: install-deps
	mdbook build

install-deps:
	cargo install mdbook mdbook-mermaid@0.14.1