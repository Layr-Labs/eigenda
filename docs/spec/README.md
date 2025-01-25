# EigenDA Spec

Built using [mdBook](https://rust-lang.github.io/mdBook/index.html).

Meant to contain technical overviews, spec, and low-level implementation details related to EigenDA, as opposed to the [docs](https://docs.eigenda.xyz/) site which is meant to contain more introductory and high-level material.

## Images

`./src/assets` is a symlink to `docs/assets` in order for MdBook to be able to access the images and other assets in the `docs` directory, given that it seems to only have access to the `src` directory.

Still not quite sure what's the best directory structure. Perhaps we could simply move the assets directory entirely under book/src, but then other READMEs throughout the repo that need access to those assets would need to be updated.