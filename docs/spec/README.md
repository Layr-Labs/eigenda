# EigenDA Spec

Built using [mdBook](https://rust-lang.github.io/mdBook/index.html) and published as a github pages site at [https://layr-labs.github.io/eigenda/](https://layr-labs.github.io/eigenda/).

Meant to contain technical overviews, spec, and low-level implementation details related to EigenDA, as opposed to the [docs](https://docs.eigenda.xyz/) site which is meant to contain more introductory and high-level material.

## Preview

To preview the book locally, run:

```bash
make serve
```

which will start a local server at `http://localhost:3000` and open your browser to preview the result.

## Github Pages

The book is automatically built and deployed to Github Pages on every push to the `main` branch.
This is done by the Github Actions workflow defined in [../../.github/workflows/mdbook.yaml](../../.github/workflows/mdbook.yaml)

## Mermaid Diagrams

We use mdbook-mermaid to render mermaid diagrams in the book. It is installed along with mdbook when running `make install-deps`. The 2 js files `mermaid-init.js` and `mermaid.min.js` were installed from `mdbook-mermaid install .` which was ran with mdbook-mermaid v0.14.1. These two files are copied into the built book and needed to render the images. Haven't found a way to only generate the images and then update the markdown files to reference the images, so we are stuck with this dependency for now.