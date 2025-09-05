# Serves the mdbook docs located in this directory on
# http://localhost:3000 and will open a browser window
# to view the docs.
serve: install-deps
	mdbook serve . --open

build: install-deps
	mdbook build

install-deps:
	cargo install mdbook mdbook-mermaid@0.14.1
	cargo install mdbook mdbook-last-changed@0.1.4