# Serves the mdbook docs located in this directory on
# http://localhost:3000 and will open a browser window
# to view the docs.
#
# TODO: Add markdown linting tests which enforce standard
#       consistency across protocol spec
serve: install-deps
	mdbook serve . --open

build: install-deps
	mdbook build

# TODO: Update mdbook to latest v0.5.1
#       https://github.com/rust-lang/mdBook/releases/tag/v0.5.1
install-deps:
	cargo install mdbook@0.4.52 mdbook-mermaid@0.14.1 mdbook-last-changed@0.1.4 mdbook-katex@0.9.4